// Package file contains utility functions that deal with files.
package file

import (
	"os"
	"path/filepath"
	"strings"
)

// IsFileExist check whether the given file or directory exists.
func IsFileExist(path string) bool {
	_, err := os.Stat(path)

	return !os.IsNotExist(err)
}

// CleanPath returns the given path trimmed and cleaned.
// Different OSes directory slashes will be handled.
func CleanPath(path string) string {
	// Trim any spaces at the edge
	path = strings.TrimSpace(path)

	// Clean the path
	path = filepath.Clean(path)

	return path
}

// IsCleanedFileExist cleans the path before checking whether the given file or directory exists.
func IsCleanedFileExist(path string) bool {
	path = CleanPath(path)

	return IsFileExist(path)
}
