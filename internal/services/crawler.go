package services

import (
	"fmt"
	"io"
	"net/http"

	"github.com/lgmontenegro/webcrawler/internal/domain"
)

type Crawler struct {
	Job []domain.SiteContent
}

func (c *Crawler) Execute(urls []string) {
	//c.setup(urls)
	c.crawlUrls(urls)

	for _, t := range c.Job {
		fmt.Println("URL: ", t.SiteURL)
		fmt.Println("content: ", string(t.Content))
	}
}
/*
func (c *Crawler) setup(urls []string) {
	for _, url := range urls {
		c.Job = append(c.Job, domain.SiteContent{
			SiteURL: url,
			Content: []byte{},
		})
	}
}*/

func (c *Crawler) crawlUrls(urls []string) {
	/*urlChannel := make(chan string, len(c.Job))
	contentToRet := make(chan []byte, len(c.Job))
	errorChannel := make(chan error, len(c.Job))*/

	for i, job := range c.Job {
		urlChannel <- job.SiteURL
		go func() {
			urlChannelString := <-urlChannel
			content, err := fetchURL(urlChannelString)

			contentToRet <- content
			errorChannel <- err
		}()

		err := <-errorChannel
		if err != nil {
			fmt.Println(err)
		}

		contentToRetByte := <-contentToRet
		c.Job[i].Content = contentToRetByte
	}
}

func fetchURL(url string) (body []byte, err error) {

	resp, err := http.Get(url)
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
