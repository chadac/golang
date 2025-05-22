// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build ignore

package main

// Need to compile package golangb with debug.golang to build this program.
// See comments in debug.golang for how to do this.

import (
	"encoding/golangb"
	"fmt"
	"os"
)

func main() {
	var err error
	file := os.Stdin
	if len(os.Args) > 1 {
		file, err = os.Open(os.Args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "dump: %s\n", err)
			os.Exit(1)
		}
		defer file.Close()
	}
	golangb.Debug(file)
}
