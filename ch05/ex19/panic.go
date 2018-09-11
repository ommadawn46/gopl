package main

import (
	"fmt"
)

func echo(s string) (r string) {
	defer func() {
		if p := recover(); p != nil {
			r = p.(string)
		}
	}()
	panic(s)
}

func main() {
	fmt.Println(echo("Panic and recovered"))
}
