// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/certificate/checks/qscd/QSCDByTL.java (DSS 6.5.RC1).
package qualification

import (
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/utils"
)

// qscdByTL extracts QCSD status from a Trusted Service.
type qscdByTL struct {
	// trustService is the Trusted Service to extract QSCD status from.
	trustService *diagnostic.TrustServiceWrapper

	// qualified is the qualification status of the certificate.
	qualified enumerations.CertificateQualifiedStatus

	// qscdFromCertificate is the QSCD strategy to be used.
	qscdFromCertificate QSCDStrategy
}

// newQSCDByTL is the default constructor. Port of
// QSCDByTL(TrustServiceWrapper, CertificateQualifiedStatus, QSCDStrategy).
func newQSCDByTL(trustService *diagnostic.TrustServiceWrapper, qualified enumerations.CertificateQualifiedStatus,
	qscdFromCertificate QSCDStrategy) *qscdByTL {
	return &qscdByTL{trustService: trustService, qualified: qualified, qscdFromCertificate: qscdFromCertificate}
}

// QSCDStatus extracts the QSCD status from the Trusted Service's captured
// qualifiers, falling back to the certificate-derived strategy when the
// service defers to the certificate. Port of the overridden
// getQSCDStatus().
func (q *qscdByTL) QSCDStatus() enumerations.QSCDStatus {
	if q.trustService == nil || !enumerations.CertificateQualifiedStatusIsQC(q.qualified) {
		return enumerations.QSCDStatusNotQSCD
	}

	capturedQualifiers := q.trustService.CapturedQualifierUris()

	// If overrules
	if utils.IsCollectionNotEmpty(capturedQualifiers) {
		if IsPostEIDAS(q.trustService.StartDate) {
			if enumerations.ServiceQualificationIsQcWithQSCD(capturedQualifiers) || enumerations.ServiceQualificationIsQcQSCDManagedOnBehalf(capturedQualifiers) {
				return enumerations.QSCDStatusQSCD
			} else if enumerations.ServiceQualificationIsQcQSCDStatusAsInCert(capturedQualifiers) {
				return q.qscdFromCertificate.QSCDStatus()
			} else if enumerations.ServiceQualificationIsQcNoQSCD(capturedQualifiers) {
				return enumerations.QSCDStatusNotQSCD
			}
		} else { // pre eIDAS
			if enumerations.ServiceQualificationIsQcWithSSCD(capturedQualifiers) {
				return enumerations.QSCDStatusQSCD
			} else if enumerations.ServiceQualificationIsQcSSCDStatusAsInCert(capturedQualifiers) {
				return q.qscdFromCertificate.QSCDStatus()
			} else if enumerations.ServiceQualificationIsQcNoSSCD(capturedQualifiers) {
				return enumerations.QSCDStatusNotQSCD
			}
		}
	}

	return q.qscdFromCertificate.QSCDStatus()
}
