package downloadManager

import (
	"context"
	"fmt"
	"net/http"
)

// HTTP Header : Accept-Ranges
// If absence or value is 'no', partial request / pause download is not supported
//
// Client add HTTP Header : Range bytes=0-999
// Server returns:
// HTTP/1.0 206 Partial Content
// Accept-Ranges: bytes
// Content-Length: 1000
// Content-Range: bytes 0-999/2200

func (d *DownloadManager) Download() error {
	// Initialize download and get response header
	err := d.InitializeDownload()
	if err != nil {
		return err
	}

	// Start the download
	err = d.StartConcurrentDownload()
	if err != nil {
		return err
	}

	return nil
}

func (d *DownloadManager) InitializeDownload() error {
	// Setup new context for download stoppage
	d.ctx, d.ctxCancel = context.WithCancel(context.Background())

	// Setup request with instance's context
	req, err := http.NewRequestWithContext(d.ctx,
		http.MethodGet,
		d.downloadUrl.String(),
		nil)
	if err != nil {
		return err
	}

	// Make the request to get the response header
	d.response, err = http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	// Check if download can be split to concurrent download
	// and can it be potentially paused
	acceptRanges, ok := d.response.Header["Accept-Ranges"]

	// Concurrent download is allowed and download might be able to be paused
	if ok && len(acceptRanges) > 0 && acceptRanges[0] != "no" {
		d.isConcurrentAllowed = true
		d.isPausedAllowed = true
	}

	// Check if download can be paused and resumed
	if !d.isPausedAllowed || d.response.ContentLength == -1 {
		// File size unknown and cannot be paused
		d.isPausedAllowed = false
	} else {
		// Download can be paused
		d.isPausedAllowed = true
	}

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

func (d *DownloadManager) DebugUrl() {
	// Debug print URL
	fmt.Println("URL scheme:", d.downloadUrl.Scheme)
	fmt.Println("URL host:", d.downloadUrl.Host)
	fmt.Println("URL Path:", d.downloadUrl.Path)
}
