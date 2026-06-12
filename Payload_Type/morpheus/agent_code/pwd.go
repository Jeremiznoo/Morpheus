//go:build windows
// +build windows

package main

import "os"

func cmdPwd() string {
	wd, err := os.Getwd()
	if err != nil {
		return "error: " + err.Error()
	}
	return wd
}
