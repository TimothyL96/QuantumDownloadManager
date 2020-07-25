package manager

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	fileUtils "github.com/ttimt/QuantumDownloadManager/internal/app/downloader/util/file"
)

// InitializeDownload initialize the new download by sending a request to the download URL
// and updating the download fields value by processing the received header.
func (d *Download) InitializeDownload() error {
	if d.isDownloadRunning {
		return errors.New("download is currently running")
	}

	// Setup new context for stopping download
	ctx, ctxCancel := context.WithCancel(context.Background())
	_ = d.setCtx(ctx)
	_ = d.setCtxCancel(ctxCancel)

	// Setup request with the newly created instance's context
	req, err := http.NewRequestWithContext(d.ctx,
		http.MethodGet,
		d.downloadURL.String(),
		nil)
	if err != nil {
		return err
	}

	// Make the request to get the response header
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	_ = d.setResponse(response)

	// Process the received request header
	err = d.processRequestHeader()
	if err != nil {
		return err
	}

	// Set download as initialized
	return d.setIsDownloadInitialized(true)
}

func (d *Download) processRequestHeader() error {
	// Partial download reference:
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Accept-Ranges
	//
	// Content-Length header - If value is -1, we do not know the file size
	// and unable to split it to download concurrently or resume download
	//
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Length
	//
	// Known content length
	if d.response.ContentLength > 0 {
		// If header "Accept-Ranges" exists and value its not none
		// Then partial request (concurrent download) / pause is supported
		//
		// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Accept-Ranges
		acceptRanges, ok := d.response.Header["Accept-Ranges"]

		// Accept-range exists
		if ok {
			// Concurrent download and pause is likely not allowed
			if acceptRanges[0] == "none" {
				_ = d.setIsConcurrentConnectionAllowed(notAllowed)
				_ = d.setIsPausedAllowed(notAllowed)
			} else {
				// Allowed
				_ = d.setIsConcurrentConnectionAllowed(allowed)
				_ = d.setIsPausedAllowed(allowed)
			}
		} else {
			// Unknown
			_ = d.setIsConcurrentConnectionAllowed(unknown)
			_ = d.setIsPausedAllowed(unknown)
		}
	} else {
		// Unknown content length
		_ = d.setIsConcurrentConnectionAllowed(notAllowed)
		_ = d.setIsPausedAllowed(notAllowed)
	}

	// Get suggested default file name from header - Content-Disposition
	// d.setDefaultFileName(...)

	return nil
}

// Download starts the download first time.
func (d *Download) Download() error {
	// Test
	_ = d.setIsConcurrentConnectionAllowed(notAllowed)

	// Block if download has started before
	if d.IsDownloadStarted() {
		return errors.New("download has already started before. Did you mean resume download ")
	}

	// Flag the download has started
	_ = d.setIsDownloadStarted(true)

	// Start the download
	if d.IsConcurrentConnectionAllowed() == allowed {
		return d.StartConcurrentDownload()
	}

	// Non concurrent single connection download
	return d.StartAtomicDownload()
}

// StartAtomicDownload download the file without any concurrent connection.
func (d *Download) StartAtomicDownload() error {
	// Check if download has been started before, and resume the last pause state
	// To be done when implementing pause feature

	// Create downloader single temporary file
	tempFile, err := d.CreateTemporaryFile()
	if err != nil {
		return err
	}

	// Set the download as running
	_ = d.setIsDownloadRunning(true)

	// Write the data to disk
	written, err := io.Copy(tempFile, d.response.Body)
	if err != nil {
		return err
	}

	log.Printf("DEBUG: Written %d out of %d bytes", written, d.response.ContentLength)

	return nil
}

// StartConcurrentDownload download the file part by part concurrently with the number of concurrent connection set.
func (d *Download) StartConcurrentDownload() error {
	// data := make([]byte, res.ContentLength)
	//
	// fmt.Println("Now start streaming")
	// resBodySize, err := io.ReadFull(res.Body, data)
	// fmt.Println("Streaming complete")

	// if err != nil {
	// 	fmt.Println("Read body size", resBodySize)
	// 	fmt.Println("Error:", err)
	// }
	//
	// downloadedFile, err := os.OpenFile("D:\\Timothy\\Desktop\\download2.mp4", os.O_CREATE|os.O_RDWR, os.ModePerm)
	// if err != nil {
	// 	panic(err)
	// }
	//
	// writeSize, err := downloadedFile.Write(data)
	//
	// if err != nil {
	// 	fmt.Println("Write size", writeSize)
	// 	panic(err)
	// }
	//
	// cancel()

	return nil
}

// StreamData xxx.
func (d *Download) StreamData() error {
	return nil
}

// Pause will pause the current download if it is currently running.
func (d *Download) Pause() error {
	return nil
}

// Resume will continue the current download if it is paused.
func (d *Download) Resume() error {
	return nil
}

// Abort will cancel and clear the current download.
func (d *Download) Abort() error {
	return nil
}

// CreateTemporaryFile xxx.
func (d *Download) CreateTemporaryFile() (*os.File, error) {
	// Increment temporary file number
	d.incrementTempFileAppender()

	// Create a new temporary file path
	tempFilePath := d.SaveFullPath() +
		".temp" +
		strconv.Itoa(d.tempFileNameAppender) +
		"." +
		TempFileFileExtension

	// Check if name for temporary file exists
	isFileExists := fileUtils.IsFileExist(tempFilePath)
	if isFileExists {
		// If file name already exists
		// Increment the temporary file appender
		return d.CreateTemporaryFile()

		// Check infinite loop?
	}

	// Check if there's enough storage to store the file
	// Use os.truncate to increase size of the new temp file

	// Create the file
	file, err := os.OpenFile(d.SaveFullPath(), os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}

	// Add this to the temporary file list
	d.appendToTempFileList(tempFilePath)

	return file, nil
}

func (d *Download) incrementTempFileAppender() {
	_ = d.setTempFileNameAppender(d.tempFileNameAppender + 1)
}

func (d *Download) appendToTempFileList(tempFilePath string) {
	_ = d.setTempFileList(append(d.tempFileList, tempFilePath))
}

// DEBUG

// DebugUrl prints download URL details.
func (d *Download) DebugUrl() {
	// Debug print URL
	fmt.Println("URL DEBUG:")
	fmt.Println("URL scheme:", d.downloadURL.Scheme)
	fmt.Println("URL host:", d.downloadURL.Host)
	fmt.Println("URL Path:", d.downloadURL.Path)
	fmt.Println()
}

// DebugHeader prints download header details.
func (d *Download) DebugHeader() {
	// Debug print response header
	fmt.Println("HEADER DEBUG:")
	fmt.Println("Is download initialized:", d.isDownloadInitialized)
	fmt.Println("Response content length:", d.response.ContentLength)
	fmt.Println("Response headers:", d.response.Header)
	fmt.Println("Response status code:", d.response.StatusCode)
	fmt.Println()
}
