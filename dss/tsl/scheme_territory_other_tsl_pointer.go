// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/function/SchemeTerritoryOtherTSLPointer.java (DSS 6.5.RC1).
package tsl

import "github.com/ryftcore/dss-go/dss/trustedlist/jaxb"

// schemeTerritoryOtherTSLPointerExpectedTagName is the private static EXPECTED_TAG_NAME.
const schemeTerritoryOtherTSLPointerExpectedTagName = "{http://uri.etsi.org/02231/v2#}SchemeTerritory"

// SchemeTerritoryOtherTSLPointer filters certain TLs by the accepted country codes.
type SchemeTerritoryOtherTSLPointer struct {
	// countryCodes is the collection of country codes to be accepted.
	countryCodes map[string]struct{}
}

var _ OtherTSLPointerPredicate = (*SchemeTerritoryOtherTSLPointer)(nil)

// NewSchemeTerritoryOtherTSLPointer is the constructor allowing to filter a single country code.
// Port of SchemeTerritoryOtherTSLPointer(String).
func NewSchemeTerritoryOtherTSLPointer(countryCode string) *SchemeTerritoryOtherTSLPointer {
	return NewSchemeTerritoryOtherTSLPointerCollection([]string{countryCode})
}

// NewSchemeTerritoryOtherTSLPointerCollection is the constructor allowing to filter a collection
// of country codes. Port of SchemeTerritoryOtherTSLPointer(Collection).
func NewSchemeTerritoryOtherTSLPointerCollection(countryCodes []string) *SchemeTerritoryOtherTSLPointer {
	set := make(map[string]struct{}, len(countryCodes))
	for _, c := range countryCodes {
		set[c] = struct{}{}
	}
	return &SchemeTerritoryOtherTSLPointer{countryCodes: set}
}

// Test ports test(OtherTSLPointerType).
func (p *SchemeTerritoryOtherTSLPointer) Test(o *jaxb.OtherTSLPointerType) bool {
	extracted := extractAdditionalInformation(o)
	schemeTerritory, _ := extracted[schemeTerritoryOtherTSLPointerExpectedTagName].(string)
	_, ok := p.countryCodes[schemeTerritory]
	return ok
}
