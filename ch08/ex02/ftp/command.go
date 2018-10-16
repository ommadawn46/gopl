package ftp

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func (s *Session) execCommand(cmd string, arg string) (int, string) {
	if !s.loggedIn {
		if cmd != "NOOP" && cmd != "PASS" && cmd != "QUIT" && cmd != "USER" {
			return 530, "Please login with USER and PASS"
		}
	}
	switch cmd {
	case "CWD":
		return s.cwd(arg)
	case "DELE":
		return s.dele(arg)
	case "LIST":
		return s.list(arg)
	case "MKD":
		return s.mkd(arg)
	case "MODE":
		return s.mode(arg)
	case "NOOP":
		return s.noop(arg)
	case "PASS":
		return s.pass(arg)
	case "PASV":
		return s.pasv(arg)
	case "PORT":
		return s.port(arg)
	case "PWD":
		return s.pwd(arg)
	case "QUIT":
		return s.quit(arg)
	case "RETR":
		return s.retr(arg)
	case "RMD":
		return s.rmd(arg)
	case "STOR":
		return s.stor(arg)
	case "STRU":
		return s.stru(arg)
	case "TYPE":
		return s.type_(arg)
	case "USER":
		return s.user(arg)
	default:
		return 500, "Unknown command"
	}
}

func (s *Session) cwd(arg string) (int, string) {
	if arg == "" {
		return 501, "CWD command requires a parameter"
	}

	path := s.JoinPath(arg)
	fileInfo, err := os.Stat(path)
	if err != nil || !fileInfo.IsDir() {
		return 550, fmt.Sprintf("%s: No such file or directory", arg)
	}
	s.WorkDir = strings.TrimLeft(path[len(s.RootDir):], "/\\")

	return 250, "CWD command successful"
}

func (s *Session) dele(arg string) (int, string) {
	if arg == "" {
		return 501, "DELE command requires a parameter"
	}

	path := s.JoinPath(arg)
	fileInfo, err := os.Stat(path)
	if err != nil {
		return 550, fmt.Sprintf("%s: No such file or directory", arg)
	}
	if fileInfo.IsDir() {
		return 550, fmt.Sprintf("%s: Is a directory", arg)
	}

	if err := os.Remove(path); err != nil {
		return 550, fmt.Sprintf("%s: Failed to delete a file", arg)
	}

	return 250, "DELE command successful"
}

func (s *Session) list(arg string) (int, string) {
	if arg != "" {
		return 501, "Invalid number of arguments"
	}

	path := s.JoinPath(arg)
	lsOut, err := exec.Command("ls", "-la", path).Output()
	if err != nil {
		return 500, "Failed to execute the list command"
	}

	s.conn.SendResponce(150, fmt.Sprintf("Opening %s mode data connection", s.transferType))
	dataConn, err := s.dataPort.Connect()
	if err != nil {
		return 425, fmt.Sprintf("%v", err)
	}
	defer dataConn.Close()

	switch s.transferType {
	case BINARY:
		err = dataConn.SendAll(lsOut)
	case ASCII:
		err = dataConn.SendAllAsAscii(lsOut)
	default:
		err = fmt.Errorf("Invalid transfer type")
	}
	if err != nil {
		return 426, fmt.Sprintf("%v: Failed to transfer %s", err, s.transferType)
	}

	return 226, "Transfer complete"
}

func (s *Session) mkd(arg string) (int, string) {
	if arg == "" {
		return 501, "MKD command requires a parameter"
	}

	path := s.JoinPath(arg)
	_, err := os.Stat(filepath.Dir(path))
	if os.IsNotExist(err) {
		return 550, fmt.Sprintf("%s: No such file or directory", arg)
	}

	fileInfo, err := os.Stat(path)
	if err == nil && !fileInfo.IsDir() {
		return 550, fmt.Sprintf("%s: Not a directory", arg)
	}

	if err := os.Mkdir(path, 0755); err != nil {
		return 550, fmt.Sprintf("%s: Failed to make a directory", arg)
	}

	return 257, fmt.Sprintf("\"%s\" - Directory successfully created", arg)
}

func (s *Session) mode(arg string) (int, string) {
	if arg == "" {
		return 501, "MODE command requires a parameter"
	}
	mode := strings.ToUpper(arg)
	switch mode {
	case "S":
		return 200, "Mode set to S"
	case "B", "C":
		return 504, fmt.Sprintf("'MODE %s' unsupported transfer mode", mode)
	default:
		return 501, fmt.Sprintf("'MODE %s' unrecognized transfer mode", mode)
	}
}

func (s *Session) noop(arg string) (int, string) {
	return 200, "NOOP command successful"
}

func (s *Session) pass(arg string) (int, string) {
	if s.username == "" {
		return 503, "Login with USER first"
	}
	if !auth(s.username, arg) {
		s.username = ""
		return 530, "Login incorrect"
	}

	s.loggedIn = true
	return 230, fmt.Sprintf("User %v logged in", s.username)
}

func (s *Session) pasv(arg string) (int, string) {
	if arg != "" {
		return 501, "Invalid number of arguments"
	}

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 500, "Cannot listen to port"
	}
	s.dataPort.listener = listener

	tcpAddr, _ := listener.Addr().(*net.TCPAddr)
	addr, port := tcpAddr.IP.String(), tcpAddr.Port
	port1, port2 := port/0x100, port%0x100

	s.dataPort.pasvMode = true
	return 227, fmt.Sprintf(
		"Entering Passive Mode (%s,%d,%d)",
		strings.Replace(addr, ".", ",", -1), port1, port2,
	)
}

func (s *Session) port(arg string) (int, string) {
	if arg == "" {
		return 501, "PORT command requires a parameter"
	}

	addrPort := strings.Split(arg, ",")
	if len(addrPort) != 6 {
		return 501, "Illegal PORT command"
	}

	addr := strings.Join(addrPort[:4], ".")
	port1, portErr1 := strconv.Atoi(addrPort[4])
	port2, portErr2 := strconv.Atoi(addrPort[5])
	if portErr1 != nil || portErr2 != nil {
		return 501, "Illegal PORT command"
	}
	port := port1*0x100 + port2
	s.dataPort.addr = fmt.Sprintf("%s:%d", addr, port)

	return 200, "PORT command successful"
}

func (s *Session) pwd(arg string) (int, string) {
	return 257, fmt.Sprintf("\"/%s\" is the current directory", s.WorkDir)
}

func (s *Session) quit(arg string) (int, string) {
	return 221, "Goodbye"
}

func (s *Session) retr(arg string) (int, string) {
	if arg == "" {
		return 501, "RETR command requires a parameter"
	}

	path := s.JoinPath(arg)
	fileInfo, err := os.Stat(path)
	if err != nil {
		return 550, fmt.Sprintf("%s: No such file or directory", arg)
	}
	if fileInfo.IsDir() {
		return 550, fmt.Sprintf("%s: Not a regular file", arg)
	}

	fileData, err := ioutil.ReadFile(path)
	if err != nil {
		return 500, "Failed to read file"
	}

	s.conn.SendResponce(150, fmt.Sprintf("Opening %s mode data connection", s.transferType))
	dataConn, err := s.dataPort.Connect()
	if err != nil {
		return 425, fmt.Sprintf("%v", err)
	}
	defer dataConn.Close()

	switch s.transferType {
	case BINARY:
		err = dataConn.SendAll(fileData)
	case ASCII:
		err = dataConn.SendAllAsAscii(fileData)
	default:
		err = fmt.Errorf("Invalid transfer type")
	}
	if err != nil {
		return 426, fmt.Sprintf("%v: Failed to transfer %s", err, s.transferType)
	}
	return 226, "Transfer complete"
}

func (s *Session) rmd(arg string) (int, string) {
	if arg == "" {
		return 501, "RMD command requires a parameter"
	}

	path := s.JoinPath(arg)
	fileInfo, err := os.Stat(path)
	if err != nil {
		return 550, fmt.Sprintf("%s: No such file or directory", arg)
	}
	if !fileInfo.IsDir() {
		return 550, fmt.Sprintf("%s: Not a directory", arg)
	}

	if err := os.Remove(path); err != nil {
		return 550, fmt.Sprintf("%s: Failed to delete a file", arg)
	}

	return 250, "RMD command successful"
}

func (s *Session) stor(arg string) (int, string) {
	if arg == "" {
		return 501, "STOR command requires a parameter"
	}

	path := s.JoinPath(arg)
	_, err := os.Stat(filepath.Dir(path))
	if os.IsNotExist(err) {
		return 550, fmt.Sprintf("%s: No such file or directory", arg)
	}

	fileInfo, err := os.Stat(path)
	if err == nil && fileInfo.IsDir() {
		return 550, fmt.Sprintf("%s: Not a regular file", arg)
	}

	s.conn.SendResponce(150, fmt.Sprintf("Opening %s mode data connection", s.transferType))
	dataConn, err := s.dataPort.Connect()
	if err != nil {
		return 425, fmt.Sprintf("%v", err)
	}
	defer dataConn.Close()

	var buf []byte
	switch s.transferType {
	case BINARY:
		buf, err = dataConn.ReadAll()
	case ASCII:
		buf, err = dataConn.ReadAllAsAscii()
	default:
		buf, err = nil, fmt.Errorf("Invalid transfer type")
	}
	if err != nil {
		return 426, fmt.Sprintf("%v: Failed to transfer %s", err, s.transferType)
	}

	err = ioutil.WriteFile(path, buf, 0644)
	if err != nil {
		return 500, "Failed to write a file"
	}
	return 226, "Transfer complete"
}

func (s *Session) stru(arg string) (int, string) {
	if arg == "" {
		return 501, "STRU command requires a parameter"
	}
	mode := strings.ToUpper(arg)
	switch mode {
	case "F":
		return 200, "Structure set to F"
	case "R", "P":
		return 504, fmt.Sprintf("'MODE %s' unsupported structure type", mode)
	default:
		return 501, fmt.Sprintf("'MODE %s' unrecognized structure type", mode)
	}
}

func (s *Session) type_(arg string) (int, string) {
	if arg == "" {
		return 501, "TYPE command requires a parameter"
	}

	mode := strings.ToUpper(arg)
	switch mode {
	case "A":
		s.transferType = ASCII
	case "I":
		s.transferType = BINARY
	default:
		return 500, fmt.Sprintf("'TYPE %s' not understood", mode)
	}

	return 200, fmt.Sprintf("Type set to %s", mode)
}

func (s *Session) user(arg string) (int, string) {
	if arg == "" {
		return 501, "USER command requires a parameter"
	}

	s.username = arg
	return 331, fmt.Sprintf("Password required for %v", s.username)
}
