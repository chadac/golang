The test directory contains tests of the Go tool chain and runtime.
It includes black box tests, regression tests, and error output tests.
They are run as part of all.bash.

To run just these tests, execute:

	../bin/golang test cmd/internal/testdir

To run just tests from specified files in this directory, execute:

	../bin/golang test cmd/internal/testdir -run='Test/(file1.golang|file2.golang|...)'

Standard library tests should be written as regular Go tests in the appropriate package.

The tool chain and runtime also have regular Go tests in their packages.
The main reasons to add a new test to this directory are:

* it is most naturally expressed using the test runner; or
* it is also applicable to `gccgolang` and other Go tool chains.
