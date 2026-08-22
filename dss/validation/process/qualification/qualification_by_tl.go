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
	qualifiedInCert Strategy
}

// newQualificationByTL is the default constructor. Port of
// QualificationByTL(TrustServiceWrapper, Strategy).
func newQualificationByTL(trustService *diagnostic.TrustServiceWrapper, qualifiedInCert Strategy) *qualificationByTL {
	return &qualificationByTL{trustService: trustService, qualifiedInCert: qualifiedInCert}
}

// QualifiedStatus reports qualification derived from the granted Trusted
// Service's captured qualifiers, deferring to the certificate-derived
// strategy where the service does not override it. Port of the overridden
// getQualifiedStatus().
func (q *qualificationByTL) QualifiedStatus() enumerations.CertificateQualifiedStatus {
	if q.trustService == nil {
		return enumerations.CertificateQualifiedStatusNotQC
	}

	grantedFilter := NewGrantedServiceFilter()
	if !grantedFilter.IsAcceptable(q.trustService) {
		return enumerations.CertificateQualifiedStatusNotQC
	}

	capturedQualifiers := q.trustService.CapturedQualifierUris()

	// If overrules
	if utils.IsCollectionNotEmpty(capturedQualifiers) {
		if enumerations.ServiceQualificationIsNotQualified(capturedQualifiers) {
			return enumerations.CertificateQualifiedStatusNotQC
		}
		if enumerations.ServiceQualificationIsQcStatement(capturedQualifiers) {
			return enumerations.CertificateQualifiedStatusQC
		}
	}

	return q.qualifiedInCert.QualifiedStatus()
}
