package main

import (
	"flag"
	"fmt"
	"io"
	"runtime/debug"
)

// cmdVersion implements "esig version".
func cmdVersion(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("version", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() {
		fmt.Fprint(stderr, "Usage: esig version\n\nPrints the esig build's module version and Go toolchain version.\n")
	}
	if err := fs.Parse(args); err != nil {
		return exitUsage
	}
	if fs.NArg() != 0 {
		fmt.Fprintf(stderr, "esig version: takes no arguments\n\n")
		fs.Usage()
		return exitUsage
	}

	info, ok := debug.ReadBuildInfo()
	if !ok {
		fmt.Fprintln(stdout, "esig: build information unavailable (not built with module support)")
		return exitOK
	}

	fmt.Fprintf(stdout, "esig %s\n", moduleVersion(info))
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
