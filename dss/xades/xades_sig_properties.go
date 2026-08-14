// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/validation/XAdESSigProperties.java (DSS 6.5.RC1).
package xades

import (
	"github.com/utain/esig/dss/internal/xmldom"
	"github.com/utain/esig/dss/xades/definition"
)

// XAdESSigProperties represents a list of XAdESAttributes. Port of the abstract class
// XAdESSigProperties, implementing spi/validation.SignatureProperties[*XAdESAttribute].
type XAdESSigProperties struct {
	// signaturePropertiesDom is the signature properties element.
	signaturePropertiesDom *xmldom.Node

	// xadesPaths is the XAdES XPaths.
	xadesPaths definition.XAdESPath
}

// newXAdESSigProperties is the port of the package-private
// XAdESSigProperties(Element, XAdESPath) constructor.
func newXAdESSigProperties(signatureProperties *xmldom.Node, xadesPaths definition.XAdESPath) XAdESSigProperties {
	return XAdESSigProperties{signaturePropertiesDom: signatureProperties, xadesPaths: xadesPaths}
}

// IsExist is the port of isExist().
func (p *XAdESSigProperties) IsExist() bool {
	return p.signaturePropertiesDom != nil
}

// Attributes is the port of getAttributes(). Elements() already restricts the traversal to
// element children, which is what the Java loop's isElementNode(Node) guard achieves by hand.
func (p *XAdESSigProperties) Attributes() []*XAdESAttribute {
	unsignedAttributes := make([]*XAdESAttribute, 0)
	if p.signaturePropertiesDom != nil {
		for _, node := range p.signaturePropertiesDom.Elements() {
			unsignedAttributes = append(unsignedAttributes, newXAdESAttribute(node, p.xadesPaths))
		}
	}
	return unsignedAttributes
}
