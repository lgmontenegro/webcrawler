package service

import "net/http"

type Services struct {
	Installed []map[string]any
}

func Setup(neededServices []string) (services Services) {

	for _, service := range neededServices {
		switch service {
		case "crawler":
			serviceInstalled := Crawler{
				Name: "MyCrawler",
			}
			services.Installed = append(services.Installed, map[string]any{service: serviceInstalled})
		case "httpClient":
			serviceInstalled := http.Client{}
			services.Installed = append(services.Installed, map[string]any{service: serviceInstalled})
		default:
			services.Installed = append(services.Installed, map[string]any{"": nil})
			return services
		}
	}

	return services
}
