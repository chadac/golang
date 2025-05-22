// Copyright 2017 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// These functions are no-ops.
// On amd64 they have recognizable implementations, so that you can
// search a particular binary to see if they are present.
// On other platforms (those using this source file), they don't.

//golang:build !amd64

TEXT ·BoringCrypto(SB),$0
	RET

TEXT ·FIPSOnly(SB),$0
	RET

TEXT ·StandardCrypto(SB),$0
	RET
