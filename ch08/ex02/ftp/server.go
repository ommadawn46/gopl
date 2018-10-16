package ftp

import (
	"encoding/csv"
	"io"
	"log"
	"net"
	"os"
)

var passwdPath = "./passwd"

type User struct {
	name string
	salt string
	hash string
}

func auth(login, pass string) bool {
	passwdFile, err := os.Open(passwdPath)
	if err != nil {
		log.Println(err)
		return false
	}
	defer passwdFile.Close()

	reader := csv.NewReader(passwdFile)
	reader.Comma = ':'
	for {
		record, err := reader.Read()
		if err == io.EOF {
			return false
		}
		if err != nil {
			log.Println(err)
			return false
		}
		user := User{record[0], record[1], record[2]}
		if user.name == login && user.hash == hash(user.salt+pass) {
			return true
		}
	}
}

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
		session := newSession(conn, rootDir)
		go session.run()
	}
	return nil
}
