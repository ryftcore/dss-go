package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/ryftcore/dss-go/dss/simplereport"
)

// cmdReport implements "esig report".
func cmdReport(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("report", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() {
		fmt.Fprint(stderr, `Usage: esig report <simple-report.xml> -render

Renders a SimpleReport XML file - as "esig validate -format simple" writes
it - as the same human summary "esig validate" prints by default.

  -render
    	required; reserved for future rendering flavours
`)
	}
	var renderFlag bool
	fs.BoolVar(&renderFlag, "render", false, "")
	leading, hadLeading, rest := splitPositional(args)
	if err := fs.Parse(rest); err != nil {
		return exitUsage
	}
	file, ok := resolveOnePositional(fs, hadLeading, leading)
	if !ok {
		fmt.Fprintf(stderr, "esig report: exactly one SimpleReport XML file is required\n\n")
		fs.Usage()
		return exitUsage
	}
	if !renderFlag {
		fmt.Fprintf(stderr, "esig report: -render is required\n\n")
		fs.Usage()
		return exitUsage
	}

	f, err := os.Open(file)
	if err != nil {
		fmt.Fprintf(stderr, "esig report: %v\n", err)
		return exitRuntime
	}
	defer f.Close()

	jaxbReport, err := simplereport.NewFacade().Unmarshal(f)
	if err != nil {
		fmt.Fprintf(stderr, "esig report: parsing %s as a SimpleReport: %v\n", file, err)
		return exitRuntime
	}

	printVerdicts(stdout, simpleReportVerdicts(simplereport.NewSimpleReport(jaxbReport)))
	return exitOK
}
