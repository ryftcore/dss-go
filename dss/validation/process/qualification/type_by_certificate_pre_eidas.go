// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/certificate/checks/type/TypeByCertificatePreEIDAS.java (DSS 6.5.RC1).
package qualification

import (
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// typeByCertificatePreEIDAS gets certificate approval status type before
// eIDAS.
type typeByCertificatePreEIDAS struct {
	// signingCertificate is the certificate to get usage type for.
	signingCertificate *diagnostic.CertificateWrapper
}

// newTypeByCertificatePreEIDAS is the default constructor. Port of
// TypeByCertificatePreEIDAS(CertificateWrapper).
func newTypeByCertificatePreEIDAS(signingCertificate *diagnostic.CertificateWrapper) *typeByCertificatePreEIDAS {
	return &typeByCertificatePreEIDAS{signingCertificate: signingCertificate}
}

// Type is the port of the overridden getType().
func (t *typeByCertificatePreEIDAS) Type() enumerations.CertificateType {
	if t.signingCertificate.IsQcCompliance() || process.IsQCP(t.signingCertificate) || process.IsQCPPlus(t.signingCertificate) {
		return enumerations.CertificateType_ESIGN // if QC -> ESign
	}
	return enumerations.CertificateType_UNKNOWN
}
