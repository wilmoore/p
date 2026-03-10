package ui

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

type selectorAdapter[T any] struct {
	title        string
	emptyMessage string
	renderRow    func(item T, width int) string
	summary      func(item T, width int) string
	searchText   func(item T) string
	directSelect func(query string) (*T, bool)
}

type selectorItem[T any] struct {
	value  T
	search string
}

const (
	selectorIndent              = 3
	selectorReservedNoSummary   = 4
	selectorReservedWithSummary = 6

	defaultTermWidth  = 80
	defaultTermHeight = 24
)

func runSelector[T any](items []T, adapter selectorAdapter[T]) (*T, error) {
	if len(items) == 0 {
		return nil, fmt.Errorf("no items")
	}

	prepared := make([]selectorItem[T], len(items))
	for i, it := range items {
		prepared[i] = selectorItem[T]{
			value:  it,
			search: strings.ToLower(adapter.searchText(it)),
		}
	}

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return nil, fmt.Errorf("failed to set raw mode: %w", err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	fmt.Print(ansiEnterAltScreen)
	defer fmt.Print(ansiExitAltScreen)

	query := ""
	selected := 0
	for {
		filtered := filterSelectorItems(prepared, query)
		selected = clampSelected(selected, len(filtered))
		renderSelector(filtered, query, selected, adapter)

		ev, err := readKeyEvent()
		if err != nil {
			return nil, fmt.Errorf("failed to read input: %w", err)
		}
		switch ev.kind {
		case keyCancel:
			return nil, nil
		case keyEnter:
			if len(filtered) > 0 {
				chosen := filtered[selected].value
				return &chosen, nil
			}
		case keyBackspace:
			if len(query) > 0 {
				query = query[:len(query)-1]
				selected = 0
			}
		case keyDown:
			if selected < len(filtered)-1 {
				selected++
			}
		case keyUp:
			if selected > 0 {
				selected--
			}
		case keyRune:
			query += string(ev.r)
			selected = 0
			if adapter.directSelect != nil {
				if chosen, ok := adapter.directSelect(query); ok {
					return chosen, nil
				}
			}
		}
	}
}

func filterSelectorItems[T any](items []selectorItem[T], query string) []selectorItem[T] {
	if query == "" {
		return items
	}
	q := strings.ToLower(query)
	filtered := make([]selectorItem[T], 0, len(items))
	for _, it := range items {
		if strings.Contains(it.search, q) {
			filtered = append(filtered, it)
		}
	}
	return filtered
}

func clampSelected(selected int, length int) int {
	if length <= 0 {
		return 0
	}
	if selected >= length {
		return length - 1
	}
	if selected < 0 {
		return 0
	}
	return selected
}

type termSize struct {
	width  int
	height int
}

func currentTermSize() termSize {
	w, h, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || w <= 0 || h <= 0 {
		return termSize{width: defaultTermWidth, height: defaultTermHeight}
	}
	return termSize{width: w, height: h}
}

func renderSelector[T any](items []selectorItem[T], query string, selected int, adapter selectorAdapter[T]) {
	size := currentTermSize()
	indent := selectorIndent

	fmt.Print(ansiClearScreen)
	fmt.Print(adapter.title)
	fmt.Print(crlf)
	fmt.Print(crlf)

	reserved := selectorReservedNoSummary
	if adapter.summary != nil {
		reserved = selectorReservedWithSummary
	}
	maxRows := size.height - reserved
	if maxRows < 1 {
		maxRows = 1
	}

	if len(items) == 0 {
		empty := adapter.emptyMessage
		if empty == "" {
			empty = uiNoMatches
		}
		fmt.Print(empty)
		fmt.Print(crlf)
		fmt.Print(crlf)
		fmt.Print(uiPrompt)
		fmt.Print(query)
		return
	}

	start, end := visibleRange(len(items), selected, maxRows)
	for i := start; i < end; i++ {
		rowWidth := size.width - indent
		if rowWidth < 0 {
			rowWidth = 0
		}
		line := adapter.renderRow(items[i].value, rowWidth)
		if i == selected {
			fmt.Print(spaces(indent - 1))
			fmt.Print(ansiInvert)
			fmt.Print(" ")
			fmt.Print(truncateRight(line, rowWidth))
			fmt.Print(" ")
			fmt.Print(ansiReset)
			fmt.Print(crlf)
		} else {
			fmt.Print(spaces(indent))
			fmt.Print(truncateRight(line, rowWidth))
			fmt.Print(crlf)
		}
	}

	if adapter.summary != nil {
		fmt.Print(crlf)
		if selected >= 0 && selected < len(items) {
			label := adapter.summary(items[selected].value, size.width)
			if label == "" {
				label = uiSelectedNone
			}
			max := size.width - len(uiSelected) - 1
			if max < 0 {
				max = 0
			}
			fmt.Print(uiSelected)
			fmt.Print(" ")
			fmt.Print(truncateRight(label, max))
			fmt.Print(crlf)
		} else {
			fmt.Print(uiSelected)
			fmt.Print(" ")
			fmt.Print(uiSelectedNone)
			fmt.Print(crlf)
		}
	}
	fmt.Print(uiPrompt)
	fmt.Print(query)
}

func visibleRange(total, selected, maxRows int) (int, int) {
	if total <= 0 {
		return 0, 0
	}
	if maxRows <= 0 {
		maxRows = 1
	}
	if total <= maxRows {
		return 0, total
	}
	half := maxRows / 2
	start := selected - half
	if start < 0 {
		start = 0
	}
	end := start + maxRows
	if end > total {
		end = total
		start = end - maxRows
		if start < 0 {
			start = 0
		}
	}
	return start, end
}
