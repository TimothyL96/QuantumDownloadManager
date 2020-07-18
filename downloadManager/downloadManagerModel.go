package downloadManager

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	MaxNrOfConcurrentDownload = 64
)

type DownloadManager struct {
	downloadUrl            *url.URL
	nrOfConcurrentDownload int
	saveDirectory          string
	saveFileName           string
	saveFullPath           string
	ctx                    context.Context
	ctxCancel              func()
	response               *http.Response
	fileSize               int64
	isPausedAllowed        bool
	isConcurrentAllowed    bool
}

func (d *DownloadManager) GetDownloadUrl() string {
	return d.downloadUrl.String()
}

func (d *DownloadManager) SetDownloadUrl(downloadUrl string) error {
	if len(downloadUrl) == 0 {
		return errors.New("download URL is empty")
	}

	parsedDownloadUrl, err := url.Parse(downloadUrl)
	if err != nil {
		return err
	}

	d.downloadUrl = parsedDownloadUrl

	return nil
}

func (d *DownloadManager) GetNrOfConcurrentDownload() int {
	return d.nrOfConcurrentDownload
}

func (d *DownloadManager) SetNrOfConcurrentDownload(nrOfConcurrentDownload int) error {
	if nrOfConcurrentDownload > MaxNrOfConcurrentDownload {
		return errors.New("number of concurrent download given exceeded maximum allowed (" +
			strconv.Itoa(MaxNrOfConcurrentDownload) +
			")")
	} else if nrOfConcurrentDownload < 1 {
		return errors.New("number of concurrent download given is less than 1")
	}

	d.nrOfConcurrentDownload = nrOfConcurrentDownload

	return nil
}

func (d *DownloadManager) GetSaveDirectory() string {
	return d.saveDirectory
}

func (d *DownloadManager) SetSaveDirectory(saveDirectory string) error {
	if len(saveDirectory) == 0 {
		return errors.New("save directory cannot be empty")
	}

	saveDirectory = strings.TrimSpace(saveDirectory)

	// Convert the slashes
	saveDirectory = filepath.FromSlash(saveDirectory)

	// Check if directory exists
	_, err := os.Stat(saveDirectory)
	if os.IsNotExist(err) {
		// Currently does not automatically create the missing directory
		return errors.New("the given save directory does not exists")
	}

	d.saveDirectory = saveDirectory

	return nil
}

func (d *DownloadManager) GetSaveFileName() string {
	return d.saveFileName
}

func (d *DownloadManager) SetSaveFileName(saveFileName string) error {
	if len(saveFileName) == 0 {
		return errors.New("file name cannot be empty")
	}

	// Convert slashes to '/'
	saveFileName = filepath.ToSlash(saveFileName)

	// Forward slash '/' not allowed in name
	if strings.ContainsRune(saveFileName, '/') {
		return errors.New("forward slash '/' is not allowed in file name")
	}

	_, err := os.Stat(saveFileName)
	if !os.IsNotExist(err) {
		// File already exists
		return errors.New("the given file name to save already exists")
	}

	d.saveFileName = saveFileName

	return nil
}

func (d *DownloadManager) GetSaveFullPath() string {
	return d.saveFullPath
}

func (d *DownloadManager) setSaveFullPath() error {
	// Join the directory and file name together
	d.saveFullPath = filepath.Join(d.saveDirectory, d.saveFileName)

	if len(d.saveFullPath) <= 0 {
		return errors.New("error in combining save directory and file path")
	}

	return nil
}

func (d *DownloadManager) String() string {
	sb := strings.Builder{}

	sb.WriteString("Download URL: ")
	sb.WriteString(d.GetDownloadUrl())
	sb.WriteString("\n")

	sb.WriteString("Number of concurrent download: ")
	sb.WriteString(strconv.Itoa(d.GetNrOfConcurrentDownload()))
	sb.WriteString("\n")

	sb.WriteString("Save directory: ")
	sb.WriteString(d.GetSaveDirectory())
	sb.WriteString("\n")

	sb.WriteString("Save file name: ")
	sb.WriteString(d.GetSaveFileName())
	sb.WriteString("\n")

	sb.WriteString("Save full path:")
	sb.WriteString(d.GetSaveFullPath())
	sb.WriteString("\n")

	return sb.String()
}
