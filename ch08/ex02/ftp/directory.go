package ftp

import (
	"path/filepath"
	"strings"
)

type Directory struct {
	RootDir string
	WorkDir string
}

func (d *Directory) JoinPath(path string) string {
	newPath := filepath.Clean(path)
	if strings.HasPrefix("/", newPath) {
		newPath = filepath.Join(d.RootDir, newPath)
	} else {
		newPath = filepath.Join(d.RootDir, d.WorkDir, newPath)
	}
	if !strings.HasPrefix(newPath, d.RootDir) {
		newPath = d.RootDir
	}
	return newPath
}
