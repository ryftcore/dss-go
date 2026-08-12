package model

import (
	"testing"
	"time"

	"github.com/utain/esig/dss/enumerations"
)

// abstractSerializableSignatureParametersSigningDate pins the signing date, which
// BLevelParameters otherwise initialises to "now" (Java: `new Date()`), so that two
// independently constructed parameter sets can be compared at all.
var abstractSerializableSignatureParametersSigningDate = time.Date(2024, time.January, 2, 3, 4, 5, 0, time.UTC)

func abstractSerializableSignatureParametersNew(t *testing.T) *AbstractSerializableSignatureParameters[*TimestampParameters] {
	t.Helper()
	p := NewAbstractSerializableSignatureParameters[*TimestampParameters]()
	p.BLevel().SetSigningDate(&abstractSerializableSignatureParametersSigningDate)
	return &p
}

// TestAbstractSerializableSignatureParametersDefaultsMatchJava pins the field
// initialisers of the Java class.
func TestAbstractSerializableSignatureParametersDefaultsMatchJava(t *testing.T) {
	p := abstractSerializableSignatureParametersNew(t)
	if p.CheckCertificateRevocation() {
		t.Errorf("checkCertificateRevocation must default to false")
	}
	if p.GenerateTBSWithoutCertificate() {
		t.Errorf("generateTBSWithoutCertificate must default to false")
	}
	if got, want := p.SignatureAlgorithm(), enumerations.SignatureAlgorithm_RSA_SSA_PSS_SHA512_MGF1; got != want {
		t.Errorf("SignatureAlgorithm() = %q, want %q", got, want)
	}
	if got, want := p.DigestAlgorithm(), enumerations.SignatureAlgorithm_RSA_SSA_PSS_SHA512_MGF1.DigestAlgorithm(); got != want {
		t.Errorf("DigestAlgorithm() = %q, want %q", got, want)
	}
	if got, want := p.EncryptionAlgorithm(), enumerations.SignatureAlgorithm_RSA_SSA_PSS_SHA512_MGF1.EncryptionAlgorithm(); got != want {
		t.Errorf("EncryptionAlgorithm() = %q, want %q", got, want)
	}
	if p.BLevel() == nil {
		t.Errorf("bLevelParams must default to a new BLevelParameters")
	}
}

// TestAbstractSerializableSignatureParametersEqualsMatchesJava pins
// AbstractSerializableSignatureParameters#equals field by field, including the
// fields upstream deliberately leaves out.
func TestAbstractSerializableSignatureParametersEqualsMatchesJava(t *testing.T) {
	if !abstractSerializableSignatureParametersNew(t).Equals(abstractSerializableSignatureParametersNew(t)) {
		t.Fatalf("two freshly constructed parameter sets must be equal")
	}

	compared := []struct {
		name   string
		mutate func(*AbstractSerializableSignatureParameters[*TimestampParameters])
	}{
		{"checkCertificateRevocation", func(p *AbstractSerializableSignatureParameters[*TimestampParameters]) {
			p.SetCheckCertificateRevocation(true)
		}},
		{"generateTBSWithoutCertificate", func(p *AbstractSerializableSignatureParameters[*TimestampParameters]) {
			p.SetGenerateTBSWithoutCertificate(true)
		}},
		{"signatureLevel", func(p *AbstractSerializableSignatureParameters[*TimestampParameters]) {
			p.SetSignatureLevel(enumerations.SignatureLevel_XAdES_BASELINE_B)
		}},
		{"signaturePackaging", func(p *AbstractSerializableSignatureParameters[*TimestampParameters]) {
			p.SetSignaturePackaging(enumerations.SignaturePackaging_ENVELOPED)
		}},
		{"digestAlgorithm", func(p *AbstractSerializableSignatureParameters[*TimestampParameters]) {
			p.SetDigestAlgorithm(enumerations.DigestAlgorithm_SHA256)
		}},
		{"referenceDigestAlgorithm", func(p *AbstractSerializableSignatureParameters[*TimestampParameters]) {
			p.SetReferenceDigestAlgorithm(enumerations.DigestAlgorithm_SHA512)
		}},
		{"bLevelParams", func(p *AbstractSerializableSignatureParameters[*TimestampParameters]) {
			bLevel := NewBLevelParameters()
			bLevel.SetSigningDate(&abstractSerializableSignatureParametersSigningDate)
			bLevel.SetClaimedSignerRoles([]string{"role"})
			p.SetBLevelParams(bLevel)
		}},
		{"contentTimestampParameters", func(p *AbstractSerializableSignatureParameters[*TimestampParameters]) {
			tsp := NewTimestampParametersWithDigestAlgorithm(enumerations.DigestAlgorithm_SHA512)
			p.SetContentTimestampParameters(&tsp)
		}},
	}
	for _, tc := range compared {
		t.Run(tc.name+" is compared", func(t *testing.T) {
			mutated := abstractSerializableSignatureParametersNew(t)
			tc.mutate(mutated)
			if mutated.Equals(abstractSerializableSignatureParametersNew(t)) {
				t.Errorf("%s must take part in Equals", tc.name)
			}
		})
	}

	// Upstream leaves validationDataEncapsulationStrategy out of equals(),
	// hashCode() and toString(); the omission is reproduced, not "fixed".
	ignored := abstractSerializableSignatureParametersNew(t)
	ignored.SetValidationDataEncapsulationStrategy(enumerations.ValidationDataEncapsulationStrategy_ANY_VALIDATION_DATA_ONLY)
	if !ignored.Equals(abstractSerializableSignatureParametersNew(t)) {
		t.Errorf("validationDataEncapsulationStrategy must NOT take part in Equals (upstream omits it)")
	}
	if got := ignored.String(); got != abstractSerializableSignatureParametersNew(t).String() {
		t.Errorf("validationDataEncapsulationStrategy must NOT appear in String() (upstream omits it)")
	}

	p := abstractSerializableSignatureParametersNew(t)
	if p.Equals(nil) {
		t.Errorf("a parameter set must never equal nil")
	}
	if !p.Equals(p) {
		t.Errorf("a parameter set must equal itself")
	}
}
