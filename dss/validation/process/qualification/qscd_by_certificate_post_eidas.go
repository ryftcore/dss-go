// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/certificate/checks/qscd/QSCDByCertificatePostEIDAS.java (DSS 6.5.RC1).
package qualification

import (
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
)

// qscdByCertificatePostEIDAS returns QSCD status for a certificate after
// eIDAS.
type qscdByCertificatePostEIDAS struct {
	// certificate is the certificate to get QSCD status for.
	certificate *diagnostic.CertificateWrapper
}

// newQSCDByCertificatePostEIDAS is the default constructor. Port of
// QSCDByCertificatePostEIDAS(CertificateWrapper).
func newQSCDByCertificatePostEIDAS(certificate *diagnostic.CertificateWrapper) *qscdByCertificatePostEIDAS {
	return &qscdByCertificatePostEIDAS{certificate: certificate}
}

// QSCDStatus reports QSCD from the certificate's QC-statement QSCD flag
// only. Port of the overridden getQSCDStatus().
func (q *qscdByCertificatePostEIDAS) QSCDStatus() enumerations.QSCDStatus {
	// checks only in QC statement extension
	if q.certificate.IsSupportedByQSCD() {
		return enumerations.QSCDStatusQSCD
	}
	return enumerations.QSCDStatusNotQSCD
}
