package files

import (
	"encoding/gob"
	"os"
	"path/filepath"
	"strings"
)

// FileSystemRepository is used to load/store/find all of the previously configured FileSystems.
type FileSystemRepository interface {
	List() []FileSystem
	FindByName(name string) []FileSystem
	Store(fileSystem FileSystem) error
}

type fileSystemRepository struct {
	path       string
	registered []FileSystem
}

func (f *fileSystemRepository) List() []FileSystem {
	return f.registered
}

func (f *fileSystemRepository) FindByName(name string) []FileSystem {
	filtered := []FileSystem{}

	for _, value := range f.registered {
		if strings.Contains(strings.ToLower(value.Name), strings.ToLower(name)) {
			filtered = append(filtered, value)
		}
	}
	return filtered
}

func (f *fileSystemRepository) Store(fileSystem FileSystem) error {
	found := false
	for index := range f.registered {
		if f.registered[index].Name == fileSystem.Name {
			f.registered[index] = fileSystem
			found = true
			break
		}
	}
	if !found {
		f.registered = append(f.registered, fileSystem)
	}
	return f.save()
}

func (f *fileSystemRepository) save() error {
	if file, err := os.Create(f.path); err != nil {
		return nil
	} else {
		defer file.Close()
		enc := gob.NewEncoder(file)
		for _, value := range f.registered {
			err = enc.Encode(value)
			if err != nil {
				return err
			}
		}
	}
	return nil
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
	if file, err := GetOrCreateFile(path); err != nil {
		return nil, err
	} else {
		defer file.Close()
		dec := gob.NewDecoder(file)
		var fileSystem *FileSystem

		values := []FileSystem{}
		for {
			err := dec.Decode(fileSystem)
			if err != nil {
				break
			}
			values = append(values, *fileSystem)
		}
		return &fileSystemRepository{path, values}, nil
	}
}
