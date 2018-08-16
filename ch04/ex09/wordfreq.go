package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	counter := make(map[string]int)

	var input = bufio.NewScanner(os.Stdin)
	input.Split(bufio.ScanWords)
	for input.Scan() {
		w := input.Text()
		counter[w]++
	}
	fmt.Printf("\nword                \tcount\n")
	for w, n := range counter {
		fmt.Printf("%-20s\t%d\n", w, n)
	}
}
