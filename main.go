package main

import (
	"fmt"

	"github.com/ttimt/QuantumDownloadManager/downloadManager"
	"github.com/ttimt/QuantumDownloadManager/userSettings"
)

func main() {
	Test()
}

func Test() {
	// Test
	const (
		nrOfConcurrentDownload = 1
		fileName               = "download.mp4"
		directory              = "D:\\Timothy/Desktop\\"
	)

	userSetting1, err := userSettings.NewUserSettings(
		userSettings.NrOfConcurrentDownload(nrOfConcurrentDownload))
	if err != nil {
		panic(err)
	}

	var url string

	// Get URL from user
	// fmt.Println("Enter download URL:")
	// reader := bufio.NewScanner(os.Stdin)
	// reader.Scan()
	// url = reader.Text()
	// fmt.Println("URL given is:", url)

	// Sample download URLs:
	// https://file-examples-com.github.io/uploads/2017/10/file-sample_150kB.pdf
	// https://file-examples-com.github.io/uploads/2017/02/file-sample_100kB.doc
	// https://file-examples-com.github.io/uploads/2017/04/file_example_MP4_1920_18MG.mp4

	// Test: Set URL in code
	url = "https://file-examples-com.github.io/uploads/2017/04/file_example_MP4_1920_18MG.mp4"
	fmt.Printf("URL is: %s\n\n", url)

	// Initialize a new download
	downloader, err := downloadManager.NewDownloadManager(
		downloadManager.DownloadUrl(url),
		downloadManager.NrOfConcurrentDownload(userSetting1.GetNrOfConcurrentDownload()),
		downloadManager.SaveDirectory(directory),
		downloadManager.SaveFileName(fileName))
	if err != nil {
		panic(err)
	}

	// Retrieve download details
	err = downloader.RetrieveDownloadDetails()
	if err != nil {
		panic(err)
	}

	// Start the download
	// err = downloader.Download()
	// if err != nil {
	// 	panic(err)
	// }
}
