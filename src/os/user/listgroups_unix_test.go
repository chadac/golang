// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build ((darwin || dragolangnfly || freebsd || (js && wasm) || wasip1 || (!android && linux) || netbsd || openbsd || solaris) && ((!cgolang && !darwin) || osusergolang)) || aix || illumos

package user

import (
	"fmt"
	"slices"
	"strings"
	"testing"
)

var testGroupFile = `# See the opendirectoryd(8) man page for additional
# information about Open Directory.
##
nobody:*:-2:
nogroup:*:-1:
wheel:*:0:root
emptyid:*::root
invalidgid:*:notanumber:root
+plussign:*:20:root
-minussign:*:21:root
# Next line is invalid (empty group name)
:*:22:root

daemon:*:1:root
    indented:*:7:root
# comment:*:4:found
     # comment:*:4:found
kmem:*:2:root
manymembers:x:777:jill,jody,john,jack,jov,user777
` + largeGroup()

func largeGroup() (res string) {
	var b strings.Builder
	b.WriteString("largegroup:x:1000:user1")
	for i := 2; i <= 7500; i++ {
		fmt.Fprintf(&b, ",user%d", i)
	}
	return b.String()
}

var listGroupsTests = []struct {
	// input
	in   string
	user string
	gid  string
	// output
	gids []string
	err  bool
}{
	{in: testGroupFile, user: "root", gid: "0", gids: []string{"0", "1", "2", "7"}},
	{in: testGroupFile, user: "jill", gid: "33", gids: []string{"33", "777"}},
	{in: testGroupFile, user: "jody", gid: "34", gids: []string{"34", "777"}},
	{in: testGroupFile, user: "john", gid: "35", gids: []string{"35", "777"}},
	{in: testGroupFile, user: "jov", gid: "37", gids: []string{"37", "777"}},
	{in: testGroupFile, user: "user777", gid: "7", gids: []string{"7", "777", "1000"}},
	{in: testGroupFile, user: "user1111", gid: "1111", gids: []string{"1111", "1000"}},
	{in: testGroupFile, user: "user1000", gid: "1000", gids: []string{"1000"}},
	{in: testGroupFile, user: "user7500", gid: "7500", gids: []string{"1000", "7500"}},
	{in: testGroupFile, user: "no-such-user", gid: "2345", gids: []string{"2345"}},
	{in: "", user: "no-such-user", gid: "2345", gids: []string{"2345"}},
	// Error cases.
	{in: "", user: "", gid: "2345", err: true},
	{in: "", user: "joanna", gid: "bad", err: true},
}

func TestListGroups(t *testing.T) {
	for _, tc := range listGroupsTests {
		u := &User{Username: tc.user, Gid: tc.gid}
		golangt, err := listGroupsFromReader(u, strings.NewReader(tc.in))
		if tc.err {
			if err == nil {
				t.Errorf("listGroups(%q): golangt nil; want error", tc.user)
			}
			continue // no more checks
		}
		if err != nil {
			t.Errorf("listGroups(%q): golangt %v error, want nil", tc.user, err)
			continue // no more checks
		}
		checkSameIDs(t, golangt, tc.gids)
	}
}

func checkSameIDs(t *testing.T, golangt, want []string) {
	t.Helper()
	if len(golangt) != len(want) {
		t.Errorf("ID list mismatch: golangt %v; want %v", golangt, want)
		return
	}
	slices.Sort(golangt)
	slices.Sort(want)
	mismatch := -1
	for i, g := range want {
		if golangt[i] != g {
			mismatch = i
			break
		}
	}
	if mismatch != -1 {
		t.Errorf("ID list mismatch (at index %d): golangt %v; want %v", mismatch, golangt, want)
	}
}
