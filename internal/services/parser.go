package services

import (
	"strings"

	"golang.org/x/net/html"
)

func ParseLinks(content string) (links []string, err error) {
	doc, err := html.Parse(strings.NewReader(content))
	if err != nil {
		return []string{}, err
	}

	var readNodes func(*html.Node) (links []string)
	readNodes = func(node *html.Node) (links []string) {
		if node.Type == html.ElementNode && node.Data == "a" {
			for _, link := range node.Attr {
				if link.Key == "href" {
					links = append(links, link.Val)
					break
				}
			}
		}
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			readNodes(c)
		}

		return links
	}

	links = readNodes(doc)

	return links, nil
}
