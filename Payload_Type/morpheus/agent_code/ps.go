//go:build windows
// +build windows

package main

import (
	"fmt"
	"strings"
)

func cmdPs() string {
	procs, err := GetProcessList()
	if err != nil {
		return "error: " + err.Error()
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%-8s %-40s\n", "PID", "Name"))
	sb.WriteString(strings.Repeat("-", 48) + "\n")

	for _, p := range procs {
		sb.WriteString(fmt.Sprintf("%-8d %-40s\n", p.PID, p.Name))
	}

	return sb.String()
}
