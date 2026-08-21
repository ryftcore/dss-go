// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/function/NonEmptyTSPInformation.java (DSS 6.5.RC1).
package tsl

import "github.com/ryftcore/dss-go/dss/trustedlist/jaxb"

// NonEmptyTSPInformation filters non-empty TSPInformation element.
type NonEmptyTSPInformation struct{}

var _ TrustServiceProviderPredicate = (*NonEmptyTSPInformation)(nil)

// NewNonEmptyTSPInformation is the default constructor. Port of NonEmptyTSPInformation().
func NewNonEmptyTSPInformation() *NonEmptyTSPInformation {
	return &NonEmptyTSPInformation{}
}

// Test ports test(TSPType).
func (p *NonEmptyTSPInformation) Test(t *jaxb.TSPType) bool {
	return t.TSPInformation != nil
}
