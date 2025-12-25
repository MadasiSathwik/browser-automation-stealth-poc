package storage

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	db *sql.DB
}

type ConnectionRequest struct {
	ID            int64
	ProfileID     string
	Name          string
	Title         string
	Company       string
	Message       string
	Status        string
	SentAt        time.Time
	AcceptedAt    *time.Time
	LastMessageAt *time.Time
}

type Message struct {
	ID        int64
	ProfileID string
	Content   string
	SentAt    time.Time
}

type DailyStats struct {
	Date             string
	ConnectionsSent  int
	MessagesS sent    int
	ConnectionsLimit int
	MessagesLimit    int
}

func NewDatabase(path string) (*Database, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	d := &Database{db: db}

	if err := d.initSchema(); err != nil {
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return d, nil
}

func (d *Database) initSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS connection_requests (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		profile_id TEXT UNIQUE NOT NULL,
		name TEXT NOT NULL,
		title TEXT,
		company TEXT,
		message TEXT,
		status TEXT DEFAULT 'pending',
		sent_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		accepted_at TIMESTAMP,
		last_message_at TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS messages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		profile_id TEXT NOT NULL,
		content TEXT NOT NULL,
		sent_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(profile_id) REFERENCES connection_requests(profile_id)
	);

	CREATE TABLE IF NOT EXISTS daily_stats (
		date TEXT PRIMARY KEY,
		connections_sent INTEGER DEFAULT 0,
		messages_sent INTEGER DEFAULT 0,
		connections_limit INTEGER DEFAULT 50,
		messages_limit INTEGER DEFAULT 30
	);

	CREATE TABLE IF NOT EXISTS session_state (
		key TEXT PRIMARY KEY,
		value TEXT,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_connection_status ON connection_requests(status);
	CREATE INDEX IF NOT EXISTS idx_connection_sent_at ON connection_requests(sent_at);
	CREATE INDEX IF NOT EXISTS idx_messages_profile ON messages(profile_id);
	CREATE INDEX IF NOT EXISTS idx_messages_sent_at ON messages(sent_at);
	`

	_, err := d.db.Exec(schema)
	return err
}

func (d *Database) SaveConnectionRequest(req *ConnectionRequest) error {
	query := `
		INSERT INTO connection_requests (profile_id, name, title, company, message, status)
		VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT(profile_id) DO UPDATE SET
			name = excluded.name,
			title = excluded.title,
			company = excluded.company,
			message = excluded.message,
			status = excluded.status
	`

	_, err := d.db.Exec(query, req.ProfileID, req.Name, req.Title, req.Company, req.Message, req.Status)
	return err
}

func (d *Database) GetConnectionRequest(profileID string) (*ConnectionRequest, error) {
	query := `
		SELECT id, profile_id, name, title, company, message, status, sent_at, accepted_at, last_message_at
		FROM connection_requests
		WHERE profile_id = ?
	`

	var req ConnectionRequest
	err := d.db.QueryRow(query, profileID).Scan(
		&req.ID, &req.ProfileID, &req.Name, &req.Title, &req.Company,
		&req.Message, &req.Status, &req.SentAt, &req.AcceptedAt, &req.LastMessageAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &req, nil
}

func (d *Database) HasProcessedProfile(profileID string) bool {
	var count int
	query := `SELECT COUNT(*) FROM connection_requests WHERE profile_id = ?`
	d.db.QueryRow(query, profileID).Scan(&count)
	return count > 0
}

func (d *Database) GetAcceptedConnections() ([]*ConnectionRequest, error) {
	query := `
		SELECT id, profile_id, name, title, company, message, status, sent_at, accepted_at, last_message_at
		FROM connection_requests
		WHERE status = 'accepted' AND last_message_at IS NULL
		ORDER BY accepted_at DESC
	`

	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var connections []*ConnectionRequest
	for rows.Next() {
		var req ConnectionRequest
		err := rows.Scan(
			&req.ID, &req.ProfileID, &req.Name, &req.Title, &req.Company,
			&req.Message, &req.Status, &req.SentAt, &req.AcceptedAt, &req.LastMessageAt,
		)
		if err != nil {
			return nil, err
		}
		connections = append(connections, &req)
	}

	return connections, nil
}

func (d *Database) SaveMessage(msg *Message) error {
	query := `INSERT INTO messages (profile_id, content) VALUES (?, ?)`
	_, err := d.db.Exec(query, msg.ProfileID, msg.Content)

	if err == nil {
		now := time.Now()
		updateQuery := `UPDATE connection_requests SET last_message_at = ? WHERE profile_id = ?`
		d.db.Exec(updateQuery, now, msg.ProfileID)
	}

	return err
}

func (d *Database) GetTodayStats() (*DailyStats, error) {
	today := time.Now().Format("2006-01-02")

	query := `
		SELECT date, connections_sent, messages_sent, connections_limit, messages_limit
		FROM daily_stats
		WHERE date = ?
	`

	var stats DailyStats
	err := d.db.QueryRow(query, today).Scan(
		&stats.Date, &stats.ConnectionsSent, &stats.MessagesSent,
		&stats.ConnectionsLimit, &stats.MessagesLimit,
	)

	if err == sql.ErrNoRows {
		return &DailyStats{
			Date:             today,
			ConnectionsSent:  0,
			MessagesSent:     0,
			ConnectionsLimit: 50,
			MessagesLimit:    30,
		}, nil
	}

	if err != nil {
		return nil, err
	}

	return &stats, nil
}

func (d *Database) IncrementConnectionCount() error {
	today := time.Now().Format("2006-01-02")

	query := `
		INSERT INTO daily_stats (date, connections_sent, connections_limit)
		VALUES (?, 1, 50)
		ON CONFLICT(date) DO UPDATE SET connections_sent = connections_sent + 1
	`

	_, err := d.db.Exec(query, today)
	return err
}

func (d *Database) IncrementMessageCount() error {
	today := time.Now().Format("2006-01-02")

	query := `
		INSERT INTO daily_stats (date, messages_sent, messages_limit)
		VALUES (?, 1, 30)
		ON CONFLICT(date) DO UPDATE SET messages_sent = messages_sent + 1
	`

	_, err := d.db.Exec(query, today)
	return err
}

func (d *Database) SaveSessionState(key, value string) error {
	query := `
		INSERT INTO session_state (key, value, updated_at)
		VALUES (?, ?, CURRENT_TIMESTAMP)
		ON CONFLICT(key) DO UPDATE SET value = excluded.value, updated_at = CURRENT_TIMESTAMP
	`

	_, err := d.db.Exec(query, key, value)
	return err
}

func (d *Database) GetSessionState(key string) (string, error) {
	var value string
	query := `SELECT value FROM session_state WHERE key = ?`
	err := d.db.QueryRow(query, key).Scan(&value)

	if err == sql.ErrNoRows {
		return "", nil
	}

	return value, err
}

func (d *Database) Close() error {
	return d.db.Close()
}
