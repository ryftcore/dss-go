// Command crossgen is the GO -> UPSTREAM direction of the cross-validation harness (task #12,
// ASiC-with-CAdES extension): it builds ASiC-S and ASiC-E containers with this package's own
// ASiCWithCAdESService, at both baseline B and T, using real crypto (an RSA PKCS#12 test key for
// the signer, an EC PKCS#12 test key as a self-hosted TSA via
// spi/validation.KeyEntityTSPSource - the same two key stores cades/testdata/crossgen already
// uses, copied here rather than shared across packages since go.mod has no notion of a "testdata"
// import), and writes the resulting containers to files. The point is a real,
// independently-verifiable artifact: the sibling asic_downstream_cross_validation_test.go runs
// this program and then hands its output to CrossGenValidator.java, which loads each container
// with upstream DSS 6.5.RC1's own SignedDocumentValidator/SignedDocumentDiagnosticDataBuilder and
// asserts the container type is recognized, the signature is intact, the signing certificate is
// identified, and the level is recognized - upstream DSS accepting what this port produced is the
// actual proof of compatibility, not another Go-side assertion.
//
// Usage: go run . <output directory>
//
// Writes <outdir>/asics-cades-b.scs, <outdir>/asics-cades-t.scs, <outdir>/asice-cades-b.sce and
// <outdir>/asice-cades-t.sce.
package main

import (
	"fmt"
	"os"
	"path/filepath"

	asiccades "github.com/utain/esig/dss/asic/cades"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/spi/validation"
	"github.com/utain/esig/dss/token"
)

// sampleContent* are the fixed payloads every generated container covers. Their own text records
// what generated them, so a file found on disk explains itself. ASiC-E carries two entries so the
// container-level manifest this format needs is exercised for real, not degenerate to one entry.
var (
	sampleContentSingle = []byte("DSS Go port cross-validation sample content - github.com/utain/esig/dss " +
		"asic/cades/testdata/crossgen (ASiC-S). Signed by the Go port's own ASiCWithCAdESService, " +
		"verified by upstream DSS 6.5.RC1's SignedDocumentValidator.")
	sampleContentMultiA = []byte("DSS Go port cross-validation sample content - github.com/utain/esig/dss " +
		"asic/cades/testdata/crossgen (ASiC-E, entry A).")
	sampleContentMultiB = []byte("DSS Go port cross-validation sample content - github.com/utain/esig/dss " +
		"asic/cades/testdata/crossgen (ASiC-E, entry B).")
)

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

	singleDoc := []model.DSSDocument{model.NewInMemoryDocumentWithName(sampleContentSingle, "sample.txt")}
	multiDocs := []model.DSSDocument{
		model.NewInMemoryDocumentWithName(sampleContentMultiA, "sample-a.txt"),
		model.NewInMemoryDocumentWithName(sampleContentMultiB, "sample-b.txt"),
	}

	cases := []struct {
		name          string
		documents     []model.DSSDocument
		containerType enumerations.ASiCContainerType
		level         enumerations.SignatureLevel
		tsp           validation.TSPSource
	}{
		{"asics-cades-b.scs", singleDoc, enumerations.ASiCContainerType_ASiC_S, enumerations.SignatureLevel_CAdES_BASELINE_B, nil},
		{"asics-cades-t.scs", singleDoc, enumerations.ASiCContainerType_ASiC_S, enumerations.SignatureLevel_CAdES_BASELINE_T, tspSource},
		{"asice-cades-b.sce", multiDocs, enumerations.ASiCContainerType_ASiC_E, enumerations.SignatureLevel_CAdES_BASELINE_B, nil},
		{"asice-cades-t.sce", multiDocs, enumerations.ASiCContainerType_ASiC_E, enumerations.SignatureLevel_CAdES_BASELINE_T, tspSource},
	}

	for _, testCase := range cases {
		if err := generate(outDir, testCase.name, testCase.documents, testCase.containerType,
			testCase.level, signerEntry, testCase.tsp); err != nil {
			fail(fmt.Errorf("generating %s: %w", testCase.name, err))
		}
		fmt.Println("wrote", filepath.Join(outDir, testCase.name))
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

// generate builds and signs an ASiC container of the given type/level over documents, writing it
// to outDir/name.
func generate(outDir, name string, documents []model.DSSDocument, containerType enumerations.ASiCContainerType,
	level enumerations.SignatureLevel, signerEntry token.DSSPrivateKeyEntry, tspSource validation.TSPSource) error {
	parameters := asiccades.NewASiCWithCAdESSignatureParameters()
	parameters.SetSignatureLevel(level)
	parameters.SetDigestAlgorithm(enumerations.DigestAlgorithm_SHA256)
	parameters.SetSigningCertificate(signerEntry.Certificate())
	parameters.SetCertificateChainFromTokens(signerEntry.CertificateChain()...)
	parameters.ASiC().SetContainerType(containerType)

	service := asiccades.NewASiCWithCAdESService(validation.NewCommonCertificateVerifier())
	if tspSource != nil {
		service.TspSource = tspSource
	}

	dataToSign := service.GetDataToSignMultiple(documents, parameters)
	signatureToken, err := reopenSignatureToken()
	if err != nil {
		return err
	}
	signatureValue, err := signatureToken.Sign(dataToSign, parameters.DigestAlgorithm(), signerEntry)
	if err != nil {
		return err
	}
	signedDocument := service.SignDocumentMultiple(documents, parameters, signatureValue)
	return writeDocument(signedDocument, filepath.Join(outDir, name))
}

// reopenSignatureToken opens the signer's key store afresh. ASiCWithCAdESSignatureParameters
// carries no live token/session, only the certificate/chain, so signing always goes back through
// a SignatureTokenConnection.
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
