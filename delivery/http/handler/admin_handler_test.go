package handler

import (
	"testing"
	"net/http/httptest"
	"net/http"
	"html/template"
)

func TestAdminPage(t *testing.T) {
	var templ = template.Must(template.ParseGlob("../../../ui/templates/*.html"))
	adminHandler := NewAdminHandler(templ,nil,nil)
	httprr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/admin", nil)
	if err != nil {
		t.Fatal(err)
	}
	adminHandler.AdminPage(httprr, req)
	resp := httprr.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, resp.StatusCode)
	}
}
