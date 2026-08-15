// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/trust/consistency/TrustServiceStatusPreEIDASConsistency.java (DSS 6.5.RC1).
package qualification

import "github.com/utain/esig/dss/diagnostic"

// trustServiceStatusPreEIDASConsistency verifies status of a trusted service created
// before eIDAS.
type trustServiceStatusPreEIDASConsistency struct{}

// newTrustServiceStatusPreEIDASConsistency is the default constructor.
func newTrustServiceStatusPreEIDASConsistency() *trustServiceStatusPreEIDASConsistency {
	return &trustServiceStatusPreEIDASConsistency{}
}

// IsConsistent is the port of the overridden isConsistent(TrustServiceWrapper).
func (c *trustServiceStatusPreEIDASConsistency) IsConsistent(trustService *diagnostic.TrustServiceWrapper) bool {
	startDate := trustService.StartDate
	if IsPreEIDAS(startDate) {
		status := trustService.Status
		return TrustServiceStatus_GRANTED.URI() != status && TrustServiceStatus_WITHDRAWN.URI() != status
	}
	return true
}
