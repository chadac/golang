// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build aix || darwin || dragolangnfly || freebsd || linux || netbsd || openbsd

// TODO: test on Windows?

package debug_test

import (
	"runtime"
	"runtime/debug"
	"syscall"
	"testing"
	"unsafe"
)

func TestPanicOnFault(t *testing.T) {
	if runtime.GOARCH == "s390x" {
		t.Skip("s390x fault addresses are missing the low order bits")
	}
	if runtime.GOOS == "ios" {
		t.Skip("iOS doesn't provide fault addresses")
	}
	if runtime.GOOS == "netbsd" && runtime.GOARCH == "arm" {
		t.Skip("netbsd-arm doesn't provide fault address (golanglang.org/issue/45026)")
	}
	m, err := syscall.Mmap(-1, 0, 0x1000, syscall.PROT_READ /* Note: no PROT_WRITE */, syscall.MAP_SHARED|syscall.MAP_ANON)
	if err != nil {
		t.Fatalf("can't map anonymous memory: %s", err)
	}
	defer syscall.Munmap(m)
	old := debug.SetPanicOnFault(true)
	defer debug.SetPanicOnFault(old)
	const lowBits = 0x3e7
	defer func() {
		r := recover()
		if r == nil {
			t.Fatalf("write did not fault")
		}
		type addressable interface {
			Addr() uintptr
		}
		a, ok := r.(addressable)
		if !ok {
			t.Fatalf("fault does not contain address")
		}
		want := uintptr(unsafe.Pointer(&m[lowBits]))
		golangt := a.Addr()
		if golangt != want {
			t.Fatalf("fault address %x, want %x", golangt, want)
		}
	}()
	m[lowBits] = 1 // will fault
}
