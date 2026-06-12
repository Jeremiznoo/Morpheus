//go:build windows
// +build windows

package main

import (
	"encoding/json"
	"fmt"
)

type MakeTokenParams struct {
	Username string `json:"username"`
	Domain   string `json:"domain"`
	Password string `json:"password"`
}

func cmdMakeToken(raw json.RawMessage) string {
	var p MakeTokenParams
	if err := json.Unmarshal(raw, &p); err != nil {
		return "error: " + err.Error()
	}

	if err := MakeToken(p.Username, p.Domain, p.Password); err != nil {
		return fmt.Sprintf("error: %s", err.Error())
	}

	return fmt.Sprintf("token created for %s\\%s", p.Domain, p.Username)
}
