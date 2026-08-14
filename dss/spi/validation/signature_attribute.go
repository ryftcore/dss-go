// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/validation/SignatureAttribute.java (DSS 6.5.RC1).
//
// FORWARD DEPENDENCY (flagged per S2B_BRIEF.md): SignatureAttributeIdentifier (Java
// spi.validation.identifier.SignatureAttributeIdentifier) is assigned to sibling chunk
// ANALYZER (s2b_ANALYZER.txt), which lands it at dss/spi/validation/identifier (pkg
// identifier), a 1:1 subpackage per S2B_BRIEF.md's package-layout section.
package validation

import "github.com/utain/esig/dss/spi/validation/identifier"

// SignatureAttribute defines a child of "signed-signature-properties" or
// "unsigned-signature-properties" element.
type SignatureAttribute interface {
	// Identifier gets the attribute identifier. Port of getIdentifier().
	Identifier() identifier.SignatureAttributeIdentifier
}
