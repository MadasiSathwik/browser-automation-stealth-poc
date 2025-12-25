# Testing Guide

This document explains how to test the browser automation POC.

## Test Structure

```
project/
├── stealth/
│   └── mouse_test.go          # Unit tests for mouse movement
├── config_test.go             # Configuration validation tests
├── storage_test.go            # Database tests
└── integration_test.go        # End-to-end tests
```

## Running Tests

### All Tests

```bash
go test ./...
```

### Specific Package

```bash
go test ./stealth
go test ./storage
go test ./auth
```

### Verbose Output

```bash
go test -v ./...
```

### With Coverage

```bash
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Benchmarks

```bash
go test -bench=. ./stealth
```

## Test Categories

### 1. Unit Tests

Test individual functions in isolation.

**Example: Mouse Movement**

```go
func TestGenerateBezierPath(t *testing.T) {
    mc := NewMouseController(&config.MouseMovementConfig{
        BezierCurves: true,
    })

    points := mc.generateBezierPath(0, 0, 100, 100)

    // Assertions
    assert.Greater(t, len(points), 20)
    assert.Equal(t, 0.0, points[0].X)
    assert.Equal(t, 100.0, points[len(points)-1].X)
}
```

**Run:**
```bash
go test ./stealth -run TestGenerateBezierPath
```

### 2. Integration Tests

Test multiple components working together.

**Example: Authentication Flow**

```go
func TestLoginAndSessionPersistence(t *testing.T) {
    // Setup
    db := storage.NewTestDatabase()
    cfg := config.DefaultConfig()
    log := logger.New()

    // Create service
    authService := auth.NewService(cfg, db, log)

    // Test login
    err := authService.Login(ctx, page)
    assert.NoError(t, err)

    // Test session save
    err = authService.SaveSession(page)
    assert.NoError(t, err)

    // Test session load
    valid, err := authService.LoadSession(page)
    assert.True(t, valid)
    assert.NoError(t, err)
}
```

### 3. End-to-End Tests

Test complete workflows.

**Example: Full Automation Flow**

```bash
# Run against mock site
go run cmd/main.go --config test-config.yaml
```

## Mock Setup

### Mock Database

```go
func NewTestDatabase() *Database {
    db, _ := NewDatabase(":memory:")  // In-memory SQLite
    return db
}
```

### Mock Browser

```go
func NewMockPage() *MockPage {
    return &MockPage{
        elements: make(map[string]*MockElement),
    }
}

func (m *MockPage) Element(selector string) (*MockElement, error) {
    elem, exists := m.elements[selector]
    if !exists {
        return nil, fmt.Errorf("element not found")
    }
    return elem, nil
}
```

### Mock Configuration

```yaml
# test-config.yaml
safe_mode: true
demo_mode: true

auth:
  login_url: "http://localhost:8080/test"

search:
  base_url: "http://localhost:8080/search"
```

## Test Data

### Sample Profiles

```go
var testProfiles = []*search.Profile{
    {
        ID:      "test-1",
        Name:    "Alice Johnson",
        Title:   "Senior Software Engineer",
        Company: "TechCorp",
        URL:     "http://localhost:8080/profile/alice",
    },
    {
        ID:      "test-2",
        Name:    "Bob Smith",
        Title:   "Engineering Manager",
        Company: "InnovateLabs",
        URL:     "http://localhost:8080/profile/bob",
    },
}
```

### Sample Messages

```go
var testTemplates = []messaging.Template{
    {
        Name:      "test_template",
        Content:   "Hi {{name}}, test message from {{company}}",
        Variables: []string{"name", "company"},
    },
}
```

## Testing Stealth Features

### Bézier Curves

**Visual Test:**

```go
func TestBezierVisual(t *testing.T) {
    mc := NewMouseController(cfg)
    points := mc.generateBezierPath(0, 0, 500, 500)

    // Output for visualization
    for _, p := range points {
        fmt.Printf("%f,%f\n", p.X, p.Y)
    }
}
```

**Visualize with Python:**

```python
import matplotlib.pyplot as plt
import numpy as np

data = np.loadtxt("bezier_output.csv", delimiter=",")
plt.plot(data[:, 0], data[:, 1])
plt.title("Bézier Curve Mouse Path")
plt.show()
```

### Timing Randomness

**Statistical Test:**

```go
func TestTimingRandomness(t *testing.T) {
    tc := NewTimingController(cfg)

    delays := make([]time.Duration, 1000)
    for i := 0; i < 1000; i++ {
        delays[i] = tc.RandomDelay()
    }

    // Calculate statistics
    mean := calculateMean(delays)
    stdDev := calculateStdDev(delays)

    // Assert proper distribution
    assert.InRange(t, mean, cfg.MinDelay, cfg.MaxDelay)
    assert.Greater(t, stdDev, 0)
}
```

### Typing Simulation

**Test Typo Probability:**

```go
func TestTypoFrequency(t *testing.T) {
    ts := NewTypingSimulator(cfg)

    typoCount := 0
    iterations := 1000

    for i := 0; i < iterations; i++ {
        // Mock test to see if typo path is taken
        if rand.Float64() < 0.15 {
            typoCount++
        }
    }

    // Should be approximately 15%
    expectedMin := iterations * 0.10
    expectedMax := iterations * 0.20

    assert.InRange(t, typoCount, expectedMin, expectedMax)
}
```

## Performance Testing

### Benchmark Tests

```go
func BenchmarkMouseMovement(b *testing.B) {
    mc := NewMouseController(cfg)

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        mc.generateBezierPath(0, 0, 1000, 1000)
    }
}

func BenchmarkTypingSimulation(b *testing.B) {
    ts := NewTypingSimulator(cfg)

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        ts.getTypingDelay('a')
    }
}
```

**Run:**
```bash
go test -bench=. -benchmem ./stealth
```

**Example Output:**
```
BenchmarkMouseMovement-8        50000    35421 ns/op    8192 B/op    12 allocs/op
BenchmarkTypingSimulation-8     1000000  1234 ns/op     0 B/op       0 allocs/op
```

## Database Testing

### Schema Tests

```go
func TestDatabaseSchema(t *testing.T) {
    db := NewTestDatabase()

    // Test tables exist
    tables := []string{
        "connection_requests",
        "messages",
        "daily_stats",
        "session_state",
    }

    for _, table := range tables {
        var count int
        err := db.db.QueryRow(
            "SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name=?",
            table,
        ).Scan(&count)

        assert.NoError(t, err)
        assert.Equal(t, 1, count, "Table %s should exist", table)
    }
}
```

### CRUD Tests

```go
func TestConnectionRequestCRUD(t *testing.T) {
    db := NewTestDatabase()

    // Create
    req := &ConnectionRequest{
        ProfileID: "test-1",
        Name:      "Test User",
        Status:    "pending",
    }
    err := db.SaveConnectionRequest(req)
    assert.NoError(t, err)

    // Read
    loaded, err := db.GetConnectionRequest("test-1")
    assert.NoError(t, err)
    assert.Equal(t, "Test User", loaded.Name)

    // Update
    req.Status = "accepted"
    err = db.SaveConnectionRequest(req)
    assert.NoError(t, err)

    // Verify update
    loaded, _ = db.GetConnectionRequest("test-1")
    assert.Equal(t, "accepted", loaded.Status)
}
```

## Mock Server for Testing

### Simple Test Server

```go
// test/mock_server.go
func StartMockServer() *httptest.Server {
    handler := http.NewServeMux()

    // Login page
    handler.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
        html := `
            <input id="email" />
            <input id="password" type="password" />
            <button type="submit">Login</button>
        `
        w.Write([]byte(html))
    })

    // Search page
    handler.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
        html := `
            <div class="profile-card">
                <span class="profile-name">Test User</span>
                <span class="profile-title">Engineer</span>
                <a class="profile-link" href="/profile/test">View</a>
            </div>
        `
        w.Write([]byte(html))
    })

    return httptest.NewServer(handler)
}
```

**Usage:**

```go
func TestWithMockServer(t *testing.T) {
    server := StartMockServer()
    defer server.Close()

    cfg := config.DefaultConfig()
    cfg.Auth.LoginURL = server.URL + "/login"
    cfg.Search.BaseURL = server.URL + "/search"

    // Run tests against mock server
}
```

## Continuous Integration

### GitHub Actions Example

```yaml
# .github/workflows/test.yml
name: Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'

      - name: Install dependencies
        run: go mod download

      - name: Run tests
        run: go test -v -race -coverprofile=coverage.txt ./...

      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          file: ./coverage.txt
```

## Test Coverage Goals

| Package | Target Coverage |
|---------|----------------|
| stealth | 90%+ |
| storage | 85%+ |
| config  | 80%+ |
| auth    | 75%+ |
| search  | 75%+ |
| connections | 70%+ |
| messaging | 70%+ |

## Common Test Patterns

### Table-Driven Tests

```go
func TestCalculations(t *testing.T) {
    tests := []struct {
        name     string
        input    int
        expected int
    }{
        {"zero", 0, 0},
        {"positive", 5, 25},
        {"negative", -3, 9},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := square(tt.input)
            assert.Equal(t, tt.expected, result)
        })
    }
}
```

### Subtests

```go
func TestAuthService(t *testing.T) {
    service := setupTestService()

    t.Run("Login", func(t *testing.T) {
        err := service.Login(ctx, page)
        assert.NoError(t, err)
    })

    t.Run("SaveSession", func(t *testing.T) {
        err := service.SaveSession(page)
        assert.NoError(t, err)
    })

    t.Run("LoadSession", func(t *testing.T) {
        valid, err := service.LoadSession(page)
        assert.True(t, valid)
        assert.NoError(t, err)
    })
}
```

### Test Helpers

```go
func setupTestService() *Service {
    return NewService(
        config.DefaultConfig(),
        storage.NewTestDatabase(),
        logger.New(),
    )
}

func createTestProfile(id string) *Profile {
    return &Profile{
        ID:      id,
        Name:    "Test User",
        Title:   "Engineer",
        Company: "TestCorp",
    }
}
```

## Debugging Tests

### Verbose Test Output

```bash
go test -v ./stealth -run TestBezier
```

### Run Specific Test

```bash
go test ./stealth -run TestGenerateBezierPath
```

### Test with Race Detector

```bash
go test -race ./...
```

### Profile Tests

```bash
go test -cpuprofile=cpu.prof ./stealth
go tool pprof cpu.prof
```

## Manual Testing Checklist

- [ ] Demo mode runs without errors
- [ ] Safe mode validation works
- [ ] Configuration loads correctly
- [ ] Database schema creates successfully
- [ ] Mouse movement appears natural (visual inspection)
- [ ] Typing has realistic delays
- [ ] Random timing is properly distributed
- [ ] Session persistence works across restarts
- [ ] Rate limits enforce correctly
- [ ] Error messages are helpful
- [ ] Logs are structured and informative

## Resources

- [Go Testing Package](https://pkg.go.dev/testing)
- [Testify Assertions](https://github.com/stretchr/testify)
- [Go Test Comments](https://github.com/golang/go/wiki/TestComments)
- [Table Driven Tests](https://github.com/golang/go/wiki/TableDrivenTests)
