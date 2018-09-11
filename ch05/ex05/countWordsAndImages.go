package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func CountWordsAndImages(url string) (words, images int, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		err = fmt.Errorf("parsing HTML: %s", err)
		return
	}
	words, images = countWordsAndImages(doc)
	return
}

func countWordsAndImages(n *html.Node) (words, images int) {
	if n == nil {
		return
	}
	switch n.Type {
	case html.ElementNode:
		if n.Data == "img" {
			images++
		}
	case html.TextNode:
		words += countWords(n.Data)
	}
	if n.Type != html.ElementNode || n.Data != "script" && n.Data != "style" {
		w, i := countWordsAndImages(n.FirstChild)
		words, images = words+w, images+i
	}
	w, i := countWordsAndImages(n.NextSibling)
	words, images = words+w, images+i
	return
}

func countWords(s string) (words int) {
	var text = bufio.NewScanner(strings.NewReader(s))
	text.Split(bufio.ScanWords)
	for text.Scan() {
		words++
	}
	return
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s url", os.Args[0])
		os.Exit(1)
	}
	url := os.Args[1]
	words, images, err := CountWordsAndImages(url)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	fmt.Printf("words:\t%d\nimages:\t%d\n", words, images)
}
