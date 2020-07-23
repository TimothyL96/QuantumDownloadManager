// Package setting manages the settings of user
package setting

import (
	"strconv"
	"strings"

	"github.com/getlantern/errors"

	"github.com/ttimt/QuantumDownloadManager/internal/app/downloader/downloadManager"
)

// Setting stores the settings of a user.
type Setting struct {
	nrOfConcurrentConnection int
}

// NrOfConcurrentConnection returns the number of concurrent connection set in user setting.
func (s *Setting) NrOfConcurrentConnection() int {
	return s.nrOfConcurrentConnection
}

// SetNrOfConcurrentConnection updates the user setting with the number of concurrent download.
//
// If the number is over maximum limit, the maximum concurrent connection will be set and error will not be nil.
// Similarly, if given number is less than 1, it will be defaulted to 1 and the return error will not be nil.
func (s *Setting) SetNrOfConcurrentConnection(nrOfConcurrentConnection int) error {
	var err error

	s.nrOfConcurrentConnection = nrOfConcurrentConnection

	if nrOfConcurrentConnection > downloadManager.MaxNrOfConcurrentConnection {
		// The given number exceeds the maximum allowed connection
		// Defaulting to maximum concurrent connection
		err = errors.Wrap(errors.New("defaulting to the maximum allowed concurrent connection (" +
			strconv.Itoa(downloadManager.MaxNrOfConcurrentConnection) +
			") as the given number exceeded maximum allowed"))

		s.nrOfConcurrentConnection = downloadManager.MaxNrOfConcurrentConnection
	} else if nrOfConcurrentConnection < 1 {
		// The given number is below 1
		// Defaulting to 1
		err = errors.Wrap(errors.New("defaulting to 1 concurrent connection as the given number is below 1"))

		s.nrOfConcurrentConnection = 1
	}

	return err
}

func (s *Setting) String() string {
	sb := strings.Builder{}

	sb.WriteString("Number of concurrent connection: ")
	sb.WriteString(strconv.Itoa(s.NrOfConcurrentConnection()))
	sb.WriteString("\n")

	return sb.String()
}
