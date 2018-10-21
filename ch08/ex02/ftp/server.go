package ftp

import (
	"fmt"
	"net"
	"strings"
	"log"

	app "github.com/ommadawn46/the_go_programming_language-training/ch08/ex02/ftp/application"
	pre "github.com/ommadawn46/the_go_programming_language-training/ch08/ex02/ftp/presentation"
	"github.com/ommadawn46/the_go_programming_language-training/ch08/ex02/ftp/usermanager"
	"github.com/ommadawn46/the_go_programming_language-training/ch08/ex02/ftp/util"
)

func ListenAndServe(addr, rootDir, passwdPath string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	log.Printf("LISTENING %v", listener.Addr())
	return Serve(listener.(*net.TCPListener), rootDir, passwdPath)
}

func Serve(listener *net.TCPListener, rootDir, passwdPath string) error {
	if !util.ExistsPath(rootDir) {
		return fmt.Errorf("directory does not exist: %s", rootDir)
	}
	if !util.IsDirectory(rootDir) {
		return fmt.Errorf("not a directory: %s", rootDir)
	}
	userMgr, err := usermanager.NewUserManager(passwdPath)
	if err != nil {
		return err
	}
	for {
		c, err := listener.Accept()
		if err != nil {
			return err
		}
		log.Printf("ACCEPT %v", c.RemoteAddr())
		conn := pre.NewCtrlConn(c.(*net.TCPConn))
		worker := app.NewWorker(userMgr, &rootDir)
		session := Session{conn, worker}
		go session.Run()
	}
}

type Session struct {
	conn   *pre.CtrlConn
	worker *app.Worker
}

func (s *Session) Run() {
	defer s.conn.Close()
	s.conn.SendResponce(220, "Ready")
	for {
		cmdName, arg, err := s.conn.RecvCommand()
		if err != nil {
			log.Printf("CLOSE %v", s.conn.RemoteAddr())
			return
		}
		log.Printf("[%v] %s %s", s.conn.RemoteAddr(), cmdName, arg)
		go func() {
			code, message := s.Dispatch(cmdName, arg)
			s.conn.SendResponce(code, message)
			if code == 221 {
				s.conn.Close() // QUIT
			}
		}()
	}
}

func (s *Session) Dispatch(cmdName, arg string) (int, string) {
	cmd, ok := app.COMMANDS[strings.ToUpper(cmdName)]
	if !ok {
		return 500, "Unknown command"
	}
	if s.worker.LoggedIn && cmd.HasAttribute(app.MustNotLogin) {
		return 503, "You are already logged in"
	}
	if !s.worker.LoggedIn && cmd.HasAttribute(app.NeedsLogin) {
		return 530, "Please login with USER and PASS"
	}
	if arg == "" && cmd.HasAttribute(app.NeedsArg) {
		return 501, "Command requires a parameter"
	}
	var recvData []byte
	if cmd.HasAttribute(app.RecvDataBefore) {
		if err := s.worker.CheckReadyForTransfer(); err != nil {
			return 503, err.Error()
		}
		s.conn.SendResponce(150, fmt.Sprintf("Opening %s mode data connection", s.worker.TransType))
		dataConn, err := s.NewDataConn()
		if err != nil {
			return 425, err.Error()
		}
		if recvData, err = s.RecvData(dataConn); err != nil {
			return 426, err.Error()
		}
	}
	code, message, sendData := s.worker.Call(cmd, arg, recvData)
	if cmd.HasAttribute(app.SendDataAfter) {
		if err := s.worker.CheckReadyForTransfer(); err != nil {
			return 503, err.Error()
		}
		s.conn.SendResponce(150, fmt.Sprintf("Opening %s mode data connection", s.worker.TransType))
		dataConn, err := s.NewDataConn()
		if err != nil {
			return 425, err.Error()
		}
		if err = s.SendData(dataConn, sendData); err != nil {
			return 426, err.Error()
		}
	}
	return code, message
}

func (s *Session) NewDataConn() (*pre.DataConn, error) {
	var dataconn *pre.DataConn
	var err error
	if s.worker.PasvMode {
		dataconn, err = pre.AcceptNewDataConn(s.worker.DataListener)
		s.worker.PasvMode, s.worker.DataListener = false, nil
	} else {
		dataconn, err = pre.DialNewDataConn(s.worker.DataAddr)
		s.worker.DataAddr = ""
	}
	return dataconn, err
}

func (s *Session) RecvData(dataConn *pre.DataConn) ([]byte, error) {
	defer dataConn.Close()

	var buf []byte
	var err error
	switch s.worker.TransType {
	case app.BINARY:
		buf, err = dataConn.ReadAll()
	case app.ASCII:
		buf, err = dataConn.ReadAllAsAscii()
	default:
		buf, err = nil, fmt.Errorf("Invalid transfer type")
	}
	if err != nil {
		return nil, fmt.Errorf("%v: Failed to transfer %s", err, s.worker.TransType)
	}
	return buf, nil
}

func (s *Session) SendData(dataConn *pre.DataConn, data []byte) error {
	defer dataConn.Close()

	var err error
	switch s.worker.TransType {
	case app.BINARY:
		err = dataConn.SendAll(data)
	case app.ASCII:
		err = dataConn.SendAllAsAscii(data)
	default:
		err = fmt.Errorf("Invalid transfer type")
	}
	if err != nil {
		return fmt.Errorf("%v: Failed to transfer %s", err, s.worker.TransType)
	}
	return nil
}
