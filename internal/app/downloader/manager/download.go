package manager

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"sync"
)

// sendHTTPRequest sends a HTTP request with custom header from parameter
// and stores the response in downloader.Response
func (d *Download) sendHTTPRequest(header map[string]string) error {
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

	// Add custom header to the request
	for k, v := range header {
		req.Header.Add(k, v)
	}

	// DEBUG: Print request header
	fmt.Println("Request header")
	fmt.Println(req.Header)
	fmt.Println()

	// Make the request to get the response header
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	_ = d.setResponse(response)

	return nil
}

func (d *Download) processRequestHeader() error {
	// Partial download reference:
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Range_requests
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
				_ = d.setIsPauseAllowed(notAllowed)
			} else {
				// Allowed
				_ = d.setIsConcurrentConnectionAllowed(allowed)
				_ = d.setIsPauseAllowed(allowed)
			}
		} else {
			// Unknown
			_ = d.setIsConcurrentConnectionAllowed(unknown)
			_ = d.setIsPauseAllowed(unknown)
		}

		// Update file size
		_ = d.setFileSize(d.response.ContentLength)
	} else {
		// Unknown content length
		_ = d.setIsConcurrentConnectionAllowed(notAllowed)
		_ = d.setIsPauseAllowed(notAllowed)
	}

	// Get suggested default file name from header - Content-Disposition
	// d.setDefaultFileName(...)

	return nil
}

// startSequentialDownload downloads the file without any concurrent connection.
func (d *Download) startSequentialDownload() error {
	// Check if download has been started before, and resume the last pause state
	// To be done when implementing pause feature

	fmt.Println("Starting sequential download")

	// Set the download as running
	_ = d.setIsDownloadRunning(true)

	// Create downloader single temporary file
	tempFile, err := d.createTemporaryFile()
	if err != nil {
		d.Abort()
		return err
	}

	// Close temp file
	defer tempFile.Close()

	// Write the data to disk
	written, err := io.Copy(tempFile, d.response.Body)
	if err != nil {
		d.Abort()
		return err
	}

	// Rename temp file to download save file name
	if err = os.Rename(tempFile.Name(), d.SaveFullPath()); err != nil {
		d.Abort()
		return err
	}

	log.Printf("DEBUG: File downloaded with a single connection."+
		"\nWritten %d out of %d bytes", written, d.response.ContentLength)

	// Set download completed
	d.complete()

	return nil
}

// startConcurrentDownload download the file part by part concurrently with the number of concurrent connection set.
func (d *Download) startConcurrentDownload() error {
	fmt.Println("Starting concurrent download")

	// Set the download as running
	_ = d.setIsDownloadRunning(true)

	contentLength := d.response.ContentLength
	var currentByte int64 = 0

	// Sync wait group
	var wg sync.WaitGroup

	for i := d.MaxNrOfConcurrentConnection(); i > 0; i-- {
		// Create a new downloader with custom header to specify the custom bytes range for concurrent download
		downloader, err := NewDownload(
			SaveDirectory(d.SaveDirectory()),
			SaveFileName(d.SaveFileName()),
			NrOfConcurrentDownload(1),
			DownloadURL(d.DownloadURL()))
		if err != nil {
			return err
		}

		// Get a temporary file name
		file, err := downloader.createTemporaryFile()
		if err != nil {
			downloader.Abort()
			d.Abort()
			return err
		}

		// Calculate bytes to get per concurrent connection
		var bytesToGet int64

		// Set a minimum bytes per temporary file / concurrent connection ?

		// If this is not the last concurrent connection
		if i != 1 {
			bytesToGet = int64(math.Floor(float64(contentLength) / float64(i)))
		} else {
			// Append remaining bytes to this request for the last concurrent connection
			bytesToGet = contentLength
		}

		// Send a HTTP request with custom header to get the new response header
		if err = downloader.sendHTTPRequest(map[string]string{"Range": "bytes=" +
			strconv.FormatInt(currentByte, 10) +
			"-" +
			strconv.FormatInt(currentByte+(bytesToGet-1), 10)}); err != nil {
			return err
		}

		// Update remaining content length and current byte
		contentLength -= bytesToGet
		currentByte += bytesToGet

		// Check status code is 206 Partial Content before moving to next concurrent connection
		// 200 - Partial download not supported
		// 206 - Successful request
		// 416 - Requested Range Not Satisfiable (Not of the requested range values overlap the available range)
		if downloader.response.StatusCode != 206 {
			downloader.Abort()
			d.Abort()
			return errors.New("Return status code is not 206 partial download but:" +
				strconv.Itoa(downloader.response.StatusCode))
		}

		// Add to temp file list
		d.appendToTempFileList(file.Name())

		// Start the concurrent download with the new bytes range calculated above
		// Use a closure in the below anonymous function / goroutine : (i int)
		// to store the concurrent connection index (Same value as 'i' in the current for loop)
		go func(i int) {
			fmt.Println("***** Starting concurrent download:", i)

			// DEBUG
			downloader.DebugHeader()

			// Track current running goroutine to its completion
			wg.Add(1)

			// Write the specific data range to disk
			//
			// If file size is 1 GB and user has 1.2 GB disk space left,
			// This might cause the space used to become 2 GB with 1 GB for save file and 1 GB for other temporary files.
			// Could implement a read line by line and removing the read line from temporary files
			// to allow user with 1.2 GB disk space download a 1 GB file concurrently
			_, _ = io.Copy(file, downloader.response.Body)

			// Close the temporary file
			_ = file.Close()

			fmt.Println("Closing concurrent download:", i)

			// Set current goroutine as completed
			wg.Done()
		}(i)
	}

	// Wait for all tracked goroutines to be completed
	wg.Wait()

	// Combine files
	fmt.Println("Temp file list:", d.tempFileList)

	// Open the download save file to save and combine the completed download parts
	// Currently, replace file if exist, can append a value later on: filename_1
	file, err := os.OpenFile(d.SaveFullPath(), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		panic(err)
	}

	fmt.Println("Writing to file:", file.Name())

	// Go through each temporary file
	for _, v := range d.tempFileList {
		file1, err := os.OpenFile(v, os.O_RDONLY, os.ModePerm)
		if err != nil {
			panic(err)
		}

		// Copy all bytes from the temporary file to the final download file
		_, _ = io.Copy(file, file1)

		// Close the temporary file after copying
		_ = file1.Close()

		// Remove the temporary file that has been copied
		// If combine all bytes before deleting all temporary files, file size could go double.
		if err = os.Remove(file1.Name()); err != nil {
			fmt.Println("Failed to delete temporary file:", err)
		}
	}

	// Set download as completed
	_ = d.setIsDownloadComplete(true)

	return nil
}

// streamData xxx.
func (d *Download) streamData() error {
	return nil
}
