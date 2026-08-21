// Command 06-asice-container signs two documents into a single ASiC-E
// container (ETSI EN 319 162) carrying an XAdES signature - the shape to
// reach for when one signature has to cover several files at once, such as
// a payload and its metadata.
//
// Run it from anywhere:
//
//	go run ./examples/06-asice-container
//
// An ASiC container is a zip file under the hood, so this example also opens
// it back up with the standard archive/zip package to show what is inside -
// no ASiC-specific API is needed for that part.
package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
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

	docs := []dss.Document{
		dss.NewDocument("payload.bin", []byte("the payload being signed")),
		dss.NewDocument("metadata.json", []byte(`{"amount":1250.00,"currency":"EUR"}`)),
	}

	// SignMultiple is what Sign calls for a single document too; only the
	// two ASiC formats accept more than one. ContainerType defaults to
	// ContainerASiCE for several documents (ContainerASiCS holds exactly
	// one).
	container, err := dss.SignMultiple(docs, signer, dss.SignOptions{
		Format: dss.FormatASiCWithXAdES,
		Level:  dss.LevelB,
	})
	if err != nil {
		log.Fatalf("signing: %v", err)
	}
	fmt.Println("container:", container.Name())

	stream, err := container.OpenStream()
	if err != nil {
		log.Fatalf("reading the container: %v", err)
	}
	defer stream.Close()
	raw, err := io.ReadAll(stream)
	if err != nil {
		log.Fatalf("reading the container: %v", err)
	}

	zr, err := zip.NewReader(bytes.NewReader(raw), int64(len(raw)))
	if err != nil {
		log.Fatalf("opening the container as a zip archive: %v", err)
	}
	fmt.Println("entries:")
	for _, f := range zr.File {
		fmt.Printf("  %s (%d bytes)\n", f.Name, f.UncompressedSize64)
	}

	reports, err := dss.Validate(container, dss.ValidateOptions{
		TrustedCertificates: signer.CertificateChain(),
	})
	if err != nil {
		log.Fatalf("validating: %v", err)
	}
	verdict := reports.Verdicts()[0]
	fmt.Printf("verdict: %s, indication=%s\n", verdict.SignatureLevel, verdict.Indication)
}
