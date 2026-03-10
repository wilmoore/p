# 013. Shared Selector Engine for Terminal UI

Date: 2026-03-10

## Status

Accepted

## Context

`p` contains multiple interactive terminal selectors (session selector and history selector).

The initial implementation duplicated:

- raw terminal mode setup/restore
- alternate screen buffer handling
- key parsing (Esc/Ctrl+C, arrows, Ctrl+J/K/N/P, Enter, Backspace)
- filtering and selection clamping
- rendering loops

This duplication increased the chance of inconsistent behavior and made UI improvements harder to apply uniformly.

## Decision

Introduce a shared selector engine used by all interactive selectors:

- `internal/ui/selector_engine.go`: common event loop, filtering, viewport rendering
- `internal/ui/keys.go`: shared key parsing and escape-sequence handling
- shared ANSI/text constants moved into `internal/ui/ansi.go` and `internal/ui/text.go`

Selectors provide an adapter describing:

- how to render a row
- how to produce search text
- optional summary line and direct-select behavior

## Consequences

### Positive

- Consistent navigation behavior across selectors.
- Reduced code duplication and easier future UI changes.
- History UI can become width-aware without re-implementing the input loop.

### Negative

- Adds abstraction layers that must remain small and well-tested.
- Changes to the engine affect all selectors.

## Alternatives Considered

1. **Keep separate implementations** - rejected due to ongoing drift and maintenance cost.
2. **Adopt a third-party TUI framework** - rejected to keep dependencies and complexity minimal.

## Related

- Planning: `doc/.plan/.done/refactor-p-log-history-ui/`
- ADR-011: `doc/decisions/011-session-history-ledger.md`
