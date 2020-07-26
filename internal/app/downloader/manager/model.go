package manager

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/ttimt/QuantumDownloadManager/internal/app/downloader/util/file"
)

const (
	// MaxNrOfConcurrentConnectionAllowed is the maximum number of connection allowed.
	MaxNrOfConcurrentConnectionAllowed = 64

	// TempFileFileExtension is the file extension for temporary download file.
	TempFileFileExtension = "qdm"
)

// Download is a session of a download.
type Download struct {
	// Download details
	downloadURL                 *url.URL
	maxNrOfConcurrentConnection int
	saveDirectory               string
	saveFullPath                string
	saveFileName                string
	defaultFileName             string

	// Flags
	isPausedAllowed               FlagState
	isConcurrentConnectionAllowed FlagState

	// Download status
	isDownloadInitialized bool
	isDownloadStarted     bool // Has the download been started once before
	isDownloadRunning     bool
	isDownloadCompleted   bool
	isDownloadAborted     bool

	// Temporary files variables
	tempFileNameAppender int
	tempFileList         []string

	// Response
	response *http.Response
	fileSize int64

	// Context
	ctx       context.Context
	ctxCancel func()
}

// DownloadURL returns current download URL.
func (d *Download) DownloadURL() string {
	return d.downloadURL.String()
}

// SetDownloadURL sets the download URL with the input
// and returns a non nil error if failed to set the download URL.
func (d *Download) SetDownloadURL(downloadUrl string) error {
	if len(downloadUrl) == 0 {
		return errors.New("download URL is empty")
	}

	parsedDownloadUrl, err := url.Parse(downloadUrl)
	if err != nil {
		return err
	}

	d.downloadURL = parsedDownloadUrl

	return nil
}

// MaxNrOfConcurrentConnection returns the number of concurrent connection set.
func (d *Download) MaxNrOfConcurrentConnection() int {
	return d.maxNrOfConcurrentConnection
}

// SetMaxNrOfConcurrentConnection set the maximum number of concurrent connection for the current download
// and return non nil error if failed to set.
func (d *Download) SetMaxNrOfConcurrentConnection(nrOfConcurrentConnection int) error {
	if nrOfConcurrentConnection > MaxNrOfConcurrentConnectionAllowed {
		return errors.New("number of concurrent connection given exceeded maximum allowed (" +
			strconv.Itoa(MaxNrOfConcurrentConnectionAllowed) +
			")")
	} else if nrOfConcurrentConnection < 1 {
		return errors.New("number of concurrent connection given is less than 1")
	}

	d.maxNrOfConcurrentConnection = nrOfConcurrentConnection

	return nil
}

// SaveDirectory returns the directory to save the download file.
func (d *Download) SaveDirectory() string {
	return d.saveDirectory
}

// SetSaveDirectory set the download save directory to the input
// and returns a non nil error if failed to set.
func (d *Download) SetSaveDirectory(directory string) error {
	if len(directory) == 0 {
		return errors.New("save directory cannot be empty")
	}

	// Check if cleaned directory path exists
	isFileExists := file.IsCleanedFileExist(directory)
	if !isFileExists {
		// Currently does not automatically create the missing directory
		return errors.New("the given save directory does not exists")
	}

	d.saveDirectory = directory

	return nil
}

// SaveFileName returns the file name to be used for the current download.
func (d *Download) SaveFileName() string {
	return d.saveFileName
}

// setFileName is a helper method for setting file name
func (d *Download) setFileName(fileName string) error {
	if len(fileName) == 0 {
		return errors.New("file name cannot be empty")
	}

	// Clean file name
	fileName = filepath.Clean(fileName)

	// File name cannot have slash as it will change the save directory path
	if strings.ContainsRune(fileName, '/') {
		return errors.New("forward slash '/' is not allowed in file name")
	}

	// When combining and writing to file:
	// If a file with the same name as current download save file name exists
	// append a unique number behind the file name
	//
	// Make sure file name and path does not clash with other incomplete downloads in this application!

	d.saveFileName = fileName

	return nil
}

// SetSaveFileName set the file name to be used for the current download
// and returns a non nil error if failed to set.
// If not set, save file name defaults to server provided name if available.
func (d *Download) SetSaveFileName(fileName string) error {
	err := d.setFileName(fileName)
	if err != nil {
		return err
	}

	d.saveFileName = fileName

	return nil
}

// DefaultFileName returns the default file name provided by the request server
func (d *Download) DefaultFileName() string {
	return d.defaultFileName
}

// setDefaultFileName set the default file name provided by the server
func (d *Download) setDefaultFileName(fileName string) error {
	err := d.setFileName(fileName)
	if err != nil {
		return err
	}

	d.defaultFileName = fileName

	return nil
}

// SaveFullPath returns the full path including directory and file name.
func (d *Download) SaveFullPath() string {
	return d.saveFullPath
}

// setSaveFullPath receives no parameter
// and will attempt to combine current download save directory and save file name
// and will return a non nil error if error occurred.
func (d *Download) setSaveFullPath() error {
	saveFileName := ""

	if d.SaveFileName() != "" {
		saveFileName = d.SaveFileName()
	} else if d.DefaultFileName() != "" {
		saveFileName = d.DefaultFileName()
	} else {
		return errors.New("save file name not set")
	}

	// Join the directory and file name together
	d.saveFullPath = filepath.Join(d.saveDirectory, saveFileName)

	if len(d.saveFullPath) <= 0 {
		return errors.New("error in combining save directory and save file path")
	}

	return nil
}

// FileSize returns the file size for the download.
// If download details is not retrieved, file size should be 0.
func (d *Download) FileSize() int64 {
	return d.fileSize
}

func (d *Download) setFileSize(fileSize int64) error {
	d.fileSize = fileSize

	return nil
}

// IsPausedAllowed returns a state indicating if pausing the download is supported.
func (d *Download) IsPausedAllowed() FlagState {
	return d.isPausedAllowed
}

func (d *Download) setIsPausedAllowed(isPausedAllowed FlagState) error {
	d.isPausedAllowed = isPausedAllowed

	return nil
}

// IsConcurrentConnectionAllowed returns a state indicating if concurrent connection is supported.
func (d *Download) IsConcurrentConnectionAllowed() FlagState {
	return d.isConcurrentConnectionAllowed
}

func (d *Download) setIsConcurrentConnectionAllowed(isConcurrentConnectionAllowed FlagState) error {
	d.isConcurrentConnectionAllowed = isConcurrentConnectionAllowed

	return nil
}

// IsDownloadInitialized returns a boolean indicating if the download has been initialized before.
func (d *Download) IsDownloadInitialized() bool {
	return d.isDownloadInitialized
}

func (d *Download) setIsDownloadInitialized(isDownloadInitialized bool) error {
	d.isDownloadInitialized = isDownloadInitialized

	return nil
}

// IsDownloadStarted returns a boolean indicating if the download has been started before.
func (d *Download) IsDownloadStarted() bool {
	return d.isDownloadStarted
}

func (d *Download) setIsDownloadStarted(isDownloadStarted bool) error {
	d.isDownloadStarted = isDownloadStarted

	return nil
}

// IsDownloadRunning returns a boolean indicating whether the download is currently running.
func (d *Download) IsDownloadRunning() bool {
	return d.isDownloadRunning
}

func (d *Download) setIsDownloadRunning(isDownloadRunning bool) error {
	d.isDownloadRunning = isDownloadRunning

	return nil
}

// IsDownloadCompleted returns a boolean indicating whether the download has been completed.
func (d *Download) IsDownloadCompleted() bool {
	return d.isDownloadCompleted
}

func (d *Download) setIsDownloadCompleted(isDownloadCompleted bool) error {
	d.isDownloadCompleted = isDownloadCompleted

	return nil
}

func (d *Download) setTempFileNameAppender(fileAppender int) error {
	d.tempFileNameAppender = fileAppender

	return nil
}

func (d *Download) setTempFileList(tempFileList []string) error {
	d.tempFileList = tempFileList

	return nil
}

func (d *Download) setCtx(ctx context.Context) error {
	d.ctx = ctx

	return nil
}

func (d *Download) setCtxCancel(ctxCancel func()) error {
	d.ctxCancel = ctxCancel

	return nil
}

func (d *Download) setResponse(response *http.Response) error {
	d.response = response

	return nil
}

func (d *Download) String() string {
	sb := strings.Builder{}

	sb.WriteString("RetrieveDownloadDetails URL: ")
	sb.WriteString(d.DownloadURL())
	sb.WriteString("\n")

	sb.WriteString("Number of concurrent download: ")
	sb.WriteString(strconv.Itoa(d.MaxNrOfConcurrentConnection()))
	sb.WriteString("\n")

	sb.WriteString("Save directory: ")
	sb.WriteString(d.SaveDirectory())
	sb.WriteString("\n")

	sb.WriteString("Save file name: ")
	sb.WriteString(d.SaveFileName())
	sb.WriteString("\n")

	sb.WriteString("Save full path:")
	sb.WriteString(d.SaveFullPath())
	sb.WriteString("\n")

	sb.WriteString("Is paused allowed: ")
	sb.WriteString(d.IsPausedAllowed().String())
	sb.WriteString("\n")

	sb.WriteString("Is concurrent connection allowed: ")
	sb.WriteString(d.IsConcurrentConnectionAllowed().String())
	sb.WriteString("\n")

	sb.WriteString("Is download initialized:")
	sb.WriteString(strconv.FormatBool(d.isDownloadInitialized))
	sb.WriteString("\n")

	sb.WriteString("Is download started:")
	sb.WriteString(strconv.FormatBool(d.IsDownloadStarted()))
	sb.WriteString("\n")

	sb.WriteString("Is download running:")
	sb.WriteString(strconv.FormatBool(d.IsDownloadRunning()))
	sb.WriteString("\n")

	sb.WriteString("Is download completed:")
	sb.WriteString(strconv.FormatBool(d.IsDownloadCompleted()))
	sb.WriteString("\n")

	return sb.String()
}
