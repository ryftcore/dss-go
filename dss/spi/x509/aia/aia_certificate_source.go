// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/x509/aia/AIACertificateSource.java (DSS 6.5.RC1).
package aia

import (
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi"
	"github.com/ryftcore/dss-go/dss/utils"
)

// AIACertificateSource is the certificate source requesting issuer certificates by AIA.
type AIACertificateSource struct {
	spi.CommonCertificateSource

	// certificate is the certificate token to get the issuer for.
	certificate *model.CertificateToken
}

// newAIACertificateSource is the constructor creating an AIA certificate source for
// certificate. Ports the protected AIACertificateSource(CertificateToken) constructor. Panics
// if certificate is nil (Java Objects.requireNonNull).
func newAIACertificateSource(certificate *model.CertificateToken) *AIACertificateSource {
	if certificate == nil {
		panic("The certificate cannot be null")
	}
	return &AIACertificateSource{
		CommonCertificateSource: spi.NewCommonCertificateSource(),
		certificate:             certificate,
	}
}

// AIACertificateSourceForCertificateToken retrieves an AIA.caIssuers for the given certificate
// using aiaSource. NOTE: this function performs an AIA URI request on instantiation.
func AIACertificateSourceForCertificateToken(certificate *model.CertificateToken, aiaSource AIASource) *AIACertificateSource {
	aiaCertificateSource := newAIACertificateSource(certificate)

	func() {
		// Java catches any Exception raised while retrieving/walking the AIA chain and
		// downgrades it to a dropped LOG.warn instead of propagating it; mirrored here with a
		// deferred recover, since aiaSource.CertificatesByAIA may panic on failure.
		defer func() {
			recover()
		}()

		extractedCertificates := aiaCertificateSourceDedup(aiaSource.CertificatesByAIA(certificate))
		if utils.IsCollectionNotEmpty(extractedCertificates) {
			currentCertificate := certificate
			for currentCertificate != nil {
				issuer := aiaCertificateSourceGetIssuer(currentCertificate, extractedCertificates)
				if aiaCertificateSourceContains(aiaCertificateSource.Certificates(), issuer) {
					// break for processed certificates
					break
				} else if issuer != nil {
					// add issuer for processing
					aiaCertificateSource.AddCertificate(issuer)
				}
				currentCertificate = issuer
			}

			// if no certificates have been extracted -> add all
			if utils.IsCollectionEmpty(aiaCertificateSource.Certificates()) {
				for _, certificateToken := range extractedCertificates {
					aiaCertificateSource.AddCertificate(certificateToken)
				}
			}
		}
		// else: "No AIA certificates have been retrieved..." LOG.warn dropped, not
		// load-bearing per PORTING.md.
	}()

	return aiaCertificateSource
}

// aiaCertificateSourceGetIssuer gets the issuer certificate for certificate from the given
// collection of candidates. Returns nil if no suitable issuer was found. Ports the protected
// static getIssuer(CertificateToken, Collection<CertificateToken>).
func aiaCertificateSourceGetIssuer(certificate *model.CertificateToken, candidates []*model.CertificateToken) *model.CertificateToken {
	issuer := spi.NewTokenIssuerSelector(certificate, candidates).Issuer()
	if issuer != nil && certificate.IsSignedByToken(issuer) {
		return issuer
	}
	return nil
}

// aiaCertificateSourceContains ports List#contains(CertificateToken), which relies on
// CertificateToken#equals (identity by DSSIDAsString()).
func aiaCertificateSourceContains(certificates []*model.CertificateToken, candidate *model.CertificateToken) bool {
	if candidate == nil {
		return false
	}
	for _, certificate := range certificates {
		if certificate.DSSIDAsString() == candidate.DSSIDAsString() {
			return true
		}
	}
	return false
}

// aiaCertificateSourceDedup de-duplicates certificateTokens by DSSIDAsString() while preserving
// order, standing in for `new LinkedHashSet<>(...)`.
func aiaCertificateSourceDedup(certificateTokens []*model.CertificateToken) []*model.CertificateToken {
	seen := make(map[string]struct{}, len(certificateTokens))
	result := make([]*model.CertificateToken, 0, len(certificateTokens))
	for _, certificateToken := range certificateTokens {
		id := certificateToken.DSSIDAsString()
		if _, found := seen[id]; found {
			continue
		}
		seen[id] = struct{}{}
		result = append(result, certificateToken)
	}
	return result
}

// IssuerFromAIA gets the issuer's certificate from Authority Information Access through the
// id-ad-caIssuers extension. Returns nil if not found.
func (s *AIACertificateSource) IssuerFromAIA() *model.CertificateToken {
	candidates := s.Certificates()
	if utils.IsCollectionNotEmpty(candidates) {
		// The potential issuers might support 3 known scenarios:
		// - issuer certificate with single entry
		// - issuer certificate is a collection of bridge certificates (all having the
		// same public key)
		// - full certification path (up to the root of the chain)
		// In case the issuer is a collection of bridge certificates, only one of the
		// bridge certificates needs to be verified
		issuer := aiaCertificateSourceGetIssuer(s.certificate, candidates)
		// "The retrieved certificate(s) using AIA do not sign the certificate..." LOG.warn
		// dropped when issuer == nil, not load-bearing per PORTING.md.
		return issuer
	}
	return nil
}

// CertificateSourceType returns CertificateSourceTypeAIA.
func (s *AIACertificateSource) CertificateSourceType() enumerations.CertificateSourceType {
	return enumerations.CertificateSourceTypeAIA
}
