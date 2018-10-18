package presentation

import (
	"fmt"
	"io"
	"net"
	"strings"
	"time"
)

type DataConn struct {
	net.Conn
}

var ASCII_REPLACER = strings.NewReplacer("\r\n", "\r\n", "\r", "\r\n", "\n", "\r\n")

func (c *DataConn) ReadAll() ([]byte, error) {
	var buf []byte
	for {
		tmp := make([]byte, 65536)
		n, err := c.Read(tmp)
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		buf = append(buf, tmp[:n]...)
	}
	return buf, nil
}

func (c *DataConn) SendAll(buf []byte) error {
	for len(buf) > 0 {
		n, err := c.Write(buf)
		if err != nil {
			return err
		}
		buf = buf[n:]
	}
	return nil
}

func (c *DataConn) ReadAllAsAscii() ([]byte, error) {
	buf, err := c.ReadAll()
	if err != nil {
		return nil, err
	}
	ascii := ASCII_REPLACER.Replace(string(buf))
	return []byte(ascii), nil
}

func (c *DataConn) SendAllAsAscii(buf []byte) error {
	ascii := ASCII_REPLACER.Replace(string(buf))
	buf = []byte(ascii)
	if err := c.SendAll(buf); err != nil {
		return err
	}
	return nil
}

func AcceptNewDataConn(listener *net.TCPListener) (*DataConn, error) {
	listener.SetDeadline(time.Now().Add(time.Second * 10))
	conn, err := listener.Accept()
	if err != nil {
		return nil, fmt.Errorf("Failed to establish connection")
	}
	return &DataConn{conn}, nil
}

func DialNewDataConn(addr string) (*DataConn, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("Failed to establish connection")
	}
	return &DataConn{conn}, nil
}
