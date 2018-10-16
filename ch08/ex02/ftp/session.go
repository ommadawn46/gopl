package ftp

import (
	"log"
	"net"
	"path/filepath"
	"strings"
)

type Session struct {
	Status

	conn     Connection
	dataPort DataPort
}

type Status struct {
	Directory

	username     string
	loggedIn     bool
	transferType TransferType
}

type Directory struct {
	RootDir string
	WorkDir string
}

func (d *Directory) JoinPath(path string) string {
	newPath := filepath.Clean(path)
	if strings.HasPrefix("/", newPath) {
		newPath = filepath.Join(d.RootDir, newPath)
	} else {
		newPath = filepath.Join(d.RootDir, d.WorkDir, newPath)
	}
	if !strings.HasPrefix(newPath, d.RootDir) {
		newPath = d.RootDir
	}
	return newPath
}

type TransferType int

const (
	ASCII TransferType = iota
	BINARY
)

func (t TransferType) String() string {
	switch t {
	case ASCII:
		return "ASCII"
	case BINARY:
		return "BINARY"
	default:
		return "UNKNOWN"
	}
}

func (s *Session) Run() {
	defer s.conn.Close()
	s.conn.SendResponce(220, "Ready")
	for {
		recv, err := s.conn.Readline()
		if err != nil {
			return
		}
		recv = strings.TrimSpace(recv)
		log.Printf("[%v] %s", s.conn.RemoteAddr(), recv)

		var cmd, arg string
		spaceIdx := strings.Index(recv, " ")
		if spaceIdx != -1 {
			cmd, arg = strings.ToUpper(recv[:spaceIdx]), recv[spaceIdx+1:]
		} else {
			cmd, arg = strings.ToUpper(recv), ""
		}
		code, message := s.execCommand(cmd, arg)

		s.conn.SendResponce(code, message)
		if code == 221 {
			// QUIT
			break
		}
	}
}

func NewSession(conn net.Conn, rootDir string) *Session {
	status := Status{}
	status.transferType = ASCII
	status.RootDir = rootDir

	connection := Connection{conn}
	dataPort := DataPort{}

	return &Session{status, connection, dataPort}
}
