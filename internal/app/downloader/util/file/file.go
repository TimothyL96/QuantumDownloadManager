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

// FixPath returns the given path trimmed and cleaned.
// Different OSes slashes will be handled.
func FixPath(path string) string {
	// Trim any spaces at the edge
	path = strings.TrimSpace(path)

	// Clean the path
	path = filepath.Clean(path)

	return path
}
