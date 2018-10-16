package ftp

import (
	"crypto/sha512"
	"fmt"
	"os"
)

func hash(s string) string {
	return fmt.Sprintf("%x", sha512.Sum512([]byte(s)))
}

func existsPath(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func isDirectory(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fileInfo.IsDir()
}
