package router

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"webhook-tower/internal/config"
)

func TestHealthCheck(t *testing.T) {
	cfg := &config.Config{}
	router := NewRouter(cfg)

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

func TestRouteMatching(t *testing.T) {
	cfg := &config.Config{
		Routes: []config.Route{
			{
				Path:   "/webhook",
				Method: "POST",
				Headers: []config.Header{
					{Key: "X-Test", Value: "passed"},
				},
				Rules: []config.Rule{
					{Field: "event", Operator: "==", Value: "push"},
				},
			},
		},
	}
	router := NewRouter(cfg)

	tests := []struct {
		name           string
		method         string
		path           string
		headers        map[string]string
		payload        string
		expectedStatus int
	}{
		{
			name:           "Valid request",
			method:         "POST",
			path:           "/webhook",
			headers:        map[string]string{"X-Test": "passed"},
			payload:        `{"event": "push"}`,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid payload",
			method:         "POST",
			path:           "/webhook",
			headers:        map[string]string{"X-Test": "passed"},
			payload:        `{"event": "pull_request"}`,
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "Invalid method",
			method:         "GET",
			path:           "/webhook",
			headers:        map[string]string{"X-Test": "passed"},
			payload:        `{"event": "push"}`,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Field not found",
			method:         "POST",
			path:           "/webhook",
			headers:        map[string]string{"X-Test": "passed"},
			payload:        `{"wrong_field": "push"}`,
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tt.method, tt.path, strings.NewReader(tt.payload))
			for k, v := range tt.headers {
				req.Header.Set(k, v)
			}
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}
