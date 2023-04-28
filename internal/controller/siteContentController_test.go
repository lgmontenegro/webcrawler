package controller

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lgmontenegro/webcrawler/internal/domain"
	"github.com/lgmontenegro/webcrawler/internal/service"
)

func TestSetup(t *testing.T) {
	t.Run("setup a site content controller", func(t *testing.T) {

		services := service.Setup([]string{"crawler", "httpClient"})
		sitesContent := make([]*domain.Site, 1)
		sitesContent = append(sitesContent, &domain.Site{SiteURL: "url"})
		gotSiteController := Setup(services, sitesContent)

		if gotSiteController.Service["crawler"].(service.Crawler).Name != "MyCrawler" {
			t.Errorf("got service %v, want %v", gotSiteController.Service["crawler"].(service.Crawler).Name, "MyCrawler")
		}
	})
}

func TestSiteContentController_ReturnContents(t *testing.T) {
	mockServer := makeServer("/", `<html>
	<body>
	<a href="test">Test</a>
	</body>
	</html>`)
	services := service.Setup([]string{"crawler", "httpClient"})
	sitesContent := make([]*domain.Site,0)
	sitesContent = append(sitesContent, &domain.Site{SiteURL: mockServer.URL})
	gotSiteController := Setup(services, sitesContent)
	gotSiteController.Service["httpClient"] = *mockServer.Client()

	gotSiteController.ReturnContents()

	for _, site := range gotSiteController.Sites {
		if len(site.Links) != 1 {
			t.Errorf("got %v, want %v", len(site.Links), 1)
		}
	}
}

func makeServer(path, body string) (server httptest.Server) {
	mux := http.NewServeMux()
	mux.HandleFunc(path, func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte(body))
	})

	server = *httptest.NewServer(mux)

	return server
}
