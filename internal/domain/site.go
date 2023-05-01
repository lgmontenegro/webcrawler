package domain

import (
	"fmt"
	"net/http"
	"net/url"
)

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
	return s.filterLinks()
}

func (s *Site) filterLinks()(err error){
	filteredLinks := []string{}
	siteURLParsed, err := url.Parse(s.SiteURL)
	if err != nil {
		return err
	}

	for _, link := range s.Links {
		u, err := url.Parse(link)

		if err != nil {
			return err
		}

		if u.Host == "" {
			link = fmt.Sprintf("%v://%v/%v", siteURLParsed.Scheme, siteURLParsed.Host, link)
			filteredLinks = append(filteredLinks, link)
		}

		if u.Host == siteURLParsed.Host{
			filteredLinks = append(filteredLinks, link)
		}
	}

	s.Links = filteredLinks

	return nil
}