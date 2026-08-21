// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/function/NonEmptyServiceInformation.java (DSS 6.5.RC1).
package tsl

import "github.com/utain/esig/dss/trustedlist/jaxb"

// NonEmptyServiceInformation filters non-empty ServiceInformation element.
type NonEmptyServiceInformation struct{}

var _ TrustServicePredicate = (*NonEmptyServiceInformation)(nil)

// NewNonEmptyServiceInformation is the default constructor. Port of NonEmptyServiceInformation().
func NewNonEmptyServiceInformation() *NonEmptyServiceInformation {
	return &NonEmptyServiceInformation{}
}

// Test ports test(TSPServiceType).
func (p *NonEmptyServiceInformation) Test(tspServiceType *jaxb.TSPServiceType) bool {
	return tspServiceType.ServiceInformation != nil && tspServiceType.ServiceInformation.StatusStartingTime != nil
}
