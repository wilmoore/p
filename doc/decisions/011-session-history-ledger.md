# 011. Session History Ledger and Duplicate Attach

Date: 2026-03-04

## Status

Accepted

## Context

Users needed better ergonomics when moving between tmux sessions managed by `p`:

- Creating a session that already exists fails with "duplicate session" even when targeting the same directory.
- After a crash it is difficult to remember which sessions were open.
- The CLI needed a way to relaunch previous workspaces without retyping long paths or names.

These gaps prompted enhancements to session creation, logging, and history browsing.

## Decision

- Detect tmux "duplicate session" errors during creation, resolve the existing session path via `tmux display-message #{session_path}`, normalize both canonical and absolute paths, and automatically attach when the directories match.
- Accept user-provided session names via `p <path> --name <custom>` while preserving absolute path validation.
- Persist every successful attach/create into a JSON Lines ledger stored at `${XDG_STATE_HOME:-~/.local/state}/p/session-log.jsonl` (overridable via `P_HISTORY_PATH`). Each entry records timestamp, action (`create`, `attach`, `attach-existing`), invocation directory, target directory, and session name.
- Serialize ledger updates with file locks and atomic temp-file writes, trimming to the 200 most recent entries.
- Add a history mode (`p --log`) that reuses the selector UI to list ledger entries (session · action · timestamp · cwd→target). Selecting an entry replays the recorded launch after re-validating the target directory.

## Consequences

### Positive

- Duplicate session collisions now result in a seamless attach instead of an error when the directories match.
- Users gain a durable, inspectable history of launches for post-crash recovery.
- The history selector keeps the UI consistent and supports fuzzy filtering of past sessions.

### Negative

- Maintaining the ledger introduces file-locking and atomic-write complexity.
- The session log consumes disk space (bounded to 200 entries but still present) and depends on filesystem semantics.
- CLI argument parsing grew to handle multiple modes/flags.

## Alternatives Considered

1. **Keep duplicate errors** – rejected because attaching automatically removes friction in the common case of re-opening the same directory.
2. **Store history in tmux environment variables** – rejected due to volatility and lack of persistence across crashes.
3. **Plain-text log with no locking** – rejected to avoid corruption when multiple `p` processes run concurrently.

## Related

- Planning: `.plan/.done/feat-session-enhancements/`
