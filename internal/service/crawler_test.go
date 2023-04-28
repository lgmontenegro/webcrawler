package service

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestCrawler_CrawlerContent(t *testing.T) {
	crawlerService := &Crawler{
		Name: "MyCrawler",
	}

	anyCrawlerService := &Crawler{
		Name: "anyCrawler",
	}

	mockServer := makeServer("/", "ok")

	tests := []struct {
		name        string
		c           *Crawler
		url         string
		httpClient  *http.Client
		wantContent []byte
		wantErr     bool
	}{
		{
			name:        "address not found",
			c:           crawlerService,
			url:         "http://test",
			httpClient:  &http.Client{},
			wantContent: []byte{},
			wantErr:     true,
		},
		{
			name:        "use MyCrawler service",
			c:           crawlerService,
			url:         mockServer.URL,
			httpClient:  mockServer.Client(),
			wantContent: []byte("ok"),
			wantErr:     false,
		},
		{
			name:        "use default service",
			c:           anyCrawlerService,
			url:         mockServer.URL,
			httpClient:  mockServer.Client(),
			wantContent: []byte{},
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotContent, err := tt.c.CrawlerContent(tt.url, tt.httpClient)
			if (err != nil) != tt.wantErr {
				t.Errorf("Crawler.CrawlerContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotContent, tt.wantContent) {
				t.Errorf("Crawler.CrawlerContent() = %v, want %v", gotContent, tt.wantContent)
			}
		})
	}
}

func TestCrawler_LinkExtractor(t *testing.T) {
	tests := []struct {
		name      string
		c         *Crawler
		content   []byte
		wantLinks []string
		wantErr   bool
	}{
		{
			name: "find test link",
			c: &Crawler{
				Name: "MyCrawler",
			},
			content: []byte(`<html>
			<body>
			<a href="test">Test</a>
			</body>
			</html>`),
			wantLinks: []string{"test"},
			wantErr:   false,
		},
		{
			name: "find 2 more links",
			c: &Crawler{
				Name: "MyCrawler",
			},
			content: []byte(`<html>
			<body>
			<a href="test">Test</a>
			<a href="alpha">Beta</a>
			<a href="beta">Alpha</a>
			</body>
			</html>`),
			wantLinks: []string{"test", "alpha", "beta"},
			wantErr:   false,
		},
		{
			name: "use default service",
			c: &Crawler{
				Name: "default",
			},
			content: []byte(`<html>
			<body>
			<a href="test">Test</a>
			<a href="alpha">Beta</a>
			<a href="beta">Alpha</a>
			</body>
			</html>`),
			wantLinks: []string{},
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLinks, err := tt.c.LinkExtractor(tt.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("Crawler.LinkExtractor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotLinks, tt.wantLinks) {
				t.Errorf("Crawler.LinkExtractor() = %v, want %v", gotLinks, tt.wantLinks)
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
