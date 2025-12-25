package auth

import (
	"fmt"
	"os"

	"github.com/automation-poc/browser-automation/config"
	"github.com/automation-poc/browser-automation/logger"
	"github.com/automation-poc/browser-automation/storage"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

type Service struct {
	cfg          *config.Config
	db           *storage.Database
	log          *logger.Logger
	loginHandler *LoginHandler
}

func NewService(cfg *config.Config, db *storage.Database, log *logger.Logger) *Service {
	return &Service{
		cfg:          cfg,
		db:           db,
		log:          log,
		loginHandler: NewLoginHandler(cfg, log),
	}
}

func (s *Service) Login(ctx context.Context, page *rod.Page) error {
	email := os.Getenv("LOGIN_EMAIL")
	password := os.Getenv("LOGIN_PASSWORD")

	if email == "" || password == "" {
		return fmt.Errorf("LOGIN_EMAIL and LOGIN_PASSWORD environment variables must be set")
	}

	if s.cfg.SafeMode {
		s.log.Warn("SAFE MODE: Login simulation only (no real credentials used)")
		email = "test@example.com"
		password = "test-password"
	}

	if s.cfg.DemoMode {
		s.log.Info("[DEMO] Would attempt login with provided credentials")
		return nil
	}

	return s.loginHandler.Login(ctx, page, email, password)
}

func (s *Service) SaveSession(page *rod.Page) error {
	s.log.Debug("Saving session cookies")

	cookies, err := page.Cookies([]string{})
	if err != nil {
		return fmt.Errorf("failed to get cookies: %w", err)
	}

	var storageCookies []storage.Cookie
	for _, cookie := range cookies {
		storageCookies = append(storageCookies, storage.Cookie{
			Name:     cookie.Name,
			Value:    cookie.Value,
			Domain:   cookie.Domain,
			Path:     cookie.Path,
			Expires:  cookie.Expires.Time(),
			HTTPOnly: cookie.HTTPOnly,
			Secure:   cookie.Secure,
		})
	}

	if err := s.db.SaveCookies(storageCookies); err != nil {
		return fmt.Errorf("failed to save cookies: %w", err)
	}

	s.log.Info("Session saved successfully")
	return nil
}

func (s *Service) LoadSession(page *rod.Page) (bool, error) {
	s.log.Debug("Attempting to load saved session")

	cookies, err := s.db.LoadCookies()
	if err != nil {
		s.log.Warnf("Failed to load cookies: %v", err)
		return false, err
	}

	if len(cookies) == 0 {
		s.log.Debug("No saved session found")
		return false, nil
	}

	if err := page.Navigate(s.cfg.Auth.LoginURL); err != nil {
		return false, fmt.Errorf("failed to navigate for cookie injection: %w", err)
	}

	for _, cookie := range cookies {
		protoCookie := &proto.NetworkCookieParam{
			Name:     cookie.Name,
			Value:    cookie.Value,
			Domain:   cookie.Domain,
			Path:     cookie.Path,
			Expires:  proto.TimeSinceEpoch(cookie.Expires.Unix()),
			HTTPOnly: cookie.HTTPOnly,
			Secure:   cookie.Secure,
		}

		if err := page.SetCookies([]*proto.NetworkCookieParam{protoCookie}); err != nil {
			s.log.Warnf("Failed to set cookie %s: %v", cookie.Name, err)
		}
	}

	if err := page.Navigate(s.cfg.Auth.LoginURL); err != nil {
		return false, fmt.Errorf("failed to navigate after cookie injection: %w", err)
	}

	if err := page.WaitLoad(); err != nil {
		return false, fmt.Errorf("failed to wait for page load: %w", err)
	}

	isLoggedIn := s.loginHandler.IsLoggedIn(page)

	if !isLoggedIn {
		s.log.Debug("Session is invalid, login required")
		s.db.InvalidateSession()
		return false, nil
	}

	s.log.Info("Session restored successfully")
	return true, nil
}
