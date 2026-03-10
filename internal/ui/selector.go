package ui

import (
	"fmt"
	"strconv"

	"github.com/wilmoore/p/internal/tmux"
)

// ShowSelector displays an fzf-like session selector.
// Supports both numeric selection and text filtering.
func ShowSelector(sessions []tmux.Session) (*tmux.Session, error) {
	if len(sessions) == 0 {
		return nil, fmt.Errorf("no sessions available")
	}

	adapter := selectorAdapter[tmux.Session]{
		title:        uiTitleSessions,
		emptyMessage: uiNoMatches,
		renderRow: func(s tmux.Session, width int) string {
			return truncateRight(s.Name, width)
		},
		searchText: func(s tmux.Session) string {
			return s.Name
		},
		directSelect: func(query string) (*tmux.Session, bool) {
			idx, err := strconv.Atoi(query)
			if err != nil {
				return nil, false
			}
			if idx < 1 || idx > len(sessions) {
				return nil, false
			}
			chosen := sessions[idx-1]
			return &chosen, true
		},
	}

	return runSelector(sessions, adapter)
}
