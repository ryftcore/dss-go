// Command 04-cades-detached produces a detached CAdES-BASELINE-B signature
// over an arbitrary binary payload: the signature is a separate .p7s file
// that never touches the original content, which is the usual shape for
// signing files you cannot or do not want to modify (archives, media,
// already-published documents).
//
// Run it from anywhere:
//
//	go run ./examples/04-cades-detached
package main

import (
	"fmt"
	"log"

	"github.com/ryftcore/dss-go/dss"
	"github.com/ryftcore/dss-go/dss/examples/internal/fixtures"
)

func main() {
	signer, err := dss.OpenPKCS12(fixtures.Path("signer_rsa.p12"), "testpassword")
	if err != nil {
		log.Fatalf("opening key store: %v", err)
	}
	defer signer.Close()

	content := dss.NewDocument("payload.bin", []byte("this is the content being signed, byte for byte"))

	signature, err := dss.Sign(content, signer, dss.SignOptions{
		Format:    dss.FormatCAdES,
		Level:     dss.LevelB,
		Packaging: dss.PackagingDetached,
	})
	if err != nil {
		log.Fatalf("signing: %v", err)
	}
	fmt.Println("detached signature:", signature.Name())

	// A detached signature carries no copy of what it covers, so validating
	// it needs the original content back - unlike the enveloped and
	// enveloping cases, which carry their own content and need no
	// DetachedContents.
	reports, err := dss.Validate(signature, dss.ValidateOptions{
		DetachedContents:    []dss.Document{content},
		TrustedCertificates: signer.CertificateChain(),
	})
	if err != nil {
		log.Fatalf("validating: %v", err)
	}
	verdict := reports.Verdicts()[0]
	fmt.Printf("verdict: %s, indication=%s\n", verdict.SignatureLevel, verdict.Indication)

	// Validating without DetachedContents finds the signature but cannot
	// check what it covers.
	blind, err := dss.Validate(signature, dss.ValidateOptions{
		TrustedCertificates: signer.CertificateChain(),
	})
	if err != nil {
		log.Fatalf("validating without the detached content: %v", err)
	}
	blindVerdict := blind.Verdicts()[0]
	fmt.Printf("without the original content: indication=%s sub_indication=%s\n",
		blindVerdict.Indication, blindVerdict.SubIndication)
}
