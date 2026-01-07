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

func TestHasSubdirectories(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a directory with subdirectories
	parent := filepath.Join(tmpDir, "parent")
	if err := os.Mkdir(parent, 0755); err != nil {
		t.Fatalf("failed to create parent directory: %v", err)
	}
	if err := os.Mkdir(filepath.Join(parent, "child"), 0755); err != nil {
		t.Fatalf("failed to create child directory: %v", err)
	}

	// Create a directory without subdirectories
	empty := filepath.Join(tmpDir, "empty")
	if err := os.Mkdir(empty, 0755); err != nil {
		t.Fatalf("failed to create empty directory: %v", err)
	}

	// Test discovery includes HasSubdirs field
	dirs, err := DiscoverDirectories(tmpDir)
	if err != nil {
		t.Fatalf("DiscoverDirectories failed: %v", err)
	}

	for _, d := range dirs {
		switch d.Name {
		case "parent":
			if !d.HasSubdirs {
				t.Errorf("expected parent to have HasSubdirs=true")
			}
		case "empty":
			if d.HasSubdirs {
				t.Errorf("expected empty to have HasSubdirs=false")
			}
		}
	}
}

func TestHasSubdirectoriesIgnoresHidden(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a directory with only hidden subdirectories
	parent := filepath.Join(tmpDir, "parent")
	if err := os.Mkdir(parent, 0755); err != nil {
		t.Fatalf("failed to create parent directory: %v", err)
	}
	if err := os.Mkdir(filepath.Join(parent, ".hidden"), 0755); err != nil {
		t.Fatalf("failed to create hidden directory: %v", err)
	}

	dirs, err := DiscoverDirectories(tmpDir)
	if err != nil {
		t.Fatalf("DiscoverDirectories failed: %v", err)
	}

	if len(dirs) != 1 {
		t.Fatalf("expected 1 directory, got %d", len(dirs))
	}

	// Should not have subdirs since only hidden exists
	if dirs[0].HasSubdirs {
		t.Errorf("expected HasSubdirs=false when only hidden subdirectories exist")
	}
}

func TestGetSubdirectories(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a parent directory with subdirectories
	parent := &Directory{
		Name: "parent",
		Path: tmpDir,
	}

	// Create subdirectories
	subdirs := []string{"sub1", "sub2", "sub3"}
	for _, s := range subdirs {
		if err := os.Mkdir(filepath.Join(tmpDir, s), 0755); err != nil {
			t.Fatalf("failed to create subdirectory: %v", err)
		}
	}

	// Create a nested subdirectory in sub1
	if err := os.Mkdir(filepath.Join(tmpDir, "sub1", "nested"), 0755); err != nil {
		t.Fatalf("failed to create nested directory: %v", err)
	}

	// Create hidden directory that should be skipped
	if err := os.Mkdir(filepath.Join(tmpDir, ".hidden"), 0755); err != nil {
		t.Fatalf("failed to create hidden directory: %v", err)
	}

	// Create a file that should be skipped
	if err := os.WriteFile(filepath.Join(tmpDir, "file.txt"), []byte("test"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	result, err := GetSubdirectories(parent)
	if err != nil {
		t.Fatalf("GetSubdirectories failed: %v", err)
	}

	if len(result) != 3 {
		t.Errorf("expected 3 subdirectories, got %d", len(result))
	}

	// Verify all expected subdirs are found and HasSubdirs is set correctly
	found := make(map[string]bool)
	for _, d := range result {
		found[d.Name] = true
		if d.Name == "sub1" && !d.HasSubdirs {
			t.Errorf("expected sub1 to have HasSubdirs=true")
		}
		if (d.Name == "sub2" || d.Name == "sub3") && d.HasSubdirs {
			t.Errorf("expected %s to have HasSubdirs=false", d.Name)
		}
	}

	for _, s := range subdirs {
		if !found[s] {
			t.Errorf("expected to find subdirectory %q", s)
		}
	}
}
