# 012. Separate Session Create from Attach

Date: 2026-03-05

## Status

Accepted

## Context

`p` maintains a session history ledger (ADR-011) used by `p --log`.

The tmux attach/switch path uses `syscall.Exec` to replace the current process with `tmux`. Any code after an attach/switch call does not run, which caused launches to never be written to the ledger when logging happened after attaching.

## Decision

Split session lifecycle into two explicit steps:

1. `tmux.CreateSession(...)` prepares a session (create detached, configure, default windows) and returns a `LaunchAction` describing whether the session was created or already existed.
2. `tmux.AttachToSession(...)` attaches/switches (and continues to use `syscall.Exec`).

Callers are responsible for logging (and any other side effects) before invoking attach/switch.

## Consequences

- The history ledger is written reliably for creates and attaches.
- Session attachment remains a true process replacement (no extra wrapper process and no terminal handoff complexity).
- Call sites become slightly more explicit (create + log + attach).

## Alternatives Considered

- Replace `syscall.Exec` with `exec.Command(...).Run()` for attach/switch.
  - Rejected: would require additional terminal/TTY handling and may change interactive behavior.
- Write the history entry from a tmux hook.
  - Rejected: increases tmux configuration surface area and complicates portability.

## Related

- ADR-011: `doc/decisions/011-session-history-ledger.md`
- Planning: `doc/.plan/.done/fix-p-log-empty-history/`
