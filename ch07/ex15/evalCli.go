package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"

	"../ex14"
)

func main() {
	stdin := bufio.NewScanner(os.Stdin)
	fmt.Printf("expr > ")
	stdin.Scan()

	exprStr := stdin.Text()
	expr, err := eval.Parse(exprStr)
	if err != nil {
		log.Fatal(err)
	}

	env := eval.Env{}
	for _, v := range eval.ExtractAllVar(expr) {
		fmt.Printf("%s = ", v)
		stdin.Scan()
		envStr := stdin.Text()
		val, err := strconv.ParseFloat(envStr, 64)
		if err != nil {
			log.Fatal(err)
		}
		env[v] = val
	}

	fmt.Println(expr.Eval(env))
}
