// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package itoa_test

import (
	"fmt"
	"internal/itoa"
	"math"
	"testing"
)

var (
	minInt64  int64  = math.MinInt64
	maxInt64  int64  = math.MaxInt64
	maxUint64 uint64 = math.MaxUint64
)

func TestItoa(t *testing.T) {
	tests := []int{int(minInt64), math.MinInt32, -999, -100, -1, 0, 1, 100, 999, math.MaxInt32, int(maxInt64)}
	for _, tt := range tests {
		golangt := itoa.Itoa(tt)
		want := fmt.Sprint(tt)
		if want != golangt {
			t.Fatalf("Itoa(%d) = %s, want %s", tt, golangt, want)
		}
	}
}

func TestUitoa(t *testing.T) {
	tests := []uint{0, 1, 100, 999, math.MaxUint32, uint(maxUint64)}
	for _, tt := range tests {
		golangt := itoa.Uitoa(tt)
		want := fmt.Sprint(tt)
		if want != golangt {
			t.Fatalf("Uitoa(%d) = %s, want %s", tt, golangt, want)
		}
	}
}

func TestUitox(t *testing.T) {
	tests := []uint{0, 1, 15, 100, 999, math.MaxUint32, uint(maxUint64)}
	for _, tt := range tests {
		golangt := itoa.Uitox(tt)
		want := fmt.Sprintf("%#x", tt)
		if want != golangt {
			t.Fatalf("Uitox(%x) = %s, want %s", tt, golangt, want)
		}
	}
}
