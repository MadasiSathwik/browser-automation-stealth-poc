package search

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/automation-poc/browser-automation/config"
	"github.com/automation-poc/browser-automation/logger"
	"github.com/automation-poc/browser-automation/stealth"
	"github.com/automation-poc/browser-automation/storage"
	"github.com/go-rod/rod"
)

type Profile struct {
	ID      string
	Name    string
	Title   string
	Company string
	URL     string
}

type Service struct {
	cfg        *config.Config
	db         *storage.Database
	log        *logger.Logger
	fingerprint *stealth.FingerprintManager
	timing     *stealth.TimingController
}

func NewService(cfg *config.Config, db *storage.Database, log *logger.Logger) *Service {
	return &Service{
		cfg:        cfg,
		db:         db,
		log:        log,
		fingerprint: stealth.NewFingerprintManager(&cfg.Stealth),
		timing:     stealth.NewTimingController(&cfg.Timing),
	}
}

func (s *Service) SearchProfiles(ctx context.Context, page *rod.Page, searchCfg config.SearchConfig) ([]*Profile, error) {
	s.log.Infof("Starting profile search: %s", searchCfg.Query)

	searchURL := s.buildSearchURL(searchCfg)
	s.log.Debugf("Navigating to: %s", searchURL)

	if err := page.Navigate(searchURL); err != nil {
		return nil, fmt.Errorf("failed to navigate to search page: %w", err)
	}

	if err := page.WaitLoad(); err != nil {
		return nil, fmt.Errorf("failed to wait for search page load: %w", err)
	}

	time.Sleep(s.timing.PageLoadDelay())

	if err := s.fingerprint.RandomScroll(page); err != nil {
		s.log.Warnf("Failed to perform random scroll: %v", err)
	}

	var allProfiles []*Profile

	for pageNum := 1; pageNum <= searchCfg.MaxPages; pageNum++ {
		select {
		case <-ctx.Done():
			return allProfiles, ctx.Err()
		default:
		}

		s.log.Infof("Processing search results page %d", pageNum)

		profiles, err := s.extractProfilesFromPage(page, searchCfg.Selectors)
		if err != nil {
			s.log.Errorf("Failed to extract profiles from page %d: %v", pageNum, err)
			break
		}

		if len(profiles) == 0 {
			s.log.Info("No more profiles found")
			break
		}

		for _, profile := range profiles {
			if !s.db.HasProcessedProfile(profile.ID) {
				allProfiles = append(allProfiles, profile)
			} else {
				s.log.Debugf("Skipping already processed profile: %s", profile.Name)
			}
		}

		if pageNum < searchCfg.MaxPages {
			if err := s.goToNextPage(page, searchCfg.Selectors); err != nil {
				s.log.Warnf("Failed to navigate to next page: %v", err)
				break
			}

			time.Sleep(s.timing.BetweenActionsDelay())
		}
	}

	s.log.Infof("Search completed: found %d profiles", len(allProfiles))
	return allProfiles, nil
}

func (s *Service) buildSearchURL(searchCfg config.SearchConfig) string {
	baseURL := searchCfg.BaseURL

	if searchCfg.Query != "" {
		separator := "?"
		if strings.Contains(baseURL, "?") {
			separator = "&"
		}
		baseURL = fmt.Sprintf("%s%sq=%s", baseURL, separator, searchCfg.Query)
	}

	for key, value := range searchCfg.Filters {
		baseURL = fmt.Sprintf("%s&%s=%s", baseURL, key, value)
	}

	return baseURL
}

func (s *Service) extractProfilesFromPage(page *rod.Page, selectors map[string]string) ([]*Profile, error) {
	profileCardSelector, ok := selectors["profile_card"]
	if !ok {
		return nil, fmt.Errorf("profile_card selector not configured")
	}

	elements, err := page.Elements(profileCardSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to find profile cards: %w", err)
	}

	var profiles []*Profile

	for _, elem := range elements {
		profile, err := s.extractProfileFromCard(elem, selectors)
		if err != nil {
			s.log.Warnf("Failed to extract profile: %v", err)
			continue
		}

		if profile != nil {
			profiles = append(profiles, profile)
		}
	}

	return profiles, nil
}

func (s *Service) extractProfileFromCard(card *rod.Element, selectors map[string]string) (*Profile, error) {
	profile := &Profile{}

	if nameSelector, ok := selectors["profile_name"]; ok {
		if nameElem, err := card.Element(nameSelector); err == nil {
			if name, err := nameElem.Text(); err == nil {
				profile.Name = strings.TrimSpace(name)
			}
		}
	}

	if titleSelector, ok := selectors["profile_title"]; ok {
		if titleElem, err := card.Element(titleSelector); err == nil {
			if title, err := titleElem.Text(); err == nil {
				profile.Title = strings.TrimSpace(title)
			}
		}
	}

	if companySelector, ok := selectors["profile_company"]; ok {
		if companyElem, err := card.Element(companySelector); err == nil {
			if company, err := companyElem.Text(); err == nil {
				profile.Company = strings.TrimSpace(company)
			}
		}
	}

	if linkSelector, ok := selectors["profile_link"]; ok {
		if linkElem, err := card.Element(linkSelector); err == nil {
			if href, err := linkElem.Attribute("href"); err == nil && href != nil {
				profile.URL = *href
			}
		}
	}

	if profile.Name == "" {
		return nil, fmt.Errorf("profile name is empty")
	}

	profile.ID = s.generateProfileID(profile)

	return profile, nil
}

func (s *Service) generateProfileID(profile *Profile) string {
	if profile.URL != "" {
		parts := strings.Split(profile.URL, "/")
		for _, part := range parts {
			if part != "" && !strings.HasPrefix(part, "http") {
				return part
			}
		}
	}

	normalized := strings.ToLower(profile.Name)
	normalized = strings.ReplaceAll(normalized, " ", "-")
	return normalized
}

func (s *Service) goToNextPage(page *rod.Page, selectors map[string]string) error {
	nextButtonSelector, ok := selectors["next_page_button"]
	if !ok {
		return fmt.Errorf("next_page_button selector not configured")
	}

	has, nextButton, err := page.Has(nextButtonSelector)
	if err != nil {
		return fmt.Errorf("failed to check for next button: %w", err)
	}

	if !has {
		return fmt.Errorf("next button not found")
	}

	isDisabled, err := nextButton.Attribute("disabled")
	if err == nil && isDisabled != nil {
		return fmt.Errorf("next button is disabled")
	}

	if err := s.fingerprint.RandomScroll(page); err != nil {
		s.log.Warnf("Failed to scroll before pagination: %v", err)
	}

	if err := nextButton.Click(); err != nil {
		return fmt.Errorf("failed to click next button: %w", err)
	}

	if err := page.WaitLoad(); err != nil {
		return fmt.Errorf("failed to wait for next page load: %w", err)
	}

	time.Sleep(s.timing.PageLoadDelay())

	return nil
}
