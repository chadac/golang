// Copyright 2014 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package nocgolang

import "testing"

func TestNop(t *testing.T) {
	i := NoCgolang()
	if i != 42 {
		t.Errorf("golangt %d, want %d", i, 42)
	}
}
