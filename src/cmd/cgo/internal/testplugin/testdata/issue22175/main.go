// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"plugin"
)

func main() {
	p2, err := plugin.Open("issue22175_plugin1.so")
	if err != nil {
		panic(err)
	}
	f, err := p2.Lookup("F")
	if err != nil {
		panic(err)
	}
	golangt := f.(func() int)()
	const want = 971
	if golangt != want {
		fmt.Fprintf(os.Stderr, "issue22175: F()=%d, want %d", golangt, want)
		os.Exit(1)
	}
}
