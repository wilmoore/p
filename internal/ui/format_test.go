package ui

import (
	"testing"
	"time"

	"github.com/wilmoore/p/internal/history"
)

func TestTruncateLeft(t *testing.T) {
	s := "abcdefghijklmnopqrstuvwxyz"
	got := truncateLeft(s, 10)
	want := "...tuvwxyz"
	if got != want {
		t.Fatalf("truncateLeft: got %q want %q", got, want)
	}
}

func TestTruncateRight(t *testing.T) {
	s := "abcdefghijklmnopqrstuvwxyz"
	got := truncateRight(s, 10)
	want := "abcdefg..."
	if got != want {
		t.Fatalf("truncateRight: got %q want %q", got, want)
	}
}

func TestVisibleRange(t *testing.T) {
	start, end := visibleRange(100, 0, 10)
	if start != 0 || end != 10 {
		t.Fatalf("start/end mismatch: got %d/%d", start, end)
	}

	start, end = visibleRange(100, 50, 10)
	if start != 45 || end != 55 {
		t.Fatalf("start/end mismatch: got %d/%d", start, end)
	}

	start, end = visibleRange(100, 99, 10)
	if start != 90 || end != 100 {
		t.Fatalf("start/end mismatch: got %d/%d", start, end)
	}
}

func TestFormatHistoryRowDoesNotExceedWidth(t *testing.T) {
	entry := history.Entry{
		Timestamp:   time.Date(2026, 3, 10, 11, 52, 0, 0, time.Local),
		Action:      history.ActionAttachExisting,
		SessionName: "book-from-developer-to-ceo-with-an-excessively-long-session-name",
		InvokeDir:   "/Users/example/Documents/src/really/long/path/that/should/be/truncated/on/render",
		TargetDir:   "/Users/example/Documents/src/another/really/long/path/that/should/be/truncated/on/render",
	}

	for _, width := range []int{120, 80, 60, 40, 20} {
		row := formatHistoryRow(entry, width)
		if len(row) > width {
			t.Fatalf("row exceeds width %d: len=%d row=%q", width, len(row), row)
		}
	}
}
