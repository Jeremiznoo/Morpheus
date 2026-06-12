//go:build windows
// +build windows

package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type MkdirParams struct {
	Path string `json:"path"`
}

func cmdMkdir(raw json.RawMessage) string {
	var p MkdirParams
	if err := json.Unmarshal(raw, &p); err != nil {
		return "error: " + err.Error()
	}

	if err := os.MkdirAll(p.Path, 0755); err != nil {
		return "error: " + err.Error()
	}

	return fmt.Sprintf("created directory: %s", p.Path)
}
