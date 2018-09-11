package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"./links"
)

func breadthFirst(f func(item string) []string, worklist []string) {
	seen := make(map[string]bool)
	for len(worklist) > 0 {
		items := worklist
		worklist = nil
		for _, item := range items {
			if !seen[item] {
				seen[item] = true
				worklist = append(worklist, f(item)...)
			}
		}
	}
}

func crawl(url string) []string {
	fmt.Println(url)
	err := save(url)
	if err != nil {
		log.Fatal(err)
	}
	list, err := links.Extract(url)
	if err != nil {
		log.Fatal(err)
	}
	return list
}

var targetHost string

func save(s_url string) error {
	u_url, err := url.Parse(s_url)
	if err != nil {
		return err
	}
	if targetHost == "" {
		targetHost = u_url.Host
	}
	if u_url.Host != targetHost {
		return nil
	}

	savepath := filepath.Join("data", u_url.Host, u_url.Path)
	if filepath.Ext(u_url.Path) == "" {
		savepath = filepath.Join(savepath, "index.html")
	}
	if err = os.MkdirAll(filepath.Dir(savepath), 0755); err != nil {
		return err
	}
	resp, err := http.Get(s_url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	file, err := os.Create(savepath)
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err = io.Copy(file, resp.Body); err != nil {
		return err
	}
	return nil
}

func main() {
	breadthFirst(crawl, os.Args[1:])
}
