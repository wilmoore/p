## Plan: `p --log` shows no session history

- **Branch**: fix/p-log-empty-history
- **Bug**: After creating/attaching two tmux sessions, running `p --log` prints `No session history yet`.
- **Context**: Feature implemented via ADR 011. Ledger stored under `${XDG_STATE_HOME:-~/.local/state}/p/session-log.jsonl`. Need to confirm entries are written whenever sessions are created/attached.

### Bug Details

1. Run `p browser-extension-add-to-attio`
2. Run `p browser-extension-conversation-titles-chatgpt`
3. Run `p --log`

**Expected**: History selector lists at least the two sessions just created.

**Actual**: CLI prints `No session history yet.`

**Environment**: macOS (per screenshot) running dev build (`p --version` -> `dev`).

**Severity**: Blocks work (history mode unusable).

### Notes

- ADR 011 defines ledger persistence and `p --log` behavior; fixes must respect JSONL storage, locking, and UI structure described there.
- Need to inspect ledger file under `${XDG_STATE_HOME:-~/.local/state}/p/session-log.jsonl` and verify writes occur when `p` launches sessions.
- Root cause hypothesis: `logLaunch(...)` is invoked after `tmux.CreateSession` / `tmux.AttachToSession`. Those functions call `syscall.Exec`, replacing the current process with `tmux`, so control never returns to execute the logging statements. Therefore no ledger entries are ever written.
- Fix direction: ensure logging happens before any exec-based attach/switch. Refactor session lifecycle so the caller can log, then attach.

### Related ADRs

- ADR-012: `doc/decisions/012-separate-session-create-from-attach.md`
