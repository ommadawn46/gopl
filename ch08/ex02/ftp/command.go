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

func (s *Session) execCommand(cmd string, args []string) {
	if !s.loggedIn {
		if cmd != "USER" && cmd != "PASS" && cmd != "NOOP" && cmd != "QUIT" {
			s.conn.Sendline("530 Please login with USER and PASS")
			return
		}
	}
	switch cmd {
	case "USER":
		s.user(args)
	case "PASS":
		s.pass(args)
	case "NOOP":
		s.noop(args)
	case "QUIT":
		s.quit(args)
	case "PASV":
		s.pasv(args)
	case "PORT":
		s.port(args)
	case "TYPE":
		s.type_(args)
	case "PWD":
		s.pwd(args)
	case "CWD":
		s.cwd(args)
	case "MKD":
		s.mkd(args)
	case "LIST":
		s.list(args)
	case "RETR":
		s.retr(args)
	case "STOR":
		s.stor(args)
	case "DELE":
		s.dele(args)
	case "RMD":
		s.rmd(args)
	case "MODE":
		s.mode(args)
	case "STRU":
		s.stru(args)
	default:
		s.conn.Sendline("500 Unknown command")
	}
}

func (s *Session) user(args []string) {
	if len(args) < 1 {
		s.conn.Sendline("500 USER: command requires a parameter")
		return
	}

	s.username = args[0]
	s.conn.Sendline(fmt.Sprintf("331 Password required for %v", s.username))
}

func (s *Session) pass(args []string) {
	if s.username == "" {
		s.conn.Sendline("503 Login with USER first")
		return
	}
	if len(args) < 1 || s.username != _TESTUSER || args[0] != _TESTPASS {
		s.username = ""
		s.conn.Sendline("530 Login incorrect")
		return
	}

	s.loggedIn = true
	s.conn.Sendline(fmt.Sprintf("230 User %v logged in", s.username))
}

func (s *Session) noop(args []string) {
	s.conn.Sendline("200 NOOP command successful")
}

func (s *Session) quit(args []string) {
	s.conn.Sendline("221 Goodbye")
	s.conn.Close()
}

func (s *Session) pasv(args []string) {
	if len(args) != 0 {
		s.conn.Sendline("501 Invalid number of arguments")
		return
	}

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		s.conn.Sendline("500 Cannot listen to port")
		return
	}
	s.dataPort.listener = listener

	tcpAddr, _ := listener.Addr().(*net.TCPAddr)
	addr, port := tcpAddr.IP.String(), tcpAddr.Port
	port1, port2 := port/0x100, port%0x100

	s.dataPort.pasvMode = true
	s.conn.Sendline(fmt.Sprintf(
		"227 Entering Passive Mode (%s,%d,%d)",
		strings.Replace(addr, ".", ",", -1), port1, port2,
	))
}

func (s *Session) port(args []string) {
	if len(args) != 1 {
		s.conn.Sendline("501 Invalid number of arguments")
		return
	}

	addrPort := strings.Split(args[0], ",")
	if len(addrPort) != 6 {
		s.conn.Sendline("501 Illegal PORT command")
		return
	}

	addr := strings.Join(addrPort[:4], ".")
	port1, portErr1 := strconv.Atoi(addrPort[4])
	port2, portErr2 := strconv.Atoi(addrPort[5])
	if portErr1 != nil || portErr2 != nil {
		s.conn.Sendline("501 Illegal PORT command")
		return
	}
	port := port1*0x100 + port2
	s.dataPort.addr = fmt.Sprintf("%s:%d", addr, port)

	s.conn.Sendline("200 PORT command successful")
}

func (s *Session) type_(args []string) {
	if len(args) == 0 {
		s.conn.Sendline("500 TYPE: command requires a parameter")
		return
	}

	mode := strings.ToUpper(args[0])
	switch mode {
	case "A":
		s.binaryMode = false
	case "I":
		s.binaryMode = true
	default:
		s.conn.Sendline(fmt.Sprintf("500 'TYPE %s' not understood", mode))
		return
	}

	s.conn.Sendline(fmt.Sprintf("200 Type set to %s", mode))
}

func (s *Session) pwd(args []string) {
	s.conn.Sendline(fmt.Sprintf("257 \"/%s\" is the current directory", s.WorkDir))
}

func (s *Session) cwd(args []string) {
	if len(args) != 1 {
		s.conn.Sendline("501 Invalid number of arguments")
		return
	}

	path := s.JoinPath(args[0])
	fileInfo, err := os.Stat(path)
	if err != nil || !fileInfo.IsDir() {
		s.conn.Sendline(fmt.Sprintf("550 %s: No such file or directory", args[0]))
		return
	}
	s.WorkDir = strings.TrimLeft(path[len(s.RootDir):], "/\\")

	s.conn.Sendline("250 CWD command successful")
}

func (s *Session) mkd(args []string) {
	if len(args) != 1 {
		s.conn.Sendline("501 Invalid number of arguments")
		return
	}

	path := s.JoinPath(args[0])
	_, err := os.Stat(filepath.Dir(path))
	if os.IsNotExist(err) {
		s.conn.Sendline(fmt.Sprintf("550 %s: No such file or directory", args[0]))
		return
	}

	fileInfo, err := os.Stat(path)
	if err == nil && !fileInfo.IsDir() {
		s.conn.Sendline(fmt.Sprintf("550 %s: Not a directory", args[0]))
		return
	}

	if err := os.Mkdir(path, 0755); err != nil {
		s.conn.Sendline(fmt.Sprintf("550 %s: Failed to make a directory", args[0]))
		return
	}

	s.conn.Sendline(fmt.Sprintf("257 \"%s\" - Directory successfully created", args[0]))
}

func (s *Session) list(args []string) {
	if len(args) > 1 {
		s.conn.Sendline("501 Invalid number of arguments")
		return
	}

	var path string
	if len(args) == 0 {
		path = s.JoinPath("")
	} else {
		path = s.JoinPath(args[0])
	}

	lsOut, err := exec.Command("ls", "-la", path).Output()
	if err != nil {
		s.conn.Sendline("500 Failed to execute the list command")
		return
	}

	err = s.sendToDataPort(lsOut)
	if err != nil {
		s.conn.Sendline(err.Error())
	}
}

func (s *Session) retr(args []string) {
	if len(args) != 1 {
		s.conn.Sendline("501 Invalid number of arguments")
		return
	}

	path := s.JoinPath(args[0])
	fileInfo, err := os.Stat(path)
	if err != nil {
		s.conn.Sendline(fmt.Sprintf("550 %s: No such file or directory", args[0]))
		return
	}
	if fileInfo.IsDir() {
		s.conn.Sendline(fmt.Sprintf("550 %s: Not a regular file", args[0]))
		return
	}

	fileData, err := ioutil.ReadFile(path)
	if err != nil {
		s.conn.Sendline("500 Failed to read file")
		return
	}

	err = s.sendToDataPort(fileData)
	if err != nil {
		s.conn.Sendline(err.Error())
	}
}

func (s *Session) stor(args []string) {
	if len(args) != 1 {
		s.conn.Sendline("501 Invalid number of arguments")
		return
	}

	path := s.JoinPath(args[0])
	_, err := os.Stat(filepath.Dir(path))
	if os.IsNotExist(err) {
		s.conn.Sendline(fmt.Sprintf("550 %s: No such file or directory", args[0]))
		return
	}

	fileInfo, err := os.Stat(path)
	if err == nil && fileInfo.IsDir() {
		s.conn.Sendline(fmt.Sprintf("550 %s: Not a regular file", args[0]))
		return
	}

	recvData, err := s.recvFromDataPort()
	if err != nil {
		s.conn.Sendline(err.Error())
	}

	err = ioutil.WriteFile(path, recvData, 0644)
	if err != nil {
		s.conn.Sendline("500 Failed to execute the command")
		return
	}
}

func (s *Session) dele(args []string) {
	if len(args) != 1 {
		s.conn.Sendline("501 Invalid number of arguments")
		return
	}

	path := s.JoinPath(args[0])
	fileInfo, err := os.Stat(path)
	if err != nil {
		s.conn.Sendline(fmt.Sprintf("550 %s: No such file or directory", args[0]))
		return
	}
	if fileInfo.IsDir() {
		s.conn.Sendline(fmt.Sprintf("550 %s: Is a directory", args[0]))
		return
	}

	if err := os.Remove(path); err != nil {
		s.conn.Sendline(fmt.Sprintf("550 %s: Failed to delete a file", args[0]))
		return
	}

	s.conn.Sendline("250 DELE command successful")
}

func (s *Session) rmd(args []string) {
	if len(args) != 1 {
		s.conn.Sendline("501 Invalid number of arguments")
		return
	}

	path := s.JoinPath(args[0])
	fileInfo, err := os.Stat(path)
	if err != nil {
		s.conn.Sendline(fmt.Sprintf("550 %s: No such file or directory", args[0]))
		return
	}
	if !fileInfo.IsDir() {
		s.conn.Sendline(fmt.Sprintf("550 %s: Not a directory", args[0]))
		return
	}

	if err := os.Remove(path); err != nil {
		s.conn.Sendline(fmt.Sprintf("550 %s: Failed to delete a file", args[0]))
		return
	}

	s.conn.Sendline("250 RMD command successful")
}

func (s *Session) mode(args []string) {
	if len(args) != 1 {
		s.conn.Sendline("501 Invalid number of arguments")
		return
	}
	mode := strings.ToUpper(args[0])
	switch mode {
	case "S":
		s.conn.Sendline("200 Mode set to S")
	case "B", "C":
		s.conn.Sendline(fmt.Sprintf("504 'MODE %s' unsupported transfer mode", mode))
	default:
		s.conn.Sendline(fmt.Sprintf("501 'MODE %s' unrecognized transfer mode", mode))
	}
}

func (s *Session) stru(args []string) {
	if len(args) != 1 {
		s.conn.Sendline("501 Invalid number of arguments")
		return
	}
	mode := strings.ToUpper(args[0])
	switch mode {
	case "F":
		s.conn.Sendline("200 Structure set to F")
	case "R", "P":
		s.conn.Sendline(fmt.Sprintf("504 'MODE %s' unsupported structure type", mode))
	default:
		s.conn.Sendline(fmt.Sprintf("501 'MODE %s' unrecognized structure type", mode))
	}
}
