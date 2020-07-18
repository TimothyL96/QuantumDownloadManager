package userSettings

import (
	"errors"
	"strconv"
	"strings"

	"github.com/ttimt/QuantumDownloadManager/downloadManager"
)

type UserSettings struct {
	nrOfConcurrentDownload int
}

func (us *UserSettings) GetNrOfConcurrentDownload() int {
	return us.nrOfConcurrentDownload
}

func (us *UserSettings) SetNrOfConcurrentDownload(nrOfConcurrentDownload int) error {
	var err error

	if nrOfConcurrentDownload > downloadManager.MaxNrOfConcurrentDownload {
		// Defaulting to max
		err = errors.New("defaulting to the maximum allowed concurrent download (" +
			strconv.Itoa(downloadManager.MaxNrOfConcurrentDownload) +
			") as given number of concurrent download exceeded maximum allowed")

		us.nrOfConcurrentDownload = downloadManager.MaxNrOfConcurrentDownload
	} else if nrOfConcurrentDownload < 1 {
		// Defaulting to 1
		err = errors.New("defaulting to 1 concurrent download as given number of concurrent download is set below 1")

		us.nrOfConcurrentDownload = 1
	} else {
		us.nrOfConcurrentDownload = nrOfConcurrentDownload
	}

	return err
}

func (us *UserSettings) String() string {
	sb := strings.Builder{}

	sb.WriteString("Number of concurrent download: ")
	sb.WriteString(strconv.Itoa(us.GetNrOfConcurrentDownload()))
	sb.WriteString("\n")

	return sb.String()
}
