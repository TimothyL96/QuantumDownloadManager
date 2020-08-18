package manager

import (
	"errors"
	"fmt"
	"io"
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
		strconv.Itoa(d.tempFileNameSuffix) +
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
	_ = d.setTempFileNameAppender(d.tempFileNameSuffix + 1)
}

func (d *Download) appendToTempFileList(tempFilePath string) {
	_ = d.setTempFileList(append(d.tempFileList, tempFilePath))
}

// combineFiles combines all temporary files together to form the final download file.
func (d *Download) combineFiles() error {
	// Combine files
	fmt.Println("Writing to file:")

	// Must have at least 1 temporary file
	if len(d.tempFileList) < 1 {
		return errors.New("must have at least 1 temporary file")
	}

	// Open the first file to append other data onto it
	firstTempFile, err := os.OpenFile(d.tempFileList[0], os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		return err
	}

	// Loop through other temporary files and put the data into the first file
	for _, v := range d.tempFileList[1:] {
		// Open current file
		f, err := os.OpenFile(v, os.O_RDONLY, os.ModePerm)
		if err != nil {
			return err
		}

		_, err = io.Copy(firstTempFile, f)
		if err != nil {
			return err
		}

		// Close the file at the end
		f.Close()
	}

	// Close the file at the end
	firstTempFile.Close()

	// Rename the file to final download file
	if err = os.Rename(firstTempFile.Name(), d.SaveFullPath()); err != nil {
		d.Abort()
		return err
	}

	// Delete all temporary files
	for _, v := range d.tempFileList[1:] {
		if err := os.Remove(v); err != nil {
			return err
		}
	}

	fmt.Println("Combine temporary files done")

	return nil
}
