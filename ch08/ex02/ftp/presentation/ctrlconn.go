package presentation

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

type CtrlConn struct {
	net.Conn
}

func (c *CtrlConn) RecvCommand() (string, string, error) {
	recv, err := c.Readline()
	if err != nil {
		return "", "", nil
	}
	recv = strings.TrimSpace(recv)
	log.Printf("[%v] %s", c.RemoteAddr(), recv)

	var cmdName, arg string
	spaceIdx := strings.Index(recv, " ")
	if spaceIdx != -1 {
		cmdName, arg = recv[:spaceIdx], recv[spaceIdx+1:]
	} else {
		cmdName, arg = recv, ""
	}
	return cmdName, arg, nil
}

func (c *CtrlConn) Readline() (string, error) {
	return bufio.NewReader(c).ReadString('\n')
}

func (c *CtrlConn) Sendline(str string) (int, error) {
	return io.WriteString(c, str+"\n")
}

func (c *CtrlConn) SendResponce(code int, message string) (int, error) {
	lines := strings.Split(message, "\n")
	length := len(lines)
	if length > 1 {
		lines[0] = fmt.Sprintf("%d-%s", code, lines[0])
	}
	lines[length-1] = fmt.Sprintf("%d %s", code, lines[length-1])
	return c.Sendline(strings.Join(lines, "\n"))
}

func NewCtrlConn(conn net.Conn) *CtrlConn {
	return &CtrlConn{conn}
}
