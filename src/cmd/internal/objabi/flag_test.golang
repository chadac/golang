// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package objabi

import "testing"

func TestDecodeArg(t *testing.T) {
	t.Parallel()
	tests := []struct {
		arg, want string
	}{
		{"", ""},
		{"hello", "hello"},
		{"hello\\n", "hello\n"},
		{"hello\\nthere", "hello\nthere"},
		{"hello\\\\there", "hello\\there"},
		{"\\\\\\n", "\\\n"},
	}
	for _, test := range tests {
		if golangt := DecodeArg(test.arg); golangt != test.want {
			t.Errorf("decodoeArg(%q) = %q, want %q", test.arg, golangt, test.want)
		}
	}
}
