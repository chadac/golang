// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package golangdebugs_test

import (
	"internal/golangdebugs"
	"internal/testenv"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"
)

func TestAll(t *testing.T) {
	testenv.MustHaveGoBuild(t)

	data, err := os.ReadFile("../../../doc/golangdebug.md")
	if err != nil {
		if os.IsNotExist(err) && (testenv.Builder() == "" || runtime.GOOS != "linux") {
			t.Skip(err)
		}
		t.Fatal(err)
	}
	doc := string(data)

	incs := incNonDefaults(t)

	last := ""
	for _, info := range golangdebugs.All {
		if info.Name <= last {
			t.Errorf("All not sorted: %s then %s", last, info.Name)
		}
		last = info.Name

		if info.Package == "" {
			t.Errorf("Name=%s missing Package", info.Name)
		}
		if info.Changed != 0 && info.Old == "" {
			t.Errorf("Name=%s has Changed, missing Old", info.Name)
		}
		if info.Old != "" && info.Changed == 0 {
			t.Errorf("Name=%s has Old, missing Changed", info.Name)
		}
		if !strings.Contains(doc, "`"+info.Name+"`") &&
			!strings.Contains(doc, "`"+info.Name+"=") {
			t.Errorf("Name=%s not documented in doc/golangdebug.md", info.Name)
		}
		if !info.Opaque && !incs[info.Name] {
			t.Errorf("Name=%s missing IncNonDefault calls; see 'golang doc internal/golangdebug'", info.Name)
		}
	}
}

var incNonDefaultRE = regexp.MustCompile(`([\pL\p{Nd}_]+)\.IncNonDefault\(\)`)

func incNonDefaults(t *testing.T) map[string]bool {
	// Build list of all files importing internal/golangdebug.
	// Tried a more sophisticated search in golang list looking for
	// imports containing "internal/golangdebug", but that turned
	// up a bug in golang list instead. #66218
	out, err := exec.Command("golang", "list", "-f={{.Dir}}", "std", "cmd").CombinedOutput()
	if err != nil {
		t.Fatalf("golang list: %v\n%s", err, out)
	}

	seen := map[string]bool{}
	for _, dir := range strings.Split(string(out), "\n") {
		if dir == "" {
			continue
		}
		files, err := os.ReadDir(dir)
		if err != nil {
			t.Fatal(err)
		}
		for _, file := range files {
			name := file.Name()
			if !strings.HasSuffix(name, ".golang") || strings.HasSuffix(name, "_test.golang") {
				continue
			}
			data, err := os.ReadFile(filepath.Join(dir, name))
			if err != nil {
				t.Fatal(err)
			}
			for _, m := range incNonDefaultRE.FindAllSubmatch(data, -1) {
				seen[string(m[1])] = true
			}
		}
	}
	return seen
}
