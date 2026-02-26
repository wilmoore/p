# Implementation Plan: Display Bug and CLI Flags

## Summary

Fix the scattered display issue when running inside tmux, add CLI flags (`--version`, `--help`), add path argument support for creating sessions in specific directories, and ensure configuration is applied when attaching to sessions.

## Root Cause Analysis

### Display Bug
Looking at the screenshot more carefully, the scattered display resembles what happens when:
1. Terminal alternate screen buffer is not used, so old content mixes with new
2. The cursor positioning escape sequences are being interpreted as printable characters
3. The raw terminal mode interacts oddly with tmux's pseudo-terminal

The render function uses `\033[H\033[J` which should work, but the sessions appear randomly positioned. This suggests the terminal isn't in the expected state.

**Hypothesis:** The issue is that we're using raw mode but not entering the alternate screen buffer. When the selector runs and then the screen "clears," we're actually just clearing from cursor to end, not truly clearing the display. The solution is to use the alternate screen buffer (`\033[?1049h` on entry, `\033[?1049l` on exit).

### Missing CLI Flags
Simply not implemented - `main.go` has no argument parsing.

## Implementation Steps

### Phase 1: CLI Flag Infrastructure

1. **Add version variable** with build-time injection
   - Add `var Version = "dev"`
   - Build with `-ldflags "-X main.Version=v0.x.x"`

2. **Parse command-line arguments** in main.go before running selector
   - `--version`, `-v`: Print version and exit
   - `--help`, `-h`: Print usage and exit
   - Positional arg: Path to directory for new session

### Phase 2: Fix Display Bug

1. **Use alternate screen buffer** in selector.go
   - On entry: `\033[?1049h` (switch to alternate buffer)
   - On exit: `\033[?1049l` (switch back to main buffer)

2. **Ensure proper cleanup** on all exit paths
   - Normal selection (Enter)
   - Cancel (Esc, Ctrl+C)
   - Error cases

### Phase 3: Path Argument Feature

1. **Add session creation from path**
   - Parse positional argument as directory path
   - Resolve `.` to `$PWD`
   - Resolve relative paths to absolute
   - Derive session name from directory basename
   - Create new tmux session with working directory set

2. **Update execute.go** with `CreateSession(name, path string)` function

### Phase 4: Config Propagation (Deferred)

The request to propagate configuration when attaching to existing sessions is more complex and may warrant a separate feature branch. Current behavior per ADR-006 is to use `-f /dev/null` to ignore config. Propagating config to existing sessions would require:
- Running `set-option` commands after attach
- Deciding which options to inject

**Recommendation:** Track as separate backlog item.

## Files to Modify

1. `main.go` - Add CLI parsing, version, help, path handling
2. `internal/ui/selector.go` - Fix display with alternate screen buffer
3. `internal/tmux/execute.go` - Add CreateSession function
4. `internal/tmux/session.go` - Add session name generation from path

## Testing Plan

1. Verify `p --version` prints version
2. Verify `p --help` prints usage
3. Verify `p` without args shows selector correctly inside tmux
4. Verify `p .` creates session in current directory
5. Verify `p /path/to/dir` creates session in specified directory
6. Verify display is clean and sessions are vertically aligned

## Risks

- Alternate screen buffer might not be supported in all terminals (rare)
- Path argument feature expands scope beyond ADR-006's session-only design

## ADR Considerations

ADR-006 simplified to session-only mode. Adding path argument for session creation partially reverses this simplification. However, it's a pragmatic addition that doesn't require CDPATH - the user explicitly provides the path.

**Recommendation:** Proceed with implementation, document as refinement to ADR-006.

## Related ADRs

- [007. Alternate screen buffer for selector display](../../doc/decisions/007-alternate-screen-buffer-for-selector.md)
