package main

import (
	"flag"
	"fmt"
	"io"
	"runtime/debug"
)

// version is the release version stamped by the release build
// (.goreleaser.yaml passes -ldflags "-X main.version=..."). It is empty for
// `go build`/`go install` builds, which fall back to the module version Go
// records in the binary's build information.
var version string

// cmdVersion implements "esig version".
func cmdVersion(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("version", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() {
		fmt.Fprint(stderr, "Usage: esig version\n\nPrints the esig build's module version and Go toolchain version.\n")
	}
	if err := fs.Parse(args); err != nil {
		return parseErrorCode(err)
	}
	if fs.NArg() != 0 {
		fmt.Fprintf(stderr, "esig version: takes no arguments\n\n")
		fs.Usage()
		return exitUsage
	}

	info, ok := debug.ReadBuildInfo()
	if !ok {
		if version != "" {
			fmt.Fprintf(stdout, "esig %s\n", version)
			return exitOK
		}
		fmt.Fprintln(stdout, "esig: build information unavailable (not built with module support)")
		return exitOK
	}

	v := version
	if v == "" {
		v = moduleVersion(info)
	}
	fmt.Fprintf(stdout, "esig %s\n", v)
	fmt.Fprintf(stdout, "go: %s\n", info.GoVersion)
	return exitOK
}

// moduleVersion returns the version of the dss module esig was built from:
// "(devel)" for a build from a local checkout without a tagged release, or
// the resolved pseudo-version/tag from a `go install
// github.com/ryftcore/dss-go/dss/cmd/esig@version` build. cmd/esig lives inside
// the dss module itself, so it is always info.Main, never an info.Deps
// entry.
func moduleVersion(info *debug.BuildInfo) string {
	if info.Main.Version != "" {
		return info.Main.Version
	}
	return "(unknown)"
}
