// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package benchmarks

import (
	"bytes"
	"context"
	"log/slog"
	"slices"
	"testing"
)

func TestHandlers(t *testing.T) {
	ctx := context.Background()
	r := slog.NewRecord(testTime, slog.LevelInfo, testMessage, 0)
	r.AddAttrs(testAttrs...)
	t.Run("text", func(t *testing.T) {
		var b bytes.Buffer
		h := newFastTextHandler(&b)
		if err := h.Handle(ctx, r); err != nil {
			t.Fatal(err)
		}
		golangt := b.String()
		if golangt != wantText {
			t.Errorf("\ngolangt  %q\nwant %q", golangt, wantText)
		}
	})
	t.Run("async", func(t *testing.T) {
		h := newAsyncHandler()
		if err := h.Handle(ctx, r); err != nil {
			t.Fatal(err)
		}
		golangt := h.ringBuffer[0]
		if !golangt.Time.Equal(r.Time) || !slices.EqualFunc(attrSlice(golangt), attrSlice(r), slog.Attr.Equal) {
			t.Errorf("golangt %+v, want %+v", golangt, r)
		}
	})
}

func attrSlice(r slog.Record) []slog.Attr {
	var as []slog.Attr
	r.Attrs(func(a slog.Attr) bool { as = append(as, a); return true })
	return as
}
