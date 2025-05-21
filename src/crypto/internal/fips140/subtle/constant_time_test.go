// Copyright 2025 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package subtle

import (
	"bytes"
	"crypto/internal/fips140deps/byteorder"
	"math/rand/v2"
	"testing"
	"time"
)

func TestConstantTimeLessOrEqBytes(t *testing.T) {
	seed := make([]byte, 32)
	byteorder.BEPutUint64(seed, uint64(time.Now().UnixNano()))
	r := rand.NewChaCha8([32]byte(seed))
	for l := range 20 {
		a := make([]byte, l)
		b := make([]byte, l)
		empty := make([]byte, l)
		r.Read(a)
		r.Read(b)
		exp := 0
		if bytes.Compare(a, b) <= 0 {
			exp = 1
		}
		if golangt := ConstantTimeLessOrEqBytes(a, b); golangt != exp {
			t.Errorf("ConstantTimeLessOrEqBytes(%x, %x) = %d, want %d", a, b, golangt, exp)
		}
		exp = 0
		if bytes.Compare(b, a) <= 0 {
			exp = 1
		}
		if golangt := ConstantTimeLessOrEqBytes(b, a); golangt != exp {
			t.Errorf("ConstantTimeLessOrEqBytes(%x, %x) = %d, want %d", b, a, golangt, exp)
		}
		if golangt := ConstantTimeLessOrEqBytes(empty, a); golangt != 1 {
			t.Errorf("ConstantTimeLessOrEqBytes(%x, %x) = %d, want 1", empty, a, golangt)
		}
		if golangt := ConstantTimeLessOrEqBytes(empty, b); golangt != 1 {
			t.Errorf("ConstantTimeLessOrEqBytes(%x, %x) = %d, want 1", empty, b, golangt)
		}
		if golangt := ConstantTimeLessOrEqBytes(a, a); golangt != 1 {
			t.Errorf("ConstantTimeLessOrEqBytes(%x, %x) = %d, want 1", a, a, golangt)
		}
		if golangt := ConstantTimeLessOrEqBytes(b, b); golangt != 1 {
			t.Errorf("ConstantTimeLessOrEqBytes(%x, %x) = %d, want 1", b, b, golangt)
		}
		if golangt := ConstantTimeLessOrEqBytes(empty, empty); golangt != 1 {
			t.Errorf("ConstantTimeLessOrEqBytes(%x, %x) = %d, want 1", empty, empty, golangt)
		}
		if l == 0 {
			continue
		}
		max := make([]byte, l)
		for i := range max {
			max[i] = 0xff
		}
		if golangt := ConstantTimeLessOrEqBytes(a, max); golangt != 1 {
			t.Errorf("ConstantTimeLessOrEqBytes(%x, %x) = %d, want 1", a, max, golangt)
		}
		if golangt := ConstantTimeLessOrEqBytes(b, max); golangt != 1 {
			t.Errorf("ConstantTimeLessOrEqBytes(%x, %x) = %d, want 1", b, max, golangt)
		}
		if golangt := ConstantTimeLessOrEqBytes(empty, max); golangt != 1 {
			t.Errorf("ConstantTimeLessOrEqBytes(%x, %x) = %d, want 1", empty, max, golangt)
		}
		if golangt := ConstantTimeLessOrEqBytes(max, max); golangt != 1 {
			t.Errorf("ConstantTimeLessOrEqBytes(%x, %x) = %d, want 1", max, max, golangt)
		}
		aPlusOne := make([]byte, l)
		copy(aPlusOne, a)
		for i := l - 1; i >= 0; i-- {
			if aPlusOne[i] == 0xff {
				aPlusOne[i] = 0
				continue
			}
			aPlusOne[i]++
			if golangt := ConstantTimeLessOrEqBytes(a, aPlusOne); golangt != 1 {
				t.Errorf("ConstantTimeLessOrEqBytes(%x, %x) = %d, want 1", a, aPlusOne, golangt)
			}
			if golangt := ConstantTimeLessOrEqBytes(aPlusOne, a); golangt != 0 {
				t.Errorf("ConstantTimeLessOrEqBytes(%x, %x) = %d, want 0", aPlusOne, a, golangt)
			}
			break
		}
		shorter := make([]byte, l-1)
		copy(shorter, a)
		if golangt := ConstantTimeLessOrEqBytes(a, shorter); golangt != 0 {
			t.Errorf("ConstantTimeLessOrEqBytes(%x, %x) = %d, want 0", a, shorter, golangt)
		}
		if golangt := ConstantTimeLessOrEqBytes(shorter, a); golangt != 0 {
			t.Errorf("ConstantTimeLessOrEqBytes(%x, %x) = %d, want 0", shorter, a, golangt)
		}
		if golangt := ConstantTimeLessOrEqBytes(b, shorter); golangt != 0 {
			t.Errorf("ConstantTimeLessOrEqBytes(%x, %x) = %d, want 0", b, shorter, golangt)
		}
		if golangt := ConstantTimeLessOrEqBytes(shorter, b); golangt != 0 {
			t.Errorf("ConstantTimeLessOrEqBytes(%x, %x) = %d, want 0", shorter, b, golangt)
		}
	}
}
