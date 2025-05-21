// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build !linux

package net

import (
	"context"
)

func (sd *sysDialer) dialMPTCP(ctx context.Context, laddr, raddr *TCPAddr) (*TCPConn, error) {
	return sd.dialTCP(ctx, laddr, raddr)
}

func (sl *sysListener) listenMPTCP(ctx context.Context, laddr *TCPAddr) (*TCPListener, error) {
	return sl.listenTCP(ctx, laddr)
}

func isUsingMultipathTCP(fd *netFD) bool {
	return false
}
