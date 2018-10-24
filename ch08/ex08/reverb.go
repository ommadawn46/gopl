package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

func echo(c net.Conn, shout string, delay time.Duration, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Fprintln(c, "\t", strings.ToUpper(shout))
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", shout)
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", strings.ToLower(shout))
}

func handleConn(conn net.Conn) {
	lineCh := make(chan string)

	go func(ch chan string) {
		input := bufio.NewScanner(conn)
		for input.Scan() {
			ch <- input.Text()
		}
		close(ch)
	}(lineCh)

	timeOut := false
	wg := sync.WaitGroup{}
	for !timeOut {
		ticker := time.NewTicker(10 * time.Second)
		select {
		case <-ticker.C:
			timeOut = true
		case text := <-lineCh:
			wg.Add(1)
			go echo(conn, text, 1*time.Second, &wg)
		}
		ticker.Stop()
	}
	wg.Wait()
	conn.Close()
}

func main() {
	l, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}
