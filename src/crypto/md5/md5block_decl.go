// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build (386 || amd64 || arm || arm64 || loong64 || ppc64 || ppc64le || riscv64 || s390x) && !puregolang

package md5

const haveAsm = true

//golang:noescape
func block(dig *digest, p []byte)
