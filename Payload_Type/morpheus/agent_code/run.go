//go:build windows
// +build windows

package main

import (
	"encoding/json"
	"os/exec"
	"strings"
	"syscall"
)

type RunParams struct {
	Path string   `json:"path"`
	Args []string `json:"args"`
}

func cmdRun(raw json.RawMessage) string {
	var p RunParams
	if err := json.Unmarshal(raw, &p); err != nil {
		return "error: " + err.Error()
	}

	cmd := exec.Command(p.Path, p.Args...)
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
