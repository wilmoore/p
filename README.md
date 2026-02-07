# p

[![Go Reference](https://pkg.go.dev/badge/github.com/wilmoore/p.svg)](https://pkg.go.dev/github.com/wilmoore/p)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

Minimal tmux session switcher. List sessions, filter, attach.

## Installation

```bash
go install github.com/wilmoore/p@latest
```

Requires Go 1.24+ and tmux.

## Usage

```bash
p
```

Shows an interactive selector of existing tmux sessions:

```
Sessions:

   api-server
   frontend
 ▸ my-project

> _
```

- **Type** to filter sessions
- **Arrow keys** to navigate
- **Enter** to attach (or switch if already inside tmux)
- **Esc** or **Ctrl+C** to cancel

## How It Works

1. Lists existing tmux sessions
2. Presents an interactive selector with type-to-filter
3. Attaches to the selected session (or switches if already in tmux)

All tmux invocations use `-f /dev/null` to ignore global configuration.

## License

[MIT](LICENSE)
