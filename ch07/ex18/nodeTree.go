package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
)

type Node interface{ String() string } // CharData OR *Element

type CharData string

func (c CharData) String() string {
	return string(c)
}

type Element struct {
	Type     xml.Name
	Attr     []xml.Attr
	Children []Node
}

func (e Element) String() string {
	childrenStr := ""
	for _, child := range e.Children {
		childrenStr += child.String()
	}
	return fmt.Sprintf("<%s>%s</%s>", e.Type.Local, childrenStr, e.Type.Local)
}

func main() {
	dec := xml.NewDecoder(os.Stdin)

	var makeNodeTree func() []Node
	makeNodeTree = func() (nodes []Node) {
		tok, err := dec.Token()
		for err == nil {
			switch tok := tok.(type) {
			case xml.StartElement:
				child := &Element{tok.Name, tok.Attr, makeNodeTree()}
				fmt.Printf("parent:   %v\n", child)
				nodes = append(nodes, Node(child))
			case xml.EndElement:
				fmt.Printf("children: %v\n", nodes)
				return
			case xml.CharData:
				nodes = append(nodes, Node(CharData(string(tok))))
			}
			tok, err = dec.Token()
		}
		if err != io.EOF {
			log.Fatal(err)
		}
		fmt.Printf("children: %v\n", nodes)
		return
	}

	roots := makeNodeTree()
	s := ""
	for _, root := range roots {
		s += root.String()
	}
	fmt.Printf("root:     %s\n", s)
}
