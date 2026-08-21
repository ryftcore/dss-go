package main

import (
	"flag"
	"fmt"
	"io"

	"github.com/ryftcore/dss-go/dss"
)

// cmdExtend implements "esig extend".
func cmdExtend(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("extend", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() {
		fmt.Fprint(stderr, `Usage: esig extend <file> -format <format> -level <level> [flags]

Raises every signature in <file> to a higher baseline level. The signing key
is not needed: extension only adds time-stamps and validation data.

  -format string
    	format of the signatures in <file>: cades, xades, pades, jades, asice or asics (required)
  -level string
    	level to reach: T, LT or LTA (required)
  -tsa string
    	RFC 3161 time-stamp authority URL (required: every reachable level is above B)
  -asic-format string
    	for -format asice/asics only: the signature format the container carries, cades or xades (default xades)
  -detached string
    	original document a detached signature covers; repeatable
  -out string
    	output file; defaults to the extended document's own computed name, written to the current directory
`)
	}

	var format, level, tsaURL, asicFormat, out string
	var detached stringList
	fs.StringVar(&format, "format", "", "")
	fs.StringVar(&level, "level", "", "")
	fs.StringVar(&tsaURL, "tsa", "", "")
	fs.StringVar(&asicFormat, "asic-format", "", "")
	fs.Var(&detached, "detached", "")
	fs.StringVar(&out, "out", "", "")
	leading, hadLeading, rest := splitPositional(args)
	if err := fs.Parse(rest); err != nil {
		return exitUsage
	}
	file, ok := resolveOnePositional(fs, hadLeading, leading)
	if !ok {
		fmt.Fprintf(stderr, "esig extend: exactly one input file is required\n\n")
		fs.Usage()
		return exitUsage
	}
	if format == "" || level == "" {
		fmt.Fprintf(stderr, "esig extend: -format and -level are required\n\n")
		fs.Usage()
		return exitUsage
	}

	targetFormat, err := parseExtendFormat(format, asicFormat)
	if err != nil {
		fmt.Fprintf(stderr, "esig extend: %v\n\n", err)
		fs.Usage()
		return exitUsage
	}
	lvl, err := parseLevel(level)
	if err != nil {
		fmt.Fprintf(stderr, "esig extend: %v\n\n", err)
		fs.Usage()
		return exitUsage
	}
	if lvl == dss.LevelB {
		fmt.Fprintf(stderr, "esig extend: -level B is not a valid extension target (nothing is below it)\n\n")
		fs.Usage()
		return exitUsage
	}
	if tsaURL == "" {
		fmt.Fprintf(stderr, "esig extend: -tsa is required (every extension target is above level B)\n\n")
		fs.Usage()
		return exitUsage
	}

	doc, err := dss.OpenDocument(file)
	if err != nil {
		fmt.Fprintf(stderr, "esig extend: %v\n", err)
		return exitRuntime
	}
	var detachedDocs []dss.Document
	if len(detached) > 0 {
		detachedDocs, err = loadDocuments(detached)
		if err != nil {
			fmt.Fprintf(stderr, "esig extend: %v\n", err)
			return exitRuntime
		}
	}

	extended, err := dss.Extend(doc, dss.ExtendOptions{
		Format:           targetFormat,
		Level:            lvl,
		TSPSource:        newHTTPTSPSource(tsaURL),
		DetachedContents: detachedDocs,
	})
	if err != nil {
		fmt.Fprintf(stderr, "esig extend: %v\n", err)
		return exitRuntime
	}

	outPath := out
	if outPath == "" {
		outPath = extended.Name()
	}
	if err := extended.Save(outPath); err != nil {
		fmt.Fprintf(stderr, "esig extend: writing %s: %v\n", outPath, err)
		return exitRuntime
	}
	fmt.Fprintf(stdout, "%s\n", outPath)
	return exitOK
}
