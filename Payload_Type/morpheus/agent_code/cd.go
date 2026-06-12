//go:build windows
// +build windows

package main

import (
	"encoding/json"
	"os"
)

type CdParams struct {
	Path string `json:"path"`
}

func cmdCd(raw json.RawMessage) string {
	var p CdParams
	if err := json.Unmarshal(raw, &p); err != nil {
		return "error: " + err.Error()
	}
	if err := os.Chdir(p.Path); err != nil {
		return "error: " + err.Error()
	}
	wd, _ := os.Getwd()
	return wd
}
