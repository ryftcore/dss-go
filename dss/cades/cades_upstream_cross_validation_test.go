// Cross-validation harness, direction UPSTREAM -> GO (task #12): parses the CAdES signatures
// checked into testdata/upstream/ with this package's own CMSDocumentAnalyzer/CAdESSignature
// and compares the result against testdata/upstream-cross-validation.json, ground truth dumped
// straight from upstream DSS 6.5.RC1's CMSDocumentAnalyzer/CAdESSignature (see
// testdata/gen/CrossValidationOracle.java for how to regenerate it, and its header for exactly
// what it asserts and why). This is the first direction of the compatibility proof: signatures
// upstream DSS produced/accepts must parse in the Go port to the same signature count, signing
// certificate, claimed signing time, detected level, counter-signature structure, and (wherever
// the original content is available to check it) reference/message-digest intactness.
//
// The reverse direction - Go-produced signatures validated by upstream DSS itself - lives in
// cades_downstream_cross_validation_test.go (testdata/crossgen/).
package cades

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/corpustest"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/spi/validation"
	"github.com/ryftcore/dss-go/dss/spi/validation/analyzer"
)

type xvalSignature struct {
	SigningCertificateFound  bool            `json:"signingCertificateFound"`
	SigningCertificateSHA256 *string         `json:"signingCertificateSHA256"`
	ClaimedSigningTimeMillis *int64          `json:"claimedSigningTimeMillis"`
	DataFoundUpToLevel       string          `json:"dataFoundUpToLevel"`
	SignerInformationDigest  *string         `json:"signerInformationDigestSHA256"`
	SignerIDIssuerSerial     *string         `json:"signerIdIssuerSerial"`
	SignerIDSubjectKeyID     *string         `json:"signerIdSubjectKeyIdentifier"`
	IsCounterSignature       bool            `json:"isCounterSignature"`
	ReferenceDataFound       bool            `json:"referenceDataFound"`
	ReferenceDataIntact      bool            `json:"referenceDataIntact"`
	SignatureIntact          bool            `json:"signatureIntact"`
	MessageDigestValueHex    *string         `json:"messageDigestValueHex"`
	CounterSignatureCount    int             `json:"counterSignatureCount"`
	CounterSignatures        []xvalSignature `json:"counterSignatures"`
}

type xvalFile struct {
	Path           string          `json:"path"`
	SignatureCount int             `json:"signatureCount"`
	Signatures     []xvalSignature `json:"signatures"`
}

type xvalGolden struct {
	Files []xvalFile `json:"files"`
}

// xvalDetachedContent mirrors gen/CrossValidationOracle.java's DETACHED_CONTENT map: the one
// genuinely detached signature in the fixture set, and the original content file to feed it.
var xvalDetachedContent = map[string]string{
	"validation/dss-1188/Test.bin.sig": "validation/dss-1188/Test.bin",
}

// NOTE ON THE FORMER LEVEL-DETECTION GAP: six fixtures here used to report a lower
// DataFoundUpToLevel() than upstream (CAdES-BASELINE-T instead of -LTA, CAdES-T instead of
// CAdES-A, CAdES-C instead of CAdES-A), which was first written off as a Phase-8
// revocation-matching difference. It was in fact a one-line mis-port in spi/validation:
// SignatureValidationContext's membership test for processedRevocations compared DSS Ids only,
// while Java's RevocationToken#equals compares the DSS Id *and* the related certificate, so a
// single CRL covering several certificates of a chain collapsed into one entry and every
// certificate but the first lost its revocation data. Fixed in
// signatureValidationContextContainsRevocation; all six fixtures are asserted strictly below,
// with no quarantine list left in place to re-populate.

// NOTE ON THE FORMER SIGNATURE-INTEGRITY GAP: two fixtures here (cades-lta-en319122.cms,
// cades-lta-ts101733.cms) are signed with a DigestInfo whose digest AlgorithmIdentifier omits
// the SHA-256 OID's trailing NULL parameter - a non-canonical but real-world encoding. Go's
// crypto/rsa.VerifyPKCS1v15 only accepts the canonical spelling, so this port used to report
// those signatures as broken while upstream DSS (through BouncyCastle's RSADigestSigner, which
// has an explicit "NULL left out" fallback) reported them intact. That was a real fidelity
// defect, not an unfixable standard-library limitation, and it is now fixed in
// spi.SignerInformationVerifier#Verify; both fixtures are asserted strictly below like every
// other one, with no quarantine list left in place to re-populate.

// TestUpstreamCrossValidation replays testdata/upstream-cross-validation.json against this
// package's own parse of testdata/upstream/.
func TestUpstreamCrossValidation(t *testing.T) {
	goldenBytes, err := os.ReadFile(corpustest.Path(t, "upstream-cross-validation.json"))
	if err != nil {
		t.Fatal(err)
	}
	var golden xvalGolden
	if err := json.Unmarshal(goldenBytes, &golden); err != nil {
		t.Fatalf("parsing golden JSON: %v", err)
	}
	if len(golden.Files) == 0 {
		t.Fatal("golden JSON carries no files")
	}

	for _, gf := range golden.Files {
		gf := gf
		t.Run(gf.Path, func(t *testing.T) {
			fullPath := cadesFixturePath(t, filepath.Join("upstream", filepath.FromSlash(gf.Path)))
			document, err := model.NewFileDocument(fullPath)
			if err != nil {
				t.Fatalf("NewFileDocument(%s): %v", fullPath, err)
			}

			a, err := NewCMSDocumentAnalyzerFromDocument(document)
			if err != nil {
				t.Fatalf("NewCMSDocumentAnalyzerFromDocument: %v", err)
			}
			a.SetCertificateVerifier(validation.NewCommonCertificateVerifier())

			if detachedRel, ok := xvalDetachedContent[gf.Path]; ok {
				detachedPath := cadesFixturePath(t, filepath.Join("upstream", filepath.FromSlash(detachedRel)))
				detachedDoc, err := model.NewFileDocument(detachedPath)
				if err != nil {
					t.Fatalf("NewFileDocument(%s): %v", detachedPath, err)
				}
				a.SetDetachedContents([]model.DSSDocument{detachedDoc})
			}

			signatures := a.Signatures()
			if len(signatures) != gf.SignatureCount {
				t.Fatalf("signature count = %d, golden wants %d", len(signatures), gf.SignatureCount)
			}
			if len(signatures) != len(gf.Signatures) {
				t.Fatalf("golden carries %d signature entries for a count of %d", len(gf.Signatures), gf.SignatureCount)
			}

			for i, sig := range signatures {
				// Top-level signatures are already initialized by CMSDocumentAnalyzer's own
				// BuildSignatures(); only recursed-into counter signatures need it (see
				// checkXvalSignature's isCounterSig parameter).
				checkXvalSignature(t, sig.(*CAdESSignature), gf.Signatures[i], i, false)
			}
		})
	}
}

func checkXvalSignature(t *testing.T, sig *CAdESSignature, want xvalSignature, index int, isCounterSig bool) {
	t.Helper()

	// A nested counter signature is not initialized by CMSDocumentAnalyzer's BuildSignatures()
	// (only the top-level ones are), and DataFoundUpToLevel() requires it - exactly the same
	// reason gen/CrossValidationOracle.java's dumpSignature() calls the Java equivalent only for
	// what it recurses into. Calling it a second time on an already-initialized top-level
	// signature is not part of its contract, so it is skipped there.
	if isCounterSig {
		sig.InitBaselineRequirementsChecker(validation.NewCommonCertificateVerifier())
	}

	certificateToken := sig.SigningCertificateToken()
	gotFound := certificateToken != nil
	if gotFound != want.SigningCertificateFound {
		t.Errorf("signature[%d]: signing certificate found = %v, want %v", index, gotFound, want.SigningCertificateFound)
	}
	if want.SigningCertificateSHA256 != nil && certificateToken != nil {
		gotDigest, err := spi.DSSUtilsDigest(enumerations.DigestAlgorithm_SHA256, certificateToken.Encoded())
		if err != nil {
			t.Fatalf("signature[%d]: digesting signing certificate: %v", index, err)
		}
		if gotHex := hexEncode(gotDigest); gotHex != *want.SigningCertificateSHA256 {
			t.Errorf("signature[%d]: signing certificate SHA-256 = %s, want %s", index, gotHex, *want.SigningCertificateSHA256)
		}
	}

	gotSigningTime := sig.SigningTime()
	switch {
	case want.ClaimedSigningTimeMillis == nil && gotSigningTime != nil:
		t.Errorf("signature[%d]: signing time = %v, want none", index, gotSigningTime)
	case want.ClaimedSigningTimeMillis != nil && gotSigningTime == nil:
		t.Errorf("signature[%d]: signing time = none, want %d", index, *want.ClaimedSigningTimeMillis)
	case want.ClaimedSigningTimeMillis != nil && gotSigningTime != nil:
		if gotMillis := gotSigningTime.UnixMilli(); gotMillis != *want.ClaimedSigningTimeMillis {
			t.Errorf("signature[%d]: signing time = %d ms, want %d ms", index, gotMillis, *want.ClaimedSigningTimeMillis)
		}
	}

	if gotLevel := string(sig.DataFoundUpToLevel()); gotLevel != want.DataFoundUpToLevel {
		t.Errorf("signature[%d]: DataFoundUpToLevel() = %s, want %s", index, gotLevel, want.DataFoundUpToLevel)
	}

	if gotCounter := sig.IsCounterSignature(); gotCounter != want.IsCounterSignature {
		t.Errorf("signature[%d]: IsCounterSignature() = %v, want %v", index, gotCounter, want.IsCounterSignature)
	}

	verification := sig.SignatureCryptographicVerification()
	if got := verification.IsReferenceDataFound(); got != want.ReferenceDataFound {
		t.Errorf("signature[%d]: ReferenceDataFound = %v, want %v", index, got, want.ReferenceDataFound)
	}
	if got := verification.IsReferenceDataIntact(); got != want.ReferenceDataIntact {
		t.Errorf("signature[%d]: ReferenceDataIntact = %v, want %v", index, got, want.ReferenceDataIntact)
	}
	if gotSignatureIntact := verification.IsSignatureIntact(); gotSignatureIntact != want.SignatureIntact {
		t.Errorf("signature[%d]: SignatureIntact = %v, want %v", index, gotSignatureIntact, want.SignatureIntact)
	}

	if want.MessageDigestValueHex != nil {
		gotDigest := sig.MessageDigestValue()
		if gotHex := hexEncode(gotDigest); gotHex != *want.MessageDigestValueHex {
			t.Errorf("signature[%d]: MessageDigestValue() = %s, want %s", index, gotHex, *want.MessageDigestValueHex)
		}
	}

	// DSS-Id stability: building the identifier twice from the same, already-parsed signature
	// must yield the same string both times (DefaultAdvancedSignature caches it on first use, so
	// this also exercises that the cache does not itself introduce drift).
	id1 := sig.ID()
	id2 := sig.ID()
	if id1 != id2 || id1 == "" {
		t.Errorf("signature[%d]: DSS-Id not stable: %q then %q", index, id1, id2)
	}

	counterSignatures := sig.CounterSignatures()
	if len(counterSignatures) != want.CounterSignatureCount {
		t.Fatalf("signature[%d]: counter-signature count = %d, want %d", index, len(counterSignatures), want.CounterSignatureCount)
	}
	if len(counterSignatures) != len(want.CounterSignatures) {
		t.Fatalf("signature[%d]: golden carries %d counter-signature entries for a count of %d", index, len(want.CounterSignatures), want.CounterSignatureCount)
	}
	for i, cs := range counterSignatures {
		checkXvalSignature(t, cs.(*CAdESSignature), want.CounterSignatures[i], i, true)
	}
}

func hexEncode(b []byte) string {
	const hexdigits = "0123456789abcdef"
	out := make([]byte, len(b)*2)
	for i, v := range b {
		out[i*2] = hexdigits[v>>4]
		out[i*2+1] = hexdigits[v&0xf]
	}
	return string(out)
}

// compile-time reminder that analyzer.DocumentAnalyzer is the interface CMSDocumentAnalyzer
// implements, kept honest for the a.SetCertificateVerifier/SetDetachedContents/Signatures calls
// above (all promoted from analyzer.DefaultDocumentAnalyzer).
var _ analyzer.DocumentAnalyzer = (*CMSDocumentAnalyzer)(nil)
