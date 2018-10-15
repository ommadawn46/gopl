package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

var padding = 2
var time_len = len("--:--:--")

type Clock struct {
	zone string
	host string
	time string
}

func (c Clock) fmtStr() string {
	if zone_len := len(c.zone); zone_len > time_len {
		return "%" + strconv.Itoa(zone_len+padding) + "v"
	} else {
		return "%" + strconv.Itoa(time_len+padding) + "v"
	}
}

func (c Clock) Now() string {
	return fmt.Sprintf(c.fmtStr(), c.time)
}

func (c Clock) Zone() string {
	return fmt.Sprintf(c.fmtStr(), c.zone)
}

func (c *Clock) Connect() {
	conn, err := net.Dial("tcp", c.host)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		conn.Close()
	}()
	reader := bufio.NewReader(conn)
	for {
		now, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		c.time = strings.TrimSpace(now)
	}
}

func main() {
	if len(os.Args) <= 1 {
		fmt.Printf("Usage: %v TimeZone=Host\n", os.Args[0])
		os.Exit(1)
	}

	var clocks []*Clock
	for _, arg := range os.Args[1:] {
		equal_idx := strings.Index(arg, "=")
		clocks = append(clocks, &Clock{
			zone: arg[:equal_idx],
			host: arg[equal_idx+1:],
			time: "--:--:--",
		})
	}

	for _, clock := range clocks {
		fmt.Print(clock.Zone())
		go clock.Connect()
	}
	fmt.Print("\n")

	for {
		for _, clock := range clocks {
			fmt.Print(clock.Now())
		}
		fmt.Print("\r")
		time.Sleep(100 * time.Millisecond)
	}
}
