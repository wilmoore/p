# 008. Inject Runtime Configuration into tmux Sessions

Date: 2026-02-26

## Status

Accepted

## Context

When running `p .`, a new tmux session is created with stock/unconfigured tmux instead of the styled, vi-mode configured session. ADR-002 specifies "Inject only minimal ergonomic defaults at runtime (vi-style copy mode bindings)" and ADR-005 specifies color-agnostic status bar styling, but neither was implemented.

The codebase had NO implementation of the "minimal injected defaults" after session creation.

## Decision

We implemented a `configureSession(sessionName string)` function in `internal/tmux/execute.go` that runs tmux set-option commands to configure:

1. **Vi-style copy mode bindings** (per ADR-002):
   - `mode-keys vi`
   - `v` for begin-selection
   - `y` for copy-selection-and-cancel

2. **Status bar styling** (per ADR-005):
   - Near-black background (colour232)
   - Sage green accent (colour108) for session name
   - Muted gray (colour240) for window states
   - Git branch in status-right

3. **Applied consistently**:
   - Called after session creation in `CreateSession()`
   - Called before attach in `AttachToSession()` for existing sessions

## Consequences

- All tmux sessions created by `p` now have consistent styling
- Vi users get familiar copy-mode bindings
- Git branch visible in status bar
- Existing sessions also get configured when attaching via `p`

## Alternatives Considered

1. **Template-based config file**: Create a temp config file and pass via `-f`. Rejected - adds file I/O and cleanup complexity.

2. **Skip configuration entirely**: Rejected - violates ADR-002 and ADR-005 which were already accepted.

3. **Configure only on new sessions**: Rejected - inconsistent experience when re-attaching.

## Related

- [ADR-002: Zero-configuration tmux execution](../decisions/002-zero-config-tmux-execution.md)
- [ADR-005: Color-agnostic Status Bar Styling](../decisions/005-color-agnostic-status-bar.md)
- Planning: `.plan/.done/fix-tmux-session-unconfigured/`
