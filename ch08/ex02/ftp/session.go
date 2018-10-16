package ftp

import (
	"fmt"
	"log"
	"net"
	"strings"
)

type Session struct {
	worker Worker
	conn   Connection
}

func (s *Session) run() {
	defer s.conn.Close()
	s.conn.sendResponce(220, "Ready")
	for {
		recv, err := s.conn.readline()
		if err != nil {
			return
		}
		recv = strings.TrimSpace(recv)
		log.Printf("[%v] %s", s.conn.RemoteAddr(), recv)

		go func() {
			var cmdName, arg string
			spaceIdx := strings.Index(recv, " ")
			if spaceIdx != -1 {
				cmdName, arg = recv[:spaceIdx], recv[spaceIdx+1:]
			} else {
				cmdName, arg = recv, ""
			}
			cmd, ok := _COMMANDS[strings.ToUpper(cmdName)]
			if !ok {
				s.conn.sendResponce(500, "Unknown command")
				return
			}
			if cmd.hasAttribute(useDataPort) {
				s.conn.sendResponce(150, fmt.Sprintf("Opening %s mode data connection", s.worker.transferType))
			}
			code, message := s.worker.call(cmd, arg)
			s.conn.sendResponce(code, message)
			if code == 221 {
				s.conn.Close() // QUIT
			}
		}()
	}
}

func newSession(conn net.Conn, rootDir string) *Session {
	directory := Directory{}
	directory.rootDir = rootDir
	worker := Worker{
		directory,
		"",
		"",
		false,
		_ASCII,
		DataPort{},
	}
	connection := Connection{conn}
	return &Session{worker, connection}
}
