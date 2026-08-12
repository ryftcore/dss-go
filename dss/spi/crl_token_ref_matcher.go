// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/x509/revocation/crl/CRLTokenRefMatcher.java (DSS 6.5.RC1).
package spi

import (
	"bytes"

	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/model/x509/revocation"
)

// CRLTokenRefMatcher matches a CRL with a reference.
type CRLTokenRefMatcher struct{}

// NewCRLTokenRefMatcher builds the matcher. Port of the default constructor.
func NewCRLTokenRefMatcher() *CRLTokenRefMatcher {
	return &CRLTokenRefMatcher{}
}

// Match reports whether the reference refers to the given token.
// Port of match(RevocationToken, RevocationRef).
//
// Java's (CRLToken) and (CRLRef) casts are kept as type assertions: a mismatch is a
// programmer error and panics, as Java's ClassCastException does.
func (m *CRLTokenRefMatcher) Match(token RevocationToken[revocation.CRL],
	reference RevocationRef[revocation.CRL]) (bool, error) {
	crlToken, ok := token.(*CRLToken)
	if !ok {
		panic("ClassCastException : the revocation token is not a CRLToken")
	}
	crlRef, ok := reference.(*CRLRef)
	if !ok {
		panic("ClassCastException : the revocation reference is not a CRLRef")
	}

	if crlRef.Digest().Value() != nil {
		return crlTokenRefMatcherMatchByDigest(crlToken, crlRef.Digest())
	}
	return false, model.NewDSSError("Digest is mandatory for comparison")
}

// MatchBinary reports whether the reference refers to the given encapsulated identifier.
// Port of the match(EncapsulatedRevocationTokenIdentifier, RevocationRef) overload, which Go
// cannot give the same name.
func (m *CRLTokenRefMatcher) MatchBinary(identifier EncapsulatedRevocationTokenIdentifier[revocation.CRL],
	reference RevocationRef[revocation.CRL]) (bool, error) {
	crlRef, ok := reference.(*CRLRef)
	if !ok {
		panic("ClassCastException : the revocation reference is not a CRLRef")
	}

	if crlRef.Digest().Value() != nil {
		return crlTokenRefMatcherMatchBinaryByDigest(identifier, crlRef.Digest())
	}
	return false, model.NewDSSError("Digest is mandatory for comparison")
}

// crlTokenRefMatcherMatchByDigest ports the private matchByDigest(RevocationToken, Digest).
func crlTokenRefMatcherMatchByDigest(token RevocationToken[revocation.CRL], digest model.Digest) (bool, error) {
	value, err := token.Digest(digest.Algorithm())
	if err != nil {
		return false, err
	}
	return bytes.Equal(digest.Value(), value), nil
}

// crlTokenRefMatcherMatchBinaryByDigest ports the private
// matchByDigest(EncapsulatedRevocationTokenIdentifier, Digest).
func crlTokenRefMatcherMatchBinaryByDigest(identifier EncapsulatedRevocationTokenIdentifier[revocation.CRL],
	digest model.Digest) (bool, error) {
	value, err := identifier.DigestValue(digest.Algorithm())
	if err != nil {
		return false, err
	}
	return bytes.Equal(digest.Value(), value), nil
}

// compile-time assertion: a CRLTokenRefMatcher matches CRL references.
var _ RevocationTokenRefMatcher[revocation.CRL] = (*CRLTokenRefMatcher)(nil)
