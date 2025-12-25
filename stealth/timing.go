package stealth

import (
	"math/rand"
	"time"

	"github.com/automation-poc/browser-automation/config"
)

type TimingController struct {
	cfg    *config.TimingConfig
	random *rand.Rand
}

func NewTimingController(cfg *config.TimingConfig) *TimingController {
	return &TimingController{
		cfg:    cfg,
		random: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (tc *TimingController) RandomDelay() time.Duration {
	seconds := tc.cfg.MinDelay + tc.random.Intn(tc.cfg.MaxDelay-tc.cfg.MinDelay+1)
	jitter := tc.random.Intn(1000)
	return time.Duration(seconds)*time.Second + time.Duration(jitter)*time.Millisecond
}

func (tc *TimingController) ThinkDelay() time.Duration {
	ms := tc.cfg.ThinkMin + tc.random.Intn(tc.cfg.ThinkMax-tc.cfg.ThinkMin+1)
	return time.Duration(ms) * time.Millisecond
}

func (tc *TimingController) BetweenActionsDelay() time.Duration {
	baseDelay := time.Duration(tc.cfg.BetweenActions) * time.Second

	jitterPercent := 0.2
	jitterRange := int(float64(tc.cfg.BetweenActions) * jitterPercent * 1000)
	jitter := tc.random.Intn(jitterRange*2) - jitterRange

	return baseDelay + time.Duration(jitter)*time.Millisecond
}

func (tc *TimingController) ScrollDelay() time.Duration {
	ms := 100 + tc.random.Intn(400)
	return time.Duration(ms) * time.Millisecond
}

func (tc *TimingController) PageLoadDelay() time.Duration {
	seconds := 2 + tc.random.Intn(3)
	return time.Duration(seconds) * time.Second
}

func (tc *TimingController) HumanizedDelay(action string) time.Duration {
	baseDelays := map[string]int{
		"click":         300,
		"scroll":        200,
		"read":          1000,
		"form_fill":     500,
		"navigation":    800,
		"hover":         150,
	}

	base, exists := baseDelays[action]
	if !exists {
		base = 500
	}

	variance := int(float64(base) * 0.4)
	jitter := tc.random.Intn(variance*2) - variance

	return time.Duration(base+jitter) * time.Millisecond
}

func (tc *TimingController) ExponentialBackoff(attempt int, maxDelay time.Duration) time.Duration {
	base := 1 * time.Second
	delay := time.Duration(1<<uint(attempt)) * base

	if delay > maxDelay {
		delay = maxDelay
	}

	jitter := time.Duration(tc.random.Intn(1000)) * time.Millisecond
	return delay + jitter
}

func (tc *TimingController) WaitForCondition(condition func() bool, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	checkInterval := 500 * time.Millisecond

	for time.Now().Before(deadline) {
		if condition() {
			return true
		}
		time.Sleep(checkInterval)
	}

	return false
}

func (tc *TimingController) ShouldActNow(businessHoursOnly bool) bool {
	if !businessHoursOnly {
		return true
	}

	return config.IsBusinessHours()
}

func (tc *TimingController) NextActionTime(businessHoursOnly bool) time.Time {
	if !businessHoursOnly {
		delay := tc.BetweenActionsDelay()
		return time.Now().Add(delay)
	}

	now := time.Now()

	if !config.IsBusinessHours() {
		nextBusinessDay := now
		for {
			nextBusinessDay = nextBusinessDay.Add(1 * time.Hour)
			if config.IsBusinessHours() {
				break
			}
			if nextBusinessDay.Sub(now) > 48*time.Hour {
				break
			}
		}
		return nextBusinessDay
	}

	delay := tc.BetweenActionsDelay()
	return now.Add(delay)
}
