// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/tsl/TrustPropertiesCertificateSource.java (DSS 6.5.RC1).
package tsl

import "github.com/ryftcore/dss-go/dss/model"

// TrustPropertiesCertificateSource provides an abstraction for a certificate source containing
// information about a validation status of Trusted Lists and corresponding trust properties.
type TrustPropertiesCertificateSource interface {
	TrustedCertificateSourceWithTime

	// Summary gets TL Validation job summary.
	Summary() *TLValidationJobSummary
	// SetSummary sets TL Validation job summary.
	SetSummary(summary *TLValidationJobSummary)
	// TrustServices returns TrustProperties for the given certificate, when applicable.
	TrustServices(token *model.CertificateToken) []*TrustProperties
	// SetTrustPropertiesByCertificates allows filling the CertificateSource.
	SetTrustPropertiesByCertificates(trustPropertiesByCerts map[*model.CertificateToken][]*TrustProperties)
	// SetTrustTimeByCertificates allows filling the CertificateSource with trusted time
	// periods.
	SetTrustTimeByCertificates(trustTimeByCertificate map[*model.CertificateToken][]*CertificateTrustTime)
}
