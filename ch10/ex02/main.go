package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/ommadawn46/the_go_programming_language-training/ch10/ex02/unarchive"
	_ "github.com/ommadawn46/the_go_programming_language-training/ch10/ex02/unarchive/tar"
	_ "github.com/ommadawn46/the_go_programming_language-training/ch10/ex02/unarchive/zip"
)

func extract(src, dest string) error {
	fmt.Println("Archive:", src)
	file, err := os.Open(src)
	if err != nil {
		return err
	}
	defer file.Close()

	files, kind, err := unarchive.Unarchive(file)
	if err != nil {
		return err
	}
	fmt.Println("Input format =", kind)

	for _, f := range files {
		path := filepath.Join(dest, f.Path())

		if f.FileInfo().IsDir() {
			fmt.Printf("    creating: %s\n", path)
			os.MkdirAll(path, os.ModePerm)
		} else {
			fmt.Printf("  extracting: %s\n", path)
			rc, err := f.Open()
			if err != nil {
				return err
			}
			defer rc.Close()

			outFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.FileInfo().Mode())
			if err != nil {
				return err
			}
			defer outFile.Close()

			_, err = io.Copy(outFile, rc)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func main() {
	var src, dest string
	switch len(os.Args) {
	case 2:
		src = os.Args[1]
		dest = "./"
	case 3:
		src = os.Args[1]
		dest = os.Args[2]
	default:
		fmt.Fprintf(os.Stderr, "Usage: %v srcfile destdir\n", os.Args[0])
		os.Exit(1)
	}

	if err := extract(src, dest); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
