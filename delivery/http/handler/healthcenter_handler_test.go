package handler

import (
	"testing"
	"net/http/httptest"
	"net/http"
	"html/template"
	"fmt"
)

func TestHealthCenters(t *testing.T) {
	var templ = template.Must(template.ParseGlob("../../../ui/templates/*.html"))
	hcHandler := NewUserHandler(templ, nil, nil)
	httprr := httptest.NewRecorder()
	url := fmt.Sprintf("/%s", "healthcenters")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}
	hcHandler.Healthcenters(httprr, req)
	resp := httprr.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, resp.StatusCode)
	}
}

func TestHealthCenterLogin(t *testing.T) {
	var templ = template.Must(template.ParseGlob("../../../ui/templates/*.html"))
	hcHandler := NewHealthCenterHandler(templ, nil, nil)
	httprr := httptest.NewRecorder()
	url := fmt.Sprintf("/%s", "healthcenter/login")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(url)
	hcHandler.HealthCenterLogin(httprr, req)
	resp := httprr.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, resp.StatusCode)
	}
}
