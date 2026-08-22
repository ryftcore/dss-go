package main

import (
	"flag"
	"fmt"
	"io"

	"github.com/ryftcore/dss-go/dss"
	"github.com/ryftcore/dss-go/dss/diagnostic"
)

// cmdInspect implements "esig inspect".
func cmdInspect(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("inspect", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() {
		fmt.Fprint(stderr, `Usage: esig inspect <file> [flags]

Prints a signature/timestamp/certificate summary read from <file>'s
diagnostic data - the facts the validation process reasoned over, gathered
without regard to whether any signature actually validates. Use "esig
validate" for a verdict.

  -detached string
    	original document a detached signature covers; repeatable
`)
	}
	var detached stringList
	fs.Var(&detached, "detached", "")
	leading, hadLeading, rest := splitPositional(args)
	if err := fs.Parse(rest); err != nil {
		return exitUsage
	}
	file, ok := resolveOnePositional(fs, hadLeading, leading)
	if !ok {
		fmt.Fprintf(stderr, "esig inspect: exactly one input file is required\n\n")
		fs.Usage()
		return exitUsage
	}

	doc, err := dss.OpenDocument(file)
	if err != nil {
		fmt.Fprintf(stderr, "esig inspect: %v\n", err)
		return exitRuntime
	}
	opts := dss.ValidateOptions{}
	if len(detached) > 0 {
		docs, err := loadDocuments(detached)
		if err != nil {
			fmt.Fprintf(stderr, "esig inspect: %v\n", err)
			return exitRuntime
		}
		opts.DetachedContents = docs
	}

	reports, err := dss.Validate(doc, opts)
	if err != nil {
		fmt.Fprintf(stderr, "esig inspect: %v\n", err)
		return exitRuntime
	}
	data := reports.GetDiagnosticData()

	fmt.Fprintf(stdout, "document: %s\n", data.DocumentName())

	signatureIDs := data.SignatureIdList()
	fmt.Fprintf(stdout, "signatures: %d\n", len(signatureIDs))
	for _, id := range signatureIDs {
		fmt.Fprintf(stdout, "  %s\n", id)
		fmt.Fprintf(stdout, "    format: %s\n", data.SignatureFormat(id))
		if t := data.SignatureDate(id); t != nil {
			fmt.Fprintf(stdout, "    claimed signing time: %s\n", t.Format(rfc3339Display))
		}
		fmt.Fprintf(stdout, "    digest algorithm: %s\n", data.SignatureDigestAlgorithm(id))
		if certID := data.SigningCertificateId(id); certID != "" {
			fmt.Fprintf(stdout, "    signing certificate: %s\n", certificateSummary(data, certID))
		} else {
			fmt.Fprintf(stdout, "    signing certificate: not identified\n")
		}
		for _, tsID := range data.TimestampIdListForSignature(id) {
			fmt.Fprintf(stdout, "    timestamp %s: %s\n", tsID, timestampSummary(data, tsID))
		}
	}

	timestampIDs := data.TimestampIdList()
	fmt.Fprintf(stdout, "detached timestamps: %d\n", len(timestampIDs))
	for _, id := range timestampIDs {
		fmt.Fprintf(stdout, "  %s: %s\n", id, timestampSummary(data, id))
	}

	certs := data.UsedCertificates()
	fmt.Fprintf(stdout, "certificates used: %d\n", len(certs))
	for _, cert := range certs {
		fmt.Fprintf(stdout, "  %s: %s\n", cert.Id(), cert.CertificateDN())
	}

	return exitOK
}

// rfc3339Display is the timestamp layout inspect prints times in.
const rfc3339Display = "2006-01-02T15:04:05Z07:00"

// certificateSummary renders a one-line description of a certificate by its
// diagnostic-data id.
func certificateSummary(data *diagnostic.Data, certID string) string {
	cert := data.UsedCertificateByIdNullSafe(certID)
	if cert == nil {
		return certID
	}
	return fmt.Sprintf("%s (serial %s)", cert.CertificateDN(), cert.SerialNumber())
}

// timestampSummary renders a one-line description of a timestamp by its
// diagnostic-data id.
func timestampSummary(data *diagnostic.Data, tsID string) string {
	ts := data.TimestampById(tsID)
	if ts == nil {
		return tsID
	}
	summary := string(ts.Type())
	if t := ts.ProductionTime(); t != nil {
		summary += " produced " + t.Format(rfc3339Display)
	}
	return summary
}
