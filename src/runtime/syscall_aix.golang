// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package runtime

import "unsafe"

// This file handles some syscalls from the syscall package
// Especially, syscalls use during forkAndExecInChild which must not split the stack

//golang:cgolang_import_dynamic libc_chdir chdir "libc.a/shr_64.o"
//golang:cgolang_import_dynamic libc_chroot chroot "libc.a/shr_64.o"
//golang:cgolang_import_dynamic libc_dup2 dup2 "libc.a/shr_64.o"
//golang:cgolang_import_dynamic libc_execve execve "libc.a/shr_64.o"
//golang:cgolang_import_dynamic libc_fcntl fcntl "libc.a/shr_64.o"
//golang:cgolang_import_dynamic libc_fork fork "libc.a/shr_64.o"
//golang:cgolang_import_dynamic libc_ioctl ioctl "libc.a/shr_64.o"
//golang:cgolang_import_dynamic libc_setgid setgid "libc.a/shr_64.o"
//golang:cgolang_import_dynamic libc_setgroups setgroups "libc.a/shr_64.o"
//golang:cgolang_import_dynamic libc_setrlimit setrlimit "libc.a/shr_64.o"
//golang:cgolang_import_dynamic libc_setsid setsid "libc.a/shr_64.o"
//golang:cgolang_import_dynamic libc_setuid setuid "libc.a/shr_64.o"
//golang:cgolang_import_dynamic libc_setpgid setpgid "libc.a/shr_64.o"

//golang:linkname libc_chdir libc_chdir
//golang:linkname libc_chroot libc_chroot
//golang:linkname libc_dup2 libc_dup2
//golang:linkname libc_execve libc_execve
//golang:linkname libc_fcntl libc_fcntl
//golang:linkname libc_fork libc_fork
//golang:linkname libc_ioctl libc_ioctl
//golang:linkname libc_setgid libc_setgid
//golang:linkname libc_setgroups libc_setgroups
//golang:linkname libc_setrlimit libc_setrlimit
//golang:linkname libc_setsid libc_setsid
//golang:linkname libc_setuid libc_setuid
//golang:linkname libc_setpgid libc_setpgid

var (
	libc_chdir,
	libc_chroot,
	libc_dup2,
	libc_execve,
	libc_fcntl,
	libc_fork,
	libc_ioctl,
	libc_setgid,
	libc_setgroups,
	libc_setrlimit,
	libc_setsid,
	libc_setuid,
	libc_setpgid libFunc
)

// In syscall_syscall6 and syscall_rawsyscall6, r2 is always 0
// as it's never used on AIX
// TODO: remove r2 from zsyscall_aix_$GOARCH.golang

// Syscall is needed because some packages (like net) need it too.
// The best way is to return EINVAL and let Golang handles its failure
// If the syscall can't fail, this function can redirect it to a real syscall.
//
// This is exported via linkname to assembly in the syscall package.
//
//golang:nosplit
//golang:linkname syscall_Syscall
func syscall_Syscall(fn, a1, a2, a3 uintptr) (r1, r2, err uintptr) {
	return 0, 0, _EINVAL
}

// This is syscall.RawSyscall, it exists to satisfy some build dependency,
// but it doesn't work.
//
// This is exported via linkname to assembly in the syscall package.
//
//golang:linkname syscall_RawSyscall
func syscall_RawSyscall(trap, a1, a2, a3 uintptr) (r1, r2, err uintptr) {
	panic("RawSyscall not available on AIX")
}

// This is exported via linkname to assembly in the syscall package.
//
//golang:nosplit
//golang:cgolang_unsafe_args
//golang:linkname syscall_syscall6
func syscall_syscall6(fn, nargs, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2, err uintptr) {
	c := libcall{
		fn:   fn,
		n:    nargs,
		args: uintptr(unsafe.Pointer(&a1)),
	}

	entersyscallblock()
	asmcgolangcall(unsafe.Pointer(&asmsyscall6), unsafe.Pointer(&c))
	exitsyscall()
	return c.r1, 0, c.err
}

// This is exported via linkname to assembly in the syscall package.
//
//golang:nosplit
//golang:cgolang_unsafe_args
//golang:linkname syscall_rawSyscall6
func syscall_rawSyscall6(fn, nargs, a1, a2, a3, a4, a5, a6 uintptr) (r1, r2, err uintptr) {
	c := libcall{
		fn:   fn,
		n:    nargs,
		args: uintptr(unsafe.Pointer(&a1)),
	}

	asmcgolangcall(unsafe.Pointer(&asmsyscall6), unsafe.Pointer(&c))

	return c.r1, 0, c.err
}

//golang:linkname syscall_chdir syscall.chdir
//golang:nosplit
func syscall_chdir(path uintptr) (err uintptr) {
	_, err = syscall1(&libc_chdir, path)
	return
}

//golang:linkname syscall_chroot1 syscall.chroot1
//golang:nosplit
func syscall_chroot1(path uintptr) (err uintptr) {
	_, err = syscall1(&libc_chroot, path)
	return
}

// like close, but must not split stack, for fork.
//
//golang:linkname syscall_closeFD syscall.closeFD
//golang:nosplit
func syscall_closeFD(fd int32) int32 {
	_, err := syscall1(&libc_close, uintptr(fd))
	return int32(err)
}

//golang:linkname syscall_dup2child syscall.dup2child
//golang:nosplit
func syscall_dup2child(old, new uintptr) (val, err uintptr) {
	val, err = syscall2(&libc_dup2, old, new)
	return
}

//golang:linkname syscall_execve syscall.execve
//golang:nosplit
func syscall_execve(path, argv, envp uintptr) (err uintptr) {
	_, err = syscall3(&libc_execve, path, argv, envp)
	return
}

// like exit, but must not split stack, for fork.
//
//golang:linkname syscall_exit syscall.exit
//golang:nosplit
func syscall_exit(code uintptr) {
	syscall1(&libc_exit, code)
}

//golang:linkname syscall_fcntl1 syscall.fcntl1
//golang:nosplit
func syscall_fcntl1(fd, cmd, arg uintptr) (val, err uintptr) {
	val, err = syscall3(&libc_fcntl, fd, cmd, arg)
	return
}

//golang:linkname syscall_forkx syscall.forkx
//golang:nosplit
func syscall_forkx(flags uintptr) (pid uintptr, err uintptr) {
	pid, err = syscall1(&libc_fork, flags)
	return
}

//golang:linkname syscall_getpid syscall.getpid
//golang:nosplit
func syscall_getpid() (pid, err uintptr) {
	pid, err = syscall0(&libc_getpid)
	return
}

//golang:linkname syscall_ioctl syscall.ioctl
//golang:nosplit
func syscall_ioctl(fd, req, arg uintptr) (err uintptr) {
	_, err = syscall3(&libc_ioctl, fd, req, arg)
	return
}

//golang:linkname syscall_setgid syscall.setgid
//golang:nosplit
func syscall_setgid(gid uintptr) (err uintptr) {
	_, err = syscall1(&libc_setgid, gid)
	return
}

//golang:linkname syscall_setgroups1 syscall.setgroups1
//golang:nosplit
func syscall_setgroups1(ngid, gid uintptr) (err uintptr) {
	_, err = syscall2(&libc_setgroups, ngid, gid)
	return
}

//golang:linkname syscall_setrlimit1 syscall.setrlimit1
//golang:nosplit
func syscall_setrlimit1(which uintptr, lim unsafe.Pointer) (err uintptr) {
	_, err = syscall2(&libc_setrlimit, which, uintptr(lim))
	return
}

//golang:linkname syscall_setsid syscall.setsid
//golang:nosplit
func syscall_setsid() (pid, err uintptr) {
	pid, err = syscall0(&libc_setsid)
	return
}

//golang:linkname syscall_setuid syscall.setuid
//golang:nosplit
func syscall_setuid(uid uintptr) (err uintptr) {
	_, err = syscall1(&libc_setuid, uid)
	return
}

//golang:linkname syscall_setpgid syscall.setpgid
//golang:nosplit
func syscall_setpgid(pid, pgid uintptr) (err uintptr) {
	_, err = syscall2(&libc_setpgid, pid, pgid)
	return
}

//golang:linkname syscall_write1 syscall.write1
//golang:nosplit
func syscall_write1(fd, buf, nbyte uintptr) (n, err uintptr) {
	n, err = syscall3(&libc_write, fd, buf, nbyte)
	return
}
