# Session Handoff Ledger

Updated: 2026-03-05T15:13:35.977Z
Current session: session-2026-03-05T15-13-35-847Z-c97f8d0f

## Outstanding Snapshots (1)

1. [pending] session-2026-03-05T15-13-35-847Z-c97f8d0f — fix/p-log-empty-history (dirty)
   File: doc/.plan/session-handoff/sessions/session-2026-03-05T15-13-35-847Z-c97f8d0f.md
   Updated: 2026-03-05T15:13:35.847Z

## Recent Activity

- None

## Commands

- `node bin/session-handoff.mjs list` — show pending snapshots
- `node bin/session-handoff.mjs ack <id> [--note "done"]` — mark complete
- `node bin/session-handoff.mjs dismiss <id> --reason "why"` — abandon work
- `node bin/session-handoff.mjs write --trigger "/pro:session.handoff"` — capture a fresh snapshot

All snapshots live under `doc/.plan/session-handoff/sessions/`. Review each file before acknowledging or dismissing it.
