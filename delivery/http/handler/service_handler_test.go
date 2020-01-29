package handler

import (
	"net/http/httptest"
	"net/http"
	"testing"
	"html/template"
	"fmt"
)

func TestService(t *testing.T) {
	var templ = template.Must(template.ParseGlob("../../../ui/templates/*.html"))
	serviceHandler := NewServiceHandler(templ)
	httprr := httptest.NewRecorder()
	url := fmt.Sprintf("/%s/%s/%s","healthcenter","service","editservice")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(url)
	serviceHandler.EditService(httprr, req)
	resp := httprr.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, resp.StatusCode)
	}
}