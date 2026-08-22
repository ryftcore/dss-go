// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/function/TrustServiceProviderByTSPName.java (DSS 6.5.RC1).
package tsl

import (
	"github.com/ryftcore/dss-go/dss/trustedlist/jaxb"
	"github.com/ryftcore/dss-go/dss/utils"
)

// TrustServiceProviderByTSPName filters TrustServicesProviders by TSP name.
type TrustServiceProviderByTSPName struct {
	// tspName is the name to filter by.
	tspName string
}

var _ TrustServiceProviderPredicate = (*TrustServiceProviderByTSPName)(nil)

// NewTrustServiceProviderByTSPName is the default constructor. Port of
// TrustServiceProviderByTSPName(String).
func NewTrustServiceProviderByTSPName(tspName string) *TrustServiceProviderByTSPName {
	return &TrustServiceProviderByTSPName{tspName: tspName}
}

// Test ports test(TSPType).
func (p *TrustServiceProviderByTSPName) Test(trustServiceProvider *jaxb.TSPType) bool {
	if trustServiceProvider != nil && utils.IsStringNotEmpty(p.tspName) {
		tspInformation := trustServiceProvider.TSPInformation
		if tspInformation == nil || tspInformation.TSPName == nil {
			return false
		}
		for _, name := range tspInformation.TSPName.Name {
			if utils.AreStringsEqualIgnoreCase(p.tspName, name.Value) {
				return true
			}
		}
	}
	return false
}
