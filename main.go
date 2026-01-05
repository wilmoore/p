package main

import (
	"fmt"
	"os"

	"github.com/wilmooreiii/p/internal/tmux"
	"github.com/wilmooreiii/p/internal/ui"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {
	// Check CDPATH first
	cdpath := os.Getenv("CDPATH")
	if cdpath == "" {
		return fmt.Errorf("CDPATH is not set.\n\np uses CDPATH to discover project directories.\nSet CDPATH to enable directory-based projects.\n\nExiting")
	}

	// Get existing tmux sessions
	sessions, err := tmux.ListSessions()
	if err != nil && !tmux.IsNoServerError(err) {
		return fmt.Errorf("failed to list tmux sessions: %w", err)
	}

	// Get directories from CDPATH
	dirs, err := tmux.DiscoverDirectories(cdpath)
	if err != nil {
		return fmt.Errorf("failed to discover directories: %w", err)
	}

	// Show selector and get user choice
	choice, err := ui.ShowSelector(sessions, dirs)
	if err != nil {
		return err
	}

	// Execute the user's choice
	return tmux.Execute(choice)
}
