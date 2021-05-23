package files

import (
	"os"
	"path/filepath"
)

// FileExists validates that the specified path exists.
// Path may be a directory or file.
func FileExists(path string) (bool, error) {
	if _, err := os.Stat(path); err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	} else {
		return false, err
	}
}

// GetOrCreateFile gets a pointer to the file at the path specified.
// If the path does not exist, it is created.
// If the parent directory or directories don't exist, they are also created.
func GetOrCreateFile(path string) (*os.File, error) {
	parent := filepath.Dir(path)
	if exists, _ := FileExists(parent); !exists {
		err := os.MkdirAll(parent, 0644)
		if err != nil {
			return nil, err
		}
	}
	if exists, _ := FileExists(path); !exists {
		return os.Create(path)
	}
	return os.Open(path)
}
