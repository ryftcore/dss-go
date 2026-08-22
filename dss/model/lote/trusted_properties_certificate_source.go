// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/lote/TrustedPropertiesCertificateSource.java (DSS 6.5.RC1).
package lote

import (
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/model/tsl"
)

// TrustedPropertiesCertificateSource contains trusted certificates and related trusted
// properties.
type TrustedPropertiesCertificateSource interface {
	tsl.TrustedCertificateSourceWithTime

	// Summary gets TL Validation job summary.
	Summary() *ValidationJobSummary
	// SetSummary sets TL Validation job summary.
	SetSummary(summary *ValidationJobSummary)
	// TrustedProperties returns TrustedProperties for the given certificate, when applicable.
	TrustedProperties(token *model.CertificateToken) []*TrustedProperties
	// SetTrustedPropertiesByCertificates allows filling the CertificateSource.
	SetTrustedPropertiesByCertificates(trustedPropertiesByCerts map[*model.CertificateToken][]*TrustedProperties)
	// SetTrustedTimeByCertificates allows filling the CertificateSource with trusted time
	// periods.
	SetTrustedTimeByCertificates(trustedTimeByCertificate map[*model.CertificateToken][]*tsl.CertificateTrustTime)
}
