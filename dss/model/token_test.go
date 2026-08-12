package model

import (
	"testing"

	"github.com/utain/esig/dss/enumerations"
)

func TestTokenBaseDefaults(t *testing.T) {
	base := NewTokenBase()

	// The Java field initialisers.
	if got := base.SignatureValidity(); got != enumerations.SignatureValidity_NOT_EVALUATED {
		t.Errorf("SignatureValidity() = %v, want NOT_EVALUATED", got)
	}
	if got := base.InvalidityReason(); got != "" {
		t.Errorf("InvalidityReason() = %q, want empty", got)
	}
	if base.IsSelfSigned() {
		t.Error("a plain token is never self-signed")
	}
	if got := base.Abbreviation(); got != "?" {
		t.Errorf("Abbreviation() = %q, want \"?\"", got)
	}
	if base.PublicKeyOfTheSigner() != nil {
		t.Error("PublicKeyOfTheSigner() starts nil")
	}
	if base.IsSignatureIntact() || base.IsValid() {
		t.Error("a token with NOT_EVALUATED validity is not intact")
	}
	if got := base.SignatureAlgorithm(); got != "" {
		t.Errorf("SignatureAlgorithm() = %v, want the zero value", got)
	}
}

func TestTokenBaseProtectedFieldSetters(t *testing.T) {
	// Java subclasses assign the protected fields directly; Go subclasses in other packages
	// go through these setters instead.
	base := NewTokenBase()

	base.SetSignatureAlgorithm(enumerations.SignatureAlgorithm_RSA_SHA512)
	if got := base.SignatureAlgorithm(); got != enumerations.SignatureAlgorithm_RSA_SHA512 {
		t.Errorf("SignatureAlgorithm() = %v", got)
	}
	base.SetSignatureValidity(enumerations.SignatureValidity_VALID)
	if !base.IsSignatureIntact() || !base.IsValid() {
		t.Error("a VALID signature must be reported as intact")
	}
	base.SetInvalidityReason("SignatureException : boom")
	if got := base.InvalidityReason(); got != "SignatureException : boom" {
		t.Errorf("InvalidityReason() = %q", got)
	}
	key := NewPublicKeyFromEncoded([]byte("spki"), nil)
	base.SetPublicKeyOfTheSigner(key)
	if base.PublicKeyOfTheSigner() != key {
		t.Error("SetPublicKeyOfTheSigner did not take effect")
	}
}

func TestTokenBasePanicsWhenNotInitialised(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("a token that never called InitToken must panic on dispatch")
		}
	}()
	base := NewTokenBase()
	base.DSSIDAsString()
}

func TestTokenBaseIssuerEntityKeyNeedsASigner(t *testing.T) {
	leaf := certificateTokenFixture(t, leafCertificateBase64)
	// getIssuerEntityKey() returns null until checkIsSignedBy has established the signer.
	if leaf.TokenBase.IssuerEntityKey() != nil {
		t.Error("IssuerEntityKey() must be nil without a signer")
	}
	root := certificateTokenFixture(t, rootCertificateBase64)
	leaf.SetPublicKeyOfTheSigner(root.PublicKey())
	if got, want := leaf.TokenBase.IssuerEntityKey().AsXmlID(), root.EntityKey().AsXmlID(); got != want {
		t.Errorf("IssuerEntityKey() = %q, want %q", got, want)
	}
}

func TestTokenInterfaceIsSatisfied(t *testing.T) {
	// The Token interface is what TokenComparator and the validation engine consume.
	var token Token = certificateTokenFixture(t, rootCertificateBase64)
	if token.DSSIDAsString() == "" {
		t.Error("DSSIDAsString() must not be empty")
	}
	if token.IssuerX500Principal() == nil {
		t.Error("IssuerX500Principal() must not be nil")
	}
	if token.CreationDate().IsZero() {
		t.Error("CreationDate() must not be the zero time")
	}
	if len(token.Encoded()) == 0 {
		t.Error("Encoded() must not be empty")
	}
	if token.ToString("") != token.String() {
		t.Error("String() must be ToString(\"\")")
	}
}
