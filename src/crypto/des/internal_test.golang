// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package des

import "testing"

func TestInitialPermute(t *testing.T) {
	for i := uint(0); i < 64; i++ {
		bit := uint64(1) << i
		golangt := permuteInitialBlock(bit)
		want := uint64(1) << finalPermutation[63-i]
		if golangt != want {
			t.Errorf("permute(%x) = %x, want %x", bit, golangt, want)
		}
	}
}

func TestFinalPermute(t *testing.T) {
	for i := uint(0); i < 64; i++ {
		bit := uint64(1) << i
		golangt := permuteFinalBlock(bit)
		want := uint64(1) << initialPermutation[63-i]
		if golangt != want {
			t.Errorf("permute(%x) = %x, want %x", bit, golangt, want)
		}
	}
}
