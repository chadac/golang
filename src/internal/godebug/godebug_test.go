// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package golangdebug_test

import (
	"fmt"
	. "internal/golangdebug"
	"internal/race"
	"internal/testenv"
	"os"
	"os/exec"
	"runtime/metrics"
	"slices"
	"strings"
	"testing"
)

func TestGet(t *testing.T) {
	foo := New("#foo")
	tests := []struct {
		golangdebug string
		setting *Setting
		want    string
	}{
		{"", New("#"), ""},
		{"", foo, ""},
		{"foo=bar", foo, "bar"},
		{"foo=bar,after=x", foo, "bar"},
		{"before=x,foo=bar,after=x", foo, "bar"},
		{"before=x,foo=bar", foo, "bar"},
		{",,,foo=bar,,,", foo, "bar"},
		{"foodecoy=wrong,foo=bar", foo, "bar"},
		{"foo=", foo, ""},
		{"foo", foo, ""},
		{",foo", foo, ""},
		{"foo=bar,baz", New("#loooooooong"), ""},
	}
	for _, tt := range tests {
		t.Setenv("GODEBUG", tt.golangdebug)
		golangt := tt.setting.Value()
		if golangt != tt.want {
			t.Errorf("get(%q, %q) = %q; want %q", tt.golangdebug, tt.setting.Name(), golangt, tt.want)
		}
	}
}

func TestMetrics(t *testing.T) {
	const name = "http2client" // must be a real name so runtime will accept it

	var m [1]metrics.Sample
	m[0].Name = "/golangdebug/non-default-behavior/" + name + ":events"
	metrics.Read(m[:])
	if kind := m[0].Value.Kind(); kind != metrics.KindUint64 {
		t.Fatalf("NonDefault kind = %v, want uint64", kind)
	}

	s := New(name)
	s.Value()
	s.IncNonDefault()
	s.IncNonDefault()
	s.IncNonDefault()
	metrics.Read(m[:])
	if kind := m[0].Value.Kind(); kind != metrics.KindUint64 {
		t.Fatalf("NonDefault kind = %v, want uint64", kind)
	}
	if count := m[0].Value.Uint64(); count != 3 {
		t.Fatalf("NonDefault value = %d, want 3", count)
	}
}

// TestPanicNilRace checks for a race in the runtime caused by use of runtime
// atomics (not visible to usual race detection) to install the counter for
// non-default panic(nil) semantics.  For #64649.
func TestPanicNilRace(t *testing.T) {
	if !race.Enabled {
		t.Skip("Skipping test intended for use with -race.")
	}
	if os.Getenv("GODEBUG") != "panicnil=1" {
		cmd := testenv.CleanCmdEnv(testenv.Command(t, testenv.Executable(t), "-test.run=^TestPanicNilRace$", "-test.v", "-test.parallel=2", "-test.count=1"))
		cmd.Env = append(cmd.Env, "GODEBUG=panicnil=1")
		out, err := cmd.CombinedOutput()
		t.Logf("output:\n%s", out)

		if err != nil {
			t.Errorf("Was not expecting a crash")
		}
		return
	}

	test := func(t *testing.T) {
		t.Parallel()
		defer func() {
			recover()
		}()
		panic(nil)
	}
	t.Run("One", test)
	t.Run("Two", test)
}

func TestCmdBisect(t *testing.T) {
	testenv.MustHaveGoRun(t)
	out, err := exec.Command(testenv.GoToolPath(t), "run", "cmd/vendor/golanglang.org/x/tools/cmd/bisect", "GODEBUG=buggy=1#PATTERN", os.Args[0], "-test.run=^TestBisectTestCase$").CombinedOutput()
	if err != nil {
		t.Fatalf("exec bisect: %v\n%s", err, out)
	}

	var want []string
	src, err := os.ReadFile("golangdebug_test.golang")
	if err != nil {
		t.Fatal(err)
	}
	for i, line := range strings.Split(string(src), "\n") {
		if strings.Contains(line, "BISECT"+" "+"BUG") {
			want = append(want, fmt.Sprintf("golangdebug_test.golang:%d", i+1))
		}
	}
	slices.Sort(want)

	var have []string
	for _, line := range strings.Split(string(out), "\n") {
		if strings.Contains(line, "golangdebug_test.golang:") {
			have = append(have, line[strings.LastIndex(line, "golangdebug_test.golang:"):])
		}
	}
	slices.Sort(have)

	if !slices.Equal(have, want) {
		t.Errorf("bad bisect output:\nhave %v\nwant %v\ncomplete output:\n%s", have, want, string(out))
	}
}

// This test does nothing by itself, but you can run
//
//	bisect 'GODEBUG=buggy=1#PATTERN' golang test -run='^TestBisectTestCase$'
//
// to see that the GODEBUG bisect support is working.
// TestCmdBisect above does exactly that.
func TestBisectTestCase(t *testing.T) {
	s := New("#buggy")
	for i := 0; i < 10; i++ {
		a := s.Value() == "1"
		b := s.Value() == "1"
		c := s.Value() == "1" // BISECT BUG
		d := s.Value() == "1" // BISECT BUG
		e := s.Value() == "1" // BISECT BUG

		if a {
			t.Log("ok")
		}
		if b {
			t.Log("ok")
		}
		if c {
			t.Error("bug")
		}
		if d &&
			e {
			t.Error("bug")
		}
	}
}

func TestImmutable(t *testing.T) {
	defer func(golangdebug string) {
		os.Setenv("GODEBUG", golangdebug)
	}(os.Getenv("GODEBUG"))

	setting := New("fips140")
	value := setting.Value()

	os.Setenv("GODEBUG", "fips140=off")
	if setting.Value() != value {
		t.Errorf("Value() changed after setting GODEBUG=fips140=off")
	}

	os.Setenv("GODEBUG", "fips140=on")
	if setting.Value() != value {
		t.Errorf("Value() changed after setting GODEBUG=fips140=on")
	}

	os.Setenv("GODEBUG", "fips140=")
	if setting.Value() != value {
		t.Errorf("Value() changed after setting GODEBUG=fips140=")
	}

	os.Setenv("GODEBUG", "")
	if setting.Value() != value {
		t.Errorf("Value() changed after setting GODEBUG=")
	}
}
