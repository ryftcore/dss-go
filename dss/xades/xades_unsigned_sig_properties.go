// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/validation/XAdESUnsignedSigProperties.java
// (DSS 6.5.RC1).
package xades

import (
	"github.com/ryftcore/dss-go/dss/internal/xmldom"
	"github.com/ryftcore/dss-go/dss/xades/definition"
	"github.com/ryftcore/dss-go/dss/xml/utils"
)

// XAdESUnsignedSigProperties represents unsigned XAdES signature properties. Port of the class
// UnsignedSigProperties, extending SigProperties.
type UnsignedSigProperties struct {
	SigProperties
}

// NewXAdESUnsignedSigProperties is the port of the public
// UnsignedSigProperties(Element, XAdESPath) constructor.
func NewUnsignedSigProperties(unsignedSignatureProperties *xmldom.Node, xadesPaths definition.XAdESPath) *UnsignedSigProperties {
	return &UnsignedSigProperties{
		SigProperties: newSigProperties(unsignedSignatureProperties, xadesPaths),
	}
}

// XAdESUnsignedSigPropertiesBuild builds a XAdESUnsignedSigProperties. Port of the static
// build(Element, XAdESPath).
func UnsignedSigPropertiesBuild(signatureElement *xmldom.Node, xadesPaths definition.XAdESPath) *UnsignedSigProperties {
	unsignedSignatureProperties := xadesUnsignedSigPropertiesGetUnsignedSignaturePropertiesDom(signatureElement, xadesPaths)
	return NewUnsignedSigProperties(unsignedSignatureProperties, xadesPaths)
}

// xadesUnsignedSigPropertiesGetUnsignedSignaturePropertiesDom gets the
// xades:UnsignedSignatureProperties element. Port of the protected static
// getUnsignedSignaturePropertiesDom(Element, XAdESPath).
func xadesUnsignedSigPropertiesGetUnsignedSignaturePropertiesDom(signatureElement *xmldom.Node, xadesPaths definition.XAdESPath) *xmldom.Node {
	element, err := utils.XPathUtilsGetElement(signatureElement, xadesPaths.UnsignedSignaturePropertiesPath())
	if err != nil {
		return nil
	}
	return element
}
