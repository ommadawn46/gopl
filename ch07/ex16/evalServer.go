package main

import (
	"fmt"
	"log"
	"net/http"

	"../ex14"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		exprStr := r.FormValue("expr")
		if exprStr == "" {
			http.Error(w, "no expr", http.StatusBadRequest)
			return
		}
		expr, err := eval.Parse(exprStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Fprintf(w, "%s = %v", exprStr, expr.Eval(eval.Env{}))
	})
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}
