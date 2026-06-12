//go:build windows

package main

var (
	maskState uint32
	maskInit  bool
)

func MaskSleep(key []byte) {
	if !maskInit {
		return
	}
	for i := range key {
		key[i] ^= byte(i) ^ byte(maskState)
	}
}

func UnmaskSleep(key []byte) {
	if !maskInit {
		return
	}
	for i := range key {
		key[i] ^= byte(i) ^ byte(maskState)
	}
	maskState++
}

func InitSleepMask() {
	maskState = 0xAAAAAAAA
	maskInit = true
}
