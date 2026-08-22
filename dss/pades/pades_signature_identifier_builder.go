// Ported from dss-pades/src/main/java/eu/europa/esig/dss/pades/validation/PAdESSignatureIdentifierBuilder.java
// (DSS 6.5.RC1).
package pades

import (
	"bytes"

	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi/validation"
)

// SignatureIdentifierBuilder builds a signature identifier for a PAdES signature.
type SignatureIdentifierBuilder struct {
	validation.AbstractSignatureIdentifierBuilder

	// padesSignature is the signature to build the identifier for, typed as the concrete PAdES
	// signature so SignaturePosition needs no runtime cast (Java casts its inherited `signature`
	// field instead, since its base class only knows AdvancedSignature).
	padesSignature *Signature
}

// NewPAdESSignatureIdentifierBuilder is the port of the constructor
// SignatureIdentifierBuilder(Signature).
func NewSignatureIdentifierBuilder(signature *Signature) *SignatureIdentifierBuilder {
	b := &SignatureIdentifierBuilder{
		AbstractSignatureIdentifierBuilder: validation.NewAbstractSignatureIdentifierBuilderBase(signature),
		padesSignature:                     signature,
	}
	b.InitAbstractSignatureIdentifierBuilder(b)
	return b
}

// CounterSignaturePosition is not supported in PAdES. Port of the protected
// getCounterSignaturePosition(AdvancedSignature) override.
//
// Panics with the Java message, matching UnsupportedOperationException.
func (b *SignatureIdentifierBuilder) CounterSignaturePosition(masterSignature validation.AdvancedSignature) any {
	panic("Not supported in PAdES!")
}

// SignaturePosition returns a position of a signature in the provided file. Port of the
// protected getSignaturePosition() override.
func (b *SignatureIdentifierBuilder) SignaturePosition() any {
	pdfRevision := b.padesSignature.PdfRevision()
	var buffer bytes.Buffer
	for _, signatureField := range pdfRevision.Fields() {
		buffer.WriteString(signatureField.FieldName())
	}
	return buffer.String()
}

// BuildBinaries re-implements spi/validation.AbstractSignatureIdentifierBuilder.BuildBinaries():
// Go has no virtual dispatch from an embedded base back into an overriding embedder, so
// PositionId's calls to CounterSignaturePosition/SignaturePosition have to route through this
// type's own overrides explicitly (see PORTING.md precedent in
// analyzer/default_document_analyzer.go's "Virtual dispatch" note, and
// cades/cades_signature_identifier_builder.go for the same pattern applied to a sibling
// identifier builder). PAdES does not override WriteSignedProperties, so this simply reproduces
// the base's BuildBinaries body over the base's WriteSignaturePosition (which itself dispatches
// through PositionId to this type's SignaturePosition/CounterSignaturePosition, correctly).
func (b *SignatureIdentifierBuilder) BuildBinaries() []byte {
	buffer := &bytes.Buffer{}
	b.WriteSignedProperties(buffer)
	b.WriteSignaturePosition(buffer)
	return buffer.Bytes()
}

// BuildSignatureIdentifier builds the SignatureIdentifier for the provided signature. Port of
// build() with its Java return type; shadows the base to route through this type's
// BuildBinaries().
func (b *SignatureIdentifierBuilder) BuildSignatureIdentifier() *validation.SignatureIdentifier {
	return validation.NewSignatureIdentifier(b.BuildBinaries())
}

// Build builds the SignatureIdentifier for the provided signature, satisfying
// model.IdentifierBuilder. Port of build(); shadows the base for the same reason as
// BuildSignatureIdentifier.
func (b *SignatureIdentifierBuilder) Build() model.Identifier {
	return b.BuildSignatureIdentifier()
}
