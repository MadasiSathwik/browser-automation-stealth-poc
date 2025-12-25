# Quick Start Demo Guide

This guide will help you run the browser automation POC in demo mode within 5 minutes.

## Prerequisites

- Go 1.21+ installed ([download here](https://golang.org/dl/))
- Terminal access

## Quick Start (3 Steps)

### 1. Install Dependencies

```bash
go mod download
```

This downloads Rod and other required packages.

### 2. Run Demo Mode

```bash
go run cmd/main.go --demo
```

Or using Make:

```bash
make demo
```

### 3. Observe Output

You'll see simulated actions printed to console:

```
INFO[2024-01-15 10:23:45] Starting Browser Automation POC
INFO[2024-01-15 10:23:45] Demo Mode: true | Safe Mode: true
INFO[2024-01-15 10:23:45] === DEMO MODE: Simulated Execution ===

[ACTION] Initializing browser with stealth configuration
  - User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64)...
  - Viewport: 1366x768 (randomized)
  - WebDriver flag: hidden

[ACTION] Navigating to login page
  - URL: https://mock-professional-network.example.com/login
  - Applying human-like mouse movement (Bézier curve)
  - Random scroll: 120px with velocity variance

[ACTION] Entering credentials
  - Typing email with realistic delays (45-120ms per char)
  - Simulated typo: 'usre@' -> backspace -> 'user@'
  - Typing password with increased delay (security field)

[ACTION] Searching for profiles
  - Query: 'Software Engineer at Tech Companies'
  - Found 15 profiles on page 1

[ACTION] Processing profile 1: Alice Johnson - Senior Software Engineer @ TechCorp
  - Hovering over profile card (250ms)
  - Moving mouse to 'Connect' button with Bézier curve
  - Clicking 'Connect'
  - Adding personalized note
  - Connection request recorded in database
  - Cooldown: 47s (randomized)

[SUMMARY] Demo completed successfully
  - Total profiles processed: 3
  - Connection requests sent: 3
  - Follow-up messages sent: 1
  - Average action delay: 45s
  - Stealth techniques applied: 8
```

## What's Happening?

The demo mode:

✅ **Prints actions** instead of executing them
✅ **Shows stealth techniques** in action
✅ **Demonstrates timing** and delays
✅ **Simulates realistic flow** from login to messaging
✅ **Safe to record** for demonstration videos

## Next Steps

### View Configuration

```bash
cat config.yaml
```

Key settings to notice:
- `safe_mode: true` - prevents real-world execution
- `demo_mode: false` - can be set to true here too
- Configurable DOM selectors
- Timing and stealth parameters

### Check Database Schema

```bash
# After running demo, check what would be stored
sqlite3 automation.db
.schema
.quit
```

### Build Binary

```bash
make build
./browser-automation --demo
```

### Explore Code

Start with these files to understand the architecture:

1. **cmd/main.go** - Entry point, shows overall flow
2. **stealth/mouse.go** - Bézier curve implementation
3. **stealth/typing.go** - Human-like typing with typos
4. **config/config.go** - Configuration structure

## Advanced Demo: With Mock Site

If you want to see the browser actually open:

1. Create a simple HTML test page:

```html
<!DOCTYPE html>
<html>
<head><title>Mock Network</title></head>
<body>
    <input id="email" placeholder="Email">
    <input id="password" type="password" placeholder="Password">
    <button type="submit">Login</button>

    <div class="profile-card">
        <span class="profile-name">John Doe</span>
        <span class="profile-title">Software Engineer</span>
        <button>Connect</button>
    </div>
</body>
</html>
```

2. Serve it locally:

```bash
# Python
python3 -m http.server 8000

# Or Node.js
npx http-server -p 8000
```

3. Update config.yaml:

```yaml
auth:
  login_url: "http://localhost:8000/test.html"
```

4. Run with browser visible:

```bash
# Edit config.yaml: browser.headless: false
go run cmd/main.go
```

Watch the browser interact with the page using human-like movements!

## Stealth Techniques Demonstrated

During demo mode, you'll see these techniques mentioned:

1. **Bézier Curves** - Smooth mouse paths
2. **Velocity Variance** - Speed changes during movement
3. **Micro-corrections** - Small adjustments
4. **Typo Simulation** - Realistic typing errors
5. **Think Delays** - Pauses before actions
6. **Random Scrolling** - Natural page interaction
7. **Viewport Randomization** - Different window sizes
8. **User-Agent Rotation** - Browser identification

## Safety Features

The demo is completely safe:

- ✅ No network requests to real sites
- ✅ No actual login attempts
- ✅ Prints actions only
- ✅ Safe to run repeatedly
- ✅ No credentials required

## Troubleshooting

**"command not found: go"**
- Install Go from https://golang.org/dl/

**"cannot find package"**
- Run: `go mod download`

**"make: command not found"**
- Use: `go run cmd/main.go --demo` instead

## Recording a Demo Video

This is perfect for screen recording:

```bash
# Terminal 1: Start recording
# Then run:
make demo

# Output is colorized and formatted for demos
# Shows realistic timing (simulated)
# Demonstrates all stealth techniques
```

## Understanding the Code

### Key Files for Review

1. **stealth/mouse.go** (Line 45-75): Bézier curve generation
2. **stealth/typing.go** (Line 30-60): Typo simulation
3. **stealth/timing.go** (Line 15-40): Randomized delays
4. **connections/connect.go** (Line 90-120): Connection flow
5. **cmd/main.go** (Line 60-150): Demo mode logic

### Architecture Highlights

```
main.go
  └─> Config Loading (config/)
  └─> Database Init (storage/)
  └─> Browser Launch (Rod)
  └─> Auth Service (auth/)
       └─> Login Handler + Session Manager
  └─> Search Service (search/)
       └─> Profile extraction + Pagination
  └─> Connection Service (connections/)
       └─> Stealth controller + Rate limiter
  └─> Messaging Service (messaging/)
       └─> Template engine + Delivery
```

## Questions?

Check the main README.md for:
- Detailed architecture explanation
- Full configuration options
- Troubleshooting guide
- Legal/ethical considerations

---

**Total Time: 2-3 minutes** ⚡

Enjoy exploring the code!
