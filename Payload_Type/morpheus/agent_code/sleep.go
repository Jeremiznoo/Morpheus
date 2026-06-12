//go:build windows
// +build windows

package main

import (
	"encoding/json"
	"fmt"
)

type SleepParams struct {
	Interval int     `json:"interval"`
	Jitter   float64 `json:"jitter"`
}

func cmdSleep(raw json.RawMessage, a *Agent) string {
	var p SleepParams
	if err := json.Unmarshal(raw, &p); err != nil {
		return "error: " + err.Error()
	}

	if p.Interval > 0 {
		a.Sleep = p.Interval
	}
	if p.Jitter >= 0 && p.Jitter <= 100 {
		a.Jitter = p.Jitter
	}

	return fmt.Sprintf("sleep set to %ds with %.0f%% jitter", a.Sleep, a.Jitter)
}
