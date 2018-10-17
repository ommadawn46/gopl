package ftp

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type Attribute int

const (
	needsArg     Attribute = iota // 要引数
	needsLogin                    // 要ログイン
	mustNotLogin                  // ログイン前のみ実行可能
	useDataPort                   // データポートを使用する
)

type Command struct {
	exec   func(*Worker, string) (int, string)
	attrs  []Attribute
	syntax string
}

func (c *Command) hasAttribute(t Attribute) bool {
	for _, a := range c.attrs {
		if a == t {
			return true
		}
	}
	return false
}

var _COMMANDS map[string]Command

func init() {
	_COMMANDS = map[string]Command{
		"CWD": {
			exec: func(w *Worker, arg string) (int, string) {
				path := w.joinPath(arg)
				if !existsPath(path) {
					return 550, fmt.Sprintf("%s: No such file or directory", arg)
				}
				if !isDirectory(path) {
					return 550, fmt.Sprintf("%s: Not a directory", arg)
				}
				w.workDir = strings.TrimLeft(path[len(w.rootDir):], "/\\")
				return 250, "CWD command successful"
			},
			attrs: []Attribute{
				needsLogin,
				needsArg,
			},
			syntax: "CWD <SP> <pathname>",
		},
		"DELE": {
			exec: func(w *Worker, arg string) (int, string) {
				path := w.joinPath(arg)
				if !existsPath(path) {
					return 550, fmt.Sprintf("%s: No such file or directory", arg)
				}
				if isDirectory(path) {
					return 550, fmt.Sprintf("%s: Is a directory", arg)
				}
				if err := os.Remove(path); err != nil {
					return 550, fmt.Sprintf("%s: Failed to delete a file", arg)
				}
				return 250, "DELE command successful"
			},
			attrs: []Attribute{
				needsLogin,
				needsArg,
			},
			syntax: "DELE <SP> <pathname>",
		},
		"HELP": {
			exec: func(w *Worker, arg string) (int, string) {
				if arg == "" {
					var cmdNames []string
					for cmdName := range _COMMANDS {
						cmdNames = append(cmdNames, cmdName)
					}
					sort.Strings(cmdNames)
					var helpMessages []string
					for i := 0; i < len(cmdNames); i += 8 {
						helpMessages = append(helpMessages, strings.Join(cmdNames[i:i+8], "\t"))
					}
					return 214, fmt.Sprintf(
						"The following commands are recognized:\n%s\nHELP command successful",
						strings.Join(helpMessages, "\n"),
					)
				} else {
					cmd, ok := _COMMANDS[strings.ToUpper(arg)]
					if !ok {
						return 502, fmt.Sprintf("Unknown command '%s'", arg)
					}
					return 214, fmt.Sprintf("Syntax: %s", cmd.syntax)
				}
			},
			attrs:  []Attribute{},
			syntax: "HELP [<sp> command]",
		},
		"LIST": {
			exec: func(w *Worker, arg string) (int, string) {
				return w.sendOsCommandOutput("ls", "-la", w.joinPath(arg))
			},
			attrs: []Attribute{
				needsLogin,
				useDataPort,
			},
			syntax: "LIST [<SP> <pathname>]",
		},
		"MDTM": {
			exec: func(w *Worker, arg string) (int, string) {
				path := w.joinPath(arg)
				if !existsPath(path) {
					return 550, fmt.Sprintf("%s: No such file or directory", arg)
				}
				if isDirectory(path) {
					return 550, fmt.Sprintf("%s: Is a directory", arg)
				}
				fileInfo, _ := os.Stat(path)
				return 213, fileInfo.ModTime().Format("20060102150405")
			},
			attrs: []Attribute{
				needsLogin,
				needsArg,
			},
			syntax: "MDTM <SP> <pathname>",
		},
		"MKD": {
			exec: func(w *Worker, arg string) (int, string) {
				path := w.joinPath(arg)
				if !existsPath(filepath.Dir(path)) {
					return 550, fmt.Sprintf("%s: No such file or directory", arg)
				}
				if existsPath(path) && !isDirectory(path) {
					return 550, fmt.Sprintf("%s: Not a directory", arg)
				}
				if err := os.Mkdir(path, 0755); err != nil {
					return 550, fmt.Sprintf("%s: Failed to make a directory", arg)
				}
				return 257, fmt.Sprintf("\"%s\" - Directory successfully created", arg)
			},
			attrs: []Attribute{
				needsLogin,
				needsArg,
			},
			syntax: "MKD <SP> <pathname>",
		},
		"MODE": {
			exec: func(w *Worker, arg string) (int, string) {
				mode := strings.ToUpper(arg)
				switch mode {
				case "S":
					return 200, "Mode set to S"
				case "B", "C":
					return 504, fmt.Sprintf("'MODE %s' unsupported transfer mode", mode)
				default:
					return 501, fmt.Sprintf("'MODE %s' unrecognized transfer mode", mode)
				}
			},
			attrs: []Attribute{
				needsLogin,
				needsArg,
			},
			syntax: "MODE <SP> <mode-code>",
		},
		"NLST": {
			exec: func(w *Worker, arg string) (int, string) {
				return w.sendOsCommandOutput("ls", w.joinPath(arg))
			},
			attrs: []Attribute{
				needsLogin,
				useDataPort,
			},
			syntax: "NLST [<SP> <pathname>]",
		},
		"NOOP": {
			exec: func(w *Worker, arg string) (int, string) {
				return 200, "NOOP command successful"
			},
			attrs:  []Attribute{},
			syntax: "NOOP",
		},
		"PASS": {
			exec: func(w *Worker, arg string) (int, string) {
				if w.username == "" {
					return 503, "Login with USER first"
				}
				if !auth(w.username, arg) {
					w.username = ""
					return 530, "Login incorrect"
				}
				w.loggedIn = true
				return 230, fmt.Sprintf("User %v logged in", w.username)
			},
			attrs: []Attribute{
				mustNotLogin,
				needsArg,
			},
			syntax: "PASS <SP> <password>",
		},
		"PASV": {
			exec: func(w *Worker, arg string) (int, string) {
				listener, err := net.Listen("tcp", "127.0.0.1:0")
				if err != nil {
					return 500, "Cannot listen to port"
				}
				w.dataPort.listener = listener

				tcpAddr, _ := listener.Addr().(*net.TCPAddr)
				addr, port := tcpAddr.IP.String(), tcpAddr.Port
				port1, port2 := port/0x100, port%0x100

				w.dataPort.pasvMode = true
				return 227, fmt.Sprintf(
					"Entering Passive Mode (%s,%d,%d)",
					strings.Replace(addr, ".", ",", -1), port1, port2,
				)
			},
			attrs: []Attribute{
				needsLogin,
			},
			syntax: "PASV",
		},
		"PORT": {
			exec: func(w *Worker, arg string) (int, string) {
				addrPort := strings.Split(arg, ",")
				if len(addrPort) != 6 {
					return 501, "Illegal PORT command"
				}

				addr := strings.Join(addrPort[:4], ".")
				port1, portErr1 := strconv.Atoi(addrPort[4])
				port2, portErr2 := strconv.Atoi(addrPort[5])
				if portErr1 != nil || portErr2 != nil {
					return 501, "Illegal PORT command"
				}
				port := port1*0x100 + port2
				w.dataPort.addr = fmt.Sprintf("%s:%d", addr, port)

				return 200, "PORT command successful"
			},
			attrs: []Attribute{
				needsLogin,
				needsArg,
			},
			syntax: "PORT <SP> <host-port>",
		},
		"PWD": {
			exec: func(w *Worker, arg string) (int, string) {
				return 257, fmt.Sprintf("\"/%s\" is the current directory", w.workDir)
			},
			attrs: []Attribute{
				needsLogin,
			},
			syntax: "PWD",
		},
		"QUIT": {
			exec: func(w *Worker, arg string) (int, string) {
				return 221, "Goodbye"
			},
			attrs:  []Attribute{},
			syntax: "QUIT",
		},
		"RETR": {
			exec: func(w *Worker, arg string) (int, string) {
				dataConn, err := w.dataPort.connect()
				if err != nil {
					return 425, err.Error()
				}
				defer dataConn.Close()

				path := w.joinPath(arg)
				if !existsPath(path) {
					return 550, fmt.Sprintf("%s: No such file or directory", arg)
				}
				if isDirectory(path) {
					return 550, fmt.Sprintf("%s: Not a regular file", arg)
				}

				fileData, err := ioutil.ReadFile(path)
				if err != nil {
					return 500, "Failed to read file"
				}

				switch w.transferType {
				case _BINARY:
					err = dataConn.sendAll(fileData)
				case _ASCII:
					err = dataConn.sendAllAsAscii(fileData)
				default:
					err = fmt.Errorf("Invalid transfer type")
				}
				if err != nil {
					return 426, fmt.Sprintf("%v: Failed to transfer %s", err, w.transferType)
				}
				return 226, "Transfer complete"
			},
			attrs: []Attribute{
				needsLogin,
				needsArg,
				useDataPort,
			},
			syntax: "RETR <SP> <pathname>",
		},
		"RMD": {
			exec: func(w *Worker, arg string) (int, string) {
				path := w.joinPath(arg)
				if !existsPath(path) {
					return 550, fmt.Sprintf("%s: No such file or directory", arg)
				}
				if !isDirectory(path) {
					return 550, fmt.Sprintf("%s: Not a directory", arg)
				}
				if err := os.Remove(path); err != nil {
					return 550, fmt.Sprintf("%s: Failed to delete a directory", arg)
				}
				return 250, "RMD command successful"
			},
			attrs: []Attribute{
				needsLogin,
				needsArg,
			},
			syntax: "RMD <SP> <pathname>",
		},
		"RNFR": {
			exec: func(w *Worker, arg string) (int, string) {
				path := w.joinPath(arg)
				if !existsPath(path) {
					return 550, fmt.Sprintf("%s: No such file or directory", arg)
				}
				w.renameFrom = path
				return 350, "File or directory exists, ready for destination name"
			},
			attrs: []Attribute{
				needsLogin,
				needsArg,
			},
			syntax: "RNFR <SP> <pathname>",
		},
		"RNTO": {
			exec: func(w *Worker, arg string) (int, string) {
				path := w.joinPath(arg)
				if !existsPath(filepath.Dir(path)) {
					return 550, fmt.Sprintf("%s: No such file or directory", arg)
				}
				if isDirectory(path) {
					return 550, fmt.Sprintf("%s: Is a directory", arg)
				}
				if err := os.Rename(w.renameFrom, path); err != nil {
					return 550, fmt.Sprintf("%s: Failed to rename a file", arg)
				}
				w.renameFrom = ""
				return 250, "Rename successful"
			},
			attrs: []Attribute{
				needsLogin,
				needsArg,
			},
			syntax: "RNTO <SP> <pathname>",
		},
		"SIZE": {
			exec: func(w *Worker, arg string) (int, string) {
				if w.transferType == _ASCII {
					return 550, "SIZE not allowed in ASCII mode"
				}
				path := w.joinPath(arg)
				if !existsPath(path) {
					return 550, fmt.Sprintf("%s: No such file or directory", arg)
				}
				if isDirectory(path) {
					return 550, fmt.Sprintf("%s: Is a directory", arg)
				}
				fileInfo, _ := os.Stat(path)
				return 213, fmt.Sprintf("%d", fileInfo.Size())
			},
			attrs: []Attribute{
				needsLogin,
				needsArg,
			},
			syntax: "SIZE <SP> <pathname>",
		},
		"STOR": {
			exec: func(w *Worker, arg string) (int, string) {
				dataConn, err := w.dataPort.connect()
				if err != nil {
					return 425, err.Error()
				}
				defer dataConn.Close()

				path := w.joinPath(arg)
				if !existsPath(filepath.Dir(path)) {
					return 550, fmt.Sprintf("%s: No such file or directory", arg)
				}
				if isDirectory(path) {
					return 550, fmt.Sprintf("%s: Not a regular file", arg)
				}

				var buf []byte
				switch w.transferType {
				case _BINARY:
					buf, err = dataConn.readAll()
				case _ASCII:
					buf, err = dataConn.readAllAsAscii()
				default:
					buf, err = nil, fmt.Errorf("Invalid transfer type")
				}
				if err != nil {
					return 426, fmt.Sprintf("%v: Failed to transfer %s", err, w.transferType)
				}

				if err := ioutil.WriteFile(path, buf, 0644); err != nil {
					return 500, "Failed to write a file"
				}
				return 226, "Transfer complete"
			},
			attrs: []Attribute{
				needsLogin,
				needsArg,
				useDataPort,
			},
			syntax: "STOR <SP> <pathname>",
		},
		"STRU": {
			exec: func(w *Worker, arg string) (int, string) {
				mode := strings.ToUpper(arg)
				switch mode {
				case "F":
					return 200, "Structure set to F"
				case "R", "P":
					return 504, fmt.Sprintf("'MODE %s' unsupported structure type", mode)
				default:
					return 501, fmt.Sprintf("'MODE %s' unrecognized structure type", mode)
				}
			},
			attrs: []Attribute{
				needsLogin,
				needsArg,
			},
			syntax: "STRU <SP> <structure-code>",
		},
		"TYPE": {
			exec: func(w *Worker, arg string) (int, string) {
				mode := strings.ToUpper(arg)
				switch mode {
				case "A":
					w.transferType = _ASCII
				case "I":
					w.transferType = _BINARY
				default:
					return 500, fmt.Sprintf("'TYPE %s' not understood", mode)
				}
				return 200, fmt.Sprintf("Type set to %s", mode)
			},
			attrs: []Attribute{
				needsLogin,
				needsArg,
			},
			syntax: "TYPE <SP> <type-code>",
		},
		"USER": {
			exec: func(w *Worker, arg string) (int, string) {
				w.username = arg
				return 331, fmt.Sprintf("Password required for %v", w.username)
			},
			attrs: []Attribute{
				mustNotLogin,
				needsArg,
			},
			syntax: "USER <SP> <username>",
		},
	}
}
