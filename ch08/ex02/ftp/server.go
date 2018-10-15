package ftp

import (
	"log"
	"net"
)

func ListenAndServe(addr string, rootDir string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		ftpSession := NewSession(conn, rootDir)
		go ftpSession.Run()
	}
	return nil
}
