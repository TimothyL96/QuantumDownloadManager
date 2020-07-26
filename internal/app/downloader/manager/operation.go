package manager

import (
	"errors"
)

// Initialize initialize the new download by sending a request to the download URL
// and updating the download fields value by processing the received header.
func (d *Download) Initialize() error {
	if d.isDownloadRunning {
		return errors.New("download is currently running")
	}

	// Send a HTTP request to get the request header
	if err := d.sendHTTPRequest(nil); err != nil {
		return err
	}

	// Process the received request header
	if err := d.processRequestHeader(); err != nil {
		return err
	}

	// Set download as initialized
	return d.setIsDownloadInitialized(true)
}

// Download starts the download first time.
func (d *Download) Download() error {
	// Test
	// _ = d.setIsConcurrentConnectionAllowed(notAllowed)

	// Block if download has started before
	if d.IsDownloadStarted() {
		return errors.New("download has already started before. Did you mean resume download ")
	}

	// Flag the download has started
	_ = d.setIsDownloadStarted(true)

	// Create a place holder file
	d.createPlaceHolderFile()

	// Start the download concurrently
	if (d.IsConcurrentConnectionAllowed() == allowed || d.IsConcurrentConnectionAllowed() == unknown) && d.MaxNrOfConcurrentConnection() > 1 {
		go d.startConcurrentDownload()
	}

	// Non concurrent single connection download
	go d.startAtomicDownload()

	return nil
}

// Pause will pause the current download if it is currently running.
func (d *Download) Pause() error {
	_ = d.setIsDownloadRunning(false)

	return nil
}

// Resume will continue the current download if it is paused.
func (d *Download) Resume() error {
	_ = d.setIsDownloadRunning(true)

	return nil
}

// Abort will cancel the current download.
func (d *Download) Abort() {
	d.ctxCancel()
	_ = d.setIsDownloadAborted(true)
	_ = d.setIsDownloadRunning(false)
}

// complete will update the download status to complete.
func (d *Download) complete() {
	_ = d.setIsDownloadCompleted(true)
	_ = d.setIsDownloadRunning(false)
}
