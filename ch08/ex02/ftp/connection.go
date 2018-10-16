package ftp

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
)

type Connection struct {
	net.Conn
}

var _ASCII_REPLACER = strings.NewReplacer("\r\n", "\r\n", "\r", "\r\n", "\n", "\r\n")

func (c *Connection) Readline() (string, error) {
	return bufio.NewReader(c).ReadString('\n')
}

func (c *Connection) Sendline(str string) (int, error) {
	return io.WriteString(c, str+"\n")
}

func (c *Connection) SendResponce(code int, message string) (int, error) {
	return c.Sendline(fmt.Sprintf("%d %s", code, message))
}

func (c *Connection) ReadAll() ([]byte, error) {
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

func (c *Connection) SendAll(buf []byte) error {
	for len(buf) > 0 {
		n, err := c.Write(buf)
		if err != nil {
			return err
		}
		buf = buf[n:]
	}
	return nil
}

func (c *Connection) ReadAllAsAscii() ([]byte, error) {
	buf, err := c.ReadAll()
	if err != nil {
		return nil, err
	}
	ascii := _ASCII_REPLACER.Replace(string(buf))
	return []byte(ascii), nil
}

func (c *Connection) SendAllAsAscii(buf []byte) error {
	ascii := _ASCII_REPLACER.Replace(string(buf))
	buf = []byte(ascii)
	if err := c.SendAll(buf); err != nil {
		return err
	}
	return nil
}
