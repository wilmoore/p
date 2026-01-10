# p

[![Go Reference](https://pkg.go.dev/badge/github.com/wilmoore/p.svg)](https://pkg.go.dev/github.com/wilmoore/p)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

Zero-configuration tmux session launcher. Discover projects from CDPATH, attach to existing sessions, and create new ones with sensible defaults.

## Features

- **Session discovery** - Lists existing tmux sessions alongside CDPATH directories
- **Drill-down navigation** - Navigate into subdirectories before creating a session
- **Collision-free naming** - Deterministic session names with hash suffixes when needed
- **Vi-style copy mode** - `v` to select, `y` to yank (pre-configured)
- **Git branch in status bar** - Shows current branch, updates every 5 seconds
- **Zero config** - Runs tmux with `-f /dev/null`, injects only essential defaults

## Installation

```bash
go install github.com/wilmoore/p@latest
```

Requires Go 1.22+ and tmux.

## Usage

Set `CDPATH` to your project directories and run `p`:

```bash
export CDPATH="$HOME/projects:$HOME/work"
p
```

The selector shows existing sessions and discovered directories:

```
Sessions:
  [1] api-server
  [2] frontend

Projects:
  [3] new-project >
  [4] another-project
```

- Select a session to attach (or switch if already inside tmux)
- Select a directory to create a new session
- Directories with `>` contain subdirectories you can drill into

### Drill-down

When you select a directory with subdirectories:

```
new-project contains subdirectories:

  [d] drill down
  [c] create session here
  [q] cancel
```

Navigate as deep as you want, then create a session at any level.

## How It Works

1. Reads `CDPATH` to discover project directories
2. Lists existing tmux sessions
3. Presents an interactive selector
4. Creates sessions with deterministic names derived from the directory path
5. Injects minimal config: vi copy mode, status bar with git branch
6. Attaches (or switches) to the selected session

All tmux invocations use `-f /dev/null` to ignore global configuration.

## License

[MIT](LICENSE)
