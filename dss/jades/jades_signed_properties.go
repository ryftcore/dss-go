// Ported from dss-jades/src/main/java/eu/europa/esig/dss/jades/validation/JAdESSignedProperties.java
// (DSS 6.5.RC1).
package jades

import (
	"github.com/ryftcore/dss-go/dss/internal/jose"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi/validation"
)

// JAdESSignedProperties represents a list of JAdES signed properties (protected header). Port of
// the class SignedProperties, implementing validation.SignatureProperties[*Attribute].
type SignedProperties struct {
	// headers represent the protected header map. Port of the private final Headers headers
	// field.
	headers *jose.Headers
}

// NewJAdESSignedProperties is the default constructor. Port of the public
// SignedProperties(Headers) constructor.
func NewJAdESSignedProperties(headers *jose.Headers) *SignedProperties {
	return &SignedProperties{headers: headers}
}

// IsExist checks if "unsigned-signature-properties" exists and can be processed. Port of
// isExist().
func (p *SignedProperties) IsExist() bool {
	return p.headers != nil
}

// Attributes returns a list of children contained in the element. Port of getAttributes().
//
// Panics with the Java message wrapping the underlying error when the headers cannot be
// re-parsed as a map: getMapKeyValues() throws an unchecked DSSException, propagating out of
// getAttributes() uncaught (no AdvancedSignature accessor reached from here has an error
// return to use instead).
func (p *SignedProperties) Attributes() []*Attribute {
	headerMap := p.mapKeyValues()

	var attributes []*Attribute
	for _, key := range headerMap.Keys() {
		attributes = append(attributes, NewJAdESAttribute(key, headerMap.Value(key)))
	}
	return attributes
}

// mapKeyValues ports the private getMapKeyValues().
//
// TODO avoid to parse (upstream's own comment, reproduced verbatim).
func (p *SignedProperties) mapKeyValues() *jose.Object {
	headerMap, err := DSSJsonUtilsParseJSONStringToMap(p.headers.FullHeaderAsJSONString())
	if err != nil {
		panic(model.NewDSSErrorMessageCause("Unable to retrieve the map from the headers", err))
	}
	return headerMap
}

// compile-time assertion: a SignedProperties satisfies
// validation.SignatureProperties[*Attribute].
var _ validation.SignatureProperties[*Attribute] = (*SignedProperties)(nil)
