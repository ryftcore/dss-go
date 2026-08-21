package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/utain/esig/dss"
	"github.com/utain/esig/dss/simplereport"
)

// cmdValidate implements "esig validate".
func cmdValidate(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("validate", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() {
		fmt.Fprint(stderr, `Usage: esig validate <file> [flags]

Validates every signature, time-stamp and evidence record in <file> and
prints a summary, or renders one of the DSS reports.

  -detached string
    	original document a detached signature covers; repeatable
  -policy string
    	custom validation policy XML; defaults to the bundled ETSI policy
  -trust string
    	DER or PEM trust anchor certificate; repeatable
  -tl-cache string
    	trusted-list cache directory from a previous "esig tl refresh -cache"
  -at string
    	validation time, RFC 3339 (e.g. 2025-01-15T00:00:00Z); defaults to now
  -format string
    	report to render instead of the human summary: simple, detailed, diagnostic or etsi-vr
  -out string
    	write the rendered report here instead of stdout (ignored without -format)

Exit code: 0 if every signature reached TOTAL_PASSED, 1 otherwise, 2 on a
usage error, 3 if validation could not run at all.
`)
	}

	var policy, tlCache, at, format, out string
	var detached, trust stringList
	fs.Var(&detached, "detached", "")
	fs.StringVar(&policy, "policy", "", "")
	fs.Var(&trust, "trust", "")
	fs.StringVar(&tlCache, "tl-cache", "", "")
	fs.StringVar(&at, "at", "", "")
	fs.StringVar(&format, "format", "", "")
	fs.StringVar(&out, "out", "", "")
	leading, hadLeading, rest := splitPositional(args)
	if err := fs.Parse(rest); err != nil {
		return exitUsage
	}
	file, ok := resolveOnePositional(fs, hadLeading, leading)
	if !ok {
		fmt.Fprintf(stderr, "esig validate: exactly one input file is required\n\n")
		fs.Usage()
		return exitUsage
	}

	opts := dss.ValidateOptions{IncludeSemantics: true}

	if len(detached) > 0 {
		docs, err := loadDocuments(detached)
		if err != nil {
			fmt.Fprintf(stderr, "esig validate: %v\n", err)
			return exitRuntime
		}
		opts.DetachedContents = docs
	}
	if policy != "" {
		doc, err := dss.OpenDocument(policy)
		if err != nil {
			fmt.Fprintf(stderr, "esig validate: opening -policy: %v\n", err)
			return exitRuntime
		}
		opts.Policy = doc
	}
	if len(trust) > 0 {
		certs, err := loadTrustAnchors(trust)
		if err != nil {
			fmt.Fprintf(stderr, "esig validate: %v\n", err)
			return exitRuntime
		}
		opts.TrustedCertificates = certs
	}
	if tlCache != "" {
		certs, err := readTLCache(tlCache)
		if err != nil {
			fmt.Fprintf(stderr, "esig validate: %v\n", err)
			return exitRuntime
		}
		opts.TrustedCertificates = append(opts.TrustedCertificates, certs...)
	}
	if at != "" {
		t, err := time.Parse(time.RFC3339, at)
		if err != nil {
			fmt.Fprintf(stderr, "esig validate: -at: %v\n\n", err)
			fs.Usage()
			return exitUsage
		}
		opts.ValidationTime = &t
	}
	if format != "" {
		if _, err := reportRenderer(format); err != nil {
			fmt.Fprintf(stderr, "esig validate: %v\n\n", err)
			fs.Usage()
			return exitUsage
		}
	}

	doc, err := dss.OpenDocument(file)
	if err != nil {
		fmt.Fprintf(stderr, "esig validate: %v\n", err)
		return exitRuntime
	}

	reports, err := dss.Validate(doc, opts)
	if err != nil {
		fmt.Fprintf(stderr, "esig validate: %v\n", err)
		return exitRuntime
	}

	if format != "" {
		render, _ := reportRenderer(format)
		xml, err := render(reports)
		if err != nil {
			fmt.Fprintf(stderr, "esig validate: rendering %s report: %v\n", format, err)
			return exitRuntime
		}
		if err := writeOut(stdout, out, xml); err != nil {
			fmt.Fprintf(stderr, "esig validate: %v\n", err)
			return exitRuntime
		}
	} else {
		printVerdicts(stdout, reports.Verdicts())
	}

	if !reports.Valid() {
		return exitValidationFailed
	}
	return exitOK
}

// reportRenderer maps a -format value onto the [dss.Reports] method that
// renders it.
func reportRenderer(format string) (func(*dss.Reports) (string, error), error) {
	switch format {
	case "simple":
		return (*dss.Reports).SimpleReportXML, nil
	case "detailed":
		return (*dss.Reports).DetailedReportXML, nil
	case "diagnostic":
		return (*dss.Reports).DiagnosticDataXML, nil
	case "etsi-vr":
		return (*dss.Reports).ETSIValidationReportXML, nil
	default:
		return nil, fmt.Errorf("unsupported -format %q (want simple, detailed, diagnostic or etsi-vr)", format)
	}
}

// writeOut writes content to path, or to stdout when path is empty.
func writeOut(stdout io.Writer, path, content string) error {
	if path == "" {
		_, err := io.WriteString(stdout, content)
		return err
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return fmt.Errorf("writing %s: %w", path, err)
	}
	return nil
}

// simpleReportVerdicts is [dss.Reports.Verdicts] for a bare
// *simplereport.SimpleReport, i.e. one read back from a saved SimpleReport
// XML file by "esig report" rather than produced in-process by
// [dss.Validate]. It reads the same fields, in the same order, that
// [dss.Reports.Verdicts] does.
func simpleReportVerdicts(simple *simplereport.SimpleReport) []dss.Verdict {
	ids := simple.GetSignatureIdList()
	verdicts := make([]dss.Verdict, 0, len(ids))
	for _, id := range ids {
		verdicts = append(verdicts, dss.Verdict{
			ID:                id,
			Indication:        simple.GetIndication(id),
			SubIndication:     simple.GetSubIndication(id),
			SignatureLevel:    simple.GetSignatureFormat(id),
			Qualification:     simple.GetSignatureQualification(id),
			SignedBy:          simple.GetSignedBy(id),
			SigningTime:       simple.GetSigningTime(id),
			BestSignatureTime: simple.GetBestSignatureTime(id),
			Errors:            simple.GetAdESValidationErrors(id),
			Warnings:          simple.GetAdESValidationWarnings(id),
			Infos:             simple.GetAdESValidationInfo(id),
		})
	}
	return verdicts
}

// printVerdicts prints the default human summary: one line per signature.
func printVerdicts(w io.Writer, verdicts []dss.Verdict) {
	if len(verdicts) == 0 {
		fmt.Fprintln(w, "no signatures found")
		return
	}
	for _, v := range verdicts {
		fmt.Fprintf(w, "%s  %s", v.ID, v.Indication)
		if v.SubIndication != "" {
			fmt.Fprintf(w, " (%s)", v.SubIndication)
		}
		fmt.Fprintf(w, "  level=%s  qualification=%s  signed-by=%q\n",
			v.SignatureLevel, v.Qualification, v.SignedBy)
		for _, e := range v.Errors {
			fmt.Fprintf(w, "    error: %s\n", e.Value)
		}
	}
}
