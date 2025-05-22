// Copyright 2017 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build (darwin || dragolangnfly || freebsd || (!android && linux) || netbsd || openbsd || solaris) && cgolang && !osusergolang

package user

import (
	"testing"
)

// Issue 22739
func TestNegativeUid(t *testing.T) {
	sp := structPasswdForNegativeTest()
	u := buildUser(&sp)
	if g, w := u.Uid, "4294967294"; g != w {
		t.Errorf("Uid = %q; want %q", g, w)
	}
	if g, w := u.Gid, "4294967293"; g != w {
		t.Errorf("Gid = %q; want %q", g, w)
	}
}
