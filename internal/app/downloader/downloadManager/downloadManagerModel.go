package downloadManager

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/ttimt/QuantumDownloadManager/internal/app/downloader/utils/file"
)

const (
	MaxNrOfConcurrentDownload = 64
)

const (
	unknown FeatureStatus = iota
	allowed
	notAllowed
)

type FeatureStatus int

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
	isPausedAllowed        FeatureStatus
	isConcurrentAllowed    FeatureStatus
	tempAppender           int
	tempFileLists          []string
	isDownloadStarted      bool
	isDownloadRunning      bool
	isDownloadComplete     bool
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

	// Clean the path
	saveDirectory = filepath.Clean(saveDirectory)

	// Check if directory exists
	isFileExists := file.IsFileExist(saveDirectory)
	if !isFileExists {
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

	// Clean file name
	saveFileName = filepath.Clean(saveFileName)

	// Forward slash '/' not allowed in name
	if strings.ContainsRune(saveFileName, '/') {
		return errors.New("forward slash '/' is not allowed in file name")
	}

	// Need to also check is the file name being used by 1 of the download instances
	isFileExists := file.IsFileExist(saveFileName)
	if isFileExists {
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

func (d *DownloadManager) getCtx() context.Context {
	return d.ctx
}

func (d *DownloadManager) setCtx(ctx context.Context) error {
	d.ctx = ctx

	return nil
}

func (d *DownloadManager) getCtxCancel() func() {
	return d.ctxCancel
}

func (d *DownloadManager) setCtxCancel(ctxCancel func()) error {
	d.ctxCancel = ctxCancel

	return nil
}

func (d *DownloadManager) getResponse() *http.Response {
	return d.response
}

func (d *DownloadManager) setResponse(response *http.Response) error {
	d.response = response

	return nil
}

func (d *DownloadManager) GetFileSize() int64 {
	return d.fileSize
}

func (d *DownloadManager) setFileSize(fileSize int64) error {
	d.fileSize = fileSize

	return nil
}

func (d *DownloadManager) GetIsPausedAllowed() FeatureStatus {
	return d.isPausedAllowed
}

func (d *DownloadManager) setIsPausedAllowed(isPausedAllowed FeatureStatus) error {
	d.isPausedAllowed = isPausedAllowed

	return nil
}

func (d *DownloadManager) GetIsConcurrentAllowed() FeatureStatus {
	return d.isConcurrentAllowed
}

func (d *DownloadManager) setIsConcurrentAllowed(isConcurrentAllowed FeatureStatus) error {
	d.isConcurrentAllowed = isConcurrentAllowed

	return nil
}

func (d *DownloadManager) getTempAppender() int {
	return d.tempAppender
}

func (d *DownloadManager) incrementTempAppender() error {
	d.tempAppender += 1

	return nil
}

func (d *DownloadManager) getTempFileList() []string {
	return d.tempFileLists
}

func (d *DownloadManager) setTempFileList(tempFileName string) error {
	d.tempFileLists = append(d.tempFileLists, tempFileName)

	return nil
}

func (d *DownloadManager) GetIsDownloadStarted() bool {
	return d.isDownloadStarted
}

func (d *DownloadManager) SetIsDownloadStarted(isDownloadStarted bool) error {
	d.isDownloadStarted = isDownloadStarted

	return nil
}

func (d *DownloadManager) GetIsDownloadRunning() bool {
	return d.isDownloadRunning
}

func (d *DownloadManager) SetIsDownloadRunning(isDownloadRunning bool) error {
	d.isDownloadRunning = isDownloadRunning

	return nil
}

func (d *DownloadManager) GetIsDownloadComplete() bool {
	return d.isDownloadComplete
}

func (d *DownloadManager) SetIsDownloadComplete(isDownloadComplete bool) error {
	d.isDownloadComplete = isDownloadComplete

	return nil
}

func (d *DownloadManager) String() string {
	sb := strings.Builder{}

	sb.WriteString("RetrieveDownloadDetails URL: ")
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
