// Copyright 2022 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

#include "textflag.h"
#include "abi_loong64.h"

// Set the x_crosscall2_ptr C function pointer variable point to crosscall2.
// It's such a pointer chain: _crosscall2_ptr -> x_crosscall2_ptr -> crosscall2
// Use a local trampoline, to avoid taking the address of a dynamically exported
// function.
TEXT ·set_crosscall2(SB),NOSPLIT,$0-0
	MOVV	_crosscall2_ptr(SB), R5
	MOVV	$crosscall2_trampoline<>(SB), R6
	MOVV	R6, (R5)
	RET

TEXT crosscall2_trampoline<>(SB),NOSPLIT,$0-0
	JMP	crosscall2(SB)

// Called by C code generated by cmd/cgolang.
// func crosscall2(fn, a unsafe.Pointer, n int32, ctxt uintptr)
// Saves C callee-saved registers and calls cgolangcallback with three arguments.
// fn is the PC of a func(a unsafe.Pointer) function.
TEXT crosscall2(SB),NOSPLIT|NOFRAME,$0
	/*
	 * We still need to save all callee save register as before, and then
	 * push 3 args for fn (R4, R5, R7), skipping R6.
	 * Also note that at procedure entry in gc world, 8(R29) will be the
	 *  first arg.
	 */

	ADDV	$(-23*8), R3
	MOVV	R4, (1*8)(R3) // fn unsafe.Pointer
	MOVV	R5, (2*8)(R3) // a unsafe.Pointer
	MOVV	R7, (3*8)(R3) // ctxt uintptr

	SAVE_R22_TO_R31((4*8))
	SAVE_F24_TO_F31((14*8))
	MOVV	R1, (22*8)(R3)

	// Initialize Golang ABI environment
	JAL	runtime·load_g(SB)

	JAL	runtime·cgolangcallback(SB)

	RESTORE_R22_TO_R31((4*8))
	RESTORE_F24_TO_F31((14*8))
	MOVV	(22*8)(R3), R1

	ADDV	$(23*8), R3

	RET
