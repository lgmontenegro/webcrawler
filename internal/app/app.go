package app

import (
	"fmt"

	"github.com/lgmontenegro/webcrawler/internal/controller"
	"github.com/lgmontenegro/webcrawler/internal/domain"
	"github.com/lgmontenegro/webcrawler/internal/service"
)

type App struct {
	InputURL []string
}

func (a *App) Execute() bool {

	neededServices, siteContent := a.setup()
	siteContentController := controller.Setup(neededServices, siteContent)
	errors := siteContentController.Process()

	if len(errors) > 0 {
		fmt.Println("Some errors found: ")
		for _, err := range errors {
			if err != nil {
				fmt.Println(err)
			}
		}
		fmt.Println("-----------")
		fmt.Println()
	}

	links := siteContentController.GetLinks()
	if len(links) > 0 {
		for site, linksFound := range links {
			fmt.Println("Links found for", site, ":")
			for _, linkFound :=  range linksFound {
				fmt.Println(linkFound)
			}
			fmt.Println("-----------")
			fmt.Println()
		}		
	}	

	return true
}

func (a *App) setup() (neededServices service.Services, sitesContent []*domain.Site) {
	neededServices = Bootstrap()
	for _, url := range a.InputURL {
		sitesContent = append(sitesContent, &domain.Site{SiteURL: url})
	}

	return neededServices, sitesContent
}
