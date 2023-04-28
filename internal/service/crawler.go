package service

import (
	"net/http"

	"github.com/lgmontenegro/webcrawler/internal/service/myCrawler"
)

type Crawler struct {
	Name string
}

func (c *Crawler) CrawlerContent(url string, client *http.Client) (content []byte, err error) {
	switch c.Name {
	case "MyCrawler":
		crawler := myCrawler.MyCrawler{
			Client: client,
		}
		return crawler.CrawlURL(url)
	default:
		return []byte{}, nil
	}
}

func (c *Crawler) LinkExtractor(content []byte) (links []string, err error) {
	switch c.Name {
	case "MyCrawler":
		crawler := myCrawler.MyCrawler{}
		return crawler.FindLinks(content)
	default:
		return []string{}, nil
	}
}