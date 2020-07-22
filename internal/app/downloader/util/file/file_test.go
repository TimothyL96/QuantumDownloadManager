package file

import (
	"path/filepath"
	"testing"
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
			testCase.path = filepath.Clean(testCase.path)

			get := IsFileExist(testCase.path)

			if testCase.want != get {
				t.Errorf("Want %t, got %t", testCase.want, get)
			}
		})
	}
}
