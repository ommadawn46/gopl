package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/ommadawn46/the_go_programming_language-training/ch08/ex10/links"
)

var cancelCh = make(chan struct{}, 1)

func cancelAll() {
	close(cancelCh)
}

func crawl(url string) []string {
	fmt.Println(url)

	list, err := links.Extract(url, cancelCh)
	if err != nil {
		log.Print(err)
	}
	return list
}

func main() {
	go func() {
		for {
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				switch scanner.Text() {
				case "C", "c":
					cancelAll()
				}
			}
		}
	}()

	worklist := make(chan []string)
	unseenLinks := make(chan string)
	go func() { worklist <- os.Args[1:] }()

	for i := 0; i < 20; i++ {
		go func() {
			for link := range unseenLinks {
				foundLinks := crawl(link)
				go func() { worklist <- foundLinks }()
			}
		}()
	}

	seen := make(map[string]bool)
	for list := range worklist {
		for _, link := range list {
			if !seen[link] {
				seen[link] = true
				unseenLinks <- link
			}
		}
	}
}
