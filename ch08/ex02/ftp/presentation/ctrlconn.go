package presentation

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
)

type CtrlConn struct {
	*net.TCPConn
	scanner *bufio.Scanner
}

func (c *CtrlConn) RecvCommand() (string, string, error) {
	recv, err := c.Readline()
	if err != nil {
		return "", "", err
	}
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
	if c.scanner.Scan() {
		return c.scanner.Text(), nil
	}
	if err := c.scanner.Err(); err != nil {
		return "", err
	} else {
		return "", io.EOF
	}
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

func NewCtrlConn(conn *net.TCPConn) *CtrlConn {
	scanner := bufio.NewScanner(conn)
	return &CtrlConn{conn, scanner}
}
