# 014. Friendly Errors with Opt-in Debug Detail

Date: 2026-03-10

## Status

Accepted

## Context

During development and refactoring, errors often contain useful context (wrapping, underlying causes).
However, printing raw error chains by default can be noisy and unfriendly for normal CLI use.

We need a consistent approach that:

- shows a concise, user-facing message by default
- retains deeper diagnostics when explicitly requested

## Decision

Add a small error formatting helper:

- `internal/clierr/clierr.go` provides:
  - `Wrap(userMessage, err)` for attaching a user message to an underlying error
  - `Format(err)` for printing user-friendly output
- When `P_DEBUG=1|true|yes|on`, `Format` includes the unwrap chain as a debug section.

Additionally, move stable user-facing strings into `internal/i18n/messages.go` to reduce scattered literals.

## Consequences

### Positive

- Cleaner CLI output for common failures.
- Debug detail remains available without recompiling or changing code.
- Fewer duplicated string literals in main flows.

### Negative

- Requires discipline to wrap errors where a user message is helpful.
- Adds a small amount of new internal surface area.

## Alternatives Considered

1. **Always print full wrapped errors** - rejected as too verbose for normal use.
2. **Introduce a full logging framework** - rejected as unnecessary for current scope.

## Related

- Planning: `doc/.plan/.done/refactor-p-log-history-ui/`
