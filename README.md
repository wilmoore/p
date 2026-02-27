<p align="center">
  <img src="assets/hero.png" alt="p" width="120">
  <h1 align="center">p</h1>
  <p align="center">
    <strong>Minimal tmux session switcher</strong>
  </p>
  <p align="center">
    <a href="https://pkg.go.dev/github.com/wilmoore/p"><img src="https://pkg.go.dev/badge/github.com/wilmoore/p.svg" alt="Go Reference"></a>
    <a href="LICENSE"><img src="https://img.shields.io/badge/License-MIT-blue.svg" alt="License: MIT"></a>
    <a href="https://github.com/wilmoore/p/releases"><img src="https://img.shields.io/github/v/release/wilmoore/p" alt="Release"></a>
  </p>
</p>

---

## Why p?

Managing multiple tmux sessions shouldn't require mental overhead. Most session managers are overengineered—config files, plugins, complex workflows.

**p** does one thing well: get you into the right tmux session, fast.

```
p              # Pick a session
p .            # Start here
```

That's it.

---

## Features

- **Instant session switching** — fzf-like fuzzy filtering
- **Vim/Emacs keybindings** — `Ctrl+J/K` or `Ctrl+N/P` navigation
- **Zero configuration** — works immediately, ignores `~/.tmux.conf`
- **Styled by default** — dark theme with git branch in status bar
- **Vi copy mode** — `v` to select, `y` to yank (built-in)

---

## Installation

```bash
go install github.com/wilmoore/p@latest
```

**Requirements:** Go 1.24+, tmux

---

## Usage

### Quick Reference

```bash
p                  # Interactive session selector
p .                # Create session in current directory
p ~/projects/app   # Create session in specific directory
p --help           # Show help
p --version        # Show version
```

### Interactive Session Selector

Running `p` displays your tmux sessions:

```
Sessions:

   api-server
   frontend
   my-project      ← highlighted

> _
```

### Keybindings

| Key | Action |
|-----|--------|
| **Type** | Filter sessions |
| `↑` `↓` | Navigate |
| `Ctrl+K` / `Ctrl+P` | Navigate up (vim/emacs) |
| `Ctrl+J` / `Ctrl+N` | Navigate down (vim/emacs) |
| `Enter` | Attach to session |
| `Esc` / `Ctrl+C` / `q` | Cancel |

### Create Session from Directory

```bash
p .              # Session named after current directory
p ~/projects     # Session named "projects"
```

Session names are derived from the directory basename. If you're already inside tmux, `p` seamlessly switches clients.

---

## How It Works

1. **Lists** existing tmux sessions
2. **Filters** as you type (case-insensitive)
3. **Attaches** to selection (or switches if inside tmux)

All tmux commands use `-f /dev/null` to bypass your config. This ensures consistent behavior everywhere.

---

## tmux Configuration

`p` injects sensible defaults into every session:

### Status Bar
- Dark gray background with git branch display
- Session name on the left, current branch on the right
- Clean, minimal aesthetic

### Vi Copy Mode
Once in a `p`-managed session, copy mode uses vi bindings:
- `v` — begin selection
- `y` — copy and exit

These work regardless of your `~/.tmux.conf`.

---

## Development

```bash
make help      # Show all targets
make dev       # Build and install locally
make check     # Run fmt, vet, test
make lint      # Run golangci-lint
```

### Project Structure

```
p/
├── main.go              # CLI entry point
├── internal/
│   ├── tmux/            # tmux session management
│   │   ├── list.go      # Session listing
│   │   └── execute.go   # Attach/create/configure
│   └── ui/
│       └── selector.go  # Interactive picker
└── doc/
    └── decisions/       # Architecture Decision Records
```

---

## Philosophy

**p** follows these principles:

1. **Zero config** — No dotfiles, no setup, no surprises
2. **Fast startup** — Single binary, instant launch
3. **Predictable** — Same behavior on every machine
4. **Minimal** — One job, done well

---

## License

[MIT](LICENSE)
