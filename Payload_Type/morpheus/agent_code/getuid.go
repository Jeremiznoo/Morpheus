//go:build windows
// +build windows

package main

import (
	"fmt"
	"os"
)

func cmdGetuid() string {
	username := getUsername()
	hostname, _ := os.Hostname()
	return fmt.Sprintf("%s\\%s", hostname, username)
}
