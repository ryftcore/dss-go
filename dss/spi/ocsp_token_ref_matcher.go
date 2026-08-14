// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/x509/revocation/ocsp/OCSPTokenRefMatcher.java (DSS 6.5.RC1).
package spi

import (
	"bytes"

	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/model/x509/revocation"
)

// OCSPTokenRefMatcher checks an OCSP token reference against a token or its binaries.
type OCSPTokenRefMatcher struct{}

// NewOCSPTokenRefMatcher builds the matcher. Port of the default constructor.
func NewOCSPTokenRefMatcher() *OCSPTokenRefMatcher {
	return &OCSPTokenRefMatcher{}
}

// Match reports whether the reference refers to the given token.
// Port of match(RevocationToken, RevocationRef).
//
// Java's (OCSPToken) and (OCSPRef) casts are kept as type assertions: a mismatch is a
// programmer error and panics, as Java's ClassCastException does.
func (m *OCSPTokenRefMatcher) Match(token RevocationToken[revocation.OCSP],
	reference RevocationRef[revocation.OCSP]) (bool, error) {
	ocspToken, ok := token.(*OCSPToken)
	if !ok {
		panic("ClassCastException : the revocation token is not an OCSPToken")
	}
	ocspRef, ok := reference.(*OCSPRef)
	if !ok {
		panic("ClassCastException : the revocation reference is not an OCSPRef")
	}

	if ocspRef.Digest().Value() != nil {
		return ocspTokenRefMatcherMatchByDigest(ocspToken, ocspRef.Digest())
	}
	return ocspTokenRefMatcherMatchResponse(ocspToken.BasicOCSPResp(), ocspRef), nil
}

// MatchBinary reports whether the reference refers to the given encapsulated identifier.
// Port of the match(EncapsulatedRevocationTokenIdentifier, RevocationRef) overload, which Go
// cannot give the same name.
func (m *OCSPTokenRefMatcher) MatchBinary(identifier EncapsulatedRevocationTokenIdentifier[revocation.OCSP],
	reference RevocationRef[revocation.OCSP]) (bool, error) {
	ocspResponseBinary, ok := identifier.(*OCSPResponseBinary)
	if !ok {
		panic("ClassCastException : the revocation identifier is not an OCSPResponseBinary")
	}
	ocspRef, ok := reference.(*OCSPRef)
	if !ok {
		panic("ClassCastException : the revocation reference is not an OCSPRef")
	}

	if ocspRef.Digest().Value() != nil {
		return ocspTokenRefMatcherMatchBinaryByDigest(ocspResponseBinary, ocspRef.Digest())
	}
	return ocspTokenRefMatcherMatchResponse(ocspResponseBinary.BasicOCSPResp(), ocspRef), nil
}

// ocspTokenRefMatcherMatchByDigest ports the private matchByDigest(RevocationToken, Digest).
func ocspTokenRefMatcherMatchByDigest(token RevocationToken[revocation.OCSP], digest model.Digest) (bool, error) {
	value, err := token.Digest(digest.Algorithm())
	if err != nil {
		return false, err
	}
	return bytes.Equal(digest.Value(), value), nil
}

// ocspTokenRefMatcherMatchBinaryByDigest ports the private
// matchByDigest(EncapsulatedRevocationTokenIdentifier, Digest).
func ocspTokenRefMatcherMatchBinaryByDigest(identifier EncapsulatedRevocationTokenIdentifier[revocation.OCSP],
	digest model.Digest) (bool, error) {
	value, err := identifier.DigestValue(digest.Algorithm())
	if err != nil {
		return false, err
	}
	return bytes.Equal(digest.Value(), value), nil
}

// ocspTokenRefMatcherMatchResponse ports the private matchResponse(BasicOCSPResp, OCSPRef).
//
// Upstream wraps the comparison in a try/catch that logs and returns false; the Go port
// turns every failure - an unset production time, a responder Name that cannot be decoded -
// into the same false.
func ocspTokenRefMatcherMatchResponse(basicOCSPResp *BasicOCSPResp, ocspRef *OCSPRef) bool {
	if basicOCSPResp == nil || ocspRef.ResponderId() == nil {
		return false
	}
	if ocspRef.ProducedAt().IsZero() || !ocspRef.ProducedAt().Equal(basicOCSPResp.ProducedAt()) {
		return false
	}
	tokenResponderID := basicOCSPResp.ResponderID().ToASN1Primitive()
	refResponderID := ocspRef.ResponderId()
	return ocspTokenRefMatcherMatchByKeyHash(tokenResponderID, refResponderID) ||
		ocspTokenRefMatcherMatchByName(tokenResponderID, refResponderID)
}

// ocspTokenRefMatcherMatchByKeyHash ports the private matchByKeyHash(ResponderID, ResponderId).
func ocspTokenRefMatcherMatchByKeyHash(tokenResponderID *ResponderID, refResponderID *ResponderId) bool {
	return refResponderID.Ski() != nil && bytes.Equal(refResponderID.Ski(), tokenResponderID.KeyHash)
}

// ocspTokenRefMatcherMatchByName ports the private matchByName(ResponderID, ResponderId).
func ocspTokenRefMatcherMatchByName(tokenResponderID *ResponderID, refResponderID *ResponderId) bool {
	if refResponderID.X500Principal() == nil || tokenResponderID.Name == nil {
		return false
	}
	tokenPrincipal, err := DSSASN1UtilsToX500Principal(tokenResponderID.Name)
	if err != nil {
		return false
	}
	return refResponderID.X500Principal().Equals(tokenPrincipal)
}

// compile-time assertion: an OCSPTokenRefMatcher matches OCSP references.
var _ RevocationTokenRefMatcher[revocation.OCSP] = (*OCSPTokenRefMatcher)(nil)
