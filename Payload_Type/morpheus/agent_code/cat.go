//go:build windows
// +build windows

package main

import (
	"encoding/json"
	"os"
)

type CatParams struct {
	Path string `json:"path"`
}

func cmdCat(raw json.RawMessage) string {
	var p CatParams
	if err := json.Unmarshal(raw, &p); err != nil {
		return "error: " + err.Error()
	}

	data, err := os.ReadFile(p.Path)
	if err != nil {
		return "error: " + err.Error()
	}

	return string(data)
}
