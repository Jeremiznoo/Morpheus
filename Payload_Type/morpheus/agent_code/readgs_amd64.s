//go:build windows
// +build windows

#include "textflag.h"

TEXT ·readPEB(SB), NOSPLIT, $0-8
    MOVQ  (GS), AX
    MOVQ  0x60(AX), AX
    MOVQ  AX, ret+0(FP)
    RET
