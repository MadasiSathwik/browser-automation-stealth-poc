package messaging

import (
	"context"
	"fmt"
	"time"

	"github.com/automation-poc/browser-automation/config"
	"github.com/automation-poc/browser-automation/logger"
	"github.com/automation-poc/browser-automation/stealth"
	"github.com/automation-poc/browser-automation/storage"
	"github.com/go-rod/rod"
)

type Service struct {
	cfg            *config.Config
	db             *storage.Database
	log            *logger.Logger
	mouse          *stealth.MouseController
	typing         *stealth.TypingSimulator
	timing         *stealth.TimingController
	fingerprint    *stealth.FingerprintManager
	templateEngine *TemplateEngine
}

func NewService(cfg *config.Config, db *storage.Database, log *logger.Logger) *Service {
	return &Service{
		cfg:            cfg,
		db:             db,
		log:            log,
		mouse:          stealth.NewMouseController(&cfg.Stealth.MouseMovement),
		typing:         stealth.NewTypingSimulator(&cfg.Timing),
		timing:         stealth.NewTimingController(&cfg.Timing),
		fingerprint:    stealth.NewFingerprintManager(&cfg.Stealth),
		templateEngine: NewTemplateEngine(),
	}
}

func (s *Service) SendFollowUpMessage(ctx context.Context, page *rod.Page, conn *storage.ConnectionRequest) error {
	if s.cfg.DemoMode {
		s.log.Infof("[DEMO] Would send follow-up message to %s", conn.Name)
		return nil
	}

	s.log.Infof("Sending follow-up message to %s", conn.Name)

	templateName := s.templateEngine.SelectBestTemplate(conn)
	message, err := s.templateEngine.RenderForConnection(conn, templateName)
	if err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}

	s.log.Debugf("Using template '%s' for message", templateName)

	messageURL := fmt.Sprintf("%s/messaging/thread/%s", s.cfg.Search.BaseURL, conn.ProfileID)
	if err := page.Navigate(messageURL); err != nil {
		return fmt.Errorf("failed to navigate to message page: %w", err)
	}

	if err := page.WaitLoad(); err != nil {
		return fmt.Errorf("failed to wait for message page load: %w", err)
	}

	time.Sleep(s.timing.PageLoadDelay())

	if err := s.fingerprint.RandomScroll(page); err != nil {
		s.log.Warnf("Failed to perform random scroll: %v", err)
	}

	messageBox, err := s.findMessageBox(page)
	if err != nil {
		return fmt.Errorf("failed to find message box: %w", err)
	}

	if err := s.mouse.HoverElement(page, messageBox); err != nil {
		s.log.Warnf("Failed to hover message box: %v", err)
	}

	time.Sleep(s.timing.ThinkDelay())

	if err := s.mouse.ClickElement(page, messageBox); err != nil {
		return fmt.Errorf("failed to click message box: %w", err)
	}

	time.Sleep(s.timing.HumanizedDelay("form_fill"))

	if err := s.typing.TypeIntoElement(messageBox, message, true); err != nil {
		return fmt.Errorf("failed to type message: %w", err)
	}

	time.Sleep(s.timing.ThinkDelay())

	if err := s.sendMessage(page); err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	msg := &storage.Message{
		ProfileID: conn.ProfileID,
		Content:   message,
		SentAt:    time.Now(),
	}

	if err := s.db.SaveMessage(msg); err != nil {
		return fmt.Errorf("failed to save message: %w", err)
	}

	if err := s.db.IncrementMessageCount(); err != nil {
		s.log.Warnf("Failed to increment message count: %v", err)
	}

	s.log.Infof("Follow-up message sent successfully to %s", conn.Name)

	return nil
}

func (s *Service) findMessageBox(page *rod.Page) (*rod.Element, error) {
	selectors := []string{
		"textarea[placeholder*='message']",
		"div[contenteditable='true']",
		"textarea[name='message']",
		".message-input",
		"#message-box",
	}

	for _, selector := range selectors {
		has, elem, err := page.Has(selector)
		if err == nil && has {
			return elem, nil
		}
	}

	return nil, fmt.Errorf("message box not found")
}

func (s *Service) sendMessage(page *rod.Page) error {
	sendSelectors := []string{
		"button:contains('Send')",
		"button[aria-label*='Send']",
		"button[type='submit']",
		".send-button",
		"#send-btn",
	}

	for _, selector := range sendSelectors {
		has, elem, err := page.Has(selector)
		if err == nil && has {
			if err := s.mouse.ClickElement(page, elem); err != nil {
				return fmt.Errorf("failed to click send button: %w", err)
			}

			time.Sleep(s.timing.HumanizedDelay("click"))
			return nil
		}
	}

	has, _, err := page.Has("textarea")
	if err == nil && has {
		s.log.Debug("Attempting to send message with Enter key")
		if err := page.Keyboard.Press("Enter"); err != nil {
			return fmt.Errorf("failed to press Enter: %w", err)
		}
		time.Sleep(s.timing.HumanizedDelay("click"))
		return nil
	}

	return fmt.Errorf("send button not found")
}

func (s *Service) HasSentMessage(profileID string) bool {
	conn, err := s.db.GetConnectionRequest(profileID)
	if err != nil || conn == nil {
		return false
	}

	return conn.LastMessageAt != nil
}

func (s *Service) SendCustomMessage(ctx context.Context, page *rod.Page, conn *storage.ConnectionRequest, message string) error {
	if s.cfg.DemoMode {
		s.log.Infof("[DEMO] Would send custom message to %s: %s", conn.Name, message)
		return nil
	}

	s.log.Infof("Sending custom message to %s", conn.Name)

	messageURL := fmt.Sprintf("%s/messaging/thread/%s", s.cfg.Search.BaseURL, conn.ProfileID)
	if err := page.Navigate(messageURL); err != nil {
		return fmt.Errorf("failed to navigate to message page: %w", err)
	}

	if err := page.WaitLoad(); err != nil {
		return fmt.Errorf("failed to wait for message page load: %w", err)
	}

	time.Sleep(s.timing.PageLoadDelay())

	messageBox, err := s.findMessageBox(page)
	if err != nil {
		return fmt.Errorf("failed to find message box: %w", err)
	}

	if err := s.typing.TypeIntoElement(messageBox, message, true); err != nil {
		return fmt.Errorf("failed to type message: %w", err)
	}

	if err := s.sendMessage(page); err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	msg := &storage.Message{
		ProfileID: conn.ProfileID,
		Content:   message,
		SentAt:    time.Now(),
	}

	if err := s.db.SaveMessage(msg); err != nil {
		return fmt.Errorf("failed to save message: %w", err)
	}

	s.log.Infof("Custom message sent successfully to %s", conn.Name)

	return nil
}
