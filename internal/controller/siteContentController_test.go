package controller

import (
	"fmt"
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

func TestSiteContentController_Process(t *testing.T) {
	tests := []struct {
		name          string
		content       string
		linksQuantity int
		wantErr       bool
	}{
		{
			name: "one link test",
			content: `<html>
			<body>
			<a href="test">Test</a>
			</body>
			</html>`,
			linksQuantity: 1,
			wantErr:       false,
		},
		{
			name: "three link test",
			content: `<html>
			<body>
			<a href="test">Test</a>
			<a href="alpha">Beta</a>
			<a href="beta">Alpha</a>
			</body>
			</html>`,
			linksQuantity: 3,
			wantErr:       false,
		},
		{
			name: "wrong link",
			content: `<html>
			<body>
			<a href="test">Test</a>
			<a href="alpha">Beta</a>
			<a href="beta">Alpha</a>
			</body>
			</html>`,
			linksQuantity: 0,
			wantErr:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := makeServer("/", tt.content)
			services := service.Setup([]string{"crawler", "httpClient"})
			sitesContent := make([]*domain.Site, 0)
			url := mockServer.URL
			if tt.wantErr {
				url = "http://localhost"
			}
			sitesContent = append(sitesContent, &domain.Site{SiteURL: url})
			gotSiteController := Setup(services, sitesContent)
			gotSiteController.Service["httpClient"] = *mockServer.Client()

			errors := gotSiteController.Process()

			if tt.wantErr {
				if len(errors) == 0 {
					t.Errorf("got %v, want %v", len(errors), ">0")
				}
			}

			for _, site := range gotSiteController.Sites {
				if len(site.Links) != tt.linksQuantity {
					t.Errorf("got %v, want %v", len(site.Links), tt.linksQuantity)
				}
			}
		})
	}
}

func TestSiteContentController_GetLinks(t *testing.T) {

	tests := []struct {
		name          string
		content       string
		linksQuantity int
		wantLinks     []string
	}{
		{
			name: "one link test",
			content: `<html>
			<body>
			<a href="test">Test</a>
			</body>
			</html>`,
			linksQuantity: 1,
			wantLinks: []string{
				"test",
			},
		},
		{
			name: "three link test",
			content: `<html>
			<body>
			<a href="test">Test</a>
			<a href="alpha">Beta</a>
			<a href="beta">Alpha</a>
			</body>
			</html>`,
			linksQuantity: 3,
			wantLinks: []string{
				"test",
				"alpha",
				"beta",
			},
		},
		{
			name: "no link test",
			content: `<html>
			<body>
			<a href="http://nonexiste/test">Test</a>
			<a href="http://nonexiste/alpha">Beta</a>
			<a href="http://nonexiste/beta">Alpha</a>
			</body>
			</html>`,
			linksQuantity: 0,
			wantLinks: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := makeServer("/", tt.content)
			services := service.Setup([]string{"crawler", "httpClient"})
			sitesContent := make([]*domain.Site, 0)
			sitesContent = append(sitesContent, &domain.Site{SiteURL: mockServer.URL})
			gotSiteController := Setup(services, sitesContent)
			gotSiteController.Service["httpClient"] = *mockServer.Client()

			gotSiteController.Process()

			links := gotSiteController.GetLinks()

			if len(links[mockServer.URL]) != tt.linksQuantity {
				t.Errorf("got link quantity %v, want %v", len(links[mockServer.URL]), tt.linksQuantity)
			}

			for i, wantLink := range tt.wantLinks {
				if links[mockServer.URL][i] != fmt.Sprintf("%v/%v", mockServer.URL, wantLink) {
					t.Errorf("got link %v, want %v", links[mockServer.URL][i], fmt.Sprintf("%v/%v", mockServer.URL, wantLink))
				}
			}

		})
	}

	t.Run("multiples urls", func(t *testing.T) {
		mockServer1 := makeServer("/", tests[0].content)
		mockServer2 := makeServer("/", tests[1].content)

		services := service.Setup([]string{"crawler", "httpClient"})
		sitesContent := make([]*domain.Site, 0)
		sitesContent = append(sitesContent, &domain.Site{SiteURL: mockServer1.URL})
		sitesContent = append(sitesContent, &domain.Site{SiteURL: mockServer2.URL})
		gotSiteController := Setup(services, sitesContent)
		gotSiteController.Service["httpClient"] = http.Client{}

		gotSiteController.Process()

		links := gotSiteController.GetLinks()

		for serverURL, link := range links {
			switch serverURL {
			case mockServer1.URL:
				for i, wantLink := range tests[0].wantLinks {
					if link[i] != fmt.Sprintf("%v/%v", mockServer1.URL, wantLink) {
						t.Errorf("got link %v, want %v", links[mockServer1.URL][i], fmt.Sprintf("%v/%v", mockServer1.URL, wantLink))
					}
				}
			case mockServer2.URL:
				for i, wantLink := range tests[1].wantLinks {
					if link[i] != fmt.Sprintf("%v/%v", mockServer2.URL, wantLink) {
						t.Errorf("got link %v, want %v", links[mockServer1.URL][i], fmt.Sprintf("%v/%v", mockServer2.URL, wantLink))
					}
				}
			default:
			}
		}

	})
}

func makeServer(path, body string) (server httptest.Server) {
	mux := http.NewServeMux()
	mux.HandleFunc(path, func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte(body))
	})

	server = *httptest.NewServer(mux)

	return server
}