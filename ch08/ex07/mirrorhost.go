package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"golang.org/x/net/html"
)

var targetHost string

type Link struct {
	tag string
	url *url.URL
}

func getSavePath(uUrl *url.URL) string {
	savepath := filepath.Join("data", uUrl.Host, uUrl.Path)
	if filepath.Ext(uUrl.Path) == "" {
		savepath = filepath.Join(savepath, "index.html")
	}
	return savepath
}

func forEachNode(n *html.Node, pre, post func(n *html.Node)) {
	if pre != nil {
		pre(n)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachNode(c, pre, post)
	}
	if post != nil {
		post(n)
	}
}

func linksExtractAndRewrite(uUrl *url.URL) ([]Link, *html.Node, error) {
	resp, err := http.Get(uUrl.String())
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("getting %s: %s", uUrl, resp.Status)
	}
	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	basepath := filepath.Join("data", uUrl.Host, uUrl.Path)

	rewriteAttr := func(a html.Attribute) (html.Attribute, *url.URL, error) {
		linkUrl, err := resp.Request.URL.Parse(a.Val)
		if err != nil {
			return html.Attribute{}, nil, err
		}
		savepath := getSavePath(linkUrl)
		a.Val, err = filepath.Rel(basepath, savepath)
		if err != nil {
			return html.Attribute{}, nil, err
		}
		return a, linkUrl, nil
	}

	var links []Link
	visitAndRewriteNode := func(n *html.Node) {
		if n.Type != html.ElementNode {
			return
		}
		var attrKey string
		switch n.Data {
		case "a", "link", "base", "area":
			attrKey = "href"
		case "img", "script", "iframe", "embed":
			attrKey = "src"
		default:
			return
		}
		for i, a := range n.Attr {
			if a.Key != attrKey {
				continue
			}
			rewritedAttr, linkUrl, err := rewriteAttr(a)
			if err != nil {
				continue
			}
			n.Attr[i] = rewritedAttr
			links = append(links, Link{n.Data, linkUrl})
		}
	}

	forEachNode(doc, visitAndRewriteNode, nil)
	return links, doc, nil
}

func save(data io.Reader, uUrl *url.URL) error {
	savepath := getSavePath(uUrl)
	if err := os.MkdirAll(filepath.Dir(savepath), 0755); err != nil {
		return err
	}
	file, err := os.Create(savepath)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err = io.Copy(file, data); err != nil {
		return err
	}
	return nil
}

func saveAndCrawl(link Link) []Link {
	fmt.Println(link.url)
	links, doc, err := linksExtractAndRewrite(link.url)
	if err != nil {
		log.Print(err)
		return nil
	}
	htmlbuf := &bytes.Buffer{}
	html.Render(htmlbuf, doc)
	if err := save(htmlbuf, link.url); err != nil {
		log.Print(err)
	}
	return links
}

func main() {
	worklist := make(chan []Link)
	unseenLinks := make(chan Link)

	var links []Link
	for _, sUrl := range os.Args[1:] {
		uUrl, err := url.Parse(sUrl)
		if err != nil {
			continue
		}
		links = append(links, Link{"a", uUrl})
	}
	go func() { worklist <- links }()

	for i := 0; i < 20; i++ {
		go func() {
			for link := range unseenLinks {
				if targetHost == "" {
					targetHost = link.url.Host
				}
				if link.url.Host != targetHost {
					continue
				}
				if link.tag == "a" {
					foundLinks := saveAndCrawl(link)
					go func() { worklist <- foundLinks }()
				} else {
					resp, err := http.Get(link.url.String())
					if err != nil {
						log.Print(err)
						continue
					}
					defer resp.Body.Close()
					if resp.StatusCode != http.StatusOK {
						log.Print(fmt.Errorf("getting %s: %s", link.url, resp.Status))
						continue
					}
					if err = save(resp.Body, link.url); err != nil {
						log.Print(err)
					}
				}
			}
		}()
	}

	seen := make(map[string]bool)
	for list := range worklist {
		for _, link := range list {
			if !seen[link.url.String()] {
				seen[link.url.String()] = true
				unseenLinks <- link
			}
		}
	}
}
