# Vim/Emacs Keyboard Navigation

## Branch
`feat/vim-keyboard-navigation`

## Summary
Add vim-style (Ctrl+J/K) and emacs-style (Ctrl+N/P) keyboard navigation to the session selector, matching fzf conventions.

## Requirements
| Keybinding | Action |
|------------|--------|
| Ctrl+J | Move selection down |
| Ctrl+K | Move selection up |
| Ctrl+N | Move selection down |
| Ctrl+P | Move selection up |

Existing arrow key navigation remains unchanged.

## Implementation
- File: `internal/ui/selector.go`
- Add control character handlers in the input switch statement
- Ctrl+J = 0x0A (10), Ctrl+K = 0x0B (11), Ctrl+N = 0x0E (14), Ctrl+P = 0x10 (16)

## Related ADRs
- ADR-002: Zero-configuration tmux execution (establishes vi-style ergonomics principle)
