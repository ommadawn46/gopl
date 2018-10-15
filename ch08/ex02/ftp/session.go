package ftp

import (
	"fmt"
	"net"
	"strings"
)

var _TESTUSER = "test"
var _TESTPASS = "qwerty"

type Status struct {
	username   string
	loggedIn   bool
	binaryMode bool
}

type Session struct {
	Status
	Directory

	conn     Connection
	dataPort DataPort
}

func (s *Session) Run() {
	defer s.conn.Close()
	s.conn.Sendline("220 Ready")
	for {
		recv, err := s.conn.Readline()
		if err != nil {
			return
		}
		fmt.Printf("%v: %s", s.conn.RemoteAddr(), recv)

		args := strings.Fields(strings.TrimSpace(recv))
		cmd, cmdArgs := strings.ToUpper(args[0]), args[1:]
		s.execCommand(cmd, cmdArgs)
	}
}

func (s *Session) sendToDataPort(data []byte) error {
	dataConn, err := s.dataPort.Connect()
	if err != nil {
		return fmt.Errorf("425 %v", err)
	}
	defer dataConn.Close()

	if s.binaryMode {
		s.conn.Sendline(fmt.Sprintf("150 Opening BINARY mode data connection"))
		if err = dataConn.SendAll(data); err != nil {
			return fmt.Errorf("426 %v: Failed to transfer BINARY", err)
		}
	} else {
		s.conn.Sendline(fmt.Sprintf("150 Opening ASCII mode data connection"))
		if err = dataConn.SendAllAsAscii(data); err != nil {
			return fmt.Errorf("426 %v: Failed to transfer ASCII", err)
		}
	}
	s.conn.Sendline("226 Transfer complete")
	return nil
}

func (s *Session) recvFromDataPort() ([]byte, error) {
	dataConn, err := s.dataPort.Connect()
	if err != nil {
		return nil, err
	}
	defer dataConn.Close()

	var buf []byte
	if s.binaryMode {
		s.conn.Sendline(fmt.Sprintf("150 Opening BINARY mode data connection"))
		buf, err = dataConn.ReadAll()
		if err != nil {
			return nil, fmt.Errorf("426 %v: Failed to transfer BINARY", err)
		}
	} else {
		s.conn.Sendline(fmt.Sprintf("150 Opening ASCII mode data connection"))
		buf, err = dataConn.ReadAllAsAscii()
		if err != nil {
			return nil, fmt.Errorf("426 %v: Failed to transfer ASCII", err)
		}
	}

	s.conn.Sendline("226 Transfer complete")
	return buf, nil
}

func NewSession(conn net.Conn, RootDir string) *Session {
	status := Status{}
	directory := Directory{RootDir: RootDir}
	connection := Connection{conn}
	dataPort := DataPort{}
	return &Session{status, directory, connection, dataPort}
}
