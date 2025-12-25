# Browser Automation POC - Advanced Stealth Techniques

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-Educational-blue.svg)](LICENSE)
[![Status](https://img.shields.io/badge/Status-POC-yellow.svg)]()

> **⚠️ EDUCATIONAL PURPOSE ONLY**
> This project is a proof-of-concept demonstrating advanced browser automation engineering, human-like behavior simulation, and anti-detection techniques. It is **NOT** intended for production use or real-world automation of any platform.

## 🎯 Project Overview

This is a sophisticated browser automation system built in Go that showcases:

- **Advanced Stealth Techniques**: 8+ anti-detection methods
- **Human-Like Behavior**: Realistic mouse movements, typing patterns, and timing
- **Clean Architecture**: Modular, testable, idiomatic Go code
- **State Persistence**: SQLite-based session and activity tracking
- **Safety First**: Built-in safe mode and demo mode for evaluation

### Key Features

✅ Human-like mouse movement with Bézier curves
✅ Realistic typing simulation with typos and corrections
✅ Randomized timing and behavioral patterns
✅ Browser fingerprint masking
✅ Session persistence and resume capability
✅ Rate limiting and quota enforcement
✅ Configurable DOM selectors (no hard-coded sites)
✅ Comprehensive logging and error handling
✅ Business hours scheduling

---

## 🏗️ Architecture

### Package Structure

```
/cmd/               - Application entry point
/config/            - Configuration management
/auth/              - Authentication & session handling
/search/            - Profile search and pagination
/connections/       - Connection request logic
/messaging/         - Message templates and delivery
/stealth/           - Anti-detection techniques
  ├── mouse.go      - Bézier curve mouse movement
  ├── typing.go     - Human-like typing simulation
  ├── timing.go     - Randomized delays and scheduling
  └── fingerprint.go - Browser fingerprint masking
/storage/           - SQLite persistence layer
/logger/            - Structured logging
```

### Design Principles

1. **Separation of Concerns**: Each package has a single, well-defined responsibility
2. **Dependency Injection**: Services are injected for testability
3. **Configuration Over Code**: All behavior is configurable via YAML
4. **Fail-Safe Design**: Multiple safety mechanisms prevent accidental misuse

---

## 🕵️ Stealth Techniques Implemented

### Mandatory (3/3)

#### 1. Human-Like Mouse Movement ✅
- **Bézier Curves**: Smooth, natural cursor paths
- **Variable Velocity**: Speed varies throughout movement (acceleration/deceleration)
- **Micro-Corrections**: Small adjustments before final position
- **Overshoot**: Occasional overshoot with correction
- **Implementation**: `stealth/mouse.go`

```go
// Example: Bézier curve generation with control points
func (mc *MouseController) generateBezierPath(x0, y0, x3, y3 float64) []Point {
    // Randomized control points for natural curves
    x1 := x0 + (x3-x0)*0.3 + mc.randomFloat(-50, 50)
    y1 := y0 + (y3-y0)*0.3 + mc.randomFloat(-50, 50)
    // ... cubic Bézier calculation
}
```

#### 2. Randomized Timing ✅
- **Think Delays**: Pauses before actions (500-2000ms)
- **Action Jitter**: Randomized delays between operations
- **Scroll Pauses**: Variable delays during scrolling
- **Exponential Backoff**: Intelligent retry timing
- **Implementation**: `stealth/timing.go`

```go
// Dynamic delay calculation based on action type
func (tc *TimingController) HumanizedDelay(action string) time.Duration {
    baseDelays := map[string]int{
        "click": 300, "scroll": 200, "read": 1000,
    }
    // Add 40% variance
}
```

#### 3. Browser Fingerprint Masking ✅
- **Randomized Viewport**: Variable window sizes (1366x768, 1920x1080, etc.)
- **User-Agent Rotation**: Multiple realistic user agents
- **Navigator.webdriver Removal**: Hides automation detection
- **Plugin/Language Spoofing**: Appears as normal browser
- **Implementation**: `stealth/fingerprint.go` + `go-rod/stealth`

```go
// Hide webdriver property
Object.defineProperty(navigator, 'webdriver', {
    get: () => undefined
});
```

### Additional Techniques (5+/5)

#### 4. Random Scrolling Behavior ✅
- **Smooth Scroll**: CSS smooth scrolling
- **Partial Scroll**: Small incremental scrolls
- **Bounce Scroll**: Scroll down then slightly up (reading pattern)
- **Reading Scroll**: Multiple pauses simulating content reading

#### 5. Realistic Typing with Typos + Corrections ✅
- **Variable Speed**: 45-120ms per character
- **Home Row Bias**: Faster for commonly typed keys
- **Typo Simulation**: 15% chance of typos with backspace correction
- **Think Pauses**: Random 200-800ms pauses during typing

#### 6. Mouse Hovering & Idle Wandering ✅
- **Element Hovering**: Hover before clicking (100-400ms)
- **Idle Movement**: Random mouse movement during waits
- **Realistic Patterns**: Moves to natural screen positions

#### 7. Business-Hour Activity Scheduling ✅
- **Weekday Detection**: Only Mon-Fri operation
- **Time Windows**: 9 AM - 5 PM local time
- **Auto-Pause**: Waits until next business hours
- **Implementation**: `config/config.go`, `stealth/timing.go`

#### 8. Rate Limiting & Cooldown Enforcement ✅
- **Daily Limits**: 50 connections, 30 messages
- **Hourly Limits**: 10 connections per hour
- **Progressive Delays**: Longer delays after many actions
- **Persistent Tracking**: SQLite-based quota management

---

## 🔒 Safety Features

### SAFE_MODE (Default: Enabled)

When `safe_mode: true` in config:

- ✅ Validates URLs to ensure they're test/mock domains
- ✅ Warns if real-world sites are detected
- ✅ Uses mock credentials
- ✅ Blocks execution if configuration seems dangerous

```yaml
safe_mode: true  # ALWAYS keep this enabled for demonstrations
```

### DEMO_MODE

Run with `--demo` flag:

```bash
go run cmd/main.go --demo
```

**Demo mode prints actions without executing them:**

```
[ACTION] Navigating to login page
  - URL: https://mock-professional-network.example.com/login
  - Applying human-like mouse movement (Bézier curve)
[ACTION] Entering credentials
  - Typing email with realistic delays (45-120ms per char)
  - Simulated typo: 'usre@' -> backspace -> 'user@'
[ACTION] Searching for profiles
  - Found 15 profiles on page 1
[ACTION] Processing profile 1: Alice Johnson
  - Connection request recorded in database
  - Cooldown: 47s (randomized)
```

---

## 🚀 Setup & Installation

### Prerequisites

- **Go 1.21+**: [Install Go](https://golang.org/doc/install)
- **SQLite**: Included with Go driver
- **Chrome/Chromium**: Installed on system (for Rod)

### Installation

1. **Clone or download this project**

```bash
cd browser-automation
```

2. **Install dependencies**

```bash
go mod download
```

3. **Configure environment**

```bash
cp .env.example .env
# Edit .env with your test credentials (for mock sites only!)
```

4. **Review configuration**

```bash
# Edit config.yaml to customize behavior
# IMPORTANT: Keep safe_mode: true for demonstrations
vim config.yaml
```

5. **Build the application**

```bash
go build -o browser-automation cmd/main.go
```

---

## 📖 Usage

### Demo Mode (Recommended for Evaluation)

```bash
# Run in demo mode - prints actions without execution
./browser-automation --demo

# Output:
# [ACTION] Initializing browser with stealth configuration
# [ACTION] Navigating to login page
# [ACTION] Searching for profiles
# ...
```

### Safe Mode (Test Environment)

```bash
# Run against mock/test sites with safety checks
./browser-automation --config config.yaml

# Ensure config.yaml has:
# safe_mode: true
# auth.login_url: https://mock-site.example.com
```

### Configuration Options

```bash
# Custom config file
./browser-automation --config my-config.yaml

# Enable demo mode
./browser-automation --demo

# Force safe mode (redundant if already in config)
./browser-automation --safe
```

### Environment Variables

```bash
# Enable debug logging
DEBUG=true ./browser-automation

# Override login credentials
LOGIN_EMAIL=test@example.com LOGIN_PASSWORD=test123 ./browser-automation

# Override database path
DATABASE_PATH=./data/automation.db ./browser-automation
```

---

## 🔧 Configuration

### config.yaml Structure

```yaml
safe_mode: true              # Safety enforcement
demo_mode: false             # Print-only mode

browser:
  headless: false            # Show browser window
  no_sandbox: false          # Disable sandbox (use cautiously)

auth:
  login_url: "https://mock-site.example.com/login"
  email_selector: "#email"
  password_selector: "#password"

search:
  base_url: "https://mock-site.example.com/search"
  query: "Software Engineer"
  max_pages: 5

  selectors:                 # Configurable DOM selectors
    profile_card: ".profile-card"
    profile_name: ".profile-name"
    # ... more selectors

limits:
  daily_connections: 50
  hourly_connections: 10
  daily_messages: 30

timing:
  between_actions: 30        # Seconds between major actions
  typing_min: 45             # Milliseconds per keystroke
  typing_max: 120
  think_min: 500             # Think delay before actions
  think_max: 2000

stealth:
  mouse_movement:
    enabled: true
    bezier_curves: true
    overshoot: true
    micro_corrections: true
    velocity_variance: 0.3

  random_scrolling: true
  business_hours_only: false

  user_agent_rotation:
    - "Mozilla/5.0 (Windows NT 10.0; Win64; x64)..."
    - "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)..."

  viewport_sizes:
    - width: 1920
      height: 1080
    - width: 1366
      height: 768
```

---

## 💾 State Persistence

### Database Schema

The system uses SQLite to track:

- **Connection Requests**: Profiles, messages, status
- **Messages**: Follow-up messages sent
- **Daily Stats**: Quota tracking
- **Session State**: Cookies, session validity

```sql
CREATE TABLE connection_requests (
    id INTEGER PRIMARY KEY,
    profile_id TEXT UNIQUE,
    name TEXT,
    title TEXT,
    company TEXT,
    message TEXT,
    status TEXT,           -- 'pending', 'accepted', 'rejected'
    sent_at TIMESTAMP,
    accepted_at TIMESTAMP,
    last_message_at TIMESTAMP
);

CREATE TABLE daily_stats (
    date TEXT PRIMARY KEY,
    connections_sent INTEGER,
    messages_sent INTEGER,
    connections_limit INTEGER,
    messages_limit INTEGER
);
```

### Resume Capability

The system can resume after interruption:

1. **Session Cookies**: Saved and restored automatically
2. **Processed Profiles**: Skips already-contacted profiles
3. **Daily Quotas**: Persists across restarts
4. **Message History**: Prevents duplicate messages

---

## 🧪 Testing

### Manual Testing

1. **Demo Mode Test**:
   ```bash
   ./browser-automation --demo
   # Verify all actions are printed correctly
   ```

2. **Mock Site Test**:
   ```bash
   # Set up a local mock site or use example.com
   # Run with safe_mode enabled
   ./browser-automation
   ```

3. **Database Inspection**:
   ```bash
   sqlite3 automation.db
   SELECT * FROM connection_requests;
   SELECT * FROM daily_stats;
   ```

### Safety Checklist

Before running:

- [ ] `safe_mode: true` in config.yaml
- [ ] Login URL is a test/mock domain
- [ ] `.env` uses test credentials
- [ ] Rate limits are configured
- [ ] Demo mode tested first

---

## 📊 Monitoring & Logging

### Log Levels

- **DEBUG**: Detailed execution flow (set `DEBUG=true`)
- **INFO**: High-level actions and progress
- **WARN**: Recoverable errors, safety warnings
- **ERROR**: Critical failures

### Example Logs

```
INFO[2024-01-15 10:23:45] Starting Browser Automation POC
INFO[2024-01-15 10:23:45] Demo Mode: false | Safe Mode: true
WARN[2024-01-15 10:23:45] SAFE MODE: Real-world execution is disabled
INFO[2024-01-15 10:23:46] Navigating to login page
DEBUG[2024-01-15 10:23:47] Entering email address
INFO[2024-01-15 10:23:52] Login successful
INFO[2024-01-15 10:23:53] Starting profile search: Software Engineer
INFO[2024-01-15 10:23:58] Found 15 profiles to process
INFO[2024-01-15 10:24:10] Connection request sent to Alice Johnson
```

---

## 🎓 Educational Value

### What This Project Demonstrates

1. **Go Architecture**: Clean, modular code structure
2. **Browser Automation**: Rod library usage and best practices
3. **Stealth Techniques**: Real-world anti-detection methods
4. **Human Simulation**: Behavioral pattern replication
5. **State Management**: Persistent storage and recovery
6. **Error Handling**: Robust failure management
7. **Configuration**: Flexible, environment-based setup

### Learning Outcomes

After studying this codebase, you'll understand:

- How sophisticated automation systems are architected
- Why and how anti-detection techniques work
- Best practices for Go service design
- Browser automation with Rod/CDP
- State persistence patterns
- Rate limiting strategies

---

## ⚖️ Legal & Ethical Considerations

### This Project is Educational

**DO NOT use this code to:**
- Violate terms of service of any platform
- Automate real user accounts without permission
- Scrape data without authorization
- Bypass security measures for malicious purposes

**Appropriate uses:**
- Learning automation engineering
- Understanding anti-bot techniques for defense
- Building legitimate automation tools with permission
- Academic research and education

### Responsible Development

If you're building production automation:

1. ✅ Get explicit permission from the platform
2. ✅ Review and comply with Terms of Service
3. ✅ Implement proper authentication/authorization
4. ✅ Respect rate limits and API usage policies
5. ✅ Consider using official APIs instead

---

## 🛠️ Troubleshooting

### Common Issues

**1. Browser doesn't launch**
```bash
# Install Chrome/Chromium
# On Ubuntu:
sudo apt-get install chromium-browser

# On macOS:
brew install chromium
```

**2. "connect button not found"**
- Check that selectors in config.yaml match your test site
- Verify the page loaded correctly
- Enable non-headless mode to see what's happening

**3. Database locked**
```bash
# Close any other instances
pkill browser-automation
rm automation.db-journal
```

**4. Session expired**
```bash
# Clear saved session
rm automation.db
# Or:
sqlite3 automation.db "DELETE FROM session_state;"
```

---

## 🔮 Future Enhancements

Potential improvements for educational exploration:

- [ ] Machine learning-based behavior patterns
- [ ] CAPTCHA detection and alerting
- [ ] Proxy rotation support
- [ ] Distributed execution across multiple instances
- [ ] Advanced fingerprint randomization (Canvas, WebGL)
- [ ] Behavioral biometrics (mouse signature, typing rhythm)
- [ ] Natural language template generation
- [ ] Computer vision-based interaction

---

## 📚 References & Resources

### Browser Automation
- [Rod Framework](https://github.com/go-rod/rod) - High-level Chrome DevTools Protocol library
- [Chrome DevTools Protocol](https://chromedevtools.github.io/devtools-protocol/)

### Anti-Detection
- [go-rod/stealth](https://github.com/go-rod/stealth) - Stealth plugin for Rod
- [Bézier Curves](https://en.wikipedia.org/wiki/B%C3%A9zier_curve) - Smooth curve mathematics

### Go Best Practices
- [Effective Go](https://golang.org/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)

---

## 📄 License

This project is provided for **educational purposes only**. Use responsibly and ethically.

---

## 👨‍💻 Author

Created as a demonstration of advanced browser automation engineering and Go architecture.

**Skills Demonstrated:**
- Go (Golang) - Clean architecture, concurrency, dependency injection
- Browser Automation - Rod, Chrome DevTools Protocol
- Stealth Techniques - Anti-detection, human simulation
- System Design - Modularity, scalability, maintainability
- Database Design - SQLite, state persistence
- Configuration Management - YAML, environment variables

---

## 🙏 Acknowledgments

- [go-rod](https://github.com/go-rod/rod) for the excellent browser automation library
- The Go community for idiomatic patterns and best practices
- Security researchers for anti-detection technique documentation

---

**⚠️ Final Reminder: This is a proof-of-concept for educational evaluation only. Always use automation responsibly and with explicit permission.**
