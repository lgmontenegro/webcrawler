package services

import (
	"fmt"
	"io"
	"net/http"
	//"sync"
)

func ExecuteCrawler(urls []string) bool {
	return crawlUrls(urls)
}

func crawlUrls(urls []string) bool {
	urlChannel := make(chan string, len(urls))	
	parsedLinks := make(chan []string, len(urls))
	errorChannel := make(chan error, len(urls))

	//var wg sync.WaitGroup	

	
	for _, url := range urls {		
		//wg.Add(1)
		urlChannel <- url

		go func() {
			links := fetchURL(urlChannel, errorChannel)
				//, &wg)
			parsedLinks <-links
		}()
	}
	//wg.Wait()

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
	errorChannel chan error,	
	//wg *sync.WaitGroup,	
) (urlContent []string) {
	req, err := http.NewRequest("GET", <- urlChannel, nil)	
	if err != nil {
		urlContent = []string{}
		errorChannel <- err
		//defer wg.Done()
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {		
		errorChannel <- err
		//defer wg.Done()
		return []string{}
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {		
		errorChannel <- err
		//defer wg.Done()
		return []string{}
	}

	parsedLinks, err := ParseLinks(string(body))
	if err != nil {		
		errorChannel <- err
		return []string{}
	}

	//defer wg.Done()
	return parsedLinks
}

/*func parseLinks(	
	urlContent string,
	errorChannel chan error,		
)(parsedLinks []string) {	
	parsedLinks, err := ParseLinks(urlContent)
	if err != nil {		
		errorChannel <- err
		return []string{}
	}

	errorChannel <-err
	return parsedLinks
}*/

func showErrors(errorChannel chan error) {
	if len(errorChannel) > 0 {
		for len(errorChannel) >0  {
			if err := <- errorChannel; err != nil {
				fmt.Println("erro:", <- errorChannel)
			}
			
		}
	}
}