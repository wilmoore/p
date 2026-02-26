# Bug Fix: tmux Session Unconfigured

## Bug Summary

Running `p .` creates a new tmux session with stock/unconfigured tmux instead of the styled, vi-mode configured session.

## Steps to Reproduce

1. Run: `p .` (from any directory)
2. Observe: tmux session created with default green status bar

## Expected Behavior

Per ADR-002 and ADR-005:
- Color-agnostic status bar (muted gray on default background)
- Vi-style copy mode bindings (v for selection, y for copy)
- Git branch displayed in status bar

## Actual Behavior

- Stock tmux green status bar
- No custom styling or bindings

## Root Cause Analysis

### What's Working

`execute.go` correctly passes `-f /dev/null` to all tmux invocations:

```go
// Line 42 in execTmux
fullArgs := append([]string{"tmux", "-f", "/dev/null"}, args...)

// Line 48 in runTmux
fullArgs := append([]string{"-f", "/dev/null"}, args...)
```

This disables user config as intended per ADR-002.

### What's Missing

**No runtime configuration injection happens after session creation.**

The codebase has NO implementation of the "minimal injected defaults" specified in:

1. **ADR-002** - "Inject only minimal ergonomic defaults at runtime (vi-style copy mode bindings)"
2. **ADR-005** - Color-agnostic status bar styling:
   - `status-style: bg=default,fg=default`
   - `status-left-style: fg=colour240`
   - `status-right-style: fg=colour240`
   - `window-status-style: fg=colour240`
   - `window-status-current-style: fg=default,bold`
   - Git branch in status-right

The `CreateSession()` function only creates the session, it doesn't configure it.

## Related ADRs

- [ADR-002: Zero-configuration tmux execution](../../doc/decisions/002-zero-config-tmux-execution.md)
- [ADR-005: Color-agnostic Status Bar Styling](../../doc/decisions/005-color-agnostic-status-bar.md)

## Fix Implementation Plan

1. Create a new function `configureSession(sessionName string)` in `execute.go`
2. This function should run tmux set-option commands to configure:
   - Vi-style copy mode bindings
   - Status bar styling
   - Git branch display
3. Call `configureSession()` after session creation in `CreateSession()`
4. Also call it in `AttachToSession()` for existing sessions (to ensure consistent experience)

## Configuration Commands to Inject

```bash
# Vi-style copy mode (ADR-002)
tmux set-option -t SESSION mode-keys vi
tmux bind-key -T copy-mode-vi v send-keys -X begin-selection
tmux bind-key -T copy-mode-vi y send-keys -X copy-selection-and-cancel

# Status bar styling (ADR-005)
tmux set-option -t SESSION status-style 'bg=default,fg=default'
tmux set-option -t SESSION status-left-style 'fg=colour240'
tmux set-option -t SESSION status-right-style 'fg=colour240'
tmux set-option -t SESSION status-right '#(git -C #{pane_current_path} rev-parse --abbrev-ref HEAD 2>/dev/null)'
tmux set-option -t SESSION status-interval 5
tmux set-window-option -t SESSION window-status-style 'fg=colour240'
tmux set-window-option -t SESSION window-status-current-style 'fg=default,bold'
```

## Verification

After fix:
1. Run `p .` in a directory
2. Verify: muted gray status bar (not green)
3. Verify: git branch appears in status-right
4. Enter copy mode (prefix + [), verify vi bindings work

## Related ADRs

- [ADR-002: Zero-configuration tmux execution](../../doc/decisions/002-zero-config-tmux-execution.md)
- [ADR-005: Color-agnostic Status Bar Styling](../../doc/decisions/005-color-agnostic-status-bar.md)
- [ADR-008: Inject Runtime Configuration into tmux Sessions](../../doc/decisions/008-inject-runtime-configuration.md)
