package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"

	"github.com/ommadawn46/the_go_programming_language-training/ch08/ex02/ftp"
	"github.com/ommadawn46/the_go_programming_language-training/ch08/ex02/ftp/usermanager"
)

var portFlag *int = flag.Int("port", 21, "listening port")
var rootFlag *string = flag.String("root", "./", "root directory")
var passwdFlag *string = flag.String("passwd", "./passwd", "passwd path")

var addUserFlag *bool = flag.Bool("adduser", false, "add user mode")
var userFlag *string = flag.String("user", "", "username to add")
var passFlag *string = flag.String("pass", "", "user's password")

func addUserMode() {
	passwdPath, err := filepath.Abs(*passwdFlag)
	if err != nil {
		log.Fatal(err)
	}
	username, password := *userFlag, *passFlag
	usermanager, err := usermanager.NewUserManager(passwdPath)
	if err != nil {
		log.Fatal(err)
	}
	if err := usermanager.AddUser(username, password); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Successfully added user: %s\n", username)
}

func serverMode() {
	passwdPath, err := filepath.Abs(*passwdFlag)
	if err != nil {
		log.Fatal(err)
	}
	rootDir, err := filepath.Abs(*rootFlag)
	if err != nil {
		log.Fatal(err)
	}
	port := *portFlag
	addr := fmt.Sprintf(":%d", port)
	log.Fatal(ftp.ListenAndServe(addr, rootDir, passwdPath))
}

func main() {
	flag.Parse()
	if *addUserFlag {
		addUserMode()
	} else {
		serverMode()
	}
}
