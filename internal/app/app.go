package app

import "github.com/lgmontenegro/webcrawler/internal/services"

type App struct {
	InputURL []string
}

func (a *App) Execute()(bool) {
	return services.ExecuteCrawler(a.InputURL)
}
