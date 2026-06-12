//go:build windows

#include "textflag.h"

// func syscall4(ssn uintptr, a1, a2, a3, a4 uintptr) uintptr
TEXT ·syscall4(SB), NOSPLIT, $0-48
	MOVQ  ssn+0(FP), AX
	MOVQ  a1+8(FP), CX
	MOVQ  a2+16(FP), DX
	MOVQ  a3+24(FP), R8
	MOVQ  a4+32(FP), R9
	MOVQ  CX, R10
	SYSCALL
	MOVQ  AX, ret+40(FP)
	RET

// func syscall5(ssn uintptr, a1, a2, a3, a4, a5 uintptr) uintptr
TEXT ·syscall5(SB), NOSPLIT, $8-56
	MOVQ  a5+40(FP), DI
	PUSHQ DI
	MOVQ  a1+8(FP), CX
	MOVQ  a2+16(FP), DX
	MOVQ  a3+24(FP), R8
	MOVQ  a4+32(FP), R9
	MOVQ  ssn+0(FP), AX
	MOVQ  CX, R10
	SYSCALL
	POPQ  DI
	MOVQ  AX, ret+48(FP)
	RET
