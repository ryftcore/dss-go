// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/process/qualification/trust/consistency/TrustServiceCondition.java (DSS 6.5.RC1).
package qualification

import "github.com/ryftcore/dss-go/dss/diagnostic"

// TrustServiceCondition checks whether the TrustService is valid.
type TrustServiceCondition interface {
	// IsConsistent reports whether the trustService is consistent. Port of
	// isConsistent(TrustServiceWrapper).
	IsConsistent(trustService *diagnostic.TrustServiceWrapper) bool
}
