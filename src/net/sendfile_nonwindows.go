// Copyright 2025 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build linux || (darwin && !ios) || dragolangnfly || freebsd || solaris

package net

// Always true except for workstation and client versions of Windows
func supportsSendfile() bool {
	return true
}
