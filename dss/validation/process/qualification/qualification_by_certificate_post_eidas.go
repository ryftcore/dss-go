// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/certificate/checks/qualified/QualificationByCertificatePostEIDAS.java (DSS 6.5.RC1).
package qualification

import (
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/utils"
)

// qualificationByCertificatePostEIDAS gets certificate qualification status
// after eIDAS.
type qualificationByCertificatePostEIDAS struct {
	// signingCertificate is the certificate to get qualification status for.
	signingCertificate *diagnostic.CertificateWrapper
}

// newQualificationByCertificatePostEIDAS is the default constructor. Port of
// QualificationByCertificatePostEIDAS(CertificateWrapper).
func newQualificationByCertificatePostEIDAS(signingCertificate *diagnostic.CertificateWrapper) *qualificationByCertificatePostEIDAS {
	return &qualificationByCertificatePostEIDAS{signingCertificate: signingCertificate}
}

// QualifiedStatus is the port of the overridden getQualifiedStatus().
func (q *qualificationByCertificatePostEIDAS) QualifiedStatus() enumerations.CertificateQualifiedStatus {
	if q.signingCertificate.IsQcCompliance() &&
		utils.IsCollectionEmpty(q.signingCertificate.QcLegislationCountryCodes()) {
		return enumerations.CertificateQualifiedStatus_QC
	}
	return enumerations.CertificateQualifiedStatus_NOT_QC
}
