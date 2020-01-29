package handler

import (
	"testing"
	"net/http/httptest"
	"net/http"
	"html/template"
)

// TestAbout tests GET /about request handler
func TestUserPage(t *testing.T) {
	var templ = template.Must(template.ParseGlob("../../../ui/templates/*.html"))
	userHandler := NewUserHandler(templ, nil,nil )
	httprr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	userHandler.Index(httprr, req)
	resp := httprr.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, resp.StatusCode)
	}
}
