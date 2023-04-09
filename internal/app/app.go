package app

import "github.com/lgmontenegro/webcrawler/internal/services"

type App struct {
	InputURL []string
}

func (a *App) Execute() {
	crawler := services.Crawler{}
	crawler.Execute(a.InputURL)
}
