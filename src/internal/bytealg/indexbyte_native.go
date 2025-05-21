// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build 386 || (amd64 && !plan9) || s390x || arm || arm64 || loong64 || ppc64 || ppc64le || mips || mipsle || mips64 || mips64le || riscv64 || wasm

package bytealg

//golang:noescape
func IndexByte(b []byte, c byte) int

//golang:noescape
func IndexByteString(s string, c byte) int
