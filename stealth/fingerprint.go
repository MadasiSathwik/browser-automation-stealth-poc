package stealth

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/automation-poc/browser-automation/config"
	"github.com/go-rod/rod"
	"github.com/go-rod/stealth"
)

type FingerprintManager struct {
	cfg    *config.StealthConfig
	random *rand.Rand
}

func NewFingerprintManager(cfg *config.StealthConfig) *FingerprintManager {
	return &FingerprintManager{
		cfg:    cfg,
		random: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (fm *FingerprintManager) ApplyStealthSettings(page *rod.Page) error {
	stealth.MustPage(page)

	if err := fm.randomizeViewport(page); err != nil {
		return fmt.Errorf("failed to set viewport: %w", err)
	}

	if err := fm.setUserAgent(page); err != nil {
		return fmt.Errorf("failed to set user agent: %w", err)
	}

	if err := fm.hideWebDriver(page); err != nil {
		return fmt.Errorf("failed to hide webdriver: %w", err)
	}

	if err := fm.injectHumanBehaviorMarkers(page); err != nil {
		return fmt.Errorf("failed to inject behavior markers: %w", err)
	}

	return nil
}

func (fm *FingerprintManager) randomizeViewport(page *rod.Page) error {
	if len(fm.cfg.ViewportSizes) == 0 {
		return nil
	}

	viewport := fm.cfg.ViewportSizes[fm.random.Intn(len(fm.cfg.ViewportSizes))]

	jitterWidth := fm.random.Intn(20) - 10
	jitterHeight := fm.random.Intn(20) - 10

	finalWidth := viewport.Width + jitterWidth
	finalHeight := viewport.Height + jitterHeight

	return page.SetViewport(&rod.Viewport{
		Width:  finalWidth,
		Height: finalHeight,
	})
}

func (fm *FingerprintManager) setUserAgent(page *rod.Page) error {
	if len(fm.cfg.UserAgentRotation) == 0 {
		return nil
	}

	userAgent := fm.cfg.UserAgentRotation[fm.random.Intn(len(fm.cfg.UserAgentRotation))]

	return page.SetUserAgent(&rod.UserAgent{
		UserAgent: userAgent,
	})
}

func (fm *FingerprintManager) hideWebDriver(page *rod.Page) error {
	script := `
		Object.defineProperty(navigator, 'webdriver', {
			get: () => undefined
		});

		window.navigator.chrome = {
			runtime: {}
		};

		Object.defineProperty(navigator, 'plugins', {
			get: () => [1, 2, 3, 4, 5]
		});

		Object.defineProperty(navigator, 'languages', {
			get: () => ['en-US', 'en']
		});

		const originalQuery = window.navigator.permissions.query;
		window.navigator.permissions.query = (parameters) => (
			parameters.name === 'notifications' ?
				Promise.resolve({ state: Notification.permission }) :
				originalQuery(parameters)
		);
	`

	_, err := page.Eval(script)
	return err
}

func (fm *FingerprintManager) injectHumanBehaviorMarkers(page *rod.Page) error {
	script := `
		window.__humanBehaviorMarkers = {
			mouseMovements: 0,
			keystrokes: 0,
			scrolls: 0,
			timestamp: Date.now()
		};

		document.addEventListener('mousemove', () => {
			window.__humanBehaviorMarkers.mouseMovements++;
		});

		document.addEventListener('keydown', () => {
			window.__humanBehaviorMarkers.keystrokes++;
		});

		document.addEventListener('scroll', () => {
			window.__humanBehaviorMarkers.scrolls++;
		});
	`

	_, err := page.Eval(script)
	return err
}

func (fm *FingerprintManager) RandomScroll(page *rod.Page) error {
	if !fm.cfg.RandomScrolling {
		return nil
	}

	scrollTypes := []string{"smooth", "partial", "bounce", "read"}
	scrollType := scrollTypes[fm.random.Intn(len(scrollTypes))]

	switch scrollType {
	case "smooth":
		return fm.smoothScroll(page)
	case "partial":
		return fm.partialScroll(page)
	case "bounce":
		return fm.bounceScroll(page)
	case "read":
		return fm.readingScroll(page)
	}

	return nil
}

func (fm *FingerprintManager) smoothScroll(page *rod.Page) error {
	scrollDistance := 200 + fm.random.Intn(400)
	_, err := page.Eval(fmt.Sprintf(`window.scrollBy({top: %d, behavior: 'smooth'})`, scrollDistance))
	return err
}

func (fm *FingerprintManager) partialScroll(page *rod.Page) error {
	scrollDistance := 50 + fm.random.Intn(150)
	_, err := page.Eval(fmt.Sprintf(`window.scrollBy(0, %d)`, scrollDistance))
	time.Sleep(time.Duration(100+fm.random.Intn(300)) * time.Millisecond)
	return err
}

func (fm *FingerprintManager) bounceScroll(page *rod.Page) error {
	scrollDown := 300 + fm.random.Intn(200)
	_, err := page.Eval(fmt.Sprintf(`window.scrollBy(0, %d)`, scrollDown))
	if err != nil {
		return err
	}

	time.Sleep(time.Duration(200+fm.random.Intn(300)) * time.Millisecond)

	scrollUp := -(scrollDown / 4)
	_, err = page.Eval(fmt.Sprintf(`window.scrollBy(0, %d)`, scrollUp))
	return err
}

func (fm *FingerprintManager) readingScroll(page *rod.Page) error {
	numScrolls := 2 + fm.random.Intn(4)

	for i := 0; i < numScrolls; i++ {
		scrollDistance := 100 + fm.random.Intn(200)
		_, err := page.Eval(fmt.Sprintf(`window.scrollBy(0, %d)`, scrollDistance))
		if err != nil {
			return err
		}

		readPause := time.Duration(1000+fm.random.Intn(3000)) * time.Millisecond
		time.Sleep(readPause)
	}

	return nil
}

func (fm *FingerprintManager) SimulateReading(page *rod.Page, duration time.Duration) error {
	endTime := time.Now().Add(duration)

	for time.Now().Before(endTime) {
		if fm.random.Float64() < 0.7 {
			if err := fm.RandomScroll(page); err != nil {
				return err
			}
		}

		pause := time.Duration(2000+fm.random.Intn(5000)) * time.Millisecond
		time.Sleep(pause)
	}

	return nil
}

func (fm *FingerprintManager) GetRandomUserAgent() string {
	if len(fm.cfg.UserAgentRotation) == 0 {
		return "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"
	}

	return fm.cfg.UserAgentRotation[fm.random.Intn(len(fm.cfg.UserAgentRotation))]
}
