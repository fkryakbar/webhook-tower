package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	yamlContent := `
routes:
  - path: "/deploy"
    method: "POST"
    api_key: "secret123"
    headers:
      - key: "Content-Type"
        value: "application/json"
    rules:
      - field: "payload.event"
        operator: "=="
        value: "push"
    command:
      execute: "ls -la"
      async: false
`

	cfg, err := LoadConfig(yamlContent)
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	if len(cfg.Routes) != 1 {
		t.Errorf("expected 1 route, got %d", len(cfg.Routes))
	}

	route := cfg.Routes[0]
	if route.Path != "/deploy" {
		t.Errorf("expected path /deploy, got %s", route.Path)
	}

	if route.Method != "POST" {
		t.Errorf("expected method POST, got %s", route.Method)
	}

	if route.APIKey != "secret123" {
		t.Errorf("expected api_key secret123, got %s", route.APIKey)
	}

	if route.Command.Execute != "ls -la" {
		t.Errorf("expected command ls -la, got %s", route.Command.Execute)
	}
}

func TestLoadConfig_InvalidYAML(t *testing.T) {
	invalidYAML := `
routes:
  - path: "/deploy"
    invalid_key: :
`
	_, err := LoadConfig(invalidYAML)
	if err == nil {
		t.Error("expected error for invalid YAML, got nil")
	}
}

func TestLoadConfigFromFile(t *testing.T) {
	yamlContent := `
routes:
  - path: "/deploy"
    method: "POST"
`
	tmpFile := "config_test.yaml"
	if err := os.WriteFile(tmpFile, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("failed to write tmp file: %v", err)
	}
	defer os.Remove(tmpFile)

	cfg, err := LoadConfigFromFile(tmpFile)
	if err != nil {
		t.Fatalf("failed to load config from file: %v", err)
	}

	if cfg.Routes[0].Path != "/deploy" {
		t.Errorf("expected path /deploy, got %s", cfg.Routes[0].Path)
	}
}
