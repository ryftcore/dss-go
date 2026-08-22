// Command esig is a command-line front end for the dss Go module
// (github.com/ryftcore/dss-go/dss): sign, extend, validate and inspect the
// signature formats the library supports, render the reports [dss.Validate]
// produces, and refresh a local trusted-list cache.
//
// It is a thin wrapper: every subcommand delegates to the dss package facade
// (github.com/ryftcore/dss-go/dss) and is a living example of that API. Advanced
// use that the facade does not cover is reached through the underlying
// packages directly - see the dss package doc.
//
// Only the standard library is used: argument parsing is [flag], nothing
// else. See each subcommand's -h output, or [usage], for its exact flags.
//
// # Exit codes
//
//   - 0: the command completed and, for "validate", every signature was
//     TOTAL_PASSED.
//   - 1: "validate" ran to completion but at least one signature did not
//     reach TOTAL_PASSED (or the document carried none).
//   - 2: the command line could not be parsed (unknown flag, missing
//     required flag, wrong number of arguments).
//   - 3: the command was understood but failed while running (an unreadable
//     file, a signing or validation error, a network failure).
package main

import (
	"fmt"
	"io"
	"os"
)

// Exit codes. See the package doc.
const (
	exitOK               = 0
	exitValidationFailed = 1
	exitUsage            = 2
	exitRuntime          = 3
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

// run dispatches to the named subcommand and returns the process exit code.
// It never calls os.Exit itself, so it can be driven directly from tests.
func run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		usage(stderr)
		return exitUsage
	}

	cmd, rest := args[0], args[1:]
	switch cmd {
	case "validate":
		return cmdValidate(rest, stdout, stderr)
	case "sign":
		return cmdSign(rest, stdout, stderr)
	case "extend":
		return cmdExtend(rest, stdout, stderr)
	case "report":
		return cmdReport(rest, stdout, stderr)
	case "inspect":
		return cmdInspect(rest, stdout, stderr)
	case "tl":
		return cmdTL(rest, stdout, stderr)
	case "version":
		return cmdVersion(rest, stdout, stderr)
	case "-h", "-help", "--help", "help":
		usage(stdout)
		return exitOK
	default:
		fmt.Fprintf(stderr, "esig: unknown command %q\n\n", cmd)
		usage(stderr)
		return exitUsage
	}
}

// usage prints the top-level command summary.
func usage(w io.Writer) {
	fmt.Fprint(w, `esig is a command-line front end for the dss digital signature library.

Usage:

	esig <command> [arguments]

Commands:

	validate   validate a signed document and print or render its reports
	sign       sign a document, producing a new signed document
	extend     raise an existing signature to a higher baseline level
	report     render a previously saved SimpleReport XML file
	inspect    print a signature/timestamp/certificate summary
	tl         refresh a local trusted-list (TSL/LOTL) cache
	version    print the esig and module build information

Run "esig <command> -h" for a command's flags.
`)
}
