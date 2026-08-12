package model

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"

	"github.com/utain/esig/dss/model/x509/revocation"
)

func TestEncapsulatedRevocationTokenIdentifier(t *testing.T) {
	binaries := []byte("revocation binaries")
	identifier := NewEncapsulatedRevocationTokenIdentifier[revocation.CRL](binaries)

	want := sha256.Sum256(binaries)
	if got := hex.EncodeToString(identifier.DigestID().Value()); got != hex.EncodeToString(want[:]) {
		t.Errorf("DigestID().Value() = %s", got)
	}
	if got := identifier.AsXmlID()[:2]; got != "R-" {
		t.Errorf("AsXmlID() prefix = %q", got)
	}
	if got := string(identifier.Binaries()); got != string(binaries) {
		t.Errorf("Binaries() = %q", got)
	}
	// getDSSId() returns the identifier itself.
	if identifier.DSSID() != Identifier(identifier) {
		t.Error("DSSID() must return the identifier itself")
	}
	if got := identifier.String(); got != "EncapsulatedRevocationTokenIdentifier:SHA256:#"+identifier.AsXmlID()[2:] {
		t.Errorf("String() = %q", got)
	}
}

func TestEncapsulatedRevocationTokenIdentifierSubclassClassName(t *testing.T) {
	// A subclass in another package keeps its own Java simple class name, which drives
	// toString() and the class check in equals().
	crlBinary := NewEncapsulatedRevocationTokenIdentifierWithClassName[revocation.CRL]("CRLBinary", []byte("x"))
	ocspBinary := NewEncapsulatedRevocationTokenIdentifierWithClassName[revocation.OCSP]("OCSPResponseBinary", []byte("x"))

	if got := crlBinary.String()[:10]; got != "CRLBinary:" {
		t.Errorf("String() = %q", crlBinary.String())
	}
	if crlBinary.Equals(ocspBinary) {
		t.Error("identifiers with different class names must not be equal")
	}
	if !crlBinary.Equals(NewEncapsulatedRevocationTokenIdentifierWithClassName[revocation.CRL]("CRLBinary", []byte("x"))) {
		t.Error("identifiers with the same class name and binaries must be equal")
	}
}
