package handler

import (
	"testing"
	"net/http/httptest"
	"net/http"
	"html/template"
)

func TestAgent(t *testing.T) {
	var templ = template.Must(template.ParseGlob("../../../ui/templates/*.html"))
	agentHandler := NewAgentHandler(templ, nil)
	httprr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/agent", nil)
	if err != nil {
		t.Fatal(err)
	}
	agentHandler.AgentPage(httprr, req)
	resp := httprr.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, resp.StatusCode)
	}
}
