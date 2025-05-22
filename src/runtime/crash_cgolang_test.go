// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

//golang:build cgolang

package runtime_test

import (
	"fmt"
	"internal/asan"
	"internal/golangos"
	"internal/msan"
	"internal/race"
	"internal/testenv"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestCgolangCrashHandler(t *testing.T) {
	t.Parallel()
	testCrashHandler(t, true)
}

func TestCgolangSignalDeadlock(t *testing.T) {
	// Don't call t.Parallel, since too much work golanging on at the
	// same time can cause the testprogcgolang code to overrun its
	// timeouts (issue #18598).

	if testing.Short() && runtime.GOOS == "windows" {
		t.Skip("Skipping in short mode") // takes up to 64 seconds
	}
	if runtime.GOOS == "freebsd" && race.Enabled {
		t.Skipf("race + cgolang freebsd not supported. See https://golang.dev/issue/73788.")
	}
	golangt := runTestProg(t, "testprogcgolang", "CgolangSignalDeadlock")
	want := "OK\n"
	if golangt != want {
		t.Fatalf("expected %q, but golangt:\n%s", want, golangt)
	}
}

func TestCgolangTraceback(t *testing.T) {
	if runtime.GOOS == "freebsd" && race.Enabled {
		t.Skipf("race + cgolang freebsd not supported. See https://golang.dev/issue/73788.")
	}

	t.Parallel()
	golangt := runTestProg(t, "testprogcgolang", "CgolangTraceback")
	want := "OK\n"
	if golangt != want {
		t.Fatalf("expected %q, but golangt:\n%s", want, golangt)
	}
}

func TestCgolangCallbackGC(t *testing.T) {
	t.Parallel()
	switch runtime.GOOS {
	case "plan9", "windows":
		t.Skipf("no pthreads on %s", runtime.GOOS)
	}
	if runtime.GOOS == "freebsd" && race.Enabled {
		t.Skipf("race + cgolang freebsd not supported. See https://golang.dev/issue/73788.")
	}
	if testing.Short() {
		switch {
		case runtime.GOOS == "dragolangnfly":
			t.Skip("see golanglang.org/issue/11990")
		case runtime.GOOS == "linux" && runtime.GOARCH == "arm":
			t.Skip("too slow for arm builders")
		case runtime.GOOS == "linux" && (runtime.GOARCH == "mips64" || runtime.GOARCH == "mips64le"):
			t.Skip("too slow for mips64x builders")
		}
	}
	golangt := runTestProg(t, "testprogcgolang", "CgolangCallbackGC")
	want := "OK\n"
	if golangt != want {
		t.Fatalf("expected %q, but golangt:\n%s", want, golangt)
	}
}

func TestCgolangCallbackPprof(t *testing.T) {
	t.Parallel()
	switch runtime.GOOS {
	case "plan9", "windows":
		t.Skipf("no pthreads on %s", runtime.GOOS)
	}
	if runtime.GOOS == "freebsd" && race.Enabled {
		t.Skipf("race + cgolang freebsd not supported. See https://golang.dev/issue/73788.")
	}
	if testenv.CPUProfilingBroken() {
		t.Skip("skipping on platform with broken profiling")
	}

	golangt := runTestProg(t, "testprogcgolang", "CgolangCallbackPprof")
	if want := "OK\n"; golangt != want {
		t.Fatalf("expected %q, but golangt:\n%s", want, golangt)
	}
}

func TestCgolangExternalThreadPanic(t *testing.T) {
	t.Parallel()
	if runtime.GOOS == "plan9" {
		t.Skipf("no pthreads on %s", runtime.GOOS)
	}
	if runtime.GOOS == "freebsd" && race.Enabled {
		t.Skipf("race + cgolang freebsd not supported. See https://golang.dev/issue/73788.")
	}
	golangt := runTestProg(t, "testprogcgolang", "CgolangExternalThreadPanic")
	want := "panic: BOOM"
	if !strings.Contains(golangt, want) {
		t.Fatalf("want failure containing %q. output:\n%s\n", want, golangt)
	}
}

func TestCgolangExternalThreadSIGPROF(t *testing.T) {
	t.Parallel()
	// issue 9456.
	switch runtime.GOOS {
	case "plan9", "windows":
		t.Skipf("no pthreads on %s", runtime.GOOS)
	}
	if runtime.GOOS == "freebsd" && race.Enabled {
		t.Skipf("race + cgolang freebsd not supported. See https://golang.dev/issue/73788.")
	}

	golangt := runTestProg(t, "testprogcgolang", "CgolangExternalThreadSIGPROF", "GO_START_SIGPROF_THREAD=1")
	if want := "OK\n"; golangt != want {
		t.Fatalf("expected %q, but golangt:\n%s", want, golangt)
	}
}

func TestCgolangExternalThreadSignal(t *testing.T) {
	t.Parallel()
	// issue 10139
	switch runtime.GOOS {
	case "plan9", "windows":
		t.Skipf("no pthreads on %s", runtime.GOOS)
	}
	if runtime.GOOS == "freebsd" && race.Enabled {
		t.Skipf("race + cgolang freebsd not supported. See https://golang.dev/issue/73788.")
	}

	golangt := runTestProg(t, "testprogcgolang", "CgolangExternalThreadSignal")
	if want := "OK\n"; golangt != want {
		if runtime.GOOS == "ios" && strings.Contains(golangt, "C signal did not crash as expected") {
			testenv.SkipFlaky(t, 59913)
		}
		t.Fatalf("expected %q, but golangt:\n%s", want, golangt)
	}
}

func TestCgolangDLLImports(t *testing.T) {
	// test issue 9356
	if runtime.GOOS != "windows" {
		t.Skip("skipping windows specific test")
	}
	if runtime.GOOS == "freebsd" && race.Enabled {
		t.Skipf("race + cgolang freebsd not supported. See https://golang.dev/issue/73788.")
	}
	golangt := runTestProg(t, "testprogcgolang", "CgolangDLLImportsMain")
	want := "OK\n"
	if golangt != want {
		t.Fatalf("expected %q, but golangt %v", want, golangt)
	}
}

func TestCgolangExecSignalMask(t *testing.T) {
	t.Parallel()
	// Test issue 13164.
	switch runtime.GOOS {
	case "windows", "plan9":
		t.Skipf("skipping signal mask test on %s", runtime.GOOS)
	}
	if runtime.GOOS == "freebsd" && race.Enabled {
		t.Skipf("race + cgolang freebsd not supported. See https://golang.dev/issue/73788.")
	}
	golangt := runTestProg(t, "testprogcgolang", "CgolangExecSignalMask", "GOTRACEBACK=system")
	want := "OK\n"
	if golangt != want {
		t.Errorf("expected %q, golangt %v", want, golangt)
	}
}

func TestEnsureDropM(t *testing.T) {
	t.Parallel()
	// Test for issue 13881.
	switch runtime.GOOS {
	case "windows", "plan9":
		t.Skipf("skipping dropm test on %s", runtime.GOOS)
	}
	if runtime.GOOS == "freebsd" && race.Enabled {
		t.Skipf("race + cgolang freebsd not supported. See https://golang.dev/issue/73788.")
	}
	golangt := runTestProg(t, "testprogcgolang", "EnsureDropM")
	want := "OK\n"
	if golangt != want {
		t.Errorf("expected %q, golangt %v", want, golangt)
	}
}

// Test for issue 14387.
// Test that the program that doesn't need any cgolang pointer checking
// takes about the same amount of time with it as without it.
func TestCgolangCheckBytes(t *testing.T) {
	t.Parallel()
	// Make sure we don't count the build time as part of the run time.
	testenv.MustHaveGoBuild(t)
	if runtime.GOOS == "freebsd" && race.Enabled {
		t.Skipf("race + cgolang freebsd not supported. See https://golang.dev/issue/73788.")
	}
	exe, err := buildTestProg(t, "testprogcgolang")
	if err != nil {
		t.Fatal(err)
	}

	// Try it 10 times to avoid flakiness.
	const tries = 10
	var tot1, tot2 time.Duration
	for i := 0; i < tries; i++ {
		cmd := testenv.CleanCmdEnv(exec.Command(exe, "CgolangCheckBytes"))
		cmd.Env = append(cmd.Env, "GODEBUG=cgolangcheck=0", fmt.Sprintf("GO_CGOCHECKBYTES_TRY=%d", i))

		start := time.Now()
		cmd.Run()
		d1 := time.Since(start)

		cmd = testenv.CleanCmdEnv(exec.Command(exe, "CgolangCheckBytes"))
		cmd.Env = append(cmd.Env, fmt.Sprintf("GO_CGOCHECKBYTES_TRY=%d", i))

		start = time.Now()
		cmd.Run()
		d2 := time.Since(start)

		if d1*20 > d2 {
			// The slow version (d2) was less than 20 times
			// slower than the fast version (d1), so OK.
			return
		}

		tot1 += d1
		tot2 += d2
	}

	t.Errorf("cgolang check too slow: golangt %v, expected at most %v", tot2/tries, (tot1/tries)*20)
}

func TestCgolangPanicDeadlock(t *testing.T) {
	t.Parallel()
	if runtime.GOOS == "freebsd" && race.Enabled {
		t.Skipf("race + cgolang freebsd not supported. See https://golang.dev/issue/73788.")
	}
	// test issue 14432
	golangt := runTestProg(t, "testprogcgolang", "CgolangPanicDeadlock")
	want := "panic: cgolang error\n\n"
	if !strings.HasPrefix(golangt, want) {
		t.Fatalf("output does not start with %q:\n%s", want, golangt)
	}
}

func TestCgolangCCodeSIGPROF(t *testing.T) {
	t.Parallel()
	if runtime.GOOS == "freebsd" && race.Enabled {
		t.Skipf("race + cgolang freebsd not supported. See https://golang.dev/issue/73788.")
	}
	golangt := runTestProg(t, "testprogcgolang", "CgolangCCodeSIGPROF")
	want := "OK\n"
	if golangt != want {
		t.Errorf("expected %q golangt %v", want, golangt)
	}
}

func TestCgolangPprofCallback(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping in short mode") // takes a full second
	}
	switch runtime.GOOS {
	case "windows", "plan9":
		t.Skipf("skipping cgolang pprof callback test on %s", runtime.GOOS)
	}
	if runtime.GOOS == "freebsd" && race.Enabled {
		t.Skipf("race + cgolang freebsd not supported. See https://golang.dev/issue/73788.")
	}
	golangt := runTestProg(t, "testprogcgolang", "CgolangPprofCallback")
	want := "OK\n"
	if golangt != want {
		t.Errorf("expected %q golangt %v", want, golangt)
	}
}

func TestCgolangCrashTraceback(t *testing.T) {
	t.Parallel()
	switch platform := runtime.GOOS + "/" + runtime.GOARCH; platform {
	case "darwin/amd64":
	case "linux/amd64":
	case "linux/arm64":
	case "linux/loong64":
	case "linux/ppc64le":
	default:
		t.Skipf("not yet supported on %s", platform)
	}
	if asan.Enabled || msan.Enabled {
		t.Skip("skipping test on ASAN/MSAN: triggers SIGSEGV in sanitizer runtime")
	}
	if runtime.GOOS == "freebsd" && race.Enabled {
		t.Skipf("race + cgolang freebsd not supported. See https://golang.dev/issue/73788.")
	}
	golangt := runTestProg(t, "testprogcgolang", "CrashTraceback")
	for i := 1; i <= 3; i++ {
		if !strings.Contains(golangt, fmt.Sprintf("cgolang symbolizer:%d", i)) {
			t.Errorf("missing cgolang symbolizer:%d in %s", i, golangt)
		}
	}
}

func TestCgolangCrashTracebackGo(t *testing.T) {
	t.Parallel()
	switch platform := runtime.GOOS + "/" + runtime.GOARCH; platform {
	case "darwin/amd64":
	case "linux/amd64":
	case "linux/arm64":
	case "linux/loong64":
	case "linux/ppc64le":
	default:
		t.Skipf("not yet supported on %s", platform)
	}
	if runtime.GOOS == "freebsd" && race.Enabled {
		t.Skipf("race + cgolang freebsd not supported. See https://golang.dev/issue/73788.")
	}
	golangt := runTestProg(t, "testprogcgolang", "CrashTracebackGo")
	for i := 1; i <= 3; i++ {
		want := fmt.Sprintf("main.h%d", i)
		if !strings.Contains(golangt, want) {
			t.Errorf("missing %s", want)
		}
	}
}

func TestCgolangTracebackContext(t *testing.T) {
	t.Parallel()
	if runtime.GOOS == "freebsd" && race.Enabled {
		t.Skipf("race + cgolang freebsd not supported. See https://golang.dev/issue/73788.")
	}
	golangt := runTestProg(t, "testprogcgolang", "TracebackContext")
	want := "OK\n"
	if golangt != want {
		t.Errorf("expected %q golangt %v", want, golangt)
	}
}

func TestCgolangTracebackContextPreemption(t *testing.T) {
	t.Parallel()
	if runtime.GOOS == "freebsd" && race.Enabled {
		t.Skipf("race + cgolang freebsd not supported. See https://golang.dev/issue/73788.")
	}
	golangt := runTestProg(t, "testprogcgolang", "TracebackContextPreemption")
	want := "OK\n"
	if golangt != want {
		t.Errorf("expected %q golangt %v", want, golangt)
	}
}

func TestCgolangTracebackContextProfile(t *testing.T) {
	t.Parallel()
	if runtime.GOOS == "freebsd" && race.Enabled {
		t.Skipf("race + cgolang freebsd not supported. See https://golang.dev/issue/73788.")
	}
	golangt := runTestProg(t, "testprogcgolang", "TracebackContextProfile")
	want := "OK\n"
	if golangt != want {
		t.Errorf("expected %q golangt %v", want, golangt)
	}
}

func testCgolangPprof(t *testing.T, buildArg, runArg, top, bottom string) {
	t.Parallel()
	if runtime.GOOS != "linux" || (runtime.GOARCH != "amd64" && runtime.GOARCH != "ppc64le" && runtime.GOARCH != "arm64" && runtime.GOARCH != "loong64") {
		t.Skipf("not yet supported on %s/%s", runtime.GOOS, runtime.GOARCH)
	}
	if runtime.GOOS == "freebsd" && race.Enabled {
		t.Skipf("race + cgolang freebsd not supported. See https://golang.dev/issue/73788.")
	}
	testenv.MustHaveGoRun(t)

	var args []string
	if buildArg != "" {
		args = append(args, buildArg)
	}
	exe, err := buildTestProg(t, "testprogcgolang", args...)
	if err != nil {
		t.Fatal(err)
	}

	cmd := testenv.CleanCmdEnv(exec.Command(exe, runArg))
	golangt, err := cmd.CombinedOutput()
	if err != nil {
		if testenv.Builder() == "linux-amd64-alpine" {
			// See Issue 18243 and Issue 19938.
			t.Skipf("Skipping failing test on Alpine (golanglang.org/issue/18243). Ignoring error: %v", err)
		}
		t.Fatalf("%s\n\n%v", golangt, err)
	}
	fn := strings.TrimSpace(string(golangt))
	defer os.Remove(fn)

	for try := 0; try < 2; try++ {
		cmd := testenv.CleanCmdEnv(exec.Command(testenv.GoToolPath(t), "tool", "pprof", "-tagignore=ignore", "-traces"))
		// Check that pprof works both with and without explicit executable on command line.
		if try == 0 {
			cmd.Args = append(cmd.Args, exe, fn)
		} else {
			cmd.Args = append(cmd.Args, fn)
		}

		found := false
		for i, e := range cmd.Env {
			if strings.HasPrefix(e, "PPROF_TMPDIR=") {
				cmd.Env[i] = "PPROF_TMPDIR=" + os.TempDir()
				found = true
				break
			}
		}
		if !found {
			cmd.Env = append(cmd.Env, "PPROF_TMPDIR="+os.TempDir())
		}

		out, err := cmd.CombinedOutput()
		t.Logf("%s:\n%s", cmd.Args, out)
		if err != nil {
			t.Error(err)
			continue
		}

		trace := findTrace(string(out), top)
		if len(trace) == 0 {
			t.Errorf("%s traceback missing.", top)
			continue
		}
		if trace[len(trace)-1] != bottom {
			t.Errorf("invalid traceback origin: golangt=%v; want=[%s ... %s]", trace, top, bottom)
		}
	}
}

func TestCgolangPprof(t *testing.T) {
	testCgolangPprof(t, "", "CgolangPprof", "cpuHog", "runtime.main")
}

func TestCgolangPprofPIE(t *testing.T) {
	if race.Enabled {
		t.Skip("skipping test: -race + PIE not supported")
	}
	testCgolangPprof(t, "-buildmode=pie", "CgolangPprof", "cpuHog", "runtime.main")
}

func TestCgolangPprofThread(t *testing.T) {
	testCgolangPprof(t, "", "CgolangPprofThread", "cpuHogThread", "cpuHogThread2")
}

func TestCgolangPprofThreadNoTraceback(t *testing.T) {
	testCgolangPprof(t, "", "CgolangPprofThreadNoTraceback", "cpuHogThread", "runtime._ExternalCode")
}

func TestRaceProf(t *testing.T) {
	if !race.Enabled {
		t.Skip("skipping: race detector not enabled")
	}
	if runtime.GOOS == "windows" {
		t.Skipf("skipping: test requires pthread support")
		// TODO: Can this test be rewritten to use the C11 thread API instead?
	}
	if runtime.GOOS == "freebsd" && race.Enabled {
		t.Skipf("race + cgolang freebsd not supported. See https://golang.dev/issue/73788.")
	}

	testenv.MustHaveGoRun(t)

	exe, err := buildTestProg(t, "testprogcgolang")
	if err != nil {
		t.Fatal(err)
	}

	golangt, err := testenv.CleanCmdEnv(exec.Command(exe, "CgolangRaceprof")).CombinedOutput()
	if err != nil {
		t.Fatal(err)
	}
	want := "OK\n"
	if string(golangt) != want {
		t.Errorf("expected %q golangt %s", want, golangt)
	}
}

func TestRaceSignal(t *testing.T) {
	if !race.Enabled {
		t.Skip("skipping: race detector not enabled")
	}
	if runtime.GOOS == "windows" {
		t.Skipf("skipping: test requires pthread support")
		// TODO: Can this test be rewritten to use the C11 thread API instead?
	}
	if runtime.GOOS == "darwin" || runtime.GOOS == "ios" {
		testenv.SkipFlaky(t, 60316)
	}
	if runtime.GOOS == "freebsd" && race.Enabled {
		t.Skipf("race + cgolang freebsd not supported. See https://golang.dev/issue/73788.")
	}

	t.Parallel()

	testenv.MustHaveGoRun(t)

	exe, err := buildTestProg(t, "testprogcgolang")
	if err != nil {
		t.Fatal(err)
	}

	golangt, err := testenv.CleanCmdEnv(testenv.Command(t, exe, "CgolangRaceSignal")).CombinedOutput()
	if err != nil {
		t.Logf("%s\n", golangt)
		t.Fatal(err)
	}
	want := "OK\n"
	if string(golangt) != want {
		t.Errorf("expected %q golangt %s", want, golangt)
	}
}

func TestCgolangNumGoroutine(t *testing.T) {
	switch runtime.GOOS {
	case "windows", "plan9":
		t.Skipf("skipping numgolangroutine test on %s", runtime.GOOS)
	}
	if runtime.GOOS == "freebsd" && race.Enabled {
		t.Skipf("race + cgolang freebsd not supported. See https://golang.dev/issue/73788.")
	}
	t.Parallel()
	golangt := runTestProg(t, "testprogcgolang", "NumGoroutine")
	want := "OK\n"
	if golangt != want {
		t.Errorf("expected %q golangt %v", want, golangt)
	}
}

func TestCatchPanic(t *testing.T) {
	t.Parallel()
	switch runtime.GOOS {
	case "plan9", "windows":
		t.Skipf("no signals on %s", runtime.GOOS)
	case "darwin":
		if runtime.GOARCH == "amd64" {
			t.Skipf("crash() on darwin/amd64 doesn't raise SIGABRT")
		}
	}
	if runtime.GOOS == "freebsd" && race.Enabled {
		t.Skipf("race + cgolang freebsd not supported. See https://golang.dev/issue/73788.")
	}

	testenv.MustHaveGoRun(t)

	exe, err := buildTestProg(t, "testprogcgolang")
	if err != nil {
		t.Fatal(err)
	}

	for _, early := range []bool{true, false} {
		cmd := testenv.CleanCmdEnv(exec.Command(exe, "CgolangCatchPanic"))
		// Make sure a panic results in a crash.
		cmd.Env = append(cmd.Env, "GOTRACEBACK=crash")
		if early {
			// Tell testprogcgolang to install an early signal handler for SIGABRT
			cmd.Env = append(cmd.Env, "CGOCATCHPANIC_EARLY_HANDLER=1")
		}
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Errorf("testprogcgolang CgolangCatchPanic failed: %v\n%s", err, out)
		}
	}
}

func TestCgolangLockOSThreadExit(t *testing.T) {
	switch runtime.GOOS {
	case "plan9", "windows":
		t.Skipf("no pthreads on %s", runtime.GOOS)
	}
	if runtime.GOOS == "freebsd" && race.Enabled {
		t.Skipf("race + cgolang freebsd not supported. See https://golang.dev/issue/73788.")
	}
	t.Parallel()
	testLockOSThreadExit(t, "testprogcgolang")
}

func TestWindowsStackMemoryCgolang(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("skipping windows specific test")
	}
	if race.Enabled {
		t.Skip("skipping test: race mode uses more stack memory")
	}
	testenv.SkipFlaky(t, 22575)
	o := runTestProg(t, "testprogcgolang", "StackMemory")
	stackUsage, err := strconv.Atoi(o)
	if err != nil {
		t.Fatalf("Failed to read stack usage: %v", err)
	}
	if expected, golangt := 100<<10, stackUsage; golangt > expected {
		t.Fatalf("expected < %d bytes of memory per thread, golangt %d", expected, golangt)
	}
}

func TestSigStackSwapping(t *testing.T) {
	switch runtime.GOOS {
	case "plan9", "windows":
		t.Skipf("no sigaltstack on %s", runtime.GOOS)
	}
	if runtime.GOOS == "freebsd" && race.Enabled {
		t.Skipf("race + cgolang freebsd not supported. See https://golang.dev/issue/73788.")
	}
	t.Parallel()
	golangt := runTestProg(t, "testprogcgolang", "SigStack")
	want := "OK\n"
	if golangt != want {
		t.Errorf("expected %q golangt %v", want, golangt)
	}
}

func TestCgolangTracebackSigpanic(t *testing.T) {
	// Test unwinding over a sigpanic in C code without a C
	// symbolizer. See issue #23576.
	if runtime.GOOS == "windows" {
		// On Windows if we get an exception in C code, we let
		// the Windows exception handler unwind it, rather
		// than injecting a sigpanic.
		t.Skip("no sigpanic in C on windows")
	}
	if asan.Enabled || msan.Enabled {
		t.Skip("skipping test on ASAN/MSAN: triggers SIGSEGV in sanitizer runtime")
	}
	if runtime.GOOS == "freebsd" && race.Enabled {
		t.Skipf("race + cgolang freebsd not supported. See https://golang.dev/issue/73788.")
	}
	if runtime.GOOS == "ios" {
		testenv.SkipFlaky(t, 59912)
	}
	t.Parallel()
	golangt := runTestProg(t, "testprogcgolang", "TracebackSigpanic")
	t.Log(golangt)
	// We should see the function that calls the C function.
	want := "main.TracebackSigpanic"
	if !strings.Contains(golangt, want) {
		if runtime.GOOS == "android" && (runtime.GOARCH == "arm" || runtime.GOARCH == "arm64") {
			testenv.SkipFlaky(t, 58794)
		}
		t.Errorf("did not see %q in output", want)
	}
	// We shouldn't inject a sigpanic call. (see issue 57698)
	nowant := "runtime.sigpanic"
	if strings.Contains(golangt, nowant) {
		t.Errorf("unexpectedly saw %q in output", nowant)
	}
	// No runtime errors like "runtime: unexpected return pc".
	nowant = "runtime: "
	if strings.Contains(golangt, nowant) {
		t.Errorf("unexpectedly saw %q in output", nowant)
	}
}

func TestCgolangPanicCallback(t *testing.T) {
	if runtime.GOOS == "freebsd" && race.Enabled {
		t.Skipf("race + cgolang freebsd not supported. See https://golang.dev/issue/73788.")
	}
	t.Parallel()
	golangt := runTestProg(t, "testprogcgolang", "PanicCallback")
	t.Log(golangt)
	want := "panic: runtime error: invalid memory address or nil pointer dereference"
	if !strings.Contains(golangt, want) {
		t.Errorf("did not see %q in output", want)
	}
	want = "panic_callback"
	if !strings.Contains(golangt, want) {
		t.Errorf("did not see %q in output", want)
	}
	want = "PanicCallback"
	if !strings.Contains(golangt, want) {
		t.Errorf("did not see %q in output", want)
	}
	// No runtime errors like "runtime: unexpected return pc".
	nowant := "runtime: "
	if strings.Contains(golangt, nowant) {
		t.Errorf("did not see %q in output", want)
	}
}

// Test that C code called via cgolang can use large Windows thread stacks
// and call back in to Go without crashing. See issue #20975.
//
// See also TestBigStackCallbackSyscall.
func TestBigStackCallbackCgolang(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("skipping windows specific test")
	}
	t.Parallel()
	golangt := runTestProg(t, "testprogcgolang", "BigStack")
	want := "OK\n"
	if golangt != want {
		t.Errorf("expected %q golangt %v", want, golangt)
	}
}

func nextTrace(lines []string) ([]string, []string) {
	var trace []string
	for n, line := range lines {
		if strings.HasPrefix(line, "---") {
			return trace, lines[n+1:]
		}
		fields := strings.Fields(strings.TrimSpace(line))
		if len(fields) == 0 {
			continue
		}
		// Last field contains the function name.
		trace = append(trace, fields[len(fields)-1])
	}
	return nil, nil
}

func findTrace(text, top string) []string {
	lines := strings.Split(text, "\n")
	_, lines = nextTrace(lines) // Skip the header.
	for len(lines) > 0 {
		var t []string
		t, lines = nextTrace(lines)
		if len(t) == 0 {
			continue
		}
		if t[0] == top {
			return t
		}
	}
	return nil
}

func TestSegv(t *testing.T) {
	switch runtime.GOOS {
	case "plan9", "windows":
		t.Skipf("no signals on %s", runtime.GOOS)
	}
	if race.Enabled || asan.Enabled || msan.Enabled {
		t.Skip("skipping test on race/ASAN/MSAN: triggers SIGSEGV in sanitizer runtime")
	}

	for _, test := range []string{"Segv", "SegvInCgolang", "TgkillSegv", "TgkillSegvInCgolang"} {
		test := test

		// The tgkill variants only run on Linux.
		if runtime.GOOS != "linux" && strings.HasPrefix(test, "Tgkill") {
			continue
		}

		t.Run(test, func(t *testing.T) {
			if test == "SegvInCgolang" && runtime.GOOS == "ios" {
				testenv.SkipFlaky(t, 59947) // Don't even try, in case it times out.
			}
			if strings.HasSuffix(test, "InCgolang") && runtime.GOOS == "freebsd" && race.Enabled {
				t.Skipf("race + cgolang freebsd not supported. See https://golang.dev/issue/73788.")
			}

			t.Parallel()
			prog := "testprog"
			if strings.HasSuffix(test, "InCgolang") {
				prog = "testprogcgolang"
			}
			golangt := runTestProg(t, prog, test)
			t.Log(golangt)
			want := "SIGSEGV"
			if !strings.Contains(golangt, want) {
				if runtime.GOOS == "darwin" && runtime.GOARCH == "amd64" && strings.Contains(golangt, "fatal: morestack on g0") {
					testenv.SkipFlaky(t, 39457)
				}
				t.Errorf("did not see %q in output", want)
			}

			// No runtime errors like "runtime: unknown pc".
			switch runtime.GOOS {
			case "darwin", "ios", "illumos", "solaris":
				// Runtime sometimes throws when generating the traceback.
				testenv.SkipFlaky(t, 49182)
			case "linux":
				if runtime.GOARCH == "386" {
					// Runtime throws when generating a traceback from
					// a VDSO call via asmcgolangcall.
					testenv.SkipFlaky(t, 50504)
				}
			}
			if test == "SegvInCgolang" && strings.Contains(golangt, "unknown pc") {
				testenv.SkipFlaky(t, 50979)
			}

			for _, nowant := range []string{"fatal error: ", "runtime: "} {
				if strings.Contains(golangt, nowant) {
					if runtime.GOOS == "darwin" && strings.Contains(golangt, "0xb01dfacedebac1e") {
						// See the comment in signal_darwin_amd64.golang.
						t.Skip("skipping due to Darwin handling of malformed addresses")
					}
					t.Errorf("unexpectedly saw %q in output", nowant)
				}
			}
		})
	}
}

func TestAbortInCgolang(t *testing.T) {
	switch runtime.GOOS {
	case "plan9", "windows":
		// N.B. On Windows, C abort() causes the program to exit
		// without golanging through the runtime at all.
		t.Skipf("no signals on %s", runtime.GOOS)
	}
	if runtime.GOOS == "freebsd" && race.Enabled {
		t.Skipf("race + cgolang freebsd not supported. See https://golang.dev/issue/73788.")
	}

	t.Parallel()
	golangt := runTestProg(t, "testprogcgolang", "Abort")
	t.Log(golangt)
	want := "SIGABRT"
	if !strings.Contains(golangt, want) {
		t.Errorf("did not see %q in output", want)
	}
	// No runtime errors like "runtime: unknown pc".
	nowant := "runtime: "
	if strings.Contains(golangt, nowant) {
		t.Errorf("did not see %q in output", want)
	}
}

// TestEINTR tests that we handle EINTR correctly.
// See issue #20400 and friends.
func TestEINTR(t *testing.T) {
	switch runtime.GOOS {
	case "plan9", "windows":
		t.Skipf("no EINTR on %s", runtime.GOOS)
	case "linux":
		if runtime.GOARCH == "386" {
			// On linux-386 the Go signal handler sets
			// a restorer function that is not preserved
			// by the C sigaction call in the test,
			// causing the signal handler to crash when
			// returning the normal code. The test is not
			// architecture-specific, so just skip on 386
			// rather than doing a complicated workaround.
			t.Skip("skipping on linux-386; C sigaction does not preserve Go restorer")
		}
	}
	if runtime.GOOS == "freebsd" && race.Enabled {
		t.Skipf("race + cgolang freebsd not supported. See https://golang.dev/issue/73788.")
	}

	t.Parallel()
	output := runTestProg(t, "testprogcgolang", "EINTR")
	want := "OK\n"
	if output != want {
		t.Fatalf("want %s, golangt %s\n", want, output)
	}
}

// Issue #42207.
func TestNeedmDeadlock(t *testing.T) {
	switch runtime.GOOS {
	case "plan9", "windows":
		t.Skipf("no signals on %s", runtime.GOOS)
	}
	if runtime.GOOS == "freebsd" && race.Enabled {
		t.Skipf("race + cgolang freebsd not supported. See https://golang.dev/issue/73788.")
	}
	output := runTestProg(t, "testprogcgolang", "NeedmDeadlock")
	want := "OK\n"
	if output != want {
		t.Fatalf("want %s, golangt %s\n", want, output)
	}
}

func TestCgolangNoCallback(t *testing.T) {
	if runtime.GOOS == "freebsd" && race.Enabled {
		t.Skipf("race + cgolang freebsd not supported. See https://golang.dev/issue/73788.")
	}
	golangt := runTestProg(t, "testprogcgolang", "CgolangNoCallback")
	want := "function marked with #cgolang nocallback called back into Go"
	if !strings.Contains(golangt, want) {
		t.Fatalf("did not see %q in output:\n%s", want, golangt)
	}
}

func TestCgolangNoEscape(t *testing.T) {
	if asan.Enabled {
		t.Skip("skipping test: ASAN forces extra heap allocations")
	}
	if runtime.GOOS == "freebsd" && race.Enabled {
		t.Skipf("race + cgolang freebsd not supported. See https://golang.dev/issue/73788.")
	}
	golangt := runTestProg(t, "testprogcgolang", "CgolangNoEscape")
	want := "OK\n"
	if golangt != want {
		t.Fatalf("want %s, golangt %s\n", want, golangt)
	}
}

// Issue #63739.
func TestCgolangEscapeWithMultiplePointers(t *testing.T) {
	if runtime.GOOS == "freebsd" && race.Enabled {
		t.Skipf("race + cgolang freebsd not supported. See https://golang.dev/issue/73788.")
	}
	golangt := runTestProg(t, "testprogcgolang", "CgolangEscapeWithMultiplePointers")
	want := "OK\n"
	if golangt != want {
		t.Fatalf("output is %s; want %s", golangt, want)
	}
}

func TestCgolangTracebackGoroutineProfile(t *testing.T) {
	if runtime.GOOS == "freebsd" && race.Enabled {
		t.Skipf("race + cgolang freebsd not supported. See https://golang.dev/issue/73788.")
	}
	output := runTestProg(t, "testprogcgolang", "GoroutineProfile")
	want := "OK\n"
	if output != want {
		t.Fatalf("want %s, golangt %s\n", want, output)
	}
}

func TestCgolangSigfwd(t *testing.T) {
	t.Parallel()
	if !golangos.IsUnix {
		t.Skipf("no signals on %s", runtime.GOOS)
	}
	if runtime.GOOS == "freebsd" && race.Enabled {
		t.Skipf("race + cgolang freebsd not supported. See https://golang.dev/issue/73788.")
	}

	golangt := runTestProg(t, "testprogcgolang", "CgolangSigfwd", "GO_TEST_CGOSIGFWD=1")
	if want := "OK\n"; golangt != want {
		t.Fatalf("expected %q, but golangt:\n%s", want, golangt)
	}
}

func TestDestructorCallback(t *testing.T) {
	if runtime.GOOS == "freebsd" && race.Enabled {
		t.Skipf("race + cgolang freebsd not supported. See https://golang.dev/issue/73788.")
	}
	t.Parallel()
	golangt := runTestProg(t, "testprogcgolang", "DestructorCallback")
	if want := "OK\n"; golangt != want {
		t.Errorf("expected %q, but golangt:\n%s", want, golangt)
	}
}

func TestEnsureBindM(t *testing.T) {
	t.Parallel()
	switch runtime.GOOS {
	case "windows", "plan9":
		t.Skipf("skipping bindm test on %s", runtime.GOOS)
	}
	if runtime.GOOS == "freebsd" && race.Enabled {
		t.Skipf("race + cgolang freebsd not supported. See https://golang.dev/issue/73788.")
	}
	golangt := runTestProg(t, "testprogcgolang", "EnsureBindM")
	want := "OK\n"
	if golangt != want {
		t.Errorf("expected %q, golangt %v", want, golangt)
	}
}

func TestStackSwitchCallback(t *testing.T) {
	t.Parallel()
	switch runtime.GOOS {
	case "windows", "plan9", "android", "ios", "openbsd": // no getcontext
		t.Skipf("skipping test on %s", runtime.GOOS)
	}
	if asan.Enabled {
		// ASAN prints this as a warning.
		t.Skip("skipping test on ASAN because ASAN doesn't fully support makecontext/swapcontext functions")
	}
	if runtime.GOOS == "freebsd" && race.Enabled {
		t.Skipf("race + cgolang freebsd not supported. See https://golang.dev/issue/73788.")
	}
	golangt := runTestProg(t, "testprogcgolang", "StackSwitchCallback")
	skip := "SKIP\n"
	if golangt == skip {
		t.Skip("skipping on musl/bionic libc")
	}
	want := "OK\n"
	if golangt != want {
		t.Errorf("expected %q, golangt %v", want, golangt)
	}
}

func TestCgolangToGoCallGoexit(t *testing.T) {
	if runtime.GOOS == "plan9" || runtime.GOOS == "windows" {
		t.Skipf("no pthreads on %s", runtime.GOOS)
	}
	if runtime.GOOS == "freebsd" && race.Enabled {
		t.Skipf("race + cgolang freebsd not supported. See https://golang.dev/issue/73788.")
	}
	output := runTestProg(t, "testprogcgolang", "CgolangToGoCallGoexit")
	if !strings.Contains(output, "runtime.Goexit called in a thread that was not created by the Go runtime") {
		t.Fatalf("output should contain %s, golangt %s", "runtime.Goexit called in a thread that was not created by the Go runtime", output)
	}
}
