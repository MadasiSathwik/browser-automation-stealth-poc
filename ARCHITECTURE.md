# Architecture Documentation

This document provides an in-depth technical explanation of the browser automation POC's architecture, design decisions, and implementation details.

## Table of Contents

1. [System Overview](#system-overview)
2. [Package Architecture](#package-architecture)
3. [Design Patterns](#design-patterns)
4. [Stealth Implementation](#stealth-implementation)
5. [State Management](#state-management)
6. [Error Handling Strategy](#error-handling-strategy)
7. [Performance Considerations](#performance-considerations)
8. [Security Model](#security-model)

---

## System Overview

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                         cmd/main.go                          │
│                    (Entry Point & Orchestration)             │
└───────────────────┬─────────────────────────────────────────┘
                    │
        ┌───────────┼───────────┐
        │           │           │
        ▼           ▼           ▼
   ┌────────┐  ┌────────┐  ┌─────────┐
   │ Config │  │ Logger │  │ Storage │
   │        │  │        │  │ (SQLite)│
   └────────┘  └────────┘  └─────────┘
        │           │           │
        └───────────┼───────────┘
                    │
        ┌───────────┼───────────┐
        │           │           │
        ▼           ▼           ▼
   ┌────────┐  ┌────────┐  ┌──────────┐
   │  Auth  │  │ Search │  │Connection│
   │Service │  │Service │  │  Service │
   └────┬───┘  └───┬────┘  └────┬─────┘
        │          │            │
        └──────────┼────────────┘
                   │
                   ▼
          ┌────────────────┐
          │ Stealth Layer  │
          │ ┌────────────┐ │
          │ │   Mouse    │ │
          │ │   Typing   │ │
          │ │   Timing   │ │
          │ │Fingerprint │ │
          │ └────────────┘ │
          └────────────────┘
                   │
                   ▼
          ┌────────────────┐
          │   Rod/CDP      │
          │   (Browser)    │
          └────────────────┘
```

### Component Responsibilities

| Component | Responsibility | Key Files |
|-----------|---------------|-----------|
| **cmd** | Application entry, flag parsing, orchestration | main.go |
| **config** | Configuration loading, validation, defaults | config.go |
| **logger** | Structured logging with levels | logger.go |
| **storage** | SQLite persistence, state management | sqlite.go, state.go |
| **auth** | Login, session management, cookie handling | login.go, session.go |
| **search** | Profile discovery, pagination, extraction | search.go, pagination.go |
| **connections** | Connection requests, rate limiting | connect.go, limits.go |
| **messaging** | Template rendering, message delivery | templates.go, messenger.go |
| **stealth** | Anti-detection, human simulation | mouse.go, typing.go, timing.go, fingerprint.go |

---

## Package Architecture

### Dependency Graph

```
cmd
 ├─> config (no deps)
 ├─> logger (no deps)
 ├─> storage (sqlite3)
 ├─> auth
 │    ├─> config
 │    ├─> logger
 │    ├─> storage
 │    └─> stealth
 ├─> search
 │    ├─> config
 │    ├─> logger
 │    ├─> storage
 │    └─> stealth
 ├─> connections
 │    ├─> config
 │    ├─> logger
 │    ├─> storage
 │    ├─> search
 │    └─> stealth
 └─> messaging
      ├─> config
      ├─> logger
      ├─> storage
      └─> stealth
```

### Design Principles

1. **Layered Architecture**: Clear separation between presentation (cmd), business logic (services), and data (storage)

2. **Dependency Injection**: Services receive dependencies through constructors
   ```go
   func NewService(cfg *config.Config, db *storage.Database, log *logger.Logger) *Service
   ```

3. **Interface Segregation**: Small, focused interfaces
   ```go
   type RateLimiter interface {
       CanSendConnection() bool
       CanSendMessage() bool
   }
   ```

4. **Single Responsibility**: Each package has one reason to change
   - `auth` - authentication changes
   - `search` - search algorithm changes
   - `stealth` - anti-detection changes

---

## Design Patterns

### 1. Service Pattern

Each domain has a Service that encapsulates business logic:

```go
type Service struct {
    cfg    *config.Config
    db     *storage.Database
    log    *logger.Logger
    // domain-specific dependencies
}

func NewService(...) *Service {
    return &Service{...}
}

func (s *Service) DomainOperation(ctx context.Context, ...) error {
    // Implementation
}
```

**Benefits:**
- Centralized business logic
- Easy to test with mocks
- Clear API boundaries

### 2. Strategy Pattern (Stealth Techniques)

Different stealth behaviors are encapsulated in separate controllers:

```go
type MouseController struct {
    cfg *MouseMovementConfig
}

type TypingSimulator struct {
    cfg *TimingConfig
}

// Services use composition
type ConnectionService struct {
    mouse  *MouseController
    typing *TypingSimulator
}
```

**Benefits:**
- Interchangeable strategies
- Easy to add new techniques
- Configurable behavior

### 3. Template Method Pattern (Messaging)

Message templates with variable substitution:

```go
type Template struct {
    Name      string
    Content   string   // "Hi {{name}}, ..."
    Variables []string
}

func (te *TemplateEngine) RenderTemplate(name string, vars map[string]string) (string, error)
```

**Benefits:**
- Reusable message templates
- Type-safe variable substitution
- Easy to add new templates

### 4. Repository Pattern (Storage)

Database operations abstracted behind clean interface:

```go
type Database struct {
    db *sql.DB
}

func (d *Database) SaveConnectionRequest(req *ConnectionRequest) error
func (d *Database) GetAcceptedConnections() ([]*ConnectionRequest, error)
```

**Benefits:**
- Decoupled from SQL details
- Easy to swap storage backends
- Mockable for testing

---

## Stealth Implementation

### Mouse Movement: Bézier Curves

**Algorithm: Cubic Bézier Curve**

```
B(t) = (1-t)³P₀ + 3(1-t)²tP₁ + 3(1-t)t²P₂ + t³P₃

Where:
- P₀ = start point
- P₁, P₂ = control points (randomized)
- P₃ = end point
- t ∈ [0, 1]
```

**Implementation:**

```go
func (mc *MouseController) cubicBezier(t, x0, y0, x1, y1, x2, y2, x3, y3 float64) Point {
    u := 1 - t
    tt := t * t
    uu := u * u
    uuu := uu * u
    ttt := tt * t

    x := uuu*x0 + 3*uu*t*x1 + 3*u*tt*x2 + ttt*x3
    y := uuu*y0 + 3*uu*t*y1 + 3*u*tt*y2 + ttt*y3

    return Point{X: x, Y: y}
}
```

**Control Point Randomization:**

```go
// Add randomness to control points for natural curves
x1 := x0 + (x3-x0)*0.3 + randomFloat(-50, 50)
y1 := y0 + (y3-y0)*0.3 + randomFloat(-50, 50)
x2 := x0 + (x3-x0)*0.7 + randomFloat(-50, 50)
y2 := y0 + (y3-y0)*0.7 + randomFloat(-50, 50)
```

**Velocity Profile:**

```go
func (mc *MouseController) calculateVelocity(current, total int) float64 {
    baseVelocity := 2.0
    progress := float64(current) / float64(total)

    // Ease-in: accelerate at start
    if progress < 0.2 {
        baseVelocity *= (1.0 + progress*2)
    }
    // Ease-out: decelerate at end
    else if progress > 0.8 {
        baseVelocity *= (1.0 + (1.0-progress)*2)
    }
    // Maintain speed in middle
    else {
        baseVelocity *= 2.5
    }

    // Add variance
    variance := 1.0 + randomFloat(-0.3, 0.3)
    return baseVelocity * variance
}
```

### Typing Simulation: Typos & Corrections

**Typo Injection Algorithm:**

1. Select random position in text (15% probability)
2. Type wrong character
3. Pause (recognition delay: 200-600ms)
4. Press backspace
5. Continue with correct character

```go
func (ts *TypingSimulator) typeWithTypo(element *rod.Element, text string) error {
    typoPosition := 1 + random.Intn(len(text)-2)

    for i, char := range text {
        if i == typoPosition {
            // Type wrong character
            element.Input(string(typoChar))
            time.Sleep(getTypingDelay(typoChar))

            // Recognition delay
            time.Sleep(200-400ms)

            // Correction
            element.Type(rod.KeyBackspace)
            time.Sleep(100ms)
        }

        // Type correct character
        element.Input(string(char))
        time.Sleep(getTypingDelay(char))
    }
}
```

**Key-Specific Timing:**

```go
func getTypingDelay(char rune) time.Duration {
    baseDelay := 45-120ms

    // Home row keys: faster
    if char in "asdfjkl;" {
        baseDelay *= 0.8
    }

    // Awkward keys: slower
    if char in "qwertyzxcvb" {
        baseDelay *= 1.2
    }

    // Uppercase: slower (shift key)
    if isUppercase(char) {
        baseDelay *= 1.3
    }

    return baseDelay
}
```

### Fingerprint Masking

**navigator.webdriver Removal:**

```javascript
Object.defineProperty(navigator, 'webdriver', {
    get: () => undefined
});
```

**Plugin Array Spoofing:**

```javascript
Object.defineProperty(navigator, 'plugins', {
    get: () => [1, 2, 3, 4, 5]  // Non-empty array
});
```

**Chrome Object Addition:**

```javascript
window.navigator.chrome = {
    runtime: {}
};
```

**Permissions API Patching:**

```javascript
const originalQuery = window.navigator.permissions.query;
window.navigator.permissions.query = (parameters) => (
    parameters.name === 'notifications' ?
        Promise.resolve({ state: Notification.permission }) :
        originalQuery(parameters)
);
```

---

## State Management

### Database Schema Design

**Entity Relationship:**

```
connection_requests
    ├── 1:N messages
    └── linked to daily_stats (by date)

session_state
    └── independent key-value store
```

**Connection Request Lifecycle:**

```
pending → [user accepts] → accepted → [message sent] → messaged
        → [user rejects] → rejected
        → [timeout]      → expired
```

**State Transitions:**

```go
func (db *Database) SaveConnectionRequest(req *ConnectionRequest) error {
    query := `
        INSERT INTO connection_requests (...)
        VALUES (...)
        ON CONFLICT(profile_id) DO UPDATE SET
            status = excluded.status,
            accepted_at = excluded.accepted_at
    `
}
```

### Session Persistence

**Cookie Storage:**

```go
type SessionState struct {
    Cookies      []Cookie
    LastActivity time.Time
    Valid        bool
}

// Serialized as JSON in session_state table
```

**Session Validation:**

```go
func (d *Database) LoadCookies() ([]Cookie, error) {
    state := loadSessionState()

    // Check expiration
    if time.Since(state.LastActivity) > 7*24*time.Hour {
        return nil, fmt.Errorf("session expired")
    }

    // Check validity flag
    if !state.Valid {
        return nil, fmt.Errorf("session invalidated")
    }

    return state.Cookies, nil
}
```

---

## Error Handling Strategy

### Error Hierarchy

```
Level 1: Fatal Errors (exit immediately)
  - Configuration validation failure
  - Database initialization failure
  - Critical dependency missing

Level 2: Recoverable Errors (retry with backoff)
  - Network timeouts
  - Element not found (temporary)
  - Rate limit exceeded

Level 3: Warnings (log and continue)
  - Failed to hover element
  - Optional feature unavailable
  - Performance degradation
```

### Retry Logic

**Exponential Backoff:**

```go
func (tc *TimingController) ExponentialBackoff(attempt int, maxDelay time.Duration) time.Duration {
    base := 1 * time.Second
    delay := time.Duration(1<<uint(attempt)) * base  // 1s, 2s, 4s, 8s...

    if delay > maxDelay {
        delay = maxDelay
    }

    jitter := random(0, 1000ms)
    return delay + jitter
}
```

**Usage Example:**

```go
func (s *Service) connectWithRetry(profile *Profile) error {
    maxAttempts := 3

    for attempt := 0; attempt < maxAttempts; attempt++ {
        err := s.connect(profile)
        if err == nil {
            return nil
        }

        if !isRetryable(err) {
            return err
        }

        delay := s.timing.ExponentialBackoff(attempt, 30*time.Second)
        s.log.Warnf("Attempt %d failed, retrying in %v", attempt+1, delay)
        time.Sleep(delay)
    }

    return fmt.Errorf("max retries exceeded")
}
```

---

## Performance Considerations

### Resource Management

**Browser Lifecycle:**

```go
browser := rod.New().ControlURL(url).MustConnect()
defer browser.Close()  // Always cleanup

page := browser.MustPage()
defer page.Close()
```

**Database Connection Pooling:**

```go
db, _ := sql.Open("sqlite3", path)
db.SetMaxOpenConns(1)  // SQLite: single writer
db.SetMaxIdleConns(1)
defer db.Close()
```

### Memory Optimization

**Streaming Results:**

```go
func (s *Service) SearchProfiles(...) ([]*Profile, error) {
    var profiles []*Profile

    for pageNum := 1; pageNum <= maxPages; pageNum++ {
        pageProfiles := extractFromPage(page)

        // Filter already processed
        for _, p := range pageProfiles {
            if !s.db.HasProcessedProfile(p.ID) {
                profiles = append(profiles, p)
            }
        }

        // Process immediately instead of accumulating
    }

    return profiles, nil
}
```

### Concurrency Considerations

**Why Single-Threaded?**

1. Browser automation is inherently sequential
2. Human behavior is not parallel
3. Rate limiting requires serialization
4. Simpler error handling

**Future: Parallel Sessions:**

```go
// Potential enhancement
func runParallelSessions(numSessions int) {
    var wg sync.WaitGroup

    for i := 0; i < numSessions; i++ {
        wg.Add(1)
        go func(sessionID int) {
            defer wg.Done()
            browser := createBrowser(sessionID)
            runAutomation(browser)
        }(i)
    }

    wg.Wait()
}
```

---

## Security Model

### Threat Model

**Assumptions:**

1. Configuration files are trusted
2. Environment variables are secure
3. Database is not tampered with
4. Local system is not compromised

**Protections:**

1. **Safe Mode Validation**
   ```go
   func (c *Config) Validate() error {
       if c.SafeMode && !isTestURL(c.Auth.LoginURL) {
           return fmt.Errorf("safe mode violation")
       }
   }
   ```

2. **Credential Isolation**
   ```go
   // Never log credentials
   log.Info("Attempting login")  // ✅
   log.Infof("Login: %s / %s", email, password)  // ❌
   ```

3. **SQL Injection Prevention**
   ```go
   // Always use parameterized queries
   query := `SELECT * FROM users WHERE id = ?`
   db.Query(query, userID)  // ✅

   query := fmt.Sprintf(`SELECT * FROM users WHERE id = %s`, userID)  // ❌
   ```

4. **Session Encryption** (Future Enhancement)
   ```go
   // TODO: Encrypt cookies at rest
   encryptedCookies := encrypt(cookies, key)
   db.SaveSessionState("session", encryptedCookies)
   ```

### Safety Mechanisms

1. **URL Whitelist:**
   ```go
   testDomains := []string{"localhost", "127.0.0.1", "example.com", "test.com", "mock", "demo"}
   ```

2. **Rate Limit Enforcement:**
   ```go
   if stats.ConnectionsSent >= cfg.Limits.DailyConnections {
       return fmt.Errorf("daily limit reached")
   }
   ```

3. **Demo Mode:**
   ```go
   if cfg.DemoMode {
       log.Info("[DEMO] Would send connection request")
       return nil
   }
   ```

---

## Extension Points

### Adding New Stealth Techniques

1. Create new file in `stealth/` package
2. Implement as configurable controller
3. Add configuration to `config.yaml`
4. Inject into services via constructor

Example:

```go
// stealth/captcha.go
type CaptchaDetector struct {
    cfg *config.CaptchaConfig
}

func (cd *CaptchaDetector) Detect(page *rod.Page) (bool, error) {
    // Implementation
}

// connections/connect.go
type Service struct {
    captcha *stealth.CaptchaDetector
}
```

### Adding New Message Templates

```go
// messaging/templates.go
var CustomTemplates = []Template{
    {
        Name: "follow_up_event",
        Content: "Hi {{name}}, great meeting you at {{event}}!",
        Variables: []string{"name", "event"},
    },
}

func init() {
    DefaultTemplates = append(DefaultTemplates, CustomTemplates...)
}
```

---

## Testing Strategy

### Unit Testing

```go
// stealth/mouse_test.go
func TestBezierCurveGeneration(t *testing.T) {
    mc := NewMouseController(&config.MouseMovementConfig{
        BezierCurves: true,
    })

    points := mc.generateBezierPath(0, 0, 100, 100)

    assert.Greater(t, len(points), 20)
    assert.Equal(t, 0.0, points[0].X)
    assert.Equal(t, 100.0, points[len(points)-1].X)
}
```

### Integration Testing

```go
// auth/session_test.go
func TestSessionPersistence(t *testing.T) {
    db := storage.NewTestDatabase()
    service := NewService(cfg, db, log)

    // Save session
    service.SaveSession(page)

    // Load session
    valid, err := service.LoadSession(page)
    assert.True(t, valid)
    assert.NoError(t, err)
}
```

### End-to-End Testing

```bash
# Run against mock site
go run cmd/main.go --config test-config.yaml
```

---

## Deployment Considerations

### Configuration Management

```
Development:  config.dev.yaml   (headless: false, demo: true)
Testing:      config.test.yaml  (mock URLs, safe_mode: true)
Production:   config.prod.yaml  (real URLs, strict limits)
```

### Monitoring

**Key Metrics:**

- Connections sent per day
- Message delivery rate
- Error rates by type
- Average action duration
- Session validity rate

**Log Aggregation:**

```go
// Use structured logging for parsing
log.WithFields(logrus.Fields{
    "action": "connection_request",
    "profile_id": profile.ID,
    "duration_ms": duration.Milliseconds(),
}).Info("Connection sent")
```

---

## Conclusion

This architecture demonstrates:

✅ **Clean Code**: Modular, testable, maintainable
✅ **Advanced Techniques**: Bézier curves, typo simulation, fingerprint masking
✅ **Safety First**: Multiple layers of protection
✅ **Production-Ready Patterns**: Proper error handling, state management, logging
✅ **Extensibility**: Easy to add features without breaking existing code

The codebase serves as both a technical demonstration and a learning resource for advanced browser automation engineering.
