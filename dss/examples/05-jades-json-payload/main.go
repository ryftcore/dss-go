// Command 05-jades-json-payload signs a JSON payload with a compact
// JAdES-BASELINE-B signature (ETSI TS 119 182) - the JOSE-family format for
// JSON APIs, where the signature travels as a single base64url string
// instead of an XML or CMS structure.
//
// Run it from anywhere:
//
//	go run ./examples/05-jades-json-payload
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

	payload := dss.NewDocument("order.json", []byte(`{"orderId":"ORD-42","amount":1250.00,"currency":"EUR"}`))

	// JWSSerialization defaults to JWSCompact, the single-string
	// "header.payload.signature" form. JWSJSON and JWSFlattenedJSON are the
	// other two ETSI TS 119 182 serializations; a detached or a
	// multi-signature JAdES needs one of them instead.
	signed, err := dss.Sign(payload, signer, dss.SignOptions{
		Format: dss.FormatJAdES,
		Level:  dss.LevelB,
	})
	if err != nil {
		log.Fatalf("signing: %v", err)
	}

	body, err := signed.OpenStream()
	if err != nil {
		log.Fatalf("reading the signed JWS: %v", err)
	}
	defer body.Close()
	buf := make([]byte, 4096)
	n, _ := body.Read(buf)
	compact := string(buf[:n])
	fmt.Printf("compact JWS (%d bytes): %.80s...\n", n, compact)

	reports, err := dss.Validate(signed, dss.ValidateOptions{
		TrustedCertificates: signer.CertificateChain(),
	})
	if err != nil {
		log.Fatalf("validating: %v", err)
	}
	verdict := reports.Verdicts()[0]
	fmt.Printf("verdict: %s, indication=%s\n", verdict.SignatureLevel, verdict.Indication)
}
