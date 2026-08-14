// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/validation/scope/XmlFullSignatureScope.java (DSS 6.5.RC1).
package xades

import (
	"github.com/utain/esig/dss/model"
	mscope "github.com/utain/esig/dss/model/scope"
)

// XmlFullSignatureScope defines a full XML document signature scope. Port of the class
// XmlFullSignatureScope, extending XmlRootSignatureScope.
//
// Java declares this class final; there is no Go language-level equivalent, so it is
// documented here instead - no other file in this package embeds XmlFullSignatureScope.
type XmlFullSignatureScope struct {
	XmlRootSignatureScope
}

// newXmlFullSignatureScope is the port of the protected XmlFullSignatureScope(String,
// DSSDocument, List<String>) constructor.
func newXmlFullSignatureScope(name string, document model.DSSDocument, transformations []string) *XmlFullSignatureScope {
	return &XmlFullSignatureScope{
		XmlRootSignatureScope: *newXmlRootSignatureScopeWithName(name, document, transformations),
	}
}

// compile-time interface assertion.
var _ mscope.SignatureScope = (*XmlFullSignatureScope)(nil)
