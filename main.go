package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/wilmoore/p/internal/tmux"
	"github.com/wilmoore/p/internal/ui"
)

// Version is set at build time via -ldflags "-X main.Version=vX.Y.Z"
var Version = "dev"

const usage = `p - minimal tmux session switcher

Usage:
  p              Show interactive session selector
  p <path>       Create new session in directory (use . for current directory)
  p --version    Show version information
  p --help       Show this help message

Navigation:
  Type           Filter sessions by name
  Arrow keys     Navigate up/down
  Enter          Attach to selected session
  Esc/Ctrl+C     Cancel

Examples:
  p              Select from existing sessions
  p .            Create session in current directory
  p ~/projects   Create session in ~/projects
`

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {
	args := os.Args[1:]

	// Handle flags
	if len(args) > 0 {
		switch args[0] {
		case "--version", "-v":
			fmt.Println(Version)
			return nil
		case "--help", "-h":
			fmt.Print(usage)
			return nil
		}

		// Treat argument as path for new session
		if args[0] != "" && args[0][0] != '-' {
			return createSessionFromPath(args[0])
		}

		// Unknown flag
		return fmt.Errorf("unknown option: %s\nRun 'p --help' for usage", args[0])
	}

	// No arguments: show session selector
	sessions, err := tmux.ListSessions()
	if err != nil && !tmux.IsNoServerError(err) {
		return fmt.Errorf("failed to list tmux sessions: %w", err)
	}

	if len(sessions) == 0 {
		return fmt.Errorf("no tmux sessions available")
	}

	// Show selector and get user choice
	choice, err := ui.ShowSelector(sessions)
	if err != nil {
		return err
	}

	// User cancelled (Ctrl+C, Esc, q)
	if choice == nil {
		return nil
	}

	// Attach to the selected session
	return tmux.AttachToSession(choice.Name)
}

// createSessionFromPath creates a new tmux session in the specified directory.
func createSessionFromPath(path string) error {
	// Resolve path
	var absPath string
	var err error

	if path == "." {
		absPath, err = os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}
	} else {
		// Expand ~ to home directory
		if len(path) > 0 && path[0] == '~' {
			home, err := os.UserHomeDir()
			if err != nil {
				return fmt.Errorf("failed to get home directory: %w", err)
			}
			path = filepath.Join(home, path[1:])
		}

		absPath, err = filepath.Abs(path)
		if err != nil {
			return fmt.Errorf("failed to resolve path: %w", err)
		}
	}

	// Verify directory exists
	info, err := os.Stat(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("directory does not exist: %s", absPath)
		}
		return fmt.Errorf("failed to stat path: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("not a directory: %s", absPath)
	}

	// Derive session name from directory basename
	sessionName := filepath.Base(absPath)

	// Create and attach to the session
	return tmux.CreateSession(sessionName, absPath)
}
