// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/x509/revocation/RevocationCertificateSource.java (DSS 6.5.RC1).
//
// OCSPCertificateSource (chunk CRLOCSP, a sibling of this phase 2a chunk) already embeds
// RevocationCertificateSourceBase (built with NewRevocationCertificateSourceBase()) and is
// asserted to satisfy the RevocationCertificateSource interface; this file defines both to
// match that shape.
package spi

// RevocationCertificateSource represents a certificate source present into a revocation token.
//
// Java's abstract class adds no members over TokenCertificateSource (its body is `// empty`);
// it exists purely to narrow the upstream type hierarchy
// (CertificateSource <- CommonCertificateSource <- TokenCertificateSource <- RevocationCertificateSource)
// so a revocation token's embedded certificate source can be typed distinctly from a signature's.
// The Go port keeps that distinction as its own named interface, embedding CertificateSource
// unchanged.
type RevocationCertificateSource interface {
	CertificateSource
}

// RevocationCertificateSourceBase carries the state of a RevocationCertificateSource; a
// concrete source (OCSPCertificateSource) embeds it. Port of the protected
// RevocationCertificateSource() constructor - upstream's body is empty, so this simply carries
// TokenCertificateSource's own zero value.
type RevocationCertificateSourceBase struct {
	TokenCertificateSource
}

// NewRevocationCertificateSourceBase builds the base state of a revocation certificate source.
// Calling it is not required - the zero value is already usable, just like
// TokenCertificateSource's - it exists for symmetry with the rest of this port's New.../Init...
// pattern and so a concrete source has an explicit, self-documenting field initializer to embed.
func NewRevocationCertificateSourceBase() RevocationCertificateSourceBase {
	return RevocationCertificateSourceBase{}
}

// compile-time assertion: a RevocationCertificateSourceBase is a RevocationCertificateSource.
var _ RevocationCertificateSource = (*RevocationCertificateSourceBase)(nil)
