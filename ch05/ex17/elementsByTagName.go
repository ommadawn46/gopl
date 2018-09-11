package main

import (
  "os"
  "fmt"

  "golang.org/x/net/html"
  "../ex07/pretty"
)

func main() {
  if len(os.Args) != 2 {
    fmt.Fprintln(os.Stderr, "Usage: %s tagName", os.Args[0])
  }
  doc, err := html.Parse(os.Stdin)
  if err != nil {
    fmt.Fprintf(os.Stderr, "elementsByTagName: %v\n", err)
  }
  elements := ElementsByTagName(doc, os.Args[1])
  for _, element := range elements{
    fmt.Println(pretty.PrettyNode(element))
  }
}

func ElementsByTagName(doc *html.Node, tagName string) []*html.Node {
  pre := func(n *html.Node) bool {
		return n.Type == html.ElementNode && n.Data == tagName
  }
  return forEachNode(doc, pre, nil)
}

func forEachNode(n *html.Node, pre, post func(n *html.Node)bool) (elements []*html.Node) {
  if pre != nil {
    if pre(n) {
      elements = append(elements, n)
    }
  }

  for c := n.FirstChild; c != nil; c = c.NextSibling {
    if r := forEachNode(c, pre, post); r != nil {
      elements = append(elements, r...)
    }
  }

  if post != nil {
    if post(n) {
      elements = append(elements, n)
    }
  }
  return
}
