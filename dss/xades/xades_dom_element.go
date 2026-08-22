// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/dom/XAdESDOMElement.java (DSS 6.5.RC1).
package xades

import (
	"github.com/ryftcore/dss-go/dss/internal/xmldom"
	"github.com/ryftcore/dss-go/dss/xades/definition"
)

// DOMElement represents a wrapper for an *xmldom.Node Element object for a XAdES signature.
type DOMElement struct {
	// ownerDocument is the owner document.
	ownerDocument *DOMDocument

	// element is the XML DOM element.
	element *xmldom.Node
}

// NewDOMElement is the default constructor.
func NewDOMElement(element *xmldom.Node, ownerDocument *DOMDocument) *DOMElement {
	return &DOMElement{element: element, ownerDocument: ownerDocument}
}

// Element gets the XML DOM Element. Ports getElement().
func (e *DOMElement) Element() *xmldom.Node {
	return e.element
}

// OwnerDocument gets the owner document. Ports getOwnerDocument().
func (e *DOMElement) OwnerDocument() *DOMDocument {
	return e.ownerDocument
}

// XAdESPathHolders gets a list of registered XAdES Path holders. Ports getXAdESPathHolders().
func (e *DOMElement) XAdESPathHolders() []definition.XAdESPath {
	return e.ownerDocument.XAdESPathHolders()
}
