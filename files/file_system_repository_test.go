package files

import (
	"io/ioutil"
	"os"
	"testing"
)

func validateLength(t *testing.T, fsr FileSystemRepository, expectedLength uint8) {
	if fsList, e := fsr.List(); e != nil {
		t.Errorf("Could not get listing: %v", e.Error())
	} else {
		if len(fsList) != int(expectedLength) {
			t.Errorf("Should have a length of %d but was %d", expectedLength, len(fsList))
		}
	}
}

func TestFileSystemRepository(t *testing.T) {
	if file, err := ioutil.TempFile("", "file*.repo"); err == nil {
		path := file.Name()
		file.Close()
		os.Remove(file.Name())

		if fsr, e := GetFileSystemRepositoryFromPath(path); e == nil {
			validateLength(t, fsr, 0)

			if b := fsr.Store(FileSystem{
				Name:        "My FS",
				Description: "This is a description",
				Addresses:   []string{"localhost:9000"},
			}); b != nil {
				t.Errorf("Could not store record: %v", b.Error())
			}

			validateLength(t, fsr, 1)

			if value, el := fsr.FindByName("My FS"); el != nil {
				t.Errorf("Could not find value in FileSystemRepository: %v", el.Error())
			} else if value.Description != "This is a description" {
				t.Error("Found an incorrect value")
			}
		} else {
			t.Errorf("Error creating File System: %v", e.Error())
		}

		if fsr, e := GetFileSystemRepositoryFromPath(path); e == nil {
			validateLength(t, fsr, 1)

			if value, b := fsr.FindByName("My FS"); b != nil {
				t.Error("Could not find value in FileSystemRepository")
			} else if value.Description != "This is a description" {
				t.Error("Found an incorrect value")
			}

			if b := fsr.Remove("My FS"); b != nil {
				t.Errorf("Error removing record: %v", b.Error())
			}
		} else {
			t.Errorf("Error creating File System: %v", e.Error())
		}
		if fsr, e := GetFileSystemRepositoryFromPath(path); e == nil {
			validateLength(t, fsr, 0)
		} else {
			t.Errorf("Error creating File System: %v", e.Error())
		}
	} else {
		t.Errorf("Error when creating temp file: %v", err.Error())
	}
}
