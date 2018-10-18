package usermanager

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/ommadawn46/the_go_programming_language-training/ch08/ex02/ftp/util"
)

type User struct {
	name string
	salt string
	hash string
}

type UserManager struct {
	passwdPath string
}

func (u *UserManager) loadUsers() ([]User, error) {
	passwdFile, err := os.Open(u.passwdPath)
	if err != nil {
		return nil, err
	}
	defer passwdFile.Close()

	var users []User
	reader := csv.NewReader(passwdFile)
	reader.Comma = ':'
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		users = append(users, User{record[0], record[1], record[2]})
	}
	return users, nil
}

func (u *UserManager) saveUser(user User) error {
	passwdFile, err := os.OpenFile(u.passwdPath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}
	defer passwdFile.Close()

	newLine := fmt.Sprintf("%s:%s:%s\n", user.name, user.salt, user.hash)
	_, err = passwdFile.WriteString(newLine)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserManager) Auth(login, pass string) bool {
	users, err := u.loadUsers()
	if err != nil {
		log.Println(err)
		return false
	}
	for _, user := range users {
		if user.name == login && user.hash == util.Hash(user.salt+pass) {
			return true
		}
	}
	return false
}

func (u *UserManager) AddUser(username, pass string) error {
	if username == "" {
		return fmt.Errorf("Username is empty")
	}
	users, err := u.loadUsers()
	if err != nil {
		return err
	}
	for _, user := range users {
		if user.name == username {
			return fmt.Errorf("%s already exists", username)
		}
	}
	salt := util.GenerateSalt()
	newUser := User{username, salt, util.Hash(salt + pass)}
	return u.saveUser(newUser)
}

func NewUserManager(passwdPath string) (*UserManager, error) {
	if !util.ExistsPath(passwdPath) {
		return nil, fmt.Errorf("file does not exist: %s", passwdPath)
	}
	if util.IsDirectory(passwdPath) {
		return nil, fmt.Errorf("not a regular file: %s", passwdPath)
	}
	return &UserManager{passwdPath}, nil
}
