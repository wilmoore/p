## Plan: Session Launch Quality-of-Life

- **Branch**: feat/session-enhancements
- **Goals**:
  - Allow `p <path> --name custom` to set deterministic, descriptive session names.
  - When tmux reports `duplicate session`, verify directories and automatically attach if they match.
  - Persist session launch history and surface it via `p --log` in a selector-style UI for inspection or relaunch.
- **Out of scope**:
  - Previously planned `p -` shortcut and Gemini session-history prototype.
  - Automatic relaunch of every logged session (history is informational + opt-in relaunch only).
- **Key decisions**:
  - Store history as JSONL under `${XDG_STATE_HOME:-~/.local/state}/p/session-log.jsonl` (override path via `P_HISTORY_PATH`).
  - `p --log` uses the same selector, rendering single-line entries (session · type · timestamp · cwd→target) with a concise footer.
  - CLI parser enforces mutually exclusive modes (path creation vs `--log`).

### Implementation Notes

1. **Cleanup**: remove Gemini artifacts (`internal/sessionhistory`, `.plan/session-handoff*`, duplicated `ListSessions`, `p -` references).
2. **Argument parsing**: add helper to interpret args into structured commands (selector, create-with-optional-name, history).
3. **Duplicate handling**: catch tmux duplicate errors, fetch existing session path via `tmux display-message -p -t name "#{session_path}"`, normalize (abs + `filepath.EvalSymlinks`). Prefer canonical matches; fall back to comparing absolute paths with warnings. Refuse auto-attach when neither comparison is trustworthy.
4. **History ledger**:
   - Fields: timestamp, action (`create` = brand-new session, `attach` = selector attach, `attach-existing` = duplicate auto-attach), invocation cwd, session name, target dir.
   - Append on every successful attach/create; keep max 200 entries by trimming oldest under an exclusive file lock (serialize read+append+rewrite).
5. **History UI**: `p --log` loads ledger entries (latest first) into selector; confirming replays recorded launch (respecting duplicate logic) while Esc simply exits. Replays reuse the same directory validation logic (`createSessionFromPath`) so missing or inaccessible directories fail fast with a clear message.
6. **Docs/tests**: update README/help text; add unit tests for arg parsing, ledger read/write, and duplicate detection helper.

## Related ADRs

- [011. Session History Ledger and Duplicate Attach](../../../doc/decisions/011-session-history-ledger.md)
