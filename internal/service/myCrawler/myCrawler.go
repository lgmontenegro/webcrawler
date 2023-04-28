package myCrawler

import (
	"io"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

type MyCrawler struct {
	Client *http.Client
}

func (m *MyCrawler) CrawlURL(url string) (body []byte, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []byte{}, err
	}

	resp, err := m.Client.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	return body, nil
}

func (m *MyCrawler) FindLinks(content []byte) (links []string, err error) {

	doc, err := html.Parse(strings.NewReader(string(content)))
	if err != nil {
		return []string{}, err
	}

	var readNodes func(*html.Node) (links []string)
	readNodes = func(node *html.Node) (link []string) {
		if node.Type == html.ElementNode && node.Data == "a" {
			for _, link := range node.Attr {
				if link.Key == "href" {
					return []string{link.Val}
				}
			}
		}

		for c := node.FirstChild; c != nil; c = c.NextSibling {
			link = append(link, readNodes(c)...)
		}

		return link
	}

	links = readNodes(doc)
	return links, nil
}
