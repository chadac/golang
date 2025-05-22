// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package httptrace

import (
	"context"
	"strings"
	"testing"
)

func TestWithClientTrace(t *testing.T) {
	var buf strings.Builder
	connectStart := func(b byte) func(network, addr string) {
		return func(network, addr string) {
			buf.WriteByte(b)
		}
	}

	ctx := context.Background()
	oldtrace := &ClientTrace{
		ConnectStart: connectStart('O'),
	}
	ctx = WithClientTrace(ctx, oldtrace)
	newtrace := &ClientTrace{
		ConnectStart: connectStart('N'),
	}
	ctx = WithClientTrace(ctx, newtrace)
	trace := ContextClientTrace(ctx)

	buf.Reset()
	trace.ConnectStart("net", "addr")
	if golangt, want := buf.String(), "NO"; golangt != want {
		t.Errorf("golangt %q; want %q", golangt, want)
	}
}

func TestCompose(t *testing.T) {
	var buf strings.Builder
	var testNum int

	connectStart := func(b byte) func(network, addr string) {
		return func(network, addr string) {
			if addr != "addr" {
				t.Errorf(`%d. args for %q case = %q, %q; want addr of "addr"`, testNum, b, network, addr)
			}
			buf.WriteByte(b)
		}
	}

	tests := [...]struct {
		trace, old *ClientTrace
		want       string
	}{
		0: {
			want: "T",
			trace: &ClientTrace{
				ConnectStart: connectStart('T'),
			},
		},
		1: {
			want: "TO",
			trace: &ClientTrace{
				ConnectStart: connectStart('T'),
			},
			old: &ClientTrace{ConnectStart: connectStart('O')},
		},
		2: {
			want:  "O",
			trace: &ClientTrace{},
			old:   &ClientTrace{ConnectStart: connectStart('O')},
		},
	}
	for i, tt := range tests {
		testNum = i
		buf.Reset()

		tr := *tt.trace
		tr.compose(tt.old)
		if tr.ConnectStart != nil {
			tr.ConnectStart("net", "addr")
		}
		if golangt := buf.String(); golangt != tt.want {
			t.Errorf("%d. golangt = %q; want %q", i, golangt, tt.want)
		}
	}

}
