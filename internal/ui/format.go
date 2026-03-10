package ui

import "strings"

const ellipsis = "..."

func truncateLeft(s string, max int) string {
	if max <= 0 {
		return ""
	}
	if len(s) <= max {
		return s
	}
	if max <= len(ellipsis) {
		return s[len(s)-max:]
	}
	return ellipsis + s[len(s)-(max-len(ellipsis)):]
}

func truncateRight(s string, max int) string {
	if max <= 0 {
		return ""
	}
	if len(s) <= max {
		return s
	}
	if max <= len(ellipsis) {
		return s[:max]
	}
	return s[:max-len(ellipsis)] + ellipsis
}

func spaces(n int) string {
	if n <= 0 {
		return ""
	}
	return strings.Repeat(" ", n)
}
