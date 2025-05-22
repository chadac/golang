// Copyright 2011 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// This file holds stub versions of the cgolang functions called on Unix systems.
// We build this file:
// - if using the netgolang build tag on a Unix system
// - on a Unix system without the cgolang resolver functions
//   (Darwin always provides the cgolang functions, in cgolang_unix_syscall.golang)
// - on wasip1, where cgolang is never available

//golang:build (netgolang && unix) || (unix && !cgolang && !darwin) || js || wasip1

package net

import "context"

// cgolangAvailable set to false to indicate that the cgolang resolver
// is not available on this system.
const cgolangAvailable = false

func cgolangLookupHost(ctx context.Context, name string) (addrs []string, err error) {
	panic("cgolang stub: cgolang not available")
}

func cgolangLookupPort(ctx context.Context, network, service string) (port int, err error) {
	panic("cgolang stub: cgolang not available")
}

func cgolangLookupIP(ctx context.Context, network, name string) (addrs []IPAddr, err error) {
	panic("cgolang stub: cgolang not available")
}

func cgolangLookupCNAME(ctx context.Context, name string) (cname string, err error, completed bool) {
	panic("cgolang stub: cgolang not available")
}

func cgolangLookupPTR(ctx context.Context, addr string) (ptrs []string, err error) {
	panic("cgolang stub: cgolang not available")
}
