// Command 01-sign-pdf-pades signs a PDF with a PAdES-BASELINE-B signature,
// then raises a second copy to PAdES-BASELINE-T by adding a trusted
// time-stamp over the signature value.
//
// Run it from anywhere:
//
//	go run ./examples/01-sign-pdf-pades
//
// It writes the two signed PDFs to a temporary directory and prints their
// paths; nothing here needs the network.
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/ryftcore/dss-go/dss"
	"github.com/ryftcore/dss-go/dss/examples/internal/fixtures"
	spivalidation "github.com/ryftcore/dss-go/dss/spi/validation"
)

func main() {
	// OpenPKCS12 opens the key store and selects its one key entry. Close it
	// once every signature is produced.
	signer, err := dss.OpenPKCS12(fixtures.Path("signer_rsa.p12"), "testpassword")
	if err != nil {
		log.Fatalf("opening key store: %v", err)
	}
	defer signer.Close()

	doc, err := dss.OpenDocument(fixtures.Path("sample.pdf"))
	if err != nil {
		log.Fatalf("opening document: %v", err)
	}

	outDir, err := os.MkdirTemp("", "dss-example-01-")
	if err != nil {
		log.Fatalf("creating output directory: %v", err)
	}
	fmt.Println("writing signed PDFs to", outDir)

	// Level B: the signature itself, no time-stamp, no revocation data.
	// PAdES ignores SignOptions.Packaging - a PDF signature is always
	// embedded in an incremental update of the document, never detached or
	// enveloping.
	signedB, err := dss.Sign(doc, signer, dss.SignOptions{
		Format: dss.FormatPAdES,
		Level:  dss.LevelB,
	})
	if err != nil {
		log.Fatalf("signing at level B: %v", err)
	}
	if err := signedB.Save(filepath.Join(outDir, signedB.Name())); err != nil {
		log.Fatalf("saving: %v", err)
	}
	fmt.Println("PAdES-BASELINE-B:", signedB.Name())

	// Level T adds a time-stamp over the signature value, so it needs a
	// TSPSource. The dss module ships no HTTP TSA client (see the dss
	// package doc's "Network access" section; the esig CLI carries its own)
	// - only
	// spi/validation.KeyEntityTSPSource, which issues RFC 3161 tokens from a
	// local key. That is exactly what a real deployment does NOT want: here
	// it stands in for a real TSA (a self-hosted one, or a commercial HTTP
	// TSA reached through a TSPSource you implement) so this example needs
	// no network access. Swap it for your TSA's TSPSource and the rest of
	// this call is unchanged.
	tsa, err := spivalidation.NewKeyEntityTSPSourceFromKeyStorePath(
		fixtures.Path("tsa_ec.p12"), "PKCS12", "testpassword", "", "testpassword")
	if err != nil {
		log.Fatalf("building the local test TSA: %v", err)
	}
	tsa.SetTsaPolicy("1.2.3.4.5.6.7.8.9") // unregistered placeholder OID, test fixture only

	signedT, err := dss.Sign(doc, signer, dss.SignOptions{
		Format:    dss.FormatPAdES,
		Level:     dss.LevelT,
		TSPSource: tsa,
	})
	if err != nil {
		log.Fatalf("signing at level T: %v", err)
	}
	if err := signedT.Save(filepath.Join(outDir, signedT.Name())); err != nil {
		log.Fatalf("saving: %v", err)
	}
	fmt.Println("PAdES-BASELINE-T:", signedT.Name())

	// LT and LTA go further still: LT embeds the certificates and
	// revocation data a verifier needs long after the fact (which means a
	// CertificateVerifier carrying CRL/OCSP sources - see
	// dss.SignOptions.CertificateVerifier), and LTA adds an archival
	// time-stamp on top. See example 09 for reading back what a validator
	// makes of a signature like these.
}
