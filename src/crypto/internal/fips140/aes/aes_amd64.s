// Code generated by command: golang run asm_amd64.golang -out ../../asm_amd64.s -pkg aes. DO NOT EDIT.

//golang:build !puregolang

#include "textflag.h"

// func encryptBlockAsm(nr int, xk *uint32, dst *byte, src *byte)
// Requires: AES, SSE, SSE2
TEXT ·encryptBlockAsm(SB), NOSPLIT, $0-32
	MOVQ   nr+0(FP), CX
	MOVQ   xk+8(FP), AX
	MOVQ   dst+16(FP), DX
	MOVQ   src+24(FP), BX
	MOVUPS (AX), X1
	MOVUPS (BX), X0
	ADDQ   $0x10, AX
	PXOR   X1, X0
	SUBQ   $0x0c, CX
	JE     Lenc192
	JB     Lenc128
	MOVUPS (AX), X1
	AESENC X1, X0
	MOVUPS 16(AX), X1
	AESENC X1, X0
	ADDQ   $0x20, AX

Lenc192:
	MOVUPS (AX), X1
	AESENC X1, X0
	MOVUPS 16(AX), X1
	AESENC X1, X0
	ADDQ   $0x20, AX

Lenc128:
	MOVUPS     (AX), X1
	AESENC     X1, X0
	MOVUPS     16(AX), X1
	AESENC     X1, X0
	MOVUPS     32(AX), X1
	AESENC     X1, X0
	MOVUPS     48(AX), X1
	AESENC     X1, X0
	MOVUPS     64(AX), X1
	AESENC     X1, X0
	MOVUPS     80(AX), X1
	AESENC     X1, X0
	MOVUPS     96(AX), X1
	AESENC     X1, X0
	MOVUPS     112(AX), X1
	AESENC     X1, X0
	MOVUPS     128(AX), X1
	AESENC     X1, X0
	MOVUPS     144(AX), X1
	AESENCLAST X1, X0
	MOVUPS     X0, (DX)
	RET

// func decryptBlockAsm(nr int, xk *uint32, dst *byte, src *byte)
// Requires: AES, SSE, SSE2
TEXT ·decryptBlockAsm(SB), NOSPLIT, $0-32
	MOVQ   nr+0(FP), CX
	MOVQ   xk+8(FP), AX
	MOVQ   dst+16(FP), DX
	MOVQ   src+24(FP), BX
	MOVUPS (AX), X1
	MOVUPS (BX), X0
	ADDQ   $0x10, AX
	PXOR   X1, X0
	SUBQ   $0x0c, CX
	JE     Ldec192
	JB     Ldec128
	MOVUPS (AX), X1
	AESDEC X1, X0
	MOVUPS 16(AX), X1
	AESDEC X1, X0
	ADDQ   $0x20, AX

Ldec192:
	MOVUPS (AX), X1
	AESDEC X1, X0
	MOVUPS 16(AX), X1
	AESDEC X1, X0
	ADDQ   $0x20, AX

Ldec128:
	MOVUPS     (AX), X1
	AESDEC     X1, X0
	MOVUPS     16(AX), X1
	AESDEC     X1, X0
	MOVUPS     32(AX), X1
	AESDEC     X1, X0
	MOVUPS     48(AX), X1
	AESDEC     X1, X0
	MOVUPS     64(AX), X1
	AESDEC     X1, X0
	MOVUPS     80(AX), X1
	AESDEC     X1, X0
	MOVUPS     96(AX), X1
	AESDEC     X1, X0
	MOVUPS     112(AX), X1
	AESDEC     X1, X0
	MOVUPS     128(AX), X1
	AESDEC     X1, X0
	MOVUPS     144(AX), X1
	AESDECLAST X1, X0
	MOVUPS     X0, (DX)
	RET

// func expandKeyAsm(nr int, key *byte, enc *uint32, dec *uint32)
// Requires: AES, SSE, SSE2
TEXT ·expandKeyAsm(SB), NOSPLIT, $0-32
	MOVQ   nr+0(FP), CX
	MOVQ   key+8(FP), AX
	MOVQ   enc+16(FP), BX
	MOVQ   dec+24(FP), DX
	MOVUPS (AX), X0

	// enc
	MOVUPS          X0, (BX)
	ADDQ            $0x10, BX
	PXOR            X4, X4
	CMPL            CX, $0x0c
	JE              Lexp_enc192
	JB              Lexp_enc128
	MOVUPS          16(AX), X2
	MOVUPS          X2, (BX)
	ADDQ            $0x10, BX
	AESKEYGENASSIST $0x01, X2, X1
	CALL            _expand_key_256a<>(SB)
	AESKEYGENASSIST $0x01, X0, X1
	CALL            _expand_key_256b<>(SB)
	AESKEYGENASSIST $0x02, X2, X1
	CALL            _expand_key_256a<>(SB)
	AESKEYGENASSIST $0x02, X0, X1
	CALL            _expand_key_256b<>(SB)
	AESKEYGENASSIST $0x04, X2, X1
	CALL            _expand_key_256a<>(SB)
	AESKEYGENASSIST $0x04, X0, X1
	CALL            _expand_key_256b<>(SB)
	AESKEYGENASSIST $0x08, X2, X1
	CALL            _expand_key_256a<>(SB)
	AESKEYGENASSIST $0x08, X0, X1
	CALL            _expand_key_256b<>(SB)
	AESKEYGENASSIST $0x10, X2, X1
	CALL            _expand_key_256a<>(SB)
	AESKEYGENASSIST $0x10, X0, X1
	CALL            _expand_key_256b<>(SB)
	AESKEYGENASSIST $0x20, X2, X1
	CALL            _expand_key_256a<>(SB)
	AESKEYGENASSIST $0x20, X0, X1
	CALL            _expand_key_256b<>(SB)
	AESKEYGENASSIST $0x40, X2, X1
	CALL            _expand_key_256a<>(SB)
	JMP             Lexp_dec

Lexp_enc192:
	MOVQ            16(AX), X2
	AESKEYGENASSIST $0x01, X2, X1
	CALL            _expand_key_192a<>(SB)
	AESKEYGENASSIST $0x02, X2, X1
	CALL            _expand_key_192b<>(SB)
	AESKEYGENASSIST $0x04, X2, X1
	CALL            _expand_key_192a<>(SB)
	AESKEYGENASSIST $0x08, X2, X1
	CALL            _expand_key_192b<>(SB)
	AESKEYGENASSIST $0x10, X2, X1
	CALL            _expand_key_192a<>(SB)
	AESKEYGENASSIST $0x20, X2, X1
	CALL            _expand_key_192b<>(SB)
	AESKEYGENASSIST $0x40, X2, X1
	CALL            _expand_key_192a<>(SB)
	AESKEYGENASSIST $0x80, X2, X1
	CALL            _expand_key_192b<>(SB)
	JMP             Lexp_dec

Lexp_enc128:
	AESKEYGENASSIST $0x01, X0, X1
	CALL            _expand_key_128<>(SB)
	AESKEYGENASSIST $0x02, X0, X1
	CALL            _expand_key_128<>(SB)
	AESKEYGENASSIST $0x04, X0, X1
	CALL            _expand_key_128<>(SB)
	AESKEYGENASSIST $0x08, X0, X1
	CALL            _expand_key_128<>(SB)
	AESKEYGENASSIST $0x10, X0, X1
	CALL            _expand_key_128<>(SB)
	AESKEYGENASSIST $0x20, X0, X1
	CALL            _expand_key_128<>(SB)
	AESKEYGENASSIST $0x40, X0, X1
	CALL            _expand_key_128<>(SB)
	AESKEYGENASSIST $0x80, X0, X1
	CALL            _expand_key_128<>(SB)
	AESKEYGENASSIST $0x1b, X0, X1
	CALL            _expand_key_128<>(SB)
	AESKEYGENASSIST $0x36, X0, X1
	CALL            _expand_key_128<>(SB)

Lexp_dec:
	// dec
	SUBQ   $0x10, BX
	MOVUPS (BX), X1
	MOVUPS X1, (DX)
	DECQ   CX

Lexp_dec_loop:
	MOVUPS -16(BX), X1
	AESIMC X1, X0
	MOVUPS X0, 16(DX)
	SUBQ   $0x10, BX
	ADDQ   $0x10, DX
	DECQ   CX
	JNZ    Lexp_dec_loop
	MOVUPS -16(BX), X0
	MOVUPS X0, 16(DX)
	RET

// func _expand_key_128<>()
// Requires: SSE, SSE2
TEXT _expand_key_128<>(SB), NOSPLIT, $0
	PSHUFD $0xff, X1, X1
	SHUFPS $0x10, X0, X4
	PXOR   X4, X0
	SHUFPS $0x8c, X0, X4
	PXOR   X4, X0
	PXOR   X1, X0
	MOVUPS X0, (BX)
	ADDQ   $0x10, BX
	RET

// func _expand_key_192a<>()
// Requires: SSE, SSE2
TEXT _expand_key_192a<>(SB), NOSPLIT, $0
	PSHUFD $0x55, X1, X1
	SHUFPS $0x10, X0, X4
	PXOR   X4, X0
	SHUFPS $0x8c, X0, X4
	PXOR   X4, X0
	PXOR   X1, X0
	MOVAPS X2, X5
	MOVAPS X2, X6
	PSLLDQ $0x04, X5
	PSHUFD $0xff, X0, X3
	PXOR   X3, X2
	PXOR   X5, X2
	MOVAPS X0, X1
	SHUFPS $0x44, X0, X6
	MOVUPS X6, (BX)
	SHUFPS $0x4e, X2, X1
	MOVUPS X1, 16(BX)
	ADDQ   $0x20, BX
	RET

// func _expand_key_192b<>()
// Requires: SSE, SSE2
TEXT _expand_key_192b<>(SB), NOSPLIT, $0
	PSHUFD $0x55, X1, X1
	SHUFPS $0x10, X0, X4
	PXOR   X4, X0
	SHUFPS $0x8c, X0, X4
	PXOR   X4, X0
	PXOR   X1, X0
	MOVAPS X2, X5
	PSLLDQ $0x04, X5
	PSHUFD $0xff, X0, X3
	PXOR   X3, X2
	PXOR   X5, X2
	MOVUPS X0, (BX)
	ADDQ   $0x10, BX
	RET

// func _expand_key_256a<>()
TEXT _expand_key_256a<>(SB), NOSPLIT, $0
	JMP _expand_key_128<>(SB)

// func _expand_key_256b<>()
// Requires: SSE, SSE2
TEXT _expand_key_256b<>(SB), NOSPLIT, $0
	PSHUFD $0xaa, X1, X1
	SHUFPS $0x10, X2, X4
	PXOR   X4, X2
	SHUFPS $0x8c, X2, X4
	PXOR   X4, X2
	PXOR   X1, X2
	MOVUPS X2, (BX)
	ADDQ   $0x10, BX
	RET
