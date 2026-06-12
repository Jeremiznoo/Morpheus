//go:build windows
// +build windows

package main

func cmdRev2self() string {
	if err := Rev2Self(); err != nil {
		return "error: " + err.Error()
	}
	return "reverted to self"
}
