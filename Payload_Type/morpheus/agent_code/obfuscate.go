//go:build windows
// +build windows

package main

import (
	"errors"
	"unsafe"
)

var (
	errCannotGetPeb = errors.New(xorDecrypt([]byte{0xc9, 0xca, 0xc6, 0xc7, 0xc1, 0xdb, 0x8c, 0xca, 0xc7, 0xd7, 0x80, 0xf1, 0xe3, 0xe5}))
	errCannotGetLdr = errors.New(xorDecrypt([]byte{0xc9, 0xca, 0xc6, 0xc7, 0xc1, 0xdb, 0x8c, 0xca, 0xc7, 0xd7, 0x80, 0xed, 0xc2, 0xd5}))
	errApiNotFound  = errors.New(xorDecrypt([]byte{0xeb, 0xfb, 0xe1, 0x89, 0xc0, 0xc0, 0xd8, 0x8d, 0xc4, 0xcc, 0xd5, 0xcf, 0xc2}))
)

func readPEB() uintptr

func xorDecrypt(data []byte) string {
	if len(data) == 0 {
		return ""
	}
	r := make([]byte, len(data))
	for i, b := range data {
		r[i] = b ^ byte(i) ^ 0xAA
	}
	return string(r)
}

func xorEncrypt(s string) []byte {
	r := make([]byte, len(s))
	for i, b := range []byte(s) {
		r[i] = b ^ byte(i) ^ 0xAA
	}
	return r
}

func hashUpper(s string) uint32 {
	var h uint32 = 5381
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'a' && c <= 'z' {
			c -= 32
		}
		h = ((h << 5) + h) + uint32(c)
	}
	return h
}

type peb struct {
	_      [2]byte
	Bebug  byte
	_      [5]byte
	Ldr    uintptr // +0x18 on x64
}

type pebLdr struct {
	_       [4]byte
	_       [4]byte
	_       [24]byte
	ModList listEntry // +0x20
}

type listEntry struct {
	Flink uintptr
	Blink uintptr
}

type unicodeStr struct {
	Len    uint16
	MaxLen uint16
	Buf    uintptr
}

type ldrEntry struct {
	InLoadOrder   [2]uintptr // +0x00
	InMemoryOrder [2]uintptr // +0x10
	InInitOrder   [2]uintptr // +0x20
	DllBase       uintptr    // +0x30
	EntryPoint    uintptr    // +0x38
	SizeOfImage   uintptr    // +0x40
	FullDllName   unicodeStr // +0x48
	BaseDllName   unicodeStr // +0x58
}

type imgDos struct {
	Magic  uint16
	_      [58]byte
	Lfanew int32
}

type imgExp struct {
	_    [4]byte
	_    [4]byte
	_    [4]byte
	_    [4]byte
	_    [4]byte
	Base uint32
	NFn  uint32
	NNm  uint32
	AFn  uint32
	ANm  uint32
	AOrd uint32
}

func getPeb() uintptr {
	return readPEB()
}

func findExportFn(dll uintptr, fnHash uint32) uintptr {
	if dll == 0 {
		return 0
	}
	dos := (*imgDos)(unsafe.Pointer(dll))
	if dos.Magic != 0x5A4D {
		return 0
	}
	ntHdr := dll + uintptr(dos.Lfanew)
	if *(*uint32)(unsafe.Pointer(ntHdr)) != 0x00004550 {
		return 0
	}
	expVA := *(*uint32)(unsafe.Pointer(ntHdr + 0x88))
	if expVA == 0 {
		return 0
	}
	exp := (*imgExp)(unsafe.Pointer(dll + uintptr(expVA)))
	funcArr := dll + uintptr(exp.AFn)
	nameArr := dll + uintptr(exp.ANm)
	ordArr := dll + uintptr(exp.AOrd)

	for i := uint32(0); i < exp.NNm; i++ {
		nrva := *(*uint32)(unsafe.Pointer(nameArr + uintptr(i*4)))
		if nrva == 0 {
			continue
		}
		np := dll + uintptr(nrva)
		h := uint32(5381)
		for j := 0; ; j++ {
			b := *(*byte)(unsafe.Pointer(np + uintptr(j)))
			if b == 0 {
				break
			}
			if b >= 'a' && b <= 'z' {
				b -= 32
			}
			h = ((h << 5) + h) + uint32(b)
		}
		if h == fnHash {
			ord := *(*uint16)(unsafe.Pointer(ordArr + uintptr(i*2)))
			rva := *(*uint32)(unsafe.Pointer(funcArr + uintptr(ord)*4))
			return dll + uintptr(rva)
		}
	}
	return 0
}

func resolveAPI(modHash, fnHash uint32) (uintptr, error) {
	p := getPeb()
	if p == 0 {
		return 0, errCannotGetPeb
	}
	ldr := *(*uintptr)(unsafe.Pointer(p + 0x18))
	if ldr == 0 {
		return 0, errCannotGetLdr
	}
	flink := *(*uintptr)(unsafe.Pointer(ldr + 0x20))
	startFlink := flink

	for {
		le := (*ldrEntry)(unsafe.Pointer(flink))
		db := le.DllBase
		if db != 0 && le.BaseDllName.Len > 0 && le.BaseDllName.Buf != 0 {
			ml := int(le.BaseDllName.Len) / 2
			h := uint32(5381)
			for i := 0; i < ml; i++ {
				c := byte(*(*uint16)(unsafe.Pointer(le.BaseDllName.Buf + uintptr(i*2))))
				if c >= 'a' && c <= 'z' {
					c -= 32
				}
				h = ((h << 5) + h) + uint32(c)
			}
			if h == modHash {
				addr := findExportFn(db, fnHash)
				if addr != 0 {
					return addr, nil
				}
			}
		}
		flink = *(*uintptr)(unsafe.Pointer(flink))
		if flink == startFlink {
			break
		}
	}
	return 0, errApiNotFound
}
