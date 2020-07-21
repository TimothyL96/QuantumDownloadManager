package file

import (
	"path/filepath"
	"testing"
)

func TestIsFileExist(t *testing.T) {
	path := "D:\\Timothy\\Desktop\\download.mp4"

	path = filepath.Clean(path)

	t.Log("Path:", path)

	get := IsFileExist(path)

	if !get {
		t.Errorf("Want %t, got %t", !get, get)
	}
}

func TestIsFileExist1(t *testing.T) {
	path := "D:\\Timothy\\Desktop\\downloa/d.mp4"

	path = filepath.Clean(path)

	t.Log("Path:", path)

	get := IsFileExist(path)

	if get {
		t.Errorf("Want %t, got %t", !get, get)
	}
}

func TestIsFileExist2(t *testing.T) {
	path := "D:\\Timothy\\Desktop\\downloa\\d.mp4"

	path = filepath.Clean(path)

	t.Log("Path:", path)

	get := IsFileExist(path)

	if get {
		t.Errorf("Want %t, got %t", !get, get)
	}
}

func TestIsFileExist3(t *testing.T) {
	path := "D:\\Timothy\\Desktop\\download.mp4\\"

	path = filepath.Clean(path)

	t.Log("Path:", path)

	get := IsFileExist(path)

	if !get {
		t.Errorf("Want %t, got %t", !get, get)
	}
}

func TestIsFileExist4(t *testing.T) {
	path := "D:\\Timothy\\Desktop\\download.mp4/"

	path = filepath.Clean(path)

	t.Log("Path:", path)

	get := IsFileExist(path)

	if !get {
		t.Errorf("Want %t, got %t", !get, get)
	}
}
