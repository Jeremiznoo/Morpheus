//go:build windows
// +build windows

package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type MoveParams struct {
	Source string `json:"source"`
	Dest   string `json:"dest"`
}

func cmdMv(raw json.RawMessage) string {
	var p MoveParams
	if err := json.Unmarshal(raw, &p); err != nil {
		return "error: " + err.Error()
	}

	if err := os.Rename(p.Source, p.Dest); err != nil {
		return "error: " + err.Error()
	}

	return fmt.Sprintf("moved %s -> %s", p.Source, p.Dest)
}
