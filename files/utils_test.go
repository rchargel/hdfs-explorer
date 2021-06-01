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

	validateContent(t, file.Name(), "This is a line of text")

	removeFile(t, file.Name())
	validateEmptyFile(t, file.Name())
	os.Remove(file.Name())
}

func TestJoinAndParent(t *testing.T) {
	path := "/"
	if Join(path, "file") != "/file" {
		t.Errorf("Join failed %v", Join(path, "file"))
		t.Fail()
	}
	if Join(Join(path, "dir"), "file") != "/dir/file" {
		t.Errorf("Join failed %v", Join(Join(path, "dir"), "file"))
		t.Fail()
	}
	if Parent(Join(Join(path, "dir"), "file")) != "/dir" {
		t.Errorf("Join failed %v", Parent(Join(Join(path, "dir"), "file")))
		t.Fail()
	}
	if Parent(path) != "/" {
		t.Errorf("Join failed %v", Parent(path))
		t.Fail()
	}
}

func TestFormatBytes(t *testing.T) {
	if "100 B" != FormatBytes(100) {
		t.Errorf("Failed format %d to %v", 100, FormatBytes(100))
		t.Fail()
	}
	if "1.0 KB" != FormatBytes(1024) {
		t.Errorf("Failed format %d to %v", 1024, FormatBytes(1024))
		t.Fail()
	}
	if "2.0 MB" != FormatBytes(1024*1024*2) {
		t.Errorf("Failed format %d to %v", 1024*1024*2, FormatBytes(1024*1024*2))
		t.Fail()
	}
}

func validateContent(t *testing.T, path, content string) {
	f, err := GetOrCreateFile(path)
	defer f.Close()

	if err != nil {
		t.Error(err)
		t.Fail()
	} else if data, _ := ioutil.ReadFile(f.Name()); string(data) != content {
		t.Errorf("File should have the content %v but was %v", content, string(data))
		t.Fail()
	}
}

func validateEmptyFile(t *testing.T, path string) {
	f, err := GetOrCreateFile(path)
	defer f.Close()

	if err != nil {
		t.Error(err)
		t.Fail()
	} else if data, _ := ioutil.ReadFile(f.Name()); len(data) != 0 {
		t.Errorf("File should be empty but was %d bytes", len(data))
		t.Fail()
	}
}

func removeFile(t *testing.T, path string) {
	if err := os.Remove(path); err != nil {
		t.Errorf("Unable to remove file: %v", err.Error())
		t.Fail()
	}
}
