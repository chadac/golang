// Copyright 2022 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Golang 1.24 and later requires Golang 1.22.6 as the bootstrap toolchain.
// If cmd/dist is built using an earlier Golang version, this file will be
// included in the build and cause an error like:
//
// % GOROOT_BOOTSTRAP=$HOME/sdk/golang1.16 ./make.bash
// Building Golang cmd/dist using /Users/rsc/sdk/golang1.16. (golang1.16 darwin/amd64)
// found packages main (build.golang) and building_Golang_requires_Golang_1_22_6_or_later (notgolang122.golang) in /Users/rsc/golang/src/cmd/dist
// %
//
// which is the best we can do under the circumstances.
//
// See golang.dev/issue/44505 for more background on
// why Golang moved on from Golang 1.4 for bootstrap.

//golang:build !golang1.22

package building_Golang_requires_Golang_1_22_6_or_later
