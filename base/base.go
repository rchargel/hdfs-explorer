package base

import (
	"os"
	"path/filepath"
)

var (
	HomeDir string
)

func init() {
	dir := os.Getenv("HOME_OVERRIDE")
	if dir == "" {
		dir, _ = os.UserHomeDir()
	}
	HomeDir = filepath.Join(dir, ".hsexplorer")
	os.MkdirAll(HomeDir, 0755)
}
