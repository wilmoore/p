package main

import (
	"fmt"
	"os"

	"github.com/wilmoore/p/internal/tmux"
	"github.com/wilmoore/p/internal/ui"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {
	// Get existing tmux sessions
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

	// Attach to the selected session
	return tmux.AttachToSession(choice.Name)
}
