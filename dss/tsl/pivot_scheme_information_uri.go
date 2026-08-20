// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/function/PivotSchemeInformationURI.java (DSS 6.5.RC1).
package tsl

import (
	"strings"

	"github.com/utain/esig/dss/trustedlist/jaxb"
)

// pivotSchemeInformationURIPivotSuffix is the defined condition in (draft) ETSI TS 119 615.
const pivotSchemeInformationURIPivotSuffix = ".xml"

// PivotSchemeInformationURI is the Pivot scheme information URI filter predicate.
type PivotSchemeInformationURI struct{}

// NewPivotSchemeInformationURI is the default constructor. Port of PivotSchemeInformationURI().
func NewPivotSchemeInformationURI() *PivotSchemeInformationURI {
	return &PivotSchemeInformationURI{}
}

// Test ports test(NonEmptyMultiLangURIType).
func (p *PivotSchemeInformationURI) Test(t *jaxb.NonEmptyMultiLangURIType) bool {
	if t != nil && t.Value != "" {
		return strings.HasSuffix(t.Value, pivotSchemeInformationURIPivotSuffix)
	}
	return false
}
