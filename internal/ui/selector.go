package ui

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/wilmoore/p/internal/tmux"
	"golang.org/x/term"
)

// ShowSelector displays an fzf-like session selector.
// Supports both numeric selection and text filtering.
func ShowSelector(sessions []tmux.Session) (*tmux.Session, error) {
	if len(sessions) == 0 {
		return nil, fmt.Errorf("no sessions available")
	}

	// Put terminal in raw mode
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return nil, fmt.Errorf("failed to set raw mode: %w", err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	// Enter alternate screen buffer for clean display
	fmt.Print("\033[?1049h")
	defer fmt.Print("\033[?1049l")

	var query string
	selected := 0

	for {
		// Filter sessions based on query
		filtered := filterSessions(sessions, query)

		// Clamp selection
		if selected >= len(filtered) {
			selected = len(filtered) - 1
		}
		if selected < 0 {
			selected = 0
		}

		// Render
		render(filtered, query, selected)

		// Read keypress
		b := make([]byte, 3)
		n, err := os.Stdin.Read(b)
		if err != nil {
			return nil, fmt.Errorf("failed to read input: %w", err)
		}

		// Handle input
		switch {
		case n == 1 && b[0] == 3: // Ctrl+C
			return nil, nil // Clean exit

		case n == 1 && b[0] == 27: // Escape
			return nil, nil // Clean exit

		case n == 1 && b[0] == 'q' && query == "": // q to quit (only when not filtering)
			return nil, nil // Clean exit

		case n == 1 && b[0] == 13: // Enter
			if len(filtered) > 0 {
				return &filtered[selected], nil
			}

		case n == 1 && b[0] == 127: // Backspace
			if len(query) > 0 {
				query = query[:len(query)-1]
				selected = 0
			}

		case n == 1 && (b[0] == 10 || b[0] == 14): // Ctrl+J (0x0A) or Ctrl+N (0x0E) - move down
			if selected < len(filtered)-1 {
				selected++
			}

		case n == 1 && (b[0] == 11 || b[0] == 16): // Ctrl+K (0x0B) or Ctrl+P (0x10) - move up
			if selected > 0 {
				selected--
			}

		case n == 3 && b[0] == 27 && b[1] == 91: // Arrow keys
			switch b[2] {
			case 65: // Up
				if selected > 0 {
					selected--
				}
			case 66: // Down
				if selected < len(filtered)-1 {
					selected++
				}
			}

		case n == 1 && b[0] >= 32 && b[0] < 127: // Printable character
			query += string(b[0])
			selected = 0

			// Check if query is a valid number for direct selection
			if idx, err := strconv.Atoi(query); err == nil {
				if idx >= 1 && idx <= len(sessions) {
					return &sessions[idx-1], nil
				}
			}
		}
	}
}

// filterSessions returns sessions matching the query (case-insensitive substring).
func filterSessions(sessions []tmux.Session, query string) []tmux.Session {
	if query == "" {
		return sessions
	}

	query = strings.ToLower(query)
	var filtered []tmux.Session
	for _, s := range sessions {
		if strings.Contains(strings.ToLower(s.Name), query) {
			filtered = append(filtered, s)
		}
	}
	return filtered
}

// render displays the current state.
// Note: In raw mode, \n only moves down; \r\n is needed to return to column 0.
func render(sessions []tmux.Session, query string, selected int) {
	// Move cursor to top-left and clear screen
	fmt.Print("\033[H\033[J")

	fmt.Print("Sessions:\r\n")
	fmt.Print("\r\n")

	if len(sessions) == 0 {
		fmt.Print("  (no matches)\r\n")
	} else {
		for i, s := range sessions {
			if i == selected {
				fmt.Printf("  \033[7m %s \033[0m\r\n", s.Name)
			} else {
				fmt.Printf("   %s\r\n", s.Name)
			}
		}
	}

	fmt.Print("\r\n")
	fmt.Printf("> %s", query)
}
