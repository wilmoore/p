# Session Handoff Ledger

Updated: 2026-03-10T12:37:57.023Z
Current session: 

## Outstanding Snapshots (0)

- None

## Recent Activity

- acked: session-2026-03-05T15-13-35-847Z-c97f8d0f — merged in #17

## Commands

- `node bin/session-handoff.mjs list` — show pending snapshots
- `node bin/session-handoff.mjs ack <id> [--note "done"]` — mark complete
- `node bin/session-handoff.mjs dismiss <id> --reason "why"` — abandon work
- `node bin/session-handoff.mjs write --trigger "/pro:session.handoff"` — capture a fresh snapshot

All snapshots live under `doc/.plan/session-handoff/sessions/`. Review each file before acknowledging or dismissing it.
