package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/ommadawn46/the_go_programming_language-training/ch05/ex13/links"
)

var depthOpt = flag.Int("depth", 3, "crawl depth")

type Target struct {
	url   string
	depth int
}

func crawl(target Target) []string {
	fmt.Println(target.depth, target.url)
	if target.depth >= *depthOpt {
		return nil
	}
	list, err := links.Extract(target.url)
	if err != nil {
		log.Print(err)
	}
	return list
}

func main() {
	flag.Parse()
	worklist := make(chan []Target)
	unseenTargets := make(chan Target)

	go func() {
		var targets []Target
		for _, url := range flag.Args() {
			targets = append(targets, Target{url, 0})
		}
		worklist <- targets
	}()

	for i := 0; i < 20; i++ {
		go func() {
			for target := range unseenTargets {
				foundLinks := crawl(target)
				var foundTargets []Target
				for _, url := range foundLinks {
					foundTargets = append(foundTargets, Target{url, target.depth + 1})
				}
				go func() { worklist <- foundTargets }()
			}
		}()
	}

	seen := make(map[string]bool)
	for list := range worklist {
		for _, target := range list {
			if !seen[target.url] {
				seen[target.url] = true
				unseenTargets <- target
			}
		}
	}
}
