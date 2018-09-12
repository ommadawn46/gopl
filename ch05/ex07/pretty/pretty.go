package pretty

import (
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

func Pretty(s string) (string, error) {
	doc, err := html.Parse(strings.NewReader(s))
	if err != nil {
		return "", err
	}
	return PrettyNode(doc), nil
}

func PrettyNode(n *html.Node) string {
	return forEachNode(n, startElement, endElement)
}

func forEachNode(n *html.Node, pre, post func(n *html.Node) string) (result string) {
	if pre != nil {
		result += pre(n)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result += forEachNode(c, pre, post)
	}

	if post != nil {
		result += post(n)
	}
	return
}

var depth int

func startElement(n *html.Node) (result string) {
	if n.Type == html.ElementNode {
		closing := ">"
		if n.FirstChild == nil {
			closing = "/>"
		}
		attrs := ""
		for _, a := range n.Attr {
			attrs += fmt.Sprintf(" %s=%q", a.Key, a.Val)
		}
		result += fmt.Sprintf("%*s<%s%s%s\n", depth*2, "", n.Data, attrs, closing)
		depth++
	}
	if n.Type == html.TextNode {
		text := strings.TrimSpace(n.Data)
		for _, line := range strings.Split(text, "\n") {
			if line != "" {
				result += fmt.Sprintf("%*s%s\n", depth*2, "", line)
			}
		}
	}
	return
}

func endElement(n *html.Node) (result string) {
	if n.Type == html.ElementNode {
		depth--
		if n.FirstChild != nil {
			result += fmt.Sprintf("%*s</%s>\n", depth*2, "", n.Data)
		}
	}
	return
}
