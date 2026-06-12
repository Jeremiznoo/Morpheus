//go:build windows
// +build windows

package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type RmParams struct {
	Path string `json:"path"`
}

func cmdRm(raw json.RawMessage) string {
	var p RmParams
	if err := json.Unmarshal(raw, &p); err != nil {
		return "error: " + err.Error()
	}

	if err := os.RemoveAll(p.Path); err != nil {
		return "error: " + err.Error()
	}

	return fmt.Sprintf("removed: %s", p.Path)
}
