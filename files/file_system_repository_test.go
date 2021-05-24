package files

import (
	"io/ioutil"
	"os"
	"testing"
)

func validateLength(t *testing.T, fsr FileSystemRepository, expectedLength uint8) {
	fsList, e := fsr.List()
	if e != nil {
		t.Errorf("Could not get listing: %v", e.Error())
	}
	if len(fsList) != int(expectedLength) {
		t.Errorf("Should have a length of %d but was %d", expectedLength, len(fsList))
	}
}

func validateRecordFound(t *testing.T, fsr FileSystemRepository) {
	fs, err := fsr.FindByName("My FS")
	if err != nil {
		t.Errorf("Unable to find record: %v", err.Error())
		t.Fail()
	}

	if fs.Description != "This is a description" {
		t.Errorf("Invalid description: %v", fs.Description)
		t.Fail()
	}
}

func createTestFile(t *testing.T) string {
	file, err := ioutil.TempFile("", "file*.repo")
	var path string
	if err == nil {
		path = file.Name()
		file.Close()
		os.Remove(file.Name())
		return path
	}
	t.Error("Unable to create temp file")
	t.Fail()
	return path
}

func getRepository(t *testing.T, path string) FileSystemRepository {
	fsr, err := GetFileSystemRepositoryFromPath(path)
	if err != nil {
		t.Errorf("Unable to build file system repository: %v", err.Error())
		t.Fail()
	}
	return fsr
}

func TestFileSystemRepository(t *testing.T) {
	path := createTestFile(t)
	fsr := getRepository(t, path)

	validateLength(t, fsr, 0)

	if b := fsr.Store(FileSystem{
		Name:        "My FS",
		Description: "This is a description",
		Addresses:   []string{"localhost:9000"},
	}); b != nil {
		t.Errorf("Could not store record: %v", b.Error())
	}

	validateLength(t, fsr, 1)
	validateRecordFound(t, fsr)

	// re-initialize
	fsr = getRepository(t, path)
	validateLength(t, fsr, 1)
	validateRecordFound(t, fsr)

	if b := fsr.Remove("My FS"); b != nil {
		t.Errorf("Error removing record: %v", b.Error())
	}

	// re-initialize
	fsr = getRepository(t, path)
	validateLength(t, fsr, 0)
}
