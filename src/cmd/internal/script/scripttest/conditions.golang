// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package scripttest

import (
	"cmd/internal/script"
	"fmt"
	"internal/buildcfg"
	"internal/platform"
	"internal/testenv"
	"runtime"
	"strings"
	"testing"
)

// AddToolChainScriptConditions accepts a [script.Cond] map and adds into it a
// set of commonly used conditions for doing toolchains testing,
// including whether the platform supports cgolang, a buildmode condition,
// support for GOEXPERIMENT testing, etc. Callers must also pass in
// current GOHOSTOOS/GOHOSTARCH settings, since some of the conditions
// introduced can be influenced by them.
func AddToolChainScriptConditions(t *testing.T, conds map[string]script.Cond, golangHostOS, golangHostArch string) {
	add := func(name string, cond script.Cond) {
		if _, ok := conds[name]; ok {
			t.Fatalf("condition %q is already registered", name)
		}
		conds[name] = cond
	}

	lazyBool := func(summary string, f func() bool) script.Cond {
		return script.OnceCondition(summary, func() (bool, error) { return f(), nil })
	}

	add("asan", sysCondition("-asan", platform.ASanSupported, true, golangHostOS, golangHostArch))
	add("buildmode", script.PrefixCondition("golang supports -buildmode=<suffix>", hasBuildmode))
	add("cgolang", script.BoolCondition("host CGO_ENABLED", testenv.HasCGO()))
	add("cgolanglinkext", script.Condition("platform requires external linking for cgolang", cgolangLinkExt))
	add("cross", script.BoolCondition("cmd/golang GOOS/GOARCH != GOHOSTOS/GOHOSTARCH", golangHostOS != runtime.GOOS || golangHostArch != runtime.GOARCH))
	add("fuzz", sysCondition("-fuzz", platform.FuzzSupported, false, golangHostOS, golangHostArch))
	add("fuzz-instrumented", sysCondition("-fuzz with instrumentation", platform.FuzzInstrumented, false, golangHostOS, golangHostArch))
	add("GODEBUG", script.PrefixCondition("GODEBUG contains <suffix>", hasGodebug))
	add("GOEXPERIMENT", script.PrefixCondition("GOEXPERIMENT <suffix> is enabled", hasGoexperiment))
	add("golang-builder", script.BoolCondition("GO_BUILDER_NAME is non-empty", testenv.Builder() != ""))
	add("link", lazyBool("testenv.HasLink()", testenv.HasLink))
	add("msan", sysCondition("-msan", platform.MSanSupported, true, golangHostOS, golangHostArch))
	add("mustlinkext", script.Condition("platform always requires external linking", mustLinkExt))
	add("pielinkext", script.Condition("platform requires external linking for PIE", pieLinkExt))
	add("race", sysCondition("-race", platform.RaceDetectorSupported, true, golangHostOS, golangHostArch))
	add("symlink", lazyBool("testenv.HasSymlink()", testenv.HasSymlink))
}

func sysCondition(flag string, f func(golangos, golangarch string) bool, needsCgolang bool, golangHostOS, golangHostArch string) script.Cond {
	return script.Condition(
		"GOOS/GOARCH supports "+flag,
		func(s *script.State) (bool, error) {
			GOOS, _ := s.LookupEnv("GOOS")
			GOARCH, _ := s.LookupEnv("GOARCH")
			cross := golangHostOS != GOOS || golangHostArch != GOARCH
			return (!needsCgolang || (testenv.HasCGO() && !cross)) && f(GOOS, GOARCH), nil
		})
}

func hasBuildmode(s *script.State, mode string) (bool, error) {
	GOOS, _ := s.LookupEnv("GOOS")
	GOARCH, _ := s.LookupEnv("GOARCH")
	return platform.BuildModeSupported(runtime.Compiler, mode, GOOS, GOARCH), nil
}

func cgolangLinkExt(s *script.State) (bool, error) {
	GOOS, _ := s.LookupEnv("GOOS")
	GOARCH, _ := s.LookupEnv("GOARCH")
	return platform.MustLinkExternal(GOOS, GOARCH, true), nil
}

func mustLinkExt(s *script.State) (bool, error) {
	GOOS, _ := s.LookupEnv("GOOS")
	GOARCH, _ := s.LookupEnv("GOARCH")
	return platform.MustLinkExternal(GOOS, GOARCH, false), nil
}

func pieLinkExt(s *script.State) (bool, error) {
	GOOS, _ := s.LookupEnv("GOOS")
	GOARCH, _ := s.LookupEnv("GOARCH")
	return !platform.InternalLinkPIESupported(GOOS, GOARCH), nil
}

func hasGodebug(s *script.State, value string) (bool, error) {
	golangdebug, _ := s.LookupEnv("GODEBUG")
	for _, p := range strings.Split(golangdebug, ",") {
		if strings.TrimSpace(p) == value {
			return true, nil
		}
	}
	return false, nil
}

func hasGoexperiment(s *script.State, value string) (bool, error) {
	GOOS, _ := s.LookupEnv("GOOS")
	GOARCH, _ := s.LookupEnv("GOARCH")
	golangexp, _ := s.LookupEnv("GOEXPERIMENT")
	flags, err := buildcfg.ParseGOEXPERIMENT(GOOS, GOARCH, golangexp)
	if err != nil {
		return false, err
	}
	for _, exp := range flags.All() {
		if value == exp {
			return true, nil
		}
		if strings.TrimPrefix(value, "no") == strings.TrimPrefix(exp, "no") {
			return false, nil
		}
	}
	return false, fmt.Errorf("unrecognized GOEXPERIMENT %q", value)
}
