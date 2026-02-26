# 007. Use Alternate Screen Buffer for Selector Display

Date: 2026-02-26

## Status

Accepted

## Context

When running `p` inside tmux, the session selector displayed scattered text across the screen instead of a clean vertical list. The original implementation used `\033[H\033[J` (cursor home + clear to end) to clear the screen before each render, but this approach failed to properly isolate the selector UI from the existing terminal content.

The visual artifacts occurred because:
1. The main terminal buffer retained previous content
2. Clearing from cursor position mixed with tmux's pseudo-terminal handling
3. Raw mode input processing interfered with screen state

## Decision

Use the alternate screen buffer (`\033[?1049h` on entry, `\033[?1049l` on exit) instead of manual screen clearing. This provides a completely separate buffer for the selector UI that:

1. Preserves the user's terminal content when the selector exits
2. Guarantees a clean slate for rendering
3. Follows standard terminal application conventions (vim, less, htop all use this)

Additionally, removed the `clearScreen()` function and all explicit screen clearing calls, relying on the alternate buffer's clean state and the deferred restore.

## Consequences

**Positive:**
- Clean, consistent display regardless of terminal state
- Terminal content preserved after selector exits
- Simpler code with fewer explicit screen management calls
- Works correctly inside tmux and in direct terminal

**Negative:**
- Very old terminals (pre-1999) may not support alternate screen buffer
- Some terminal multiplexers with non-standard implementations might behave unexpectedly

The trade-off is acceptable because modern terminals universally support this feature.

## Alternatives Considered

1. **Enhanced clear sequence** (`\033[2J\033[H`) - Still has issues with buffer mixing
2. **Double buffering in application** - Complex and reinvents terminal functionality
3. **ncurses/termbox** - External dependency for a simple UI

## Related

- Planning: `.plan/.done/fix-output-display-and-cli-flags/`
- ADR-006: Simplify to Session-Only Mode
