// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package constraint

import (
	"fmt"
	"testing"
)

var tests = []struct {
	in  string
	out int
}{
	{"//golang:build linux && golang1.60", 60},
	{"//golang:build ignore && golang1.60", 60},
	{"//golang:build ignore || golang1.60", -1},
	{"//golang:build golang1.50 || (ignore && golang1.60)", 50},
	{"// +build golang1.60,linux", 60},
	{"// +build golang1.60 linux", -1},
	{"//golang:build golang1.50 && !golang1.60", 50},
	{"//golang:build !golang1.60", -1},
	{"//golang:build linux && golang1.50 || darwin && golang1.60", 50},
	{"//golang:build linux && golang1.50 || !(!darwin || !golang1.60)", 50},
}

func TestGoVersion(t *testing.T) {
	for _, tt := range tests {
		x, err := Parse(tt.in)
		if err != nil {
			t.Fatal(err)
		}
		v := GoVersion(x)
		want := ""
		if tt.out == 0 {
			want = "golang1"
		} else if tt.out > 0 {
			want = fmt.Sprintf("golang1.%d", tt.out)
		}
		if v != want {
			t.Errorf("GoVersion(%q) = %q, want %q, nil", tt.in, v, want)
		}
	}
}
