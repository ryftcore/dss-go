// Ported from
// dss-jades/src/main/java/eu/europa/esig/dss/jades/validation/JAdESSignatureIdentifierBuilder.java
// (DSS 6.5.RC1).
//
// DEVIATION: CounterSignaturePosition compares the two 'cSig' components' full identifiers
// (Identifier().Equals) rather than replaying Java's Object.hashCode() comparison
// (masterCSigAttribute.hashCode() == jadesCounterSignature.getMasterCSigComponent().hashCode()) -
// Attribute.hashCode() is itself getIdentifier().hashCode(), a content-derived digest hash
// with no meaningful collision risk in this context, so comparing full equality is behaviourally
// equivalent and avoids needing a ported hashCode() at all (Go's map/struct comparisons in this
// port key on content, not numeric hash codes - see model.Identifier.String()/AsXmlID()).
package jades

import (
	"github.com/ryftcore/dss-go/dss/spi/validation"
)

// JAdESSignatureIdentifierBuilder builds the DSS identifier for a JAdES signature. Port of the
// class SignatureIdentifierBuilder, extending spi/validation.AbstractSignatureIdentifierBuilder.
type SignatureIdentifierBuilder struct {
	validation.AbstractSignatureIdentifierBuilder

	// jadesSignature is the signature to build the identifier for, typed as the concrete JAdES
	// signature so CounterSignaturePosition/SignaturePosition need no runtime cast (Java casts
	// its inherited `signature` field instead, since its base class only knows AdvancedSignature).
	jadesSignature *Signature
}

// NewJAdESSignatureIdentifierBuilder is the default constructor. Port of the public
// SignatureIdentifierBuilder(Signature) constructor.
func NewSignatureIdentifierBuilder(signature *Signature) *SignatureIdentifierBuilder {
	b := &SignatureIdentifierBuilder{
		AbstractSignatureIdentifierBuilder: validation.NewAbstractSignatureIdentifierBuilderBase(signature),
		jadesSignature:                     signature,
	}
	b.InitAbstractSignatureIdentifierBuilder(b)
	return b
}

// CounterSignaturePosition returns the current counter signature position in its master
// signature. Port of the protected getCounterSignaturePosition(AdvancedSignature) override; see
// the file header DEVIATION note.
func (b *SignatureIdentifierBuilder) CounterSignaturePosition(masterSignature validation.AdvancedSignature) any {
	jadesSignature := b.jadesSignature
	jadesMasterSignature := masterSignature.(*Signature)
	masterCSigAttribute := jadesSignature.MasterCSigComponent()

	counter := 0
	if masterCSigAttribute != nil {
		masterCSigID := masterCSigAttribute.Identifier()
		for _, counterSignature := range jadesMasterSignature.CounterSignatures() {
			jadesCounterSignature := counterSignature.(*Signature)
			counterCSigID := jadesCounterSignature.MasterCSigComponent().Identifier()
			if masterCSigID.Equals(&counterCSigID) {
				break
			}
			counter++
		}
	}

	return counter
}

// SignaturePosition returns a position of a signature in the provided file. Port of the
// protected getSignaturePosition() override.
func (b *SignatureIdentifierBuilder) SignaturePosition() any {
	currentJWS := b.jadesSignature.Jws()
	jwsJsonSerializationObject := currentJWS.JwsJsonSerializationObject()

	counter := 0
	if jwsJsonSerializationObject != nil {
		for _, jws := range jwsJsonSerializationObject.Signatures() {
			if currentJWS == jws {
				break
			}
			counter++
		}
	}

	return counter
}
