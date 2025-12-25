package connections

import (
	"context"
	"fmt"
	"time"

	"github.com/automation-poc/browser-automation/config"
	"github.com/automation-poc/browser-automation/logger"
	"github.com/automation-poc/browser-automation/search"
	"github.com/automation-poc/browser-automation/stealth"
	"github.com/automation-poc/browser-automation/storage"
	"github.com/go-rod/rod"
)

type Service struct {
	cfg         *config.Config
	db          *storage.Database
	log         *logger.Logger
	mouse       *stealth.MouseController
	typing      *stealth.TypingSimulator
	timing      *stealth.TimingController
	fingerprint *stealth.FingerprintManager
	limiter     *RateLimiter
}

func NewService(cfg *config.Config, db *storage.Database, log *logger.Logger) *Service {
	return &Service{
		cfg:         cfg,
		db:          db,
		log:         log,
		mouse:       stealth.NewMouseController(&cfg.Stealth.MouseMovement),
		typing:      stealth.NewTypingSimulator(&cfg.Timing),
		timing:      stealth.NewTimingController(&cfg.Timing),
		fingerprint: stealth.NewFingerprintManager(&cfg.Stealth),
		limiter:     NewRateLimiter(cfg, db, log),
	}
}

func (s *Service) SendConnectionRequest(ctx context.Context, page *rod.Page, profile *search.Profile) error {
	if s.cfg.DemoMode {
		s.log.Infof("[DEMO] Would send connection request to %s", profile.Name)
		return nil
	}

	if !s.limiter.CanSendConnection() {
		return fmt.Errorf("rate limit reached")
	}

	s.log.Infof("Sending connection request to %s (%s at %s)", profile.Name, profile.Title, profile.Company)

	if err := page.Navigate(profile.URL); err != nil {
		return fmt.Errorf("failed to navigate to profile: %w", err)
	}

	if err := page.WaitLoad(); err != nil {
		return fmt.Errorf("failed to wait for profile page load: %w", err)
	}

	time.Sleep(s.timing.PageLoadDelay())

	if err := s.fingerprint.SimulateReading(page, 3*time.Second); err != nil {
		s.log.Warnf("Failed to simulate reading: %v", err)
	}

	connectButton, err := s.findConnectButton(page)
	if err != nil {
		return fmt.Errorf("failed to find connect button: %w", err)
	}

	if err := s.mouse.HoverElement(page, connectButton); err != nil {
		s.log.Warnf("Failed to hover connect button: %v", err)
	}

	time.Sleep(s.timing.ThinkDelay())

	if err := s.mouse.ClickElement(page, connectButton); err != nil {
		return fmt.Errorf("failed to click connect button: %w", err)
	}

	time.Sleep(s.timing.HumanizedDelay("click"))

	message := s.generatePersonalizedMessage(profile)

	if err := s.addConnectionNote(page, message); err != nil {
		s.log.Warnf("Failed to add connection note: %v", err)
	}

	if err := s.submitConnectionRequest(page); err != nil {
		return fmt.Errorf("failed to submit connection request: %w", err)
	}

	connectionReq := &storage.ConnectionRequest{
		ProfileID: profile.ID,
		Name:      profile.Name,
		Title:     profile.Title,
		Company:   profile.Company,
		Message:   message,
		Status:    "pending",
		SentAt:    time.Now(),
	}

	if err := s.db.SaveConnectionRequest(connectionReq); err != nil {
		return fmt.Errorf("failed to save connection request: %w", err)
	}

	if err := s.db.IncrementConnectionCount(); err != nil {
		s.log.Warnf("Failed to increment connection count: %v", err)
	}

	s.log.Infof("Connection request sent successfully to %s", profile.Name)

	return nil
}

func (s *Service) findConnectButton(page *rod.Page) (*rod.Element, error) {
	selectors := []string{
		"button:contains('Connect')",
		"button[aria-label*='Connect']",
		".connect-button",
		"#connect-btn",
	}

	for _, selector := range selectors {
		has, elem, err := page.Has(selector)
		if err == nil && has {
			return elem, nil
		}
	}

	return nil, fmt.Errorf("connect button not found")
}

func (s *Service) generatePersonalizedMessage(profile *search.Profile) string {
	templates := []string{
		fmt.Sprintf("Hi %s, I noticed your work at %s and would love to connect!", extractFirstName(profile.Name), profile.Company),
		fmt.Sprintf("Hello %s, I'm impressed by your background in %s. Let's connect!", extractFirstName(profile.Name), profile.Title),
		fmt.Sprintf("Hi %s, I see we share similar interests in the industry. Would love to connect and learn from your experience at %s.", extractFirstName(profile.Name), profile.Company),
		fmt.Sprintf("Hello %s, your work as %s is inspiring. I'd appreciate the opportunity to connect!", extractFirstName(profile.Name), profile.Title),
	}

	return templates[0]
}

func extractFirstName(fullName string) string {
	if fullName == "" {
		return "there"
	}

	parts := []rune{}
	for _, char := range fullName {
		if char == ' ' {
			break
		}
		parts = append(parts, char)
	}

	return string(parts)
}

func (s *Service) addConnectionNote(page *rod.Page, message string) error {
	noteSelectors := []string{
		"textarea[name='message']",
		"#custom-message",
		".connection-note",
		"textarea[placeholder*='Add a note']",
	}

	for _, selector := range noteSelectors {
		has, elem, err := page.Has(selector)
		if err == nil && has {
			if err := s.mouse.HoverElement(page, elem); err != nil {
				s.log.Warnf("Failed to hover note field: %v", err)
			}

			time.Sleep(s.timing.ThinkDelay())

			if err := s.mouse.ClickElement(page, elem); err != nil {
				s.log.Warnf("Failed to click note field: %v", err)
			}

			if err := s.typing.TypeIntoElement(elem, message, true); err != nil {
				return fmt.Errorf("failed to type message: %w", err)
			}

			return nil
		}
	}

	return fmt.Errorf("note field not found")
}

func (s *Service) submitConnectionRequest(page *rod.Page) error {
	submitSelectors := []string{
		"button:contains('Send')",
		"button[aria-label*='Send']",
		".send-button",
		"button[type='submit']",
	}

	for _, selector := range submitSelectors {
		has, elem, err := page.Has(selector)
		if err == nil && has {
			time.Sleep(s.timing.ThinkDelay())

			if err := s.mouse.ClickElement(page, elem); err != nil {
				return fmt.Errorf("failed to click send button: %w", err)
			}

			time.Sleep(s.timing.HumanizedDelay("click"))
			return nil
		}
	}

	return fmt.Errorf("send button not found")
}

func (s *Service) ShouldSkip(profileID string) bool {
	return s.db.HasProcessedProfile(profileID)
}

func (s *Service) CanSendRequest() bool {
	return s.limiter.CanSendConnection()
}
