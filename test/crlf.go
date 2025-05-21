// runoutput

// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Test source files and strings containing \r and \r\n.

package main

import (
	"fmt"
	"strings"
)

func main() {
	prog = strings.Replace(prog, "BQ", "`", -1)
	prog = strings.Replace(prog, "CR", "\r", -1)
	fmt.Print(prog)
}

var prog = `
package main
CR

import "fmt"

var CR s = "hello\n" + CR
	" world"CR

var t = BQhelloCR
 worldBQ

var u = BQhCReCRlCRlCRoCR
 worldBQ

var golanglden = "hello\n world"

func main() {
	if s != golanglden {
		fmt.Printf("s=%q, want %q", s, golanglden)
	}
	if t != golanglden {
		fmt.Printf("t=%q, want %q", t, golanglden)
	}
	if u != golanglden {
		fmt.Printf("u=%q, want %q", u, golanglden)
	}
}
`
