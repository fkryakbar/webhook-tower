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
				APIKey: "secret-key",
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
		queryParams    string
		payload        string
		expectedStatus int
	}{
		{
			name:           "Valid request with header key",
			method:         "POST",
			path:           "/webhook",
			headers:        map[string]string{"X-Test": "passed", "X-API-Key": "secret-key"},
			payload:        `{"event": "push"}`,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Valid request with query key",
			method:         "POST",
			path:           "/webhook?api_key=secret-key",
			headers:        map[string]string{"X-Test": "passed"},
			payload:        `{"event": "push"}`,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid API key",
			method:         "POST",
			path:           "/webhook",
			headers:        map[string]string{"X-Test": "passed", "X-API-Key": "wrong-key"},
			payload:        `{"event": "push"}`,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Missing API key",
			method:         "POST",
			path:           "/webhook",
			headers:        map[string]string{"X-Test": "passed"},
			payload:        `{"event": "push"}`,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Payload rule mismatch",
			method:         "POST",
			path:           "/webhook",
			headers:        map[string]string{"X-Test": "passed", "X-API-Key": "secret-key"},
			payload:        `{"event": "pull_request"}`,
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "Field not found in payload",
			method:         "POST",
			path:           "/webhook",
			headers:        map[string]string{"X-Test": "passed", "X-API-Key": "secret-key"},
			payload:        `{"wrong_field": "push"}`,
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			path := tt.path
			req, _ := http.NewRequest(tt.method, path, strings.NewReader(tt.payload))
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
