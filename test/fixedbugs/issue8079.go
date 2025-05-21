// compile

// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Issue 8079: gccgolang crashes when compiling interface with blank type name.

package p

type _ interface{}
