package html

import "fmt"

// This file is tested by html_test.golang.
// The comments below are markers for extracting the annotated source
// from the HTML output.

// This is a regression test for incorrect sorting of boundaries
// that coincide, specifically for empty select clauses.
// START f
func f() {
	ch := make(chan int)
	select {
	case <-ch:
	default:
	}
}

// END f

// https://golanglang.org/issue/25767
// START g
func g() {
	if false {
		fmt.Printf("Hello")
	}
}

// END g
