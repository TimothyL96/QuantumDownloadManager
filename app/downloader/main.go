package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi"

	"github.com/ttimt/QuantumDownloadManager/internal/app/downloader/manager"
	"github.com/ttimt/QuantumDownloadManager/internal/app/downloader/user/setting"
)

type downloadDTO struct {
	Url string
}

func main() {
	// Testing concurrent download
	// test()

	fmt.Println("Running")
	r := chi.NewRouter()

	r.Post("/download", func(w http.ResponseWriter, r *http.Request) {
		var url downloadDTO
		err := json.NewDecoder(r.Body).Decode(&url)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("received", url.Url)
		_, _ = io.WriteString(w, url.Url)
	})

	output := TestOutput{
		Status:     1,
		Percentage: 100.1,
	}

	r.Get("/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(output)
	})

	log.Fatal(http.ListenAndServe(":3333", r))
}

// TestOutput for GET API
type TestOutput struct {
	Status     int     `json:"status"`
	Percentage float32 `json:"percentage"`
}

func test() {
	// Test
	const (
		nrOfConcurrentDownload = 8
		fileName               = "download.mp4"
		directory              = `D:\Timothy/Desktop\`
	)

	userSetting1, err := setting.NewSetting(
		setting.NrOfConcurrentConnection(nrOfConcurrentDownload))
	if err != nil {
		panic(err)
	}

	// Sample download URLs:
	// https://file-examples-com.github.io/uploads/2017/10/file-sample_150kB.pdf
	// https://file-examples-com.github.io/uploads/2017/02/file-sample_100kB.doc
	// https://file-examples-com.github.io/uploads/2017/04/file_example_MP4_1920_18MG.mp4

	// Test: Set URL in code
	url := "https://file-examples-com.github.io/uploads/2017/04/file_example_MP4_1920_18MG.mp4"
	fmt.Printf("URL is: %s\n\n", url)

	// Initialize new download
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

	fmt.Println(downloader)
	downloader.DebugHeader()
	downloader.DebugFileSize()

	// Start the download
	err = downloader.Start()
	if err != nil {
		panic(err)
	}

	fmt.Println("Download complete:", downloader.IsDownloadComplete())
}
