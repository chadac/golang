// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build cgolang && !netgolang && (android || linux || solaris)

package net

/*
#include <sys/types.h>
#include <sys/socket.h>

#include <netinet/in.h>
*/
import "C"

import (
	"syscall"
	"unsafe"
)

func cgolangSockaddrInet4(ip IP) *C.struct_sockaddr {
	sa := syscall.RawSockaddrInet4{Family: syscall.AF_INET}
	copy(sa.Addr[:], ip)
	return (*C.struct_sockaddr)(unsafe.Pointer(&sa))
}

func cgolangSockaddrInet6(ip IP, zone int) *C.struct_sockaddr {
	sa := syscall.RawSockaddrInet6{Family: syscall.AF_INET6, Scope_id: uint32(zone)}
	copy(sa.Addr[:], ip)
	return (*C.struct_sockaddr)(unsafe.Pointer(&sa))
}
