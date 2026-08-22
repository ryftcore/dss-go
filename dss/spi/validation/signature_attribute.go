// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/validation/SignatureAttribute.java (DSS 6.5.RC1).
//
// SignatureAttributeIdentifier (Java spi.validation.identifier.SignatureAttributeIdentifier)
// lands at dss/spi/validation/identifier (pkg identifier), a 1:1 subpackage.
package validation

import "github.com/ryftcore/dss-go/dss/spi/validation/identifier"

// SignatureAttribute defines a child of "signed-signature-properties" or
// "unsigned-signature-properties" element.
type SignatureAttribute interface {
	// Identifier gets the attribute identifier. Port of getIdentifier().
	Identifier() identifier.SignatureAttributeIdentifier
}
