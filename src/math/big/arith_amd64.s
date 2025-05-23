// Copyright 2025 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Code generated by 'golang generate' (with ./internal/asmgen). DO NOT EDIT.

//golang:build !math_big_pure_golang

#include "textflag.h"

// func addVV(z, x, y []Word) (c Word)
TEXT ·addVV(SB), NOSPLIT, $0
	MOVQ z_len+8(FP), BX
	MOVQ x_base+24(FP), SI
	MOVQ y_base+48(FP), DI
	MOVQ z_base+0(FP), R8
	// compute unrolled loop lengths
	MOVQ BX, R9
	ANDQ $3, R9
	SHRQ $2, BX
	MOVQ $0, R10	// clear saved carry
loop1:
	TESTQ R9, R9; JZ loop1done
loop1cont:
	// unroll 1X
	ADDQ R10, R10	// restore carry
	MOVQ 0(SI), R10
	ADCQ 0(DI), R10
	MOVQ R10, 0(R8)
	SBBQ R10, R10	// save carry
	LEAQ 8(SI), SI	// ADD $8, SI
	LEAQ 8(DI), DI	// ADD $8, DI
	LEAQ 8(R8), R8	// ADD $8, R8
	SUBQ $1, R9; JNZ loop1cont
loop1done:
loop4:
	TESTQ BX, BX; JZ loop4done
loop4cont:
	// unroll 4X
	ADDQ R10, R10	// restore carry
	MOVQ 0(SI), R9
	MOVQ 8(SI), R10
	MOVQ 16(SI), R11
	MOVQ 24(SI), R12
	ADCQ 0(DI), R9
	ADCQ 8(DI), R10
	ADCQ 16(DI), R11
	ADCQ 24(DI), R12
	MOVQ R9, 0(R8)
	MOVQ R10, 8(R8)
	MOVQ R11, 16(R8)
	MOVQ R12, 24(R8)
	SBBQ R10, R10	// save carry
	LEAQ 32(SI), SI	// ADD $32, SI
	LEAQ 32(DI), DI	// ADD $32, DI
	LEAQ 32(R8), R8	// ADD $32, R8
	SUBQ $1, BX; JNZ loop4cont
loop4done:
	NEGQ R10	// convert add carry
	MOVQ R10, c+72(FP)
	RET

// func subVV(z, x, y []Word) (c Word)
TEXT ·subVV(SB), NOSPLIT, $0
	MOVQ z_len+8(FP), BX
	MOVQ x_base+24(FP), SI
	MOVQ y_base+48(FP), DI
	MOVQ z_base+0(FP), R8
	// compute unrolled loop lengths
	MOVQ BX, R9
	ANDQ $3, R9
	SHRQ $2, BX
	MOVQ $0, R10	// clear saved carry
loop1:
	TESTQ R9, R9; JZ loop1done
loop1cont:
	// unroll 1X
	ADDQ R10, R10	// restore carry
	MOVQ 0(SI), R10
	SBBQ 0(DI), R10
	MOVQ R10, 0(R8)
	SBBQ R10, R10	// save carry
	LEAQ 8(SI), SI	// ADD $8, SI
	LEAQ 8(DI), DI	// ADD $8, DI
	LEAQ 8(R8), R8	// ADD $8, R8
	SUBQ $1, R9; JNZ loop1cont
loop1done:
loop4:
	TESTQ BX, BX; JZ loop4done
loop4cont:
	// unroll 4X
	ADDQ R10, R10	// restore carry
	MOVQ 0(SI), R9
	MOVQ 8(SI), R10
	MOVQ 16(SI), R11
	MOVQ 24(SI), R12
	SBBQ 0(DI), R9
	SBBQ 8(DI), R10
	SBBQ 16(DI), R11
	SBBQ 24(DI), R12
	MOVQ R9, 0(R8)
	MOVQ R10, 8(R8)
	MOVQ R11, 16(R8)
	MOVQ R12, 24(R8)
	SBBQ R10, R10	// save carry
	LEAQ 32(SI), SI	// ADD $32, SI
	LEAQ 32(DI), DI	// ADD $32, DI
	LEAQ 32(R8), R8	// ADD $32, R8
	SUBQ $1, BX; JNZ loop4cont
loop4done:
	NEGQ R10	// convert sub carry
	MOVQ R10, c+72(FP)
	RET

// func lshVU(z, x []Word, s uint) (c Word)
TEXT ·lshVU(SB), NOSPLIT, $0
	MOVQ z_len+8(FP), BX
	TESTQ BX, BX; JZ ret0
	MOVQ s+48(FP), CX
	MOVQ x_base+24(FP), SI
	MOVQ z_base+0(FP), DI
	// run loop backward
	LEAQ (SI)(BX*8), SI
	LEAQ (DI)(BX*8), DI
	// shift first word into carry
	MOVQ -8(SI), R8
	MOVQ $0, R9
	SHLQ CX, R8, R9
	MOVQ R9, c+56(FP)
	// shift remaining words
	SUBQ $1, BX
	// compute unrolled loop lengths
	MOVQ BX, R9
	ANDQ $3, R9
	SHRQ $2, BX
loop1:
	TESTQ R9, R9; JZ loop1done
loop1cont:
	// unroll 1X
	MOVQ -16(SI), R10
	SHLQ CX, R10, R8
	MOVQ R8, -8(DI)
	MOVQ R10, R8
	LEAQ -8(SI), SI	// ADD $-8, SI
	LEAQ -8(DI), DI	// ADD $-8, DI
	SUBQ $1, R9; JNZ loop1cont
loop1done:
loop4:
	TESTQ BX, BX; JZ loop4done
loop4cont:
	// unroll 4X
	MOVQ -16(SI), R9
	MOVQ -24(SI), R10
	MOVQ -32(SI), R11
	MOVQ -40(SI), R12
	SHLQ CX, R9, R8
	SHLQ CX, R10, R9
	SHLQ CX, R11, R10
	SHLQ CX, R12, R11
	MOVQ R8, -8(DI)
	MOVQ R9, -16(DI)
	MOVQ R10, -24(DI)
	MOVQ R11, -32(DI)
	MOVQ R12, R8
	LEAQ -32(SI), SI	// ADD $-32, SI
	LEAQ -32(DI), DI	// ADD $-32, DI
	SUBQ $1, BX; JNZ loop4cont
loop4done:
	// store final shifted bits
	SHLQ CX, R8
	MOVQ R8, -8(DI)
	RET
ret0:
	MOVQ $0, c+56(FP)
	RET

// func rshVU(z, x []Word, s uint) (c Word)
TEXT ·rshVU(SB), NOSPLIT, $0
	MOVQ z_len+8(FP), BX
	TESTQ BX, BX; JZ ret0
	MOVQ s+48(FP), CX
	MOVQ x_base+24(FP), SI
	MOVQ z_base+0(FP), DI
	// shift first word into carry
	MOVQ 0(SI), R8
	MOVQ $0, R9
	SHRQ CX, R8, R9
	MOVQ R9, c+56(FP)
	// shift remaining words
	SUBQ $1, BX
	// compute unrolled loop lengths
	MOVQ BX, R9
	ANDQ $3, R9
	SHRQ $2, BX
loop1:
	TESTQ R9, R9; JZ loop1done
loop1cont:
	// unroll 1X
	MOVQ 8(SI), R10
	SHRQ CX, R10, R8
	MOVQ R8, 0(DI)
	MOVQ R10, R8
	LEAQ 8(SI), SI	// ADD $8, SI
	LEAQ 8(DI), DI	// ADD $8, DI
	SUBQ $1, R9; JNZ loop1cont
loop1done:
loop4:
	TESTQ BX, BX; JZ loop4done
loop4cont:
	// unroll 4X
	MOVQ 8(SI), R9
	MOVQ 16(SI), R10
	MOVQ 24(SI), R11
	MOVQ 32(SI), R12
	SHRQ CX, R9, R8
	SHRQ CX, R10, R9
	SHRQ CX, R11, R10
	SHRQ CX, R12, R11
	MOVQ R8, 0(DI)
	MOVQ R9, 8(DI)
	MOVQ R10, 16(DI)
	MOVQ R11, 24(DI)
	MOVQ R12, R8
	LEAQ 32(SI), SI	// ADD $32, SI
	LEAQ 32(DI), DI	// ADD $32, DI
	SUBQ $1, BX; JNZ loop4cont
loop4done:
	// store final shifted bits
	SHRQ CX, R8
	MOVQ R8, 0(DI)
	RET
ret0:
	MOVQ $0, c+56(FP)
	RET

// func mulAddVWW(z, x []Word, m, a Word) (c Word)
TEXT ·mulAddVWW(SB), NOSPLIT, $0
	MOVQ m+48(FP), BX
	MOVQ a+56(FP), SI
	MOVQ z_len+8(FP), DI
	MOVQ x_base+24(FP), R8
	MOVQ z_base+0(FP), R9
	// compute unrolled loop lengths
	MOVQ DI, R10
	ANDQ $3, R10
	SHRQ $2, DI
loop1:
	TESTQ R10, R10; JZ loop1done
loop1cont:
	// unroll 1X in batches of 1
	MOVQ 0(R8), AX
	// multiply
	MULQ BX
	ADDQ SI, AX
	MOVQ DX, SI
	ADCQ $0, SI
	MOVQ AX, 0(R9)
	LEAQ 8(R8), R8	// ADD $8, R8
	LEAQ 8(R9), R9	// ADD $8, R9
	SUBQ $1, R10; JNZ loop1cont
loop1done:
loop4:
	TESTQ DI, DI; JZ loop4done
loop4cont:
	// unroll 4X in batches of 1
	MOVQ 0(R8), AX
	// multiply
	MULQ BX
	ADDQ SI, AX
	MOVQ DX, SI
	ADCQ $0, SI
	MOVQ AX, 0(R9)
	MOVQ 8(R8), AX
	// multiply
	MULQ BX
	ADDQ SI, AX
	MOVQ DX, SI
	ADCQ $0, SI
	MOVQ AX, 8(R9)
	MOVQ 16(R8), AX
	// multiply
	MULQ BX
	ADDQ SI, AX
	MOVQ DX, SI
	ADCQ $0, SI
	MOVQ AX, 16(R9)
	MOVQ 24(R8), AX
	// multiply
	MULQ BX
	ADDQ SI, AX
	MOVQ DX, SI
	ADCQ $0, SI
	MOVQ AX, 24(R9)
	LEAQ 32(R8), R8	// ADD $32, R8
	LEAQ 32(R9), R9	// ADD $32, R9
	SUBQ $1, DI; JNZ loop4cont
loop4done:
	MOVQ SI, c+64(FP)
	RET

// func addMulVVWW(z, x, y []Word, m, a Word) (c Word)
TEXT ·addMulVVWW(SB), NOSPLIT, $0
	CMPB ·hasADX(SB), $0; JNZ altcarry
	MOVQ m+72(FP), BX
	MOVQ a+80(FP), SI
	MOVQ z_len+8(FP), DI
	MOVQ x_base+24(FP), R8
	MOVQ y_base+48(FP), R9
	MOVQ z_base+0(FP), R10
	// compute unrolled loop lengths
	MOVQ DI, R11
	ANDQ $3, R11
	SHRQ $2, DI
loop1:
	TESTQ R11, R11; JZ loop1done
loop1cont:
	// unroll 1X in batches of 1
	MOVQ 0(R9), AX
	// multiply
	MULQ BX
	ADDQ SI, AX
	MOVQ DX, SI
	ADCQ $0, SI
	// add
	ADDQ 0(R8), AX
	ADCQ $0, SI
	MOVQ AX, 0(R10)
	LEAQ 8(R8), R8	// ADD $8, R8
	LEAQ 8(R9), R9	// ADD $8, R9
	LEAQ 8(R10), R10	// ADD $8, R10
	SUBQ $1, R11; JNZ loop1cont
loop1done:
loop4:
	TESTQ DI, DI; JZ loop4done
loop4cont:
	// unroll 4X in batches of 1
	MOVQ 0(R9), AX
	// multiply
	MULQ BX
	ADDQ SI, AX
	MOVQ DX, SI
	ADCQ $0, SI
	// add
	ADDQ 0(R8), AX
	ADCQ $0, SI
	MOVQ AX, 0(R10)
	MOVQ 8(R9), AX
	// multiply
	MULQ BX
	ADDQ SI, AX
	MOVQ DX, SI
	ADCQ $0, SI
	// add
	ADDQ 8(R8), AX
	ADCQ $0, SI
	MOVQ AX, 8(R10)
	MOVQ 16(R9), AX
	// multiply
	MULQ BX
	ADDQ SI, AX
	MOVQ DX, SI
	ADCQ $0, SI
	// add
	ADDQ 16(R8), AX
	ADCQ $0, SI
	MOVQ AX, 16(R10)
	MOVQ 24(R9), AX
	// multiply
	MULQ BX
	ADDQ SI, AX
	MOVQ DX, SI
	ADCQ $0, SI
	// add
	ADDQ 24(R8), AX
	ADCQ $0, SI
	MOVQ AX, 24(R10)
	LEAQ 32(R8), R8	// ADD $32, R8
	LEAQ 32(R9), R9	// ADD $32, R9
	LEAQ 32(R10), R10	// ADD $32, R10
	SUBQ $1, DI; JNZ loop4cont
loop4done:
	MOVQ SI, c+88(FP)
	RET
altcarry:
	MOVQ m+72(FP), DX
	MOVQ a+80(FP), BX
	MOVQ z_len+8(FP), SI
	MOVQ $0, DI
	MOVQ x_base+24(FP), R8
	MOVQ y_base+48(FP), R9
	MOVQ z_base+0(FP), R10
	// compute unrolled loop lengths
	MOVQ SI, R11
	ANDQ $7, R11
	SHRQ $3, SI
alt1:
	TESTQ R11, R11; JZ alt1done
alt1cont:
	// unroll 1X
	// multiply and add
	TESTQ AX, AX	// clear carry
	TESTQ AX, AX	// clear carry
	MULXQ 0(R9), R13, R12
	ADCXQ BX, R13
	ADOXQ 0(R8), R13
	MOVQ R13, 0(R10)
	MOVQ R12, BX
	ADCXQ DI, BX
	ADOXQ DI, BX
	LEAQ 8(R8), R8	// ADD $8, R8
	LEAQ 8(R9), R9	// ADD $8, R9
	LEAQ 8(R10), R10	// ADD $8, R10
	SUBQ $1, R11; JNZ alt1cont
alt1done:
alt8:
	TESTQ SI, SI; JZ alt8done
alt8cont:
	// unroll 8X in batches of 2
	// multiply and add
	TESTQ AX, AX	// clear carry
	TESTQ AX, AX	// clear carry
	MULXQ 0(R9), R13, R11
	ADCXQ BX, R13
	ADOXQ 0(R8), R13
	MULXQ 8(R9), R14, BX
	ADCXQ R11, R14
	ADOXQ 8(R8), R14
	MOVQ R13, 0(R10)
	MOVQ R14, 8(R10)
	MULXQ 16(R9), R13, R11
	ADCXQ BX, R13
	ADOXQ 16(R8), R13
	MULXQ 24(R9), R14, BX
	ADCXQ R11, R14
	ADOXQ 24(R8), R14
	MOVQ R13, 16(R10)
	MOVQ R14, 24(R10)
	MULXQ 32(R9), R13, R11
	ADCXQ BX, R13
	ADOXQ 32(R8), R13
	MULXQ 40(R9), R14, BX
	ADCXQ R11, R14
	ADOXQ 40(R8), R14
	MOVQ R13, 32(R10)
	MOVQ R14, 40(R10)
	MULXQ 48(R9), R13, R11
	ADCXQ BX, R13
	ADOXQ 48(R8), R13
	MULXQ 56(R9), R14, BX
	ADCXQ R11, R14
	ADOXQ 56(R8), R14
	MOVQ R13, 48(R10)
	MOVQ R14, 56(R10)
	ADCXQ DI, BX
	ADOXQ DI, BX
	LEAQ 64(R8), R8	// ADD $64, R8
	LEAQ 64(R9), R9	// ADD $64, R9
	LEAQ 64(R10), R10	// ADD $64, R10
	SUBQ $1, SI; JNZ alt8cont
alt8done:
	MOVQ BX, c+88(FP)
	RET
