// Ported from dss-tsl-validation/src/main/java/eu/europa/esig/dss/tsl/function/MimetypeOtherTSLPointer.java (DSS 6.5.RC1).
package tsl

import (
	"strings"

	"github.com/ryftcore/dss-go/dss/trustedlist/jaxb"
)

// mimetypeOtherTSLPointerExpectedTagName is the private static EXPECTED_TAG_NAME.
const mimetypeOtherTSLPointerExpectedTagName = "{http://uri.etsi.org/02231/v2/additionaltypes#}MimeType"

// MimetypeOtherTSLPointer allows filtering of TSL pointers by a MimeType.
type MimetypeOtherTSLPointer struct {
	// expectedMimeType is the MimeType to filter by.
	expectedMimeType string
}

var _ OtherTSLPointerPredicate = (*MimetypeOtherTSLPointer)(nil)

// NewMimetypeOtherTSLPointer is the default constructor. Port of MimetypeOtherTSLPointer(String).
//
// Panics with the Java message when expectedMimeType is empty (Objects.requireNonNull).
func NewMimetypeOtherTSLPointer(expectedMimeType string) *MimetypeOtherTSLPointer {
	if expectedMimeType == "" {
		panic("Expected MimeType must be defined")
	}
	return &MimetypeOtherTSLPointer{expectedMimeType: expectedMimeType}
}

// Test ports test(OtherTSLPointerType).
func (p *MimetypeOtherTSLPointer) Test(o *jaxb.OtherTSLPointerType) bool {
	extracted := extractAdditionalInformation(o)
	mimeType, _ := extracted[mimetypeOtherTSLPointerExpectedTagName].(string)
	return strings.EqualFold(p.expectedMimeType, mimeType)
}
