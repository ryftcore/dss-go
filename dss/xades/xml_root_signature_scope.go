// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/validation/scope/XmlRootSignatureScope.java (DSS 6.5.RC1).
package xades

import (
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/model"
	mscope "github.com/utain/esig/dss/model/scope"
)

// XmlRootSignatureScope defines a root XML document signature scope. Port of the class
// XmlRootSignatureScope, extending XmlElementSignatureScope.
type XmlRootSignatureScope struct {
	XmlElementSignatureScope
}

// newXmlRootSignatureScope is the port of the protected XmlRootSignatureScope(DSSDocument,
// List<String>) constructor.
func newXmlRootSignatureScope(document model.DSSDocument, transformations []string) *XmlRootSignatureScope {
	return newXmlRootSignatureScopeWithName("Full XML File", document, transformations)
}

// newXmlRootSignatureScopeWithName is the port of the protected XmlRootSignatureScope(String,
// DSSDocument, List<String>) constructor.
func newXmlRootSignatureScopeWithName(name string, document model.DSSDocument, transformations []string) *XmlRootSignatureScope {
	return &XmlRootSignatureScope{
		XmlElementSignatureScope: *newXmlElementSignatureScope(name, document, transformations),
	}
}

// Description returns the XmlRootSignatureScope description. Port of
// getDescription(TokenIdentifierProvider).
func (s *XmlRootSignatureScope) Description(tokenIdentifierProvider model.TokenIdentifierProvider) string {
	return s.AddTransformationIfNeeded("The full XML file")
}

// Type returns the type of the signature scope. Port of getType().
func (s *XmlRootSignatureScope) Type() enumerations.SignatureScopeType {
	return enumerations.SignatureScopeType_FULL
}

// compile-time interface assertion.
var _ mscope.SignatureScope = (*XmlRootSignatureScope)(nil)
