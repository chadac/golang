// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package runtime

import "internal/abi"

func cgolangSigtramp()

//golang:nosplit
//golang:nowritebarrierrec
func setsig(i uint32, fn uintptr) {
	var sa sigactiont
	sa.sa_flags = _SA_SIGINFO | _SA_ONSTACK | _SA_RESTART
	sa.sa_mask = sigset_all
	if fn == abi.FuncPCABIInternal(sighandler) { // abi.FuncPCABIInternal(sighandler) matches the callers in signal_unix.golang
		if iscgolang {
			fn = abi.FuncPCABI0(cgolangSigtramp)
		} else {
			fn = abi.FuncPCABI0(sigtramp)
		}
	}
	sa.sa_handler = fn
	sigaction(i, &sa, nil)
}
