package services

import (
	"fmt"
	"io"
	"net/http"
	"sync"
)

func ExecuteCrawler(urls []string) bool {
	return crawlUrls(urls)
}

func crawlUrls(urls []string) bool {
	urlChannel := make(chan string, len(urls))
	urlContent := make(chan string, len(urls))
	parsedLinks := make(chan []string, len(urls))
	errorChannel := make(chan error, len(urls))

	var wg sync.WaitGroup
	var m sync.Mutex

	for _, url := range urls {
		wg.Add(1)
		urlChannel <- url

		go func() {
			fetchURL(urlChannel, urlContent, errorChannel, &wg, &m)
			parseLinks(<- urlContent, errorChannel, parsedLinks, &wg)
		}()
	}
	wg.Wait()

	showErrors(errorChannel)

	fmt.Println("link total:", len(parsedLinks))
	if len(parsedLinks) > 0 {		
		for len(parsedLinks) >0  {
			links := <- parsedLinks
			fmt.Println(len(links))
			for i, link := range links {
				fmt.Println(i, link)
			}
			
		}
	}
		
	return true
}

func fetchURL(
	urlChannel chan string,
	urlContent chan string,
	errorChannel chan error,	
	wg *sync.WaitGroup,
	m *sync.Mutex,
) {
	m.Lock()
	req, err := http.NewRequest("GET", <- urlChannel, nil)	
	if err != nil {
		urlContent <- ""
		errorChannel <- err
		defer wg.Done()
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		urlContent <- ""
		errorChannel <- err
		defer wg.Done()
		return
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		urlContent <- ""
		errorChannel <- err
		defer wg.Done()
		return
	}

	urlContent <- string(body)
	m.Unlock()
}

func parseLinks(	
	urlContent string,
	errorChannel chan error,
	parsedLinks chan []string,
	wg *sync.WaitGroup,
) {	
	body := urlContent
	links, err := ParseLinks(body)
	if err != nil {
		parsedLinks <- []string{}
		errorChannel <- err
		defer wg.Done()
		return
	}

	parsedLinks <- links
	errorChannel <-err
	defer wg.Done()
}

func showErrors(errorChannel chan error) {
	if len(errorChannel) > 0 {
		for len(errorChannel) >0  {
			if err := <- errorChannel; err != nil {
				fmt.Println("erro:", <- errorChannel)
			}
			
		}
	}
}