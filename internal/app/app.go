package app

import (
	"github.com/lgmontenegro/webcrawler/internal/controller"
	"github.com/lgmontenegro/webcrawler/internal/domain"
	"github.com/lgmontenegro/webcrawler/internal/service"
)

type App struct {
	InputURL []string
}

func (a *App) Execute() bool {
	
	neededServices, siteContent := a.setup()
	siteController := controller.Setup(neededServices, siteContent)
	siteController.ReturnContents()

	return true
}

func (a *App) setup()(neededServices service.Services, sitesContent []domain.Site){
	neededServices = Bootstrap()
	for _, url := range a.InputURL {
		sitesContent = append(sitesContent, domain.Site{SiteURL: url})		
	}

	return neededServices, sitesContent
}


