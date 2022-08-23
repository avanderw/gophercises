package main

import (
	"io"

	"golang.org/x/net/html"
)

type Link struct {
	Href string
	Text string
}

// Parse will take in an HTML  document and return a slice of links parsed from it.
func Parse(r io.Reader) ([]Link, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}
	return visit(nil, doc), nil
}

func visit(links []Link, n *html.Node) []Link {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				links = append(links, Link{
					Href: a.Val,
					Text: getText(n),
				})
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = visit(links, c)
	}
	return links
}

func getText(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	if n.Type != html.ElementNode {
		return ""
	}
	var ret string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ret += getText(c)
	}
	return ret
}
