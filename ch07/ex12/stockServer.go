package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func main() {
	db := database{"shoes": 50, "socks": 5}
	mux := http.NewServeMux()
	mux.HandleFunc("/list", db.list)
	mux.HandleFunc("/price", db.price)
	mux.HandleFunc("/create", db.create)
	mux.HandleFunc("/update", db.update)
	mux.HandleFunc("/delete", db.delete)
	log.Fatal(http.ListenAndServe("localhost:8000", mux))
}

type database map[string]int

type Item struct {
	Name  string
	Price int
}

var itemTable = template.Must(template.New("itemTable").Parse(`
<html>
   <body>
      <h1>ItemTable</h1>
      <table border="1">
         <thead>
            <tr>
               <th>Name</th>
               <th>Price</th>
            </tr>
         </thead>
         <tbody>
            {{range.}}
            <tr>
               <td>{{.Name}}</td>
               <td>{{.Price}}</td>
            </tr>
            {{end}}
         </tbody>
      </table>
   </body>
</html>
`))

func (db database) list(w http.ResponseWriter, req *http.Request) {
	items := []*Item{}
	for name, price := range db {
		items = append(items, &Item{name, price})
	}
	if err := itemTable.Execute(w, items); err != nil {
		log.Fatal(err)
	}
}

func (db database) price(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	price, ok := db[item]
	if ok {
		fmt.Fprintf(w, "$%d\n", price)
	} else {
		w.WriteHeader(http.StatusNotFound) // 404
		fmt.Fprintf(w, "no such item: %q\n", item)
	}
}

func (db database) create(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	priceStr := req.URL.Query().Get("price")
	price, err := strconv.Atoi(priceStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "invalid price: %s\n", priceStr)
	} else {
		db[item] = price
		fmt.Fprintf(w, "created %s: $%d\n", item, price)
	}
}

func (db database) update(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	priceStr := req.URL.Query().Get("price")
	price, err := strconv.Atoi(priceStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "invalid price: %s\n", priceStr)
	} else {
		if _, ok := db[item]; !ok {
			fmt.Fprintf(w, "%s does not exist\n", item)
		} else {
			db[item] = price
			fmt.Fprintf(w, "updated %s: $%d\n", item, price)
		}
	}
}

func (db database) delete(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	if _, ok := db[item]; !ok {
		fmt.Fprintf(w, "%s does not exist\n", item)
	} else {
		delete(db, item)
		fmt.Fprintf(w, "deleted %s\n", item)
	}
}
