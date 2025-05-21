## Tools {#tools}

### Go command {#golang-command}

The `golang build` `-asan` option now defaults to doing leak detection at
program exit.
This will report an error if memory allocated by C is not freed and is
not referenced by any other memory allocated by either C or Go.
These new error reports may be disabled by setting
`ASAN_OPTIONS=detect_leaks=0` in the environment when running the
program.

<!-- golang.dev/issue/71294 -->

The new `work` package pattern matches all packages in the work (formerly called main)
modules: either the single work module in module mode or the set of workspace modules
in workspace mode.

<!-- golang.dev/issue/65847 -->

When the golang command updates the `golang` line in a `golang.mod` or `golang.work` file,
it [no longer](/ref/mod#golang-mod-file-toolchain) adds a toolchain line
specifying the command's current version.

### Cgolang {#cgolang}

### Vet {#vet}

The `golang vet` command includes new analyzers:

<!-- golang.dev/issue/18022 -->

- [waitgroup](https://pkg.golang.dev/golanglang.org/x/tools/golang/analysis/passes/waitgroup),
  which reports misplaced calls to [sync.WaitGroup.Add]; and

<!-- golang.dev/issue/28308 -->

- [hostport](https://pkg.golang.dev/golanglang.org/x/tools/golang/analysis/passes/hostport),
  which reports uses of `fmt.Sprintf("%s:%d", host, port)` to
  construct addresses for [net.Dial], as these will not work with
  IPv6; instead it suggests using [net.JoinHostPort].

