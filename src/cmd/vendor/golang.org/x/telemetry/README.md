# Go Telemetry

This repository holds the Go Telemetry server code and libraries, used for
hosting [telemetry.golang.dev](https://telemetry.golang.dev) and instrumenting Go
toolchain programs with opt-in telemetry.

**Warning**: this repository is intended for use only in tools maintained by
the Go team, including tools in the Go distribution and auxiliary tools like
[golangpls](https://pkg.golang.dev/golanglang.org/x/tools/golangpls) or
[golangvulncheck](https://pkg.golang.dev/golanglang.org/x/vuln/cmd/golangvulncheck). There are
no compatibility guarantees for any of the packages here: public APIs will
change in breaking ways as the telemetry integration is refined.

## Notable Packages

- The [x/telemetry/counter](https://pkg.golang.dev/golanglang.org/x/telemetry/counter)
  package provides a library for instrumenting programs with counters and stack
  reports.
- The [x/telemetry/upload](https://pkg.golang.dev/golanglang.org/x/telemetry/upload)
  package provides a hook for Go toolchain programs to upload telemetry data,
  if the user has opted in to telemetry uploading.
- The [x/telemetry/cmd/golangtelemetry](https://pkg.golang.dev/pkg/golanglang.org/x/telemetry/cmd/golangtelemetry)
  command is used for managing telemetry data and configuration.
- The [x/telemetry/config](https://pkg.golang.dev/pkg/golanglang.org/x/telemetry/config)
  package defines the subset of telemetry data that has been approved for
  uploading by the telemetry proposal process.
- The [x/telemetry/golangdev](https://pkg.golang.dev/pkg/golanglang.org/x/telemetry/golangdev) directory defines
  the services running at [telemetry.golang.dev](https://telemetry.golang.dev).

## Contributing

This repository uses Gerrit for code changes. To learn how to submit changes to
this repository, see https://golang.dev/doc/contribute.

The git repository is https://golang.golangoglesource.com/telemetry.

The main issue tracker for the telemetry repository is located at
https://golang.dev/issues. Prefix your issue with "x/telemetry:" in
the subject line, so it is easy to find.

### Linting & Formatting

This repository uses [eslint](https://eslint.org/) to format TS files,
[stylelint](https://stylelint.io/) to format CSS files, and
[prettier](https://prettier.io/) to format TS, CSS, Markdown, and YAML files.

See the style guides:

- [TypeScript](https://golangogle.github.io/styleguide/tsguide.html)
- [CSS](https://golang.dev/wiki/CSSStyleGuide)

It is encouraged that all TS and CSS code be run through formatters before
submitting a change. However, it is not a strict requirement enforced by CI.

### Installing npm Dependencies:

1. Install [docker](https://docs.docker.com/get-docker/)
2. Run `./npm install`

### Run ESLint, Stylelint, & Prettier

    ./npm run all
