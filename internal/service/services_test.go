package service

import (
	"net/http"
	"reflect"
	"testing"
)

func TestSetup(t *testing.T) {
	serviceInstalled := Crawler{
		Name: "MyCrawler",
	}
	
	tests := []struct {
		name           string
		neededServices []string
		wantServices   map[string]any
	}{
		{
			name:           "get crawler",
			neededServices: []string{"crawler"},
			wantServices:   map[string]any{"crawler": serviceInstalled},
		},
		{
			name:           "get http client",
			neededServices: []string{"httpClient"},
			wantServices:   map[string]any{"httpClient": http.Client{}},
		},
		{
			name:           "default services",
			neededServices: []string{""},
			wantServices:   map[string]any{"": nil},
		},
	}
	for _, tt := range tests {
		services := Services{}
		services.Installed = append(services.Installed, tt.wantServices)

		t.Run(tt.name, func(t *testing.T) {
			if gotServices := Setup(tt.neededServices); !reflect.DeepEqual(gotServices, services) {
				t.Errorf("Setup() = %v, want %v", gotServices, services)
			}
		})
	}
}
