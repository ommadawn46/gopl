package util

import (
	"crypto/sha512"
	"fmt"
	"math/rand"
	"os"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func randomBytes(n int) []byte {
	buf := make([]byte, n)
	rand.Read(buf)
	return buf
}

func GenerateSalt() string {
	salt := ""
	buf := randomBytes(16)
	for _, b := range buf {
		salt += fmt.Sprintf("%02x", b)
	}
	return salt
}

func Hash(s string) string {
	return fmt.Sprintf("%x", sha512.Sum512([]byte(s)))
}

func ExistsPath(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func IsDirectory(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fileInfo.IsDir()
}
