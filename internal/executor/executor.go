package executor

import (
	"bytes"
	"os/exec"
	"runtime"
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
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("powershell.exe", "-Command", command)
	} else {
		cmd = exec.Command("sh", "-c", command)
	}

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
