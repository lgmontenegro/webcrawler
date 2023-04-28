package domain

import "net/http"

type Site struct {
	SiteURL string
	Content []byte
	Links []string
}

type Crawler interface {
	CrawlerContent(url string, client *http.Client)(content []byte, err error)
	LinkExtractor(content []byte)(links []string, err error)
}

func (s *Site) ProcessURL(crawler Crawler, client *http.Client)(err error){
	content, err := crawler.CrawlerContent(s.SiteURL, client)
	if err != nil {
		return err
	}
	s.Content = content

	links, err := crawler.LinkExtractor(s.Content)
	if err != nil {
		return err
	}
	s.Links = links

	return nil
}