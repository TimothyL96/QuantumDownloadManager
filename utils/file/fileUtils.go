package file

import (
	"os"
	"path/filepath"
)

// Check the given file or directory exists
func IsFileExist(path string) (bool, error) {
	path = filepath.Clean(path)

	_, err := os.Stat(path)

	// Return error if there are other errors
	if err != nil && err != os.ErrNotExist {
		return false, err
	}

	return !os.IsNotExist(err), nil
}
