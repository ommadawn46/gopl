package counter

import (
	"bufio"
	"bytes"
)

type LineCounter int

func (c *LineCounter) Write(p []byte) (n int, err error) {
	s := bufio.NewScanner(bytes.NewReader(p))
	s.Split(bufio.ScanLines)
	for s.Scan() {
		*c++
	}
	return len(p), nil
}

type WordCounter int

func (c *WordCounter) Write(p []byte) (n int, err error) {
	s := bufio.NewScanner(bytes.NewReader(p))
	s.Split(bufio.ScanWords)
	for s.Scan() {
		*c++
	}
	return len(p), nil
}
