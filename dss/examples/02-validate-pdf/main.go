// Command 02-validate-pdf signs a PDF and then validates it twice: once with
// no trust anchor configured, and once trusting the signer's certificate.
// The point of the example is the difference between the two - an
// unanchored chain is a correct INDETERMINATE, not a bug.
//
// Run it from anywhere:
//
//	go run ./examples/02-validate-pdf
package main

import (
	"fmt"
	"log"

	"github.com/utain/esig/dss"
	"github.com/utain/esig/dss/examples/internal/fixtures"
)

func main() {
	signer, err := dss.OpenPKCS12(fixtures.Path("signer_rsa.p12"), "testpassword")
	if err != nil {
		log.Fatalf("opening key store: %v", err)
	}
	defer signer.Close()

	doc, err := dss.OpenDocument(fixtures.Path("sample.pdf"))
	if err != nil {
		log.Fatalf("opening document: %v", err)
	}

	signed, err := dss.Sign(doc, signer, dss.SignOptions{
		Format: dss.FormatPAdES,
		Level:  dss.LevelB,
	})
	if err != nil {
		log.Fatalf("signing: %v", err)
	}

	// Validate auto-detects the format - PDF here - from the document
	// itself, so it is never named explicitly.
	untrusted, err := dss.Validate(signed, dss.ValidateOptions{})
	if err != nil {
		log.Fatalf("validating (untrusted): %v", err)
	}
	printVerdict("no trust anchor configured", untrusted)

	// Trusting the certificate that produced the signature - which is what
	// a real deployment does with its own root/intermediate CAs, or with
	// the certificate source a TSL validation job produces (example 07) -
	// changes the verdict.
	trusted, err := dss.Validate(signed, dss.ValidateOptions{
		TrustedCertificates: signer.CertificateChain(),
	})
	if err != nil {
		log.Fatalf("validating (trusted): %v", err)
	}
	printVerdict("signer certificate trusted", trusted)
}

func printVerdict(label string, reports *dss.Reports) {
	fmt.Println(label + ":")
	for _, v := range reports.Verdicts() {
		fmt.Printf("  %s  indication=%s", v.SignatureLevel, v.Indication)
		if v.SubIndication != "" {
			fmt.Printf(" sub_indication=%s", v.SubIndication)
		}
		fmt.Println()
		for _, e := range v.Errors {
			fmt.Println("    error:", e.Value)
		}
	}
}
