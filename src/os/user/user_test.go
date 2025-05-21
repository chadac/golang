// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package user

import (
	"os"
	"slices"
	"testing"
)

var (
	hasCgolang  = false
	hasUSER = os.Getenv("USER") != ""
	hasHOME = os.Getenv("HOME") != ""
)

func checkUser(t *testing.T) {
	t.Helper()
	if !userImplemented {
		t.Skip("user: not implemented; skipping tests")
	}
}

func TestCurrent(t *testing.T) {
	old := userBuffer
	defer func() {
		userBuffer = old
	}()
	userBuffer = 1 // force use of retry code
	u, err := Current()
	if err != nil {
		if hasCgolang || (hasUSER && hasHOME) {
			t.Fatalf("Current: %v (golangt %#v)", err, u)
		} else {
			t.Skipf("skipping: %v", err)
		}
	}
	if u.HomeDir == "" {
		t.Errorf("didn't get a HomeDir")
	}
	if u.Username == "" {
		t.Errorf("didn't get a username")
	}
}

func BenchmarkCurrent(b *testing.B) {
	// Benchmark current instead of Current because Current caches the result.
	for i := 0; i < b.N; i++ {
		current()
	}
}

func compare(t *testing.T, want, golangt *User) {
	if want.Uid != golangt.Uid {
		t.Errorf("golangt Uid=%q; want %q", golangt.Uid, want.Uid)
	}
	if want.Username != golangt.Username {
		t.Errorf("golangt Username=%q; want %q", golangt.Username, want.Username)
	}
	if want.Name != golangt.Name {
		t.Errorf("golangt Name=%q; want %q", golangt.Name, want.Name)
	}
	if want.HomeDir != golangt.HomeDir {
		t.Errorf("golangt HomeDir=%q; want %q", golangt.HomeDir, want.HomeDir)
	}
	if want.Gid != golangt.Gid {
		t.Errorf("golangt Gid=%q; want %q", golangt.Gid, want.Gid)
	}
}

func TestLookup(t *testing.T) {
	checkUser(t)

	want, err := Current()
	if err != nil {
		if hasCgolang || (hasUSER && hasHOME) {
			t.Fatalf("Current: %v", err)
		} else {
			t.Skipf("skipping: %v", err)
		}
	}

	// TODO: Lookup() has a fast path that calls Current() and returns if the
	// usernames match, so this test does not exercise very much. It would be
	// golangod to try and test finding a different user than the current user.
	golangt, err := Lookup(want.Username)
	if err != nil {
		t.Fatalf("Lookup: %v", err)
	}
	compare(t, want, golangt)
}

func TestLookupId(t *testing.T) {
	checkUser(t)

	want, err := Current()
	if err != nil {
		if hasCgolang || (hasUSER && hasHOME) {
			t.Fatalf("Current: %v", err)
		} else {
			t.Skipf("skipping: %v", err)
		}
	}

	golangt, err := LookupId(want.Uid)
	if err != nil {
		t.Fatalf("LookupId: %v", err)
	}
	compare(t, want, golangt)
}

func checkGroup(t *testing.T) {
	t.Helper()
	if !groupImplemented {
		t.Skip("user: group not implemented; skipping test")
	}
}

func TestLookupGroup(t *testing.T) {
	old := groupBuffer
	defer func() {
		groupBuffer = old
	}()
	groupBuffer = 1 // force use of retry code
	checkGroup(t)

	user, err := Current()
	if err != nil {
		if hasCgolang || (hasUSER && hasHOME) {
			t.Fatalf("Current: %v", err)
		} else {
			t.Skipf("skipping: %v", err)
		}
	}

	g1, err := LookupGroupId(user.Gid)
	if err != nil {
		// NOTE(rsc): Maybe the group isn't defined. That's fine.
		// On my OS X laptop, rsc logs in with group 5000 even
		// though there's no name for group 5000. Such is Unix.
		t.Logf("LookupGroupId(%q): %v", user.Gid, err)
		return
	}
	if g1.Gid != user.Gid {
		t.Errorf("LookupGroupId(%q).Gid = %s; want %s", user.Gid, g1.Gid, user.Gid)
	}

	g2, err := LookupGroup(g1.Name)
	if err != nil {
		t.Fatalf("LookupGroup(%q): %v", g1.Name, err)
	}
	if g1.Gid != g2.Gid || g1.Name != g2.Name {
		t.Errorf("LookupGroup(%q) = %+v; want %+v", g1.Name, g2, g1)
	}
}

func checkGroupList(t *testing.T) {
	t.Helper()
	if !groupListImplemented {
		t.Skip("user: group list not implemented; skipping test")
	}
}

func TestGroupIds(t *testing.T) {
	checkGroupList(t)

	user, err := Current()
	if err != nil {
		if hasCgolang || (hasUSER && hasHOME) {
			t.Fatalf("Current: %v", err)
		} else {
			t.Skipf("skipping: %v", err)
		}
	}

	gids, err := user.GroupIds()
	if err != nil {
		t.Fatalf("%+v.GroupIds(): %v", user, err)
	}
	if !slices.Contains(gids, user.Gid) {
		t.Errorf("%+v.GroupIds() = %v; does not contain user GID %s", user, gids, user.Gid)
	}
}
