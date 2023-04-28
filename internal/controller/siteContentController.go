package controller

import (
	"net/http"

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

func (s *SiteContentController) ReturnContents() {	
	crawlerService, _ := s.Service["crawler"].(srvc.Crawler)
	httpClient, _ := s.Service["httpClient"].(http.Client)

	for _, site := range s.Sites {
		site.ProcessURL(&crawlerService, &httpClient)
	}
}
