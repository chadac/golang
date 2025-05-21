// run

// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package main

//golang:noinline
func f32(_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, x int32) uint64 {
	return uint64(uint32(x))
}

//golang:noinline
func f16(_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, x int16) uint64 {
	return uint64(uint16(x))
}

//golang:noinline
func f8(_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, x int8) uint64 {
	return uint64(uint8(x))
}

//golang:noinline
func g32(_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, x uint32) int64 {
	return int64(int32(x))
}

//golang:noinline
func g16(_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, x uint16) int64 {
	return int64(int16(x))
}

//golang:noinline
func g8(_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, x uint8) int64 {
	return int64(int8(x))
}

func main() {
	if golangt := f32(0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, -1); golangt != 0xffffffff {
		println("bad f32", golangt)
	}
	if golangt := f16(0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, -1); golangt != 0xffff {
		println("bad f16", golangt)
	}
	if golangt := f8(0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, -1); golangt != 0xff {
		println("bad f8", golangt)
	}
	if golangt := g32(0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xffffffff); golangt != -1 {
		println("bad g32", golangt)
	}
	if golangt := g16(0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xffff); golangt != -1 {
		println("bad g16", golangt)
	}
	if golangt := g8(0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xff); golangt != -1 {
		println("bad g8", golangt)
	}
}
