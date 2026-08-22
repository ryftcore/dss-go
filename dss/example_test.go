package dss_test

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ryftcore/dss-go/dss"
	spivalidation "github.com/ryftcore/dss-go/dss/spi/validation"
	"github.com/ryftcore/dss-go/dss/token"
)

// The examples sign and validate in the same process, using the small
// fixtures under testdata (see testdata/README.md): a self-signed RSA test key
// in a PKCS#12 store, a minimal PDF, and an EC key acting as a local RFC 3161
// time-stamp authority so that nothing here needs the network.

// exampleSigner opens the test key store the examples sign with.
func exampleSigner() *dss.Signer {
	signer, err := dss.OpenPKCS12("testdata/signer_rsa.p12", "testpassword")
	if err != nil {
		log.Fatal(err)
	}
	return signer
}

// exampleTSA turns the test EC key into a local time-stamp authority. A real
// deployment points a TSPSource at its TSA instead; see the package doc on
// what the port ships.
func exampleTSA() dss.TSPSource {
	tsa, err := spivalidation.NewKeyEntityTSPSourceFromKeyStorePath(
		"testdata/tsa_ec.p12", "PKCS12", "testpassword", "", "testpassword")
	if err != nil {
		log.Fatal(err)
	}
	tsa.SetTsaPolicy("1.2.3.4.5.6.7.8.9")
	return tsa
}

// exampleSignedXML signs a small XML invoice at level B, for the examples that
// are about the reports rather than about signing.
func exampleSignedXML() (dss.Document, *dss.Signer) {
	signer := exampleSigner()
	doc := dss.NewDocument("invoice.xml", []byte("<invoice><total>42</total></invoice>"))
	signed, err := dss.Sign(doc, signer, dss.SignOptions{Format: dss.FormatXAdES, Level: dss.LevelB})
	if err != nil {
		log.Fatal(err)
	}
	return signed, signer
}

// Sign an XML document with an enveloped XAdES-BASELINE-B signature.
func ExampleSign() {
	signer, err := dss.OpenPKCS12("testdata/signer_rsa.p12", "testpassword")
	if err != nil {
		log.Fatal(err)
	}
	defer signer.Close()

	doc := dss.NewDocument("invoice.xml", []byte("<invoice><total>42</total></invoice>"))

	signed, err := dss.Sign(doc, signer, dss.SignOptions{
		Format: dss.FormatXAdES,
		Level:  dss.LevelB,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(signed.Name())
	// Output: invoice-signed-xades-baseline-b.xml
}

// Sign a PDF. PAdES ignores SignOptions.Packaging: a PDF signature is always
// embedded in an incremental update of the document itself.
func ExampleSign_pdf() {
	signer := exampleSigner()
	defer signer.Close()

	doc, err := dss.OpenDocument("testdata/sample.pdf")
	if err != nil {
		log.Fatal(err)
	}

	signed, err := dss.Sign(doc, signer, dss.SignOptions{
		Format: dss.FormatPAdES,
		Level:  dss.LevelB,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(signed.Name())
	// Output: sample-signed-pades-baseline-b.pdf
}

// Sign detached: the signature is a separate file, and validating it later
// needs the original content back through ValidateOptions.DetachedContents.
func ExampleSign_detached() {
	signer := exampleSigner()
	defer signer.Close()

	content := dss.NewDocument("payload.bin", []byte("payload"))

	signature, err := dss.Sign(content, signer, dss.SignOptions{
		Format:    dss.FormatCAdES,
		Level:     dss.LevelB,
		Packaging: dss.PackagingDetached,
	})
	if err != nil {
		log.Fatal(err)
	}

	reports, err := dss.Validate(signature, dss.ValidateOptions{
		DetachedContents:    []dss.Document{content},
		TrustedCertificates: signer.CertificateChain(),
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(signature.Name(), reports.Verdicts()[0].Indication)
	// Output: payload-signed-cades-baseline-b.p7s TOTAL_PASSED
}

// Sign at level T, which adds a time-stamp over the signature value and
// therefore needs a TSPSource.
func ExampleSign_timestamped() {
	signer := exampleSigner()
	defer signer.Close()

	doc := dss.NewDocument("payload.bin", []byte("payload"))

	signed, err := dss.Sign(doc, signer, dss.SignOptions{
		Format:    dss.FormatCAdES,
		Level:     dss.LevelT,
		TSPSource: exampleTSA(),
	})
	if err != nil {
		log.Fatal(err)
	}

	reports, err := dss.Validate(signed, dss.ValidateOptions{
		TrustedCertificates: signer.CertificateChain(),
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(reports.Verdicts()[0].SignatureLevel)
	// Output: CAdES-BASELINE-T
}

// Package several documents into one signed ASiC-E container.
func ExampleSignMultiple() {
	signer := exampleSigner()
	defer signer.Close()

	docs := []dss.Document{
		dss.NewDocument("payload.bin", []byte("payload")),
		dss.NewDocument("metadata.json", []byte(`{"amount":42}`)),
	}

	container, err := dss.SignMultiple(docs, signer, dss.SignOptions{
		Format:        dss.FormatASiCWithXAdES,
		Level:         dss.LevelB,
		ContainerType: dss.ContainerASiCE,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(container.Name())
	// Output: container-signed-xades-baseline-b.sce
}

// Raise an existing B-level signature to T by adding a time-stamp. The signing
// key is not needed: extension never touches the signature value.
func ExampleExtend() {
	signer := exampleSigner()
	defer signer.Close()

	doc := dss.NewDocument("payload.bin", []byte("payload"))
	signed, err := dss.Sign(doc, signer, dss.SignOptions{Format: dss.FormatCAdES, Level: dss.LevelB})
	if err != nil {
		log.Fatal(err)
	}

	extended, err := dss.Extend(signed, dss.ExtendOptions{
		Format:    dss.FormatCAdES,
		Level:     dss.LevelT,
		TSPSource: exampleTSA(),
	})
	if err != nil {
		log.Fatal(err)
	}

	reports, err := dss.Validate(extended, dss.ValidateOptions{
		TrustedCertificates: signer.CertificateChain(),
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(reports.Verdicts()[0].SignatureLevel)
	// Output: CAdES-BASELINE-T
}

// Validate a signed document. With no trust anchor the chain reaches nothing
// trusted, so the verdict is INDETERMINATE rather than a pass - which is the
// correct answer, not a failure of the library.
func ExampleValidate() {
	signed, signer := exampleSignedXML()
	defer signer.Close()

	reports, err := dss.Validate(signed, dss.ValidateOptions{})
	if err != nil {
		log.Fatal(err)
	}
	verdict := reports.Verdicts()[0]
	fmt.Println(verdict.Indication, verdict.SubIndication)
	// Output: INDETERMINATE NO_CERTIFICATE_CHAIN_FOUND
}

// Validate with trust anchors: the same document, now anchored, passes.
func ExampleValidate_trusted() {
	signed, signer := exampleSignedXML()
	defer signer.Close()

	reports, err := dss.Validate(signed, dss.ValidateOptions{
		TrustedCertificateSources: []dss.CertificateSource{
			dss.TrustStore(signer.CertificateChain()...),
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(reports.Verdicts()[0].Indication)
	// Output: TOTAL_PASSED
}

// OpenDocument reads a document from disk, lazily.
func ExampleOpenDocument() {
	doc, err := dss.OpenDocument("testdata/sample.pdf")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(doc.Name(), doc.MimeType())
	// Output: sample.pdf PDF
}

// NewDocument wraps bytes already in memory. The name travels with the
// document into container entries and into the reports.
func ExampleNewDocument() {
	doc := dss.NewDocument("invoice.xml", []byte("<invoice/>"))
	fmt.Println(doc.Name())
	// Output: invoice.xml
}

// OpenPKCS12 opens a .p12 or .pfx key store and selects its first key entry.
func ExampleOpenPKCS12() {
	signer, err := dss.OpenPKCS12("testdata/signer_rsa.p12", "testpassword")
	if err != nil {
		log.Fatal(err)
	}
	defer signer.Close()

	fmt.Println(signer.Certificate().Certificate().Subject.CommonName)
	// Output: Go Port Test RSA
}

// OpenPKCS12Bytes is the same for a key store already in memory.
func ExampleOpenPKCS12Bytes() {
	store, err := os.ReadFile("testdata/signer_rsa.p12")
	if err != nil {
		log.Fatal(err)
	}
	signer, err := dss.OpenPKCS12Bytes(store, "testpassword")
	if err != nil {
		log.Fatal(err)
	}
	defer signer.Close()

	fmt.Println(len(signer.CertificateChain()))
	// Output: 1
}

// NewSigner pairs any token connection with the key entry to sign with - a
// smart card, an HSM, or, as here, a key store the caller opened itself and
// keeps ownership of.
func ExampleNewSigner() {
	connection, err := token.NewPkcs12SignatureTokenFromFilepath(
		"testdata/signer_rsa.p12", token.NewPasswordProtection([]byte("testpassword")))
	if err != nil {
		log.Fatal(err)
	}
	defer connection.Close()

	keys, err := connection.Keys()
	if err != nil {
		log.Fatal(err)
	}
	signer, err := dss.NewSigner(connection, keys[0])
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(signer.KeyEntry().EncryptionAlgorithm())
	// Output: RSA
}

// Certificate returns the certificate the signature will name as its signing
// certificate.
func ExampleSigner_Certificate() {
	signer := exampleSigner()
	defer signer.Close()

	fmt.Println(signer.Certificate().IsSelfSigned())
	// Output: true
}

// CertificateChain returns the chain the key store carries, which is what gets
// embedded in the signature.
func ExampleSigner_CertificateChain() {
	signer := exampleSigner()
	defer signer.Close()

	for _, certificate := range signer.CertificateChain() {
		fmt.Println(certificate.Certificate().Subject.CommonName)
	}
	// Output: Go Port Test RSA
}

// KeyEntry exposes the underlying key entry for code that drives the ported
// services directly.
func ExampleSigner_KeyEntry() {
	signer := exampleSigner()
	defer signer.Close()

	fmt.Println(signer.KeyEntry().EncryptionAlgorithm())
	// Output: RSA
}

// Token exposes the underlying connection, for instance to list the other key
// entries of the same store.
func ExampleSigner_Token() {
	signer := exampleSigner()
	defer signer.Close()

	keys, err := signer.Token().Keys()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(len(keys))
	// Output: 1
}

// Close releases a key store this package opened. It is a no-op for a Signer
// built with NewSigner.
func ExampleSigner_Close() {
	signer, err := dss.OpenPKCS12("testdata/signer_rsa.p12", "testpassword")
	if err != nil {
		log.Fatal(err)
	}
	signer.Close()
	fmt.Println("closed")
	// Output: closed
}

// BaselineLevel resolves a facade format and level onto the ETSI signature
// level the underlying services take.
func ExampleFormat_BaselineLevel() {
	level, err := dss.FormatXAdES.BaselineLevel(dss.LevelLTA)
	if err != nil {
		log.Fatal(err)
	}
	// String renders the dash spelling, as Java's toString() does; the
	// underlying value is the Java enum name.
	fmt.Println(level, string(level))

	// An ASiC container carries the levels of the signature format inside it.
	level, err = dss.FormatASiCWithCAdES.BaselineLevel(dss.LevelT)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(level)
	// Output:
	// XAdES-BASELINE-LTA XAdES_BASELINE_LTA
	// CAdES-BASELINE-T
}

// IsContainer tells the two ASiC formats - the only ones that can cover
// several documents with one signature - from the rest.
func ExampleFormat_IsContainer() {
	fmt.Println(dss.FormatASiCWithCAdES.IsContainer(), dss.FormatPAdES.IsContainer())
	// Output: true false
}

// String returns the format name.
func ExampleFormat_String() {
	fmt.Println(dss.FormatJAdES.String())
	// Output: JAdES
}

// NeedsTimestamp reports which levels require a TSPSource.
func ExampleLevel_NeedsTimestamp() {
	for _, level := range []dss.Level{dss.LevelB, dss.LevelT, dss.LevelLT, dss.LevelLTA} {
		fmt.Println(level, level.NeedsTimestamp())
	}
	// Output:
	// B false
	// T true
	// LT true
	// LTA true
}

// String returns the level name.
func ExampleLevel_String() {
	fmt.Println(dss.LevelLTA.String())
	// Output: LTA
}

// LoadCertificate reads a DER or PEM encoded certificate from disk, ready to
// be used as a trust anchor.
func ExampleLoadCertificate() {
	certificate, err := dss.LoadCertificate("testdata/signer_rsa.cer")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(certificate.Certificate().Subject.CommonName)
	// Output: Go Port Test RSA
}

// LoadCertificateBytes is the same for an encoding already in memory.
func ExampleLoadCertificateBytes() {
	der, err := os.ReadFile("testdata/signer_rsa.cer")
	if err != nil {
		log.Fatal(err)
	}
	certificate, err := dss.LoadCertificateBytes(der)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(certificate.IsSelfSigned())
	// Output: true
}

// TrustStore turns a set of certificates into trust anchors for validation.
func ExampleTrustStore() {
	certificate, err := dss.LoadCertificate("testdata/signer_rsa.cer")
	if err != nil {
		log.Fatal(err)
	}
	signed, signer := exampleSignedXML()
	defer signer.Close()

	reports, err := dss.Validate(signed, dss.ValidateOptions{
		TrustedCertificateSources: []dss.CertificateSource{dss.TrustStore(certificate)},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(reports.Valid())
	// Output: true
}

// Verdicts is the short answer for each signature in the document.
func ExampleReports_Verdicts() {
	signed, signer := exampleSignedXML()
	defer signer.Close()

	reports, err := dss.Validate(signed, dss.ValidateOptions{
		TrustedCertificates: signer.CertificateChain(),
	})
	if err != nil {
		log.Fatal(err)
	}
	for _, verdict := range reports.Verdicts() {
		fmt.Println(verdict.SignatureLevel, verdict.Indication, verdict.SignedBy, verdict.Valid())
	}
	// Output: XAdES-BASELINE-B TOTAL_PASSED Go Port Test RSA true
}

// TimestampVerdicts covers the detached time-stamp tokens a document carries;
// an ordinary signed document carries none, because its time-stamps live
// inside the signature.
func ExampleReports_TimestampVerdicts() {
	signed, signer := exampleSignedXML()
	defer signer.Close()

	reports, err := dss.Validate(signed, dss.ValidateOptions{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(len(reports.TimestampVerdicts()))
	// Output: 0
}

// Valid is the single-boolean answer: every signature passed, and there was at
// least one.
func ExampleReports_Valid() {
	signed, signer := exampleSignedXML()
	defer signer.Close()

	untrusted, err := dss.Validate(signed, dss.ValidateOptions{})
	if err != nil {
		log.Fatal(err)
	}
	trusted, err := dss.Validate(signed, dss.ValidateOptions{
		TrustedCertificates: signer.CertificateChain(),
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(untrusted.Valid(), trusted.Valid())
	// Output: false true
}

// SignatureCount and ValidSignatureCount count what was found and what passed.
func ExampleReports_SignatureCount() {
	signed, signer := exampleSignedXML()
	defer signer.Close()

	reports, err := dss.Validate(signed, dss.ValidateOptions{
		TrustedCertificates: signer.CertificateChain(),
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(reports.SignatureCount(), reports.ValidSignatureCount())
	// Output: 1 1
}

// ValidSignatureCount counts the signatures that reached TOTAL_PASSED.
func ExampleReports_ValidSignatureCount() {
	signed, signer := exampleSignedXML()
	defer signer.Close()

	reports, err := dss.Validate(signed, dss.ValidateOptions{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(reports.ValidSignatureCount())
	// Output: 0
}

// SimpleReportXML is the report an operator reads. The examples print only the
// root element, because the reports carry the validation time.
func ExampleReports_SimpleReportXML() {
	signed, signer := exampleSignedXML()
	defer signer.Close()

	reports, err := dss.Validate(signed, dss.ValidateOptions{})
	if err != nil {
		log.Fatal(err)
	}
	xml, err := reports.SimpleReportXML()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(strings.Contains(xml, "<SimpleReport"))
	// Output: true
}

// DetailedReportXML carries every EN 319 102-1 building block and check.
func ExampleReports_DetailedReportXML() {
	signed, signer := exampleSignedXML()
	defer signer.Close()

	reports, err := dss.Validate(signed, dss.ValidateOptions{})
	if err != nil {
		log.Fatal(err)
	}
	xml, err := reports.DetailedReportXML()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(strings.Contains(xml, "<DetailedReport"))
	// Output: true
}

// DiagnosticDataXML carries the raw facts the process reasoned over.
func ExampleReports_DiagnosticDataXML() {
	signed, signer := exampleSignedXML()
	defer signer.Close()

	reports, err := dss.Validate(signed, dss.ValidateOptions{})
	if err != nil {
		log.Fatal(err)
	}
	xml, err := reports.DiagnosticDataXML()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(strings.Contains(xml, "<DiagnosticData"))
	// Output: true
}

// ETSIValidationReportXML is the standardised, machine-readable report.
func ExampleReports_ETSIValidationReportXML() {
	signed, signer := exampleSignedXML()
	defer signer.Close()

	reports, err := dss.Validate(signed, dss.ValidateOptions{})
	if err != nil {
		log.Fatal(err)
	}
	xml, err := reports.ETSIValidationReportXML()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(strings.Contains(xml, "ValidationReport"))
	// Output: true
}

// Valid reports whether one signature reached TOTAL_PASSED.
func ExampleVerdict_Valid() {
	signed, signer := exampleSignedXML()
	defer signer.Close()

	reports, err := dss.Validate(signed, dss.ValidateOptions{
		TrustedCertificates: signer.CertificateChain(),
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(reports.Verdicts()[0].Valid())
	// Output: true
}
