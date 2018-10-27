package ftp

import (
	"fmt"
	"log"
	"net"
	"strings"

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
		session := Session{conn, nil, worker}
		go session.Run()
	}
}

type Session struct {
	ctrlconn  *pre.CtrlConn
	dataconns []*pre.DataConn
	worker    *app.Worker
}

func (s *Session) Run() {
	defer s.ctrlconn.Close()
	s.ctrlconn.SendResponce(220, "Ready")
	for {
		cmdName, arg, err := s.ctrlconn.RecvCommand()
		if err != nil {
			log.Printf("CLOSE %v", s.ctrlconn.RemoteAddr())
			return
		}
		log.Printf("[%v] %s %s", s.ctrlconn.RemoteAddr(), cmdName, arg)
		go func() {
			code, message := s.Dispatch(cmdName, arg)
			s.ctrlconn.SendResponce(code, message)
			if code == 221 {
				s.ctrlconn.Close() // QUIT
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
	if cmd.HasAttribute(app.CloseAllDataConns) {
		for _, dataconn := range s.dataconns {
			dataconn.Close()
		}
		s.dataconns = nil
	}
	var recvData []byte
	if cmd.HasAttribute(app.RecvDataBefore) {
		if err := s.worker.CheckReadyForTransfer(); err != nil {
			return 503, err.Error()
		}
		s.ctrlconn.SendResponce(150, fmt.Sprintf("Opening %s mode data connection", s.worker.TransType))
		dataconn, err := s.NewDataConn()
		if err != nil {
			return 425, err.Error()
		}
		defer dataconn.Close()
		if recvData, err = s.RecvData(dataconn); err != nil {
			return 426, err.Error()
		}
	}
	code, message, sendData := s.worker.Call(cmd, arg, recvData)
	if cmd.HasAttribute(app.SendDataAfter) {
		if err := s.worker.CheckReadyForTransfer(); err != nil {
			return 503, err.Error()
		}
		s.ctrlconn.SendResponce(150, fmt.Sprintf("Opening %s mode data connection", s.worker.TransType))
		dataconn, err := s.NewDataConn()
		if err != nil {
			return 425, err.Error()
		}
		defer dataconn.Close()
		if err = s.SendData(dataconn, sendData); err != nil {
			return 426, err.Error()
		}
	}
	return code, message
}

func (s *Session) NewDataConn() (*pre.DataConn, error) {
	if s.dataconns == nil {
		s.dataconns = []*pre.DataConn{}
	}
	var dataconn *pre.DataConn
	var err error
	if s.worker.PasvMode {
		dataconn, err = pre.AcceptNewDataConn(s.worker.DataListener)
		s.worker.PasvMode, s.worker.DataListener = false, nil
	} else {
		dataconn, err = pre.DialNewDataConn(s.worker.DataAddr)
		s.worker.DataAddr = ""
	}
	s.dataconns = append(s.dataconns, dataconn)
	return dataconn, err
}

func (s *Session) RecvData(dataconn *pre.DataConn) ([]byte, error) {
	var buf []byte
	var err error
	switch s.worker.TransType {
	case app.BINARY:
		buf, err = dataconn.ReadAll()
	case app.ASCII:
		buf, err = dataconn.ReadAllAsAscii()
	default:
		buf, err = nil, fmt.Errorf("Invalid transfer type")
	}
	if err != nil {
		return nil, fmt.Errorf("%v: Failed to transfer %s", err, s.worker.TransType)
	}
	return buf, nil
}

func (s *Session) SendData(dataconn *pre.DataConn, data []byte) error {
	var err error
	switch s.worker.TransType {
	case app.BINARY:
		err = dataconn.SendAll(data)
	case app.ASCII:
		err = dataconn.SendAllAsAscii(data)
	default:
		err = fmt.Errorf("Invalid transfer type")
	}
	if err != nil {
		return fmt.Errorf("%v: Failed to transfer %s", err, s.worker.TransType)
	}
	return nil
}
