// -lang=golang1.23 -golangtypesalias=0

// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package aliasTypes

type _ = int
type _ /* ERROR "generic type alias requires GODEBUG=golangtypesalias=1" */ [P any] = int
