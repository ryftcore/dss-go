// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/certificate/checks/qualified/QualificationByCertificatePreEIDAS.java (DSS 6.5.RC1).
package qualification

import (
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/utils"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// qualificationByCertificatePreEIDAS gets certificate qualification status
// before eIDAS.
type qualificationByCertificatePreEIDAS struct {
	// signingCertificate is the certificate to get qualification status for.
	signingCertificate *diagnostic.CertificateWrapper
}

// newQualificationByCertificatePreEIDAS is the default constructor. Port of
// QualificationByCertificatePreEIDAS(CertificateWrapper).
func newQualificationByCertificatePreEIDAS(signingCertificate *diagnostic.CertificateWrapper) *qualificationByCertificatePreEIDAS {
	return &qualificationByCertificatePreEIDAS{signingCertificate: signingCertificate}
}

// QualifiedStatus is the port of the overridden getQualifiedStatus().
func (q *qualificationByCertificatePreEIDAS) QualifiedStatus() enumerations.CertificateQualifiedStatus {
	if (q.signingCertificate.IsQcCompliance() || process.IsQCP(q.signingCertificate) ||
		process.IsQCPPlus(q.signingCertificate)) &&
		utils.IsCollectionEmpty(q.signingCertificate.QcLegislationCountryCodes()) {
		return enumerations.CertificateQualifiedStatusQC
	}
	return enumerations.CertificateQualifiedStatusNotQC
}
