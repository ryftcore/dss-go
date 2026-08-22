// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/certificate/checks/qscd/QSCDByCertificatePreEIDAS.java (DSS 6.5.RC1).
package qualification

import (
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/validation/process"
)

// qscdByCertificatePreEIDAS returns QSCD status for a certificate before
// eIDAS.
type qscdByCertificatePreEIDAS struct {
	// certificate is the certificate to get QSCD status for.
	certificate *diagnostic.CertificateWrapper
}

// newQSCDByCertificatePreEIDAS is the default constructor. Port of
// QSCDByCertificatePreEIDAS(CertificateWrapper).
func newQSCDByCertificatePreEIDAS(certificate *diagnostic.CertificateWrapper) *qscdByCertificatePreEIDAS {
	return &qscdByCertificatePreEIDAS{certificate: certificate}
}

// QSCDStatus reports QSCD from the pre-eIDAS QCP+ policy OID or the
// certificate's QC-statement SSCD flag. Port of the overridden
// getQSCDStatus().
func (q *qscdByCertificatePreEIDAS) QSCDStatus() enumerations.QSCDStatus {
	// checks in policy id extension
	policyIdSupportedByQSCD := process.IsQCPPlus(q.certificate)

	// checks in QC statement extension
	qcStatementSupportedByQSCD := q.certificate.IsSupportedByQSCD()

	if policyIdSupportedByQSCD || qcStatementSupportedByQSCD {
		return enumerations.QSCDStatusQSCD
	}
	return enumerations.QSCDStatusNotQSCD
}
