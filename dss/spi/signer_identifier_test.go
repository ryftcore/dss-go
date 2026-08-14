package spi

import (
	"fmt"
	"math/big"
	"sort"
	"strings"
	"testing"

	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/utils"
)

// TestSignerIdentifierKnownAnswers checks the DER IssuerSerial, the string form and the flattened
// attribute map against the BouncyCastle/JDK known answers; see
// certificate_extensions_utils_kat_test.go for how the fixture is generated.
//
// The fixture's signerId.toString uses X500Principal#getName(RFC2253) rather than
// X500Principal#toString(), matching the documented deviation of model.X500Principal.
func TestSignerIdentifierKnownAnswers(t *testing.T) {
	for index, answers := range certificateExtensionsKATLoad(t) {
		name := answers.str(t, "file")
		if certificateExtensionsKATUnparseable[name] {
			continue
		}
		t.Run(fmt.Sprintf("%02d_%s", index, name), func(t *testing.T) {
			certificateToken := certificateExtensionsKATToken(t, name)
			signerIdentifier := NewSignerIdentifier()
			signerIdentifier.SetIssuerName(certificateToken.IssuerX500Principal())
			signerIdentifier.SetSerialNumber(certificateToken.SerialNumber())

			certificateExtensionsKATAssertHex(t, "signerId.issuerSerial", signerIdentifier.IssuerSerialEncoded(), answers)
			if got, want := signerIdentifier.String(), answers.str(t, "signerId.toString"); got != want {
				t.Errorf("SignerIdentifier.String() = %q, want %q", got, want)
			}

			attributes := DSSASN1UtilsAttributeMap(certificateToken.IssuerX500Principal())
			if attributes == nil {
				t.Fatalf("DSSASN1UtilsAttributeMap failed")
			}
			flattened := make([]string, 0, len(attributes))
			for attributeType, value := range attributes {
				flattened = append(flattened, attributeType+"="+value)
			}
			sort.Strings(flattened)
			want := answers.list(t, "signerId.issuerMap")
			if strings.Join(flattened, "|") != strings.Join(want, "|") {
				t.Errorf("issuer attribute map = %v, want %v", flattened, want)
			}
		})
	}
}

// signerIdentifierTestTokens answers two distinct certificate tokens from the fixture corpus, the
// first of which carries a subjectKeyIdentifier extension.
func signerIdentifierTestTokens(t *testing.T) (*model.CertificateToken, *model.CertificateToken) {
	t.Helper()
	var withSKI, other *model.CertificateToken
	for _, answers := range certificateExtensionsKATLoad(t) {
		name := answers.str(t, "file")
		if certificateExtensionsKATUnparseable[name] {
			continue
		}
		token := certificateExtensionsKATToken(t, name)
		if withSKI == nil && !answers.isNull(t, "ski") {
			withSKI = token
			continue
		}
		if other == nil {
			other = token
		}
	}
	if withSKI == nil || other == nil {
		t.Fatalf("the fixture corpus no longer covers a certificate with a subjectKeyIdentifier")
	}
	return withSKI, other
}

func TestSignerIdentifierIsEmpty(t *testing.T) {
	signerIdentifier := NewSignerIdentifier()
	if !signerIdentifier.IsEmpty() {
		t.Errorf("IsEmpty() = false, want true for a fresh SignerIdentifier")
	}
	if signerIdentifier.IssuerSerialEncoded() != nil {
		t.Errorf("IssuerSerialEncoded() = non-nil, want nil without an issuer and a serial number")
	}
	if got, want := signerIdentifier.String(), "IssuerSerialInfo [ski=]"; got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}

	signerIdentifier.SetSki([]byte{1, 2, 3})
	if signerIdentifier.IsEmpty() {
		t.Errorf("IsEmpty() = true, want false once a SKI is set")
	}
	if got, want := signerIdentifier.String(), "IssuerSerialInfo [ski="+utils.ToBase64([]byte{1, 2, 3})+"]"; got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}

	signerIdentifier = NewSignerIdentifier()
	signerIdentifier.SetSerialNumber(big.NewInt(42))
	if signerIdentifier.IsEmpty() {
		t.Errorf("IsEmpty() = true, want false once a serial number is set")
	}
	if signerIdentifier.IssuerSerialEncoded() != nil {
		t.Errorf("IssuerSerialEncoded() = non-nil, want nil without an issuer name")
	}
	if got, want := signerIdentifier.String(), "IssuerSerialInfo [issuerName=null, serialNumber=42]"; got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}

func TestSignerIdentifierIsRelatedToCertificate(t *testing.T) {
	withSKI, other := signerIdentifierTestTokens(t)

	byIssuerAndSerial := NewSignerIdentifier()
	byIssuerAndSerial.SetIssuerName(withSKI.IssuerX500Principal())
	byIssuerAndSerial.SetSerialNumber(withSKI.SerialNumber())
	related, err := byIssuerAndSerial.IsRelatedToCertificate(withSKI)
	if err != nil {
		t.Fatalf("IsRelatedToCertificate: unexpected error %v", err)
	}
	if !related {
		t.Errorf("IsRelatedToCertificate() = false, want true for the issuer/serial pair of the token")
	}
	related, err = byIssuerAndSerial.IsRelatedToCertificate(other)
	if err != nil {
		t.Fatalf("IsRelatedToCertificate: unexpected error %v", err)
	}
	if related {
		t.Errorf("IsRelatedToCertificate() = true, want false for a different certificate")
	}

	// With neither an issuer nor a serial number the comparison falls back to the SKI.
	ski, err := CertificateExtensionsUtilsSubjectKeyIdentifier(withSKI)
	if err != nil {
		t.Fatalf("SubjectKeyIdentifier: unexpected error %v", err)
	}
	bySKI := NewSignerIdentifier()
	bySKI.SetSki(ski.Ski())
	related, err = bySKI.IsRelatedToCertificate(withSKI)
	if err != nil {
		t.Fatalf("IsRelatedToCertificate: unexpected error %v", err)
	}
	if !related {
		t.Errorf("IsRelatedToCertificate() = false, want true for the SKI of the token")
	}

	// A serial number without an issuer name also takes the SKI branch.
	partial := NewSignerIdentifier()
	partial.SetSerialNumber(withSKI.SerialNumber())
	partial.SetSki(ski.Ski())
	related, err = partial.IsRelatedToCertificate(withSKI)
	if err != nil {
		t.Fatalf("IsRelatedToCertificate: unexpected error %v", err)
	}
	if !related {
		t.Errorf("IsRelatedToCertificate() = false, want true when only the SKI can be matched")
	}
}

func TestSignerIdentifierEqualsAndEquivalence(t *testing.T) {
	withSKI, other := signerIdentifierTestTokens(t)

	first := NewSignerIdentifier()
	first.SetIssuerName(withSKI.IssuerX500Principal())
	first.SetSerialNumber(withSKI.SerialNumber())
	second := NewSignerIdentifier()
	second.SetIssuerName(withSKI.IssuerX500Principal())
	second.SetSerialNumber(withSKI.SerialNumber())

	if !first.Equals(second) {
		t.Errorf("Equals() = false, want true for two identifiers built from the same certificate")
	}
	if !first.IsEquivalent(second) {
		t.Errorf("IsEquivalent() = false, want true for two identifiers built from the same certificate")
	}
	if first.Equals(nil) {
		t.Errorf("Equals(nil) = true, want false")
	}
	if first.IsEquivalent(nil) {
		t.Errorf("IsEquivalent(nil) = true, want false")
	}

	second.SetSerialNumber(new(big.Int).Add(withSKI.SerialNumber(), big.NewInt(1)))
	if first.Equals(second) {
		t.Errorf("Equals() = true, want false for a different serial number")
	}
	if first.IsEquivalent(second) {
		t.Errorf("IsEquivalent() = true, want false for a different serial number")
	}

	// A different issuer name breaks the equivalence even when the serial number matches.
	third := NewSignerIdentifier()
	third.SetIssuerName(other.IssuerX500Principal())
	third.SetSerialNumber(withSKI.SerialNumber())
	if first.IssuerName().Equals(other.IssuerX500Principal()) {
		t.Fatalf("the two fixture certificates unexpectedly share an issuer name")
	}
	if first.IsEquivalent(third) {
		t.Errorf("IsEquivalent() = true, want false for a different issuer name")
	}

	// The SKI branch.
	fourth := NewSignerIdentifier()
	fourth.SetSki([]byte{1, 2, 3})
	fifth := NewSignerIdentifier()
	fifth.SetSki([]byte{1, 2, 3})
	if !fourth.IsEquivalent(fifth) || !fourth.Equals(fifth) {
		t.Errorf("two identifiers sharing a SKI must be equivalent and equal")
	}
	fifth.SetSki([]byte{4, 5, 6})
	if fourth.IsEquivalent(fifth) || fourth.Equals(fifth) {
		t.Errorf("two identifiers with different SKIs must be neither equivalent nor equal")
	}
}

func TestSignerIdentifierCurrentFlag(t *testing.T) {
	signerIdentifier := NewSignerIdentifier()
	if signerIdentifier.IsCurrent() {
		t.Errorf("IsCurrent() = true, want false by default")
	}
	signerIdentifier.SetCurrent(true)
	if !signerIdentifier.IsCurrent() {
		t.Errorf("IsCurrent() = false after SetCurrent(true)")
	}
}
