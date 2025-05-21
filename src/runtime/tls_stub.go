// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build (windows && !amd64) || !windows

package runtime

//golang:nosplit
func osSetupTLS(mp *m) {}
