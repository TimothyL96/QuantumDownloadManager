package manager

import (
	"context"
	"fmt"
	"io"
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
	// Get the filename and extension.
	// Ex: Content-Disposition: form-data; name="fieldName"; filename="filename.jpg"
	//
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Disposition
	//
	contentDisposition, ok := d.response.Header["Accept-Ranges"]
	if ok {
		_ = d.setDefaultFileName(contentDisposition[0])
	} else {
		// Content-Type:[video/mp4]
		// Get file extension if Content-Disposition not specified
		//
		// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Type
		//
		contentType, ok := d.response.Header["Content-Type"]
		if ok {
			_ = contentType[0] // TODO
		} else {
			// Try to get file name and extension from URL
		}
	}

	return nil
}

// startDownload starts the download.
func (d *Download) startDownload() error {
	fmt.Println("Starting download")

	// Set the download as running
	_ = d.Resume()

	contentLength := d.response.ContentLength
	var currentByte int64 = 0
	var wg sync.WaitGroup

	if d.IsConcurrentConnectionAllowed() == notAllowed {
		_ = d.SetMaxNrOfConcurrentConnection(1)
	} else if d.IsConcurrentConnectionAllowed() == unknown {
		// Unknown check if return is 206
	}

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

		// Add to temp file list
		d.appendToTempFileList(file.Name())

		// Calculate bytes to get per concurrent connection
		// Append remaining bytes to this request for the last concurrent connection
		bytesToGet := contentLength

		// Set a minimum bytes per temporary file / concurrent connection ?

		// If this is not the last concurrent connection
		if i != 1 {
			bytesToGet = int64(math.Floor(float64(contentLength) / float64(i)))
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
			fmt.Println("Return status code is not 206 partial download but:" +
				strconv.Itoa(downloader.response.StatusCode))
		}

		// Start the concurrent download with the new bytes range calculated above
		// Use a closure in the below anonymous function / goroutine : (i int)
		// to store the concurrent connection index (Same value as 'i' in the current for loop)

		// Track current running goroutine to its completion
		wg.Add(1)
		go downloader.streamDownload(file, &wg)

		// Stop adding concurrent connection if return is 200 instead of 206
		if downloader.response.StatusCode == 200 {
			break
		}
	}

	// Wait for all tracked goroutines to be completed
	wg.Wait()

	// Combine files and get the final download file
	if err := d.combineFiles(); err != nil {
		return err
	}

	// Set download as completed
	d.complete()

	return nil
}

func (d *Download) streamDownload(file *os.File, wg *sync.WaitGroup) {
	// Write the specific data range to disk
	//
	// If file size is 1 GB and user has 1.2 GB disk space left,
	// This might cause the space used to become 2 GB with 1 GB for save file and 1 GB for other temporary files.
	// Could implement a read line by line and removing the read line from temporary files
	// to allow user with 1.2 GB disk space download a 1 GB file concurrently
	_, _ = io.Copy(file, d.response.Body)

	// Close the temporary file
	_ = file.Close()

	// Set current goroutine as completed
	wg.Done()
}
