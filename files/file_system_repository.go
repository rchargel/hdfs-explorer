package files

import (
	"encoding/gob"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// FileSystemRepository is used to load/store/find all of the previously configured FileSystems.
type FileSystemRepository interface {
	List() ([]FileSystem, error)
	FindByName(name string) (*FileSystem, error)
	Remove(name string) error
	Store(fileSystem FileSystem) error
}

type fileSystemRepository struct {
	path string
}

func (f *fileSystemRepository) List() ([]FileSystem, error) {
	return f.load()
}

func (f *fileSystemRepository) FindByName(name string) (*FileSystem, error) {
	if values, err := f.load(); err == nil {
		for _, value := range values {
			if value.Name == name {
				return &value, nil
			}
		}
	} else {
		return nil, err
	}
	return nil, errors.New(fmt.Sprintf("Could not find a FileSystem named %v", name))
}

func (f *fileSystemRepository) Remove(name string) error {
	if values, err := f.load(); err == nil {
		saved := make([]FileSystem, 0)
		for _, value := range values {
			if value.Name != name {
				saved = append(saved, value)
			}
		}

		return f.save(saved)
	} else {
		return err
	}
}

func (f *fileSystemRepository) Store(fileSystem FileSystem) error {
	log.Printf("Storing file system: %v", fileSystem.Name)
	if values, err := f.load(); err == nil {
		saved := make([]FileSystem, 0)
		added := false
		for _, value := range values {
			if value.Name == fileSystem.Name {
				saved = append(saved, fileSystem)
				added = true
			} else {
				saved = append(saved, value)
			}
		}
		if !added {
			saved = append(saved, fileSystem)
		}
		return f.save(saved)
	} else {
		return err
	}
}

func (f *fileSystemRepository) load() ([]FileSystem, error) {
	fileSystems := make([]FileSystem, 0)

	if file, err := os.Open(f.path); err == nil {
		defer file.Close()
		dec := gob.NewDecoder(file)
		fs := &FileSystem{}
		for {
			err := dec.Decode(fs)
			if err != nil {
				break
			}

			fileSystems = append(fileSystems, *fs)
		}
	} else {
		return nil, err
	}

	return fileSystems, nil
}

func (f *fileSystemRepository) save(fileSystems []FileSystem) error {
	if tmp, err := ioutil.TempFile("", "filerepo*.tmp"); err == nil {
		enc := gob.NewEncoder(tmp)
		for _, value := range fileSystems {
			err = enc.Encode(value)
			if err != nil {
				return err
			}
		}
		tmp.Close()
		os.Remove(f.path)
		return os.Rename(tmp.Name(), f.path)
	} else {
		return err
	}
}

// GetFileSystemRepository loads an instance of the FileSystemRepository interface.
// This call will default to reading the data from $HOME/.fsrepo/fs.repo.
func GetFileSystemRepository() (FileSystemRepository, error) {
	dir, _ := os.UserHomeDir()
	filepath.Join(dir, ".fsrepo", "fs.repo")
	return GetFileSystemRepositoryFromPath(dir)
}

// GetFileSystemRepositoryFromPath loads an instance of the FileSystemRepository found at a specified
// location.
func GetFileSystemRepositoryFromPath(path string) (FileSystemRepository, error) {
	if file, err := GetOrCreateFile(path); err == nil {
		file.Close()
		return &fileSystemRepository{path}, nil
	} else {
		return nil, err
	}
}
