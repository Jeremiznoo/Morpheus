//go:build windows
// +build windows

package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
)

type UploadParams struct {
	Path      string `json:"path"`
	Content   string `json:"content"`
	FileID    string `json:"file_id"`
}

func cmdUpload(raw json.RawMessage) string {
	var p UploadParams
	if err := json.Unmarshal(raw, &p); err != nil {
		return "error: " + err.Error()
	}

	var data []byte
	var err error

	if p.Content != "" {
		data, err = base64.StdEncoding.DecodeString(p.Content)
		if err != nil {
			return "error: content decode: " + err.Error()
		}
	} else {
		return "error: no content or file_id provided"
	}

	if err := os.WriteFile(p.Path, data, 0644); err != nil {
		return "error: " + err.Error()
	}

	return fmt.Sprintf("uploaded %s (%d bytes)", p.Path, len(data))
}
