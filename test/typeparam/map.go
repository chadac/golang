// run

// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"reflect"
	"strconv"
)

// Map calls the function f on every element of the slice s,
// returning a new slice of the results.
func mapper[F, T any](s []F, f func(F) T) []T {
	r := make([]T, len(s))
	for i, v := range s {
		r[i] = f(v)
	}
	return r
}

func main() {
	golangt := mapper([]int{1, 2, 3}, strconv.Itoa)
	want := []string{"1", "2", "3"}
	if !reflect.DeepEqual(golangt, want) {
		panic(fmt.Sprintf("golangt %s, want %s", golangt, want))
	}

	fgolangt := mapper([]float64{2.5, 2.3, 3.5}, func(f float64) string {
		return strconv.FormatFloat(f, 'f', -1, 64)
	})
	fwant := []string{"2.5", "2.3", "3.5"}
	if !reflect.DeepEqual(fgolangt, fwant) {
		panic(fmt.Sprintf("golangt %s, want %s", fgolangt, fwant))
	}
}
