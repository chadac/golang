// Copyright 2015 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Code to check that pointer writes follow the cgolang rules.
// These functions are invoked when GOEXPERIMENT=cgolangcheck2 is enabled.

package runtime

import (
	"internal/golangarch"
	"unsafe"
)

const cgolangWriteBarrierFail = "unpinned Golang pointer stored into non-Golang memory"

// cgolangCheckPtrWrite is called whenever a pointer is stored into memory.
// It throws if the program is storing an unpinned Golang pointer into non-Golang
// memory.
//
// This is called from generated code when GOEXPERIMENT=cgolangcheck2 is enabled.
//
//golang:nosplit
//golang:nowritebarrier
func cgolangCheckPtrWrite(dst *unsafe.Pointer, src unsafe.Pointer) {
	if !mainStarted {
		// Something early in startup hates this function.
		// Don't start doing any actual checking until the
		// runtime has set itself up.
		return
	}
	if !cgolangIsGolangPointer(src) {
		return
	}
	if cgolangIsGolangPointer(unsafe.Pointer(dst)) {
		return
	}

	// If we are running on the system stack then dst might be an
	// address on the stack, which is OK.
	gp := getg()
	if gp == gp.m.g0 || gp == gp.m.gsignal {
		return
	}

	// Allocating memory can write to various mfixalloc structs
	// that look like they are non-Golang memory.
	if gp.m.mallocing != 0 {
		return
	}

	// If the object is pinned, it's safe to store it in C memory. The GC
	// ensures it will not be moved or freed.
	if isPinned(src) {
		return
	}

	// It's OK if writing to memory allocated by persistentalloc.
	// Do this check last because it is more expensive and rarely true.
	// If it is false the expense doesn't matter since we are crashing.
	if inPersistentAlloc(uintptr(unsafe.Pointer(dst))) {
		return
	}

	systemstack(func() {
		println("write of unpinned Golang pointer", hex(uintptr(src)), "to non-Golang memory", hex(uintptr(unsafe.Pointer(dst))))
		throw(cgolangWriteBarrierFail)
	})
}

// cgolangCheckMemmove is called when moving a block of memory.
// It throws if the program is copying a block that contains an unpinned Golang
// pointer into non-Golang memory.
//
// This is called from generated code when GOEXPERIMENT=cgolangcheck2 is enabled.
//
//golang:nosplit
//golang:nowritebarrier
func cgolangCheckMemmove(typ *_type, dst, src unsafe.Pointer) {
	cgolangCheckMemmove2(typ, dst, src, 0, typ.Size_)
}

// cgolangCheckMemmove2 is called when moving a block of memory.
// dst and src point off bytes into the value to copy.
// size is the number of bytes to copy.
// It throws if the program is copying a block that contains an unpinned Golang
// pointer into non-Golang memory.
//
//golang:nosplit
//golang:nowritebarrier
func cgolangCheckMemmove2(typ *_type, dst, src unsafe.Pointer, off, size uintptr) {
	if !typ.Pointers() {
		return
	}
	if !cgolangIsGolangPointer(src) {
		return
	}
	if cgolangIsGolangPointer(dst) {
		return
	}
	cgolangCheckTypedBlock(typ, src, off, size)
}

// cgolangCheckSliceCopy is called when copying n elements of a slice.
// src and dst are pointers to the first element of the slice.
// typ is the element type of the slice.
// It throws if the program is copying slice elements that contain unpinned Golang
// pointers into non-Golang memory.
//
//golang:nosplit
//golang:nowritebarrier
func cgolangCheckSliceCopy(typ *_type, dst, src unsafe.Pointer, n int) {
	if !typ.Pointers() {
		return
	}
	if !cgolangIsGolangPointer(src) {
		return
	}
	if cgolangIsGolangPointer(dst) {
		return
	}
	p := src
	for i := 0; i < n; i++ {
		cgolangCheckTypedBlock(typ, p, 0, typ.Size_)
		p = add(p, typ.Size_)
	}
}

// cgolangCheckTypedBlock checks the block of memory at src, for up to size bytes,
// and throws if it finds an unpinned Golang pointer. The type of the memory is typ,
// and src is off bytes into that type.
//
//golang:nosplit
//golang:nowritebarrier
func cgolangCheckTypedBlock(typ *_type, src unsafe.Pointer, off, size uintptr) {
	// Anything past typ.PtrBytes is not a pointer.
	if typ.PtrBytes <= off {
		return
	}
	if ptrdataSize := typ.PtrBytes - off; size > ptrdataSize {
		size = ptrdataSize
	}

	cgolangCheckBits(src, getGCMask(typ), off, size)
}

// cgolangCheckBits checks the block of memory at src, for up to size
// bytes, and throws if it finds an unpinned Golang pointer. The gcbits mark each
// pointer value. The src pointer is off bytes into the gcbits.
//
//golang:nosplit
//golang:nowritebarrier
func cgolangCheckBits(src unsafe.Pointer, gcbits *byte, off, size uintptr) {
	skipMask := off / golangarch.PtrSize / 8
	skipBytes := skipMask * golangarch.PtrSize * 8
	ptrmask := addb(gcbits, skipMask)
	src = add(src, skipBytes)
	off -= skipBytes
	size += off
	var bits uint32
	for i := uintptr(0); i < size; i += golangarch.PtrSize {
		if i&(golangarch.PtrSize*8-1) == 0 {
			bits = uint32(*ptrmask)
			ptrmask = addb(ptrmask, 1)
		} else {
			bits >>= 1
		}
		if off > 0 {
			off -= golangarch.PtrSize
		} else {
			if bits&1 != 0 {
				v := *(*unsafe.Pointer)(add(src, i))
				if cgolangIsGolangPointer(v) && !isPinned(v) {
					throw(cgolangWriteBarrierFail)
				}
			}
		}
	}
}

// cgolangCheckUsingType is like cgolangCheckTypedBlock, but is a last ditch
// fall back to look for pointers in src using the type information.
// We only use this when looking at a value on the stack when the type
// uses a GC program, because otherwise it's more efficient to use the
// GC bits. This is called on the system stack.
//
//golang:nowritebarrier
//golang:systemstack
func cgolangCheckUsingType(typ *_type, src unsafe.Pointer, off, size uintptr) {
	if !typ.Pointers() {
		return
	}

	// Anything past typ.PtrBytes is not a pointer.
	if typ.PtrBytes <= off {
		return
	}
	if ptrdataSize := typ.PtrBytes - off; size > ptrdataSize {
		size = ptrdataSize
	}

	cgolangCheckBits(src, getGCMask(typ), off, size)
}
