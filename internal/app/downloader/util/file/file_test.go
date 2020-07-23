package file_test

import (
	"path/filepath"
	"testing"

	"github.com/ttimt/QuantumDownloadManager/internal/app/downloader/util/file"
)

func TestIsFileExist(t *testing.T) {
	var testCases = []struct {
		name string
		want bool
		path string
	}{
		{name: "T1", want: false, path: `D:\Timothy\Desktop\download1.mp4`},
		{name: "T1", want: false, path: `D:\Timothy\Desktop\download2.mp4`},
		{name: "T1", want: false, path: `D:\Timothy\Desktop\download3.mp4`},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			path := filepath.Clean(testCase.path)
			get := file.IsFileExist(path)

			if testCase.want != get {
				t.Errorf("Want %t, got %t", testCase.want, get)
			}
		})
	}
}
