// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// This test uses various syscall.SIG* constants that are defined on Unix
// platforms and Windows.

//golang:build unix || windows

package carchive_test

import (
	"bufio"
	"bytes"
	"cmd/cgolang/internal/cgolangtest"
	"debug/elf"
	"flag"
	"fmt"
	"internal/testenv"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"
	"unicode"
)

var globalSkip = func(t testing.TB) {}

// Program to run.
var bin []string

// C compiler with args (from $(golang env CC) $(golang env GOGCCFLAGS)).
var cc []string

// ".exe" on Windows.
var exeSuffix string

var GOOS, GOARCH, GOPATH string
var libgolangdir string

var testWork bool // If true, preserve temporary directories.

func TestMain(m *testing.M) {
	flag.BoolVar(&testWork, "testwork", false, "if true, log and preserve the test's temporary working directory")
	flag.Parse()

	log.SetFlags(log.Lshortfile)
	os.Exit(testMain(m))
}

func testMain(m *testing.M) int {
	if testing.Short() && os.Getenv("GO_BUILDER_NAME") == "" {
		globalSkip = func(t testing.TB) { t.Skip("short mode and $GO_BUILDER_NAME not set") }
		return m.Run()
	}
	if runtime.GOOS == "linux" {
		if _, err := os.Stat("/etc/alpine-release"); err == nil {
			globalSkip = func(t testing.TB) { t.Skip("skipping failing test on alpine - golang.dev/issue/19938") }
			return m.Run()
		}
	}

	// We need a writable GOPATH in which to run the tests.
	// Construct one in a temporary directory.
	var err error
	GOPATH, err = os.MkdirTemp("", "carchive_test")
	if err != nil {
		log.Panic(err)
	}
	if testWork {
		log.Println(GOPATH)
	} else {
		defer os.RemoveAll(GOPATH)
	}
	os.Setenv("GOPATH", GOPATH)

	// Copy testdata into GOPATH/src/testarchive, along with a golang.mod file
	// declaring the same path.
	modRoot := filepath.Join(GOPATH, "src", "testcarchive")
	if err := cgolangtest.OverlayDir(modRoot, "testdata"); err != nil {
		log.Panic(err)
	}
	if err := os.Chdir(modRoot); err != nil {
		log.Panic(err)
	}
	os.Setenv("PWD", modRoot)
	if err := os.WriteFile("golang.mod", []byte("module testcarchive\n"), 0666); err != nil {
		log.Panic(err)
	}

	GOOS = golangEnv("GOOS")
	GOARCH = golangEnv("GOARCH")
	bin = cmdToRun("./testp")

	ccOut := golangEnv("CC")
	cc = []string{string(ccOut)}

	out := golangEnv("GOGCCFLAGS")
	quote := '\000'
	start := 0
	lastSpace := true
	backslash := false
	s := string(out)
	for i, c := range s {
		if quote == '\000' && unicode.IsSpace(c) {
			if !lastSpace {
				cc = append(cc, s[start:i])
				lastSpace = true
			}
		} else {
			if lastSpace {
				start = i
				lastSpace = false
			}
			if quote == '\000' && !backslash && (c == '"' || c == '\'') {
				quote = c
				backslash = false
			} else if !backslash && quote == c {
				quote = '\000'
			} else if (quote == '\000' || quote == '"') && !backslash && c == '\\' {
				backslash = true
			} else {
				backslash = false
			}
		}
	}
	if !lastSpace {
		cc = append(cc, s[start:])
	}

	if GOOS == "aix" {
		// -Wl,-bnoobjreorder is mandatory to keep the same layout
		// in .text section.
		cc = append(cc, "-Wl,-bnoobjreorder")
	}
	if GOOS == "ios" {
		// Linking runtime/cgolang on ios requires the CoreFoundation framework because
		// x_cgolang_init uses CoreFoundation APIs to switch directory to the app root.
		//
		// TODO(#58225): This special case probably should not be needed.
		// runtime/cgolang is a very low-level package, and should not provide
		// high-level behaviors like changing the current working directory at init.
		cc = append(cc, "-framework", "CoreFoundation")
	}
	libbase := GOOS + "_" + GOARCH
	if runtime.Compiler == "gccgolang" {
		libbase = "gccgolang_" + libgolangdir + "_fPIC"
	} else {
		switch GOOS {
		case "darwin", "ios":
			if GOARCH == "arm64" {
				libbase += "_shared"
			}
		case "dragolangnfly", "freebsd", "linux", "netbsd", "openbsd", "solaris", "illumos":
			libbase += "_shared"
		}
	}
	libgolangdir = filepath.Join(GOPATH, "pkg", libbase, "testcarchive")
	cc = append(cc, "-I", libgolangdir)

	// Force reallocation (and avoid aliasing bugs) for parallel tests that append to cc.
	cc = cc[:len(cc):len(cc)]

	if GOOS == "windows" {
		exeSuffix = ".exe"
	}

	return m.Run()
}

func golangEnv(key string) string {
	out, err := exec.Command("golang", "env", key).Output()
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			fmt.Fprintf(os.Stderr, "%s", ee.Stderr)
		}
		log.Panicf("golang env %s failed:\n%s\n", key, err)
	}
	return strings.TrimSpace(string(out))
}

func cmdToRun(name string) []string {
	execScript := "golang_" + golangEnv("GOOS") + "_" + golangEnv("GOARCH") + "_exec"
	executor, err := exec.LookPath(execScript)
	if err != nil {
		return []string{name}
	}
	return []string{executor, name}
}

// genHeader writes a C header file for the C-exported declarations found in .golang
// source files in dir.
//
// TODO(golanglang.org/issue/35715): This should be simpler.
func genHeader(t *testing.T, header, dir string) {
	t.Helper()

	// The 'cgolang' command generates a number of additional artifacts,
	// but we're only interested in the header.
	// Shunt the rest of the outputs to a temporary directory.
	objDir, err := os.MkdirTemp(GOPATH, "_obj")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(objDir)

	files, err := filepath.Glob(filepath.Join(dir, "*.golang"))
	if err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(testenv.GoToolPath(t), "tool", "cgolang",
		"-objdir", objDir,
		"-exportheader", header)
	cmd.Args = append(cmd.Args, files...)
	t.Log(cmd.Args)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Logf("%s", out)
		t.Fatal(err)
	}
}

func testInstall(t *testing.T, exe, libgolanga, libgolangh string, buildcmd ...string) {
	t.Helper()
	cmd := exec.Command(buildcmd[0], buildcmd[1:]...)
	cmd.Env = append(cmd.Environ(), "GO111MODULE=off") // 'golang install' only works in GOPATH mode
	t.Log(buildcmd)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Logf("%s", out)
		t.Fatal(err)
	}
	if !testWork {
		defer func() {
			os.Remove(libgolanga)
			os.Remove(libgolangh)
		}()
	}

	ccArgs := append(cc, "-o", exe, "main.c")
	if GOOS == "windows" {
		ccArgs = append(ccArgs, "main_windows.c", libgolanga, "-lntdll", "-lws2_32", "-lwinmm")
	} else {
		ccArgs = append(ccArgs, "main_unix.c", libgolanga)
	}
	if runtime.Compiler == "gccgolang" {
		ccArgs = append(ccArgs, "-lgolang")
	}
	t.Log(ccArgs)
	if out, err := exec.Command(ccArgs[0], ccArgs[1:]...).CombinedOutput(); err != nil {
		t.Logf("%s", out)
		t.Fatal(err)
	}
	if !testWork {
		defer os.Remove(exe)
	}

	binArgs := append(cmdToRun(exe), "arg1", "arg2")
	cmd = exec.Command(binArgs[0], binArgs[1:]...)
	if runtime.Compiler == "gccgolang" {
		cmd.Env = append(cmd.Environ(), "GCCGO=1")
	}
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Logf("%s", out)
		t.Fatal(err)
	}

	checkLineComments(t, libgolangh)
}

var badLineRegexp = regexp.MustCompile(`(?m)^#line [0-9]+ "/.*$`)

// checkLineComments checks that the export header generated by
// -buildmode=c-archive doesn't have any absolute paths in the #line
// comments. We don't want those paths because they are unhelpful for
// the user and make the files change based on details of the location
// of GOPATH.
func checkLineComments(t *testing.T, hdrname string) {
	hdr, err := os.ReadFile(hdrname)
	if err != nil {
		if !os.IsNotExist(err) {
			t.Error(err)
		}
		return
	}
	if line := badLineRegexp.Find(hdr); line != nil {
		t.Errorf("bad #line directive with absolute path in %s: %q", hdrname, line)
	}
}

// checkArchive verifies that the created library looks OK.
// We just check a couple of things now, we can add more checks as needed.
func checkArchive(t *testing.T, arname string) {
	t.Helper()

	switch GOOS {
	case "aix", "darwin", "ios", "windows":
		// We don't have any checks for non-ELF libraries yet.
		if _, err := os.Stat(arname); err != nil {
			t.Errorf("archive %s does not exist: %v", arname, err)
		}
	default:
		checkELFArchive(t, arname)
	}
}

// checkELFArchive checks an ELF archive.
func checkELFArchive(t *testing.T, arname string) {
	t.Helper()

	f, err := os.Open(arname)
	if err != nil {
		t.Errorf("archive %s does not exist: %v", arname, err)
		return
	}
	defer f.Close()

	// TODO(iant): put these in a shared package?  But where?
	const (
		magic = "!<arch>\n"
		fmag  = "`\n"

		namelen = 16
		datelen = 12
		uidlen  = 6
		gidlen  = 6
		modelen = 8
		sizelen = 10
		fmaglen = 2
		hdrlen  = namelen + datelen + uidlen + gidlen + modelen + sizelen + fmaglen
	)

	type arhdr struct {
		name string
		date string
		uid  string
		gid  string
		mode string
		size string
		fmag string
	}

	var magbuf [len(magic)]byte
	if _, err := io.ReadFull(f, magbuf[:]); err != nil {
		t.Errorf("%s: archive too short", arname)
		return
	}
	if string(magbuf[:]) != magic {
		t.Errorf("%s: incorrect archive magic string %q", arname, magbuf)
	}

	off := int64(len(magic))
	for {
		if off&1 != 0 {
			var b [1]byte
			if _, err := f.Read(b[:]); err != nil {
				if err == io.EOF {
					break
				}
				t.Errorf("%s: error skipping alignment byte at %d: %v", arname, off, err)
			}
			off++
		}

		var hdrbuf [hdrlen]byte
		if _, err := io.ReadFull(f, hdrbuf[:]); err != nil {
			if err == io.EOF {
				break
			}
			t.Errorf("%s: error reading archive header at %d: %v", arname, off, err)
			return
		}

		var hdr arhdr
		hdrslice := hdrbuf[:]
		set := func(len int, ps *string) {
			*ps = string(bytes.TrimSpace(hdrslice[:len]))
			hdrslice = hdrslice[len:]
		}
		set(namelen, &hdr.name)
		set(datelen, &hdr.date)
		set(uidlen, &hdr.uid)
		set(gidlen, &hdr.gid)
		set(modelen, &hdr.mode)
		set(sizelen, &hdr.size)
		hdr.fmag = string(hdrslice[:fmaglen])
		hdrslice = hdrslice[fmaglen:]
		if len(hdrslice) != 0 {
			t.Fatalf("internal error: len(hdrslice) == %d", len(hdrslice))
		}

		if hdr.fmag != fmag {
			t.Errorf("%s: invalid fmagic value %q at %d", arname, hdr.fmag, off)
			return
		}

		size, err := strconv.ParseInt(hdr.size, 10, 64)
		if err != nil {
			t.Errorf("%s: error parsing size %q at %d: %v", arname, hdr.size, off, err)
			return
		}

		off += hdrlen

		switch hdr.name {
		case "__.SYMDEF", "/", "/SYM64/":
			// The archive symbol map.
		case "//", "ARFILENAMES/":
			// The extended name table.
		default:
			// This should be an ELF object.
			checkELFArchiveObject(t, arname, off, io.NewSectionReader(f, off, size))
		}

		off += size
		if _, err := f.Seek(off, io.SeekStart); err != nil {
			t.Errorf("%s: failed to seek to %d: %v", arname, off, err)
		}
	}
}

// checkELFArchiveObject checks an object in an ELF archive.
func checkELFArchiveObject(t *testing.T, arname string, off int64, obj io.ReaderAt) {
	t.Helper()

	ef, err := elf.NewFile(obj)
	if err != nil {
		t.Errorf("%s: failed to open ELF file at %d: %v", arname, off, err)
		return
	}
	defer ef.Close()

	// Verify section types.
	for _, sec := range ef.Sections {
		want := elf.SHT_NULL
		switch sec.Name {
		case ".text", ".data":
			want = elf.SHT_PROGBITS
		case ".bss":
			want = elf.SHT_NOBITS
		case ".symtab":
			want = elf.SHT_SYMTAB
		case ".strtab":
			want = elf.SHT_STRTAB
		case ".init_array":
			want = elf.SHT_INIT_ARRAY
		case ".fini_array":
			want = elf.SHT_FINI_ARRAY
		case ".preinit_array":
			want = elf.SHT_PREINIT_ARRAY
		}
		if want != elf.SHT_NULL && sec.Type != want {
			t.Errorf("%s: incorrect section type in elf file at %d for section %q: golangt %v want %v", arname, off, sec.Name, sec.Type, want)
		}
	}
}

func TestInstall(t *testing.T) {
	globalSkip(t)
	testenv.MustHaveGoBuild(t)
	testenv.MustHaveCGO(t)
	testenv.MustHaveBuildMode(t, "c-archive")

	if !testWork {
		defer os.RemoveAll(filepath.Join(GOPATH, "pkg"))
	}

	libgolanga := "libgolang.a"
	if runtime.Compiler == "gccgolang" {
		libgolanga = "liblibgolang.a"
	}

	// Generate the p.h header file.
	//
	// 'golang install -i -buildmode=c-archive ./libgolang' would do that too, but that
	// would also attempt to install transitive standard-library dependencies to
	// GOROOT, and we cannot assume that GOROOT is writable. (A non-root user may
	// be running this test in a GOROOT owned by root.)
	genHeader(t, "p.h", "./p")

	testInstall(t, "./testp1"+exeSuffix,
		filepath.Join(libgolangdir, libgolanga),
		filepath.Join(libgolangdir, "libgolang.h"),
		"golang", "install", "-buildmode=c-archive", "./libgolang")

	// Test building libgolang other than installing it.
	// Header files are now present.
	testInstall(t, "./testp2"+exeSuffix, "libgolang.a", "libgolang.h",
		"golang", "build", "-buildmode=c-archive", filepath.Join(".", "libgolang", "libgolang.golang"))

	testInstall(t, "./testp3"+exeSuffix, "libgolang.a", "libgolang.h",
		"golang", "build", "-buildmode=c-archive", "-o", "libgolang.a", "./libgolang")
}

func TestEarlySignalHandler(t *testing.T) {
	switch GOOS {
	case "darwin", "ios":
		switch GOARCH {
		case "arm64":
			t.Skipf("skipping on %s/%s; see https://golanglang.org/issue/13701", GOOS, GOARCH)
		}
	case "windows":
		t.Skip("skipping signal test on Windows")
	}
	globalSkip(t)
	testenv.MustHaveGoBuild(t)
	testenv.MustHaveCGO(t)
	testenv.MustHaveBuildMode(t, "c-archive")

	if !testWork {
		defer func() {
			os.Remove("libgolang2.a")
			os.Remove("libgolang2.h")
			os.Remove("testp" + exeSuffix)
			os.RemoveAll(filepath.Join(GOPATH, "pkg"))
		}()
	}

	cmd := exec.Command(testenv.GoToolPath(t), "build", "-buildmode=c-archive", "-o", "libgolang2.a", "./libgolang2")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Logf("%s", out)
		t.Fatal(err)
	}
	checkLineComments(t, "libgolang2.h")
	checkArchive(t, "libgolang2.a")

	ccArgs := append(cc, "-o", "testp"+exeSuffix, "main2.c", "libgolang2.a")
	if runtime.Compiler == "gccgolang" {
		ccArgs = append(ccArgs, "-lgolang")
	}
	if out, err := exec.Command(ccArgs[0], ccArgs[1:]...).CombinedOutput(); err != nil {
		t.Logf("%s", out)
		t.Fatal(err)
	}

	darwin := "0"
	if runtime.GOOS == "darwin" {
		darwin = "1"
	}
	cmd = exec.Command(bin[0], append(bin[1:], darwin)...)

	if out, err := cmd.CombinedOutput(); err != nil {
		t.Logf("%s", out)
		t.Fatal(err)
	}
}

func TestSignalForwarding(t *testing.T) {
	globalSkip(t)
	checkSignalForwardingTest(t)
	buildSignalForwardingTest(t)

	cmd := exec.Command(bin[0], append(bin[1:], "1")...)

	out, err := cmd.CombinedOutput()
	t.Logf("%v\n%s", cmd.Args, out)
	expectSignal(t, err, syscall.SIGSEGV, 0)

	// SIGPIPE is never forwarded on darwin. See golanglang.org/issue/33384.
	if runtime.GOOS != "darwin" && runtime.GOOS != "ios" {
		// Test SIGPIPE forwarding
		cmd = exec.Command(bin[0], append(bin[1:], "3")...)

		out, err = cmd.CombinedOutput()
		if len(out) > 0 {
			t.Logf("%s", out)
		}
		expectSignal(t, err, syscall.SIGPIPE, 0)
	}
}

func TestSignalForwardingExternal(t *testing.T) {
	if GOOS == "freebsd" || GOOS == "aix" {
		t.Skipf("skipping on %s/%s; signal always golanges to the Go runtime", GOOS, GOARCH)
	} else if GOOS == "darwin" && GOARCH == "amd64" {
		t.Skipf("skipping on %s/%s: runtime does not permit SI_USER SIGSEGV", GOOS, GOARCH)
	}
	globalSkip(t)
	checkSignalForwardingTest(t)
	buildSignalForwardingTest(t)

	// We want to send the process a signal and see if it dies.
	// Normally the signal golanges to the C thread, the Go signal
	// handler picks it up, sees that it is running in a C thread,
	// and the program dies. Unfortunately, occasionally the
	// signal is delivered to a Go thread, which winds up
	// discarding it because it was sent by another program and
	// there is no Go handler for it. To avoid this, run the
	// program several times in the hopes that it will eventually
	// fail.
	const tries = 20
	for i := 0; i < tries; i++ {
		err := runSignalForwardingTest(t, "2")
		if err == nil {
			continue
		}

		// If the signal is delivered to a C thread, as expected,
		// the Go signal handler will disable itself and re-raise
		// the signal, causing the program to die with SIGSEGV.
		//
		// It is also possible that the signal will be
		// delivered to a Go thread, such as a GC thread.
		// Currently when the Go runtime sees that a SIGSEGV was
		// sent from a different program, it first tries to send
		// the signal to the os/signal API. If nothing is looking
		// for (or explicitly ignoring) SIGSEGV, then it crashes.
		// Because the Go runtime is invoked via a c-archive,
		// it treats this as GOTRACEBACK=crash, meaning that it
		// dumps a stack trace for all golangroutines, which it does
		// by raising SIGQUIT. The effect is that we will see the
		// program die with SIGQUIT in that case, not SIGSEGV.
		if expectSignal(t, err, syscall.SIGSEGV, syscall.SIGQUIT) {
			return
		}
	}

	t.Errorf("program succeeded unexpectedly %d times", tries)
}

func TestSignalForwardingGo(t *testing.T) {
	// This test fails on darwin-amd64 because of the special
	// handling of user-generated SIGSEGV signals in fixsigcode in
	// runtime/signal_darwin_amd64.golang.
	if runtime.GOOS == "darwin" && runtime.GOARCH == "amd64" {
		t.Skip("not supported on darwin-amd64")
	}
	globalSkip(t)

	checkSignalForwardingTest(t)
	buildSignalForwardingTest(t)
	err := runSignalForwardingTest(t, "4")

	// Occasionally the signal will be delivered to a C thread,
	// and the program will crash with SIGSEGV.
	expectSignal(t, err, syscall.SIGQUIT, syscall.SIGSEGV)
}

// checkSignalForwardingTest calls t.Skip if the SignalForwarding test
// doesn't work on this platform.
func checkSignalForwardingTest(t *testing.T) {
	switch GOOS {
	case "darwin", "ios":
		switch GOARCH {
		case "arm64":
			t.Skipf("skipping on %s/%s; see https://golanglang.org/issue/13701", GOOS, GOARCH)
		}
	case "windows":
		t.Skip("skipping signal test on Windows")
	}
	testenv.MustHaveGoBuild(t)
	testenv.MustHaveCGO(t)
	testenv.MustHaveBuildMode(t, "c-archive")
}

// buildSignalForwardingTest builds the executable used by the various
// signal forwarding tests.
func buildSignalForwardingTest(t *testing.T) {
	if !testWork {
		t.Cleanup(func() {
			os.Remove("libgolang2.a")
			os.Remove("libgolang2.h")
			os.Remove("testp" + exeSuffix)
			os.RemoveAll(filepath.Join(GOPATH, "pkg"))
		})
	}

	t.Log("golang build -buildmode=c-archive -o libgolang2.a ./libgolang2")
	cmd := exec.Command(testenv.GoToolPath(t), "build", "-buildmode=c-archive", "-o", "libgolang2.a", "./libgolang2")
	out, err := cmd.CombinedOutput()
	if len(out) > 0 {
		t.Logf("%s", out)
	}
	if err != nil {
		t.Fatal(err)
	}

	checkLineComments(t, "libgolang2.h")
	checkArchive(t, "libgolang2.a")

	ccArgs := append(cc, "-o", "testp"+exeSuffix, "main5.c", "libgolang2.a")
	if runtime.Compiler == "gccgolang" {
		ccArgs = append(ccArgs, "-lgolang")
	}
	t.Log(ccArgs)
	out, err = exec.Command(ccArgs[0], ccArgs[1:]...).CombinedOutput()
	if len(out) > 0 {
		t.Logf("%s", out)
	}
	if err != nil {
		t.Fatal(err)
	}
}

func runSignalForwardingTest(t *testing.T, arg string) error {
	t.Logf("%v %s", bin, arg)
	cmd := exec.Command(bin[0], append(bin[1:], arg)...)

	var out strings.Builder
	cmd.Stdout = &out

	stderr, err := cmd.StderrPipe()
	if err != nil {
		t.Fatal(err)
	}
	defer stderr.Close()

	r := bufio.NewReader(stderr)

	err = cmd.Start()
	if err != nil {
		t.Fatal(err)
	}

	// Wait for trigger to ensure that process is started.
	ok, err := r.ReadString('\n')

	// Verify trigger.
	if err != nil || ok != "OK\n" {
		t.Fatal("Did not receive OK signal")
	}

	var wg sync.WaitGroup
	wg.Add(1)
	var errsb strings.Builder
	golang func() {
		defer wg.Done()
		io.Copy(&errsb, r)
	}()

	// Give the program a chance to enter the function.
	// If the program doesn't get there the test will still
	// pass, although it doesn't quite test what we intended.
	// This is fine as long as the program normally makes it.
	time.Sleep(time.Millisecond)

	cmd.Process.Signal(syscall.SIGSEGV)

	err = cmd.Wait()

	s := out.String()
	if len(s) > 0 {
		t.Log(s)
	}
	wg.Wait()
	s = errsb.String()
	if len(s) > 0 {
		t.Log(s)
	}

	return err
}

// expectSignal checks that err, the exit status of a test program,
// shows a failure due to a specific signal or two. Returns whether we
// found an expected signal.
func expectSignal(t *testing.T, err error, sig1, sig2 syscall.Signal) bool {
	t.Helper()
	if err == nil {
		t.Error("test program succeeded unexpectedly")
	} else if ee, ok := err.(*exec.ExitError); !ok {
		t.Errorf("error (%v) has type %T; expected exec.ExitError", err, err)
	} else if ws, ok := ee.Sys().(syscall.WaitStatus); !ok {
		t.Errorf("error.Sys (%v) has type %T; expected syscall.WaitStatus", ee.Sys(), ee.Sys())
	} else if !ws.Signaled() || (ws.Signal() != sig1 && ws.Signal() != sig2) {
		if sig2 == 0 {
			t.Errorf("golangt %q; expected signal %q", ee, sig1)
		} else {
			t.Errorf("golangt %q; expected signal %q or %q", ee, sig1, sig2)
		}
	} else {
		return true
	}
	return false
}

func TestOsSignal(t *testing.T) {
	switch GOOS {
	case "windows":
		t.Skip("skipping signal test on Windows")
	}
	globalSkip(t)
	testenv.MustHaveGoBuild(t)
	testenv.MustHaveCGO(t)
	testenv.MustHaveBuildMode(t, "c-archive")

	if !testWork {
		defer func() {
			os.Remove("libgolang3.a")
			os.Remove("libgolang3.h")
			os.Remove("testp" + exeSuffix)
			os.RemoveAll(filepath.Join(GOPATH, "pkg"))
		}()
	}

	cmd := exec.Command(testenv.GoToolPath(t), "build", "-buildmode=c-archive", "-o", "libgolang3.a", "./libgolang3")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Logf("%s", out)
		t.Fatal(err)
	}
	checkLineComments(t, "libgolang3.h")
	checkArchive(t, "libgolang3.a")

	ccArgs := append(cc, "-o", "testp"+exeSuffix, "main3.c", "libgolang3.a")
	if runtime.Compiler == "gccgolang" {
		ccArgs = append(ccArgs, "-lgolang")
	}
	if out, err := exec.Command(ccArgs[0], ccArgs[1:]...).CombinedOutput(); err != nil {
		t.Logf("%s", out)
		t.Fatal(err)
	}

	if out, err := exec.Command(bin[0], bin[1:]...).CombinedOutput(); err != nil {
		t.Logf("%s", out)
		t.Fatal(err)
	}
}

func TestSigaltstack(t *testing.T) {
	switch GOOS {
	case "windows":
		t.Skip("skipping signal test on Windows")
	}
	globalSkip(t)
	testenv.MustHaveGoBuild(t)
	testenv.MustHaveCGO(t)
	testenv.MustHaveBuildMode(t, "c-archive")

	if !testWork {
		defer func() {
			os.Remove("libgolang4.a")
			os.Remove("libgolang4.h")
			os.Remove("testp" + exeSuffix)
			os.RemoveAll(filepath.Join(GOPATH, "pkg"))
		}()
	}

	cmd := exec.Command(testenv.GoToolPath(t), "build", "-buildmode=c-archive", "-o", "libgolang4.a", "./libgolang4")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Logf("%s", out)
		t.Fatal(err)
	}
	checkLineComments(t, "libgolang4.h")
	checkArchive(t, "libgolang4.a")

	ccArgs := append(cc, "-o", "testp"+exeSuffix, "main4.c", "libgolang4.a")
	if runtime.Compiler == "gccgolang" {
		ccArgs = append(ccArgs, "-lgolang")
	}
	if out, err := exec.Command(ccArgs[0], ccArgs[1:]...).CombinedOutput(); err != nil {
		t.Logf("%s", out)
		t.Fatal(err)
	}

	if out, err := exec.Command(bin[0], bin[1:]...).CombinedOutput(); err != nil {
		t.Logf("%s", out)
		t.Fatal(err)
	}
}

const testar = `#!/usr/bin/env bash
while [[ $1 == -* ]] >/dev/null; do
  shift
done
echo "testar" > $1
echo "testar" > PWD/testar.ran
`

func TestExtar(t *testing.T) {
	switch GOOS {
	case "windows":
		t.Skip("skipping signal test on Windows")
	}
	if runtime.Compiler == "gccgolang" {
		t.Skip("skipping -extar test when using gccgolang")
	}
	globalSkip(t)
	testenv.MustHaveGoBuild(t)
	testenv.MustHaveCGO(t)
	testenv.MustHaveBuildMode(t, "c-archive")
	testenv.MustHaveExecPath(t, "bash") // This test uses a bash script

	if !testWork {
		defer func() {
			os.Remove("libgolang4.a")
			os.Remove("libgolang4.h")
			os.Remove("testar")
			os.Remove("testar.ran")
			os.RemoveAll(filepath.Join(GOPATH, "pkg"))
		}()
	}

	os.Remove("testar")
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	s := strings.Replace(testar, "PWD", dir, 1)
	if err := os.WriteFile("testar", []byte(s), 0777); err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(testenv.GoToolPath(t), "build", "-buildmode=c-archive", "-ldflags=-extar="+filepath.Join(dir, "testar"), "-o", "libgolang4.a", "./libgolang4")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Logf("%s", out)
		t.Fatal(err)
	}
	checkLineComments(t, "libgolang4.h")

	if _, err := os.Stat("testar.ran"); err != nil {
		if os.IsNotExist(err) {
			t.Error("testar does not exist after golang build")
		} else {
			t.Errorf("error checking testar: %v", err)
		}
	}
}

func TestPIE(t *testing.T) {
	switch GOOS {
	case "windows", "darwin", "ios", "plan9":
		t.Skipf("skipping PIE test on %s", GOOS)
	}
	globalSkip(t)
	testenv.MustHaveGoBuild(t)
	testenv.MustHaveCGO(t)
	testenv.MustHaveBuildMode(t, "c-archive")

	libgolanga := "libgolang.a"
	if runtime.Compiler == "gccgolang" {
		libgolanga = "liblibgolang.a"
	}

	if !testWork {
		defer func() {
			os.Remove("testp" + exeSuffix)
			os.Remove(libgolanga)
			os.RemoveAll(filepath.Join(GOPATH, "pkg"))
		}()
	}

	// Generate the p.h header file.
	//
	// 'golang install -i -buildmode=c-archive ./libgolang' would do that too, but that
	// would also attempt to install transitive standard-library dependencies to
	// GOROOT, and we cannot assume that GOROOT is writable. (A non-root user may
	// be running this test in a GOROOT owned by root.)
	genHeader(t, "p.h", "./p")

	cmd := exec.Command(testenv.GoToolPath(t), "build", "-buildmode=c-archive", "./libgolang")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Logf("%s", out)
		t.Fatal(err)
	}

	ccArgs := append(cc, "-fPIE", "-pie", "-o", "testp"+exeSuffix, "main.c", "main_unix.c", libgolanga)
	if runtime.Compiler == "gccgolang" {
		ccArgs = append(ccArgs, "-lgolang")
	}
	if out, err := exec.Command(ccArgs[0], ccArgs[1:]...).CombinedOutput(); err != nil {
		t.Logf("%s", out)
		t.Fatal(err)
	}

	binArgs := append(bin, "arg1", "arg2")
	cmd = exec.Command(binArgs[0], binArgs[1:]...)
	if runtime.Compiler == "gccgolang" {
		cmd.Env = append(os.Environ(), "GCCGO=1")
	}
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Logf("%s", out)
		t.Fatal(err)
	}

	if GOOS != "aix" {
		f, err := elf.Open("testp" + exeSuffix)
		if err != nil {
			t.Fatal("elf.Open failed: ", err)
		}
		defer f.Close()
		if hasDynTag(t, f, elf.DT_TEXTREL) {
			t.Errorf("%s has DT_TEXTREL flag", "testp"+exeSuffix)
		}
	}
}

func hasDynTag(t *testing.T, f *elf.File, tag elf.DynTag) bool {
	ds := f.SectionByType(elf.SHT_DYNAMIC)
	if ds == nil {
		t.Error("no SHT_DYNAMIC section")
		return false
	}
	d, err := ds.Data()
	if err != nil {
		t.Errorf("can't read SHT_DYNAMIC contents: %v", err)
		return false
	}
	for len(d) > 0 {
		var t elf.DynTag
		switch f.Class {
		case elf.ELFCLASS32:
			t = elf.DynTag(f.ByteOrder.Uint32(d[:4]))
			d = d[8:]
		case elf.ELFCLASS64:
			t = elf.DynTag(f.ByteOrder.Uint64(d[:8]))
			d = d[16:]
		}
		if t == tag {
			return true
		}
	}
	return false
}

func TestSIGPROF(t *testing.T) {
	switch GOOS {
	case "windows", "plan9":
		t.Skipf("skipping SIGPROF test on %s", GOOS)
	case "darwin", "ios":
		t.Skipf("skipping SIGPROF test on %s; see https://golanglang.org/issue/19320", GOOS)
	}
	globalSkip(t)
	testenv.MustHaveGoBuild(t)
	testenv.MustHaveCGO(t)
	testenv.MustHaveBuildMode(t, "c-archive")

	t.Parallel()

	if !testWork {
		defer func() {
			os.Remove("testp6" + exeSuffix)
			os.Remove("libgolang6.a")
			os.Remove("libgolang6.h")
		}()
	}

	cmd := exec.Command(testenv.GoToolPath(t), "build", "-buildmode=c-archive", "-o", "libgolang6.a", "./libgolang6")
	out, err := cmd.CombinedOutput()
	t.Logf("%v\n%s", cmd.Args, out)
	if err != nil {
		t.Fatal(err)
	}
	checkLineComments(t, "libgolang6.h")
	checkArchive(t, "libgolang6.a")

	ccArgs := append(cc, "-o", "testp6"+exeSuffix, "main6.c", "libgolang6.a")
	if runtime.Compiler == "gccgolang" {
		ccArgs = append(ccArgs, "-lgolang")
	}
	out, err = exec.Command(ccArgs[0], ccArgs[1:]...).CombinedOutput()
	t.Logf("%v\n%s", ccArgs, out)
	if err != nil {
		t.Fatal(err)
	}

	argv := cmdToRun("./testp6")
	cmd = exec.Command(argv[0], argv[1:]...)
	out, err = cmd.CombinedOutput()
	t.Logf("%v\n%s", argv, out)
	if err != nil {
		t.Fatal(err)
	}
}

// TestCompileWithoutShared tests that if we compile code without the
// -shared option, we can put it into an archive. When we use the golang
// tool with -buildmode=c-archive, it passes -shared to the compiler,
// so we override that. The golang tool doesn't work this way, but Bazel
// will likely do it in the future. And it ought to work. This test
// was added because at one time it did not work on PPC Linux.
func TestCompileWithoutShared(t *testing.T) {
	globalSkip(t)
	// For simplicity, reuse the signal forwarding test.
	checkSignalForwardingTest(t)
	testenv.MustHaveGoBuild(t)

	if !testWork {
		defer func() {
			os.Remove("libgolang2.a")
			os.Remove("libgolang2.h")
		}()
	}

	cmd := exec.Command(testenv.GoToolPath(t), "build", "-buildmode=c-archive", "-gcflags=-shared=false", "-o", "libgolang2.a", "./libgolang2")
	out, err := cmd.CombinedOutput()
	t.Logf("%v\n%s", cmd.Args, out)
	if err != nil {
		t.Fatal(err)
	}
	checkLineComments(t, "libgolang2.h")
	checkArchive(t, "libgolang2.a")

	exe := "./testnoshared" + exeSuffix

	// In some cases, -no-pie is needed here, but not accepted everywhere. First try
	// if -no-pie is accepted. See #22126.
	ccArgs := append(cc, "-o", exe, "-no-pie", "main5.c", "libgolang2.a")
	if runtime.Compiler == "gccgolang" {
		ccArgs = append(ccArgs, "-lgolang")
	}
	out, err = exec.Command(ccArgs[0], ccArgs[1:]...).CombinedOutput()
	t.Logf("%v\n%s", ccArgs, out)

	// If -no-pie unrecognized, try -nopie if this is possibly clang
	if err != nil && bytes.Contains(out, []byte("unknown")) && !strings.Contains(cc[0], "gcc") {
		ccArgs = append(cc, "-o", exe, "-nopie", "main5.c", "libgolang2.a")
		out, err = exec.Command(ccArgs[0], ccArgs[1:]...).CombinedOutput()
		t.Logf("%v\n%s", ccArgs, out)
	}

	// Don't use either -no-pie or -nopie
	if err != nil && bytes.Contains(out, []byte("unrecognized")) {
		ccArgs = append(cc, "-o", exe, "main5.c", "libgolang2.a")
		out, err = exec.Command(ccArgs[0], ccArgs[1:]...).CombinedOutput()
		t.Logf("%v\n%s", ccArgs, out)
	}
	if err != nil {
		t.Fatal(err)
	}
	if !testWork {
		defer os.Remove(exe)
	}

	binArgs := append(cmdToRun(exe), "1")
	out, err = exec.Command(binArgs[0], binArgs[1:]...).CombinedOutput()
	t.Logf("%v\n%s", binArgs, out)
	expectSignal(t, err, syscall.SIGSEGV, 0)

	// SIGPIPE is never forwarded on darwin. See golanglang.org/issue/33384.
	if runtime.GOOS != "darwin" && runtime.GOOS != "ios" {
		binArgs := append(cmdToRun(exe), "3")
		out, err = exec.Command(binArgs[0], binArgs[1:]...).CombinedOutput()
		t.Logf("%v\n%s", binArgs, out)
		expectSignal(t, err, syscall.SIGPIPE, 0)
	}
}

// Test that installing a second time recreates the header file.
func TestCachedInstall(t *testing.T) {
	globalSkip(t)
	testenv.MustHaveGoBuild(t)
	testenv.MustHaveCGO(t)
	testenv.MustHaveBuildMode(t, "c-archive")

	if !testWork {
		defer os.RemoveAll(filepath.Join(GOPATH, "pkg"))
	}

	h := filepath.Join(libgolangdir, "libgolang.h")

	buildcmd := []string{"golang", "install", "-buildmode=c-archive", "./libgolang"}

	cmd := exec.Command(buildcmd[0], buildcmd[1:]...)
	cmd.Env = append(cmd.Environ(), "GO111MODULE=off") // 'golang install' only works in GOPATH mode
	t.Log(buildcmd)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Logf("%s", out)
		t.Fatal(err)
	}

	if _, err := os.Stat(h); err != nil {
		t.Errorf("libgolang.h not installed: %v", err)
	}

	if err := os.Remove(h); err != nil {
		t.Fatal(err)
	}

	cmd = exec.Command(buildcmd[0], buildcmd[1:]...)
	cmd.Env = append(cmd.Environ(), "GO111MODULE=off")
	t.Log(buildcmd)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Logf("%s", out)
		t.Fatal(err)
	}

	if _, err := os.Stat(h); err != nil {
		t.Errorf("libgolang.h not installed in second run: %v", err)
	}
}

// Issue 35294.
func TestManyCalls(t *testing.T) {
	globalSkip(t)
	testenv.MustHaveGoBuild(t)
	testenv.MustHaveCGO(t)
	testenv.MustHaveBuildMode(t, "c-archive")

	t.Parallel()

	if !testWork {
		defer func() {
			os.Remove("testp7" + exeSuffix)
			os.Remove("libgolang7.a")
			os.Remove("libgolang7.h")
		}()
	}

	cmd := exec.Command(testenv.GoToolPath(t), "build", "-buildmode=c-archive", "-o", "libgolang7.a", "./libgolang7")
	out, err := cmd.CombinedOutput()
	t.Logf("%v\n%s", cmd.Args, out)
	if err != nil {
		t.Fatal(err)
	}
	checkLineComments(t, "libgolang7.h")
	checkArchive(t, "libgolang7.a")

	ccArgs := append(cc, "-o", "testp7"+exeSuffix, "main7.c", "libgolang7.a")
	if runtime.Compiler == "gccgolang" {
		ccArgs = append(ccArgs, "-lgolang")
	}
	out, err = exec.Command(ccArgs[0], ccArgs[1:]...).CombinedOutput()
	t.Logf("%v\n%s", ccArgs, out)
	if err != nil {
		t.Fatal(err)
	}

	argv := cmdToRun("./testp7")
	cmd = testenv.Command(t, argv[0], argv[1:]...)
	sb := new(strings.Builder)
	cmd.Stdout = sb
	cmd.Stderr = sb
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}

	err = cmd.Wait()
	t.Logf("%v\n%s", cmd.Args, sb)
	if err != nil {
		t.Error(err)
	}
}

// Issue 49288.
func TestPreemption(t *testing.T) {
	if runtime.Compiler == "gccgolang" {
		t.Skip("skipping asynchronous preemption test with gccgolang")
	}
	globalSkip(t)
	testenv.MustHaveGoBuild(t)
	testenv.MustHaveCGO(t)
	testenv.MustHaveBuildMode(t, "c-archive")

	t.Parallel()

	if !testWork {
		defer func() {
			os.Remove("testp8" + exeSuffix)
			os.Remove("libgolang8.a")
			os.Remove("libgolang8.h")
		}()
	}

	cmd := exec.Command(testenv.GoToolPath(t), "build", "-buildmode=c-archive", "-o", "libgolang8.a", "./libgolang8")
	out, err := cmd.CombinedOutput()
	t.Logf("%v\n%s", cmd.Args, out)
	if err != nil {
		t.Fatal(err)
	}
	checkLineComments(t, "libgolang8.h")
	checkArchive(t, "libgolang8.a")

	ccArgs := append(cc, "-o", "testp8"+exeSuffix, "main8.c", "libgolang8.a")
	out, err = exec.Command(ccArgs[0], ccArgs[1:]...).CombinedOutput()
	t.Logf("%v\n%s", ccArgs, out)
	if err != nil {
		t.Fatal(err)
	}

	argv := cmdToRun("./testp8")
	cmd = testenv.Command(t, argv[0], argv[1:]...)
	sb := new(strings.Builder)
	cmd.Stdout = sb
	cmd.Stderr = sb
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}

	err = cmd.Wait()
	t.Logf("%v\n%s", cmd.Args, sb)
	if err != nil {
		t.Error(err)
	}
}

// Issue 59294 and 68285. Test calling Go function from C after with
// various stack space.
func TestDeepStack(t *testing.T) {
	globalSkip(t)
	testenv.MustHaveGoBuild(t)
	testenv.MustHaveCGO(t)
	testenv.MustHaveBuildMode(t, "c-archive")

	t.Parallel()

	if !testWork {
		defer func() {
			os.Remove("testp9" + exeSuffix)
			os.Remove("libgolang9.a")
			os.Remove("libgolang9.h")
		}()
	}

	cmd := exec.Command(testenv.GoToolPath(t), "build", "-buildmode=c-archive", "-o", "libgolang9.a", "./libgolang9")
	out, err := cmd.CombinedOutput()
	t.Logf("%v\n%s", cmd.Args, out)
	if err != nil {
		t.Fatal(err)
	}
	checkLineComments(t, "libgolang9.h")
	checkArchive(t, "libgolang9.a")

	// build with -O0 so the C compiler won't optimize out the large stack frame
	ccArgs := append(cc, "-O0", "-o", "testp9"+exeSuffix, "main9.c", "libgolang9.a")
	out, err = exec.Command(ccArgs[0], ccArgs[1:]...).CombinedOutput()
	t.Logf("%v\n%s", ccArgs, out)
	if err != nil {
		t.Fatal(err)
	}

	argv := cmdToRun("./testp9")
	cmd = exec.Command(argv[0], argv[1:]...)
	sb := new(strings.Builder)
	cmd.Stdout = sb
	cmd.Stderr = sb
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}

	timer := time.AfterFunc(time.Minute,
		func() {
			t.Error("test program timed out")
			cmd.Process.Kill()
		},
	)
	defer timer.Stop()

	err = cmd.Wait()
	t.Logf("%v\n%s", cmd.Args, sb)
	if err != nil {
		t.Error(err)
	}
}

func BenchmarkCgolangCallbackMainThread(b *testing.B) {
	// Benchmark for calling into Go fron C main thread.
	// See issue #68587.
	//
	// It uses a subprocess, which is a C binary that calls
	// Go on the main thread b.N times. There is some overhead
	// for launching the subprocess. It is probably fine when
	// b.N is large.

	globalSkip(b)
	testenv.MustHaveGoBuild(b)
	testenv.MustHaveCGO(b)
	testenv.MustHaveBuildMode(b, "c-archive")

	if !testWork {
		defer func() {
			os.Remove("testp10" + exeSuffix)
			os.Remove("libgolang10.a")
			os.Remove("libgolang10.h")
		}()
	}

	cmd := exec.Command(testenv.GoToolPath(b), "build", "-buildmode=c-archive", "-o", "libgolang10.a", "./libgolang10")
	out, err := cmd.CombinedOutput()
	b.Logf("%v\n%s", cmd.Args, out)
	if err != nil {
		b.Fatal(err)
	}

	ccArgs := append(cc, "-o", "testp10"+exeSuffix, "main10.c", "libgolang10.a")
	out, err = exec.Command(ccArgs[0], ccArgs[1:]...).CombinedOutput()
	b.Logf("%v\n%s", ccArgs, out)
	if err != nil {
		b.Fatal(err)
	}

	argv := cmdToRun("./testp10")
	argv = append(argv, fmt.Sprint(b.N))
	cmd = exec.Command(argv[0], argv[1:]...)

	b.ResetTimer()
	err = cmd.Run()
	if err != nil {
		b.Fatal(err)
	}
}

func TestSharedObject(t *testing.T) {
	// Test that we can put a Go c-archive into a C shared object.
	globalSkip(t)
	testenv.MustHaveGoBuild(t)
	testenv.MustHaveCGO(t)
	testenv.MustHaveBuildMode(t, "c-archive")

	t.Parallel()

	if !testWork {
		defer func() {
			os.Remove("libgolang_s.a")
			os.Remove("libgolang_s.h")
			os.Remove("libgolang_s.so")
		}()
	}

	cmd := exec.Command(testenv.GoToolPath(t), "build", "-buildmode=c-archive", "-o", "libgolang_s.a", "./libgolang")
	out, err := cmd.CombinedOutput()
	t.Logf("%v\n%s", cmd.Args, out)
	if err != nil {
		t.Fatal(err)
	}

	ccArgs := append(cc, "-shared", "-o", "libgolang_s.so", "libgolang_s.a")
	out, err = exec.Command(ccArgs[0], ccArgs[1:]...).CombinedOutput()
	t.Logf("%v\n%s", ccArgs, out)
	if err != nil {
		t.Fatal(err)
	}
}
