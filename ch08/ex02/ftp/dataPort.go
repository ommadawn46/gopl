package ftp

import (
	"fmt"
	"net"
	"time"
)

type DataPort struct {
	listener net.Listener
	addr     string
	pasvMode bool
}

func (d *DataPort) connect() (*Connection, error) {
	var conn net.Conn
	var err error

	if d.pasvMode {
		if d.listener == nil {
			return nil, fmt.Errorf("No PASV Listener")
		}
		d.listener.(*net.TCPListener).SetDeadline(time.Now().Add(time.Second * 10))
		conn, err = d.listener.Accept()
		d.listener, d.pasvMode = nil, false
	} else {
		if d.addr == "" {
			return nil, fmt.Errorf("Set port with PORT first")
		}
		conn, err = net.Dial("tcp", d.addr)
		d.addr = ""
	}

	if err != nil {
		return nil, fmt.Errorf("Failed to establish connection")
	}
	return &Connection{conn}, nil
}
