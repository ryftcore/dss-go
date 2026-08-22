// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/certificate/checks/type/TypeByTL.java (DSS 6.5.RC1).
package qualification

import (
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/utils"
)

// typeByTL gets certificate approval status type based on the information
// extracted from a TrustService.
type typeByTL struct {
	// trustService is the Trusted Service to get certificate approval status
	// type from.
	trustService *diagnostic.TrustServiceWrapper

	// qualified is the certificate qualified status.
	qualified enumerations.CertificateQualifiedStatus

	// typeInCert is the certificate's usage type extraction strategy.
	typeInCert TypeStrategy
}

// newTypeByTL is the default constructor. Port of
// TypeByTL(TrustServiceWrapper, CertificateQualifiedStatus, TypeStrategy).
func newTypeByTL(trustService *diagnostic.TrustServiceWrapper, qualified enumerations.CertificateQualifiedStatus,
	typeInCert TypeStrategy) *typeByTL {
	return &typeByTL{trustService: trustService, qualified: qualified, typeInCert: typeInCert}
}

// Type reports the certificate's usage type derived from the Trusted
// Service's captured qualifiers, falling back to the certificate-derived
// strategy where the service does not override it. Port of the overridden
// getType().
func (t *typeByTL) Type() enumerations.CertificateType {

	// overrules are only applicable when the certificate is qualified (cert + TL)
	if enumerations.CertificateQualifiedStatusIsQC(t.qualified) {

		if t.trustService == nil {
			return enumerations.CertificateTypeUnknown
		}

		if IsPreEIDAS(t.trustService.StartDate) {
			return enumerations.CertificateTypeESign
		}

		usageQualifiers := enumerations.ServiceQualificationGetUsageQualifiers(t.trustService.CapturedQualifierUris())

		if utils.CollectionSize(usageQualifiers) > 1 {
			return enumerations.CertificateTypeUnknown

		} else if utils.IsCollectionNotEmpty(usageQualifiers) {
			// If overrules

			if enumerations.ServiceQualificationIsQcForEsig(usageQualifiers) {
				return enumerations.CertificateTypeESign
			}

			if enumerations.ServiceQualificationIsQcForEseal(usageQualifiers) {
				return enumerations.CertificateTypeESeal
			}

			if enumerations.ServiceQualificationIsQcForWSA(usageQualifiers) {
				return enumerations.CertificateTypeWSA
			}
		}
	}

	return t.typeInCert.Type()
}
