# Bug Fix: Output Display and CLI Flags

## Branch
`fix/output-display-and-cli-flags`

## Status
Implementation complete. Ready for verification.

## Issues Identified

### 1. Display Bug (Critical)
- **Symptom:** Sessions appear scattered randomly across the screen instead of in a vertical list
- **Environment:** Running inside tmux
- **Root Cause Investigation:**
  - The `render()` function in `selector.go` uses `\033[H\033[J` to clear and home cursor
  - Then prints each session with `\r\n` line endings (correct for raw mode)
  - The scattered appearance suggests the terminal isn't properly handling the escape sequences
  - OR there's a race condition between clearing the screen and rendering

### 2. Missing `--version` flag
- **Expected:** `p --version` prints version info
- **Actual:** No argument parsing exists; tool tries to run selector

### 3. Missing `--help` flag
- **Expected:** `p --help` prints usage information
- **Actual:** No argument parsing exists; tool tries to run selector

### 4. Missing path argument feature
- **Expected:** `p .` or `p /path` creates new session in that directory
- **Actual:** Tool only lists and attaches to existing sessions
- **Note:** ADR-006 simplified to session-only mode, but this feature request makes sense

### 5. Missing config propagation
- **Expected:** When attaching to existing session, apply p's configuration
- **Actual:** Sessions don't get p's tmux settings when attached to

## Severity
Critical (blocks work)

## Reproduction Steps
1. Have multiple tmux sessions running
2. Run `p` from inside tmux
3. Observe: sessions displayed in scattered positions

## Root Cause Analysis

### Display Issue
The render function uses standard ANSI escape sequences that should work universally:
- `\033[H` - move cursor to home (top-left)
- `\033[J` - clear from cursor to end of screen

Possible causes:
1. Terminal alternate screen buffer not being used
2. Raw mode interaction with tmux's pseudo-terminal
3. Missing terminal capability detection

### CLI Flags
Simply not implemented - main.go has no argument parsing.

## Implementation Plan

### Phase 1: CLI Infrastructure
1. Add version string with build-time injection via `-ldflags`
2. Add argument parsing for `--version`, `--help`, path arguments
3. Define version constant/variable

### Phase 2: Display Fix
1. Investigate alternate screen buffer (`\033[?1049h/l`)
2. Ensure proper terminal initialization before raw mode
3. Add terminal capability detection if needed

### Phase 3: Path Argument Feature
1. When path argument provided, create new session in that directory
2. Handle `.` for current directory
3. Handle relative and absolute paths

### Phase 4: Config Propagation
1. When attaching to existing session, inject p's tmux configuration
2. Use `set-option` commands after attach

## Related ADRs
- ADR-006: Simplify to Session-Only Mode (current design)
- The path argument feature may need a new ADR or amendment

---

## Implementation Complete

### Changes Made

1. **internal/ui/selector.go**
   - Added alternate screen buffer (`\033[?1049h`/`\033[?1049l`) for clean display
   - Removed `clearScreen()` function (no longer needed with alternate buffer)
   - Simplified exit paths to rely on deferred cleanup

2. **main.go**
   - Added `Version` variable for build-time injection
   - Added `usage` constant with comprehensive help text
   - Implemented CLI argument parsing for `--version`, `-v`, `--help`, `-h`
   - Added `createSessionFromPath()` function for path argument handling
   - Supports `.` for current directory, `~` expansion, relative and absolute paths

3. **internal/tmux/execute.go**
   - Added `CreateSession(name, workingDir string)` function
   - Added `runTmux()` helper for non-exec tmux commands
   - Handles both inside-tmux (detached creation + switch) and outside-tmux (direct attach) cases

4. **README.md**
   - Updated usage examples with new CLI options
   - Documented path argument feature

### Testing Performed

- `p --version` -> prints "p dev"
- `p --help` -> prints usage
- `p -v` -> prints version (short form)
- `p -h` -> prints help (short form)
- `p --unknown` -> error message with help hint
- `p /nonexistent` -> error message about missing directory
- `go build` -> success
- `go vet ./...` -> no issues
- `gofmt -l` -> no formatting issues

### Deferred Work

- **Config propagation when attaching to existing sessions** - This is more complex and warrants a separate feature branch. Added to backlog considerations.
