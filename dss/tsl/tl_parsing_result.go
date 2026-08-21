// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/parsing/TLParsingResult.java (DSS 6.5.RC1).
package tsl

import tslmodel "github.com/utain/esig/dss/model/tsl"

// TLParsingResult is a parsed TL result.
type TLParsingResult struct {
	AbstractTLParsingResult

	// trustServiceProviders is the list of found trust service providers.
	trustServiceProviders []*tslmodel.TrustServiceProvider
}

// NewTLParsingResult is the default constructor. Port of TLParsingResult().
func NewTLParsingResult() *TLParsingResult {
	return &TLParsingResult{AbstractTLParsingResult: NewAbstractTLParsingResult()}
}

// TrustServiceProviders gets the trust service providers. Port of getTrustServiceProviders().
func (r *TLParsingResult) TrustServiceProviders() []*tslmodel.TrustServiceProvider {
	return r.trustServiceProviders
}

// SetTrustServiceProviders sets the trust service providers. Port of
// setTrustServiceProviders(List).
func (r *TLParsingResult) SetTrustServiceProviders(trustServiceProviders []*tslmodel.TrustServiceProvider) {
	r.trustServiceProviders = trustServiceProviders
}
