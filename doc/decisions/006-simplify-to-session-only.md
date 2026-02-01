# 006. Simplify to Session-Only Mode

Date: 2026-01-29

## Status

Accepted

Supersedes: [001. Use CDPATH for Discovery](001-use-cdpath-for-discovery.md)

## Context

The original `p` tool combined two distinct responsibilities:
1. Discovering project directories via CDPATH
2. Managing tmux sessions

This coupling introduced complexity:
- Session naming required hash-suffix collision handling
- Drill-down navigation added UI complexity
- CDPATH configuration was a prerequisite that confused users
- The directory discovery logic accounted for ~40% of the codebase

User feedback indicated most usage was attaching to existing sessions, not creating new ones from directories.

## Decision

Simplify `p` to focus solely on tmux session management:
- List existing tmux sessions
- Provide an interactive selector with type-to-filter
- Attach to or switch to the selected session

Remove all CDPATH-related functionality:
- Directory discovery
- Session creation from directories
- Collision-free naming logic
- Drill-down navigation

## Consequences

### Positive

- Dramatically reduced codebase (~60% fewer lines)
- Zero configuration required (no CDPATH dependency)
- Single responsibility: session selection
- Faster startup (no directory scanning)
- Simpler mental model for users

### Negative

- Cannot create new sessions from directories (use `tmux new-session` directly)
- Previous ADRs 001, 003, 004, 005 no longer apply to current design
- Users relying on CDPATH discovery need alternative workflow

## Alternatives Considered

1. **Keep both modes**: Session-only when CDPATH unset, full mode otherwise
   - Rejected: Adds conditional complexity, testing burden

2. **Separate tools**: `p` for sessions, `pd` for directories
   - Rejected: Fragmented UX, maintenance overhead

3. **Plugin architecture**: Load directory discovery as optional module
   - Rejected: Over-engineering for a simple CLI tool

## Related

- Planning: N/A (simplified directly on feature branch)
- Supersedes: ADR 001, ADR 003, ADR 004, ADR 005
