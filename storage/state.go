package storage

import (
	"encoding/json"
	"fmt"
	"time"
)

type SessionState struct {
	Cookies      []Cookie  `json:"cookies"`
	LastActivity time.Time `json:"last_activity"`
	Valid        bool      `json:"valid"`
}

type Cookie struct {
	Name     string    `json:"name"`
	Value    string    `json:"value"`
	Domain   string    `json:"domain"`
	Path     string    `json:"path"`
	Expires  time.Time `json:"expires"`
	HTTPOnly bool      `json:"http_only"`
	Secure   bool      `json:"secure"`
}

func (d *Database) SaveCookies(cookies []Cookie) error {
	state := SessionState{
		Cookies:      cookies,
		LastActivity: time.Now(),
		Valid:        true,
	}

	data, err := json.Marshal(state)
	if err != nil {
		return fmt.Errorf("failed to marshal session state: %w", err)
	}

	return d.SaveSessionState("session", string(data))
}

func (d *Database) LoadCookies() ([]Cookie, error) {
	data, err := d.GetSessionState("session")
	if err != nil {
		return nil, err
	}

	if data == "" {
		return nil, nil
	}

	var state SessionState
	if err := json.Unmarshal([]byte(data), &state); err != nil {
		return nil, fmt.Errorf("failed to unmarshal session state: %w", err)
	}

	if !state.Valid {
		return nil, fmt.Errorf("session is marked as invalid")
	}

	maxAge := 7 * 24 * time.Hour
	if time.Since(state.LastActivity) > maxAge {
		return nil, fmt.Errorf("session expired (last activity: %s)", state.LastActivity)
	}

	return state.Cookies, nil
}

func (d *Database) InvalidateSession() error {
	return d.SaveSessionState("session", "")
}

func (d *Database) GetAutomationState() (map[string]interface{}, error) {
	stats, err := d.GetTodayStats()
	if err != nil {
		return nil, err
	}

	state := map[string]interface{}{
		"date":              stats.Date,
		"connections_sent":  stats.ConnectionsSent,
		"messages_sent":     stats.MessagesSent,
		"connections_limit": stats.ConnectionsLimit,
		"messages_limit":    stats.MessagesLimit,
		"can_send_more":     stats.ConnectionsSent < stats.ConnectionsLimit,
	}

	return state, nil
}
