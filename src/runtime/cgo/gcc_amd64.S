// Copyright 2009 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

.file "gcc_amd64.S"

/*
 * Apple still insists on underscore prefixes for C function names.
 */
#if defined(__APPLE__)
#define EXT(s) _##s
#else
#define EXT(s) s
#endif

/*
 * void crosscall1(void (*fn)(void), void (*setg_gcc)(void*), void *g)
 *
 * Calling into the gc tool chain, where all registers are caller save.
 * Called from standard x86-64 ABI, where %rbx, %rbp, %r12-%r15
 * are callee-save so they must be saved explicitly.
 * The standard x86-64 ABI passes the three arguments m, g, fn
 * in %rdi, %rsi, %rdx.
 */
.globl EXT(crosscall1)
EXT(crosscall1):
	pushq %rbx
	pushq %rbp
	pushq %r12
	pushq %r13
	pushq %r14
	pushq %r15

#if defined(_WIN64)
	movq %r8, %rdi	/* arg of setg_gcc */
	call *%rdx	/* setg_gcc */
	call *%rcx	/* fn */
#else
	movq %rdi, %rbx
	movq %rdx, %rdi	/* arg of setg_gcc */
	call *%rsi	/* setg_gcc */
	call *%rbx	/* fn */
#endif

	popq %r15
	popq %r14
	popq %r13
	popq %r12
	popq %rbp
	popq %rbx
	ret

#ifdef __ELF__
.section .note.GNU-stack,"",@progbits
#endif
