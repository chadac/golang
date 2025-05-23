// Code generated by command: golang run ctr_amd64_asm.golang -out ../../ctr_amd64.s. DO NOT EDIT.

//golang:build !puregolang

#include "textflag.h"

// func ctrBlocks1Asm(nr int, xk *[60]uint32, dst *[16]byte, src *[16]byte, ivlo uint64, ivhi uint64)
// Requires: AES, SSE, SSE2, SSE4.1, SSSE3
TEXT ·ctrBlocks1Asm(SB), $0-48
	MOVQ   nr+0(FP), AX
	MOVQ   xk+8(FP), CX
	MOVQ   dst+16(FP), DX
	MOVQ   src+24(FP), BX
	MOVQ   ivlo+32(FP), SI
	MOVQ   ivhi+40(FP), DI
	MOVOU  bswapMask<>+0(SB), X0
	MOVQ   SI, X1
	PINSRQ $0x01, DI, X1
	PSHUFB X0, X1
	MOVUPS (CX), X0
	PXOR   X0, X1
	ADDQ   $0x10, CX
	SUBQ   $0x0c, AX
	JE     enc192
	JB     enc128
	MOVUPS (CX), X0
	AESENC X0, X1
	MOVUPS 16(CX), X0
	AESENC X0, X1
	ADDQ   $0x20, CX

enc192:
	MOVUPS (CX), X0
	AESENC X0, X1
	MOVUPS 16(CX), X0
	AESENC X0, X1
	ADDQ   $0x20, CX

enc128:
	MOVUPS     (CX), X0
	AESENC     X0, X1
	MOVUPS     16(CX), X0
	AESENC     X0, X1
	MOVUPS     32(CX), X0
	AESENC     X0, X1
	MOVUPS     48(CX), X0
	AESENC     X0, X1
	MOVUPS     64(CX), X0
	AESENC     X0, X1
	MOVUPS     80(CX), X0
	AESENC     X0, X1
	MOVUPS     96(CX), X0
	AESENC     X0, X1
	MOVUPS     112(CX), X0
	AESENC     X0, X1
	MOVUPS     128(CX), X0
	AESENC     X0, X1
	MOVUPS     144(CX), X0
	AESENCLAST X0, X1
	MOVUPS     (BX), X0
	PXOR       X1, X0
	MOVUPS     X0, (DX)
	RET

DATA bswapMask<>+0(SB)/8, $0x08090a0b0c0d0e0f
DATA bswapMask<>+8(SB)/8, $0x0001020304050607
GLOBL bswapMask<>(SB), RODATA|NOPTR, $16

// func ctrBlocks2Asm(nr int, xk *[60]uint32, dst *[32]byte, src *[32]byte, ivlo uint64, ivhi uint64)
// Requires: AES, SSE, SSE2, SSE4.1, SSSE3
TEXT ·ctrBlocks2Asm(SB), $0-48
	MOVQ   nr+0(FP), AX
	MOVQ   xk+8(FP), CX
	MOVQ   dst+16(FP), DX
	MOVQ   src+24(FP), BX
	MOVQ   ivlo+32(FP), SI
	MOVQ   ivhi+40(FP), DI
	MOVOU  bswapMask<>+0(SB), X0
	MOVQ   SI, X1
	PINSRQ $0x01, DI, X1
	PSHUFB X0, X1
	ADDQ   $0x01, SI
	ADCQ   $0x00, DI
	MOVQ   SI, X2
	PINSRQ $0x01, DI, X2
	PSHUFB X0, X2
	MOVUPS (CX), X0
	PXOR   X0, X1
	PXOR   X0, X2
	ADDQ   $0x10, CX
	SUBQ   $0x0c, AX
	JE     enc192
	JB     enc128
	MOVUPS (CX), X0
	AESENC X0, X1
	AESENC X0, X2
	MOVUPS 16(CX), X0
	AESENC X0, X1
	AESENC X0, X2
	ADDQ   $0x20, CX

enc192:
	MOVUPS (CX), X0
	AESENC X0, X1
	AESENC X0, X2
	MOVUPS 16(CX), X0
	AESENC X0, X1
	AESENC X0, X2
	ADDQ   $0x20, CX

enc128:
	MOVUPS     (CX), X0
	AESENC     X0, X1
	AESENC     X0, X2
	MOVUPS     16(CX), X0
	AESENC     X0, X1
	AESENC     X0, X2
	MOVUPS     32(CX), X0
	AESENC     X0, X1
	AESENC     X0, X2
	MOVUPS     48(CX), X0
	AESENC     X0, X1
	AESENC     X0, X2
	MOVUPS     64(CX), X0
	AESENC     X0, X1
	AESENC     X0, X2
	MOVUPS     80(CX), X0
	AESENC     X0, X1
	AESENC     X0, X2
	MOVUPS     96(CX), X0
	AESENC     X0, X1
	AESENC     X0, X2
	MOVUPS     112(CX), X0
	AESENC     X0, X1
	AESENC     X0, X2
	MOVUPS     128(CX), X0
	AESENC     X0, X1
	AESENC     X0, X2
	MOVUPS     144(CX), X0
	AESENCLAST X0, X1
	AESENCLAST X0, X2
	MOVUPS     (BX), X0
	PXOR       X1, X0
	MOVUPS     X0, (DX)
	MOVUPS     16(BX), X0
	PXOR       X2, X0
	MOVUPS     X0, 16(DX)
	RET

// func ctrBlocks4Asm(nr int, xk *[60]uint32, dst *[64]byte, src *[64]byte, ivlo uint64, ivhi uint64)
// Requires: AES, SSE, SSE2, SSE4.1, SSSE3
TEXT ·ctrBlocks4Asm(SB), $0-48
	MOVQ   nr+0(FP), AX
	MOVQ   xk+8(FP), CX
	MOVQ   dst+16(FP), DX
	MOVQ   src+24(FP), BX
	MOVQ   ivlo+32(FP), SI
	MOVQ   ivhi+40(FP), DI
	MOVOU  bswapMask<>+0(SB), X0
	MOVQ   SI, X1
	PINSRQ $0x01, DI, X1
	PSHUFB X0, X1
	ADDQ   $0x01, SI
	ADCQ   $0x00, DI
	MOVQ   SI, X2
	PINSRQ $0x01, DI, X2
	PSHUFB X0, X2
	ADDQ   $0x01, SI
	ADCQ   $0x00, DI
	MOVQ   SI, X3
	PINSRQ $0x01, DI, X3
	PSHUFB X0, X3
	ADDQ   $0x01, SI
	ADCQ   $0x00, DI
	MOVQ   SI, X4
	PINSRQ $0x01, DI, X4
	PSHUFB X0, X4
	MOVUPS (CX), X0
	PXOR   X0, X1
	PXOR   X0, X2
	PXOR   X0, X3
	PXOR   X0, X4
	ADDQ   $0x10, CX
	SUBQ   $0x0c, AX
	JE     enc192
	JB     enc128
	MOVUPS (CX), X0
	AESENC X0, X1
	AESENC X0, X2
	AESENC X0, X3
	AESENC X0, X4
	MOVUPS 16(CX), X0
	AESENC X0, X1
	AESENC X0, X2
	AESENC X0, X3
	AESENC X0, X4
	ADDQ   $0x20, CX

enc192:
	MOVUPS (CX), X0
	AESENC X0, X1
	AESENC X0, X2
	AESENC X0, X3
	AESENC X0, X4
	MOVUPS 16(CX), X0
	AESENC X0, X1
	AESENC X0, X2
	AESENC X0, X3
	AESENC X0, X4
	ADDQ   $0x20, CX

enc128:
	MOVUPS     (CX), X0
	AESENC     X0, X1
	AESENC     X0, X2
	AESENC     X0, X3
	AESENC     X0, X4
	MOVUPS     16(CX), X0
	AESENC     X0, X1
	AESENC     X0, X2
	AESENC     X0, X3
	AESENC     X0, X4
	MOVUPS     32(CX), X0
	AESENC     X0, X1
	AESENC     X0, X2
	AESENC     X0, X3
	AESENC     X0, X4
	MOVUPS     48(CX), X0
	AESENC     X0, X1
	AESENC     X0, X2
	AESENC     X0, X3
	AESENC     X0, X4
	MOVUPS     64(CX), X0
	AESENC     X0, X1
	AESENC     X0, X2
	AESENC     X0, X3
	AESENC     X0, X4
	MOVUPS     80(CX), X0
	AESENC     X0, X1
	AESENC     X0, X2
	AESENC     X0, X3
	AESENC     X0, X4
	MOVUPS     96(CX), X0
	AESENC     X0, X1
	AESENC     X0, X2
	AESENC     X0, X3
	AESENC     X0, X4
	MOVUPS     112(CX), X0
	AESENC     X0, X1
	AESENC     X0, X2
	AESENC     X0, X3
	AESENC     X0, X4
	MOVUPS     128(CX), X0
	AESENC     X0, X1
	AESENC     X0, X2
	AESENC     X0, X3
	AESENC     X0, X4
	MOVUPS     144(CX), X0
	AESENCLAST X0, X1
	AESENCLAST X0, X2
	AESENCLAST X0, X3
	AESENCLAST X0, X4
	MOVUPS     (BX), X0
	PXOR       X1, X0
	MOVUPS     X0, (DX)
	MOVUPS     16(BX), X0
	PXOR       X2, X0
	MOVUPS     X0, 16(DX)
	MOVUPS     32(BX), X0
	PXOR       X3, X0
	MOVUPS     X0, 32(DX)
	MOVUPS     48(BX), X0
	PXOR       X4, X0
	MOVUPS     X0, 48(DX)
	RET

// func ctrBlocks8Asm(nr int, xk *[60]uint32, dst *[128]byte, src *[128]byte, ivlo uint64, ivhi uint64)
// Requires: AES, SSE, SSE2, SSE4.1, SSSE3
TEXT ·ctrBlocks8Asm(SB), $0-48
	MOVQ   nr+0(FP), AX
	MOVQ   xk+8(FP), CX
	MOVQ   dst+16(FP), DX
	MOVQ   src+24(FP), BX
	MOVQ   ivlo+32(FP), SI
	MOVQ   ivhi+40(FP), DI
	MOVOU  bswapMask<>+0(SB), X0
	MOVQ   SI, X1
	PINSRQ $0x01, DI, X1
	PSHUFB X0, X1
	ADDQ   $0x01, SI
	ADCQ   $0x00, DI
	MOVQ   SI, X2
	PINSRQ $0x01, DI, X2
	PSHUFB X0, X2
	ADDQ   $0x01, SI
	ADCQ   $0x00, DI
	MOVQ   SI, X3
	PINSRQ $0x01, DI, X3
	PSHUFB X0, X3
	ADDQ   $0x01, SI
	ADCQ   $0x00, DI
	MOVQ   SI, X4
	PINSRQ $0x01, DI, X4
	PSHUFB X0, X4
	ADDQ   $0x01, SI
	ADCQ   $0x00, DI
	MOVQ   SI, X5
	PINSRQ $0x01, DI, X5
	PSHUFB X0, X5
	ADDQ   $0x01, SI
	ADCQ   $0x00, DI
	MOVQ   SI, X6
	PINSRQ $0x01, DI, X6
	PSHUFB X0, X6
	ADDQ   $0x01, SI
	ADCQ   $0x00, DI
	MOVQ   SI, X7
	PINSRQ $0x01, DI, X7
	PSHUFB X0, X7
	ADDQ   $0x01, SI
	ADCQ   $0x00, DI
	MOVQ   SI, X8
	PINSRQ $0x01, DI, X8
	PSHUFB X0, X8
	MOVUPS (CX), X0
	PXOR   X0, X1
	PXOR   X0, X2
	PXOR   X0, X3
	PXOR   X0, X4
	PXOR   X0, X5
	PXOR   X0, X6
	PXOR   X0, X7
	PXOR   X0, X8
	ADDQ   $0x10, CX
	SUBQ   $0x0c, AX
	JE     enc192
	JB     enc128
	MOVUPS (CX), X0
	AESENC X0, X1
	AESENC X0, X2
	AESENC X0, X3
	AESENC X0, X4
	AESENC X0, X5
	AESENC X0, X6
	AESENC X0, X7
	AESENC X0, X8
	MOVUPS 16(CX), X0
	AESENC X0, X1
	AESENC X0, X2
	AESENC X0, X3
	AESENC X0, X4
	AESENC X0, X5
	AESENC X0, X6
	AESENC X0, X7
	AESENC X0, X8
	ADDQ   $0x20, CX

enc192:
	MOVUPS (CX), X0
	AESENC X0, X1
	AESENC X0, X2
	AESENC X0, X3
	AESENC X0, X4
	AESENC X0, X5
	AESENC X0, X6
	AESENC X0, X7
	AESENC X0, X8
	MOVUPS 16(CX), X0
	AESENC X0, X1
	AESENC X0, X2
	AESENC X0, X3
	AESENC X0, X4
	AESENC X0, X5
	AESENC X0, X6
	AESENC X0, X7
	AESENC X0, X8
	ADDQ   $0x20, CX

enc128:
	MOVUPS     (CX), X0
	AESENC     X0, X1
	AESENC     X0, X2
	AESENC     X0, X3
	AESENC     X0, X4
	AESENC     X0, X5
	AESENC     X0, X6
	AESENC     X0, X7
	AESENC     X0, X8
	MOVUPS     16(CX), X0
	AESENC     X0, X1
	AESENC     X0, X2
	AESENC     X0, X3
	AESENC     X0, X4
	AESENC     X0, X5
	AESENC     X0, X6
	AESENC     X0, X7
	AESENC     X0, X8
	MOVUPS     32(CX), X0
	AESENC     X0, X1
	AESENC     X0, X2
	AESENC     X0, X3
	AESENC     X0, X4
	AESENC     X0, X5
	AESENC     X0, X6
	AESENC     X0, X7
	AESENC     X0, X8
	MOVUPS     48(CX), X0
	AESENC     X0, X1
	AESENC     X0, X2
	AESENC     X0, X3
	AESENC     X0, X4
	AESENC     X0, X5
	AESENC     X0, X6
	AESENC     X0, X7
	AESENC     X0, X8
	MOVUPS     64(CX), X0
	AESENC     X0, X1
	AESENC     X0, X2
	AESENC     X0, X3
	AESENC     X0, X4
	AESENC     X0, X5
	AESENC     X0, X6
	AESENC     X0, X7
	AESENC     X0, X8
	MOVUPS     80(CX), X0
	AESENC     X0, X1
	AESENC     X0, X2
	AESENC     X0, X3
	AESENC     X0, X4
	AESENC     X0, X5
	AESENC     X0, X6
	AESENC     X0, X7
	AESENC     X0, X8
	MOVUPS     96(CX), X0
	AESENC     X0, X1
	AESENC     X0, X2
	AESENC     X0, X3
	AESENC     X0, X4
	AESENC     X0, X5
	AESENC     X0, X6
	AESENC     X0, X7
	AESENC     X0, X8
	MOVUPS     112(CX), X0
	AESENC     X0, X1
	AESENC     X0, X2
	AESENC     X0, X3
	AESENC     X0, X4
	AESENC     X0, X5
	AESENC     X0, X6
	AESENC     X0, X7
	AESENC     X0, X8
	MOVUPS     128(CX), X0
	AESENC     X0, X1
	AESENC     X0, X2
	AESENC     X0, X3
	AESENC     X0, X4
	AESENC     X0, X5
	AESENC     X0, X6
	AESENC     X0, X7
	AESENC     X0, X8
	MOVUPS     144(CX), X0
	AESENCLAST X0, X1
	AESENCLAST X0, X2
	AESENCLAST X0, X3
	AESENCLAST X0, X4
	AESENCLAST X0, X5
	AESENCLAST X0, X6
	AESENCLAST X0, X7
	AESENCLAST X0, X8
	MOVUPS     (BX), X0
	PXOR       X1, X0
	MOVUPS     X0, (DX)
	MOVUPS     16(BX), X0
	PXOR       X2, X0
	MOVUPS     X0, 16(DX)
	MOVUPS     32(BX), X0
	PXOR       X3, X0
	MOVUPS     X0, 32(DX)
	MOVUPS     48(BX), X0
	PXOR       X4, X0
	MOVUPS     X0, 48(DX)
	MOVUPS     64(BX), X0
	PXOR       X5, X0
	MOVUPS     X0, 64(DX)
	MOVUPS     80(BX), X0
	PXOR       X6, X0
	MOVUPS     X0, 80(DX)
	MOVUPS     96(BX), X0
	PXOR       X7, X0
	MOVUPS     X0, 96(DX)
	MOVUPS     112(BX), X0
	PXOR       X8, X0
	MOVUPS     X0, 112(DX)
	RET
