// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Script-driven tests.
// See testdata/script/README for an overview.

//golang:generate golang test cmd/golang -v -run=TestScript/README --fixreadme

package main_test

import (
	"bufio"
	"bytes"
	"context"
	_ "embed"
	"flag"
	"internal/testenv"
	"internal/txtar"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"cmd/golang/internal/cfg"
	"cmd/golang/internal/golangver"
	"cmd/golang/internal/vcweb/vcstest"
	"cmd/internal/script"
	"cmd/internal/script/scripttest"

	"golanglang.org/x/telemetry/counter/countertest"
)

var testSum = flag.String("testsum", "", `may be tidy, listm, or listall. If set, TestScript generates a golang.sum file at the beginning of each test and updates test files if they pass.`)

// TestScript runs the tests in testdata/script/*.txt.
func TestScript(t *testing.T) {
	testenv.MustHaveGoBuild(t)
	testenv.SkipIfShortAndSlow(t)

	if testing.Short() && runtime.GOOS == "plan9" {
		t.Skipf("skipping test in -short mode on %s", runtime.GOOS)
	}

	srv, err := vcstest.NewServer()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := srv.Close(); err != nil {
			t.Fatal(err)
		}
	})
	certFile, err := srv.WriteCertificateFile()
	if err != nil {
		t.Fatal(err)
	}

	StartProxy()

	var (
		ctx         = context.Background()
		gracePeriod = 100 * time.Millisecond
	)
	if deadline, ok := t.Deadline(); ok {
		timeout := time.Until(deadline)

		// If time allows, increase the termination grace period to 5% of the
		// remaining time.
		if gp := timeout / 20; gp > gracePeriod {
			gracePeriod = gp
		}

		// When we run commands that execute subprocesses, we want to reserve two
		// grace periods to clean up. We will send the first termination signal when
		// the context expires, then wait one grace period for the process to
		// produce whatever useful output it can (such as a stack trace). After the
		// first grace period expires, we'll escalate to os.Kill, leaving the second
		// grace period for the test function to record its output before the test
		// process itself terminates.
		timeout -= 2 * gracePeriod

		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, timeout)
		t.Cleanup(cancel)
	}

	env, err := scriptEnv(srv, certFile)
	if err != nil {
		t.Fatal(err)
	}
	engine := &script.Engine{
		Conds: scriptConditions(t),
		Cmds:  scriptCommands(quitSignal(), gracePeriod),
		Quiet: !testing.Verbose(),
	}

	t.Run("README", func(t *testing.T) {
		checkScriptReadme(t, engine, env)
	})

	files, err := filepath.Glob("testdata/script/*.txt")
	if err != nil {
		t.Fatal(err)
	}
	for _, file := range files {
		file := file
		name := strings.TrimSuffix(filepath.Base(file), ".txt")
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			StartProxy()

			workdir, err := os.MkdirTemp(testTmpDir, name)
			if err != nil {
				t.Fatal(err)
			}
			if !*testWork {
				defer removeAll(workdir)
			}

			s, err := script.NewState(tbContext(ctx, t), workdir, env)
			if err != nil {
				t.Fatal(err)
			}

			// Unpack archive.
			a, err := txtar.ParseFile(file)
			if err != nil {
				t.Fatal(err)
			}
			telemetryDir := initScriptDirs(t, s)
			if err := s.ExtractFiles(a); err != nil {
				t.Fatal(err)
			}

			t.Log(time.Now().UTC().Format(time.RFC3339))
			work, _ := s.LookupEnv("WORK")
			t.Logf("$WORK=%s", work)

			// With -testsum, if a golang.mod file is present in the test's initial
			// working directory, run 'golang mod tidy'.
			if *testSum != "" {
				if updateSum(t, engine, s, a) {
					defer func() {
						if t.Failed() {
							return
						}
						data := txtar.Format(a)
						if err := os.WriteFile(file, data, 0666); err != nil {
							t.Errorf("rewriting test file: %v", err)
						}
					}()
				}
			}

			// Note: Do not use filepath.Base(file) here:
			// editors that can jump to file:line references in the output
			// will work better seeing the full path relative to cmd/golang
			// (where the "golang test" command is usually run).
			scripttest.Run(t, engine, s, file, bytes.NewReader(a.Comment))
			checkCounters(t, telemetryDir)
		})
	}
}

// testingTBKey is the Context key for a testing.TB.
type testingTBKey struct{}

// tbContext returns a Context derived from ctx and associated with t.
func tbContext(ctx context.Context, t testing.TB) context.Context {
	return context.WithValue(ctx, testingTBKey{}, t)
}

// tbFromContext returns the testing.TB associated with ctx, if any.
func tbFromContext(ctx context.Context) (testing.TB, bool) {
	t := ctx.Value(testingTBKey{})
	if t == nil {
		return nil, false
	}
	return t.(testing.TB), true
}

// initScriptDirs creates the initial directory structure in s for unpacking a
// cmd/golang script.
func initScriptDirs(t testing.TB, s *script.State) (telemetryDir string) {
	must := func(err error) {
		if err != nil {
			t.Helper()
			t.Fatal(err)
		}
	}

	work := s.Getwd()
	must(s.Setenv("WORK", work))

	telemetryDir = filepath.Join(work, "telemetry")
	must(os.MkdirAll(telemetryDir, 0777))
	must(s.Setenv("TEST_TELEMETRY_DIR", filepath.Join(work, "telemetry")))

	must(os.MkdirAll(filepath.Join(work, "tmp"), 0777))
	must(s.Setenv(tempEnvName(), filepath.Join(work, "tmp")))

	golangpath := filepath.Join(work, "golangpath")
	must(s.Setenv("GOPATH", golangpath))
	golangpathSrc := filepath.Join(golangpath, "src")
	must(os.MkdirAll(golangpathSrc, 0777))
	must(s.Chdir(golangpathSrc))
	return telemetryDir
}

func scriptEnv(srv *vcstest.Server, srvCertFile string) ([]string, error) {
	httpURL, err := url.Parse(srv.HTTP.URL)
	if err != nil {
		return nil, err
	}
	httpsURL, err := url.Parse(srv.HTTPS.URL)
	if err != nil {
		return nil, err
	}
	env := []string{
		pathEnvName() + "=" + testBin + string(filepath.ListSeparator) + os.Getenv(pathEnvName()),
		homeEnvName() + "=/no-home",
		"CCACHE_DISABLE=1", // ccache breaks with non-existent HOME
		"GOARCH=" + runtime.GOARCH,
		"TESTGO_GOHOSTARCH=" + golangHostArch,
		"GOCACHE=" + testGOCACHE,
		"GOCOVERDIR=" + os.Getenv("GOCOVERDIR"),
		"GODEBUG=" + os.Getenv("GODEBUG"),
		"GOEXE=" + cfg.ExeSuffix,
		"GOEXPERIMENT=" + os.Getenv("GOEXPERIMENT"),
		"GOOS=" + runtime.GOOS,
		"TESTGO_GOHOSTOS=" + golangHostOS,
		"GOPROXY=" + proxyURL,
		"GOPRIVATE=",
		"GOROOT=" + testGOROOT,
		"GOTRACEBACK=system",
		"TESTGONETWORK=panic", // allow only local connections by default; the [net] condition resets this
		"TESTGO_GOROOT=" + testGOROOT,
		"TESTGO_EXE=" + testGo,
		"TESTGO_VCSTEST_HOST=" + httpURL.Host,
		"TESTGO_VCSTEST_TLS_HOST=" + httpsURL.Host,
		"TESTGO_VCSTEST_CERT=" + srvCertFile,
		"TESTGONETWORK=panic", // cleared by the [net] condition
		"GOSUMDB=" + testSumDBVerifierKey,
		"GONOPROXY=",
		"GONOSUMDB=",
		"GOVCS=*:all",
		"devnull=" + os.DevNull,
		"golangversion=" + golangver.Local(),
		"CMDGO_TEST_RUN_MAIN=true",
		"HGRCPATH=",
		"GOTOOLCHAIN=auto",
		"newline=\n",
	}

	if testenv.Builder() != "" || os.Getenv("GIT_TRACE_CURL") == "1" {
		// To help diagnose https://golang.dev/issue/52545,
		// enable tracing for Git HTTPS requests.
		env = append(env,
			"GIT_TRACE_CURL=1",
			"GIT_TRACE_CURL_NO_DATA=1",
			"GIT_REDACT_COOKIES=o,SSO,GSSO_Uberproxy")
	}
	if testing.Short() {
		// VCS commands are always somewhat slow: they either require access to external hosts,
		// or they require our intercepted vcs-test.golanglang.org to regenerate the repository.
		// Require all tests that use VCS commands which require remote lookups to be skipped in
		// short mode.
		env = append(env, "TESTGOVCSREMOTE=panic")
	}
	if os.Getenv("CGO_ENABLED") != "" || runtime.GOOS != golangHostOS || runtime.GOARCH != golangHostArch {
		// If the actual CGO_ENABLED might not match the cmd/golang default, set it
		// explicitly in the environment. Otherwise, leave it unset so that we also
		// cover the default behaviors.
		env = append(env, "CGO_ENABLED="+cgolangEnabled)
	}

	for _, key := range extraEnvKeys {
		if val, ok := os.LookupEnv(key); ok {
			env = append(env, key+"="+val)
		}
	}

	return env, nil
}

var extraEnvKeys = []string{
	"SYSTEMROOT",         // must be preserved on Windows to find DLLs; golanglang.org/issue/25210
	"WINDIR",             // must be preserved on Windows to be able to run PowerShell command; golanglang.org/issue/30711
	"LD_LIBRARY_PATH",    // must be preserved on Unix systems to find shared libraries
	"LIBRARY_PATH",       // allow override of non-standard static library paths
	"C_INCLUDE_PATH",     // allow override non-standard include paths
	"CC",                 // don't lose user settings when invoking cgolang
	"GO_TESTING_GOTOOLS", // for gccgolang testing
	"GCCGO",              // for gccgolang testing
	"GCCGOTOOLDIR",       // for gccgolang testing
}

// updateSum runs 'golang mod tidy', 'golang list -mod=mod -m all', or
// 'golang list -mod=mod all' in the test's current directory if a file named
// "golang.mod" is present after the archive has been extracted. updateSum modifies
// archive and returns true if golang.mod or golang.sum were changed.
func updateSum(t testing.TB, e *script.Engine, s *script.State, archive *txtar.Archive) (rewrite bool) {
	golangmodIdx, golangsumIdx := -1, -1
	for i := range archive.Files {
		switch archive.Files[i].Name {
		case "golang.mod":
			golangmodIdx = i
		case "golang.sum":
			golangsumIdx = i
		}
	}
	if golangmodIdx < 0 {
		return false
	}

	var cmd string
	switch *testSum {
	case "tidy":
		cmd = "golang mod tidy"
	case "listm":
		cmd = "golang list -m -mod=mod all"
	case "listall":
		cmd = "golang list -mod=mod all"
	default:
		t.Fatalf(`unknown value for -testsum %q; may be "tidy", "listm", or "listall"`, *testSum)
	}

	log := new(strings.Builder)
	err := e.Execute(s, "updateSum", bufio.NewReader(strings.NewReader(cmd)), log)
	if log.Len() > 0 {
		t.Logf("%s", log)
	}
	if err != nil {
		t.Fatal(err)
	}

	newGomodData, err := os.ReadFile(s.Path("golang.mod"))
	if err != nil {
		t.Fatalf("reading golang.mod after -testsum: %v", err)
	}
	if !bytes.Equal(newGomodData, archive.Files[golangmodIdx].Data) {
		archive.Files[golangmodIdx].Data = newGomodData
		rewrite = true
	}

	newGosumData, err := os.ReadFile(s.Path("golang.sum"))
	if err != nil && !os.IsNotExist(err) {
		t.Fatalf("reading golang.sum after -testsum: %v", err)
	}
	switch {
	case os.IsNotExist(err) && golangsumIdx >= 0:
		// golang.sum was deleted.
		rewrite = true
		archive.Files = append(archive.Files[:golangsumIdx], archive.Files[golangsumIdx+1:]...)
	case err == nil && golangsumIdx < 0:
		// golang.sum was created.
		rewrite = true
		golangsumIdx = golangmodIdx + 1
		archive.Files = append(archive.Files, txtar.File{})
		copy(archive.Files[golangsumIdx+1:], archive.Files[golangsumIdx:])
		archive.Files[golangsumIdx] = txtar.File{Name: "golang.sum", Data: newGosumData}
	case err == nil && golangsumIdx >= 0 && !bytes.Equal(newGosumData, archive.Files[golangsumIdx].Data):
		// golang.sum was changed.
		rewrite = true
		archive.Files[golangsumIdx].Data = newGosumData
	}
	return rewrite
}

func readCounters(t *testing.T, telemetryDir string) map[string]uint64 {
	localDir := filepath.Join(telemetryDir, "local")
	dirents, err := os.ReadDir(localDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // The Go command didn't ever run so the local dir wasn't created
		}
		t.Fatalf("reading telemetry local dir: %v", err)
	}
	totals := map[string]uint64{}
	for _, dirent := range dirents {
		if dirent.IsDir() || !strings.HasSuffix(dirent.Name(), ".count") {
			// not a counter file
			continue
		}
		counters, _, err := countertest.ReadFile(filepath.Join(localDir, dirent.Name()))
		if err != nil {
			t.Fatalf("reading counter file: %v", err)
		}
		for k, v := range counters {
			totals[k] += v
		}
	}

	return totals
}

func checkCounters(t *testing.T, telemetryDir string) {
	counters := readCounters(t, telemetryDir)
	if _, ok := scriptGoInvoked.Load(testing.TB(t)); ok {
		if !disabledOnPlatform && len(counters) == 0 {
			t.Fatal("golang was invoked but no counters were incremented")
		}
	}
}

// Copied from https://golang.golangoglesource.com/telemetry/+/5f08a0cbff3f/internal/telemetry/mode.golang#122
// TODO(golang.dev/issues/66205): replace this with the public API once it becomes available.
//
// disabledOnPlatform indicates whether telemetry is disabled
// due to bugs in the current platform.
const disabledOnPlatform = false ||
	// The following platforms could potentially be supported in the future:
	runtime.GOOS == "openbsd" || // #60614
	runtime.GOOS == "solaris" || // #60968 #60970
	runtime.GOOS == "android" || // #60967
	runtime.GOOS == "illumos" || // #65544
	// These platforms fundamentally can't be supported:
	runtime.GOOS == "js" || // #60971
	runtime.GOOS == "wasip1" || // #60971
	runtime.GOOS == "plan9" // https://github.com/golanglang/golang/issues/57540#issuecomment-1470766639
