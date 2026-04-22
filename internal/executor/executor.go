package executor

import (
	"bytes"
	"os/exec"
	"strings"
	"text/template"
)

// Result represents the outcome of a command execution
type Result struct {
	Stdout   string
	Stderr   string
	ExitCode int
	Error    error
}

// PrepareCommand injects data into a command template
func PrepareCommand(tmplStr string, data interface{}) (string, error) {
	tmpl, err := template.New("command").Parse(tmplStr)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// Execute runs a shell command and returns the result
func Execute(command string) *Result {
	// Simple shell execution (using sh -c or cmd /c based on OS)
	// For now, let's assume a shell is available.
	// We'll use "powershell.exe -Command" for Windows since we are on win32.
	
	cmd := exec.Command("powershell.exe", "-Command", command)
	
	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	
	err := cmd.Run()
	
	exitCode := 0
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode = exitError.ExitCode()
		} else {
			exitCode = -1 // System error
		}
	}

	return &Result{
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
		ExitCode: exitCode,
		Error:    err,
	}
}
