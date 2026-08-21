// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/certificate/checks/qualified/QualificationByTL.java (DSS 6.5.RC1).
package qualification

import (
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/utils"
)

// qualificationByTL gets certificate qualification status based on
// information extracted from a TrustService.
type qualificationByTL struct {
	// trustService is the Trusted Service to get qualification status from.
	trustService *diagnostic.TrustServiceWrapper

	// qualifiedInCert is the qualification strategy to be used.
	qualifiedInCert QualificationStrategy
}

// newQualificationByTL is the default constructor. Port of
// QualificationByTL(TrustServiceWrapper, QualificationStrategy).
func newQualificationByTL(trustService *diagnostic.TrustServiceWrapper, qualifiedInCert QualificationStrategy) *qualificationByTL {
	return &qualificationByTL{trustService: trustService, qualifiedInCert: qualifiedInCert}
}

// QualifiedStatus is the port of the overridden getQualifiedStatus().
func (q *qualificationByTL) QualifiedStatus() enumerations.CertificateQualifiedStatus {
	if q.trustService == nil {
		return enumerations.CertificateQualifiedStatus_NOT_QC
	}

	grantedFilter := NewGrantedServiceFilter()
	if !grantedFilter.IsAcceptable(q.trustService) {
		return enumerations.CertificateQualifiedStatus_NOT_QC
	}

	capturedQualifiers := q.trustService.CapturedQualifierUris()

	// If overrules
	if utils.IsCollectionNotEmpty(capturedQualifiers) {
		if enumerations.ServiceQualificationIsNotQualified(capturedQualifiers) {
			return enumerations.CertificateQualifiedStatus_NOT_QC
		}
		if enumerations.ServiceQualificationIsQcStatement(capturedQualifiers) {
			return enumerations.CertificateQualifiedStatus_QC
		}
	}

	return q.qualifiedInCert.QualifiedStatus()
}
