// Command 03-xades-enveloped-invoice signs a small XML invoice with an
// enveloped XAdES-BASELINE-B signature - the signature lives inside the
// signed document itself, as a ds:Signature element next to the business
// content, which is the usual choice for XML business documents (invoices,
// orders, e-government forms).
//
// Run it from anywhere:
//
//	go run ./examples/03-xades-enveloped-invoice
package main

import (
	"fmt"
	"log"

	"github.com/utain/esig/dss"
	"github.com/utain/esig/dss/examples/internal/fixtures"
)

const invoice = `<?xml version="1.0" encoding="UTF-8"?>
<Invoice>
  <Number>INV-2026-0042</Number>
  <Total currency="EUR">1250.00</Total>
</Invoice>
`

func main() {
	signer, err := dss.OpenPKCS12(fixtures.Path("signer_rsa.p12"), "testpassword")
	if err != nil {
		log.Fatalf("opening key store: %v", err)
	}
	defer signer.Close()

	doc := dss.NewDocument("invoice.xml", []byte(invoice))

	// XAdES defaults to PackagingEnveloped, which is what SignOptions
	// leaves it at here. FormatXAdES also accepts PackagingEnveloping (the
	// signature wraps the content instead of sitting inside it) and
	// PackagingDetached (a separate signature file, see example 04 for the
	// CAdES equivalent).
	signed, err := dss.Sign(doc, signer, dss.SignOptions{
		Format: dss.FormatXAdES,
		Level:  dss.LevelB,
	})
	if err != nil {
		log.Fatalf("signing: %v", err)
	}
	fmt.Println("signed document:", signed.Name())

	reports, err := dss.Validate(signed, dss.ValidateOptions{
		TrustedCertificates: signer.CertificateChain(),
	})
	if err != nil {
		log.Fatalf("validating: %v", err)
	}
	verdict := reports.Verdicts()[0]
	fmt.Printf("verdict: %s, signed by %q, indication=%s\n",
		verdict.SignatureLevel, verdict.SignedBy, verdict.Indication)
}
