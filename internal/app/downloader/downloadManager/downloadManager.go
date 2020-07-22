package downloadManager

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

// *********** Initializer
func (d *DownloadManager) RetrieveDownloadDetails() error {
	if d.GetIsDownloadStarted() {
		return errors.New("download has already been initialized")
	}

	// Initialize download and get response header
	err := d.InitializeDownload()
	if err != nil {
		return err
	}

	// Debug
	d.DebugUrl()
	d.DebugHeader()

	return nil
}

func (d *DownloadManager) InitializeDownload() error {
	// Setup new context for stopping download
	ctx, ctxCancel := context.WithCancel(context.Background())
	_ = d.setCtx(ctx)
	_ = d.setCtxCancel(ctxCancel)

	// Setup request with the newly created instance's context
	req, err := http.NewRequestWithContext(d.ctx,
		http.MethodGet,
		d.downloadUrl.String(),
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
				_ = d.setIsConcurrentAllowed(notAllowed)
				_ = d.setIsPausedAllowed(notAllowed)
			} else {
				// Allowed
				_ = d.setIsConcurrentAllowed(allowed)
				_ = d.setIsPausedAllowed(allowed)
			}
		} else {
			// Unknown
			_ = d.setIsConcurrentAllowed(unknown)
			_ = d.setIsPausedAllowed(unknown)
		}
	} else {
		// Unknown content length
		_ = d.setIsConcurrentAllowed(notAllowed)
		_ = d.setIsPausedAllowed(notAllowed)
	}

	// Get suggested default file name from header - Content-Disposition

	return nil
}

// *********** Main downloader
func (d *DownloadManager) Download() error {
	// Test
	_ = d.setIsConcurrentAllowed(notAllowed)

	// Block if download has started before
	if d.GetIsDownloadStarted() {
		return errors.New("download has already started before. Did you mean resume download ")
	}

	// Flag the download has started
	_ = d.SetIsDownloadStarted(true)

	// Start the download
	if d.GetIsConcurrentAllowed() == allowed {
		return d.StartConcurrentDownload()
	} else {
		return d.StartAtomicDownload()
	}
}

func (d *DownloadManager) StartAtomicDownload() error {
	// Create downloader single temporary file
	tempFile, err := d.CreateTemporaryFile()
	if err != nil {
		return err
	}

	// Write the data to disk
	written, err := io.Copy(tempFile, d.response.Body)
	if err != nil {
		return err
	}

	log.Printf("Written %d out of %d bytes", written, d.getResponse().ContentLength)

	return nil
}

func (d *DownloadManager) StartConcurrentDownload() error {
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

func (d *DownloadManager) StreamData() error {
	return nil
}

func (d *DownloadManager) CreateTemporaryFile() (*os.File, error) {
	// Increment temporary file number
	_ = d.incrementTempAppender()

	// Create downloader new temporary file path
	tempFilePath := d.GetSaveFullPath() + ".temp" +
		strconv.Itoa(d.getTempAppender()) +
		".qdm"

	// Check if file exists
	isFileExists := fileUtils.IsFileExist(tempFilePath)
	if isFileExists {
		// If file name already exists
		// Increment the temporary file appender
		return d.CreateTemporaryFile()
	}

	// Check if there's enough storage to store the file
	// Use os.truncate to increase size of the new temp file

	// Create the file
	file, err := os.OpenFile(d.GetSaveFullPath(), os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}

	// Add this to the temporary file list
	_ = d.setTempFileList(tempFilePath)

	return file, nil
}

// *********** DEBUG
func (d *DownloadManager) DebugUrl() {
	// Debug print URL
	fmt.Println("URL DEBUG:")
	fmt.Println("URL scheme:", d.downloadUrl.Scheme)
	fmt.Println("URL host:", d.downloadUrl.Host)
	fmt.Println("URL Path:", d.downloadUrl.Path)
	fmt.Println()
}

func (d *DownloadManager) DebugHeader() {
	// Debug print response header
	fmt.Println("HEADER DEBUG:")
	fmt.Println("Response content length:", d.response.ContentLength)
	fmt.Println("Response headers:", d.response.Header)
	fmt.Println("Response status code:", d.response.StatusCode)
	fmt.Println()
}
