package router

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"webhook-tower/internal/config"
)

func TestIntegration_WebhookFlow(t *testing.T) {
	// 1. Setup config
	cfg := &config.Config{
		Routes: []config.Route{
			{
				Path:   "/deploy",
				Method: "POST",
				APIKey: "secret-123",
				Rules: []config.Rule{
					{Field: "branch", Operator: "==", Value: "main"},
				},
				Command: config.Command{
					Execute: "Write-Output 'deploying {{.branch}}'",
					Async:   false,
				},
			},
		},
	}
	router := NewRouter(cfg)

	// 2. Simulate valid webhook request
	payload := map[string]interface{}{
		"branch": "main",
	}
	body, _ := json.Marshal(payload)
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/deploy", bytes.NewBuffer(body))
	req.Header.Set("X-API-Key", "secret-123")
	req.Header.Set("Content-Type", "application/json")
	
	router.ServeHTTP(w, req)

	// 3. Verify response
	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)

	if resp["message"] != "executed" {
		t.Errorf("expected message 'executed', got %v", resp["message"])
	}

	stdout := resp["stdout"].(string)
	if !strings.Contains(stdout, "deploying main") {
		t.Errorf("expected stdout to contain 'deploying main', got '%s'", stdout)
	}
}
