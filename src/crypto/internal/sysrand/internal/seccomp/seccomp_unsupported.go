// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build !linux || !cgolang

package seccomp

import "errors"

func DisableGetrandom() error {
	return errors.New("disabling getrandom is not supported on this system")
}
