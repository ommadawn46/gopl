package main

import (
	"flag"
	"fmt"
	"time"
)

var nOpt = flag.Int("n", 10000, "number of goroutines")

func main() {
	flag.Parse()
	n := *nOpt

	firstCh := make(chan struct{})
	var ch1, ch2 chan struct{}

	fmt.Println("Creating goroutines start")
	start := time.Now()
	ch1, ch2 = nil, firstCh
	for i := 0; i < n; i++ {
		ch1, ch2 = ch2, make(chan struct{})
		go func(ch1 chan struct{}, ch2 chan struct{}) {
			ch1 <- <-ch2
		}(ch1, ch2)
	}
	lastCh := ch2
	end := time.Now()
	fmt.Printf("Created %d Goroutine, elapsed %.2fs.\n", n, end.Sub(start).Seconds())

	fmt.Println("Passing pipeline start")
	start = time.Now()
	lastCh <- struct{}{}
	<-firstCh
	end = time.Now()
	fmt.Printf("Passed through %d pipes, elapsed %.2fs.\n", n, end.Sub(start).Seconds())
}
