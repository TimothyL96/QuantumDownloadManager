package manager

// ConfigOption is the signature of functional option for Start.
type ConfigOption func(d *Download) error

// NewDownload create and returns a new Start instance with configurations from the parameter input.
func NewDownload(configurations ...ConfigOption) (*Download, error) {
	download := &Download{}

	for _, configuration := range configurations {
		err := configuration(download)

		if err != nil {
			return nil, err
		}
	}

	// Update save full path after updating save directory and save file name
	err := download.setSaveFullPath()
	if err != nil {
		return nil, err
	}

	return download, nil
}

// Functional options functions:

// DownloadURL allows setting the value of download URL.
func DownloadURL(downloadURL string) ConfigOption {
	return func(d *Download) error {
		return d.SetDownloadURL(downloadURL)
	}
}

// NrOfConcurrentDownload allows setting the value of number of concurrent download.
func NrOfConcurrentDownload(nrOfConcurrentDownload int) ConfigOption {
	return func(d *Download) error {
		return d.SetMaxNrOfConcurrentConnection(nrOfConcurrentDownload)
	}
}

// SaveDirectory allows setting the value of download directory.
func SaveDirectory(saveDirectory string) ConfigOption {
	return func(d *Download) error {
		return d.SetSaveDirectory(saveDirectory)
	}
}

// SaveFileName allows setting the value of download file name.
func SaveFileName(saveFileName string) ConfigOption {
	return func(d *Download) error {
		return d.SetSaveFileName(saveFileName)
	}
}
