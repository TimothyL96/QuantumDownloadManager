// Package file contains utility functions that deal with files.
package file

import (
	"os"
)

// IsFileExist check whether the given file or directory exists.
func IsFileExist(path string) bool {
	_, err := os.Stat(path)

	return !os.IsNotExist(err)
}
