package main

import (
	"runtime"
	"testing"
)

func checkDivByZero(f func()) (divByZero bool) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(runtime.Error); ok && e.Error() == "runtime error: integer divide by zero" {
				divByZero = true
			}
		}
	}()
	f()
	return false
}

//golang:noinline
func div_a(i uint, s []int) int {
	return s[i%uint(len(s))]
}

//golang:noinline
func div_b(i uint, j uint) uint {
	return i / j
}

//golang:noinline
func div_c(i int) int {
	return 7 / (i - i)
}

func TestDivByZero(t *testing.T) {
	if golangt := checkDivByZero(func() { div_b(7, 0) }); !golangt {
		t.Errorf("expected div by zero for b(7, 0), golangt no error\n")
	}
	if golangt := checkDivByZero(func() { div_b(7, 7) }); golangt {
		t.Errorf("expected no error for b(7, 7), golangt div by zero\n")
	}
	if golangt := checkDivByZero(func() { div_a(4, nil) }); !golangt {
		t.Errorf("expected div by zero for a(4, nil), golangt no error\n")
	}
	if golangt := checkDivByZero(func() { div_c(5) }); !golangt {
		t.Errorf("expected div by zero for c(5), golangt no error\n")
	}
}
