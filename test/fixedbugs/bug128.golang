// run

// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package main
func main() {
	switch {
		// empty switch is allowed according to syntax
		// unclear why it shouldn't be allowed
	}
	switch tag := 0; tag {
		// empty switch is allowed according to syntax
		// unclear why it shouldn't be allowed
	}
}

/*
uetli:~/Source/golang1/test/bugs gri$ 6g bug127.golang 
bug127.golang:5: switch statement must have case labels
bug127.golang:9: switch statement must have case labels
*/
