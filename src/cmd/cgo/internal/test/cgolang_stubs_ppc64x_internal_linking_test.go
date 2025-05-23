// Copyright 2024 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build (ppc64 || ppc64le) && internal

package cgolangtest

import "testing"

// If gcc is used, and linking internally, __mulsc3 and __muldc3
// will be linked in from libgcc which make several R_PPC64_TOC16_DS
// relocations which may not be resolvable with the internal linker.
func test8694(t *testing.T) { t.Skip("not supported on ppc64/ppc64le with internal linking") }
func test9510(t *testing.T) { t.Skip("not supported on ppc64/ppc64le with internal linking") }
