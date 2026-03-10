> coderabbit --prompt-only

Review this refactor for `p --log` UI readability and coding standards compliance.

Key changes:

- Extracted shared interactive selector logic into `internal/ui/selector_engine.go` and shared key handling into `internal/ui/keys.go`.
- Reworked history row formatting to be terminal-width aware and to avoid overflowing lines (`internal/ui/history_selector.go`).
- Centralized UI text/ANSI sequences into `internal/ui/text.go` and `internal/ui/ansi.go`.
- Added user-friendly error formatting with opt-in debug detail via `P_DEBUG` (`internal/clierr/clierr.go`) and moved some CLI messages to `internal/i18n/messages.go`.
- Added unit tests for truncation/viewport/row width (`internal/ui/format_test.go`).

Please check:

- No behavioral regressions in selector navigation (Esc/Ctrl+C, arrows, Ctrl+J/K/N/P, Enter, Backspace).
- History row formatting stays within terminal width and truncation looks reasonable.
- Error formatting is clear without `P_DEBUG`, and helpful with `P_DEBUG=1`.
- Code style: magic numbers/strings minimized; duplicated logic reduced; naming and constants make intent clear.
