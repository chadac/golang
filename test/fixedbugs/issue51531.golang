// errorcheck -lang=golang1.17

// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package p

type empty interface{}

type Foo[T empty] int // ERROR "type parameter requires golang1\.18 or later \(-lang was set to golang1\.17; check golang.mod\)"

func Bar[T empty]() {} // ERROR "type parameter requires golang1\.18 or later \(-lang was set to golang1\.17; check golang.mod\)"
