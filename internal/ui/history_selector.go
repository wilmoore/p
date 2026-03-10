package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/wilmoore/p/internal/history"
)

// ShowHistory renders the history selector UI.
func ShowHistory(entries []history.Entry) (*history.Entry, error) {
	if len(entries) == 0 {
		return nil, fmt.Errorf("no history entries")
	}

	adapter := selectorAdapter[history.Entry]{
		title:        uiTitleHistory,
		emptyMessage: uiNoMatches,
		renderRow:    formatHistoryRow,
		summary:      formatHistorySummary,
		searchText:   historySearchText,
	}
	return runSelector(entries, adapter)
}

const (
	historyStampLayout = "01/02 15:04"
	historyArrow       = " -> "

	historyActionWidth = 14 // len("attach-existing")
	historyStampWidth  = 11 // len("01/02 15:04")

	historyMinSessionWidth       = 12
	historyMaxSessionWidth       = 32
	historyMaxSessionWidthNarrow = 20
	historyNarrowWidthThreshold  = 60

	historySummaryPathNarrow         = 24
	historySummaryPathWide           = 40
	historySummaryWideWidthThreshold = 100
)

func historySearchText(entry history.Entry) string {
	stamp := entry.Timestamp.Local().Format(historyStampLayout)
	parts := []string{
		entry.SessionName,
		string(entry.Action),
		entry.InvokeDir,
		entry.TargetDir,
		stamp,
	}
	return strings.Join(parts, " ")
}

func formatHistoryRow(entry history.Entry, width int) string {
	if width <= 0 {
		return ""
	}

	stamp := entry.Timestamp.Local().Format(historyStampLayout)
	session := entry.SessionName
	action := string(entry.Action)

	left := entry.InvokeDir
	if left == "" {
		left = "-"
	}
	right := entry.TargetDir
	if right == "" {
		right = "-"
	}

	minSession := historyMinSessionWidth
	maxSession := historyMaxSessionWidth
	if width < historyNarrowWidthThreshold {
		maxSession = historyMaxSessionWidthNarrow
	}

	// Fixed-width portion: session + action + stamp + separating spaces.
	fixed := historyActionWidth + historyStampWidth + 6
	pathAvail := width - fixed
	if pathAvail < minSession {
		pathAvail = minSession
	}

	sessionWidth := width - (historyActionWidth + historyStampWidth + 4) - pathAvail
	if sessionWidth < minSession {
		sessionWidth = minSession
	}
	if sessionWidth > maxSession {
		sessionWidth = maxSession
	}

	pathAvail = width - (sessionWidth + historyActionWidth + historyStampWidth + 4)
	if pathAvail < 0 {
		pathAvail = 0
	}

	fromWidth := (pathAvail - len(historyArrow)) / 2
	toWidth := pathAvail - len(historyArrow) - fromWidth
	if fromWidth < 0 {
		fromWidth = 0
	}
	if toWidth < 0 {
		toWidth = 0
	}

	from := truncateLeft(left, fromWidth)
	to := truncateLeft(right, toWidth)
	path := ""
	if pathAvail > 0 {
		path = from + historyArrow + to
	}

	row := fmt.Sprintf(
		"%-*s %-*s %-*s %s",
		sessionWidth,
		truncateRight(session, sessionWidth),
		historyActionWidth,
		truncateRight(action, historyActionWidth),
		historyStampWidth,
		stamp,
		path,
	)
	return truncateRight(row, width)
}

func formatHistorySummary(entry history.Entry, width int) string {
	stamp := entry.Timestamp.Format(time.RFC3339)
	left := entry.InvokeDir
	if left == "" {
		left = "-"
	}
	right := entry.TargetDir
	if right == "" {
		right = "-"
	}

	// Keep the summary concise and width-friendly.
	maxPath := historySummaryPathNarrow
	if width > historySummaryWideWidthThreshold {
		maxPath = historySummaryPathWide
	}
	from := truncateLeft(left, maxPath)
	to := truncateLeft(right, maxPath)

	return fmt.Sprintf(
		"%s %s from %s to %s at %s",
		entry.Action,
		entry.SessionName,
		from,
		to,
		stamp,
	)
}
