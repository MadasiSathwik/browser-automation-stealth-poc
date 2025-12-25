package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/automation-poc/browser-automation/auth"
	"github.com/automation-poc/browser-automation/config"
	"github.com/automation-poc/browser-automation/connections"
	"github.com/automation-poc/browser-automation/logger"
	"github.com/automation-poc/browser-automation/messaging"
	"github.com/automation-poc/browser-automation/search"
	"github.com/automation-poc/browser-automation/storage"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

var (
	configPath = flag.String("config", "config.yaml", "Path to configuration file")
	demoMode   = flag.Bool("demo", false, "Run in demo mode (prints actions without executing)")
	safeMode   = flag.Bool("safe", true, "Enable safe mode (prevents real-world execution)")
)

func main() {
	flag.Parse()

	log := logger.New()
	log.Info("Starting Browser Automation POC")
	log.Infof("Demo Mode: %v | Safe Mode: %v", *demoMode, *safeMode)

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	if *demoMode {
		cfg.DemoMode = true
		log.Info("DEMO MODE: Actions will be printed but not executed")
	}

	if *safeMode {
		cfg.SafeMode = true
		log.Warn("SAFE MODE: Real-world execution is disabled")
	}

	if err := cfg.Validate(); err != nil {
		log.Fatalf("Configuration validation failed: %v", err)
	}

	db, err := storage.NewDatabase(cfg.Database.Path)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Info("Received shutdown signal, cleaning up...")
		cancel()
	}()

	if err := run(ctx, cfg, db, log); err != nil {
		log.Fatalf("Automation failed: %v", err)
	}

	log.Info("Automation completed successfully")
}

func run(ctx context.Context, cfg *config.Config, db *storage.Database, log *logger.Logger) error {
	if cfg.DemoMode {
		return runDemo(ctx, cfg, db, log)
	}

	browser, page, err := initializeBrowser(cfg, log)
	if err != nil {
		return fmt.Errorf("browser initialization failed: %w", err)
	}
	defer browser.Close()

	authService := auth.NewService(cfg, db, log)
	sessionValid, err := authService.LoadSession(page)
	if err != nil {
		log.Warnf("Failed to load session: %v", err)
	}

	if !sessionValid {
		log.Info("No valid session found, attempting login")
		if err := authService.Login(ctx, page); err != nil {
			return fmt.Errorf("login failed: %w", err)
		}

		if err := authService.SaveSession(page); err != nil {
			log.Warnf("Failed to save session: %v", err)
		}
	}

	searchService := search.NewService(cfg, db, log)
	profiles, err := searchService.SearchProfiles(ctx, page, cfg.Search)
	if err != nil {
		return fmt.Errorf("search failed: %w", err)
	}

	log.Infof("Found %d profiles to process", len(profiles))

	connectionService := connections.NewService(cfg, db, log)
	messageService := messaging.NewService(cfg, db, log)

	for _, profile := range profiles {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if connectionService.ShouldSkip(profile.ID) {
			log.Debugf("Skipping profile %s (already processed)", profile.ID)
			continue
		}

		if !connectionService.CanSendRequest() {
			log.Info("Daily connection limit reached, stopping")
			break
		}

		if err := connectionService.SendConnectionRequest(ctx, page, profile); err != nil {
			log.Errorf("Failed to send connection to %s: %v", profile.Name, err)
			continue
		}

		log.Infof("Connection request sent to %s", profile.Name)
		time.Sleep(time.Duration(cfg.Timing.BetweenActions) * time.Second)
	}

	acceptedConnections, err := db.GetAcceptedConnections()
	if err != nil {
		log.Errorf("Failed to get accepted connections: %v", err)
		return nil
	}

	for _, conn := range acceptedConnections {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if messageService.HasSentMessage(conn.ProfileID) {
			continue
		}

		if err := messageService.SendFollowUpMessage(ctx, page, conn); err != nil {
			log.Errorf("Failed to send message to %s: %v", conn.Name, err)
			continue
		}

		log.Infof("Follow-up message sent to %s", conn.Name)
		time.Sleep(time.Duration(cfg.Timing.BetweenActions) * time.Second)
	}

	return nil
}

func runDemo(ctx context.Context, cfg *config.Config, db *storage.Database, log *logger.Logger) error {
	log.Info("=== DEMO MODE: Simulated Execution ===")
	log.Info("")

	log.Info("[ACTION] Initializing browser with stealth configuration")
	log.Info("  - User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64)...")
	log.Info("  - Viewport: 1366x768 (randomized)")
	log.Info("  - WebDriver flag: hidden")
	log.Info("")

	log.Info("[ACTION] Attempting to load saved session")
	log.Info("  - No valid session found")
	log.Info("")

	log.Info("[ACTION] Navigating to login page")
	log.Info("  - URL: https://mock-professional-network.example.com/login")
	log.Info("  - Applying human-like mouse movement (Bézier curve)")
	log.Info("  - Random scroll: 120px with velocity variance")
	log.Info("")

	log.Info("[ACTION] Entering credentials")
	log.Info("  - Typing email with realistic delays (45-120ms per char)")
	log.Info("  - Simulated typo: 'usre@' -> backspace -> 'user@'")
	log.Info("  - Hovering over password field before click")
	log.Info("  - Typing password with increased delay (security field)")
	log.Info("")

	log.Info("[ACTION] Submitting login form")
	log.Info("  - Human-like delay before submit: 823ms")
	log.Info("  - Session cookies saved")
	log.Info("")

	log.Info("[ACTION] Searching for profiles")
	log.Info("  - Query: 'Software Engineer at Tech Companies'")
	log.Info("  - Applying search filters via DOM selectors")
	log.Info("  - Random scroll to trigger lazy-loading")
	log.Info("  - Found 15 profiles on page 1")
	log.Info("")

	mockProfiles := []string{
		"Alice Johnson - Senior Software Engineer @ TechCorp",
		"Bob Smith - Engineering Manager @ InnovateLabs",
		"Carol White - Full Stack Developer @ StartupXYZ",
	}

	for i, profile := range mockProfiles {
		log.Infof("[ACTION] Processing profile %d: %s", i+1, profile)
		log.Info("  - Hovering over profile card (250ms)")
		log.Info("  - Moving mouse to 'Connect' button with Bézier curve")
		log.Info("  - Random micro-correction: +3px Y-axis")
		log.Info("  - Clicking 'Connect'")
		log.Info("  - Adding personalized note:")
		log.Info("    'Hi Alice, I noticed we share an interest in distributed systems...'")
		log.Info("  - Typing note with human-like delays")
		log.Info("  - Clicking 'Send' with 1.2s think delay")
		log.Info("  - Connection request recorded in database")
		log.Info("  - Cooldown: 47s (randomized)")
		log.Info("")

		time.Sleep(500 * time.Millisecond)
	}

	log.Info("[ACTION] Checking for accepted connections")
	log.Info("  - Found 2 accepted connections from previous runs")
	log.Info("")

	log.Info("[ACTION] Sending follow-up message")
	log.Info("  - To: Bob Smith")
	log.Info("  - Template: 'Thanks for connecting! I'd love to learn more about {{company}}'")
	log.Info("  - Rendered: 'Thanks for connecting! I'd love to learn more about InnovateLabs'")
	log.Info("  - Message sent with realistic typing simulation")
	log.Info("")

	log.Info("[SUMMARY] Demo completed successfully")
	log.Info("  - Total profiles processed: 3")
	log.Info("  - Connection requests sent: 3")
	log.Info("  - Follow-up messages sent: 1")
	log.Info("  - Average action delay: 45s")
	log.Info("  - Stealth techniques applied: 8")
	log.Info("")

	return nil
}

func initializeBrowser(cfg *config.Config, log *logger.Logger) (*rod.Browser, *rod.Page, error) {
	log.Info("Initializing browser with stealth configuration")

	l := launcher.New().
		Headless(cfg.Browser.Headless).
		Leakless(true)

	if cfg.Browser.NoSandbox {
		l = l.NoSandbox(true)
	}

	url, err := l.Launch()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to launch browser: %w", err)
	}

	browser := rod.New().ControlURL(url).MustConnect()

	page := browser.MustPage()

	return browser, page, nil
}
