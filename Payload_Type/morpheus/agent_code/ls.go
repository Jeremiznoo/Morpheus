//go:build windows
// +build windows

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
)

type LsParams struct {
	Path string `json:"path"`
}

func cmdLs(raw json.RawMessage) string {
	var p LsParams
	if err := json.Unmarshal(raw, &p); err != nil {
		p.Path = "."
	}

	if p.Path == "" {
		p.Path = "."
	}

	entries, err := os.ReadDir(p.Path)
	if err != nil {
		return "error: " + err.Error()
	}

	sort.Slice(entries, func(i, j int) bool {
		return strings.ToLower(entries[i].Name()) < strings.ToLower(entries[j].Name())
	})

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Directory: %s\n\n", p.Path))

	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}

		mode := info.Mode().String()
		size := info.Size()
		modTime := info.ModTime().Format("2006-01-02 15:04:05")
		name := entry.Name()

		if entry.IsDir() {
			sb.WriteString(fmt.Sprintf("d %s %10d %s  %s/", mode, size, modTime, name))
		} else {
			sb.WriteString(fmt.Sprintf("- %s %10d %s  %s", mode, size, modTime, name))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}
