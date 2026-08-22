// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/x509/ResponderId.java (DSS 6.5.RC1).
package spi

import (
	"bytes"

	"github.com/ryftcore/dss-go/dss/model"
)

// ResponderId represents a ResponderId extracted from an OCSP response.
type ResponderId struct {
	// subjectX500Principal is the X500Principal of the OCSP issuer.
	subjectX500Principal *model.X500Principal

	// ski is the SKI of the OCSP issuer.
	ski []byte
}

// NewResponderId builds the ResponderId. Port of the default constructor.
func NewResponderId(subjectX500Principal *model.X500Principal, ski []byte) *ResponderId {
	return &ResponderId{subjectX500Principal: subjectX500Principal, ski: ski}
}

// X500Principal gets the X500Principal of the OCSP issuer. Port of getX500Principal().
func (r *ResponderId) X500Principal() *model.X500Principal {
	return r.subjectX500Principal
}

// SetX500Principal sets the X500Principal of the OCSP issuer. Port of setX500Principal(X500Principal).
func (r *ResponderId) SetX500Principal(subjectX500Principal *model.X500Principal) {
	r.subjectX500Principal = subjectX500Principal
}

// Ski gets the SKI of the issuer. Port of getSki().
func (r *ResponderId) Ski() []byte {
	return r.ski
}

// SetSki sets the SKI of the issuer. Port of setSki(byte[]).
func (r *ResponderId) SetSki(ski []byte) {
	r.ski = ski
}

// IsRelatedToCertificate checks if the ResponderId is related to the given certificateToken.
// Port of isRelatedToCertificate(CertificateToken).
func (r *ResponderId) IsRelatedToCertificate(certificateToken *model.CertificateToken) bool {
	if r.subjectX500Principal != nil {
		return DSSASN1UtilsX500PrincipalAreEquals(certificateToken.Subject().Principal(), r.subjectX500Principal)
	}
	return DSSASN1UtilsIsSkiEqual(r.ski, certificateToken)
}

// Equals reports whether both ResponderIds carry the same SKI and X500Principal.
// Port of equals(Object). Java's paired hashCode() has no Go counterpart.
func (r *ResponderId) Equals(other *ResponderId) bool {
	if r == other {
		return true
	}
	if other == nil {
		return false
	}
	if !bytes.Equal(r.ski, other.ski) {
		return false
	}
	if r.subjectX500Principal == nil {
		return other.subjectX500Principal == nil
	}
	return r.subjectX500Principal.Equals(other.subjectX500Principal)
}
