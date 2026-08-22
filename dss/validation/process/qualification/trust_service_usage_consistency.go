// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/trust/consistency/TrustServiceUsageConsistency.java (DSS 6.5.RC1).
package qualification

import (
	"github.com/ryftcore/dss-go/dss/diagnostic"
	"github.com/ryftcore/dss-go/dss/enumerations"
)

// trustServiceUsageConsistency: a Trusted Service can only have one of these values
// {QcForEsig, QcForEseal or QcForWSA} or none.
type trustServiceUsageConsistency struct{}

// newTrustServiceUsageConsistency is the default constructor.
func newTrustServiceUsageConsistency() *trustServiceUsageConsistency {
	return &trustServiceUsageConsistency{}
}

// IsConsistent reports whether the trust service carries at most one of
// QcForEsig, QcForEseal, or QcForWSA. Port of the overridden
// isConsistent(TrustServiceWrapper).
func (c *trustServiceUsageConsistency) IsConsistent(trustService *diagnostic.TrustServiceWrapper) bool {
	capturedQualifiers := trustService.CapturedQualifierUris()

	qcForEsig := enumerations.ServiceQualificationIsQcForEsig(capturedQualifiers)
	qcForEseal := enumerations.ServiceQualificationIsQcForEseal(capturedQualifiers)
	qcForWSA := enumerations.ServiceQualificationIsQcForWSA(capturedQualifiers)

	noneOfThem := !(qcForEsig || qcForEseal || qcForWSA)

	count := 0
	for _, b := range []bool{qcForEsig, qcForEseal, qcForWSA} {
		if b {
			count++
		}
	}
	onlyOneOfThem := count == 1

	return noneOfThem || onlyOneOfThem
}
