package ftp

import (
	"crypto/sha512"
	"fmt"
	"log"
	"net"
)

type User struct {
	name string
	hash string
}

var _USERS = []User{
	{"test1", "0dd3e512642c97ca3f747f9a76e374fbda73f9292823c0313be9d78add7cdd8f72235af0c553dd26797e78e1854edee0ae002f8aba074b066dfce1af114e32f8"},
	{"test2", "925f43c3cfb956bbe3c6aa8023ba7ad5cfa21d104186fffc69e768e55940d9653b1cd36fba614fba2e1844f4436da20f83750c6ec1db356da154691bdd71a9b1"},
	{"test3", "8b3ec31f18ecbd708250148799f324c631e0c07f5af4e63fa3980588f27c78d048a2b0c812ffdd779d2b4596451dcdccb3d4cecdfee4a06e244404e3f7fe569b"},
}

func hash(s string) string {
	return fmt.Sprintf("%x", sha512.Sum512([]byte(s)))
}

func auth(login, pass string) bool {
	for _, user := range _USERS {
		if user.name == login && user.hash == hash(pass) {
			return true
		}
	}
	return false
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
		session := NewSession(conn, rootDir)
		go session.Run()
	}
	return nil
}
