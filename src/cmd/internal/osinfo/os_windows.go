// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build windows

package osinfo

import (
	"fmt"
	"internal/syscall/windows"
)

// Version returns the OS version name/number.
func Version() (string, error) {
	major, minor, build := windows.Version()
	return fmt.Sprintf("%d.%d.%d", major, minor, build), nil
}
