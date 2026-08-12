package crlparser

import (
	"encoding/asn1"
	"io"
	"testing"
)

func TestNewCRLValidity_NilBinaryPanics(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatalf("expected a panic for a nil CRLBinary")
		}
		if r != "CRLBinary cannot be null!" {
			t.Errorf("panic value = %v, want the Java message", r)
		}
	}()
	NewCRLValidity(nil)
}

func TestCRLValidity_CrlBinaryAndDerEncoded(t *testing.T) {
	der := []byte{0x30, 0x03, 0x02, 0x01, 0x01}
	binary := NewCRLBinary(der)
	validity := NewCRLValidity(binary)

	if validity.CrlBinary() != binary {
		t.Errorf("CrlBinary() did not return the constructor argument")
	}
	if string(validity.DerEncoded()) != string(der) {
		t.Errorf("DerEncoded() = %v, want %v", validity.DerEncoded(), der)
	}
	read, err := io.ReadAll(validity.ToCRLInputStream())
	if err != nil {
		t.Fatalf("reading ToCRLInputStream: %v", err)
	}
	if string(read) != string(der) {
		t.Errorf("ToCRLInputStream() produced %v, want %v", read, der)
	}
}

// TestCRLValidity_IsValid_IsUnknownCriticalExtension exercises the boolean combinations from
// CRLValidity.java's isValid()/isUnknownCriticalExtension(): the four "coherence" flags, and
// the three ways a critical extension set can be flagged "unknown".
func TestCRLValidity_IsValid_IsUnknownCriticalExtension(t *testing.T) {
	newBase := func() *CRLValidity {
		v := NewCRLValidity(NewCRLBinary([]byte{0x30, 0x00}))
		v.SetIssuerX509PrincipalMatches(true)
		v.SetSignatureIntact(true)
		v.SetCrlSignKeyUsage(true)
		v.SetURL("http://example.com/a.crl")
		return v
	}

	t.Run("fully coherent, no critical extensions", func(t *testing.T) {
		v := newBase()
		if v.IsUnknownCriticalExtension() {
			t.Errorf("expected no unknown critical extension")
		}
		if !v.IsValid() {
			t.Errorf("expected IsValid() to be true")
		}
	})

	t.Run("critical extensions but URL set and reasonFlags nil and not every onlyXxx flag", func(t *testing.T) {
		v := newBase()
		v.SetCriticalExtensionsOid([]string{"2.5.29.28"})
		if v.IsUnknownCriticalExtension() {
			t.Errorf("expected no unknown critical extension (only URL/reasonFlags/onlyXxx combo trip it)")
		}
		if !v.IsValid() {
			t.Errorf("expected IsValid() to be true")
		}
	})

	t.Run("critical extensions and no URL", func(t *testing.T) {
		v := newBase()
		v.SetURL("")
		v.SetCriticalExtensionsOid([]string{"2.5.29.28"})
		if !v.IsUnknownCriticalExtension() {
			t.Errorf("expected an unknown critical extension when the URL is unset")
		}
		if v.IsValid() {
			t.Errorf("expected IsValid() to be false")
		}
	})

	t.Run("critical extensions and reasonFlags present", func(t *testing.T) {
		v := newBase()
		v.SetReasonFlags(&asn1.BitString{Bytes: []byte{0x80}, BitLength: 1})
		v.SetCriticalExtensionsOid([]string{"2.5.29.28"})
		if !v.IsUnknownCriticalExtension() {
			t.Errorf("expected an unknown critical extension when reasonFlags is present")
		}
	})

	t.Run("critical extensions and every onlyXxx flag set", func(t *testing.T) {
		v := newBase()
		v.SetOnlyUserCerts(true)
		v.SetOnlyCaCerts(true)
		v.SetOnlyAttributeCerts(true)
		v.SetIndirectCrl(true)
		v.SetCriticalExtensionsOid([]string{"2.5.29.28"})
		if !v.IsUnknownCriticalExtension() {
			t.Errorf("expected an unknown critical extension when every onlyXxx flag is set")
		}
	})

	t.Run("not coherent: issuer mismatch", func(t *testing.T) {
		v := newBase()
		v.SetIssuerX509PrincipalMatches(false)
		if v.IsValid() {
			t.Errorf("expected IsValid() to be false")
		}
	})

	t.Run("not coherent: signature not intact", func(t *testing.T) {
		v := newBase()
		v.SetSignatureIntact(false)
		if v.IsValid() {
			t.Errorf("expected IsValid() to be false")
		}
	})

	t.Run("not coherent: no cRLSign key usage", func(t *testing.T) {
		v := newBase()
		v.SetCrlSignKeyUsage(false)
		if v.IsValid() {
			t.Errorf("expected IsValid() to be false")
		}
	})
}

func TestCRLValidity_Equals(t *testing.T) {
	der := []byte{0x30, 0x03, 0x02, 0x01, 0x01}

	build := func() *CRLValidity {
		v := NewCRLValidity(NewCRLBinary(der))
		v.SetURL("http://example.com/a.crl")
		v.SetSignatureIntact(true)
		v.SetCriticalExtensionsOid([]string{"2.5.29.28"})
		return v
	}

	a := build()
	b := build()
	if !a.Equals(b) {
		t.Errorf("expected two independently built, identically configured CRLValidity values to be equal")
	}
	if !a.Equals(a) {
		t.Errorf("expected a value to equal itself")
	}
	if a.Equals(nil) {
		t.Errorf("expected a value to not equal nil")
	}

	c := build()
	c.SetURL("http://example.com/different.crl")
	if a.Equals(c) {
		t.Errorf("expected values with different URLs to not be equal")
	}
}

func TestCRLValidity_String(t *testing.T) {
	v := NewCRLValidity(NewCRLBinary([]byte{0x30, 0x00}))
	// Just check it does not panic and mentions the DSS Id, matching toString()'s shape.
	s := v.String()
	if s == "" {
		t.Errorf("expected a non-empty String()")
	}
}

func TestCRLValidity_SetReasonFlagsAndOnlyXxxHaveNoGetters(t *testing.T) {
	// Compile-time-ish smoke test documenting the deliberate asymmetry with upstream: these
	// setters exist, but (like the Java class) there is no exported getter, since only
	// IsUnknownCriticalExtension reads them, from within the package.
	v := NewCRLValidity(NewCRLBinary([]byte{0x30, 0x00}))
	v.SetOnlyUserCerts(true)
	v.SetOnlyCaCerts(true)
	v.SetOnlyAttributeCerts(true)
	v.SetIndirectCrl(true)
	v.SetReasonFlags(&asn1.BitString{Bytes: []byte{0x80}, BitLength: 1})
}
