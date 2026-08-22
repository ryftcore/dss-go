// Command crossgen is the GO -> UPSTREAM direction of the PAdES cross-validation harness (task
// #12, PAdES extension): it signs three corpus PDFs of different xref styles (a classic-xref-
// table document, and two cross-reference-stream documents from different producers) with this
// package's own Service, producing invisible PAdES-B and PAdES-T signatures with real crypto
// (an RSA PKCS#12 test key for the signer, an EC PKCS#12 test key as a self-hosted TSA via
// spi/validation.KeyEntityTSPSource - the same pattern dss/cades and dss/xades's own crossgen
// generators use), and writes them to files. The point is a real, independently-verifiable
// artifact: pades_downstream_cross_validation_test.go runs this program and then hands its output
// to CrossGenValidator.java, which loads each file with upstream DSS's own
// SignedDocumentValidator/SignedDocumentDiagnosticDataBuilder and asserts the signature is
// intact, the signing certificate is identified, and the level is recognized - upstream DSS
// accepting what this port produced is the actual proof of compatibility, not another Go-side
// assertion.
//
// Usage: go run . <output directory>
//
// Writes, for each of corpus/EmptyPage.pdf (cross-reference stream), corpus/testdoc.pdf (classic
// xref table) and corpus/pdf-xref-streams.pdf (cross-reference stream, a different producer):
// <outdir>/<stem>-b.pdf (PAdES-BASELINE-B) and <outdir>/<stem>-t.pdf (PAdES-BASELINE-T).
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/pades"
	"github.com/ryftcore/dss-go/dss/spi/validation"
	"github.com/ryftcore/dss-go/dss/token"
)

// corpusFiles are the three corpus PDFs of different xref styles this generator signs, relative
// to this file's own directory. See the package doc comment for what distinguishes each.
var corpusFiles = []string{
	filepath.Join("corpus", "EmptyPage.pdf"),
	filepath.Join("corpus", "testdoc.pdf"),
	filepath.Join("corpus", "pdf-xref-streams.pdf"),
}

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

	for _, corpusFile := range corpusFiles {
		stem := strings.TrimSuffix(filepath.Base(corpusFile), ".pdf")
		inputPath := filepath.Join(selfDir, corpusFile)

		bName := stem + "-b.pdf"
		if err := generate(inputPath, filepath.Join(outDir, bName),
			enumerations.SignatureLevelPAdESBaselineB, signerEntry, nil); err != nil {
			fail(fmt.Errorf("generating %s: %w", bName, err))
		}
		fmt.Println("wrote", filepath.Join(outDir, bName))

		tName := stem + "-t.pdf"
		if err := generate(inputPath, filepath.Join(outDir, tName),
			enumerations.SignatureLevelPAdESBaselineT, signerEntry, tspSource); err != nil {
			fail(fmt.Errorf("generating %s: %w", tName, err))
		}
		fmt.Println("wrote", filepath.Join(outDir, tName))
	}
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

// newParameters builds the SignatureParameters shared by every generated signature: the
// digest algorithm and signing certificate/chain. No SignatureImageParameters are set, so every
// generated signature is invisible (no visible signature appearance is drawn on the page) - the
// native engine has no rasteriser to draw one with anyway (internal/pdf/DESIGN.md §0.2).
func newParameters(level enumerations.SignatureLevel, signerEntry token.DSSPrivateKeyEntry) *pades.SignatureParameters {
	parameters := pades.NewPAdESSignatureParameters()
	parameters.SetSignatureLevel(level)
	parameters.SetDigestAlgorithm(enumerations.DigestAlgorithmSHA256)
	parameters.SetSigningCertificate(signerEntry.Certificate())
	parameters.SetCertificateChainFromTokens(signerEntry.CertificateChain()...)
	return parameters
}

// generate signs inputPath (an unsigned corpus PDF) at the given PAdES level and writes the
// result to outputPath.
func generate(inputPath, outputPath string, level enumerations.SignatureLevel,
	signerEntry token.DSSPrivateKeyEntry, tspSource validation.TSPSource) error {
	parameters := newParameters(level, signerEntry)

	service := pades.NewPAdESService(validation.NewCommonCertificateVerifier())
	if tspSource != nil {
		service.SetTspSource(tspSource)
	}

	toSignDocument, err := model.NewFileDocument(inputPath)
	if err != nil {
		return err
	}

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
	return writeDocument(signedDocument, outputPath)
}

// reopenSignatureToken opens the signer's key store afresh. SignatureParameters carries no
// live token/session, only the certificate/chain, so signing always goes back through a
// SignatureTokenConnection - the same reason cades/testdata/crossgen/main.go's own
// reopenSignatureToken exists.
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
