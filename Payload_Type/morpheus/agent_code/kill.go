//go:build windows
// +build windows

package main

import (
	"encoding/json"
	"fmt"
)

type KillParams struct {
	PID uint32 `json:"pid"`
}

func cmdKill(raw json.RawMessage) string {
	var p KillParams
	if err := json.Unmarshal(raw, &p); err != nil {
		return "error: " + err.Error()
	}

	if err := KillProcess(p.PID); err != nil {
		return fmt.Sprintf("error: %s", err.Error())
	}

	return fmt.Sprintf("killed pid %d", p.PID)
}
