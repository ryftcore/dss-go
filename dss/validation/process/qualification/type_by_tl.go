// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/certificate/checks/type/TypeByTL.java (DSS 6.5.RC1).
package qualification

import (
	"github.com/utain/esig/dss/diagnostic"
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/utils"
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

// Type is the port of the overridden getType().
func (t *typeByTL) Type() enumerations.CertificateType {

	// overrules are only applicable when the certificate is qualified (cert + TL)
	if enumerations.CertificateQualifiedStatusIsQC(t.qualified) {

		if t.trustService == nil {
			return enumerations.CertificateType_UNKNOWN
		}

		if IsPreEIDAS(t.trustService.StartDate) {
			return enumerations.CertificateType_ESIGN
		}

		usageQualifiers := enumerations.ServiceQualificationGetUsageQualifiers(t.trustService.CapturedQualifierUris())

		if utils.CollectionSize(usageQualifiers) > 1 {
			return enumerations.CertificateType_UNKNOWN

		} else if utils.IsCollectionNotEmpty(usageQualifiers) {
			// If overrules

			if enumerations.ServiceQualificationIsQcForEsig(usageQualifiers) {
				return enumerations.CertificateType_ESIGN
			}

			if enumerations.ServiceQualificationIsQcForEseal(usageQualifiers) {
				return enumerations.CertificateType_ESEAL
			}

			if enumerations.ServiceQualificationIsQcForWSA(usageQualifiers) {
				return enumerations.CertificateType_WSA
			}
		}
	}

	return t.typeInCert.Type()
}
