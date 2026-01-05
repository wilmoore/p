package tmux

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDiscoverDirectories(t *testing.T) {
	// Create temp directory structure
	tmpDir := t.TempDir()

	// Create some project directories
	projects := []string{"project1", "project2", "project3"}
	for _, p := range projects {
		if err := os.Mkdir(filepath.Join(tmpDir, p), 0755); err != nil {
			t.Fatalf("failed to create test directory: %v", err)
		}
	}

	// Create a hidden directory that should be skipped
	if err := os.Mkdir(filepath.Join(tmpDir, ".hidden"), 0755); err != nil {
		t.Fatalf("failed to create hidden directory: %v", err)
	}

	// Create a file that should be skipped
	if err := os.WriteFile(filepath.Join(tmpDir, "file.txt"), []byte("test"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	dirs, err := DiscoverDirectories(tmpDir)
	if err != nil {
		t.Fatalf("DiscoverDirectories failed: %v", err)
	}

	if len(dirs) != 3 {
		t.Errorf("expected 3 directories, got %d", len(dirs))
	}

	// Verify all expected projects are found
	found := make(map[string]bool)
	for _, d := range dirs {
		found[d.Name] = true
	}
	for _, p := range projects {
		if !found[p] {
			t.Errorf("expected to find directory %q", p)
		}
	}
}

func TestDiscoverDirectoriesEmpty(t *testing.T) {
	dirs, err := DiscoverDirectories("")
	if err != nil {
		t.Fatalf("DiscoverDirectories failed: %v", err)
	}
	if len(dirs) != 0 {
		t.Errorf("expected 0 directories for empty CDPATH, got %d", len(dirs))
	}
}

func TestDiscoverDirectoriesMultiplePaths(t *testing.T) {
	// Create two temp directories
	tmpDir1 := t.TempDir()
	tmpDir2 := t.TempDir()

	// Create projects in each
	if err := os.Mkdir(filepath.Join(tmpDir1, "proj1"), 0755); err != nil {
		t.Fatalf("failed to create test directory: %v", err)
	}
	if err := os.Mkdir(filepath.Join(tmpDir2, "proj2"), 0755); err != nil {
		t.Fatalf("failed to create test directory: %v", err)
	}

	cdpath := tmpDir1 + ":" + tmpDir2
	dirs, err := DiscoverDirectories(cdpath)
	if err != nil {
		t.Fatalf("DiscoverDirectories failed: %v", err)
	}

	if len(dirs) != 2 {
		t.Errorf("expected 2 directories, got %d", len(dirs))
	}
}
