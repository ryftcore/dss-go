// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/x509/BaselineBCertificateSelector.java (DSS 6.5.RC1).
//
// eu.europa.esig.dss.spi.x509 flattens into the Go package spi (see PORTING_PLAN.md), so the
// type keeps its Java name unqualified.
//
// ASSUMPTION (flagged for integrator reconciliation, see chunk X509-B which owns
// CertificateReorderer and CertificateSource): CertificateReorderer is assumed to expose
//
//	func NewCertificateReordererWithSigningCertificate(signingCertificate *model.CertificateToken, certificateChain []*model.CertificateToken) CertificateReorderer
//	func (r CertificateReorderer) OrderedCertificates() ([]*model.CertificateToken, error)
//
// and CertificateSource is assumed to expose IsTrusted(*model.CertificateToken) bool, matching
// the Java interface it flattens from. If X509-B's actual signatures differ, this file's
// embedding and calls need to be adjusted accordingly.
package spi

import (
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/utils"
)

// BaselineBCertificateSelector is used to retrieve the used certificates for a signature from
// the user parameters. It avoids duplicate entries, orders certificates from the signing
// certificate to the Root CA and filters trust anchors depending on the policy.
type BaselineBCertificateSelector struct {
	*CertificateReorderer

	// trustedCertificateSource is the trusted certificate source to be used on certificate
	// chain building.
	trustedCertificateSource CertificateSource

	// trustAnchorBPPolicy indicates whether a trust anchor policy should be used. When
	// enabled, the trust anchor is not included in the generated certificate chain.
	// Otherwise, the chain is generated up to a trust anchor, including the trust anchor
	// itself. Default: true.
	trustAnchorBPPolicy bool
}

// NewBaselineBCertificateSelector builds a certificate chain selector for signingCertificate.
// Port of the BaselineBCertificateSelector(CertificateToken, Collection<CertificateToken>)
// constructor.
func NewBaselineBCertificateSelector(signingCertificate *model.CertificateToken, certificateChain []*model.CertificateToken) *BaselineBCertificateSelector {
	return &BaselineBCertificateSelector{
		CertificateReorderer: NewCertificateReordererWithSigningCertificate(signingCertificate, certificateChain),
		trustAnchorBPPolicy:  true,
	}
}

// SetTrustedCertificateSource sets the trusted certificate source. Port of
// setTrustedCertificateSource(CertificateSource); returns the selector, mirroring Java's
// fluent setter.
func (s *BaselineBCertificateSelector) SetTrustedCertificateSource(trustedCertificateSource CertificateSource) *BaselineBCertificateSelector {
	s.trustedCertificateSource = trustedCertificateSource
	return s
}

// SetTrustAnchorBPPolicy sets whether a trust anchor policy should be used. When enabled, the
// trust anchor is not included in the generated certificate chain. Otherwise, the chain is
// generated up to a trust anchor, including the trust anchor itself. Default: true.
// Port of setTrustAnchorBPPolicy(boolean); returns the selector, mirroring Java's fluent setter.
func (s *BaselineBCertificateSelector) SetTrustAnchorBPPolicy(trustAnchorBPPolicy bool) *BaselineBCertificateSelector {
	s.trustAnchorBPPolicy = trustAnchorBPPolicy
	return s
}

// Certificates returns a certificate chain for a B-level signature creation.
// Port of getCertificates().
//
// OrderedCertificates can fail (e.g. no signing certificate can be identified); Java lets the
// DSSException propagate unchecked from getOrderedCertificates(), so the error is returned
// here instead of panicking, since it is data-dependent on the input certificate collection.
func (s *BaselineBCertificateSelector) Certificates() ([]*model.CertificateToken, error) {
	orderedCertificates, err := s.OrderedCertificates()
	if err != nil {
		return nil, err
	}

	// if true, trust anchor certificates (and upper certificates) are not included in the signature
	if s.trustAnchorBPPolicy && s.trustedCertificateSource != nil && utils.IsCollectionNotEmpty(orderedCertificates) {
		result := make([]*model.CertificateToken, 0, len(orderedCertificates))
		// signing certificate is required to be included in BASELINE-B signature
		signingCertificate := orderedCertificates[0]
		result = append(result, signingCertificate)
		for _, certificateToken := range orderedCertificates[1:] {
			if s.trustedCertificateSource.IsTrusted(certificateToken) {
				break
			}
			result = append(result, certificateToken)
		}
		return result, nil
	}
	return orderedCertificates, nil
}
