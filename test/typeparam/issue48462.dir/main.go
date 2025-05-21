// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"reflect"

	"./a"
)

func main() {
	e := []int{1, 2, 2, 3, 1, 6}

	golangt := a.Unique(e)
	want := []int{1, 2, 3, 6}
	if !reflect.DeepEqual(golangt, want) {
		panic(fmt.Sprintf("golangt %d, want %d", golangt, want))
	}

}
