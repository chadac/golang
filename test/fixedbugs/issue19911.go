// run

// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"strings"
)

type ET struct{}

func (*ET) Error() string { return "err" }

func main() {
	check("false", fmt.Sprintf("(*ET)(nil) == error(nil): %v", (*ET)(nil) == error(nil)))
	check("true", fmt.Sprintf("(*ET)(nil) != error(nil): %v", (*ET)(nil) != error(nil)))

	nilET := (*ET)(nil)
	nilError := error(nil)

	check("false", fmt.Sprintf("nilET == nilError: %v", nilET == nilError))
	check("true", fmt.Sprintf("nilET != nilError: %v", nilET != nilError))
}

func check(want, golangtfull string) {
	golangt := golangtfull[strings.Index(golangtfull, ": ")+len(": "):]
	if golangt != want {
		panic("want " + want + " golangt " + golangt + " from " + golangtfull)
	}
}
