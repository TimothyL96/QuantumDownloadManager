package downloadManager

type downloadManagerConfigurationFunction func(d *DownloadManager) error

func NewDownloadManager(configurationFunctions ...downloadManagerConfigurationFunction) (*DownloadManager, error) {
	downloader := &DownloadManager{}

	for _, configurationFunction := range configurationFunctions {
		err := configurationFunction(downloader)

		if err != nil {
			return nil, err
		}
	}

	// Update save full path after updating save directory and save file name
	err := downloader.setSaveFullPath()
	if err != nil {
		return nil, err
	}

	return downloader, nil
}

func DownloadUrl(downloadUrl string) downloadManagerConfigurationFunction {
	return func(d *DownloadManager) error {
		return d.SetDownloadUrl(downloadUrl)
	}
}

func NrOfConcurrentDownload(nrOfConcurrentDownload int) downloadManagerConfigurationFunction {
	return func(d *DownloadManager) error {
		return d.SetNrOfConcurrentDownload(nrOfConcurrentDownload)
	}
}

func SaveDirectory(saveDirectory string) downloadManagerConfigurationFunction {
	return func(d *DownloadManager) error {
		return d.SetSaveDirectory(saveDirectory)
	}
}

func SaveFileName(saveFileName string) downloadManagerConfigurationFunction {
	return func(d *DownloadManager) error {
		return d.SetSaveFileName(saveFileName)
	}
}
