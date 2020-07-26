package manager

import (
	"os"
	"strconv"

	"github.com/ttimt/QuantumDownloadManager/internal/app/downloader/util/file"
)

// createTemporaryFile creates a temporary file with a unique name appended by a number.
func (d *Download) createTemporaryFile() (*os.File, error) {
	// Increment temporary tempFile number
	d.incrementTempFileAppender()

	// Create a new temporary tempFile path
	tempFilePath := d.SaveFullPath() +
		".temp" +
		strconv.Itoa(d.tempFileNameAppender) +
		"." +
		TempFileFileExtension

	// Check if temporary tempFile exists
	isFileExists := file.IsFileExist(tempFilePath)
	if isFileExists {
		// If tempFile name already exists
		// Increment the temporary tempFile appender
		return d.createTemporaryFile()

		// Check infinite loop?
	}

	// Check if there's enough storage to store the tempFile
	// Use os.truncate to increase size of the new temp tempFile without writing to tempFile

	// Create the tempFile
	tempFile, err := os.OpenFile(tempFilePath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}

	// Add this to the temporary tempFile list
	d.appendToTempFileList(tempFilePath)

	return tempFile, nil
}

// createPlaceHolderFile creates a placeholder file
// with the name same as the download save file name.
func (d *Download) createPlaceHolderFile() {
	// Create a place holder file with the same name as the download save file and path
	// Check file exist before creating the place holder file
	placeholderFile, _ := os.Create(d.SaveFullPath())
	_ = placeholderFile.Close()
}

func (d *Download) incrementTempFileAppender() {
	_ = d.setTempFileNameAppender(d.tempFileNameAppender + 1)
}

func (d *Download) appendToTempFileList(tempFilePath string) {
	_ = d.setTempFileList(append(d.tempFileList, tempFilePath))
}
