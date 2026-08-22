// Tests for SignatureRequirementsChecker, matching
// dss-document/src/main/java/eu/europa/esig/dss/signature/SignatureRequirementsChecker.java
// upstream.
package document

import (
	"testing"

	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi/validation"
)

// fakeRequirementsSignature is a partial AdvancedSignature double (same technique as
// signatureStatusDeterminismFakeSignature in package validation): it embeds the (nil) interface
// so every method type-checks, overriding only the profile-level predicates
// SignatureRequirementsChecker actually calls.
type fakeRequirementsSignature struct {
	validation.AdvancedSignature
	hasT, hasLT, hasLTA, hasC, hasX, hasXL, hasA bool
	allSelfSigned                                bool
	embeddedEvidenceRecords                      []validation.EvidenceRecord
}

func (f *fakeRequirementsSignature) ID() string                         { return "fake-signature" }
func (f *fakeRequirementsSignature) HasTProfile() bool                  { return f.hasT }
func (f *fakeRequirementsSignature) HasLTProfile() bool                 { return f.hasLT }
func (f *fakeRequirementsSignature) HasLTAProfile() bool                { return f.hasLTA }
func (f *fakeRequirementsSignature) HasCProfile() bool                  { return f.hasC }
func (f *fakeRequirementsSignature) HasXProfile() bool                  { return f.hasX }
func (f *fakeRequirementsSignature) HasXLProfile() bool                 { return f.hasXL }
func (f *fakeRequirementsSignature) HasAProfile() bool                  { return f.hasA }
func (f *fakeRequirementsSignature) AreAllSelfSignedCertificates() bool { return f.allSelfSigned }
func (f *fakeRequirementsSignature) EmbeddedEvidenceRecords() []validation.EvidenceRecord {
	return f.embeddedEvidenceRecords
}

var _ validation.AdvancedSignature = (*fakeRequirementsSignature)(nil)

func newRequirementsChecker(t *testing.T, verifier validation.CertificateVerifier) *SignatureRequirementsChecker[*model.TimestampParameters] {
	t.Helper()
	params := NewAbstractSignatureParameters[*model.TimestampParameters]()
	return NewSignatureRequirementsChecker[*model.TimestampParameters](verifier, &params)
}

// ---- constructor guards -----------------------------------------------------

func TestNewSignatureRequirementsCheckerNilCertificateVerifierPanics(t *testing.T) {
	params := NewAbstractSignatureParameters[*model.TimestampParameters]()
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic for nil CertificateVerifier")
		}
	}()
	NewSignatureRequirementsChecker[*model.TimestampParameters](nil, &params)
}

func TestNewSignatureRequirementsCheckerNilParametersPanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic for nil signature parameters")
		}
	}()
	NewSignatureRequirementsChecker[*model.TimestampParameters](validation.NewCommonCertificateVerifier(), nil)
}

// ---- pure profile-level predicates ------------------------------------------

func TestHasLTLevelOrHigher(t *testing.T) {
	c := newRequirementsChecker(t, validation.NewCommonCertificateVerifier())

	cases := []struct {
		name string
		sig  *fakeRequirementsSignature
		want bool
	}{
		{"LTA profile alone", &fakeRequirementsSignature{hasLTA: true}, true},
		{"LT + T, not all self-signed", &fakeRequirementsSignature{hasLT: true, hasT: true}, true},
		{"C + T, not all self-signed", &fakeRequirementsSignature{hasC: true, hasT: true}, true},
		{"LT + T but all self-signed", &fakeRequirementsSignature{hasLT: true, hasT: true, allSelfSigned: true}, false},
		{"LT without T", &fakeRequirementsSignature{hasLT: true}, false},
		{"none", &fakeRequirementsSignature{}, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := c.HasLTLevelOrHigher(tc.sig); got != tc.want {
				t.Fatalf("HasLTLevelOrHigher() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestHasLTALevelOrHigherAndHasALevelOrHigher(t *testing.T) {
	c := newRequirementsChecker(t, validation.NewCommonCertificateVerifier())
	withLTA := &fakeRequirementsSignature{hasLTA: true}
	withoutLTA := &fakeRequirementsSignature{}

	if !c.HasLTALevelOrHigher(withLTA) {
		t.Fatal("HasLTALevelOrHigher() = false, want true")
	}
	if c.HasLTALevelOrHigher(withoutLTA) {
		t.Fatal("HasLTALevelOrHigher() = true, want false")
	}

	// HasALevelOrHigher is defined to delegate directly to HasLTALevelOrHigher.
	if c.HasALevelOrHigher(withLTA) != c.HasLTALevelOrHigher(withLTA) {
		t.Fatal("HasALevelOrHigher() should mirror HasLTALevelOrHigher()")
	}
}

func TestHasXLevelOrHigher(t *testing.T) {
	c := newRequirementsChecker(t, validation.NewCommonCertificateVerifier())

	cases := []struct {
		name string
		sig  *fakeRequirementsSignature
		want bool
	}{
		{"X profile alone", &fakeRequirementsSignature{hasX: true}, true},
		{"A profile alone", &fakeRequirementsSignature{hasA: true}, true},
		{"XL + T, not self-signed", &fakeRequirementsSignature{hasXL: true, hasT: true}, true},
		{"XL + T but self-signed", &fakeRequirementsSignature{hasXL: true, hasT: true, allSelfSigned: true}, false},
		{"XL without T", &fakeRequirementsSignature{hasXL: true}, false},
		{"none", &fakeRequirementsSignature{}, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := c.HasXLevelOrHigher(tc.sig); got != tc.want {
				t.Fatalf("HasXLevelOrHigher() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestHasXLLevelOrHigher(t *testing.T) {
	c := newRequirementsChecker(t, validation.NewCommonCertificateVerifier())

	cases := []struct {
		name string
		sig  *fakeRequirementsSignature
		want bool
	}{
		{"A profile alone", &fakeRequirementsSignature{hasA: true}, true},
		{"XL + T + X, not self-signed", &fakeRequirementsSignature{hasXL: true, hasT: true, hasX: true}, true},
		{"XL + T + X but self-signed", &fakeRequirementsSignature{hasXL: true, hasT: true, hasX: true, allSelfSigned: true}, false},
		{"XL + T without X", &fakeRequirementsSignature{hasXL: true, hasT: true}, false},
		{"none", &fakeRequirementsSignature{}, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := c.HasXLLevelOrHigher(tc.sig); got != tc.want {
				t.Fatalf("HasXLLevelOrHigher() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestHasEmbeddedEvidenceRecords(t *testing.T) {
	c := newRequirementsChecker(t, validation.NewCommonCertificateVerifier())
	withER := &fakeRequirementsSignature{embeddedEvidenceRecords: []validation.EvidenceRecord{nil}}
	withoutER := &fakeRequirementsSignature{}

	if !c.HasEmbeddedEvidenceRecords(withER) {
		t.Fatal("HasEmbeddedEvidenceRecords() = false, want true")
	}
	if c.HasEmbeddedEvidenceRecords(withoutER) {
		t.Fatal("HasEmbeddedEvidenceRecords() = true, want false")
	}
}

// ---- Assert*Possible: no-op when the relevant alert is explicitly disarmed --------------------

func TestAssertExtendPossibleNoOpsWithoutAugmentationAlert(t *testing.T) {
	// CommonCertificateVerifier arms every augmentation alert with ExceptionOnStatusAlert by
	// default; explicitly silence them to exercise the "alert unset -> no-op" guard clause each
	// assert*IsHighest/assertHasNoEmbeddedEvidenceRecords method starts with.
	verifier := validation.NewCommonCertificateVerifier()
	verifier.SetAugmentationAlertOnHigherSignatureLevel(nil)
	verifier.SetAlertOnInvalidSignature(nil)
	c := newRequirementsChecker(t, verifier)

	// Every one of these signatures would fail its respective check if the alert were armed;
	// with the alert disarmed, all Assert*Possible calls must be silent no-ops.
	signaturesThatWouldFailEverything := []validation.AdvancedSignature{
		&fakeRequirementsSignature{hasLTA: true, embeddedEvidenceRecords: []validation.EvidenceRecord{nil}},
	}

	c.AssertExtendToTLevelPossible(signaturesThatWouldFailEverything)
	c.AssertExtendToLTLevelPossible(signaturesThatWouldFailEverything)
	c.AssertExtendToCLevelPossible(signaturesThatWouldFailEverything)
	c.AssertExtendToXLevelPossible(signaturesThatWouldFailEverything)
	c.AssertExtendToXLLevelPossible(signaturesThatWouldFailEverything)
	c.AssertExtendToLTALevelPossible(signaturesThatWouldFailEverything)
	c.AssertSignaturesValid(signaturesThatWouldFailEverything)
	// Reaching here without panicking is the assertion.
}

// ---- Assert*Possible: panics via mustAlert (CommonCertificateVerifier's default alert is ----
// ---- already armed with ExceptionOnStatusAlert) when the signature exceeds the target level ---

func TestAssertExtendToTLevelPossiblePanicsWhenAlreadyAtHigherLevel(t *testing.T) {
	c := newRequirementsChecker(t, validation.NewCommonCertificateVerifier())

	sig := &fakeRequirementsSignature{hasLTA: true} // HasLTLevelOrHigher() == true
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic: signature already has a higher level than T")
		}
	}()
	c.AssertExtendToTLevelPossible([]validation.AdvancedSignature{sig})
}

func TestAssertExtendToLTALevelPossiblePanicsOnEmbeddedEvidenceRecords(t *testing.T) {
	c := newRequirementsChecker(t, validation.NewCommonCertificateVerifier())

	sig := &fakeRequirementsSignature{embeddedEvidenceRecords: []validation.EvidenceRecord{nil}}
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic: signature is preserved by an embedded evidence record")
		}
	}()
	c.AssertExtendToLTALevelPossible([]validation.AdvancedSignature{sig})
}

func TestAssertExtendToTLevelPossibleNoPanicWhenBelowTargetLevel(t *testing.T) {
	c := newRequirementsChecker(t, validation.NewCommonCertificateVerifier())

	sig := &fakeRequirementsSignature{} // no profile flags set: definitely below T-level
	c.AssertExtendToTLevelPossible([]validation.AdvancedSignature{sig})
	// Reaching here without panicking is the assertion.
}
