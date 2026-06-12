//go:build windows
// +build windows

package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type CopyParams struct {
	Source string `json:"source"`
	Dest   string `json:"dest"`
}

func cmdCp(raw json.RawMessage) string {
	var p CopyParams
	if err := json.Unmarshal(raw, &p); err != nil {
		return "error: " + err.Error()
	}

	input, err := os.ReadFile(p.Source)
	if err != nil {
		return "error: " + err.Error()
	}

	if err := os.WriteFile(p.Dest, input, 0644); err != nil {
		return "error: " + err.Error()
	}

	return fmt.Sprintf("copied %s -> %s (%d bytes)", p.Source, p.Dest, len(input))
}
