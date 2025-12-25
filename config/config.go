package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Browser  BrowserConfig  `yaml:"browser"`
	Auth     AuthConfig     `yaml:"auth"`
	Search   SearchConfig   `yaml:"search"`
	Limits   LimitsConfig   `yaml:"limits"`
	Timing   TimingConfig   `yaml:"timing"`
	Stealth  StealthConfig  `yaml:"stealth"`
	Database DatabaseConfig `yaml:"database"`
	DemoMode bool           `yaml:"demo_mode"`
	SafeMode bool           `yaml:"safe_mode"`
}

type BrowserConfig struct {
	Headless  bool   `yaml:"headless"`
	NoSandbox bool   `yaml:"no_sandbox"`
	UserAgent string `yaml:"user_agent"`
}

type AuthConfig struct {
	LoginURL        string            `yaml:"login_url"`
	EmailSelector   string            `yaml:"email_selector"`
	PasswordSelector string           `yaml:"password_selector"`
	SubmitSelector  string            `yaml:"submit_selector"`
	Selectors       map[string]string `yaml:"selectors"`
}

type SearchConfig struct {
	BaseURL       string            `yaml:"base_url"`
	Query         string            `yaml:"query"`
	Filters       map[string]string `yaml:"filters"`
	MaxPages      int               `yaml:"max_pages"`
	ResultsPerPage int              `yaml:"results_per_page"`
	Selectors     map[string]string `yaml:"selectors"`
}

type LimitsConfig struct {
	DailyConnections int `yaml:"daily_connections"`
	HourlyConnections int `yaml:"hourly_connections"`
	DailyMessages    int `yaml:"daily_messages"`
}

type TimingConfig struct {
	MinDelay       int `yaml:"min_delay"`
	MaxDelay       int `yaml:"max_delay"`
	BetweenActions int `yaml:"between_actions"`
	TypingMin      int `yaml:"typing_min"`
	TypingMax      int `yaml:"typing_max"`
	ThinkMin       int `yaml:"think_min"`
	ThinkMax       int `yaml:"think_max"`
}

type StealthConfig struct {
	MouseMovement     MouseMovementConfig `yaml:"mouse_movement"`
	RandomScrolling   bool                `yaml:"random_scrolling"`
	BusinessHoursOnly bool                `yaml:"business_hours_only"`
	UserAgentRotation []string            `yaml:"user_agent_rotation"`
	ViewportSizes     []ViewportSize      `yaml:"viewport_sizes"`
}

type MouseMovementConfig struct {
	Enabled          bool    `yaml:"enabled"`
	BezierCurves     bool    `yaml:"bezier_curves"`
	Overshoot        bool    `yaml:"overshoot"`
	MicroCorrections bool    `yaml:"micro_corrections"`
	VelocityVariance float64 `yaml:"velocity_variance"`
}

type ViewportSize struct {
	Width  int `yaml:"width"`
	Height int `yaml:"height"`
}

type DatabaseConfig struct {
	Path string `yaml:"path"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return DefaultConfig(), nil
		}
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	cfg.applyEnvironmentOverrides()

	return &cfg, nil
}

func (c *Config) applyEnvironmentOverrides() {
	if loginURL := os.Getenv("LOGIN_URL"); loginURL != "" {
		c.Auth.LoginURL = loginURL
	}
	if searchQuery := os.Getenv("SEARCH_QUERY"); searchQuery != "" {
		c.Search.Query = searchQuery
	}
	if dbPath := os.Getenv("DATABASE_PATH"); dbPath != "" {
		c.Database.Path = dbPath
	}
}

func (c *Config) Validate() error {
	if c.SafeMode && c.Auth.LoginURL != "" {
		if !isTestURL(c.Auth.LoginURL) {
			return fmt.Errorf("safe mode enabled but login URL appears to be a real site: %s", c.Auth.LoginURL)
		}
	}

	if c.Limits.DailyConnections <= 0 {
		return fmt.Errorf("daily_connections must be greater than 0")
	}

	if c.Timing.MinDelay < 0 || c.Timing.MaxDelay < c.Timing.MinDelay {
		return fmt.Errorf("invalid timing configuration")
	}

	if c.Database.Path == "" {
		c.Database.Path = "automation.db"
	}

	return nil
}

func isTestURL(url string) bool {
	testDomains := []string{
		"localhost",
		"127.0.0.1",
		"example.com",
		"test.com",
		"mock",
		"demo",
	}

	for _, domain := range testDomains {
		if len(url) > 0 && contains(url, domain) {
			return true
		}
	}

	return false
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func DefaultConfig() *Config {
	return &Config{
		Browser: BrowserConfig{
			Headless:  false,
			NoSandbox: false,
			UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
		},
		Auth: AuthConfig{
			LoginURL:         "https://mock-network.example.com/login",
			EmailSelector:    "#email",
			PasswordSelector: "#password",
			SubmitSelector:   "button[type='submit']",
			Selectors: map[string]string{
				"captcha_detected": ".captcha-container",
				"2fa_detected":     ".two-factor-prompt",
				"login_failed":     ".error-message",
			},
		},
		Search: SearchConfig{
			BaseURL:        "https://mock-network.example.com/search",
			Query:          "Software Engineer",
			MaxPages:       5,
			ResultsPerPage: 10,
			Selectors: map[string]string{
				"profile_card":      ".profile-card",
				"profile_name":      ".profile-name",
				"profile_title":     ".profile-title",
				"profile_company":   ".profile-company",
				"profile_link":      ".profile-link",
				"next_page_button":  ".pagination-next",
			},
		},
		Limits: LimitsConfig{
			DailyConnections:  50,
			HourlyConnections: 10,
			DailyMessages:     30,
		},
		Timing: TimingConfig{
			MinDelay:       2,
			MaxDelay:       8,
			BetweenActions: 30,
			TypingMin:      45,
			TypingMax:      120,
			ThinkMin:       500,
			ThinkMax:       2000,
		},
		Stealth: StealthConfig{
			MouseMovement: MouseMovementConfig{
				Enabled:          true,
				BezierCurves:     true,
				Overshoot:        true,
				MicroCorrections: true,
				VelocityVariance: 0.3,
			},
			RandomScrolling:   true,
			BusinessHoursOnly: false,
			UserAgentRotation: []string{
				"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
				"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36",
				"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36",
			},
			ViewportSizes: []ViewportSize{
				{Width: 1920, Height: 1080},
				{Width: 1366, Height: 768},
				{Width: 1536, Height: 864},
			},
		},
		Database: DatabaseConfig{
			Path: "automation.db",
		},
		SafeMode: true,
		DemoMode: false,
	}
}

func IsBusinessHours() bool {
	now := time.Now()
	hour := now.Hour()
	weekday := now.Weekday()

	if weekday == time.Saturday || weekday == time.Sunday {
		return false
	}

	return hour >= 9 && hour < 17
}
