package tmux

import (
	"os"
	"path/filepath"
	"strings"
)

// Directory represents a project directory discovered via CDPATH.
type Directory struct {
	Name       string
	Path       string
	HasSubdirs bool // True if directory contains subdirectories
}

// DiscoverDirectories finds all directories in the CDPATH locations.
// Only immediate children of CDPATH directories are returned (no recursion).
func DiscoverDirectories(cdpath string) ([]Directory, error) {
	if cdpath == "" {
		return nil, nil
	}

	paths := strings.Split(cdpath, ":")
	var dirs []Directory

	for _, basePath := range paths {
		basePath = strings.TrimSpace(basePath)
		if basePath == "" {
			continue
		}

		// Expand ~ to home directory
		if strings.HasPrefix(basePath, "~") {
			home, err := os.UserHomeDir()
			if err != nil {
				continue
			}
			basePath = filepath.Join(home, basePath[1:])
		}

		entries, err := os.ReadDir(basePath)
		if err != nil {
			// Skip directories that don't exist or can't be read
			continue
		}

		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}
			// Skip hidden directories
			if strings.HasPrefix(entry.Name(), ".") {
				continue
			}

			fullPath := filepath.Join(basePath, entry.Name())
			dirs = append(dirs, Directory{
				Name:       entry.Name(),
				Path:       fullPath,
				HasSubdirs: hasSubdirectories(fullPath),
			})
		}
	}

	return dirs, nil
}

// hasSubdirectories checks if a directory contains any non-hidden subdirectories.
func hasSubdirectories(path string) bool {
	entries, err := os.ReadDir(path)
	if err != nil {
		return false
	}

	for _, entry := range entries {
		if entry.IsDir() && !strings.HasPrefix(entry.Name(), ".") {
			return true
		}
	}
	return false
}

// GetSubdirectories returns the immediate subdirectories of a given directory.
// Used for drill-down navigation.
func GetSubdirectories(dir *Directory) ([]Directory, error) {
	entries, err := os.ReadDir(dir.Path)
	if err != nil {
		return nil, err
	}

	var subdirs []Directory
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		// Skip hidden directories
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		fullPath := filepath.Join(dir.Path, entry.Name())
		subdirs = append(subdirs, Directory{
			Name:       entry.Name(),
			Path:       fullPath,
			HasSubdirs: hasSubdirectories(fullPath),
		})
	}

	return subdirs, nil
}
