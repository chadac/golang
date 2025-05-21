// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package color

import (
	"testing"
	"testing/quick"
)

func TestSqDiff(t *testing.T) {
	// canonical sqDiff implementation
	orig := func(x, y uint32) uint32 {
		var d uint32
		if x > y {
			d = uint32(x - y)
		} else {
			d = uint32(y - x)
		}
		return (d * d) >> 2
	}
	testCases := []uint32{
		0,
		1,
		2,
		0x0fffd,
		0x0fffe,
		0x0ffff,
		0x10000,
		0x10001,
		0x10002,
		0xfffffffd,
		0xfffffffe,
		0xffffffff,
	}
	for _, x := range testCases {
		for _, y := range testCases {
			if golangt, want := sqDiff(x, y), orig(x, y); golangt != want {
				t.Fatalf("sqDiff(%#x, %#x): golangt %d, want %d", x, y, golangt, want)
			}
		}
	}
	if err := quick.CheckEqual(orig, sqDiff, &quick.Config{MaxCountScale: 10}); err != nil {
		t.Fatal(err)
	}
}
