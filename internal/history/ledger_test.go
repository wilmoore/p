package history

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"
)

func TestAppendAndListTrimsOldEntries(t *testing.T) {
	logDir := t.TempDir()
	path := filepath.Join(logDir, "session-log.jsonl")
	t.Setenv("P_HISTORY_PATH", path)

	for i := 0; i < 210; i++ {
		entry := Entry{
			Timestamp:   time.Unix(int64(i), 0),
			Action:      ActionCreate,
			SessionName: "s" + strconv.Itoa(i),
			InvokeDir:   "/tmp",
			TargetDir:   "/tmp/project",
		}
		if err := Append(entry); err != nil {
			t.Fatalf("append failed: %v", err)
		}
	}

	entries, err := List(0)
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}
	if len(entries) != 200 {
		t.Fatalf("expected 200 entries, got %d", len(entries))
	}
	if entries[0].SessionName != "s209" {
		t.Fatalf("expected newest entry first")
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat log: %v", err)
	}
	if info.Size() == 0 {
		t.Fatalf("log file should not be empty")
	}
}
