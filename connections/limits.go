package connections

import (
	"time"

	"github.com/automation-poc/browser-automation/config"
	"github.com/automation-poc/browser-automation/logger"
	"github.com/automation-poc/browser-automation/storage"
)

type RateLimiter struct {
	cfg *config.Config
	db  *storage.Database
	log *logger.Logger
}

func NewRateLimiter(cfg *config.Config, db *storage.Database, log *logger.Logger) *RateLimiter {
	return &RateLimiter{
		cfg: cfg,
		db:  db,
		log: log,
	}
}

func (rl *RateLimiter) CanSendConnection() bool {
	stats, err := rl.db.GetTodayStats()
	if err != nil {
		rl.log.Errorf("Failed to get daily stats: %v", err)
		return false
	}

	if stats.ConnectionsSent >= rl.cfg.Limits.DailyConnections {
		rl.log.Warnf("Daily connection limit reached: %d/%d", stats.ConnectionsSent, rl.cfg.Limits.DailyConnections)
		return false
	}

	if !rl.isWithinBusinessHours() {
		rl.log.Debug("Outside business hours")
		return false
	}

	return true
}

func (rl *RateLimiter) CanSendMessage() bool {
	stats, err := rl.db.GetTodayStats()
	if err != nil {
		rl.log.Errorf("Failed to get daily stats: %v", err)
		return false
	}

	if stats.MessagesSent >= rl.cfg.Limits.DailyMessages {
		rl.log.Warnf("Daily message limit reached: %d/%d", stats.MessagesSent, rl.cfg.Limits.DailyMessages)
		return false
	}

	return true
}

func (rl *RateLimiter) isWithinBusinessHours() bool {
	if !rl.cfg.Stealth.BusinessHoursOnly {
		return true
	}

	return config.IsBusinessHours()
}

func (rl *RateLimiter) GetRemainingConnections() int {
	stats, err := rl.db.GetTodayStats()
	if err != nil {
		return 0
	}

	remaining := rl.cfg.Limits.DailyConnections - stats.ConnectionsSent
	if remaining < 0 {
		return 0
	}

	return remaining
}

func (rl *RateLimiter) GetRemainingMessages() int {
	stats, err := rl.db.GetTodayStats()
	if err != nil {
		return 0
	}

	remaining := rl.cfg.Limits.DailyMessages - stats.MessagesSent
	if remaining < 0 {
		return 0
	}

	return remaining
}

func (rl *RateLimiter) WaitUntilNextWindow() time.Duration {
	now := time.Now()

	if rl.cfg.Stealth.BusinessHoursOnly && !config.IsBusinessHours() {
		nextBusinessDay := now
		for !config.IsBusinessHours() {
			nextBusinessDay = nextBusinessDay.Add(1 * time.Hour)
			if nextBusinessDay.Sub(now) > 48*time.Hour {
				break
			}
		}
		return nextBusinessDay.Sub(now)
	}

	tomorrow := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
	return tomorrow.Sub(now)
}

func (rl *RateLimiter) EnforceDelay(action string) {
	stats, _ := rl.db.GetTodayStats()

	baseDelay := time.Duration(rl.cfg.Timing.BetweenActions) * time.Second

	if stats != nil {
		actionsToday := stats.ConnectionsSent + stats.MessagesSent
		if actionsToday > 30 {
			baseDelay = baseDelay * 2
		}
	}

	time.Sleep(baseDelay)
}
