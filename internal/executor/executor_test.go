package executor

import (
	"strings"
	"testing"
)

func TestExecute_Success(t *testing.T) {
	result := Execute("Write-Output 'hello'")
	
	if result.ExitCode != 0 {
		t.Errorf("expected exit code 0, got %d", result.ExitCode)
	}
	
	if !strings.Contains(result.Stdout, "hello") {
		t.Errorf("expected stdout to contain 'hello', got '%s'", result.Stdout)
	}
}

func TestExecute_Failure(t *testing.T) {
	result := Execute("exit 1")
	
	if result.ExitCode != 1 {
		t.Errorf("expected exit code 1, got %d", result.ExitCode)
	}
}

func TestExecute_Stderr(t *testing.T) {
	// In PowerShell, writing to error stream
	result := Execute("Write-Error 'oops'")
	
	if !strings.Contains(result.Stderr, "oops") {
		t.Errorf("expected stderr to contain 'oops', got '%s'", result.Stderr)
	}
}

func TestPrepareCommand(t *testing.T) {
	tests := []struct {
		name     string
		template string
		data     map[string]interface{}
		expected string
	}{
		{
			name:     "Simple injection",
			template: "echo {{.name}}",
			data:     map[string]interface{}{"name": "world"},
			expected: "echo world",
		},
		{
			name:     "Multiple variables",
			template: "git checkout {{.branch}} && git pull {{.remote}}",
			data:     map[string]interface{}{"branch": "main", "remote": "origin"},
			expected: "git checkout main && git pull origin",
		},
		{
			name:     "Missing variable",
			template: "echo {{.missing}}",
			data:     map[string]interface{}{"name": "world"},
			expected: "echo <no value>", // text/template default
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := PrepareCommand(tt.template, tt.data)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}
