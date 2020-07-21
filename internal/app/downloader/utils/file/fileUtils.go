package file

import (
	"os"
)

// Check the given file or directory exists
func IsFileExist(path string) bool {
	_, err := os.Stat(path)

	return !os.IsNotExist(err)
}
