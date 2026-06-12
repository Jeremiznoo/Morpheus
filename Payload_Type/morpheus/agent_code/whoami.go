//go:build windows
// +build windows

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type WhoamiParams struct {
	Privs bool `json:"privileges"`
}

func cmdWhoami(raw json.RawMessage) string {
	var p WhoamiParams
	json.Unmarshal(raw, &p)

	username := getUsername()
	hostname, _ := os.Hostname()
	domain := GetDomainInfo()
	integrity := GetIntegrityLevel()

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("User     : %s\\%s\n", hostname, username))
	sb.WriteString(fmt.Sprintf("Hostname : %s\n", hostname))
	sb.WriteString(fmt.Sprintf("Domain   : %s\n", domain))
	sb.WriteString(fmt.Sprintf("PID      : %d\n", getProcessID()))
	sb.WriteString(fmt.Sprintf("Arch     : %s\n", GetArchitecture()))
	sb.WriteString(fmt.Sprintf("Integrity: %d\n", int(integrity)))

	return sb.String()
}
