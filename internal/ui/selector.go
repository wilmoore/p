package ui

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/wilmooreiii/p/internal/tmux"
)

// ShowSelector displays the session/directory selector and returns the user's choice.
func ShowSelector(sessions []tmux.Session, dirs []tmux.Directory) (tmux.Choice, error) {
	fmt.Println()

	// Build combined list for selection
	type item struct {
		label     string
		isSession bool
		session   *tmux.Session
		directory *tmux.Directory
	}
	var items []item

	// Add sessions first
	if len(sessions) > 0 {
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

	// Add directories
	if len(dirs) > 0 {
		fmt.Println("Projects:")
		for i := range dirs {
			idx := len(items) + 1
			fmt.Printf("  [%d] %s\n", idx, dirs[i].Name)
			items = append(items, item{
				label:     dirs[i].Name,
				isSession: false,
				directory: &dirs[i],
			})
		}
		fmt.Println()
	}

	if len(items) == 0 {
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

	// Parse selection
	idx, err := strconv.Atoi(input)
	if err != nil || idx < 1 || idx > len(items) {
		return tmux.Choice{}, fmt.Errorf("invalid selection: %s", input)
	}

	selected := items[idx-1]
	return tmux.Choice{
		IsSession: selected.isSession,
		Session:   selected.session,
		Directory: selected.directory,
	}, nil
}
