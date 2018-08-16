package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	counts := make(map[string]int)
	file_names := make(map[string]map[string]struct{})

	files := os.Args[1:]
	if len(files) == 0 {
		countLines(os.Stdin, "stdin", counts, file_names)
	} else {
		for _, arg := range files {
			f, err := os.Open(arg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "dup2: %v\n", err)
				continue
			}
			countLines(f, arg, counts, file_names)
			f.Close()
		}
	}
	for line, n := range counts {
		if n > 1 {
			keys := make([]string, 0, len(file_names[line]))
			for k := range file_names[line] {
				keys = append(keys, k)
			}
			fmt.Printf("%d\t%s\t%s\n", n, line, strings.Join(keys, " "))
		}
	}
}

func countLines(f *os.File, file_name string, counts map[string]int, file_names map[string]map[string]struct{}) {
	input := bufio.NewScanner(f)
	for input.Scan() {
		counts[input.Text()]++
		if _, ok := file_names[input.Text()]; !ok {
			file_names[input.Text()] = make(map[string]struct{})
		}
		file_names[input.Text()][file_name] = struct{}{}
	}
}
