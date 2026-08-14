// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/validation/identifier/SignatureAttributeIdentifier.java (DSS 6.5.RC1).
//
// Java's abstract Identifier resolves each instance's report-visible simple class name via
// getClass().getSimpleName() at runtime; Go has no such reflection, so the concrete class name
// is threaded through explicitly, per the phase 2a handoff convention ("Identifier subclasses
// pass their Java simple class name to identifier constructors (report-visible)"). Concrete
// subclasses of SignatureAttributeIdentifier (signature-format-specific unsigned-attribute
// identifiers, ported in later phases) call NewSignatureAttributeIdentifierBase with their own
// Java simple class name.
package identifier

import "github.com/utain/esig/dss/model"

// SignatureAttributeIdentifier identifies uniquely an unsigned attribute of a signature.
type SignatureAttributeIdentifier struct {
	model.IdentifierBase
}

// NewSignatureAttributeIdentifierBase computes an identifier from the binaries with the fixed
// "SA-" prefix. Port of the protected SignatureAttributeIdentifier(byte[]) constructor; every
// concrete subclass constructor calls this with its own Java simple class name.
func NewSignatureAttributeIdentifierBase(className string, data []byte) SignatureAttributeIdentifier {
	return SignatureAttributeIdentifier{
		IdentifierBase: model.NewIdentifierBase(className, "SA-", data),
	}
}

// compile-time assertion: a SignatureAttributeIdentifier is an identifier.
var _ model.Identifier = (*SignatureAttributeIdentifier)(nil)
