# 003. Collision-free session naming with hash suffix

Date: 2026-01-04

## Status

Accepted

## Context

Session names are derived from directory names. When multiple CDPATH entries contain directories with the same base name (e.g., `/work/api` and `/personal/api`), a naming collision occurs. tmux does not allow duplicate session names.

## Decision

Use a deterministic hash suffix to resolve collisions:
1. First directory to register a base name gets the clean name (e.g., `api`)
2. Subsequent directories with the same base name get a hash suffix (e.g., `api-a1b2c3`)
3. The hash is derived from the full path, ensuring stability across sessions
4. Hash is 6 characters (short but collision-resistant for practical use)

## Consequences

**Positive:**
- Fully automatic - no user intervention required
- Deterministic - same path always produces same name
- Short names when possible - only adds suffix when needed
- No prompts or manual resolution required

**Negative:**
- Hash suffixes are not human-readable
- Order matters - first registration wins the clean name
- Registry state is not persisted (resets each run)

## Alternatives Considered

1. **Full path as session name** - Rejected because names become unwieldy (`-home-user-projects-api`).

2. **User prompt for conflicts** - Rejected because it breaks the zero-interaction goal.

3. **Numbered suffixes (api-1, api-2)** - Rejected because numbers are not stable across sessions.

4. **Parent directory prefix (work-api)** - Considered but rejected because it makes all names longer even without conflicts.

## Related

- Planning: `.plan/.done/mvp-p/`
- Implementation: `internal/tmux/naming.go`
