package domain

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lgmontenegro/webcrawler/internal/service"
)

func TestSite_ProcessURL(t *testing.T) {

	serverOk := makeServer("/",
	`<html>
	<body>
	<a href="test">Test</a>
	</body>
	</html>`)

	tests := []struct {
		name    string
		s       *Site
		crawler Crawler
		client  *http.Client
		wantErr bool
	}{
		{
			name: "",
			s: &Site{
				SiteURL: serverOk.URL,
			},
			crawler: &service.Crawler{
				Name: "MyCrawler",
			},
			client: serverOk.Client(),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.ProcessURL(tt.crawler, tt.client); (err != nil) != tt.wantErr {
				t.Errorf("Site.ProcessURL() error = %v, wantErr %v", err, tt.wantErr)
			}

			if len(tt.s.Links) != 1 {
				t.Errorf("Site.ProcessURL() links len = %v, wantErr %v", len(tt.s.Links), 1)
			}
		})
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
