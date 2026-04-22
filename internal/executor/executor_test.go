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
