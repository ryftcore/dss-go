// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/signature/identifier/SignatureIdentifier.java (DSS 6.5.RC1).
package validation

import "github.com/ryftcore/dss-go/dss/model"

// signatureIdentifierPrefix is the identifier prefix for a signature identifier, passed to the
// Java super(String, byte[]) constructor as the literal "S-".
const signatureIdentifierPrefix = "S-"

// signatureIdentifierClassName is SignatureIdentifier.class.getSimpleName(), which
// model.NewIdentifierBase records in the identifier and which the validation reports show.
const signatureIdentifierClassName = "SignatureIdentifier"

// SignatureIdentifier is the DSS Signature identifier.
//
// Java declares the class final; Go has no such modifier, so the type is simply never embedded.
type SignatureIdentifier struct {
	model.IdentifierBase
}

// NewSignatureIdentifier builds a SignatureIdentifier over the given binaries.
// Port of the protected SignatureIdentifier(byte[]) constructor - protected in Java because only
// AbstractSignatureIdentifierBuilder.build() is meant to call it; exported here since Go has no
// protected visibility, following the convention set by model.NewTimestampTokenIdentifier for
// the analogous TimestampIdentifierBuilder-only TimestampTokenIdentifier constructor.
func NewSignatureIdentifier(bytes []byte) *SignatureIdentifier {
	return &SignatureIdentifier{model.NewIdentifierBase(signatureIdentifierClassName, signatureIdentifierPrefix, bytes)}
}

// compile-time assertion: a SignatureIdentifier is an Identifier.
var _ model.Identifier = (*SignatureIdentifier)(nil)
