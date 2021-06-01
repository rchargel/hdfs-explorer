package files

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
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

func Parent(path string) string {
	return strings.ReplaceAll(filepath.Dir(path), "\\", "/")
}

func Join(path, file string) string {
	return strings.ReplaceAll(filepath.Join(path, file), "\\", "/")
}

func FormatBytes(size uint64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := uint64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}
