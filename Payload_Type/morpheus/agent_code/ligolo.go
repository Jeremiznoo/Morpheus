//go:build windows
// +build windows

package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

type LigoloTunnel struct {
	Running  bool
	Name     string
	Listener string
}

type LigoloParams struct {
	Action   string `json:"action"`
	Tunnel   string `json:"tunnel"`
	Listener string `json:"listener"`
}

func cmdLigoloStart(raw json.RawMessage, a *Agent) string {
	var p LigoloParams
	if err := json.Unmarshal(raw, &p); err != nil {
		return "error: " + err.Error()
	}

	if a.Ligolo != nil && a.Ligolo.Running {
		return "ligolo tunnel already running"
	}

	listener := p.Listener
	if listener == "" {
		listener = "0.0.0.0:11601"
	}

	a.Ligolo = &LigoloTunnel{
		Running:  true,
		Name:     p.Tunnel,
		Listener: listener,
	}

	return fmt.Sprintf("ligolo tunnel '%s' started on %s", p.Tunnel, listener)
}

func cmdLigoloStop(raw json.RawMessage, a *Agent) string {
	if a.Ligolo == nil || !a.Ligolo.Running {
		return "no ligolo tunnel running"
	}

	a.Ligolo.Running = false
	a.Ligolo = nil
	return "ligolo tunnel stopped"
}

func cmdLigoloStatus(a *Agent) string {
	if a.Ligolo == nil || !a.Ligolo.Running {
		return "ligolo: not running"
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Ligolo Tunnel:\n"))
	sb.WriteString(fmt.Sprintf("  Name    : %s\n", a.Ligolo.Name))
	sb.WriteString(fmt.Sprintf("  Listener: %s\n", a.Ligolo.Listener))
	return sb.String()
}
