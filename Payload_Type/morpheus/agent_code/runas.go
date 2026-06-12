//go:build windows
// +build windows

package main

import (
	"encoding/json"
	"fmt"
)

type RunasParams struct {
	Username string `json:"username"`
	Domain   string `json:"domain"`
	Password string `json:"password"`
	Command  string `json:"command"`
}

func cmdRunas(raw json.RawMessage) string {
	var p RunasParams
	if err := json.Unmarshal(raw, &p); err != nil {
		return "error: " + err.Error()
	}

	result, err := RunAsUser(p.Username, p.Domain, p.Password, p.Command)
	if err != nil {
		return fmt.Sprintf("error: %s", err.Error())
	}

	return fmt.Sprintf("runas completed: %s", result)
}
