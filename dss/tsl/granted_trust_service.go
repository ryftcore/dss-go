// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/function/GrantedTrustService.java (DSS 6.5.RC1).
package tsl

import (
	"github.com/utain/esig/dss/trustedlist/jaxb"
	"github.com/utain/esig/dss/utils"
	"github.com/utain/esig/dss/validation/process/qualification"
)

// GrantedTrustService filters TrustServices by 'granted' property (supports pre- and
// post-eIDAS).
type GrantedTrustService struct{}

var _ TrustServicePredicate = (*GrantedTrustService)(nil)

// NewGrantedTrustService is the default constructor. Port of GrantedTrustService().
func NewGrantedTrustService() *GrantedTrustService {
	return &GrantedTrustService{}
}

// Test ports test(TSPServiceType).
func (p *GrantedTrustService) Test(trustService *jaxb.TSPServiceType) bool {
	if trustService == nil {
		return false
	}
	serviceInformation := trustService.ServiceInformation
	if serviceInformation == nil {
		return false
	}

	// Current status
	if qualification.TrustServiceStatusIsAcceptableStatusAfterEIDAS(serviceInformation.ServiceStatus) ||
		qualification.TrustServiceStatusIsAcceptableStatusBeforeEIDAS(serviceInformation.ServiceStatus) {
		return true
	}

	// Past
	serviceHistory := trustService.ServiceHistory
	if serviceHistory != nil && utils.IsCollectionNotEmpty(serviceHistory.ServiceHistoryInstance) {
		for _, serviceHistoryInstance := range serviceHistory.ServiceHistoryInstance {
			if qualification.TrustServiceStatusIsAcceptableStatusAfterEIDAS(serviceHistoryInstance.ServiceStatus) ||
				qualification.TrustServiceStatusIsAcceptableStatusBeforeEIDAS(serviceHistoryInstance.ServiceStatus) {
				return true
			}
		}
	}

	return false
}
