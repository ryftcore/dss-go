// Command crossgen is the GO -> UPSTREAM direction of the cross-validation harness (task #12,
// JAdES extension): it signs a fixed sample document with this package's own JAdESService,
// producing JAdES-B and JAdES-T signatures in both COMPACT and FLATTENED JSON serialization, plus
// a DETACHED (sigD, ObjectIdByURIHash mechanism) JAdES-B, with real crypto (an RSA PKCS#12 test
// key for the signer, an EC PKCS#12 test key as a self-hosted TSA via
// spi/validation.KeyEntityTSPSource - the same two key stores cades/testdata/crossgen already
// uses, copied here rather than shared across packages since go.mod has no notion of a "testdata"
// import), and writes them to files. The point is a real, independently-verifiable artifact:
// jades_downstream_cross_validation_test.go runs this program and then hands its output to
// CrossGenValidator.java (an unmodified copy of cades/testdata/crossgen's - it drives upstream
// DSS's own generic SignedDocumentValidator/SignedDocumentDiagnosticDataBuilder, so nothing in it
// is CAdES-specific), which asserts the signature is intact, the signing certificate is
// identified, and the level is recognized - upstream DSS accepting what this port produced is the
// actual proof of compatibility, not another Go-side assertion.
//
// Usage: go run . <output directory>
//
// Writes <outdir>/jades-b-compact.json, <outdir>/jades-t-compact.json,
// <outdir>/jades-b-flattened.json, <outdir>/jades-t-flattened.json, and
// <outdir>/jades-b-detached.json (+ its companion <outdir>/jades-b-detached-content.txt, named to
// match the sigD ObjectIdByURIHash reference CrossGenValidator.java's setDetachedContents feeds
// upstream).
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/jades"
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/spi/validation"
	"github.com/utain/esig/dss/token"
)

// sampleJSONContent is the fixed payload every generated signature covers. Its own text records
// what generated it, so a file found on disk explains itself.
const sampleJSONContent = `{"note":"DSS Go port cross-validation sample content - github.com/utain/esig/dss jades/testdata/crossgen. Signed by the Go port's own JAdESService, verified by upstream DSS 6.5.RC1's SignedDocumentValidator."}`

// detachedContentName is the document name the ObjectIdByURIHash sigD mechanism references (and
// the file name written to <outdir>): matters because the sigD header embeds it as a URI, and
// CrossGenValidator.java's detached FileDocument must carry the identical name for upstream to
// resolve the reference and match the digest.
const detachedContentName = "jades-b-detached-content.txt"

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

	// JWSSerializationType_COMPACT_SERIALIZATION only ever carries JAdES-BASELINE-B: compact has
	// no unprotected-header slot to hold the 'etsiU' array a -T signature-timestamp (or any
	// higher level) needs, exactly what JAdESService.GetDataToSign itself panics with
	// ("Only JAdES_BASELINE_B level is allowed for JAdES Compact Signature!") if asked to do
	// otherwise - so, unlike XAdES/CAdES, there is no compact -T fixture here.
	if err := generate(outDir, "jades-b-compact.json", enumerations.SignatureLevel_JAdES_BASELINE_B,
		enumerations.JWSSerializationType_COMPACT_SERIALIZATION, signerEntry, nil); err != nil {
		fail(fmt.Errorf("generating JAdES-B compact: %w", err))
	}
	if err := generate(outDir, "jades-b-flattened.json", enumerations.SignatureLevel_JAdES_BASELINE_B,
		enumerations.JWSSerializationType_FLATTENED_JSON_SERIALIZATION, signerEntry, nil); err != nil {
		fail(fmt.Errorf("generating JAdES-B flattened: %w", err))
	}
	if err := generate(outDir, "jades-t-flattened.json", enumerations.SignatureLevel_JAdES_BASELINE_T,
		enumerations.JWSSerializationType_FLATTENED_JSON_SERIALIZATION, signerEntry, tspSource); err != nil {
		fail(fmt.Errorf("generating JAdES-T flattened: %w", err))
	}
	if err := generate(outDir, "jades-t-full.json", enumerations.SignatureLevel_JAdES_BASELINE_T,
		enumerations.JWSSerializationType_JSON_SERIALIZATION, signerEntry, tspSource); err != nil {
		fail(fmt.Errorf("generating JAdES-T full JSON serialization: %w", err))
	}
	if err := generateDetached(outDir, signerEntry); err != nil {
		fail(fmt.Errorf("generating JAdES-B detached: %w", err))
	}

	fmt.Println("wrote", filepath.Join(outDir, "jades-b-compact.json"))
	fmt.Println("wrote", filepath.Join(outDir, "jades-b-flattened.json"))
	fmt.Println("wrote", filepath.Join(outDir, "jades-t-flattened.json"))
	fmt.Println("wrote", filepath.Join(outDir, "jades-t-full.json"))
	fmt.Println("wrote", filepath.Join(outDir, "jades-b-detached.json"))
	fmt.Println("wrote", filepath.Join(outDir, detachedContentName))
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

// newParameters builds the JAdESSignatureParameters shared by every generated signature: the
// digest algorithm, signing certificate/chain, serialization type, and (for anything above
// baseline B) the TSP source a JAdES-T needs to request its signature-timestamp.
func newParameters(level enumerations.SignatureLevel, serializationType enumerations.JWSSerializationType,
	signerEntry token.DSSPrivateKeyEntry) *jades.JAdESSignatureParameters {
	parameters := jades.NewJAdESSignatureParameters()
	parameters.SetSignatureLevel(level)
	parameters.SetSignaturePackaging(enumerations.SignaturePackaging_ENVELOPING)
	parameters.SetJwsSerializationType(serializationType)
	parameters.SetDigestAlgorithm(enumerations.DigestAlgorithm_SHA256)
	parameters.SetSigningCertificate(signerEntry.Certificate())
	parameters.SetCertificateChainFromTokens(signerEntry.CertificateChain()...)
	return parameters
}

// generate signs sampleJSONContent at the given level/serialization type and writes it to
// outDir/name.
func generate(outDir, name string, level enumerations.SignatureLevel, serializationType enumerations.JWSSerializationType,
	signerEntry token.DSSPrivateKeyEntry, tspSource validation.TSPSource) error {
	parameters := newParameters(level, serializationType, signerEntry)

	service := jades.NewJAdESService(validation.NewCommonCertificateVerifier())
	if tspSource != nil {
		service.TspSource = tspSource
	}

	toSignDocument := model.NewInMemoryDocumentWithMimeType([]byte(sampleJSONContent), "sample.json", enumerations.MimeTypeEnum_JSON)

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

// generateDetached signs sampleJSONContent as a DETACHED JAdES-B signature using the
// ObjectIdByURIHash sigD mechanism (FLATTENED JSON serialization, the form that mechanism is
// documented against), writing both the signature and the original content it covers (upstream
// needs the latter, under the same document name, to resolve the sigD reference and validate
// reference/message-digest intactness).
func generateDetached(outDir string, signerEntry token.DSSPrivateKeyEntry) error {
	parameters := jades.NewJAdESSignatureParameters()
	parameters.SetSignatureLevel(enumerations.SignatureLevel_JAdES_BASELINE_B)
	parameters.SetSignaturePackaging(enumerations.SignaturePackaging_DETACHED)
	parameters.SetSigDMechanism(enumerations.SigDMechanism_OBJECT_ID_BY_URI_HASH)
	parameters.SetJwsSerializationType(enumerations.JWSSerializationType_FLATTENED_JSON_SERIALIZATION)
	parameters.SetDigestAlgorithm(enumerations.DigestAlgorithm_SHA256)
	parameters.SetSigningCertificate(signerEntry.Certificate())
	parameters.SetCertificateChainFromTokens(signerEntry.CertificateChain()...)

	service := jades.NewJAdESService(validation.NewCommonCertificateVerifier())
	toSignDocument := model.NewInMemoryDocumentWithMimeType([]byte(sampleJSONContent), detachedContentName, enumerations.MimeTypeEnum_JSON)

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
	if err := writeDocument(signedDocument, filepath.Join(outDir, "jades-b-detached.json")); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(outDir, detachedContentName), []byte(sampleJSONContent), 0o644)
}

// reopenSignatureToken opens the signer's key store afresh. JAdESSignatureParameters carries no
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
