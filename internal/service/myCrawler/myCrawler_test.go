package myCrawler

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestMyCrawler_CrawlURL(t *testing.T) {
	serverOk := makeServer("/", "ok")

	tests := []struct {
		testName  string
		myCrawler *MyCrawler
		url       string
		wantBody  []byte
		wantErr   bool
	}{
		{
			testName: "ok",
			myCrawler: &MyCrawler{
				Client: serverOk.Client(),
			},
			url:      serverOk.URL,
			wantBody: []byte("ok"),
			wantErr:  false,
		},
		{
			testName: "wrong address",
			myCrawler: &MyCrawler{
				Client: serverOk.Client(),
			},
			url:      "http://nonexist",
			wantBody: []byte{},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			gotBody, err := tt.myCrawler.CrawlURL(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCrawler.CrawlURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotBody, tt.wantBody) {
				t.Errorf("MyCrawler.CrawlURL() = %v, want %v", gotBody, tt.wantBody)
			}
		})
	}
}

func TestMyCrawler_FindLinks(t *testing.T) {
	tests := []struct {
		name      string
		m         *MyCrawler
		content   []byte
		wantLinks []string
		wantErr   bool
	}{
		{
			name: "find test link",
			m:    &MyCrawler{},
			content: []byte(`<html>
			<body>
			<a href="test">Test</a>
			</body>
			</html>`),
			wantLinks: []string{"test"},
			wantErr: false,
		},
		{
			name: "find 2 more links",
			m:    &MyCrawler{},
			content: []byte(`<html>
			<body>
			<a href="test">Test</a>
			<a href="alpha">Beta</a>
			<a href="beta">Alpha</a>
			</body>
			</html>`),
			wantLinks: []string{"test", "alpha", "beta"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLinks, err := tt.m.FindLinks(tt.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCrawler.FindLinks() error = %v, wantErr %v", err, tt.wantErr)
				return			
			}
			
			if len(gotLinks) != len(tt.wantLinks) {
				t.Errorf("MyCrawler.FindLinks() = %v, want %v", gotLinks, tt.wantLinks)
			}

			for i, links := range gotLinks{
				if links != tt.wantLinks[i] {
					t.Errorf("Different values and positions: %v want %v", links, tt.wantLinks[i])	
				}
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
