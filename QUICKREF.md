# Quick Reference Guide

## Common Commands

### Build & Run

```bash
# Download dependencies
go mod download

# Build binary
make build
# or
go build -o browser-automation cmd/main.go

# Run demo mode (safe, prints actions)
make demo
# or
go run cmd/main.go --demo

# Run with default config
make run
# or
go run cmd/main.go --config config.yaml

# Run with safe mode
make safe
# or
go run cmd/main.go --safe --config config.yaml

# Run with custom config
go run cmd/main.go --config my-config.yaml

# Run with debug logging
DEBUG=true go run cmd/main.go --demo
```

### Testing

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests with coverage
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run specific package tests
go test ./stealth
go test ./storage
go test ./auth

# Run specific test
go test ./stealth -run TestGenerateBezierPath

# Run benchmarks
go test -bench=. ./stealth
go test -bench=BenchmarkMouseMovement ./stealth

# Run with race detector
go test -race ./...
```

### Development

```bash
# Format code
make fmt
# or
go fmt ./...

# Tidy dependencies
go mod tidy

# Verify dependencies
go mod verify

# List dependencies
go list -m all

# Clean build artifacts
make clean
# or
rm -f browser-automation automation.db

# Rebuild from scratch
make rebuild
# or
make clean && make build
```

### Database

```bash
# Open database
sqlite3 automation.db

# View tables
sqlite3 automation.db ".tables"

# View schema
sqlite3 automation.db ".schema"

# Query connections
sqlite3 automation.db "SELECT * FROM connection_requests;"

# View daily stats
sqlite3 automation.db "SELECT * FROM daily_stats;"

# Clear session
sqlite3 automation.db "DELETE FROM session_state;"

# Reset database
rm automation.db
```

### Configuration

```bash
# View current config
cat config.yaml

# Validate config
go run cmd/main.go --config config.yaml 2>&1 | grep -i "validation\|error"

# Copy example env
cp .env.example .env

# Edit environment
vim .env
```

## File Locations

### Source Code

| File | Purpose |
|------|---------|
| `cmd/main.go` | Application entry point |
| `config/config.go` | Configuration structure |
| `stealth/mouse.go` | Mouse movement (Bézier) |
| `stealth/typing.go` | Typing simulation |
| `stealth/timing.go` | Delay management |
| `stealth/fingerprint.go` | Fingerprint masking |
| `auth/login.go` | Login automation |
| `auth/session.go` | Session management |
| `search/search.go` | Profile search |
| `connections/connect.go` | Connection logic |
| `messaging/messenger.go` | Message delivery |
| `storage/sqlite.go` | Database operations |

### Documentation

| File | Purpose |
|------|---------|
| `README.md` | Main documentation |
| `ARCHITECTURE.md` | Technical details |
| `DEMO.md` | Quick start guide |
| `TESTING.md` | Testing guide |
| `STEALTH_TECHNIQUES.md` | Stealth reference |
| `PROJECT_SUMMARY.md` | Overview |
| `QUICKREF.md` | This file |

### Configuration

| File | Purpose |
|------|---------|
| `config.yaml` | Main configuration |
| `.env.example` | Environment template |
| `.env` | Environment variables (create from example) |
| `go.mod` | Go dependencies |

## Configuration Quick Reference

### Enable/Disable Features

```yaml
# Safe mode (validates URLs)
safe_mode: true

# Demo mode (print only)
demo_mode: false

# Headless browser
browser:
  headless: false

# Mouse features
stealth:
  mouse_movement:
    enabled: true
    bezier_curves: true
    overshoot: true
    micro_corrections: true

# Random scrolling
stealth:
  random_scrolling: true

# Business hours only
stealth:
  business_hours_only: false
```

### Adjust Timing

```yaml
timing:
  min_delay: 2          # Min seconds between actions
  max_delay: 8          # Max seconds between actions
  between_actions: 30   # Seconds between major actions
  typing_min: 45        # Min ms per keystroke
  typing_max: 120       # Max ms per keystroke
  think_min: 500        # Min ms think delay
  think_max: 2000       # Max ms think delay
```

### Adjust Limits

```yaml
limits:
  daily_connections: 50    # Max connections per day
  hourly_connections: 10   # Max connections per hour
  daily_messages: 30       # Max messages per day
```

### Change Selectors

```yaml
auth:
  email_selector: "#email"
  password_selector: "#password"
  submit_selector: "button[type='submit']"

search:
  selectors:
    profile_card: ".profile-card"
    profile_name: ".profile-name"
    profile_title: ".profile-title"
```

## Environment Variables

```bash
# Authentication (required for non-demo mode)
LOGIN_EMAIL=test@example.com
LOGIN_PASSWORD=test-password

# Configuration overrides
LOGIN_URL=https://mock-site.example.com/login
SEARCH_QUERY=Software Engineer
DATABASE_PATH=automation.db

# Debug mode
DEBUG=true
```

## Makefile Targets

| Command | Description |
|---------|-------------|
| `make build` | Build binary |
| `make run` | Build and run |
| `make demo` | Run demo mode |
| `make safe` | Run with safe mode |
| `make deps` | Install dependencies |
| `make test` | Run tests |
| `make clean` | Remove artifacts |
| `make rebuild` | Clean and rebuild |
| `make fmt` | Format code |
| `make help` | Show help |

## Common Flags

| Flag | Description | Example |
|------|-------------|---------|
| `--config` | Config file path | `--config custom.yaml` |
| `--demo` | Enable demo mode | `--demo` |
| `--safe` | Enable safe mode | `--safe` |

## Troubleshooting

### Browser doesn't launch

```bash
# Check if Chrome/Chromium is installed
which chromium-browser
which google-chrome

# Install Chrome (Ubuntu)
sudo apt-get install chromium-browser

# Install Chrome (macOS)
brew install chromium
```

### Database locked

```bash
# Kill running processes
pkill browser-automation

# Remove journal file
rm automation.db-journal
```

### Import errors

```bash
# Download dependencies
go mod download

# Tidy dependencies
go mod tidy
```

### Configuration errors

```bash
# Validate YAML syntax
yamllint config.yaml

# Check for typos in config
cat config.yaml | grep -i "error\|invalid"
```

## Code Examples

### Add New Stealth Technique

```go
// stealth/newtech.go
package stealth

type NewTechnique struct {
    cfg *config.Config
}

func NewNewTechnique(cfg *config.Config) *NewTechnique {
    return &NewTechnique{cfg: cfg}
}

func (nt *NewTechnique) Apply(page *rod.Page) error {
    // Implementation
    return nil
}
```

### Add New Message Template

```go
// messaging/templates.go
var NewTemplate = Template{
    Name:    "custom_template",
    Content: "Hi {{name}}, custom message about {{topic}}",
    Variables: []string{"name", "topic"},
}

func init() {
    DefaultTemplates = append(DefaultTemplates, NewTemplate)
}
```

### Add New Configuration

```yaml
# config.yaml
custom:
  feature_enabled: true
  setting_value: 42
```

```go
// config/config.go
type Config struct {
    // ... existing fields
    Custom CustomConfig `yaml:"custom"`
}

type CustomConfig struct {
    FeatureEnabled bool `yaml:"feature_enabled"`
    SettingValue   int  `yaml:"setting_value"`
}
```

## Performance Tips

### Faster Execution (Testing)

```yaml
timing:
  min_delay: 1
  max_delay: 2
  between_actions: 5

browser:
  headless: true
```

### Maximum Stealth (Production)

```yaml
timing:
  min_delay: 3
  max_delay: 10
  between_actions: 45

stealth:
  business_hours_only: true
  mouse_movement:
    enabled: true
    bezier_curves: true
    overshoot: true
    micro_corrections: true
  random_scrolling: true

limits:
  daily_connections: 30
  hourly_connections: 5
```

## Useful Scripts

### Reset Everything

```bash
#!/bin/bash
make clean
rm -f automation.db automation.db-journal
rm -f .env
cp .env.example .env
make deps
make build
```

### Run Daily

```bash
#!/bin/bash
# daily-run.sh
LOG_FILE="automation-$(date +%Y%m%d).log"
./browser-automation --config config.yaml > "$LOG_FILE" 2>&1
```

### Monitor Database

```bash
#!/bin/bash
# monitor-db.sh
watch -n 5 'sqlite3 automation.db "SELECT date, connections_sent, messages_sent FROM daily_stats ORDER BY date DESC LIMIT 5;"'
```

## Key Functions

### Mouse Movement

```go
mouseController.MoveToElement(page, element)
mouseController.HoverElement(page, element)
mouseController.ClickElement(page, element)
mouseController.IdleWander(page)
```

### Typing

```go
typingSimulator.TypeIntoElement(element, text, humanLike)
typingSimulator.TypeWithThinkDelay(element, text)
typingSimulator.PasteText(element, text, humanLike)
```

### Timing

```go
timingController.RandomDelay()
timingController.ThinkDelay()
timingController.BetweenActionsDelay()
timingController.ExponentialBackoff(attempt, maxDelay)
```

### Fingerprint

```go
fingerprintManager.ApplyStealthSettings(page)
fingerprintManager.RandomScroll(page)
fingerprintManager.SimulateReading(page, duration)
```

## Resources

### Documentation

- Main: `README.md`
- Architecture: `ARCHITECTURE.md`
- Demo: `DEMO.md`
- Testing: `TESTING.md`
- Stealth: `STEALTH_TECHNIQUES.md`

### External

- [Go Docs](https://golang.org/doc/)
- [Rod Framework](https://go-rod.github.io/)
- [SQLite](https://www.sqlite.org/docs.html)

## Getting Help

### Check Logs

```bash
# Enable debug mode
DEBUG=true go run cmd/main.go --demo

# Save logs to file
go run cmd/main.go --demo 2>&1 | tee output.log

# Search logs
grep -i "error\|warn" output.log
```

### Common Issues

1. **"go: command not found"**
   - Install Go from https://golang.org/dl/

2. **"cannot find package"**
   - Run `go mod download`

3. **"element not found"**
   - Check selectors in config.yaml
   - Verify page structure

4. **"database locked"**
   - Close other instances
   - Remove `automation.db-journal`

5. **"rate limit reached"**
   - Check `daily_stats` table
   - Wait for next day or increase limits

---

For more details, see the full documentation in README.md
