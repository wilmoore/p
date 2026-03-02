# Plan: Default Windows + Consistent Background

## Branch
`feat/default-windows-and-consistent-background`

## Requirements

### Feature 1: Default Windows on New Session
- When creating a **new** tmux session, create configurable default windows
- Default window names: `home,cmd` (window 0: "home", window 1: "cmd")
- Window 0 should remain active after creation
- Overridable via `P_WINDOWS` environment variable (comma-separated)
- Example: `P_WINDOWS=home,web,cmd p ~/path/to/project`
- Only applies to **new** sessions, not when attaching to existing ones

### Feature 2: Background Color Consistency
- **Current behavior**: Window 0 has black background, subsequent windows have grey
- **Expected**: All windows have black background (colour16)
- **Root cause**: `window-style bg=colour16` is session-level but may not apply to new windows
- **Fix**: Ensure background styling applies to all windows

## Implementation Steps

### Step 1: Fix background color for new windows
In `configureSession()`, the `window-style` setting was using session-targeted options, but `window-style` and `window-active-style` are window options that new windows don't properly inherit from session defaults.

**Root cause:** Session-targeted `set-option -t session` for window options only sets the session's default, but newly created windows may inherit from the global default instead.

**Fix in `internal/tmux/execute.go`:**
```go
// Use global options (-g) so all windows inherit the black background
runTmuxSilent("set-option", "-g", "window-style", "bg=colour16")
runTmuxSilent("set-option", "-g", "window-active-style", "bg=colour16")
```

This is consistent with `p`'s philosophy of creating a managed tmux environment (using `-f /dev/null` to ignore user config).

### Step 2: Add default window creation
Modify `CreateSession()` to:
1. Read `P_WINDOWS` env var (default: `home,cmd`)
2. Parse comma-separated window names
3. Rename window 0 to first name
4. Create additional windows for remaining names
5. Select window 0

**New function in `internal/tmux/execute.go`:**
```go
func createDefaultWindows(sessionName, workingDir string) {
    // Get window names from env or use default
    windowsEnv := os.Getenv("P_WINDOWS")
    if windowsEnv == "" {
        windowsEnv = "home,cmd"
    }

    windows := strings.Split(windowsEnv, ",")
    if len(windows) == 0 {
        return
    }

    target := "-t" + sessionName

    // Rename first window (index 0)
    runTmuxSilent("rename-window", target+":0", strings.TrimSpace(windows[0]))

    // Create additional windows
    for i := 1; i < len(windows); i++ {
        name := strings.TrimSpace(windows[i])
        runTmuxSilent("new-window", target, "-c", workingDir, "-n", name)
    }

    // Select window 0
    runTmuxSilent("select-window", target+":0")
}
```

### Step 3: Call from CreateSession
Call `createDefaultWindows()` after session creation but before `configureSession()`.

### Step 4: Update README
Document the `P_WINDOWS` environment variable.

## Testing Plan

1. **Manual test - default windows**:
   ```bash
   # Kill any existing 'test' session
   tmux kill-session -t test 2>/dev/null
   # Create new session
   go build -o p && ./p /tmp
   # Verify: Should have windows 0:home and 1:cmd
   ```

2. **Manual test - custom windows**:
   ```bash
   tmux kill-session -t test2 2>/dev/null
   P_WINDOWS=main,server,logs ./p /tmp
   # Verify: Should have 0:main, 1:server, 2:logs
   ```

3. **Manual test - background color**:
   ```bash
   # In the session, switch between windows
   # All windows should have black background, not grey
   ```

4. **Existing session test**:
   ```bash
   # Attach to existing session
   ./p  # select existing session
   # Verify: No new windows created, existing windows preserved
   ```

## Files Changed

- `internal/tmux/execute.go` - Add window creation logic, fix background
- `README.md` - Document P_WINDOWS env var

## Related ADRs

- ADR-002: Zero-configuration tmux execution
- ADR-005: Color-agnostic status bar styling
- ADR-008: Inject runtime configuration into tmux sessions
- ADR-009: Configurable default windows for new sessions
- ADR-010: Apply global window style for consistent backgrounds
