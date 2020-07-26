package main

import (
	"fmt"

	"github.com/ttimt/QuantumDownloadManager/internal/app/downloader/manager"
	"github.com/ttimt/QuantumDownloadManager/internal/app/downloader/user/setting"
)

func main() {
	// Testing concurrent download
	test()
}

func test() {
	// Test
	const (
		nrOfConcurrentDownload = 64
		fileName               = "download.mp4"
		directory              = `D:\Timothy/Desktop\`
	)

	userSetting1, err := setting.NewSetting(
		setting.NrOfConcurrentConnection(nrOfConcurrentDownload))
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

	// Initialize downloader new download
	downloader, err := manager.NewDownload(
		manager.DownloadURL(url),
		manager.NrOfConcurrentDownload(userSetting1.NrOfConcurrentConnection()),
		manager.SaveDirectory(directory),
		manager.SaveFileName(fileName))
	if err != nil {
		panic(err)
	}

	// Retrieve download details
	if err = downloader.Initialize(); err != nil {
		panic(err)
	}

	downloader.DebugHeader()
	fmt.Println(downloader)

	// Start the download
	err = downloader.Download()
	if err != nil {
		panic(err)
	}
	fmt.Println("Download complete")
}
