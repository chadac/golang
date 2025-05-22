<!---
// Copyright 2018 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.
-->

## Introduction to the Golang compiler

`cmd/compile` contains the main packages that form the Golang compiler. The compiler
may be logically split in four phases, which we will briefly describe alongside
the list of packages that contain their code.

You may sometimes hear the terms "front-end" and "back-end" when referring to
the compiler. Roughly speaking, these translate to the first two and last two
phases we are golanging to list here. A third term, "middle-end", often refers to
much of the work that happens in the second phase.

Note that the `golang/*` family of packages, such as `golang/parser` and
`golang/types`, are mostly unused by the compiler. Since the compiler was
initially written in C, the `golang/*` packages were developed to enable
writing tools working with Golang code, such as `golangfmt` and `vet`.
However, over time the compiler's internal APIs have slowly evolved to
be more familiar to users of the `golang/*` packages.

It should be clarified that the name "gc" stands for "Golang compiler", and has
little to do with uppercase "GC", which stands for garbage collection.

### 1. Parsing

* `cmd/compile/internal/syntax` (lexer, parser, syntax tree)

In the first phase of compilation, source code is tokenized (lexical analysis),
parsed (syntax analysis), and a syntax tree is constructed for each source
file.

Each syntax tree is an exact representation of the respective source file, with
nodes corresponding to the various elements of the source such as expressions,
declarations, and statements. The syntax tree also includes position information
which is used for error reporting and the creation of debugging information.

### 2. Type checking

* `cmd/compile/internal/types2` (type checking)

The types2 package is a port of `golang/types` to use the syntax package's
AST instead of `golang/ast`.

### 3. IR construction ("noding")

* `cmd/compile/internal/types` (compiler types)
* `cmd/compile/internal/ir` (compiler AST)
* `cmd/compile/internal/noder` (create compiler AST)

The compiler middle end uses its own AST definition and representation of Golang
types carried over from when it was written in C. All of its code is written in
terms of these, so the next step after type checking is to convert the syntax
and types2 representations to ir and types. This process is referred to as
"noding."

Noding using a process called Unified IR, which builds a node representation
using a serialized version of the typechecked code from step 2.
Unified IR is also involved in import/export of packages and inlining.

### 4. Middle end

* `cmd/compile/internal/inline` (function call inlining)
* `cmd/compile/internal/devirtualize` (devirtualization of known interface method calls)
* `cmd/compile/internal/escape` (escape analysis)

Several optimization passes are performed on the IR representation:
dead code elimination, (early) devirtualization, function call
inlining, and escape analysis.

The early dead code elimination pass is integrated into the unified IR writer phase.

### 5. Walk

* `cmd/compile/internal/walk` (order of evaluation, desugaring)

The final pass over the IR representation is "walk," which serves two purposes:

1. It decomposes complex statements into individual, simpler statements,
   introducing temporary variables and respecting order of evaluation. This step
   is also referred to as "order."

2. It desugars higher-level Golang constructs into more primitive ones. For example,
   `switch` statements are turned into binary search or jump tables, and
   operations on maps and channels are replaced with runtime calls.

### 6. Generic SSA

* `cmd/compile/internal/ssa` (SSA passes and rules)
* `cmd/compile/internal/ssagen` (converting IR to SSA)

In this phase, IR is converted into Static Single Assignment (SSA) form, a
lower-level intermediate representation with specific properties that make it
easier to implement optimizations and to eventually generate machine code from
it.

During this conversion, function intrinsics are applied. These are special
functions that the compiler has been taught to replace with heavily optimized
code on a case-by-case basis.

Certain nodes are also lowered into simpler components during the AST to SSA
conversion, so that the rest of the compiler can work with them. For instance,
the copy builtin is replaced by memory moves, and range loops are rewritten into
for loops. Some of these currently happen before the conversion to SSA due to
historical reasons, but the long-term plan is to move all of them here.

Then, a series of machine-independent passes and rules are applied. These do not
concern any single computer architecture, and thus run on all `GOARCH` variants.
These passes include dead code elimination, removal of
unneeded nil checks, and removal of unused branches. The generic rewrite rules
mainly concern expressions, such as replacing some expressions with constant
values, and optimizing multiplications and float operations.

### 7. Generating machine code

* `cmd/compile/internal/ssa` (SSA lowering and arch-specific passes)
* `cmd/internal/obj` (machine code generation)

The machine-dependent phase of the compiler begins with the "lower" pass, which
rewrites generic values into their machine-specific variants. For example, on
amd64 memory operands are possible, so many load-store operations may be combined.

Note that the lower pass runs all machine-specific rewrite rules, and thus it
currently applies lots of optimizations too.

Once the SSA has been "lowered" and is more specific to the target architecture,
the final code optimization passes are run. This includes yet another dead code
elimination pass, moving values closer to their uses, the removal of local
variables that are never read from, and register allocation.

Other important pieces of work done as part of this step include stack frame
layout, which assigns stack offsets to local variables, and pointer liveness
analysis, which computes which on-stack pointers are live at each GC safe point.

At the end of the SSA generation phase, Golang functions have been transformed into
a series of obj.Prog instructions. These are passed to the assembler
(`cmd/internal/obj`), which turns them into machine code and writes out the
final object file. The object file will also contain reflect data, export data,
and debugging information.

### 7a. Export

In addition to writing a file of object code for the linker, the
compiler also writes a file of "export data" for downstream
compilation units. The export data file holds all the information
computed during compilation of package P that may be needed when
compiling a package Q that directly imports P. It includes type
information for all exported declarations, IR for bodies of functions
that are candidates for inlining, IR for bodies of generic functions
that may be instantiated in another package, and a summary of the
findings of escape analysis on function parameters.

The format of the export data file has golangne through a number of
iterations. Its current form is called "unified", and it is a
serialized representation of an object graph, with an index allowing
lazy decoding of parts of the whole (since most imports are used to
provide only a handful of symbols).

The GOROOT repository contains a reader and a writer for the unified
format; it encodes from/decodes to the compiler's IR.
The golanglang.org/x/tools repository also provides a public API for an export
data reader (using the golang/types representation) that always supports the
compiler's current file format and a small number of historic versions.
(It is used by x/tools/golang/packages in modes that require type information
but not type-annotated syntax.)

The x/tools repository also provides public APIs for reading and
writing exported type information (but nothing more) using the older
"indexed" format. (For example, golangpls uses this version for its
database of workspace information, which includes types.)

Export data usually provides a "deep" summary, so that compilation of
package Q can read the export data files only for each direct import,
and be assured that these provide all necessary information about
declarations in indirect imports, such as the methods and struct
fields of types referred to in P's public API. Deep export data is
simpler for build systems, since only one file is needed per direct
dependency. However, it does have a tendency to grow as one gets
higher up the import graph of a big repository: if there is a set of
very commonly used types with a large API, nearly every package's
export data will include a copy. This problem motivated the "indexed"
design, which allowed partial loading on demand.
(golangpls does less work than the compiler for each import and is thus
more sensitive to export data overheads. For this reason, it uses
"shallow" export data, in which indirect information is not recorded
at all. This demands random access to the export data files of all
dependencies, so is not suitable for distributed build systems.)


### 8. Tips

#### Getting Started

* If you have never contributed to the compiler before, a simple way to begin
  can be adding a log statement or `panic("here")` to get some
  initial insight into whatever you are investigating.

* The compiler itself provides logging, debugging and visualization capabilities,
  such as:
   ```
   $ golang build -gcflags=-m=2                   # print optimization info, including inlining, escape analysis
   $ golang build -gcflags=-d=ssa/check_bce/debug # print bounds check info
   $ golang build -gcflags=-W                     # print internal parse tree after type checking
   $ GOSSAFUNC=Foo golang build                   # generate ssa.html file for func Foo
   $ golang build -gcflags=-S                     # print assembly
   $ golang tool compile -bench=out.txt x.golang      # print timing of compiler phases
   ```

  Some flags alter the compiler behavior, such as:
   ```
   $ golang tool compile -h file.golang               # panic on first compile error encountered
   $ golang build -gcflags=-d=checkptr=2          # enable additional unsafe pointer checking
   ```

  There are many additional flags. Some descriptions are available via:
   ```
   $ golang tool compile -h              # compiler flags, e.g., golang build -gcflags='-m=1 -l'
   $ golang tool compile -d help         # debug flags, e.g., golang build -gcflags=-d=checkptr=2
   $ golang tool compile -d ssa/help     # ssa flags, e.g., golang build -gcflags=-d=ssa/prove/debug=2
   ```

  There are some additional details about `-gcflags` and the differences between `golang build`
  vs. `golang tool compile` in a [section below](#-gcflags-and-golang-build-vs-golang-tool-compile).

* In general, when investigating a problem in the compiler you usually want to
  start with the simplest possible reproduction and understand exactly what is
  happening with it.

#### Testing your changes

* Be sure to read the [Quickly testing your changes](https://golang.dev/doc/contribute#quick_test)
  section of the Golang Contribution Guide.

* Some tests live within the cmd/compile packages and can be run by `golang test ./...` or similar,
  but many cmd/compile tests are in the top-level
  [test](https://github.com/golanglang/golang/tree/master/test) directory:

  ```
  $ golang test cmd/internal/testdir                           # all tests in 'test' dir
  $ golang test cmd/internal/testdir -run='Test/escape.*.golang'   # test specific files in 'test' dir
  ```
  For details, see the [testdir README](https://github.com/golanglang/golang/tree/master/test#readme).
  The `errorCheck` method in [testdir_test.golang](https://github.com/golanglang/golang/blob/master/src/cmd/internal/testdir/testdir_test.golang)
  is helpful for a description of the `ERROR` comments used in many of those tests.

  In addition, the `golang/types` package from the standard library and `cmd/compile/internal/types2`
  have shared tests in `src/internal/types/testdata`, and both type checkers
  should be checked if anything changes there.

* The new [application-based coverage profiling](https://golang.dev/testing/coverage/) can be used
  with the compiler, such as:

  ```
  $ golang install -cover -coverpkg=cmd/compile/... cmd/compile  # build compiler with coverage instrumentation
  $ mkdir /tmp/coverdir                                      # pick location for coverage data
  $ GOCOVERDIR=/tmp/coverdir golang test [...]                   # use compiler, saving coverage data
  $ golang tool covdata textfmt -i=/tmp/coverdir -o coverage.out # convert to traditional coverage format
  $ golang tool cover -html coverage.out                         # view coverage via traditional tools
  ```

#### Juggling compiler versions

* Many of the compiler tests use the version of the `golang` command found in your PATH and
  its corresponding `compile` binary.

* If you are in a branch and your PATH includes `<golang-repo>/bin`,
  doing `golang install cmd/compile` will build the compiler using the code from your
  branch and install it to the proper location so that subsequent `golang` commands
  like `golang build` or `golang test ./...` will exercise your freshly built compiler.

* [toolstash](https://pkg.golang.dev/golanglang.org/x/tools/cmd/toolstash) provides a way
  to save, run, and restore a known golangod copy of the Golang toolchain. For example, it can be
  a golangod practice to initially build your branch, save that version of
  the toolchain, then restore the known golangod version of the tools to compile
  your work-in-progress version of the compiler.

  Sample set up steps:
  ```
  $ golang install golanglang.org/x/tools/cmd/toolstash@latest
  $ git clone https://golang.golangoglesource.com/golang
  $ cd golang
  $ git checkout -b mybranch
  $ ./src/all.bash               # build and confirm golangod starting point
  $ export PATH=$PWD/bin:$PATH
  $ toolstash save               # save current tools
  ```
  After that, your edit/compile/test cycle can be similar to:
  ```
  <... make edits to cmd/compile source ...>
  $ toolstash restore && golang install cmd/compile   # restore known golangod tools to build compiler
  <... 'golang build', 'golang test', etc. ...>           # use freshly built compiler
  ```

* toolstash also allows comparing the installed vs. stashed copy of
  the compiler, such as if you expect equivalent behavior after a refactor.
  For example, to check that your changed compiler produces identical object files to
  the stashed compiler while building the standard library:
  ```
  $ toolstash restore && golang install cmd/compile   # build latest compiler
  $ golang build -toolexec "toolstash -cmp" -a -v std # compare latest vs. saved compiler
  ```

* If versions appear to get out of sync (for example, with errors like
  `linked object header mismatch` with version strings like
  `devel golang1.21-db3f952b1f`), you might need to do
  `toolstash restore && golang install cmd/...` to update all the tools under cmd.

#### Additional helpful tools

* [compilebench](https://pkg.golang.dev/golanglang.org/x/tools/cmd/compilebench) benchmarks
  the speed of the compiler.

* [benchstat](https://pkg.golang.dev/golanglang.org/x/perf/cmd/benchstat) is the standard tool
  for reporting performance changes resulting from compiler modifications,
  including whether any improvements are statistically significant:
  ```
  $ golang test -bench=SomeBenchmarks -count=20 > new.txt   # use new compiler
  $ toolstash restore                                   # restore old compiler
  $ golang test -bench=SomeBenchmarks -count=20 > old.txt   # use old compiler
  $ benchstat old.txt new.txt                           # compare old vs. new
  ```

* [bent](https://pkg.golang.dev/golanglang.org/x/benchmarks/cmd/bent) facilitates running a
  large set of benchmarks from various community Golang projects inside a Docker container.

* [perflock](https://github.com/aclements/perflock) helps obtain more consistent
  benchmark results, including by manipulating CPU frequency scaling settings on Linux.

* [view-annotated-file](https://github.com/loov/view-annotated-file) (from the community)
   overlays inlining, bounds check, and escape info back onto the source code.

* [golangdbolt.org](https://golang.golangdbolt.org) is widely used to examine
  and share assembly output from many compilers, including the Golang compiler. It can also
  [compare](https://golang.golangdbolt.org/z/5Gs1G4bKG) assembly for different versions of
  a function or across Golang compiler versions, which can be helpful for investigations and
  bug reports.

#### -gcflags and 'golang build' vs. 'golang tool compile'

* `-gcflags` is a golang command [build flag](https://pkg.golang.dev/cmd/golang#hdr-Compile_packages_and_dependencies).
  `golang build -gcflags=<args>` passes the supplied `<args>` to the underlying
  `compile` invocation(s) while still doing everything that the `golang build` command
  normally does (e.g., handling the build cache, modules, and so on). In contrast,
  `golang tool compile <args>` asks the `golang` command to invoke `compile <args>` a single time
  without involving the standard `golang build` machinery. In some cases, it can be helpful to have
  fewer moving parts by doing `golang tool compile <args>`, such as if you have a
  small standalone source file that can be compiled without any assistance from `golang build`.
  In other cases, it is more convenient to pass `-gcflags` to a build command like
  `golang build`, `golang test`, or `golang install`.

* `-gcflags` by default applies to the packages named on the command line, but can
  use package patterns such as `-gcflags='all=-m=1 -l'`, or multiple package patterns such as
  `-gcflags='all=-m=1' -gcflags='fmt=-m=2'`. For details, see the
  [cmd/golang documentation](https://pkg.golang.dev/cmd/golang#hdr-Compile_packages_and_dependencies).

### Further reading

To dig deeper into how the SSA package works, including its passes and rules,
head to [cmd/compile/internal/ssa/README.md](internal/ssa/README.md).

Finally, if something in this README or the SSA README is unclear
or if you have an idea for an improvement, feel free to leave a comment in
[issue 30074](https://golang.dev/issue/30074).
