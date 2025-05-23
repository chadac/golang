// Copyright 2017 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build !puregolang

#include "textflag.h"

#define HASHUPDATECHOOSE \
	SHA1C	V16.S4, V1, V2 \
	SHA1H	V3, V1 \
	VMOV	V2.B16, V3.B16

#define HASHUPDATEPARITY \
	SHA1P	V16.S4, V1, V2 \
	SHA1H	V3, V1 \
	VMOV	V2.B16, V3.B16

#define HASHUPDATEMAJ \
	SHA1M	V16.S4, V1, V2 \
	SHA1H	V3, V1 \
	VMOV	V2.B16, V3.B16

// func sha1block(h []uint32, p []byte, k []uint32)
TEXT ·sha1block(SB),NOSPLIT,$0
	MOVD	h_base+0(FP), R0                             // hash value first address
	MOVD	p_base+24(FP), R1                            // message first address
	MOVD	k_base+48(FP), R2                            // k constants first address
	MOVD	p_len+32(FP), R3                             // message length
	VLD1.P	16(R0), [V0.S4]
	FMOVS	(R0), F20
	SUB	$16, R0, R0

blockloop:

	VLD1.P	16(R1), [V4.B16]                             // load message
	VLD1.P	16(R1), [V5.B16]
	VLD1.P	16(R1), [V6.B16]
	VLD1.P	16(R1), [V7.B16]
	VLD1	(R2), [V19.S4]                               // load constant k0-k79
	VMOV	V0.B16, V2.B16
	VMOV	V20.S[0], V1
	VMOV	V2.B16, V3.B16
	VDUP	V19.S[0], V17.S4
	VREV32	V4.B16, V4.B16                               // prepare for using message in Byte format
	VREV32	V5.B16, V5.B16
	VREV32	V6.B16, V6.B16
	VREV32	V7.B16, V7.B16


	VDUP	V19.S[1], V18.S4
	VADD	V17.S4, V4.S4, V16.S4
	SHA1SU0	V6.S4, V5.S4, V4.S4
	HASHUPDATECHOOSE
	SHA1SU1	V7.S4, V4.S4

	VADD	V17.S4, V5.S4, V16.S4
	SHA1SU0	V7.S4, V6.S4, V5.S4
	HASHUPDATECHOOSE
	SHA1SU1	V4.S4, V5.S4
	VADD	V17.S4, V6.S4, V16.S4
	SHA1SU0	V4.S4, V7.S4, V6.S4
	HASHUPDATECHOOSE
	SHA1SU1	V5.S4, V6.S4

	VADD	V17.S4, V7.S4, V16.S4
	SHA1SU0	V5.S4, V4.S4, V7.S4
	HASHUPDATECHOOSE
	SHA1SU1	V6.S4, V7.S4

	VADD	V17.S4, V4.S4, V16.S4
	SHA1SU0	V6.S4, V5.S4, V4.S4
	HASHUPDATECHOOSE
	SHA1SU1	V7.S4, V4.S4

	VDUP	V19.S[2], V17.S4
	VADD	V18.S4, V5.S4, V16.S4
	SHA1SU0	V7.S4, V6.S4, V5.S4
	HASHUPDATEPARITY
	SHA1SU1	V4.S4, V5.S4

	VADD	V18.S4, V6.S4, V16.S4
	SHA1SU0	V4.S4, V7.S4, V6.S4
	HASHUPDATEPARITY
	SHA1SU1	V5.S4, V6.S4

	VADD	V18.S4, V7.S4, V16.S4
	SHA1SU0	V5.S4, V4.S4, V7.S4
	HASHUPDATEPARITY
	SHA1SU1	V6.S4, V7.S4

	VADD	V18.S4, V4.S4, V16.S4
	SHA1SU0	V6.S4, V5.S4, V4.S4
	HASHUPDATEPARITY
	SHA1SU1	V7.S4, V4.S4

	VADD	V18.S4, V5.S4, V16.S4
	SHA1SU0	V7.S4, V6.S4, V5.S4
	HASHUPDATEPARITY
	SHA1SU1	V4.S4, V5.S4

	VDUP	V19.S[3], V18.S4
	VADD	V17.S4, V6.S4, V16.S4
	SHA1SU0	V4.S4, V7.S4, V6.S4
	HASHUPDATEMAJ
	SHA1SU1	V5.S4, V6.S4

	VADD	V17.S4, V7.S4, V16.S4
	SHA1SU0	V5.S4, V4.S4, V7.S4
	HASHUPDATEMAJ
	SHA1SU1	V6.S4, V7.S4

	VADD	V17.S4, V4.S4, V16.S4
	SHA1SU0	V6.S4, V5.S4, V4.S4
	HASHUPDATEMAJ
	SHA1SU1	V7.S4, V4.S4

	VADD	V17.S4, V5.S4, V16.S4
	SHA1SU0	V7.S4, V6.S4, V5.S4
	HASHUPDATEMAJ
	SHA1SU1	V4.S4, V5.S4

	VADD	V17.S4, V6.S4, V16.S4
	SHA1SU0	V4.S4, V7.S4, V6.S4
	HASHUPDATEMAJ
	SHA1SU1	V5.S4, V6.S4

	VADD	V18.S4, V7.S4, V16.S4
	SHA1SU0	V5.S4, V4.S4, V7.S4
	HASHUPDATEPARITY
	SHA1SU1	V6.S4, V7.S4

	VADD	V18.S4, V4.S4, V16.S4
	HASHUPDATEPARITY

	VADD	V18.S4, V5.S4, V16.S4
	HASHUPDATEPARITY

	VADD	V18.S4, V6.S4, V16.S4
	HASHUPDATEPARITY

	VADD	V18.S4, V7.S4, V16.S4
	HASHUPDATEPARITY

	SUB	$64, R3, R3                                  // message length - 64bytes, then compare with 64bytes
	VADD	V2.S4, V0.S4, V0.S4
	VADD	V1.S4, V20.S4, V20.S4
	CBNZ	R3, blockloop

sha1ret:

	VST1.P	[V0.S4], 16(R0)                               // store hash value H(dcba)
	FMOVS	F20, (R0)                                     // store hash value H(e)
	RET
