// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/x509/CertificateTokenRefMatcher.java (DSS 6.5.RC1).
//
// eu.europa.esig.dss.spi.x509 flattens into the Go package spi (see PORTING_PLAN.md), so the
// type keeps its Java name unqualified.
package spi

import (
	"bytes"

	"github.com/ryftcore/dss-go/dss/model"
)

// CertificateTokenRefMatcher is used to verify if a given CertificateToken matches a
// CertificateRef.
type CertificateTokenRefMatcher struct{}

// NewCertificateTokenRefMatcher builds a matcher. Port of the default constructor.
func NewCertificateTokenRefMatcher() *CertificateTokenRefMatcher {
	return &CertificateTokenRefMatcher{}
}

// Match verifies if the given CertificateToken matches the CertificateRef.
// Port of match(CertificateToken, CertificateRef).
//
// SignerIdentifier.IsRelatedToCertificate returns an error when a certificate extension
// cannot be read; Java lets the underlying DSSException propagate unchecked from that call,
// so a matching failure here is treated as "does not match" rather than surfaced separately -
// callers that need the failure reason should call it directly. ResponderId.IsRelatedToCertificate
// has no such failure mode (it never reads a fallible extension) and returns a plain bool,
// matching its Java signature exactly.
func (m *CertificateTokenRefMatcher) Match(certificateToken *model.CertificateToken, certificateRef *CertificateRef) bool {
	// If we only know the public key, the token is null
	if certificateToken == nil {
		return false
	}

	certDigest := certificateRef.CertDigest()
	signerIdentifier := certificateRef.CertificateIdentifier()
	responderId := certificateRef.ResponderId()
	publicKey := certificateRef.PublicKey()
	if !certDigest.IsEmpty() && m.MatchByDigest(certificateToken, certificateRef) {
		return true
	} else if signerIdentifier != nil {
		if related, err := signerIdentifier.IsRelatedToCertificate(certificateToken); err == nil && related {
			return true
		}
	} else if responderId != nil {
		if responderId.IsRelatedToCertificate(certificateToken) {
			return true
		}
	} else if publicKey != nil && publicKey.Equals(certificateToken.PublicKey()) {
		return true
	}
	return false
}

// MatchByDigest verifies if only the digest within the certificateRef corresponds to
// certificateToken. Port of matchByDigest(CertificateToken, CertificateRef).
func (m *CertificateTokenRefMatcher) MatchByDigest(certificateToken *model.CertificateToken, certificateRef *CertificateRef) bool {
	certDigest := certificateRef.CertDigest()
	if !certDigest.IsEmpty() {
		currentDigest, err := certificateToken.Digest(certDigest.Algorithm())
		if err != nil {
			return false
		}
		return bytes.Equal(currentDigest, certDigest.Value())
	}
	return false
}

// MatchBySerialNumber verifies if only the serial number within the certificateRef
// corresponds to certificateToken. Port of matchBySerialNumber(CertificateToken, CertificateRef).
func (m *CertificateTokenRefMatcher) MatchBySerialNumber(certificateToken *model.CertificateToken, certificateRef *CertificateRef) bool {
	signerIdentifier := certificateRef.CertificateIdentifier()
	if signerIdentifier != nil && signerIdentifier.SerialNumber() != nil {
		return certificateToken.SerialNumber().Cmp(signerIdentifier.SerialNumber()) == 0
	}
	return false
}

// MatchByIssuerName verifies if only the issuer name within the certificateRef corresponds to
// certificateToken. Port of matchByIssuerName(CertificateToken, CertificateRef).
func (m *CertificateTokenRefMatcher) MatchByIssuerName(certificateToken *model.CertificateToken, certificateRef *CertificateRef) bool {
	signerIdentifier := certificateRef.CertificateIdentifier()
	if signerIdentifier != nil && signerIdentifier.IssuerName() != nil {
		return DSSASN1UtilsX500PrincipalAreEquals(signerIdentifier.IssuerName(), certificateToken.IssuerX500Principal())
	}
	return false
}

// MatchByResponderId verifies if only the responder Id within the certificateRef corresponds
// to certificateToken. Port of matchByResponderId(CertificateToken, CertificateRef).
func (m *CertificateTokenRefMatcher) MatchByResponderId(certificateToken *model.CertificateToken, certificateRef *CertificateRef) bool {
	responderId := certificateRef.ResponderId()
	if responderId != nil {
		return responderId.IsRelatedToCertificate(certificateToken)
	}
	return false
}
