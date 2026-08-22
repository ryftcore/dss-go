// Command crossgen is the GO -> UPSTREAM direction of the cross-validation harness (task #12,
// XAdES extension): it signs a fixed sample document with this package's own XAdESService,
// producing XAdES-B enveloped, XAdES-B enveloping, XAdES-B detached, and XAdES-T enveloping
// signatures with real crypto (an RSA PKCS#12 test key for the signer, an EC PKCS#12 test key as
// a self-hosted TSA via spi/validation.KeyEntityTSPSource - the same two key stores
// cades/testdata/crossgen already uses, copied here rather than shared across packages since
// go.mod has no notion of a "testdata" import), and writes them to files. The point is a real,
// independently-verifiable artifact: xades_downstream_cross_validation_test.go runs this program
// and then hands its output to CrossGenValidator.java (an unmodified copy of
// cades/testdata/crossgen's - it drives upstream DSS's own generic SignedDocumentValidator/
// SignedDocumentDiagnosticDataBuilder, so nothing in it is CAdES-specific), which asserts the
// signature is intact, the signing certificate is identified, and the level is recognized -
// upstream DSS accepting what this port produced is the actual proof of compatibility, not
// another Go-side assertion.
//
// Usage: go run . <output directory>
//
// Writes <outdir>/xades-b-enveloped.xml, <outdir>/xades-b-enveloping.xml,
// <outdir>/xades-t-enveloping.xml, and <outdir>/xades-b-detached.xml (+ its companion
// <outdir>/xades-b-detached-content.xml).
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi/validation"
	"github.com/ryftcore/dss-go/dss/token"
	"github.com/ryftcore/dss-go/dss/xades"
)

// sampleXMLContent is the fixed payload every generated signature covers. Its own text records
// what generated it, so a file found on disk explains itself. It has to be well-formed XML (not
// just any bytes) for SignaturePackagingEnveloped, which embeds the ds:Signature as a child of
// this document's own root element.
const sampleXMLContent = `<?xml version="1.0" encoding="UTF-8"?>` +
	`<SampleDocument xmlns="urn:dss:go:crossgen">` +
	`<Note>DSS Go port cross-validation sample content - github.com/ryftcore/dss-go/dss xades/testdata/crossgen. ` +
	`Signed by the Go port's own XAdESService, verified by upstream DSS 6.5.RC1's SignedDocumentValidator.</Note>` +
	`</SampleDocument>`

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: crossgen <output directory>")
		os.Exit(2)
	}
	outDir := os.Args[1]
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		fail(err)
	}

	selfDir, err := os.Getwd()
	if err != nil {
		fail(err)
	}

	signerEntry, err := loadKeyEntry(filepath.Join(selfDir, "signer_rsa.p12"), "testpassword")
	if err != nil {
		fail(fmt.Errorf("loading signer key: %w", err))
	}

	tspSource, err := validation.NewKeyEntityTSPSourceFromKeyStorePath(
		filepath.Join(selfDir, "tsa_ec.p12"), "PKCS12", "testpassword", "", "testpassword")
	if err != nil {
		fail(fmt.Errorf("loading self-TSA key: %w", err))
	}
	// A self-hosted TSA still has to advertise a policy OID (RFC 3161's TSTInfo.policy is
	// mandatory); this one is an arbitrary, unregistered test OID, exactly the role
	// "1.2.3.4.5.6.7.8.9" style placeholders play in DSS's own KeyEntityTSPSource unit tests.
	tspSource.SetTsaPolicy("1.2.3.4.5.6.7.8.9")

	if err := generate(outDir, "xades-b-enveloped.xml", enumerations.SignatureLevelXAdESBaselineB,
		enumerations.SignaturePackagingEnveloped, signerEntry, nil); err != nil {
		fail(fmt.Errorf("generating XAdES-B enveloped: %w", err))
	}
	if err := generate(outDir, "xades-b-enveloping.xml", enumerations.SignatureLevelXAdESBaselineB,
		enumerations.SignaturePackagingEnveloping, signerEntry, nil); err != nil {
		fail(fmt.Errorf("generating XAdES-B enveloping: %w", err))
	}
	if err := generate(outDir, "xades-t-enveloping.xml", enumerations.SignatureLevelXAdESBaselineT,
		enumerations.SignaturePackagingEnveloping, signerEntry, tspSource); err != nil {
		fail(fmt.Errorf("generating XAdES-T enveloping: %w", err))
	}
	if err := generateDetached(outDir, signerEntry); err != nil {
		fail(fmt.Errorf("generating XAdES-B detached: %w", err))
	}

	fmt.Println("wrote", filepath.Join(outDir, "xades-b-enveloped.xml"))
	fmt.Println("wrote", filepath.Join(outDir, "xades-b-enveloping.xml"))
	fmt.Println("wrote", filepath.Join(outDir, "xades-t-enveloping.xml"))
	fmt.Println("wrote", filepath.Join(outDir, "xades-b-detached.xml"))
	fmt.Println("wrote", filepath.Join(outDir, "xades-b-detached-content.xml"))
}

// loadKeyEntry opens a PKCS#12 test key store and returns its single key entry.
func loadKeyEntry(path, password string) (token.DSSPrivateKeyEntry, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	signatureToken, err := token.NewPkcs12SignatureTokenFromBytes(data, token.NewPasswordProtection([]byte(password)))
	if err != nil {
		return nil, err
	}
	keys, err := signatureToken.Keys()
	if err != nil {
		return nil, err
	}
	if len(keys) == 0 {
		return nil, fmt.Errorf("%s carries no key entries", path)
	}
	entry, ok := keys[0].(*token.KSPrivateKeyEntry)
	if !ok {
		return nil, fmt.Errorf("%s: unexpected key entry type %T", path, keys[0])
	}
	return signatureToken.KeyWithPassword(entry.Alias(), token.NewPasswordProtection([]byte(password)))
}

// newParameters builds the XAdESSignatureParameters shared by every generated signature: the
// digest algorithm, signing certificate/chain, and (for anything above baseline B) the TSP source
// a XAdES-T needs to request its signature-timestamp.
func newParameters(level enumerations.SignatureLevel, packaging enumerations.SignaturePackaging,
	signerEntry token.DSSPrivateKeyEntry) *xades.XAdESSignatureParameters {
	parameters := xades.NewXAdESSignatureParameters()
	parameters.SetSignatureLevel(level)
	parameters.SetSignaturePackaging(packaging)
	parameters.SetDigestAlgorithm(enumerations.DigestAlgorithmSHA256)
	parameters.SetSigningCertificate(signerEntry.Certificate())
	parameters.SetCertificateChainFromTokens(signerEntry.CertificateChain()...)
	return parameters
}

// generate signs sampleXMLContent at the given level/packaging and writes it to outDir/name.
func generate(outDir, name string, level enumerations.SignatureLevel, packaging enumerations.SignaturePackaging,
	signerEntry token.DSSPrivateKeyEntry, tspSource validation.TSPSource) error {
	parameters := newParameters(level, packaging, signerEntry)

	service := xades.NewXAdESService(validation.NewCommonCertificateVerifier())
	if tspSource != nil {
		service.TspSource = tspSource
	}

	toSignDocument := model.NewInMemoryDocumentWithMimeType([]byte(sampleXMLContent), "sample.xml", enumerations.MimeTypeEnumXML)

	dataToSign := service.GetDataToSign(toSignDocument, parameters)
	signatureToken, err := reopenSignatureToken()
	if err != nil {
		return err
	}
	signatureValue, err := signatureToken.Sign(dataToSign, parameters.DigestAlgorithm(), signerEntry)
	if err != nil {
		return err
	}
	signedDocument := service.SignDocument(toSignDocument, parameters, signatureValue)
	return writeDocument(signedDocument, filepath.Join(outDir, name))
}

// generateDetached signs sampleXMLContent as a DETACHED XAdES-B signature, writing both the
// signature and the original content it covers (upstream needs the latter as detached content to
// validate reference/message-digest intactness).
func generateDetached(outDir string, signerEntry token.DSSPrivateKeyEntry) error {
	parameters := newParameters(enumerations.SignatureLevelXAdESBaselineB, enumerations.SignaturePackagingDetached, signerEntry)

	service := xades.NewXAdESService(validation.NewCommonCertificateVerifier())
	toSignDocument := model.NewInMemoryDocumentWithMimeType([]byte(sampleXMLContent), "sample-detached.xml", enumerations.MimeTypeEnumXML)

	dataToSign := service.GetDataToSign(toSignDocument, parameters)
	signatureToken, err := reopenSignatureToken()
	if err != nil {
		return err
	}
	signatureValue, err := signatureToken.Sign(dataToSign, parameters.DigestAlgorithm(), signerEntry)
	if err != nil {
		return err
	}
	signedDocument := service.SignDocument(toSignDocument, parameters, signatureValue)
	if err := writeDocument(signedDocument, filepath.Join(outDir, "xades-b-detached.xml")); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(outDir, "xades-b-detached-content.xml"), []byte(sampleXMLContent), 0o644)
}

// reopenSignatureToken opens the signer's key store afresh. XAdESSignatureParameters carries no
// live token/session, only the certificate/chain, so signing always goes back through a
// SignatureTokenConnection.
func reopenSignatureToken() (*token.Pkcs12SignatureToken, error) {
	selfDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(filepath.Join(selfDir, "signer_rsa.p12"))
	if err != nil {
		return nil, err
	}
	return token.NewPkcs12SignatureTokenFromBytes(data, token.NewPasswordProtection([]byte("testpassword")))
}

func writeDocument(document model.DSSDocument, path string) error {
	return document.Save(path)
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, "crossgen:", err)
	os.Exit(1)
}
