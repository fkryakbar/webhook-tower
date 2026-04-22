package router

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
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
			req, _ := http.NewRequest(tt.method, tt.path, strings.NewReader(tt.payload))
			for k, v := range tt.headers {
				req.Header.Set(k, v)
			}
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("%s: expected status %d, got %d", tt.name, tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestGithubWebhookSignatureVerification(t *testing.T) {
	cfg := &config.Config{
		Routes: []config.Route{
			{
				Path:                "/github",
				Method:              "POST",
				GithubWebhookSecret: "my-super-secret",
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
			name:           "Valid signature",
			method:         "POST",
			path:           "/github",
			headers:        map[string]string{"X-Hub-Signature-256": "sha256=c408c0d877912f1368eb6e74e41e159144bff3de796ca195daca3962431b5ffa"}, // payload is `{"event": "push"}` and secret is `my-super-secret`
			payload:        `{"event": "push"}`,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Missing signature",
			method:         "POST",
			path:           "/github",
			headers:        map[string]string{},
			payload:        `{"event": "push"}`,
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "Invalid signature format",
			method:         "POST",
			path:           "/github",
			headers:        map[string]string{"X-Hub-Signature-256": "17ab5a242ad4d075ba5a8edcbf0e6508d72863fb72fc823bc0c14b62eb03cb8d"},
			payload:        `{"event": "push"}`,
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "Invalid signature",
			method:         "POST",
			path:           "/github",
			headers:        map[string]string{"X-Hub-Signature-256": "sha256=wrongsignature075ba5a8edcbf0e6508d72863fb72fc823bc0c14b62eb03cb8d"},
			payload:        `{"event": "push"}`,
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
				t.Errorf("%s: expected status %d, got %d", tt.name, tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestExecutionModes(t *testing.T) {
	cfg := &config.Config{
		Routes: []config.Route{
			{
				Path:   "/sync",
				Method: "POST",
				Command: config.Command{
					Execute: "echo 'sync'",
					Async:   false,
				},
			},
			{
				Path:   "/async",
				Method: "POST",
				Command: config.Command{
					Execute: "sleep 1; echo 'async'",
					Async:   true,
				},
			},
		},
	}
	router := NewRouter(cfg)

	t.Run("Synchronous", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/sync", strings.NewReader("{}"))
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
		if !strings.Contains(w.Body.String(), "executed") {
			t.Errorf("expected body to contain 'executed', got %s", w.Body.String())
		}
	})

	t.Run("Asynchronous", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/async", strings.NewReader("{}"))
		
		start := time.Now()
		router.ServeHTTP(w, req)
		duration := time.Since(start)

		if w.Code != http.StatusAccepted {
			t.Errorf("expected status 202, got %d", w.Code)
		}
		if duration > 500*time.Millisecond {
			t.Errorf("expected async response to be fast, took %v", duration)
		}
	})
}
