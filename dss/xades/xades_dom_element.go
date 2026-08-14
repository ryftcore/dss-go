// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/dom/XAdESDOMElement.java (DSS 6.5.RC1).
package xades

import (
	"github.com/utain/esig/dss/internal/xmldom"
	"github.com/utain/esig/dss/xades/definition"
)

// XAdESDOMElement represents a wrapper for an *xmldom.Node Element object for a XAdES signature.
type XAdESDOMElement struct {
	// ownerDocument is the owner document.
	ownerDocument *XAdESDOMDocument

	// element is the XML DOM element.
	element *xmldom.Node
}

// NewXAdESDOMElement is the default constructor.
func NewXAdESDOMElement(element *xmldom.Node, ownerDocument *XAdESDOMDocument) *XAdESDOMElement {
	return &XAdESDOMElement{element: element, ownerDocument: ownerDocument}
}

// Element gets the XML DOM Element. Ports getElement().
func (e *XAdESDOMElement) Element() *xmldom.Node {
	return e.element
}

// OwnerDocument gets the owner document. Ports getOwnerDocument().
func (e *XAdESDOMElement) OwnerDocument() *XAdESDOMDocument {
	return e.ownerDocument
}

// XAdESPathHolders gets a list of registered XAdES Path holders. Ports getXAdESPathHolders().
func (e *XAdESDOMElement) XAdESPathHolders() []definition.XAdESPath {
	return e.ownerDocument.XAdESPathHolders()
}
