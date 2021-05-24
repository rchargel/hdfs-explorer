package files

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestDirectoryExists(t *testing.T) {
	dir, _ := ioutil.TempDir("", "*")

	if exists, err := FileExists(dir); err != nil {
		t.Errorf("Should not be an error: %v", err.Error())
	} else {
		if !exists {
			t.Errorf("Directory should exist: %v", dir)
		}
	}

	os.RemoveAll(dir)

	if exists, err := FileExists(dir); err != nil {
		t.Errorf("Should not be an error: %v", err.Error())
	} else {
		if exists {
			t.Errorf("Directory should not exist: %v", dir)
		}
	}
}

func TestFileExists(t *testing.T) {
	file, _ := ioutil.TempFile("", "test*")
	file.Close()

	if exists, err := FileExists(file.Name()); err != nil {
		t.Errorf("Should not be an error: %v", err.Error())
	} else {
		if !exists {
			t.Errorf("File should exist: %v", file.Name())
		}
	}

	os.Remove(file.Name())

	if exists, err := FileExists(file.Name()); err != nil {
		t.Errorf("Should not be an error: %v", err.Error())
	} else {
		if exists {
			t.Errorf("File should not exist: %v", file.Name())
		}
	}
}

func TestGetOrCreateFile(t *testing.T) {
	file, _ := ioutil.TempFile("", "test*")
	file.WriteString("This is a line of text")
	file.Close()

	if f, err := GetOrCreateFile(file.Name()); err != nil {
		t.Errorf("Should not be an error: %v", err.Error())
	} else {
		data, _ := ioutil.ReadFile(f.Name())
		t.Logf("File length is %d", len(data))
		if len(data) == 0 || string(data) != "This is a line of text" {
			t.Error("File should already have content")
		}
		f.Close()
	}

	if err := os.Remove(file.Name()); err != nil {
		t.Errorf("Unable to remove file: %v", err.Error())
	}

	if f, err := GetOrCreateFile(file.Name()); err != nil {
		t.Errorf("Should not be an error: %v", err.Error())
	} else {
		data, _ := ioutil.ReadFile(f.Name())
		if len(data) != 0 {
			t.Errorf("File should be empty but was %d", len(data))
		}
		f.Close()
	}
	os.Remove(file.Name())
}
