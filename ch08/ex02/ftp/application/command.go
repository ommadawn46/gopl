package application

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/ommadawn46/the_go_programming_language-training/ch08/ex02/ftp/util"
)

type Attribute int

const (
	NeedsArg          Attribute = iota // 要引数
	NeedsLogin                         // 要ログイン
	MustNotLogin                       // ログイン前のみ実行可能
	RecvDataBefore                     // コマンド実行前にデータを受け取る
	SendDataAfter                      // コマンド実行後にデータを転送する
	CloseAllDataConns                  // データ転送を中断する
)

type Command struct {
	exec   func(*Worker, string, []byte) (int, string, []byte)
	attrs  []Attribute
	syntax string
}

func (c *Command) HasAttribute(t Attribute) bool {
	for _, a := range c.attrs {
		if a == t {
			return true
		}
	}
	return false
}

var COMMANDS map[string]Command

func init() {
	COMMANDS = map[string]Command{
		"ABOR": {
			exec: func(w *Worker, arg string, _ []byte) (int, string, []byte) {
				return 226, "ABOR command successful", nil
			},
			attrs: []Attribute{
				NeedsLogin,
				CloseAllDataConns,
			},
			syntax: "ABOR",
		},
		"CWD": {
			exec: func(w *Worker, arg string, _ []byte) (int, string, []byte) {
				path := w.joinPath(arg)
				if !util.ExistsPath(path) {
					return 550, fmt.Sprintf("%s: No such file or directory", arg), nil
				}
				if !util.IsDirectory(path) {
					return 550, fmt.Sprintf("%s: Not a directory", arg), nil
				}
				w.workDir = strings.TrimLeft(path[len(*w.rootDirPtr):], "/\\")
				return 250, "CWD command successful", nil
			},
			attrs: []Attribute{
				NeedsLogin,
				NeedsArg,
			},
			syntax: "CWD <SP> <pathname>",
		},
		"DELE": {
			exec: func(w *Worker, arg string, _ []byte) (int, string, []byte) {
				path := w.joinPath(arg)
				if !util.ExistsPath(path) {
					return 550, fmt.Sprintf("%s: No such file or directory", arg), nil
				}
				if util.IsDirectory(path) {
					return 550, fmt.Sprintf("%s: Is a directory", arg), nil
				}
				if err := os.Remove(path); err != nil {
					return 550, fmt.Sprintf("%s: Failed to delete a file", arg), nil
				}
				return 250, "DELE command successful", nil
			},
			attrs: []Attribute{
				NeedsLogin,
				NeedsArg,
			},
			syntax: "DELE <SP> <pathname>",
		},
		"EPRT": {
			exec: func(w *Worker, arg string, _ []byte) (int, string, []byte) {
				addrPort := strings.Split(arg, "|")
				if len(addrPort) != 5 {
					return 501, "Illegal EPRT command", nil
				}
				protocol, err1 := strconv.Atoi(addrPort[1])
				addr := addrPort[2]
				if protocol == 2 {
					addr = "[" + addr + "]"
				}
				port, err2 := strconv.Atoi(addrPort[3])
				if err1 != nil || err2 != nil {
					return 501, "Illegal EPRT command", nil
				}
				w.DataAddr = fmt.Sprintf("%s:%d", addr, port)

				return 200, "EPRT command successful", nil
			},
			attrs: []Attribute{
				NeedsLogin,
				NeedsArg,
			},
			syntax: "EPRT <SP> <d> <net-prt> <d> <net-addr> <d> <tcp-port> <d>",
		},
		"EPSV": {
			exec: func(w *Worker, arg string, _ []byte) (int, string, []byte) {
				listener, err := net.Listen("tcp", ":0")
				if err != nil {
					return 500, "Cannot listen to port", nil
				}
				w.DataListener = listener.(*net.TCPListener)

				tcpAddr, _ := listener.Addr().(*net.TCPAddr)
				port := tcpAddr.Port

				w.PasvMode = true
				return 229, fmt.Sprintf(
					"Entering Extended Passive Mode (|||%d|)",
					port,
				), nil
			},
			attrs: []Attribute{
				NeedsLogin,
			},
			syntax: "EPSV",
		},
		"HELP": {
			exec: func(w *Worker, arg string, _ []byte) (int, string, []byte) {
				if arg == "" {
					var cmdNames []string
					for cmdName := range COMMANDS {
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
					), nil
				} else {
					cmd, ok := COMMANDS[strings.ToUpper(arg)]
					if !ok {
						return 502, fmt.Sprintf("Unknown command '%s'", arg), nil
					}
					return 214, fmt.Sprintf("Syntax: %s", cmd.syntax), nil
				}
			},
			attrs:  []Attribute{},
			syntax: "HELP [<sp> command]",
		},
		"LIST": {
			exec: func(w *Worker, arg string, _ []byte) (int, string, []byte) {
				return w.osCommandOutput("ls", "-la", w.joinPath(arg))
			},
			attrs: []Attribute{
				NeedsLogin,
				SendDataAfter,
			},
			syntax: "LIST [<SP> <pathname>]",
		},
		"MDTM": {
			exec: func(w *Worker, arg string, _ []byte) (int, string, []byte) {
				path := w.joinPath(arg)
				if !util.ExistsPath(path) {
					return 550, fmt.Sprintf("%s: No such file or directory", arg), nil
				}
				if util.IsDirectory(path) {
					return 550, fmt.Sprintf("%s: Is a directory", arg), nil
				}
				fileInfo, _ := os.Stat(path)
				return 213, fileInfo.ModTime().Format("20060102150405"), nil
			},
			attrs: []Attribute{
				NeedsLogin,
				NeedsArg,
			},
			syntax: "MDTM <SP> <pathname>",
		},
		"MKD": {
			exec: func(w *Worker, arg string, _ []byte) (int, string, []byte) {
				path := w.joinPath(arg)
				if !util.ExistsPath(filepath.Dir(path)) {
					return 550, fmt.Sprintf("%s: No such file or directory", arg), nil
				}
				if util.ExistsPath(path) && !util.IsDirectory(path) {
					return 550, fmt.Sprintf("%s: Not a directory", arg), nil
				}
				if err := os.Mkdir(path, 0755); err != nil {
					return 550, fmt.Sprintf("%s: Failed to make a directory", arg), nil
				}
				return 257, fmt.Sprintf("\"%s\" - Directory successfully created", arg), nil
			},
			attrs: []Attribute{
				NeedsLogin,
				NeedsArg,
			},
			syntax: "MKD <SP> <pathname>",
		},
		"MODE": {
			exec: func(w *Worker, arg string, _ []byte) (int, string, []byte) {
				mode := strings.ToUpper(arg)
				switch mode {
				case "S":
					return 200, "Mode set to S", nil
				case "B", "C":
					return 504, fmt.Sprintf("'MODE %s' unsupported transfer mode", mode), nil
				default:
					return 501, fmt.Sprintf("'MODE %s' unrecognized transfer mode", mode), nil
				}
			},
			attrs: []Attribute{
				NeedsLogin,
				NeedsArg,
			},
			syntax: "MODE <SP> <mode-code>",
		},
		"NLST": {
			exec: func(w *Worker, arg string, _ []byte) (int, string, []byte) {
				return w.osCommandOutput("ls", w.joinPath(arg))
			},
			attrs: []Attribute{
				NeedsLogin,
				SendDataAfter,
			},
			syntax: "NLST [<SP> <pathname>]",
		},
		"NOOP": {
			exec: func(w *Worker, arg string, _ []byte) (int, string, []byte) {
				return 200, "NOOP command successful", nil
			},
			attrs:  []Attribute{},
			syntax: "NOOP",
		},
		"PASS": {
			exec: func(w *Worker, arg string, _ []byte) (int, string, []byte) {
				if w.loginName == "" {
					return 503, "Login with USER first", nil
				}
				if !w.userMgr.Auth(w.loginName, arg) {
					w.loginName = ""
					return 530, "Login incorrect", nil
				}
				w.LoggedIn = true
				return 230, fmt.Sprintf("User %v logged in", w.loginName), nil
			},
			attrs: []Attribute{
				MustNotLogin,
				NeedsArg,
			},
			syntax: "PASS <SP> <password>",
		},
		"PASV": {
			exec: func(w *Worker, arg string, _ []byte) (int, string, []byte) {
				listener, err := net.Listen("tcp4", "127.0.0.1:0")
				if err != nil {
					return 500, "Cannot listen to port", nil
				}
				w.DataListener = listener.(*net.TCPListener)

				tcpAddr, _ := listener.Addr().(*net.TCPAddr)
				addr, port := tcpAddr.IP.String(), tcpAddr.Port
				port1, port2 := port/0x100, port%0x100

				w.PasvMode = true
				return 227, fmt.Sprintf(
					"Entering Passive Mode (%s,%d,%d)",
					strings.Replace(addr, ".", ",", -1), port1, port2,
				), nil
			},
			attrs: []Attribute{
				NeedsLogin,
			},
			syntax: "PASV",
		},
		"PORT": {
			exec: func(w *Worker, arg string, _ []byte) (int, string, []byte) {
				addrPort := strings.Split(arg, ",")
				if len(addrPort) != 6 {
					return 501, "Illegal PORT command", nil
				}
				addr := strings.Join(addrPort[:4], ".")
				port1, portErr1 := strconv.Atoi(addrPort[4])
				port2, portErr2 := strconv.Atoi(addrPort[5])
				if portErr1 != nil || portErr2 != nil {
					return 501, "Illegal PORT command", nil
				}
				port := port1*0x100 + port2
				w.DataAddr = fmt.Sprintf("%s:%d", addr, port)

				return 200, "PORT command successful", nil
			},
			attrs: []Attribute{
				NeedsLogin,
				NeedsArg,
			},
			syntax: "PORT <SP> <host-port>",
		},
		"PWD": {
			exec: func(w *Worker, arg string, _ []byte) (int, string, []byte) {
				return 257, fmt.Sprintf("\"/%s\" is the current directory", w.workDir), nil
			},
			attrs: []Attribute{
				NeedsLogin,
			},
			syntax: "PWD",
		},
		"QUIT": {
			exec: func(w *Worker, arg string, _ []byte) (int, string, []byte) {
				return 221, "Goodbye", nil
			},
			attrs:  []Attribute{},
			syntax: "QUIT",
		},
		"RETR": {
			exec: func(w *Worker, arg string, _ []byte) (int, string, []byte) {
				path := w.joinPath(arg)
				if !util.ExistsPath(path) {
					return 550, fmt.Sprintf("%s: No such file or directory", arg), nil
				}
				if util.IsDirectory(path) {
					return 550, fmt.Sprintf("%s: Not a regular file", arg), nil
				}
				fileData, err := ioutil.ReadFile(path)
				if err != nil {
					return 500, "Failed to read file", nil
				}
				return 226, "Transfer complete", fileData
			},
			attrs: []Attribute{
				NeedsLogin,
				NeedsArg,
				SendDataAfter,
			},
			syntax: "RETR <SP> <pathname>",
		},
		"RMD": {
			exec: func(w *Worker, arg string, _ []byte) (int, string, []byte) {
				path := w.joinPath(arg)
				if !util.ExistsPath(path) {
					return 550, fmt.Sprintf("%s: No such file or directory", arg), nil
				}
				if !util.IsDirectory(path) {
					return 550, fmt.Sprintf("%s: Not a directory", arg), nil
				}
				if err := os.Remove(path); err != nil {
					return 550, fmt.Sprintf("%s: Failed to delete a directory", arg), nil
				}
				return 250, "RMD command successful", nil
			},
			attrs: []Attribute{
				NeedsLogin,
				NeedsArg,
			},
			syntax: "RMD <SP> <pathname>",
		},
		"RNFR": {
			exec: func(w *Worker, arg string, _ []byte) (int, string, []byte) {
				path := w.joinPath(arg)
				if !util.ExistsPath(path) {
					return 550, fmt.Sprintf("%s: No such file or directory", arg), nil
				}
				w.renameFrom = path
				return 350, "File or directory exists, ready for destination name", nil
			},
			attrs: []Attribute{
				NeedsLogin,
				NeedsArg,
			},
			syntax: "RNFR <SP> <pathname>",
		},
		"RNTO": {
			exec: func(w *Worker, arg string, _ []byte) (int, string, []byte) {
				path := w.joinPath(arg)
				if !util.ExistsPath(filepath.Dir(path)) {
					return 550, fmt.Sprintf("%s: No such file or directory", arg), nil
				}
				if util.IsDirectory(path) {
					return 550, fmt.Sprintf("%s: Is a directory", arg), nil
				}
				if err := os.Rename(w.renameFrom, path); err != nil {
					return 550, fmt.Sprintf("%s: Failed to rename a file", arg), nil
				}
				w.renameFrom = ""
				return 250, "Rename successful", nil
			},
			attrs: []Attribute{
				NeedsLogin,
				NeedsArg,
			},
			syntax: "RNTO <SP> <pathname>",
		},
		"SIZE": {
			exec: func(w *Worker, arg string, _ []byte) (int, string, []byte) {
				if w.TransType == ASCII {
					return 550, "SIZE not allowed in ASCII mode", nil
				}
				path := w.joinPath(arg)
				if !util.ExistsPath(path) {
					return 550, fmt.Sprintf("%s: No such file or directory", arg), nil
				}
				if util.IsDirectory(path) {
					return 550, fmt.Sprintf("%s: Is a directory", arg), nil
				}
				fileInfo, _ := os.Stat(path)
				return 213, fmt.Sprintf("%d", fileInfo.Size()), nil
			},
			attrs: []Attribute{
				NeedsLogin,
				NeedsArg,
			},
			syntax: "SIZE <SP> <pathname>",
		},
		"STOR": {
			exec: func(w *Worker, arg string, data []byte) (int, string, []byte) {
				path := w.joinPath(arg)
				if !util.ExistsPath(filepath.Dir(path)) {
					return 550, fmt.Sprintf("%s: No such file or directory", arg), nil
				}
				if util.IsDirectory(path) {
					return 550, fmt.Sprintf("%s: Not a regular file", arg), nil
				}
				if err := ioutil.WriteFile(path, data, 0644); err != nil {
					return 500, "Failed to write a file", nil
				}
				return 226, "Transfer complete", nil
			},
			attrs: []Attribute{
				NeedsLogin,
				NeedsArg,
				RecvDataBefore,
			},
			syntax: "STOR <SP> <pathname>",
		},
		"STRU": {
			exec: func(w *Worker, arg string, _ []byte) (int, string, []byte) {
				mode := strings.ToUpper(arg)
				switch mode {
				case "F":
					return 200, "Structure set to F", nil
				case "R", "P":
					return 504, fmt.Sprintf("'MODE %s' unsupported structure type", mode), nil
				default:
					return 501, fmt.Sprintf("'MODE %s' unrecognized structure type", mode), nil
				}
			},
			attrs: []Attribute{
				NeedsLogin,
				NeedsArg,
			},
			syntax: "STRU <SP> <structure-code>",
		},
		"TYPE": {
			exec: func(w *Worker, arg string, _ []byte) (int, string, []byte) {
				mode := strings.ToUpper(arg)
				switch mode {
				case "A":
					w.TransType = ASCII
				case "I":
					w.TransType = BINARY
				default:
					return 500, fmt.Sprintf("'TYPE %s' not understood", mode), nil
				}
				return 200, fmt.Sprintf("Type set to %s", mode), nil
			},
			attrs: []Attribute{
				NeedsLogin,
				NeedsArg,
			},
			syntax: "TYPE <SP> <type-code>",
		},
		"USER": {
			exec: func(w *Worker, arg string, _ []byte) (int, string, []byte) {
				w.loginName = arg
				return 331, fmt.Sprintf("Password required for %v", w.loginName), nil
			},
			attrs: []Attribute{
				MustNotLogin,
				NeedsArg,
			},
			syntax: "USER <SP> <username>",
		},
	}
}
