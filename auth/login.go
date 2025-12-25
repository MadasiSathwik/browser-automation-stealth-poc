package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/automation-poc/browser-automation/config"
	"github.com/automation-poc/browser-automation/logger"
	"github.com/automation-poc/browser-automation/stealth"
	"github.com/go-rod/rod"
)

type LoginHandler struct {
	cfg    *config.Config
	log    *logger.Logger
	mouse  *stealth.MouseController
	typing *stealth.TypingSimulator
	timing *stealth.TimingController
}

func NewLoginHandler(cfg *config.Config, log *logger.Logger) *LoginHandler {
	return &LoginHandler{
		cfg:    cfg,
		log:    log,
		mouse:  stealth.NewMouseController(&cfg.Stealth.MouseMovement),
		typing: stealth.NewTypingSimulator(&cfg.Timing),
		timing: stealth.NewTimingController(&cfg.Timing),
	}
}

func (lh *LoginHandler) Login(ctx context.Context, page *rod.Page, email, password string) error {
	lh.log.Info("Navigating to login page")

	if err := page.Navigate(lh.cfg.Auth.LoginURL); err != nil {
		return fmt.Errorf("failed to navigate to login page: %w", err)
	}

	if err := page.WaitLoad(); err != nil {
		return fmt.Errorf("failed to wait for page load: %w", err)
	}

	time.Sleep(lh.timing.PageLoadDelay())

	if err := lh.detectSecurityCheckpoints(page); err != nil {
		return err
	}

	lh.log.Debug("Entering email address")
	emailInput, err := page.Element(lh.cfg.Auth.EmailSelector)
	if err != nil {
		return fmt.Errorf("failed to find email input: %w", err)
	}

	if err := lh.mouse.HoverElement(page, emailInput); err != nil {
		lh.log.Warnf("Failed to hover email input: %v", err)
	}

	time.Sleep(lh.timing.ThinkDelay())

	if err := lh.mouse.ClickElement(page, emailInput); err != nil {
		return fmt.Errorf("failed to click email input: %w", err)
	}

	if err := lh.typing.TypeIntoElement(emailInput, email, true); err != nil {
		return fmt.Errorf("failed to type email: %w", err)
	}

	lh.log.Debug("Entering password")
	passwordInput, err := page.Element(lh.cfg.Auth.PasswordSelector)
	if err != nil {
		return fmt.Errorf("failed to find password input: %w", err)
	}

	if err := lh.mouse.HoverElement(page, passwordInput); err != nil {
		lh.log.Warnf("Failed to hover password input: %v", err)
	}

	time.Sleep(lh.timing.ThinkDelay())

	if err := lh.mouse.ClickElement(page, passwordInput); err != nil {
		return fmt.Errorf("failed to click password input: %w", err)
	}

	if err := lh.typing.TypeWithThinkDelay(passwordInput, password); err != nil {
		return fmt.Errorf("failed to type password: %w", err)
	}

	time.Sleep(lh.timing.HumanizedDelay("form_fill"))

	lh.log.Debug("Submitting login form")
	submitButton, err := page.Element(lh.cfg.Auth.SubmitSelector)
	if err != nil {
		return fmt.Errorf("failed to find submit button: %w", err)
	}

	if err := lh.mouse.ClickElement(page, submitButton); err != nil {
		return fmt.Errorf("failed to click submit button: %w", err)
	}

	if err := page.WaitLoad(); err != nil {
		return fmt.Errorf("failed to wait for post-login page load: %w", err)
	}

	time.Sleep(lh.timing.PageLoadDelay())

	if err := lh.detectLoginFailure(page); err != nil {
		return err
	}

	if err := lh.detectSecurityCheckpoints(page); err != nil {
		return err
	}

	lh.log.Info("Login successful")
	return nil
}

func (lh *LoginHandler) detectSecurityCheckpoints(page *rod.Page) error {
	if captchaSelector, exists := lh.cfg.Auth.Selectors["captcha_detected"]; exists {
		has, _, err := page.Has(captchaSelector)
		if err == nil && has {
			lh.log.Warn("CAPTCHA detected - manual intervention required")
			return fmt.Errorf("CAPTCHA challenge detected")
		}
	}

	if twoFASelector, exists := lh.cfg.Auth.Selectors["2fa_detected"]; exists {
		has, _, err := page.Has(twoFASelector)
		if err == nil && has {
			lh.log.Warn("2FA prompt detected - manual intervention required")
			return fmt.Errorf("two-factor authentication required")
		}
	}

	return nil
}

func (lh *LoginHandler) detectLoginFailure(page *rod.Page) error {
	if errorSelector, exists := lh.cfg.Auth.Selectors["login_failed"]; exists {
		has, elem, err := page.Has(errorSelector)
		if err == nil && has {
			errorText, _ := elem.Text()
			lh.log.Errorf("Login failed: %s", errorText)
			return fmt.Errorf("login failed: %s", errorText)
		}
	}

	return nil
}

func (lh *LoginHandler) IsLoggedIn(page *rod.Page) bool {
	currentURL := page.MustInfo().URL

	if currentURL == lh.cfg.Auth.LoginURL {
		return false
	}

	if loginCheck, exists := lh.cfg.Auth.Selectors["logged_in_indicator"]; exists {
		has, _, _ := page.Has(loginCheck)
		return has
	}

	return true
}
