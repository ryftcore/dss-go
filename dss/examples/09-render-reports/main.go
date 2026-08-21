// Command 09-render-reports signs and validates a document, then walks
// through the four reports Validate returns: the SimpleReport a human reads
// first, the DetailedReport behind every check that led to its verdict, the
// diagnostic data the process reasoned over, and the standardised ETSI TS
// 119 102-2 validation report. dss.Reports embeds the full upstream report
// API, so everything shown here is also reachable without the facade.
//
// Run it from anywhere:
//
//	go run ./examples/09-render-reports
package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/ryftcore/dss-go/dss"
	"github.com/ryftcore/dss-go/dss/examples/internal/fixtures"
)

func main() {
	signer, err := dss.OpenPKCS12(fixtures.Path("signer_rsa.p12"), "testpassword")
	if err != nil {
		log.Fatalf("opening key store: %v", err)
	}
	defer signer.Close()

	doc := dss.NewDocument("invoice.xml", []byte("<invoice><total>42</total></invoice>"))
	signed, err := dss.Sign(doc, signer, dss.SignOptions{Format: dss.FormatXAdES, Level: dss.LevelB})
	if err != nil {
		log.Fatalf("signing: %v", err)
	}

	reports, err := dss.Validate(signed, dss.ValidateOptions{
		TrustedCertificates: signer.CertificateChain(),
	})
	if err != nil {
		log.Fatalf("validating: %v", err)
	}

	// Reports.Verdicts() is the SimpleReport read as Go values - the
	// question "did signature X pass, and why" without walking any XML.
	fmt.Println("== SimpleReport, as Go values ==")
	for _, v := range reports.Verdicts() {
		fmt.Printf("  %s  signed by %q  indication=%s  qualification=%s  valid=%v\n",
			v.SignatureLevel, v.SignedBy, v.Indication, v.Qualification, v.Valid())
	}
	fmt.Printf("%d of %d signatures valid overall: %v\n\n",
		reports.ValidSignatureCount(), reports.SignatureCount(), reports.Valid())

	// The same SimpleReport, as the XML a report viewer or another system
	// consumes.
	simpleXML, err := reports.SimpleReportXML()
	if err != nil {
		log.Fatalf("marshalling the simple report: %v", err)
	}
	fmt.Println("== SimpleReport XML (first line) ==")
	fmt.Println(" ", firstLine(simpleXML))

	// DetailedReport: every EN 319 102-1 building block (format checking,
	// identification of the signing certificate, X.509 certificate
	// validation, cryptographic verification, ...) with its own conclusion -
	// what to read when the SimpleReport's indication needs explaining.
	detailedXML, err := reports.DetailedReportXML()
	if err != nil {
		log.Fatalf("marshalling the detailed report: %v", err)
	}
	fmt.Println("\n== DetailedReport XML ==")
	fmt.Printf("  %d bytes, root element %s\n", len(detailedXML), firstElement(detailedXML))

	// Diagnostic data: the raw facts the process reasoned over - every
	// certificate, revocation datum and signature property it looked at,
	// before any policy was applied to them. This is what a policy
	// constraint (example 08) is actually evaluated against.
	diagnosticXML, err := reports.DiagnosticDataXML()
	if err != nil {
		log.Fatalf("marshalling the diagnostic data: %v", err)
	}
	fmt.Println("\n== DiagnosticData XML ==")
	fmt.Printf("  %d bytes, root element %s\n", len(diagnosticXML), firstElement(diagnosticXML))

	// ETSI TS 119 102-2 validation report: the standardised, machine
	// readable report format other systems interoperate against.
	etsiXML, err := reports.ETSIValidationReportXML()
	if err != nil {
		log.Fatalf("marshalling the ETSI validation report: %v", err)
	}
	fmt.Println("\n== ETSI TS 119 102-2 ValidationReport XML ==")
	fmt.Printf("  %d bytes, root element %s\n", len(etsiXML), firstElement(etsiXML))

	// Everything above went through the facade's convenience wrapper.
	// dss.Reports embeds *reports.Reports, so the underlying JAXB models -
	// typed Go structs, not XML strings - are reachable the same way a
	// caller not using the facade would reach them:
	simpleReport := reports.GetSimpleReport()
	fmt.Printf("\n(underlying API) simplereport.SimpleReport: %d signature(s)\n", simpleReport.GetSignaturesCount())
}

func firstLine(s string) string {
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		return s[:i]
	}
	return s
}

// firstElement returns the name of the XML document's root element, ignoring
// the leading <?xml ...?> declaration when present.
func firstElement(s string) string {
	i := strings.IndexByte(s, '<')
	for i >= 0 && strings.HasPrefix(s[i:], "<?") {
		end := strings.Index(s[i:], "?>")
		if end < 0 {
			break
		}
		i = strings.IndexByte(s[i+end+2:], '<')
		if i < 0 {
			break
		}
		i += end + 2
	}
	if i < 0 {
		return "?"
	}
	rest := s[i+1:]
	end := strings.IndexAny(rest, " \t\n>")
	if end < 0 {
		return rest
	}
	return rest[:end]
}
