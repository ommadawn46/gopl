package ftp

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

type Directory struct {
	rootDir string
	workDir string
}

func (d *Directory) joinPath(path string) string {
	newPath := filepath.Clean(path)
	if strings.HasPrefix("/", newPath) {
		newPath = filepath.Join(d.rootDir, newPath)
	} else {
		newPath = filepath.Join(d.rootDir, d.workDir, newPath)
	}
	if !strings.HasPrefix(newPath, d.rootDir) {
		newPath = d.rootDir
	}
	return newPath
}

type TransferType int

const (
	_ASCII TransferType = iota
	_BINARY
)

func (t TransferType) String() string {
	switch t {
	case _ASCII:
		return "ASCII"
	case _BINARY:
		return "BINARY"
	default:
		return "UNKNOWN"
	}
}

type Worker struct {
	Directory
	username     string
	renameFrom   string
	loggedIn     bool
	transferType TransferType
	dataPort     DataPort
}

func (w *Worker) call(cmd Command, arg string) (int, string) {
	if w.loggedIn && cmd.hasAttribute(mustNotLogin) {
		return 503, "You are already logged in"
	}
	if !w.loggedIn && cmd.hasAttribute(needsLogin) {
		return 530, "Please login with USER and PASS"
	}
	if arg == "" && cmd.hasAttribute(needsArg) {
		return 501, "Command requires a parameter"
	}
	return cmd.exec(w, arg)
}

func (w *Worker) sendOsCommandOutput(cmd string, args ...string) (int, string) {
	dataConn, err := w.dataPort.connect()
	if err != nil {
		return 425, err.Error()
	}
	defer dataConn.Close()

	out, err := exec.Command(cmd, args...).Output()
	if err != nil {
		return 500, "Failed to execute the command"
	}

	switch w.transferType {
	case _BINARY:
		err = dataConn.sendAll(out)
	case _ASCII:
		err = dataConn.sendAllAsAscii(out)
	default:
		err = fmt.Errorf("Invalid transfer type")
	}
	if err != nil {
		return 426, fmt.Sprintf("%v: Failed to transfer %s", err, w.transferType)
	}

	return 226, "Transfer complete"
}
