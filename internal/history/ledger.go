package history

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofrs/flock"
)

// Action represents how the session was launched.
type Action string

const (
	ActionCreate         Action = "create"
	ActionAttach         Action = "attach"
	ActionAttachExisting Action = "attach-existing"
)

// Entry captures a single launch event.
type Entry struct {
	Timestamp   time.Time `json:"ts"`
	Action      Action    `json:"action"`
	SessionName string    `json:"sessionName"`
	InvokeDir   string    `json:"invokeDir"`
	TargetDir   string    `json:"targetDir"`
}

const maxEntries = 200

// Append writes a new entry to the ledger, keeping only the most recent maxEntries.
func Append(entry Entry) error {
	path, err := logFilePath()
	if err != nil {
		return err
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	lock := flock.New(path + ".lock")
	if err := lock.Lock(); err != nil {
		return err
	}
	defer lock.Unlock()

	entries, err := readEntries(path)
	if err != nil {
		return err
	}
	entries = append(entries, entry)
	if len(entries) > maxEntries {
		entries = entries[len(entries)-maxEntries:]
	}

	tmp, err := os.CreateTemp(dir, "session-log-*.tmp")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	cleanupTmp := func() {
		tmp.Close()
		os.Remove(tmpName)
	}
	writer := bufio.NewWriter(tmp)
	for _, e := range entries {
		enc, err := json.Marshal(e)
		if err != nil {
			cleanupTmp()
			return err
		}
		if _, err := writer.Write(append(enc, '\n')); err != nil {
			cleanupTmp()
			return err
		}
	}
	if err := writer.Flush(); err != nil {
		cleanupTmp()
		return err
	}
	if err := tmp.Sync(); err != nil {
		cleanupTmp()
		return err
	}
	if err := tmp.Close(); err != nil {
		os.Remove(tmpName)
		return err
	}
	if err := os.Rename(tmpName, path); err != nil {
		os.Remove(tmpName)
		return err
	}
	if dirHandle, err := os.Open(dir); err == nil {
		_ = dirHandle.Sync()
		dirHandle.Close()
	}
	return nil
}

// List returns up to limit entries, newest first.
func List(limit int) ([]Entry, error) {
	path, err := logFilePath()
	if err != nil {
		return nil, err
	}

	entries, err := readEntries(path)
	if err != nil {
		return nil, err
	}

	for i, j := 0, len(entries)-1; i < j; i, j = i+1, j-1 {
		entries[i], entries[j] = entries[j], entries[i]
	}

	if limit > 0 && len(entries) > limit {
		entries = entries[:limit]
	}
	return entries, nil
}

func readEntries(path string) ([]Entry, error) {
	file, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []Entry{}, nil
		}
		return nil, err
	}
	defer file.Close()
	return decodeEntries(file)
}

func decodeEntries(r io.Reader) ([]Entry, error) {
	scanner := bufio.NewScanner(r)
	var entries []Entry
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var entry Entry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			return nil, fmt.Errorf("failed to decode history entry: %w", err)
		}
		entries = append(entries, entry)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return entries, nil
}

func logFilePath() (string, error) {
	if override := os.Getenv("P_HISTORY_PATH"); override != "" {
		return override, nil
	}

	base := os.Getenv("XDG_STATE_HOME")
	if base == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to determine home directory: %w", err)
		}
		base = filepath.Join(home, ".local", "state")
	}
	return filepath.Join(base, "p", "session-log.jsonl"), nil
}
