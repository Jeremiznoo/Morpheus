//go:build windows
// +build windows

package main

import (
	"fmt"
	"net"
	"strings"
)

func cmdIfconfig() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "error: " + err.Error()
	}

	var sb strings.Builder
	for _, iface := range interfaces {
		addrs, _ := iface.Addrs()
		if len(addrs) == 0 {
			continue
		}

		sb.WriteString(fmt.Sprintf("%s:\n", iface.Name))
		sb.WriteString(fmt.Sprintf("  MAC:  %s\n", iface.HardwareAddr.String()))

		for _, addr := range addrs {
			sb.WriteString(fmt.Sprintf("  IP:   %s\n", addr.String()))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}
