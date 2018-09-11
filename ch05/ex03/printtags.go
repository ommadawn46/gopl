package main

import (
  "fmt"
  "os"

  "golang.org/x/net/html"
)

func main() {
  doc, err := html.Parse(os.Stdin)
  if err != nil {
    fmt.Fprintf(os.Stderr, "printtags: %v\n", err)
  }
  printtags(doc)
}

func printtags(n *html.Node) {
  if n == nil {
    return
  }
  if n.Type == html.TextNode {
      fmt.Println(n.Data)
  }
  // scriptタグとstyleタグの中身は見ない
  if n.Type != html.ElementNode || n.Data != "script" && n.Data != "style" {
    printtags(n.FirstChild)
  }
  printtags(n.NextSibling)
}
