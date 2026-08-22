// Ported from dss-cades/src/main/java/eu/europa/esig/dss/cades/validation/CAdESSignatureIdentifierBuilder.java (DSS 6.5.RC1).
package cades

import (
	"bytes"

	"github.com/ryftcore/dss-go/dss/internal/cmscore"
	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/spi/validation"
)

// CAdESSignatureIdentifierBuilder builds a SignatureIdentifier for a CAdES signature. Port of
// the class SignatureIdentifierBuilder, extending
// spi/validation.AbstractSignatureIdentifierBuilder.
type SignatureIdentifierBuilder struct {
	validation.AbstractSignatureIdentifierBuilder

	// cadesSignature is the signature to build the identifier for, typed as the concrete CAdES
	// signature so CounterSignaturePosition/SignaturePosition need no runtime cast (Java casts
	// its inherited `signature` field instead, since its base class only knows AdvancedSignature).
	cadesSignature *Signature
}

// NewCAdESSignatureIdentifierBuilder is the port of the constructor
// SignatureIdentifierBuilder(Signature).
func NewSignatureIdentifierBuilder(signature *Signature) *SignatureIdentifierBuilder {
	b := &SignatureIdentifierBuilder{
		AbstractSignatureIdentifierBuilder: validation.NewAbstractSignatureIdentifierBuilderBase(signature),
		cadesSignature:                     signature,
	}
	b.InitAbstractSignatureIdentifierBuilder(b)
	return b
}

// WriteSignedProperties shadows spi/validation.AbstractSignatureIdentifierBuilder's method of
// the same name: it appends the manifest filename after the inherited signed properties. Port
// of the protected writeSignedProperties(ByteArrayOutputStream) override.
func (b *SignatureIdentifierBuilder) WriteSignedProperties(buffer *bytes.Buffer) {
	b.AbstractSignatureIdentifierBuilder.WriteSignedProperties(buffer)
	b.WriteString(buffer, b.manifestFilename())
}

// manifestFilename ports the private getManifestFilename(); returns "" where Java returns null,
// which WriteString already treats identically to an absent value.
func (b *SignatureIdentifierBuilder) manifestFilename() string {
	manifestFile := b.Signature().ManifestFile()
	if manifestFile != nil {
		return manifestFile.Filename()
	}
	return ""
}

// BuildBinaries re-implements spi/validation.AbstractSignatureIdentifierBuilder.BuildBinaries():
// Go has no virtual dispatch from an embedded base back into an overriding embedder, so the
// call to WriteSignedProperties has to route through this type's own override explicitly (see
// PORTING.md precedent in analyzer/default_document_analyzer.go's "Virtual dispatch" note, and
// spi/validation/timestamp/timestamp_identifier_builder.go for the same pattern applied to a
// sibling identifier builder).
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

// CounterSignaturePosition returns the current counter signature position in its master
// signature. Port of the protected getCounterSignaturePosition(AdvancedSignature) override.
func (b *SignatureIdentifierBuilder) CounterSignaturePosition(masterSignature validation.AdvancedSignature) any {
	cadesMasterSignature := masterSignature.(*Signature)
	return cadesSignatureIdentifierBuilderCount(
		cadesMasterSignature.CounterSignatureStore().SignerInfos(), b.cadesSignature.SignerInformation())
}

// SignaturePosition returns a position of a signature in the provided file. Port of the
// protected getSignaturePosition() override.
func (b *SignatureIdentifierBuilder) SignaturePosition() any {
	return cadesSignatureIdentifierBuilderCount(
		b.cadesSignature.CMS().SignerInfos(), b.cadesSignature.SignerInformation())
}

// cadesSignatureIdentifierBuilderCount ports the private count(Collection<SignerInformation>,
// SignerInformation): the number of signers preceding currentSignerInformation, compared by
// identity ("compare by memory to avoid matching signers with identical content") - reproduced
// here as Go pointer identity, exactly analogous.
func cadesSignatureIdentifierBuilderCount(signerInformationStore []*cmscore.SignerInfo, currentSignerInformation *cmscore.SignerInfo) int {
	counter := 0
	for _, signerInformation := range signerInformationStore {
		if currentSignerInformation == signerInformation {
			break
		}
		counter++
	}
	return counter
}
