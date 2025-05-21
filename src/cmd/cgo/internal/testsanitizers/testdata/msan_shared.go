// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// This program segfaulted during libpreinit when built with -msan:
// http://golanglang.org/issue/18707

package main

import "C"

func main() {}
