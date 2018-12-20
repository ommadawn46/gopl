package bzip

import (
	"io"
	"os/exec"
)

type writer struct {
	cmd   *exec.Cmd
	stdin io.WriteCloser
}

var BZIP2PATH = "/usr/bin/bzip2"

func NewWriter(out io.Writer) io.WriteCloser {
	cmd := exec.Command(BZIP2PATH)
	cmd.Stdout = out
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil
	}
	if err = cmd.Start(); err != nil {
		return nil
	}
	return &writer{cmd, stdin}
}

func (w *writer) Write(data []byte) (int, error) {
	return w.stdin.Write(data)
}

func (w *writer) Close() error {
	closeErr := w.stdin.Close()
	waitErr := w.cmd.Wait()
	if closeErr != nil {
		return closeErr
	}
	if waitErr != nil {
		return waitErr
	}
	return nil
}
