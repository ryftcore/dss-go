// Cross-validation harness, direction UPSTREAM -> GO (task #12, XAdES extension): parses the
// XAdES signatures checked into testdata/upstream/ with this package's own
// XMLDocumentAnalyzer/XAdESSignature and compares the result against
// testdata/upstream-cross-validation.json, ground truth dumped straight from upstream DSS
// 6.5.RC1's own XMLDocumentAnalyzer/XAdESSignature (see testdata/gen/CrossValidationOracle.java
// for how to regenerate it, and its header for exactly what it asserts and why). This is the
// first direction of the compatibility proof: signatures upstream DSS produced/accepts must parse
// in the Go port to the same signature count, signing certificate, claimed signing time, detected
// level, counter-signature structure, and - via the full per-Reference breakdown, XAdES's own
// natural identity anchor - reference/message-digest intactness across enveloped, enveloping, and
// (where the referenced content is unavailable, as it genuinely is for several fixtures below)
// not-found packaging, several digest algorithms, and both exclusive and inclusive C14N.
//
// The reverse direction - Go-produced signatures validated by upstream DSS itself - lives in
// xades_downstream_cross_validation_test.go (testdata/crossgen/).
//
// NOTE ON dss1811-multi-algo.xml AND Signature-X-BG-1.xml/Signature-X-AT-1.xml: these three
// fixtures deliberately/incidentally carry references that fail to resolve or intentionally
// mismatch even under upstream DSS's own Apache Santuario-based verification (confirmed by
// running gen/CrossValidationOracle.java directly and inspecting its "Expected Digest .../ Actual
// Digest ..." and "Could not find a resolver for URI ..." log output) - NOT a Go-port defect. The
// golden JSON captures upstream's actual (failing) verdict for them, and this test asserts the Go
// port reproduces that same failing verdict exactly, which is the correct cross-validation
// behaviour: matching upstream, whatever upstream says, not asserting our own assumption of what
// "should" happen.
//
// DEFECTS THIS HARNESS FOUND, ALL SINCE FIXED (kept here because each one names the fixture that
// catches it, and every one of them was a silent wrong answer rather than a crash):
//
//  1. Signature-X-BE_ECON-3.xml reported XAdES_BASELINE_T where upstream reports
//     XAdES_BASELINE_LT. Cause: crypto/x509.ParseRevocationList rejects X.509 v1 CRLs
//     ("x509: unsupported crl version"), because RFC 5280 makes TBSCertList.version OPTIONAL and
//     Go insists on it. Two of this signature's four embedded CRLs (Belgium Root CA2's and
//     Belgium Root CA3's) are v1, so their revocation data vanished and the LT requirement could
//     not be met. Fixed in crlparser (crlUtilsParseRevocationList).
//  2. xades-ecc-brainpool.xml reported no signing certificate, XAdES_BES, and a broken
//     signature. Cause: crypto/x509 knows only P-224/256/384/521 and fails the entire
//     certificate parse for any other curve, so this brainpoolP256r1 certificate - and every
//     certificate on a non-NIST curve anywhere in the port - was unparseable. Fixed in
//     internal/eccurve.
//  3. Signature-X-CZ_SEF-5.xml reported XML_NOT_ETSI where upstream reports XAdES_BASELINE_LT.
//     Cause: the manifest entries of a ds:Reference[@Type=".../Manifest"] were computed and then
//     dropped, because model.ReferenceValidation had no way to attach them; the six
//     DataObjectFormat properties that reference those entries then resolved to nothing and the
//     BES-profile check failed. Fixed with model.ReferenceValidation.AddDependentValidations,
//     and the golden data now carries the dependentValidations this test asserts - without which
//     the defect stayed invisible even with the fixture present.
//
// There are NO quarantines here: no skipped fixture, no soft assertion, no t.Logf standing in
// for t.Error. Every fixture in the golden file is asserted in full, and a failure is a failure.
package xades

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

type xvalReferenceValidation struct {
	Type   string `json:"type"`
	ID     string `json:"id"`
	URI    string `json:"uri"`
	Found  bool   `json:"found"`
	Intact bool   `json:"intact"`
	// DependentValidations are the entries of a signed ds:Manifest that a
	// ds:Reference[@Type=".../Manifest"] points at. They are part of the golden data because a
	// DataObjectFormat qualifying property may reference a manifest entry instead of a
	// ds:SignedInfo/ds:Reference: dropping them makes every such signature fail its BES-profile
	// check and be reported as XML_NOT_ETSI, which is precisely what this port did until the
	// audit added Signature-X-CZ_SEF-5.xml and started asserting this field.
	DependentValidations []xvalReferenceValidation `json:"dependentValidations"`
}

type xvalSignature struct {
	SigningCertificateFound  bool                      `json:"signingCertificateFound"`
	SigningCertificateSHA256 *string                   `json:"signingCertificateSHA256"`
	ClaimedSigningTimeMillis *int64                    `json:"claimedSigningTimeMillis"`
	DataFoundUpToLevel       string                    `json:"dataFoundUpToLevel"`
	IsCounterSignature       bool                      `json:"isCounterSignature"`
	ReferenceDataFound       bool                      `json:"referenceDataFound"`
	ReferenceDataIntact      bool                      `json:"referenceDataIntact"`
	SignatureIntact          bool                      `json:"signatureIntact"`
	ReferenceValidationCount int                       `json:"referenceValidationCount"`
	ReferenceValidations     []xvalReferenceValidation `json:"referenceValidations"`
	CounterSignatureCount    int                       `json:"counterSignatureCount"`
	CounterSignatures        []xvalSignature           `json:"counterSignatures"`
}

type xvalFile struct {
	Path           string          `json:"path"`
	SignatureCount int             `json:"signatureCount"`
	Signatures     []xvalSignature `json:"signatures"`
}

type xvalGolden struct {
	Files []xvalFile `json:"files"`
}

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
			fullPath := xadesFixturePath(t, filepath.Join("upstream", filepath.FromSlash(gf.Path)))
			document, err := model.NewFileDocument(fullPath)
			if err != nil {
				t.Fatalf("NewFileDocument(%s): %v", fullPath, err)
			}

			a, err := NewXMLDocumentAnalyzer(document)
			if err != nil {
				t.Fatalf("NewXMLDocumentAnalyzer: %v", err)
			}
			a.SetCertificateVerifier(validation.NewCommonCertificateVerifier())

			signatures := a.Signatures()
			if len(signatures) != gf.SignatureCount {
				t.Fatalf("signature count = %d, golden wants %d", len(signatures), gf.SignatureCount)
			}
			if len(signatures) != len(gf.Signatures) {
				t.Fatalf("golden carries %d signature entries for a count of %d", len(gf.Signatures), gf.SignatureCount)
			}

			for i, sig := range signatures {
				// Top-level signatures are already initialized by XMLDocumentAnalyzer's own
				// BuildSignatures(); only recursed-into counter signatures need it (see
				// checkXvalSignature's isCounterSig parameter).
				checkXvalSignature(t, sig.(*XAdESSignature), gf.Signatures[i], i, false)
			}
		})
	}
}

func checkXvalSignature(t *testing.T, sig *XAdESSignature, want xvalSignature, index int, isCounterSig bool) {
	t.Helper()

	// A nested counter signature is not initialized by XMLDocumentAnalyzer's BuildSignatures()
	// (only the top-level ones are), and DataFoundUpToLevel() requires it - the same reason
	// gen/CrossValidationOracle.java's dumpSignature() calls the Java equivalent only for what it
	// recurses into.
	if isCounterSig {
		sig.InitBaselineRequirementsChecker(validation.NewCommonCertificateVerifier())
	}

	certificateToken := sig.SigningCertificateToken()
	gotFound := certificateToken != nil
	if gotFound != want.SigningCertificateFound {
		t.Errorf("signature[%d]: signing certificate found = %v, want %v", index, gotFound, want.SigningCertificateFound)
	}
	if want.SigningCertificateSHA256 != nil && certificateToken != nil {
		gotDigest, err := spi.DSSUtilsDigest(enumerations.DigestAlgorithmSHA256, certificateToken.Encoded())
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

	referenceValidations := sig.ReferenceValidations()
	if len(referenceValidations) != want.ReferenceValidationCount {
		t.Fatalf("signature[%d]: reference validation count = %d, want %d", index, len(referenceValidations), want.ReferenceValidationCount)
	}
	if len(referenceValidations) != len(want.ReferenceValidations) {
		t.Fatalf("signature[%d]: golden carries %d reference validation entries for a count of %d", index, len(want.ReferenceValidations), want.ReferenceValidationCount)
	}
	for i, rv := range referenceValidations {
		wantRV := want.ReferenceValidations[i]
		if gotType := string(rv.Type()); gotType != wantRV.Type {
			t.Errorf("signature[%d]: referenceValidations[%d].Type() = %s, want %s", index, i, gotType, wantRV.Type)
		}
		if rv.Id() != wantRV.ID {
			t.Errorf("signature[%d]: referenceValidations[%d].Id() = %q, want %q", index, i, rv.Id(), wantRV.ID)
		}
		if rv.Uri() != wantRV.URI {
			t.Errorf("signature[%d]: referenceValidations[%d].Uri() = %q, want %q", index, i, rv.Uri(), wantRV.URI)
		}
		if rv.IsFound() != wantRV.Found {
			t.Errorf("signature[%d]: referenceValidations[%d].IsFound() = %v, want %v", index, i, rv.IsFound(), wantRV.Found)
		}
		if rv.IsIntact() != wantRV.Intact {
			t.Errorf("signature[%d]: referenceValidations[%d].IsIntact() = %v, want %v", index, i, rv.IsIntact(), wantRV.Intact)
		}
		checkXvalDependentValidations(t, rv.DependentValidations(), wantRV.DependentValidations, index, i)
	}

	// DSS-Id stability: building the identifier twice from the same, already-parsed signature must
	// yield the same string both times (DefaultAdvancedSignature caches it on first use, so this
	// also exercises that the cache does not itself introduce drift).
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
		checkXvalSignature(t, cs.(*XAdESSignature), want.CounterSignatures[i], i, true)
	}
}

// checkXvalDependentValidations asserts a MANIFEST reference validation's manifest entries
// against the golden data, in order.
func checkXvalDependentValidations(t *testing.T, got []*model.ReferenceValidation,
	want []xvalReferenceValidation, signatureIndex, referenceIndex int) {
	t.Helper()
	if len(got) != len(want) {
		t.Errorf("signature[%d]: referenceValidations[%d] has %d dependent validation(s), want %d",
			signatureIndex, referenceIndex, len(got), len(want))
		return
	}
	for i, dv := range got {
		wantDV := want[i]
		if gotType := string(dv.Type()); gotType != wantDV.Type {
			t.Errorf("signature[%d]: referenceValidations[%d].DependentValidations()[%d].Type() = %s, want %s",
				signatureIndex, referenceIndex, i, gotType, wantDV.Type)
		}
		if dv.Id() != wantDV.ID {
			t.Errorf("signature[%d]: referenceValidations[%d].DependentValidations()[%d].Id() = %q, want %q",
				signatureIndex, referenceIndex, i, dv.Id(), wantDV.ID)
		}
		if dv.Uri() != wantDV.URI {
			t.Errorf("signature[%d]: referenceValidations[%d].DependentValidations()[%d].Uri() = %q, want %q",
				signatureIndex, referenceIndex, i, dv.Uri(), wantDV.URI)
		}
		if dv.IsFound() != wantDV.Found {
			t.Errorf("signature[%d]: referenceValidations[%d].DependentValidations()[%d].IsFound() = %v, want %v",
				signatureIndex, referenceIndex, i, dv.IsFound(), wantDV.Found)
		}
		if dv.IsIntact() != wantDV.Intact {
			t.Errorf("signature[%d]: referenceValidations[%d].DependentValidations()[%d].IsIntact() = %v, want %v",
				signatureIndex, referenceIndex, i, dv.IsIntact(), wantDV.Intact)
		}
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

// compile-time reminder that analyzer.DocumentAnalyzer is the interface XMLDocumentAnalyzer
// implements, kept honest for the a.SetCertificateVerifier/Signatures calls above (all promoted
// from analyzer.DefaultDocumentAnalyzer).
var _ analyzer.DocumentAnalyzer = (*XMLDocumentAnalyzer)(nil)
