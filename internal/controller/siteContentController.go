package controller

import (
	"net/http"
	"sync"

	"github.com/lgmontenegro/webcrawler/internal/domain"
	srvc "github.com/lgmontenegro/webcrawler/internal/service"
)

type SiteContentController struct {
	Service map[string]any
	Sites   []*domain.Site
}

func Setup(services srvc.Services, sitesContent []*domain.Site) (siteController SiteContentController) {
	siteController.Service = make(map[string]any, 1)
	for _, service := range services.Installed {
		for typeService, s := range service {
			siteController.Service[typeService] = s
		}
	}

	siteController.Sites = sitesContent

	return siteController
}

func (s *SiteContentController) Process() (errors []error) {
	crawlerService, _ := s.Service["crawler"].(srvc.Crawler)
	httpClient, _ := s.Service["httpClient"].(http.Client)

	var wg sync.WaitGroup
	errorCh := make(chan error)

	for _, site := range s.Sites {
		wg.Add(1)

		go func(site *domain.Site) {
			defer wg.Done()
			errorCh <- site.ProcessURL(&crawlerService, &httpClient)
		}(site)
		errors = append(errors, <-errorCh)
	}

	wg.Wait()

	return errors
}

func (s *SiteContentController) GetLinks() (links map[string][]string) {
	links = make(map[string][]string, 0)

	for _, siteCrawled := range s.Sites {
		links[siteCrawled.SiteURL] = siteCrawled.Links
	}

	return links
}
