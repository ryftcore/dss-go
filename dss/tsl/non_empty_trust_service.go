// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/function/NonEmptyTrustService.java (DSS 6.5.RC1).
package tsl

import (
	"github.com/utain/esig/dss/trustedlist/jaxb"
	"github.com/utain/esig/dss/utils"
)

// NonEmptyTrustService filters non-empty TrustServices.
type NonEmptyTrustService struct{}

var _ TrustServiceProviderPredicate = (*NonEmptyTrustService)(nil)

// NewNonEmptyTrustService is the default constructor. Port of NonEmptyTrustService().
func NewNonEmptyTrustService() *NonEmptyTrustService {
	return &NonEmptyTrustService{}
}

// Test ports test(TSPType).
func (p *NonEmptyTrustService) Test(t *jaxb.TSPType) bool {
	servicesList := t.TSPServices
	return servicesList != nil && utils.IsCollectionNotEmpty(servicesList.TSPService)
}
