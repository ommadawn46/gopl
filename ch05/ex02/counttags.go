package main

import (
  "fmt"
  "os"

  "golang.org/x/net/html"
)

func main() {
  doc, err := html.Parse(os.Stdin)
  if err != nil {
    fmt.Fprintf(os.Stderr, "counttags: %v\n", err)
  }
  for tag, n := range counttags(map[string]int{}, doc) {
    fmt.Printf("%5d\t%s\n", n, tag)
  }
}

func counttags(counter map[string]int, n *html.Node) map[string]int {
  if n == nil {
    return counter
  }
  if n.Type == html.ElementNode {
      counter[n.Data] += 1
  }
	counter = counttags(counter, n.FirstChild)
	return counttags(counter, n.NextSibling)
}
