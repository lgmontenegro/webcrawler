package app

import "github.com/lgmontenegro/webcrawler/internal/service"


func Bootstrap() (services service.Services) {

	neededServices := []string{
		"crawler",
		"httpClient",
	}

	return service.Setup(neededServices)
}
