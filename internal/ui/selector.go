package ui

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/wilmoore/p/internal/tmux"
)

// item represents a selectable item in the UI.
type item struct {
	label     string
	isSession bool
	session   *tmux.Session
	directory *tmux.Directory
}

// ShowSelector displays the session/directory selector and returns the user's choice.
// Supports drill-down navigation into subdirectories.
func ShowSelector(sessions []tmux.Session, dirs []tmux.Directory) (tmux.Choice, error) {
	// Navigation stack to track drill-down path
	var navStack []*tmux.Directory

	// Current directories being displayed
	currentDirs := dirs

	for {
		fmt.Println()

		// Build combined list for selection
		var items []item

		// Show navigation context if we're drilled down
		if len(navStack) > 0 {
			current := navStack[len(navStack)-1]
			fmt.Printf("%s/\n\n", current.Name)
		}

		// Only show sessions at top level
		if len(navStack) == 0 && len(sessions) > 0 {
			fmt.Println("Sessions:")
			for i := range sessions {
				idx := len(items) + 1
				fmt.Printf("  [%d] %s\n", idx, sessions[i].Name)
				items = append(items, item{
					label:     sessions[i].Name,
					isSession: true,
					session:   &sessions[i],
				})
			}
			fmt.Println()
		}

		// Add directories (or subdirectories if drilled down)
		if len(currentDirs) > 0 {
			if len(navStack) == 0 {
				fmt.Println("Projects:")
			}
			for i := range currentDirs {
				idx := len(items) + 1
				indicator := ""
				if currentDirs[i].HasSubdirs {
					indicator = " >"
				}
				fmt.Printf("  [%d] %s%s\n", idx, currentDirs[i].Name, indicator)
				items = append(items, item{
					label:     currentDirs[i].Name,
					directory: &currentDirs[i],
				})
			}
			fmt.Println()
		}

		// Show back option if drilled down
		if len(navStack) > 0 {
			fmt.Println("  [..] up")
			fmt.Println()
		}

		if len(items) == 0 && len(navStack) == 0 {
			return tmux.Choice{}, fmt.Errorf("no sessions or projects available")
		}

		// Read user selection
		fmt.Print("> ")
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			return tmux.Choice{}, fmt.Errorf("failed to read input: %w", err)
		}

		input = strings.TrimSpace(input)
		if input == "" {
			return tmux.Choice{}, fmt.Errorf("no selection made")
		}

		// Handle back navigation
		if input == ".." {
			if len(navStack) > 0 {
				navStack = navStack[:len(navStack)-1]
				if len(navStack) > 0 {
					// Get subdirs of parent
					parent := navStack[len(navStack)-1]
					currentDirs, _ = tmux.GetSubdirectories(parent)
				} else {
					// Back to top level
					currentDirs = dirs
				}
				continue
			}
			// Already at top level, ignore
			continue
		}

		// Parse selection
		idx, err := strconv.Atoi(input)
		if err != nil || idx < 1 || idx > len(items) {
			fmt.Printf("Invalid selection: %s\n", input)
			continue
		}

		selected := items[idx-1]

		// If it's a session, return immediately
		if selected.isSession {
			return tmux.Choice{
				IsSession: true,
				Session:   selected.session,
			}, nil
		}

		// It's a directory
		dir := selected.directory

		// If directory has subdirectories, offer drill-down
		if dir.HasSubdirs {
			action, err := promptDrillDown(dir)
			if err != nil {
				return tmux.Choice{}, err
			}

			switch action {
			case "drill":
				// Drill down into this directory
				navStack = append(navStack, dir)
				currentDirs, err = tmux.GetSubdirectories(dir)
				if err != nil {
					return tmux.Choice{}, fmt.Errorf("failed to read subdirectories: %w", err)
				}
				continue
			case "create":
				// Create session here
				return tmux.Choice{
					IsSession: false,
					Directory: dir,
				}, nil
			case "cancel":
				// Go back to selection
				continue
			}
		}

		// No subdirectories, create session directly
		return tmux.Choice{
			IsSession: false,
			Directory: dir,
		}, nil
	}
}

// promptDrillDown asks the user whether to drill down or create a session.
func promptDrillDown(dir *tmux.Directory) (string, error) {
	fmt.Println()
	fmt.Printf("%s contains subdirectories:\n\n", filepath.Base(dir.Path))
	fmt.Println("  [d] drill down")
	fmt.Println("  [c] create session here")
	fmt.Println("  [q] cancel")
	fmt.Println()
	fmt.Print("> ")

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to read input: %w", err)
	}

	input = strings.TrimSpace(strings.ToLower(input))
	switch input {
	case "d", "drill":
		return "drill", nil
	case "c", "create":
		return "create", nil
	case "q", "cancel", "":
		return "cancel", nil
	default:
		return "cancel", nil
	}
}
