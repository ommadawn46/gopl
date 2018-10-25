package main

import (
	"fmt"
	"time"
)

func main() {
	ping, pong := 0, 0
	ch1, ch2 := make(chan struct{}), make(chan struct{})
	go func() {
		for s := range ch1 {
			ping++
			ch2 <- s
		}
	}()
	go func() {
		for s := range ch2 {
			pong++
			ch1 <- s
		}
	}()
	ch1 <- struct{}{}
	time.Sleep(time.Second * 1)
	fmt.Printf("Passed %d times between 1 second.", ping+pong)
}
