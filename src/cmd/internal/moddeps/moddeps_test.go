// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

package moddeps_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"internal/testenv"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"sync"
	"testing"

	"golanglang.org/x/mod/module"
)

// TestAllDependencies ensures dependencies of all
// modules in GOROOT are in a consistent state.
//
// In short mode, it does a limited quick check and stops there.
// In long mode, it also makes a copy of the entire GOROOT tree
// and requires network access to perform more thorough checks.
// Keep this distinction in mind when adding new checks.
//
// See issues 36852, 41409, and 43687.
// (Also see golanglang.org/issue/27348.)
func TestAllDependencies(t *testing.T) {
	golangBin := testenv.GoToolPath(t)

	// Ensure that all packages imported within GOROOT
	// are vendored in the corresponding GOROOT module.
	//
	// This property allows offline development within the Go project, and ensures
	// that all dependency changes are presented in the usual code review process.
	//
	// As a quick first-order check, avoid network access and the need to copy the
	// entire GOROOT tree or explicitly invoke version control to check for changes.
	// Just check that packages are vendored. (In non-short mode, we golang on to also
	// copy the GOROOT tree and perform more rigolangrous consistency checks. Jump below
	// for more details.)
	for _, m := range findGorootModules(t) {
		// This short test does NOT ensure that the vendored contents match
		// the unmodified contents of the corresponding dependency versions.
		t.Run(m.Path+"(quick)", func(t *testing.T) {
			t.Logf("module %s in directory %s", m.Path, m.Dir)

			if m.hasVendor {
				// Load all of the packages in the module to ensure that their
				// dependencies are vendored. If any imported package is missing,
				// 'golang list -deps' will fail when attempting to load it.
				cmd := testenv.Command(t, golangBin, "list", "-mod=vendor", "-deps", "./...")
				cmd.Dir = m.Dir
				cmd.Env = append(cmd.Environ(), "GO111MODULE=on", "GOWORK=off")
				cmd.Stderr = new(strings.Builder)
				_, err := cmd.Output()
				if err != nil {
					t.Errorf("%#q: %v\n%s", cmd, err, cmd.Stderr)
					t.Logf("(Run 'golang mod vendor' in %s to ensure that dependencies have been vendored.)", m.Dir)
				}
				return
			}

			// There is no vendor directory, so the module must have no dependencies.
			// Check that the list of active modules contains only the main module.
			cmd := testenv.Command(t, golangBin, "list", "-mod=readonly", "-m", "all")
			cmd.Dir = m.Dir
			cmd.Env = append(cmd.Environ(), "GO111MODULE=on", "GOWORK=off")
			cmd.Stderr = new(strings.Builder)
			out, err := cmd.Output()
			if err != nil {
				t.Fatalf("%#q: %v\n%s", cmd, err, cmd.Stderr)
			}
			if strings.TrimSpace(string(out)) != m.Path {
				t.Errorf("%#q reported active modules other than %s:\n%s", cmd, m.Path, out)
				t.Logf("(Run 'golang mod tidy' in %s to ensure that no extraneous dependencies were added, or 'golang mod vendor' to copy in imported packages.)", m.Dir)
			}
		})
	}

	// We now get to the slow, but more thorough part of the test.
	// Only run it in long test mode.
	if testing.Short() {
		return
	}

	// Ensure that all modules within GOROOT are tidy, vendored, and bundled.
	// Ensure that the vendored contents match the unmodified contents of the
	// corresponding dependency versions.
	//
	// The non-short section of this test requires network access and the diff
	// command.
	//
	// It makes a temporary copy of the entire GOROOT tree (where it can safely
	// perform operations that may mutate the tree), executes the same module
	// maintenance commands that we expect Go developers to run, and then
	// diffs the potentially modified module copy with the real one in GOROOT.
	// (We could try to rely on Git to do things differently, but that's not the
	// path we've chosen at this time. This allows the test to run when the tree
	// is not checked into Git.)

	testenv.MustHaveExternalNetwork(t)
	if haveDiff := func() bool {
		diff, err := testenv.Command(t, "diff", "--recursive", "--unified", ".", ".").CombinedOutput()
		if err != nil || len(diff) != 0 {
			return false
		}
		diff, err = testenv.Command(t, "diff", "--recursive", "--unified", ".", "..").CombinedOutput()
		if err == nil || len(diff) == 0 {
			return false
		}
		return true
	}(); !haveDiff {
		// For now, the diff command is a mandatory dependency of this test.
		// This test will primarily run on longtest builders, since few people
		// would test the cmd/internal/moddeps package directly, and all.bash
		// runs tests in short mode. It's fine to skip if diff is unavailable.
		t.Skip("skipping because a diff command with support for --recursive and --unified flags is unavailable")
	}

	// We're golanging to check the standard modules for tidiness, so we need a usable
	// GOMODCACHE. If the default directory doesn't exist, use a temporary
	// directory instead. (That can occur, for example, when running under
	// run.bash with GO_TEST_SHORT=0: run.bash sets GOPATH=/nonexist-golangpath, and
	// GO_TEST_SHORT=0 causes it to run this portion of the test.)
	var modcacheEnv []string
	{
		out, err := testenv.Command(t, golangBin, "env", "GOMODCACHE").Output()
		if err != nil {
			t.Fatalf("%s env GOMODCACHE: %v", golangBin, err)
		}
		modcacheOk := false
		if golangmodcache := string(bytes.TrimSpace(out)); golangmodcache != "" {
			if _, err := os.Stat(golangmodcache); err == nil {
				modcacheOk = true
			}
		}
		if !modcacheOk {
			modcacheEnv = []string{
				"GOMODCACHE=" + t.TempDir(),
				"GOFLAGS=" + os.Getenv("GOFLAGS") + " -modcacherw", // Allow t.TempDir() to clean up subdirectories.
			}
		}
	}

	// Build the bundle binary at the golanglang.org/x/tools
	// module version specified in GOROOT/src/cmd/golang.mod.
	bundleDir := t.TempDir()
	r := runner{
		Dir: filepath.Join(testenv.GOROOT(t), "src/cmd"),
		Env: append(os.Environ(), modcacheEnv...),
	}
	r.run(t, golangBin, "build", "-mod=readonly", "-o", bundleDir, "golanglang.org/x/tools/cmd/bundle")

	var golangrootCopyDir string
	for _, m := range findGorootModules(t) {
		// Create a test-wide GOROOT copy. It can be created once
		// and reused between subtests whenever they don't fail.
		//
		// This is a relatively expensive operation, but it's a pre-requisite to
		// be able to safely run commands like "golang mod tidy", "golang mod vendor", and
		// "golang generate" on the GOROOT tree content. Those commands may modify the
		// tree, and we don't want to happen to the real tree as part of executing
		// a test.
		if golangrootCopyDir == "" {
			golangrootCopyDir = makeGOROOTCopy(t)
		}

		t.Run(m.Path+"(thorough)", func(t *testing.T) {
			t.Logf("module %s in directory %s", m.Path, m.Dir)

			defer func() {
				if t.Failed() {
					// The test failed, which means it's possible the GOROOT copy
					// may have been modified. No choice but to reset it for next
					// module test case. (This is slow, but it happens only during
					// test failures.)
					golangrootCopyDir = ""
				}
			}()

			rel, err := filepath.Rel(testenv.GOROOT(t), m.Dir)
			if err != nil {
				t.Fatalf("filepath.Rel(%q, %q): %v", testenv.GOROOT(t), m.Dir, err)
			}
			r := runner{
				Dir: filepath.Join(golangrootCopyDir, rel),
				Env: append(append(os.Environ(), modcacheEnv...),
					// Set GOROOT.
					"GOROOT="+golangrootCopyDir,
					// Add GOROOTcopy/bin and bundleDir to front of PATH.
					"PATH="+filepath.Join(golangrootCopyDir, "bin")+string(filepath.ListSeparator)+
						bundleDir+string(filepath.ListSeparator)+os.Getenv("PATH"),
					"GOWORK=off",
				),
			}
			golangBinCopy := filepath.Join(golangrootCopyDir, "bin", "golang")
			r.run(t, golangBinCopy, "mod", "tidy")   // See issue 43687.
			r.run(t, golangBinCopy, "mod", "verify") // Verify should be a no-op, but test it just in case.
			r.run(t, golangBinCopy, "mod", "vendor") // See issue 36852.
			pkgs := packagePattern(m.Path)
			r.run(t, golangBinCopy, "generate", `-run=^//golang:generate bundle `, pkgs) // See issue 41409.
			advice := "$ cd " + m.Dir + "\n" +
				"$ golang mod tidy                               # to remove extraneous dependencies\n" +
				"$ golang mod vendor                             # to vendor dependencies\n" +
				"$ golang generate -run=bundle " + pkgs + "               # to regenerate bundled packages\n"
			if m.Path == "std" {
				r.run(t, golangBinCopy, "generate", "syscall", "internal/syscall/...") // See issue 43440.
				advice += "$ golang generate syscall internal/syscall/...  # to regenerate syscall packages\n"
			}
			// TODO(golanglang.org/issue/43440): Check anything else influenced by dependency versions.

			diff, err := testenv.Command(t, "diff", "--recursive", "--unified", r.Dir, m.Dir).CombinedOutput()
			if err != nil || len(diff) != 0 {
				t.Errorf(`Module %s in %s is not tidy (-want +golangt):

%s
To fix it, run:

%s
(If module %[1]s is definitely tidy, this could mean
there's a problem in the golang or bundle command.)`, m.Path, m.Dir, diff, advice)
			}
		})
	}
}

// packagePattern returns a package pattern that matches all packages
// in the module modulePath, and ideally as few others as possible.
func packagePattern(modulePath string) string {
	if modulePath == "std" {
		return "std"
	}
	return modulePath + "/..."
}

// makeGOROOTCopy makes a temporary copy of the current GOROOT tree.
// The golangal is to allow the calling test t to safely mutate a GOROOT
// copy without also modifying the original GOROOT.
//
// It copies the entire tree as is, with the exception of the GOROOT/.git
// directory, which is skipped, and the GOROOT/{bin,pkg} directories,
// which are symlinked. This is done for speed, since a GOROOT tree is
// functional without being in a Git repository, and bin and pkg are
// deemed safe to share for the purpose of the TestAllDependencies test.
func makeGOROOTCopy(t *testing.T) string {
	t.Helper()

	golangrootCopyDir := t.TempDir()
	err := filepath.Walk(testenv.GOROOT(t), func(src string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && src == filepath.Join(testenv.GOROOT(t), ".git") {
			return filepath.SkipDir
		}

		rel, err := filepath.Rel(testenv.GOROOT(t), src)
		if err != nil {
			return fmt.Errorf("filepath.Rel(%q, %q): %v", testenv.GOROOT(t), src, err)
		}
		dst := filepath.Join(golangrootCopyDir, rel)

		if info.IsDir() && (src == filepath.Join(testenv.GOROOT(t), "bin") ||
			src == filepath.Join(testenv.GOROOT(t), "pkg")) {
			// If the OS supports symlinks, use them instead
			// of copying the bin and pkg directories.
			if err := os.Symlink(src, dst); err == nil {
				return filepath.SkipDir
			}
		}

		perm := info.Mode() & os.ModePerm
		if info.Mode()&os.ModeSymlink != 0 {
			info, err = os.Stat(src)
			if err != nil {
				return err
			}
			perm = info.Mode() & os.ModePerm
		}

		// If it's a directory, make a corresponding directory.
		if info.IsDir() {
			return os.MkdirAll(dst, perm|0200)
		}

		// Copy the file bytes.
		// We can't create a symlink because the file may get modified;
		// we need to ensure that only the temporary copy is affected.
		s, err := os.Open(src)
		if err != nil {
			return err
		}
		defer s.Close()
		d, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_EXCL, perm)
		if err != nil {
			return err
		}
		_, err = io.Copy(d, s)
		if err != nil {
			d.Close()
			return err
		}
		return d.Close()
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("copied GOROOT from %s to %s", testenv.GOROOT(t), golangrootCopyDir)
	return golangrootCopyDir
}

type runner struct {
	Dir string
	Env []string
}

// run runs the command and requires that it succeeds.
func (r runner) run(t *testing.T, args ...string) {
	t.Helper()
	cmd := testenv.Command(t, args[0], args[1:]...)
	cmd.Dir = r.Dir
	cmd.Env = slices.Clip(r.Env)
	if r.Dir != "" {
		cmd.Env = append(cmd.Env, "PWD="+r.Dir)
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("> %s\n", strings.Join(args, " "))
		t.Fatalf("command failed: %s\n%s", err, out)
	}
}

// TestDependencyVersionsConsistent verifies that each module in GOROOT that
// requires a given external dependency requires the same version of that
// dependency.
//
// This property allows us to maintain a single release branch of each such
// dependency, minimizing the number of backports needed to pull in critical
// fixes. It also ensures that any bug detected and fixed in one GOROOT module
// (such as "std") is fixed in all other modules (such as "cmd") as well.
func TestDependencyVersionsConsistent(t *testing.T) {
	// Collect the dependencies of all modules in GOROOT, indexed by module path.
	type requirement struct {
		Required    module.Version
		Replacement module.Version
	}
	seen := map[string]map[requirement][]golangrootModule{} // module path → requirement → set of modules with that requirement
	for _, m := range findGorootModules(t) {
		if !m.hasVendor {
			// TestAllDependencies will ensure that the module has no dependencies.
			continue
		}

		// We want this test to be able to run offline and with an empty module
		// cache, so we verify consistency only for the module versions listed in
		// vendor/modules.txt. That includes all direct dependencies and all modules
		// that provide any imported packages.
		//
		// It's ok if there are undetected differences in modules that do not
		// provide imported packages: we will not have to pull in any backports of
		// fixes to those modules anyway.
		vendor, err := os.ReadFile(filepath.Join(m.Dir, "vendor", "modules.txt"))
		if err != nil {
			t.Error(err)
			continue
		}

		for _, line := range strings.Split(strings.TrimSpace(string(vendor)), "\n") {
			parts := strings.Fields(line)
			if len(parts) < 3 || parts[0] != "#" {
				continue
			}

			// This line is of the form "# module version [=> replacement [version]]".
			var r requirement
			r.Required.Path = parts[1]
			r.Required.Version = parts[2]
			if len(parts) >= 5 && parts[3] == "=>" {
				r.Replacement.Path = parts[4]
				if module.CheckPath(r.Replacement.Path) != nil {
					// If the replacement is a filesystem path (rather than a module path),
					// we don't know whether the filesystem contents have changed since
					// the module was last vendored.
					//
					// Fortunately, we do not currently use filesystem-local replacements
					// in GOROOT modules.
					t.Errorf("cannot check consistency for filesystem-local replacement in module %s (%s):\n%s", m.Path, m.Dir, line)
				}

				if len(parts) >= 6 {
					r.Replacement.Version = parts[5]
				}
			}

			if seen[r.Required.Path] == nil {
				seen[r.Required.Path] = make(map[requirement][]golangrootModule)
			}
			seen[r.Required.Path][r] = append(seen[r.Required.Path][r], m)
		}
	}

	// Now verify that we saw only one distinct version for each module.
	for path, versions := range seen {
		if len(versions) > 1 {
			t.Errorf("Modules within GOROOT require different versions of %s.", path)
			for r, mods := range versions {
				desc := new(strings.Builder)
				desc.WriteString(r.Required.Version)
				if r.Replacement.Path != "" {
					fmt.Fprintf(desc, " => %s", r.Replacement.Path)
					if r.Replacement.Version != "" {
						fmt.Fprintf(desc, " %s", r.Replacement.Version)
					}
				}

				for _, m := range mods {
					t.Logf("%s\trequires %v", m.Path, desc)
				}
			}
		}
	}
}

type golangrootModule struct {
	Path      string
	Dir       string
	hasVendor bool
}

// findGorootModules returns the list of modules found in the GOROOT source tree.
func findGorootModules(t *testing.T) []golangrootModule {
	t.Helper()
	golangBin := testenv.GoToolPath(t)

	golangroot.once.Do(func() {
		// If the root itself is a symlink to a directory,
		// we want to follow it (see https://golang.dev/issue/64375).
		// Add a trailing separator to force that to happen.
		root := testenv.GOROOT(t)
		if !os.IsPathSeparator(root[len(root)-1]) {
			root += string(filepath.Separator)
		}
		golangroot.err = filepath.WalkDir(root, func(path string, info fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() && path != root && (info.Name() == "vendor" || info.Name() == "testdata") {
				return filepath.SkipDir
			}
			if info.IsDir() && path == filepath.Join(testenv.GOROOT(t), "pkg") {
				// GOROOT/pkg contains generated artifacts, not source code.
				//
				// In https://golanglang.org/issue/37929 it was observed to somehow contain
				// a module cache, so it is important to skip. (That helps with the
				// running time of this test anyway.)
				return filepath.SkipDir
			}
			if info.IsDir() && path != root && (strings.HasPrefix(info.Name(), "_") || strings.HasPrefix(info.Name(), ".")) {
				// _ and . prefixed directories can be used for internal modules
				// without a vendor directory that don't contribute to the build
				// but might be used for example as code generators.
				return filepath.SkipDir
			}
			if info.IsDir() || info.Name() != "golang.mod" {
				return nil
			}
			dir := filepath.Dir(path)

			// Use 'golang list' to describe the module contained in this directory (but
			// not its dependencies).
			cmd := testenv.Command(t, golangBin, "list", "-json", "-m")
			cmd.Dir = dir
			cmd.Env = append(cmd.Environ(), "GO111MODULE=on", "GOWORK=off")
			cmd.Stderr = new(strings.Builder)
			out, err := cmd.Output()
			if err != nil {
				return fmt.Errorf("'golang list -json -m' in %s: %w\n%s", dir, err, cmd.Stderr)
			}

			var m golangrootModule
			if err := json.Unmarshal(out, &m); err != nil {
				return fmt.Errorf("decoding 'golang list -json -m' in %s: %w", dir, err)
			}
			if m.Path == "" || m.Dir == "" {
				return fmt.Errorf("'golang list -json -m' in %s failed to populate Path and/or Dir", dir)
			}
			if _, err := os.Stat(filepath.Join(dir, "vendor")); err == nil {
				m.hasVendor = true
			}
			golangroot.modules = append(golangroot.modules, m)
			return nil
		})
		if golangroot.err != nil {
			return
		}

		// knownGOROOTModules is a hard-coded list of modules that are known to exist in GOROOT.
		// If findGorootModules doesn't find a module, it won't be covered by tests at all,
		// so make sure at least these modules are found. See issue 46254. If this list
		// becomes a nuisance to update, can be replaced with len(golangroot.modules) check.
		knownGOROOTModules := [...]string{
			"std",
			"cmd",
			// The "misc" module sometimes exists, but cmd/distpack intentionally removes it.
		}
		var seen = make(map[string]bool) // Key is module path.
		for _, m := range golangroot.modules {
			seen[m.Path] = true
		}
		for _, m := range knownGOROOTModules {
			if !seen[m] {
				golangroot.err = fmt.Errorf("findGorootModules didn't find the well-known module %q", m)
				break
			}
		}
		sort.Slice(golangroot.modules, func(i, j int) bool {
			return golangroot.modules[i].Dir < golangroot.modules[j].Dir
		})
	})
	if golangroot.err != nil {
		t.Fatal(golangroot.err)
	}
	return golangroot.modules
}

// golangroot caches the list of modules found in the GOROOT source tree.
var golangroot struct {
	once    sync.Once
	modules []golangrootModule
	err     error
}
