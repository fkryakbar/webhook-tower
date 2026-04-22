package router

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthCheck(t *testing.T) {
	router := NewRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	expected := `{"status":"up"}`
	if w.Body.String() != expected {
		t.Errorf("expected body %s, got %s", expected, w.Body.String())
	}
}
