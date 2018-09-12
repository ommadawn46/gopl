package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
)

type Tag struct {
	name  string
	id    string
	class string
}

func main() {
	dec := xml.NewDecoder(os.Stdin)
	var stack []Tag
	for {
		tok, err := dec.Token()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Fprintf(os.Stderr, "xmlselect: %v\n", err)
			os.Exit(1)
		}
		switch tok := tok.(type) {
		case xml.StartElement:
			tag := Tag{name: tok.Name.Local}
			for _, val := range tok.Attr {
				switch val.Name.Local {
				case "id":
					tag.id = val.Value
				case "class":
					tag.class = val.Value
				}
			}
			stack = append(stack, tag)
		case xml.EndElement:
			stack = stack[:len(stack)-1]
		case xml.CharData:
			if containsAll(stack, os.Args[1:]) {
				for _, tag := range stack {
					fmt.Printf("<%s", tag.name)
					if tag.id != "" {
						fmt.Printf(" id=%q", tag.id)
					}
					if tag.class != "" {
						fmt.Printf(" class=%q", tag.class)
					}
					fmt.Print(">")
				}
				fmt.Printf(" %s\n", tok)
			}
		}
	}
}

func containsAll(x []Tag, y []string) bool {
	for len(y) <= len(x) {
		if len(y) == 0 {
			return true
		}
		if x[0].name == y[0] || x[0].id == y[0] || x[0].class == y[0] {
			y = y[1:]
		}
		x = x[1:]
	}
	return false
}
