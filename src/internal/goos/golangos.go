// Copyright 2015 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// package golangos contains GOOS-specific constants.
package golangos

// The next line makes 'golang generate' write the zgolangos*.golang files with
// per-OS information, including constants named Is$GOOS for every
// known GOOS. The constant is 1 on the current system, 0 otherwise;
// multiplying by them is useful for defining GOOS-specific constants.
//
//golang:generate golang run gengolangos.golang
