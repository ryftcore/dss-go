// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/signature/identifier/SignatureIdentifierBuilder.java (DSS 6.5.RC1).
package validation

import "github.com/utain/esig/dss/model"

// SignatureIdentifierBuilder builds a deterministic Signature Identifier for the produced
// reports.
type SignatureIdentifierBuilder interface {
	// model.IdentifierBuilder supplies Build() model.Identifier, the generic contract Java's
	// interface inherits from IdentifierBuilder<SignatureIdentifier> (via the raw
	// IdentifierBuilder import repeated twice in the upstream source - a harmless duplicate
	// import with no Go counterpart).
	model.IdentifierBuilder

	// BuildSignatureIdentifier builds the SignatureIdentifier for the provided AdvancedSignature.
	// Port of build(), narrowed to its covariant Java return type; Go has no covariant returns,
	// so the generic Build() above and this method both exist, mirroring the
	// TimestampTokenIdentifierBuilder / TimestampIdentifierBuilder convention.
	BuildSignatureIdentifier() *SignatureIdentifier
}
