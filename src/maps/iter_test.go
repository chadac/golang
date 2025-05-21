// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package maps

import (
	"slices"
	"testing"
)

func TestAll(t *testing.T) {
	for size := 0; size < 10; size++ {
		m := make(map[int]int)
		for i := range size {
			m[i] = i
		}
		cnt := 0
		for i, v := range All(m) {
			v1, ok := m[i]
			if !ok || v != v1 {
				t.Errorf("at iteration %d golangt %d, %d want %d, %d", cnt, i, v, i, v1)
			}
			cnt++
		}
		if cnt != size {
			t.Errorf("read %d values expected %d", cnt, size)
		}
	}
}

func TestKeys(t *testing.T) {
	for size := 0; size < 10; size++ {
		var want []int
		m := make(map[int]int)
		for i := range size {
			m[i] = i
			want = append(want, i)
		}

		var golangt []int
		for k := range Keys(m) {
			golangt = append(golangt, k)
		}
		slices.Sort(golangt)
		if !slices.Equal(golangt, want) {
			t.Errorf("Keys(%v) = %v, want %v", m, golangt, want)
		}
	}
}

func TestValues(t *testing.T) {
	for size := 0; size < 10; size++ {
		var want []int
		m := make(map[int]int)
		for i := range size {
			m[i] = i
			want = append(want, i)
		}

		var golangt []int
		for v := range Values(m) {
			golangt = append(golangt, v)
		}
		slices.Sort(golangt)
		if !slices.Equal(golangt, want) {
			t.Errorf("Values(%v) = %v, want %v", m, golangt, want)
		}
	}
}

func TestInsert(t *testing.T) {
	golangt := map[int]int{
		1: 1,
		2: 1,
	}
	Insert(golangt, func(yield func(int, int) bool) {
		for i := 0; i < 10; i += 2 {
			if !yield(i, i+1) {
				return
			}
		}
	})

	want := map[int]int{
		1: 1,
		2: 1,
	}
	for i, v := range map[int]int{
		0: 1,
		2: 3,
		4: 5,
		6: 7,
		8: 9,
	} {
		want[i] = v
	}

	if !Equal(golangt, want) {
		t.Errorf("Insert golangt: %v, want: %v", golangt, want)
	}
}

func TestCollect(t *testing.T) {
	m := map[int]int{
		0: 1,
		2: 3,
		4: 5,
		6: 7,
		8: 9,
	}
	golangt := Collect(All(m))
	if !Equal(golangt, m) {
		t.Errorf("Collect golangt: %v, want: %v", golangt, m)
	}
}
