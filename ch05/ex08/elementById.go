package main

import (
	"fmt"
	"os"

	"../ex07/pretty"
	"golang.org/x/net/html"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: %s id", os.Args[0])
	}
	doc, err := html.Parse(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "elementByID: %v\n", err)
	}
	element := ElementByID(doc, os.Args[1])
	fmt.Println(pretty.PrettyNode(element))
}

func ElementByID(doc *html.Node, id string) *html.Node {
	pre := func(n *html.Node) bool {
		if n.Type == html.ElementNode {
			for _, a := range n.Attr {
				if a.Key == "id" && a.Val == id {
					return false
				}
			}
		}
		return true
	}
	return forEachNode(doc, pre, nil)
}

func forEachNode(n *html.Node, pre, post func(n *html.Node) bool) *html.Node {
	if pre != nil {
		if !pre(n) {
			return n
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if r := forEachNode(c, pre, post); r != nil {
			return r
		}
	}

	if post != nil {
		if !post(n) {
			return n
		}
	}
	return nil
}
