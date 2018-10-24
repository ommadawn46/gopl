package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

type client struct {
	ch   chan<- string
	name string
}

var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string)
)

func broadcaster() {
	clients := make(map[client]bool)
	for {
		select {
		case msg := <-messages:
			for cli := range clients {
				select {
				case cli.ch <- msg:
				default:
				}
			}

		case cli := <-entering:
			if len(clients) > 0 {
				var clientNames []string
				for oCli := range clients {
					clientNames = append(clientNames, oCli.name)
				}
				cli.ch <- "Users in chat: " + strings.Join(clientNames, ", ")
			}
			clients[cli] = true

		case cli := <-leaving:
			delete(clients, cli)
			close(cli.ch)
		}
	}
}

func handleConn(conn net.Conn) {
	ch := make(chan string)
	go clientWriter(conn, ch)

	inputs := make(chan string)
	go func() {
		input := bufio.NewScanner(conn)
		for input.Scan() {
			inputs <- input.Text()
		}
		close(inputs)
	}()

	ch <- "What is your name?: "
	who := <-inputs
	cli := client{ch, who}
	ch <- "Welcome, " + who

	messages <- who + " has arrived"
	entering <- cli

	finish := false
	ticker := time.NewTicker(5 * time.Minute)
	for !finish {
		select {
		case text, ok := <-inputs:
			if ok {
				messages <- who + ": " + text
			} else {
				finish = true
			}

		case <-ticker.C:
			finish = true
		}
	}
	ticker.Stop()

	leaving <- cli
	messages <- who + " has left"
	conn.Close()
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg)
	}
}

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	go broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}
