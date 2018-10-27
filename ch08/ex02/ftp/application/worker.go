package application

import (
	"fmt"
	"net"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/ommadawn46/the_go_programming_language-training/ch08/ex02/ftp/usermanager"
)

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

type Worker struct {
	TransType    TransferType
	PasvMode     bool
	DataListener *net.TCPListener
	DataAddr     string
	LoggedIn     bool
	cmdMutex     sync.Mutex

	rootDirPtr *string
	workDir    string
	renameFrom string

	userMgr   *usermanager.UserManager
	loginName string
}

func (w *Worker) Call(cmd Command, arg string, data []byte) (int, string, []byte) {
	w.cmdMutex.Lock()
	code, message, data := cmd.exec(w, arg, data)
	w.cmdMutex.Unlock()
	return code, message, data
}

func (w *Worker) CheckReadyForTransfer() error {
	if w.PasvMode && w.DataListener == nil {
		return fmt.Errorf("No PASV Listener")
	}
	if !w.PasvMode && w.DataAddr == "" {
		return fmt.Errorf("Set port with PORT first")
	}
	return nil
}

func (w *Worker) osCommandOutput(cmd string, args ...string) (int, string, []byte) {
	out, err := exec.Command(cmd, args...).Output()
	if err != nil {
		return 500, "Failed to execute the command", nil
	}
	return 226, "Transfer complete", out
}

func (w *Worker) joinPath(path string) string {
	newPath := filepath.Clean(path)
	if strings.HasPrefix("/", newPath) {
		newPath = filepath.Join(*w.rootDirPtr, newPath)
	} else {
		newPath = filepath.Join(*w.rootDirPtr, w.workDir, newPath)
	}
	if !strings.HasPrefix(newPath, *w.rootDirPtr) {
		newPath = *w.rootDirPtr
	}
	return newPath
}

func NewWorker(userMgr *usermanager.UserManager, rootDirPtr *string) *Worker {
	return &Worker{
		TransType:  ASCII,
		userMgr:    userMgr,
		rootDirPtr: rootDirPtr,
	}
}
