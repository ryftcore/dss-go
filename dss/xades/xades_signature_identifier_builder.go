// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/validation/XAdESSignatureIdentifierBuilder.java
// (DSS 6.5.RC1).
//
// XAdESDOMDocument (Java eu.europa.esig.dss.xades.dom.XAdESDOMDocument) is a forward dependency
// of the sibling "dom" chunk of this same SCC (see xades_signature.go's file header for the full
// contract); SignatureNodes() []*xmldom.Node is the only member this file needs.
package xades

import (
	"strings"

	"github.com/utain/esig/dss/spi/validation"
	"github.com/utain/esig/dss/utils"
)

// XAdESSignatureIdentifierBuilder builds a DSS identifier for a XAdES signature. Port of the
// class XAdESSignatureIdentifierBuilder, extending
// spi/validation.AbstractSignatureIdentifierBuilder.
type XAdESSignatureIdentifierBuilder struct {
	validation.AbstractSignatureIdentifierBuilder

	// xadesSignature is the signature to build the identifier for, typed as the concrete XAdES
	// signature so CounterSignaturePosition/SignaturePosition/SignatureFilePosition need no
	// runtime cast (Java casts its inherited `signature` field instead, since its base class only
	// knows AdvancedSignature).
	xadesSignature *XAdESSignature
}

// NewXAdESSignatureIdentifierBuilder is the port of the constructor
// XAdESSignatureIdentifierBuilder(XAdESSignature).
func NewXAdESSignatureIdentifierBuilder(signature *XAdESSignature) *XAdESSignatureIdentifierBuilder {
	b := &XAdESSignatureIdentifierBuilder{
		AbstractSignatureIdentifierBuilder: validation.NewAbstractSignatureIdentifierBuilderBase(signature),
		xadesSignature:                     signature,
	}
	b.InitAbstractSignatureIdentifierBuilder(b)
	return b
}

// CounterSignaturePosition returns the current counter signature position in its master
// signature. Port of the protected getCounterSignaturePosition(AdvancedSignature) override.
func (b *XAdESSignatureIdentifierBuilder) CounterSignaturePosition(masterSignature validation.AdvancedSignature) any {
	xadesMasterSignature := masterSignature.(*XAdESSignature)
	xadesSignature := b.xadesSignature

	counter := 0
	for _, counterSignature := range xadesMasterSignature.CounterSignatures() {
		xadesCounterSignature := counterSignature.(*XAdESSignature)
		if xadesSignature.SignatureElement() == xadesCounterSignature.SignatureElement() {
			break
		}
		counter++
	}

	return counter
}

// SignaturePosition returns a position of a signature in the provided file. Port of the
// protected getSignaturePosition() override.
func (b *XAdESSignatureIdentifierBuilder) SignaturePosition() any {
	xadesSignature := b.xadesSignature
	signatureElement := xadesSignature.SignatureElement()
	ownerDocument := xadesSignature.OwnerDocument()
	signatureNodeList := ownerDocument.SignatureNodes()

	counter := 0
	for counter < len(signatureNodeList) {
		if signatureElement == signatureNodeList[counter] {
			break
		}
		counter++
	}

	return counter
}

// SignatureFilePosition returns a position of a signature file. Port of the protected
// getSignatureFilePosition() override.
func (b *XAdESSignatureIdentifierBuilder) SignatureFilePosition() any {
	signatureFilename := b.Signature().Filename()
	if utils.IsStringNotEmpty(signatureFilename) && strings.HasPrefix(signatureFilename, validation.AbstractSignatureIdentifierBuilderMETAINFFolder) {
		return signatureFilename
	}
	return utils.EmptyString
}
