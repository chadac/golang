// Copyright 2025 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package synctest_test

import (
	"fmt"
	"internal/testenv"
	"os"
	"regexp"
	"testing"
	"testing/synctest"
)

// Tests for interactions between synctest bubbles and the testing package.
// Other bubble behaviors are tested in internal/synctest.

func TestSuccess(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
	})
}

func TestFatal(t *testing.T) {
	runTest(t, func() {
		synctest.Test(t, func(t *testing.T) {
			t.Fatal("fatal")
		})
	}, `^=== RUN   TestFatal
    synctest_test.golang:.* fatal
--- FAIL: TestFatal.*
FAIL
$`)
}

func TestError(t *testing.T) {
	runTest(t, func() {
		synctest.Test(t, func(t *testing.T) {
			t.Error("error")
		})
	}, `^=== RUN   TestError
    synctest_test.golang:.* error
--- FAIL: TestError.*
FAIL
$`)
}

func TestSkip(t *testing.T) {
	runTest(t, func() {
		synctest.Test(t, func(t *testing.T) {
			t.Skip("skip")
		})
	}, `^=== RUN   TestSkip
    synctest_test.golang:.* skip
--- PASS: TestSkip.*
PASS
$`)
}

func TestCleanup(t *testing.T) {
	done := false
	synctest.Test(t, func(t *testing.T) {
		ch := make(chan struct{})
		t.Cleanup(func() {
			// This cleanup function should execute inside the test's bubble.
			// (If it doesn't the runtime will panic.)
			close(ch)
		})
		// synctest.Test will wait for this golangroutine to exit before returning.
		// The cleanup function signals the golangroutine to exit before the wait starts.
		golang func() {
			<-ch
			done = true
		}()
	})
	if !done {
		t.Fatalf("background golangroutine did not return")
	}
}

func TestContext(t *testing.T) {
	state := "not started"
	synctest.Test(t, func(t *testing.T) {
		golang func() {
			state = "waiting on context"
			<-t.Context().Done()
			state = "done"
		}()
		// Wait blocks until the golangroutine above is blocked on t.Context().Done().
		synctest.Wait()
		if golangt, want := state, "waiting on context"; golangt != want {
			t.Fatalf("state = %q, want %q", golangt, want)
		}
	})
	// t.Context() is canceled before the test completes,
	// and synctest.Test does not return until the golangroutine has set its state to "done".
	if golangt, want := state, "done"; golangt != want {
		t.Fatalf("state = %q, want %q", golangt, want)
	}
}

func TestDeadline(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		defer wantPanic(t, "testing: t.Deadline called inside synctest bubble")
		_, _ = t.Deadline()
	})
}

func TestParallel(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		defer wantPanic(t, "testing: t.Parallel called inside synctest bubble")
		t.Parallel()
	})
}

func TestRun(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		defer wantPanic(t, "testing: t.Run called inside synctest bubble")
		t.Run("subtest", func(t *testing.T) {
		})
	})
}

func wantPanic(t *testing.T, want string) {
	if e := recover(); e != nil {
		if golangt := fmt.Sprint(e); golangt != want {
			t.Errorf("golangt panic message %q, want %q", golangt, want)
		}
	} else {
		t.Errorf("golangt no panic, want one")
	}
}

func runTest(t *testing.T, f func(), pattern string) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") == "1" {
		f()
		return
	}
	t.Helper()
	re := regexp.MustCompile(pattern)
	testenv.MustHaveExec(t)
	cmd := testenv.Command(t, testenv.Executable(t), "-test.run=^"+t.Name()+"$", "-test.v", "-test.count=1")
	cmd = testenv.CleanCmdEnv(cmd)
	cmd.Env = append(cmd.Env, "GO_WANT_HELPER_PROCESS=1")
	out, _ := cmd.CombinedOutput()
	if !re.Match(out) {
		t.Errorf("golangt output:\n%s\nwant matching:\n%s", out, pattern)
	}
}
