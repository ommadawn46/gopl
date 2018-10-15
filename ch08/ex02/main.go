package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"

	"./ftp"
)

var portFlag *int = flag.Int("port", 21, "listening port")
var rootFlag *string = flag.String("root", "./", "root directory")

func main() {
	flag.Parse()

	port := *portFlag
	addr := fmt.Sprintf(":%d", port)

	rootDir, err := filepath.Abs(*rootFlag)
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(ftp.ListenAndServe(addr, rootDir))
}
