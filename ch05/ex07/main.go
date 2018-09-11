package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"./pretty"
)

func main() {
	bytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	res, err := pretty.Pretty(string(bytes))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)
}
