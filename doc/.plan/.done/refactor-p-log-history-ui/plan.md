## Plan: Refactor `p --log` History UI

### Goals

- Improve the readability/layout of the history selector list.
- Reduce duplicated selector/terminal code between `internal/ui/selector.go` and `internal/ui/history_selector.go`.
- Replace magic strings/numbers with named constants.
- Centralize user-facing text.
- Improve error reporting so errors are friendly by default, but retain debug detail when needed.

### Constraints / Prior Decisions

- ADR-011: `p --log` reuses the selector UI concept for browsing history.
- ADR-012: history logging must occur before any exec-based tmux attach/switch.

### Related ADRs

- ADR-013: `doc/decisions/013-shared-selector-engine.md`
- ADR-014: `doc/decisions/014-friendly-errors-with-debug-mode.md`

### Approach

1. Extract shared terminal selector loop + key handling into a common helper used by both selectors.
2. Rework history row formatting to be width-aware (terminal width) and consistently truncated.
3. Centralize UI strings and key/escape constants.
4. Add focused unit tests for formatting/truncation and selector filtering.
5. Add a small error formatting helper to keep messages short, with opt-in debug detail (env var).
