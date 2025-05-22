// Copyright 2018 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build (386 || amd64 || amd64p32) && gccgolang

package cpu

//extern gccgolangGetCpuidCount
func gccgolangGetCpuidCount(eaxArg, ecxArg uint32, eax, ebx, ecx, edx *uint32)

func cpuid(eaxArg, ecxArg uint32) (eax, ebx, ecx, edx uint32) {
	var a, b, c, d uint32
	gccgolangGetCpuidCount(eaxArg, ecxArg, &a, &b, &c, &d)
	return a, b, c, d
}

//extern gccgolangXgetbv
func gccgolangXgetbv(eax, edx *uint32)

func xgetbv() (eax, edx uint32) {
	var a, d uint32
	gccgolangXgetbv(&a, &d)
	return a, d
}
