// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/validation/XAdESSigProperties.java (DSS 6.5.RC1).
package xades

import (
	"github.com/ryftcore/dss-go/dss/internal/xmldom"
	"github.com/ryftcore/dss-go/dss/xades/definition"
)

// XAdESSigProperties represents a list of XAdESAttributes. Port of the abstract class
// SigProperties, implementing spi/validation.SignatureProperties[*Attribute].
type SigProperties struct {
	// signaturePropertiesDom is the signature properties element.
	signaturePropertiesDom *xmldom.Node

	// xadesPaths is the XAdES XPaths.
	xadesPaths definition.XAdESPath
}

// newXAdESSigProperties is the port of the package-private
// SigProperties(Element, XAdESPath) constructor.
func newSigProperties(signatureProperties *xmldom.Node, xadesPaths definition.XAdESPath) SigProperties {
	return SigProperties{signaturePropertiesDom: signatureProperties, xadesPaths: xadesPaths}
}

// IsExist is the port of isExist().
func (p *SigProperties) IsExist() bool {
	return p.signaturePropertiesDom != nil
}

// Attributes is the port of getAttributes(). Elements() already restricts the traversal to
// element children, which is what the Java loop's isElementNode(Node) guard achieves by hand.
func (p *SigProperties) Attributes() []*Attribute {
	unsignedAttributes := make([]*Attribute, 0)
	if p.signaturePropertiesDom != nil {
		for _, node := range p.signaturePropertiesDom.Elements() {
			unsignedAttributes = append(unsignedAttributes, newAttribute(node, p.xadesPaths))
		}
	}
	return unsignedAttributes
}
