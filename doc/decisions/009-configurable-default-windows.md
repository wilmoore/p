# 009. Configurable Default Windows for New Sessions

Date: 2026-03-02

## Status

Accepted

## Context

New tmux sessions created by `p` always opened with a single, generically named window, forcing users to manually create their preferred layout every time. Users consistently open a home window rooted at the project path plus an additional shell (e.g., for commands or editors). The project already enforces a curated tmux environment (per ADR-002 and ADR-008), so allowing configurable default windows at session creation aligns with that direction while keeping user configuration lightweight.

## Decision

Introduce the `P_WINDOWS` environment variable (defaulting to `home,cmd`) to describe a comma-separated list of window names that should exist the moment a session is created. After starting the session we rename window `0` to the first name, create additional windows for the remaining entries (rooted at the requested working directory), and re-select window `0` so `home` remains active. This logic only runs for brand-new sessions, leaving attachments to existing sessions untouched.

## Consequences

- Positive: Users land in a predictable, ready-to-use layout without repetitive manual window creation.
- Positive: Teams can standardize on an environment variable (or wrapper script) without distributing tmux config files.
- Negative: Session creation now performs more tmux operations, slightly increasing startup time.
- Negative: Misconfigured `P_WINDOWS` values (e.g., trailing commas) can result in empty window names; we mitigate by trimming and ignoring blanks.

## Alternatives Considered

- **Do nothing** — rejected because it keeps the repetitive manual workflow the feature set aims to remove.
- **Require `.tmux.conf` or hooks** — rejected; this contradicts ADR-002/008 that `p` should own tmux configuration without relying on user files.
- **Hardcode more windows without configurability** — rejected because different teams require different layouts and this would force forks.

## Related

- Planning: `.plan/.done/feat-default-windows-and-consistent-background/`
