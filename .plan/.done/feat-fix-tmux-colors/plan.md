# Fix tmux Colors - Match Savvy AI Design

## Branch
`feat/fix-tmux-colors`

## Summary
Fix inverted tmux color scheme to match the intended Savvy AI aesthetic.

## Problem
The tmux status bar colors were inverted from the intended design:
- Main interface background: dark gray (inherited from terminal)
- Footer/status line: colour232 (near-black #050505)

## Solution
Correct the color values to match Savvy AI design:
- Main interface: pure black (colour16)
- Footer/status line: dark gray (colour235 - matches Savvy AI navbar)
- Footer text: colour240 on dark gray
- No borders (pane-border-status off)

## Changes
1. `status-style`: bg=colour232 → bg=colour235
2. `window-style`: added bg=colour16 (pure black pane background)
3. `pane-border-status`: set to off

## Related
- Backlog item #23
- ADR-005: Color-agnostic status bar styling
- ADR-008: Inject runtime configuration
