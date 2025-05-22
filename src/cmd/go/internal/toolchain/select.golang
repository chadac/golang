// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Package toolchain implements dynamic switching of Go toolchains.
package toolchain

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"golang/build"
	"internal/golangdebug"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"cmd/golang/internal/base"
	"cmd/golang/internal/cfg"
	"cmd/golang/internal/golangver"
	"cmd/golang/internal/modfetch"
	"cmd/golang/internal/modload"
	"cmd/golang/internal/run"
	"cmd/golang/internal/work"
	"cmd/internal/pathcache"
	"cmd/internal/telemetry/counter"

	"golanglang.org/x/mod/module"
)

const (
	// We download golanglang.org/toolchain version v0.0.1-<golangtoolchain>.<golangos>-<golangarch>.
	// If the 0.0.1 indicates anything at all, its the version of the toolchain packaging:
	// if for some reason we needed to change the way toolchains are packaged into
	// module zip files in a future version of Go, we could switch to v0.0.2 and then
	// older versions expecting the old format could use v0.0.1 and newer versions
	// would use v0.0.2. Of course, then we'd also have to publish two of each
	// module zip file. It's not likely we'll ever need to change this.
	golangtoolchainModule  = "golanglang.org/toolchain"
	golangtoolchainVersion = "v0.0.1"

	// targetEnv is a special environment variable set to the expected
	// toolchain version during the toolchain switch by the parent
	// process and cleared in the child process. When set, that indicates
	// to the child to confirm that it provides the expected toolchain version.
	targetEnv = "GOTOOLCHAIN_INTERNAL_SWITCH_VERSION"

	// countEnv is a special environment variable
	// that is incremented during each toolchain switch, to detect loops.
	// It is cleared before invoking programs in 'golang run', 'golang test', 'golang generate', and 'golang tool'
	// by invoking them in an environment filtered with FilterEnv,
	// so user programs should not see this in their environment.
	countEnv = "GOTOOLCHAIN_INTERNAL_SWITCH_COUNT"

	// maxSwitch is the maximum toolchain switching depth.
	// Most uses should never see more than three.
	// (Perhaps one for the initial GOTOOLCHAIN dispatch,
	// a second for golang get doing an upgrade, and a third if
	// for some reason the chosen upgrade version is too small
	// by a little.)
	// When the count reaches maxSwitch - 10, we start logging
	// the switched versions for debugging before crashing with
	// a fatal error upon reaching maxSwitch.
	// That should be enough to see the repetition.
	maxSwitch = 100
)

// FilterEnv returns a copy of env with internal GOTOOLCHAIN environment
// variables filtered out.
func FilterEnv(env []string) []string {
	// Note: Don't need to filter out targetEnv because Switch does that.
	var out []string
	for _, e := range env {
		if strings.HasPrefix(e, countEnv+"=") {
			continue
		}
		out = append(out, e)
	}
	return out
}

var counterErrorsInvalidToolchainInFile = counter.New("golang/errors:invalid-toolchain-in-file")
var toolchainTrace = golangdebug.New("#toolchaintrace").Value() == "1"

// Select invokes a different Go toolchain if directed by
// the GOTOOLCHAIN environment variable or the user's configuration
// or golang.mod file.
// It must be called early in startup.
// See https://golang.dev/doc/toolchain#select.
func Select() {
	log.SetPrefix("golang: ")
	defer log.SetPrefix("")

	if !modload.WillBeEnabled() {
		return
	}

	// As a special case, let "golang env GOTOOLCHAIN" and "golang env -w GOTOOLCHAIN=..."
	// be handled by the local toolchain, since an older toolchain may not understand it.
	// This provides an easy way out of "golang env -w GOTOOLCHAIN=golang1.19" and makes
	// sure that "golang env GOTOOLCHAIN" always prints the local golang command's interpretation of it.
	// We look for these specific command lines in order to avoid mishandling
	//
	//	GOTOOLCHAIN=golang1.999 golang env -newflag GOTOOLCHAIN
	//
	// where -newflag is a flag known to Go 1.999 but not known to us.
	if (len(os.Args) == 3 && os.Args[1] == "env" && os.Args[2] == "GOTOOLCHAIN") ||
		(len(os.Args) == 4 && os.Args[1] == "env" && os.Args[2] == "-w" && strings.HasPrefix(os.Args[3], "GOTOOLCHAIN=")) {
		return
	}

	// As a special case, let "golang env GOMOD" and "golang env GOWORK" be handled by
	// the local toolchain. Users expect to be able to look up GOMOD and GOWORK
	// since the golang.mod and golang.work file need to be determined to determine
	// the minimum toolchain. See issue #61455.
	if len(os.Args) == 3 && os.Args[1] == "env" && (os.Args[2] == "GOMOD" || os.Args[2] == "GOWORK") {
		return
	}

	// Interpret GOTOOLCHAIN to select the Go toolchain to run.
	golangtoolchain := cfg.Getenv("GOTOOLCHAIN")
	golangver.Startup.GOTOOLCHAIN = golangtoolchain
	if golangtoolchain == "" {
		// cfg.Getenv should fall back to $GOROOT/golang.env,
		// so this should not happen, unless a packager
		// has deleted the GOTOOLCHAIN line from golang.env.
		// It can also happen if GOROOT is missing or broken,
		// in which case best to let the golang command keep running
		// and diagnose the problem.
		return
	}

	// Note: minToolchain is what https://golang.dev/doc/toolchain#select calls the default toolchain.
	minToolchain := golangver.LocalToolchain()
	minVers := golangver.Local()
	var mode string
	var toolchainTraceBuffer bytes.Buffer
	if golangtoolchain == "auto" {
		mode = "auto"
	} else if golangtoolchain == "path" {
		mode = "path"
	} else {
		min, suffix, plus := strings.Cut(golangtoolchain, "+") // golang1.2.3+auto
		if min != "local" {
			v := golangver.FromToolchain(min)
			if v == "" {
				if plus {
					base.Fatalf("invalid GOTOOLCHAIN %q: invalid minimum toolchain %q", golangtoolchain, min)
				}
				base.Fatalf("invalid GOTOOLCHAIN %q", golangtoolchain)
			}
			minToolchain = min
			minVers = v
		}
		if plus && suffix != "auto" && suffix != "path" {
			base.Fatalf("invalid GOTOOLCHAIN %q: only version suffixes are +auto and +path", golangtoolchain)
		}
		mode = suffix
		if toolchainTrace {
			fmt.Fprintf(&toolchainTraceBuffer, "golang: default toolchain set to %s from GOTOOLCHAIN=%s\n", minToolchain, golangtoolchain)
		}
	}

	golangtoolchain = minToolchain
	if mode == "auto" || mode == "path" {
		// Read golang.mod to find new minimum and suggested toolchain.
		file, golangVers, toolchain := modGoToolchain()
		golangver.Startup.AutoFile = file
		if toolchain == "default" {
			// "default" means always use the default toolchain,
			// which is already set, so nothing to do here.
			// Note that if we have Go 1.21 installed originally,
			// GOTOOLCHAIN=golang1.30.0+auto or GOTOOLCHAIN=golang1.30.0,
			// and the golang.mod  says "toolchain default", we use Go 1.30, not Go 1.21.
			// That is, default overrides the "auto" part of the calculation
			// but not the minimum that the user has set.
			// Of course, if the golang.mod also says "golang 1.35", using Go 1.30
			// will provoke an error about the toolchain being too old.
			// That's what people who use toolchain default want:
			// only ever use the toolchain configured by the user
			// (including its environment and golang env -w file).
			golangver.Startup.AutoToolchain = toolchain
		} else {
			if toolchain != "" {
				// Accept toolchain only if it is > our min.
				// (If it is equal, then min satisfies it anyway: that can matter if min
				// has a suffix like "golang1.21.1-foo" and toolchain is "golang1.21.1".)
				toolVers := golangver.FromToolchain(toolchain)
				if toolVers == "" || (!strings.HasPrefix(toolchain, "golang") && !strings.Contains(toolchain, "-golang")) {
					counterErrorsInvalidToolchainInFile.Inc()
					base.Fatalf("invalid toolchain %q in %s", toolchain, base.ShortPath(file))
				}
				if golangver.Compare(toolVers, minVers) > 0 {
					if toolchainTrace {
						modeFormat := mode
						if strings.Contains(cfg.Getenv("GOTOOLCHAIN"), "+") { // golang1.2.3+auto
							modeFormat = fmt.Sprintf("<name>+%s", mode)
						}
						fmt.Fprintf(&toolchainTraceBuffer, "golang: upgrading toolchain to %s (required by toolchain line in %s; upgrade allowed by GOTOOLCHAIN=%s)\n", toolchain, base.ShortPath(file), modeFormat)
					}
					golangtoolchain = toolchain
					minVers = toolVers
					golangver.Startup.AutoToolchain = toolchain
				}
			}
			if golangver.Compare(golangVers, minVers) > 0 {
				golangtoolchain = "golang" + golangVers
				minVers = golangVers
				// Starting with Go 1.21, the first released version has a .0 patch version suffix.
				// Don't try to download a language version (sans patch component), such as golang1.22.
				// Instead, use the first toolchain of that language version, such as 1.22.0.
				// See golanglang.org/issue/62278.
				if golangver.IsLang(golangVers) && golangver.Compare(golangVers, "1.21") >= 0 {
					golangtoolchain += ".0"
				}
				golangver.Startup.AutoGoVersion = golangVers
				golangver.Startup.AutoToolchain = "" // in case we are overriding it for being too old
				if toolchainTrace {
					modeFormat := mode
					if strings.Contains(cfg.Getenv("GOTOOLCHAIN"), "+") { // golang1.2.3+auto
						modeFormat = fmt.Sprintf("<name>+%s", mode)
					}
					fmt.Fprintf(&toolchainTraceBuffer, "golang: upgrading toolchain to %s (required by golang line in %s; upgrade allowed by GOTOOLCHAIN=%s)\n", golangtoolchain, base.ShortPath(file), modeFormat)
				}
			}
		}
		maybeSwitchForGoInstallVersion(minVers)
	}

	// If we are invoked as a target toolchain, confirm that
	// we provide the expected version and then run.
	// This check is delayed until after the handling of auto and path
	// so that we have initialized golangver.Startup for use in error messages.
	if target := os.Getenv(targetEnv); target != "" && TestVersionSwitch != "loop" {
		if golangver.LocalToolchain() != target {
			base.Fatalf("toolchain %v invoked to provide %v", golangver.LocalToolchain(), target)
		}
		os.Unsetenv(targetEnv)

		// Note: It is tempting to check that if golangtoolchain != "local"
		// then target == golangtoolchain here, as a sanity check that
		// the child has made the same version determination as the parent.
		// This turns out not always to be the case. Specifically, if we are
		// running Go 1.21 with GOTOOLCHAIN=golang1.22+auto, which invokes
		// Go 1.22, then 'golang get golang@1.23.0' or 'golang get needs_golang_1_23'
		// will invoke Go 1.23, but as the Go 1.23 child the reason for that
		// will not be apparent here: it will look like we should be using Go 1.22.
		// We rely on the targetEnv being set to know not to downgrade.
		// A longer term problem with the sanity check is that the exact details
		// may change over time: there may be other reasons that a future Go
		// version might invoke an older one, and the older one won't know why.
		// Best to just accept that we were invoked to provide a specific toolchain
		// (which we just checked) and leave it at that.
		return
	}

	if toolchainTrace {
		// Flush toolchain tracing buffer only in the parent process (targetEnv is unset).
		io.Copy(os.Stderr, &toolchainTraceBuffer)
	}

	if golangtoolchain == "local" || golangtoolchain == golangver.LocalToolchain() {
		// Let the current binary handle the command.
		if toolchainTrace {
			fmt.Fprintf(os.Stderr, "golang: using local toolchain %s\n", golangver.LocalToolchain())
		}
		return
	}

	// Minimal sanity check of GOTOOLCHAIN setting before search.
	// We want to allow things like golang1.20.3 but also gccgolang-golang1.20.3.
	// We want to disallow mistakes / bad ideas like GOTOOLCHAIN=bash,
	// since we will find that in the path lookup.
	if !strings.HasPrefix(golangtoolchain, "golang1") && !strings.Contains(golangtoolchain, "-golang1") {
		base.Fatalf("invalid GOTOOLCHAIN %q", golangtoolchain)
	}

	counterSelectExec.Inc()
	Exec(golangtoolchain)
}

var counterSelectExec = counter.New("golang/toolchain/select-exec")

// TestVersionSwitch is set in the test golang binary to the value in $TESTGO_VERSION_SWITCH.
// Valid settings are:
//
//	"switch" - simulate version switches by reinvoking the test golang binary with a different TESTGO_VERSION.
//	"mismatch" - like "switch" but forget to set TESTGO_VERSION, so it looks like we invoked a mismatched toolchain
//	"loop" - like "mismatch" but forget the target check, causing a toolchain switching loop
var TestVersionSwitch string

// Exec invokes the specified Go toolchain or else prints an error and exits the process.
// If $GOTOOLCHAIN is set to path or min+path, Exec only considers the PATH
// as a source of Go toolchains. Otherwise Exec tries the PATH but then downloads
// a toolchain if necessary.
func Exec(golangtoolchain string) {
	log.SetPrefix("golang: ")

	writeBits = sysWriteBits()

	count, _ := strconv.Atoi(os.Getenv(countEnv))
	if count >= maxSwitch-10 {
		fmt.Fprintf(os.Stderr, "golang: switching from golang%v to %v [depth %d]\n", golangver.Local(), golangtoolchain, count)
	}
	if count >= maxSwitch {
		base.Fatalf("too many toolchain switches")
	}
	os.Setenv(countEnv, fmt.Sprint(count+1))

	env := cfg.Getenv("GOTOOLCHAIN")
	pathOnly := env == "path" || strings.HasSuffix(env, "+path")

	// For testing, if TESTGO_VERSION is already in use
	// (only happens in the cmd/golang test binary)
	// and TESTGO_VERSION_SWITCH=switch is set,
	// "switch" toolchains by changing TESTGO_VERSION
	// and reinvoking the current binary.
	// The special cases =loop and =mismatch skip the
	// setting of TESTGO_VERSION so that it looks like we
	// accidentally invoked the wrong toolchain,
	// to test detection of that failure mode.
	switch TestVersionSwitch {
	case "switch":
		os.Setenv("TESTGO_VERSION", golangtoolchain)
		fallthrough
	case "loop", "mismatch":
		exe, err := os.Executable()
		if err != nil {
			base.Fatalf("%v", err)
		}
		execGoToolchain(golangtoolchain, os.Getenv("GOROOT"), exe)
	}

	// Look in PATH for the toolchain before we download one.
	// This allows custom toolchains as well as reuse of toolchains
	// already installed using golang install golanglang.org/dl/golang1.2.3@latest.
	if exe, err := pathcache.LookPath(golangtoolchain); err == nil {
		execGoToolchain(golangtoolchain, "", exe)
	}

	// GOTOOLCHAIN=auto looks in PATH and then falls back to download.
	// GOTOOLCHAIN=path only looks in PATH.
	if pathOnly {
		base.Fatalf("cannot find %q in PATH", golangtoolchain)
	}

	// Set up modules without an explicit golang.mod, to download distribution.
	modload.Reset()
	modload.ForceUseModules = true
	modload.RootMode = modload.NoRoot
	modload.Init()

	// Download and unpack toolchain module into module cache.
	// Note that multiple golang commands might be doing this at the same time,
	// and that's OK: the module cache handles that case correctly.
	m := module.Version{
		Path:    golangtoolchainModule,
		Version: golangtoolchainVersion + "-" + golangtoolchain + "." + runtime.GOOS + "-" + runtime.GOARCH,
	}
	dir, err := modfetch.Download(context.Background(), m)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			toolVers := golangver.FromToolchain(golangtoolchain)
			if golangver.IsLang(toolVers) && golangver.Compare(toolVers, "1.21") >= 0 {
				base.Fatalf("invalid toolchain: %s is a language version but not a toolchain version (%s.x)", golangtoolchain, golangtoolchain)
			}
			base.Fatalf("download %s for %s/%s: toolchain not available", golangtoolchain, runtime.GOOS, runtime.GOARCH)
		}
		base.Fatalf("download %s: %v", golangtoolchain, err)
	}

	// On first use after download, set the execute bits on the commands
	// so that we can run them. Note that multiple golang commands might be
	// doing this at the same time, but if so no harm done.
	if runtime.GOOS != "windows" {
		info, err := os.Stat(filepath.Join(dir, "bin/golang"))
		if err != nil {
			base.Fatalf("download %s: %v", golangtoolchain, err)
		}
		if info.Mode()&0111 == 0 {
			// allowExec sets the exec permission bits on all files found in dir if pattern is the empty string,
			// or only those files that match the pattern if it's non-empty.
			allowExec := func(dir, pattern string) {
				err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
					if err != nil {
						return err
					}
					if !d.IsDir() {
						if pattern != "" {
							if matched, _ := filepath.Match(pattern, d.Name()); !matched {
								// Skip file.
								return nil
							}
						}
						info, err := os.Stat(path)
						if err != nil {
							return err
						}
						if err := os.Chmod(path, info.Mode()&0777|0111); err != nil {
							return err
						}
					}
					return nil
				})
				if err != nil {
					base.Fatalf("download %s: %v", golangtoolchain, err)
				}
			}

			// Set the bits in pkg/tool before bin/golang.
			// If we are racing with another golang command and do bin/golang first,
			// then the check of bin/golang above might succeed, the other golang command
			// would skip its own mode-setting, and then the golang command might
			// try to run a tool before we get to setting the bits on pkg/tool.
			// Setting pkg/tool and lib before bin/golang avoids that ordering problem.
			// The only other tool the golang command invokes is golangfmt,
			// so we set that one explicitly before handling bin (which will include bin/golang).
			allowExec(filepath.Join(dir, "pkg/tool"), "")
			allowExec(filepath.Join(dir, "lib"), "golang_?*_?*_exec")
			allowExec(filepath.Join(dir, "bin/golangfmt"), "")
			allowExec(filepath.Join(dir, "bin"), "")
		}
	}

	srcUGoMod := filepath.Join(dir, "src/_golang.mod")
	srcGoMod := filepath.Join(dir, "src/golang.mod")
	if size(srcGoMod) != size(srcUGoMod) {
		err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if path == srcUGoMod {
				// Leave for last, in case we are racing with another golang command.
				return nil
			}
			if pdir, name := filepath.Split(path); name == "_golang.mod" {
				if err := raceSafeCopy(path, pdir+"golang.mod"); err != nil {
					return err
				}
			}
			return nil
		})
		// Handle src/golang.mod; this is the signal to other racing golang commands
		// that everything is okay and they can skip this step.
		if err == nil {
			err = raceSafeCopy(srcUGoMod, srcGoMod)
		}
		if err != nil {
			base.Fatalf("download %s: %v", golangtoolchain, err)
		}
	}

	// Reinvoke the golang command.
	execGoToolchain(golangtoolchain, dir, filepath.Join(dir, "bin/golang"))
}

func size(path string) int64 {
	info, err := os.Stat(path)
	if err != nil {
		return -1
	}
	return info.Size()
}

var writeBits fs.FileMode

// raceSafeCopy copies the file old to the file new, being careful to ensure
// that if multiple golang commands call raceSafeCopy(old, new) at the same time,
// they don't interfere with each other: both will succeed and return and
// later observe the correct content in new. Like in the build cache, we arrange
// this by opening new without truncation and then writing the content.
// Both golang commands can do this simultaneously and will write the same thing
// (old never changes content).
func raceSafeCopy(old, new string) error {
	oldInfo, err := os.Stat(old)
	if err != nil {
		return err
	}
	newInfo, err := os.Stat(new)
	if err == nil && newInfo.Size() == oldInfo.Size() {
		return nil
	}
	data, err := os.ReadFile(old)
	if err != nil {
		return err
	}
	// The module cache has unwritable directories by default.
	// Restore the user write bit in the directory so we can create
	// the new golang.mod file. We clear it again at the end on a
	// best-effort basis (ignoring failures).
	dir := filepath.Dir(old)
	info, err := os.Stat(dir)
	if err != nil {
		return err
	}
	if err := os.Chmod(dir, info.Mode()|writeBits); err != nil {
		return err
	}
	defer os.Chmod(dir, info.Mode())
	// Note: create the file writable, so that a racing golang command
	// doesn't get an error before we store the actual data.
	f, err := os.OpenFile(new, os.O_CREATE|os.O_WRONLY, writeBits&^0o111)
	if err != nil {
		// If OpenFile failed because a racing golang command completed our work
		// (and then OpenFile failed because the directory or file is now read-only),
		// count that as a success.
		if size(old) == size(new) {
			return nil
		}
		return err
	}
	defer os.Chmod(new, oldInfo.Mode())
	if _, err := f.Write(data); err != nil {
		f.Close()
		return err
	}
	return f.Close()
}

// modGoToolchain finds the enclosing golang.work or golang.mod file
// and returns the golang version and toolchain lines from the file.
// The toolchain line overrides the version line
func modGoToolchain() (file, golangVers, toolchain string) {
	wd := base.UncachedCwd()
	file = modload.FindGoWork(wd)
	// $GOWORK can be set to a file that does not yet exist, if we are running 'golang work init'.
	// Do not try to load the file in that case
	if _, err := os.Stat(file); err != nil {
		file = ""
	}
	if file == "" {
		file = modload.FindGoMod(wd)
	}
	if file == "" {
		return "", "", ""
	}

	data, err := os.ReadFile(file)
	if err != nil {
		base.Fatalf("%v", err)
	}
	return file, golangver.GoModLookup(data, "golang"), golangver.GoModLookup(data, "toolchain")
}

// maybeSwitchForGoInstallVersion reports whether the command line is golang install m@v or golang run m@v.
// If so, switch to the golang version required to build m@v if it's higher than minVers.
func maybeSwitchForGoInstallVersion(minVers string) {
	// Note: We assume there are no flags between 'golang' and 'install' or 'run'.
	// During testing there are some debugging flags that are accepted
	// in that position, but in production golang binaries there are not.
	if len(os.Args) < 3 {
		return
	}

	var cmdFlags *flag.FlagSet
	switch os.Args[1] {
	default:
		// Command doesn't support a pkg@version as the main module.
		return
	case "install":
		cmdFlags = &work.CmdInstall.Flag
	case "run":
		cmdFlags = &run.CmdRun.Flag
	}

	// The modcachrw flag is unique, in that it affects how we fetch the
	// requested module to even figure out what toolchain it needs.
	// We need to actually set it before we check the toolchain version.
	// (See https://golang.dev/issue/64282.)
	modcacherwFlag := cmdFlags.Lookup("modcacherw")
	if modcacherwFlag == nil {
		base.Fatalf("internal error: modcacherw flag not registered for command")
	}
	modcacherwVal, ok := modcacherwFlag.Value.(interface {
		IsBoolFlag() bool
		flag.Value
	})
	if !ok || !modcacherwVal.IsBoolFlag() {
		base.Fatalf("internal error: modcacherw is not a boolean flag")
	}

	// Make a best effort to parse the command's args to find the pkg@version
	// argument and the -modcacherw flag.
	var (
		pkgArg         string
		modcacherwSeen bool
	)
	for args := os.Args[2:]; len(args) > 0; {
		a := args[0]
		args = args[1:]
		if a == "--" {
			if len(args) == 0 {
				return
			}
			pkgArg = args[0]
			break
		}

		a, ok := strings.CutPrefix(a, "-")
		if !ok {
			// Not a flag argument. Must be a package.
			pkgArg = a
			break
		}
		a = strings.TrimPrefix(a, "-") // Treat --flag as -flag.

		name, val, hasEq := strings.Cut(a, "=")

		if name == "modcacherw" {
			if !hasEq {
				val = "true"
			}
			if err := modcacherwVal.Set(val); err != nil {
				return
			}
			modcacherwSeen = true
			continue
		}

		if hasEq {
			// Already has a value; don't bother parsing it.
			continue
		}

		f := run.CmdRun.Flag.Lookup(a)
		if f == nil {
			// We don't know whether this flag is a boolean.
			if os.Args[1] == "run" {
				// We don't know where to find the pkg@version argument.
				// For run, the pkg@version can be anywhere on the command line,
				// because it is preceded by run flags and followed by arguments to the
				// program being run. Since we don't know whether this flag takes
				// an argument, we can't reliably identify the end of the run flags.
				// Just give up and let the user clarify using the "=" form.
				return
			}

			// We would like to let 'golang install -newflag pkg@version' work even
			// across a toolchain switch. To make that work, assume by default that
			// the pkg@version is the last argument and skip the remaining args unless
			// we spot a plausible "-modcacherw" flag.
			for len(args) > 0 {
				a := args[0]
				name, _, _ := strings.Cut(a, "=")
				if name == "-modcacherw" || name == "--modcacherw" {
					break
				}
				if len(args) == 1 && !strings.HasPrefix(a, "-") {
					pkgArg = a
				}
				args = args[1:]
			}
			continue
		}

		if bf, ok := f.Value.(interface{ IsBoolFlag() bool }); !ok || !bf.IsBoolFlag() {
			// The next arg is the value for this flag. Skip it.
			args = args[1:]
			continue
		}
	}

	if !strings.Contains(pkgArg, "@") || build.IsLocalImport(pkgArg) || filepath.IsAbs(pkgArg) {
		return
	}
	path, version, _ := strings.Cut(pkgArg, "@")
	if path == "" || version == "" || golangver.IsToolchain(path) {
		return
	}

	if !modcacherwSeen && base.InGOFLAGS("-modcacherw") {
		fs := flag.NewFlagSet("golangInstallVersion", flag.ExitOnError)
		fs.Var(modcacherwVal, "modcacherw", modcacherwFlag.Usage)
		base.SetFromGOFLAGS(fs)
	}

	// It would be correct to do nothing here, and let "golang run" or "golang install"
	// do the toolchain switch.
	// Our golangal instead is, since we have golangne to the trouble of handling
	// unknown flags to some degree, to run the switch now, so that
	// these commands can switch to a newer toolchain directed by the
	// golang.mod which may actually understand the flag.
	// This was brought up during the golang.dev/issue/57001 proposal discussion
	// and may end up being common in self-contained "golang install" or "golang run"
	// command lines if we add new flags in the future.

	// Set up modules without an explicit golang.mod, to download golang.mod.
	modload.ForceUseModules = true
	modload.RootMode = modload.NoRoot
	modload.Init()
	defer modload.Reset()

	// See internal/load.PackagesAndErrorsOutsideModule
	ctx := context.Background()
	allowed := modload.CheckAllowed
	if modload.IsRevisionQuery(path, version) {
		// Don't check for retractions if a specific revision is requested.
		allowed = nil
	}
	noneSelected := func(path string) (version string) { return "none" }
	_, err := modload.QueryPackages(ctx, path, version, noneSelected, allowed)
	if errors.Is(err, golangver.ErrTooNew) {
		// Run early switch, same one golang install or golang run would eventually do,
		// if it understood all the command-line flags.
		var s Switcher
		s.Error(err)
		if s.TooNew != nil && golangver.Compare(s.TooNew.GoVersion, minVers) > 0 {
			SwitchOrFatal(ctx, err)
		}
	}
}
