package ui

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/wilmoore/p/internal/history"
	"golang.org/x/sys/unix"
	"golang.org/x/term"
)

// ShowHistory renders the history selector UI.
func ShowHistory(entries []history.Entry) (*history.Entry, error) {
	if len(entries) == 0 {
		return nil, fmt.Errorf("no history entries")
	}
	items := make([]historyItem, len(entries))
	for i, entry := range entries {
		items[i] = newHistoryItem(entry)
	}
	return historySelector(items)
}

type historyItem struct {
	entry   history.Entry
	display string
	summary string
	search  string
}

func newHistoryItem(entry history.Entry) historyItem {
	stamp := entry.Timestamp.Local().Format("01/02 15:04")
	left := shorten(entry.InvokeDir)
	right := shorten(entry.TargetDir)
	display := fmt.Sprintf("%-18s %-14s %s %s -> %s", entry.SessionName, entry.Action, stamp, left, right)
	summary := fmt.Sprintf("%s from %s to %s at %s", entry.Action, entry.InvokeDir, entry.TargetDir, entry.Timestamp.Format(time.RFC3339))
	searchParts := []string{entry.SessionName, string(entry.Action), entry.InvokeDir, entry.TargetDir, display, stamp}
	search := strings.ToLower(strings.Join(searchParts, " "))
	return historyItem{entry: entry, display: display, summary: summary, search: search}
}

func historySelector(items []historyItem) (*history.Entry, error) {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return nil, fmt.Errorf("failed to set raw mode: %w", err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	fmt.Print("\033[?1049h")
	defer fmt.Print("\033[?1049l")

	var query string
	selected := 0
	for {
		filtered := filterHistory(items, query)
		if selected >= len(filtered) {
			selected = len(filtered) - 1
		}
		if selected < 0 {
			selected = 0
		}
		renderHistory(filtered, query, selected)

		b := make([]byte, 3)
		n, err := os.Stdin.Read(b)
		if err != nil {
			return nil, fmt.Errorf("failed to read input: %w", err)
		}

		switch {
		case n == 1 && b[0] == 3:
			return nil, nil
		case n == 1 && b[0] == 27:
			handled, err := handleEscapeSequence(&selected, len(filtered))
			if err != nil {
				return nil, err
			}
			if !handled {
				return nil, nil
			}
		case n == 1 && b[0] == 13:
			if len(filtered) > 0 {
				entry := filtered[selected].entry
				return &entry, nil
			}
		case n == 1 && b[0] == 127:
			if len(query) > 0 {
				query = query[:len(query)-1]
				selected = 0
			}
		case n == 1 && (b[0] == 10 || b[0] == 14):
			if selected < len(filtered)-1 {
				selected++
			}
		case n == 1 && (b[0] == 11 || b[0] == 16):
			if selected > 0 {
				selected--
			}
		case n == 3 && b[0] == 27 && b[1] == 91:
			switch b[2] {
			case 65:
				if selected > 0 {
					selected--
				}
			case 66:
				if selected < len(filtered)-1 {
					selected++
				}
			}
		case n == 1 && b[0] >= 32 && b[0] < 127:
			query += string(b[0])
			selected = 0
		}
	}
}

func filterHistory(items []historyItem, query string) []historyItem {
	if query == "" {
		return items
	}
	q := strings.ToLower(query)
	var filtered []historyItem
	for _, item := range items {
		if strings.Contains(item.search, q) {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

func renderHistory(items []historyItem, query string, selected int) {
	fmt.Print("\033[H\033[J")
	fmt.Print("Session History:\r\n\r\n")
	if len(items) == 0 {
		fmt.Print("  (no matches)\r\n\r\n")
		fmt.Printf("> %s", query)
		return
	}
	for i, item := range items {
		if i == selected {
			fmt.Printf("  \033[7m %s \033[0m\r\n", item.display)
		} else {
			fmt.Printf("   %s\r\n", item.display)
		}
	}
	fmt.Print("\r\n")
	if selected >= 0 && selected < len(items) {
		fmt.Printf("Selected: %s\r\n", items[selected].summary)
	} else {
		fmt.Print("Selected: -\r\n")
	}
	fmt.Printf("> %s", query)
}

func shorten(path string) string {
	if path == "" {
		return "-"
	}
	const max = 24
	if len(path) <= max {
		return path
	}
	return "..." + path[len(path)-max+3:]
}

func handleEscapeSequence(selected *int, length int) (bool, error) {
	fd := int(os.Stdin.Fd())
	ready, err := hasPendingInput(fd)
	if err != nil {
		return false, err
	}
	if !ready {
		return false, nil
	}
	extra := make([]byte, 2)
	total := 0
	for total < len(extra) {
		n, err := os.Stdin.Read(extra[total:])
		if err != nil {
			return false, err
		}
		total += n
		if n == 0 {
			break
		}
	}
	if total >= 2 && extra[0] == '[' {
		switch extra[1] {
		case 'A':
			if *selected > 0 {
				*selected--
			}
			return true, nil
		case 'B':
			if *selected < length-1 {
				*selected++
			}
			return true, nil
		}
	}
	return false, nil
}

func hasPendingInput(fd int) (bool, error) {
	var set unix.FdSet
	fdSet(fd, &set)
	var tv unix.Timeval
	n, err := unix.Select(fd+1, &set, nil, nil, &tv)
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

func fdSet(fd int, set *unix.FdSet) {
	set.Bits[fd/64] |= 1 << (uint(fd) % 64)
}
