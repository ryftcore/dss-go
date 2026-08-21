// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/function/TypeOtherTSLPointer.java (DSS 6.5.RC1).
package tsl

import (
	"strings"

	"github.com/ryftcore/dss-go/dss/trustedlist/jaxb"
)

// typeOtherTSLPointerExpectedTagName is the private static EXPECTED_TAG_NAME.
const typeOtherTSLPointerExpectedTagName = "{http://uri.etsi.org/02231/v2#}TSLType"

// TypeOtherTSLPointer allows TSL filtering by TSLType.
type TypeOtherTSLPointer struct {
	// expectedTSLType is the TSLType value to filter by.
	expectedTSLType string
}

var _ OtherTSLPointerPredicate = (*TypeOtherTSLPointer)(nil)

// NewTypeOtherTSLPointer is the default constructor. Port of TypeOtherTSLPointer(String).
//
// Panics with the Java message when expectedTSLType is empty (Objects.requireNonNull).
func NewTypeOtherTSLPointer(expectedTSLType string) *TypeOtherTSLPointer {
	if expectedTSLType == "" {
		panic("Expected TSLType must be defined")
	}
	return &TypeOtherTSLPointer{expectedTSLType: expectedTSLType}
}

// Test ports test(OtherTSLPointerType).
func (p *TypeOtherTSLPointer) Test(o *jaxb.OtherTSLPointerType) bool {
	extracted := extractAdditionalInformation(o)
	tslType, _ := extracted[typeOtherTSLPointerExpectedTagName].(string)
	return strings.EqualFold(p.expectedTSLType, tslType)
}
