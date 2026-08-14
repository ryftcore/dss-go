// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/validation/XAdESSignedDataObjectProperties.java
// (DSS 6.5.RC1).
//
// XAdESSigProperties (Java eu.europa.esig.dss.xades.validation.XAdESSigProperties, same
// "validation" SCC as this file per S4D_BRIEF.md's package-layout rule, hence the same Go
// package xades) is landed alongside this file as xades_sig_properties.go: abstract in Java,
// implementing spi.validation.SignatureProperties[XAdESAttribute] over a signature-properties DOM
// element and an XAdESPath. Its constructor is the package-private
// newXAdESSigProperties(signatureProperties *xmldom.Node, xadesPaths definition.XAdESPath)
// XAdESSigProperties, called directly below - Go has no abstract-class dispatch to register back
// into, and XAdESSigProperties calls nothing virtual on its embedder (IsExist/Attributes are
// entirely self-contained), so no Init-style registration is needed.
package xades

import (
	"github.com/utain/esig/dss/internal/xmldom"
	"github.com/utain/esig/dss/xades/definition"
	"github.com/utain/esig/dss/xml/utils"
)

// XAdESSignedDataObjectProperties builds XAdESSignedDataObjectProperties. Port of the class
// XAdESSignedDataObjectProperties, extending XAdESSigProperties.
type XAdESSignedDataObjectProperties struct {
	XAdESSigProperties
}

// newXAdESSignedDataObjectProperties is the port of the package-private
// XAdESSignedDataObjectProperties(Element, XAdESPath) constructor.
func newXAdESSignedDataObjectProperties(signatureProperties *xmldom.Node, xadesPaths definition.XAdESPath) *XAdESSignedDataObjectProperties {
	return &XAdESSignedDataObjectProperties{
		XAdESSigProperties: newXAdESSigProperties(signatureProperties, xadesPaths),
	}
}

// XAdESSignedDataObjectPropertiesBuild builds a XAdESSignedDataObjectProperties. Port of the
// static build(Element, XAdESPath).
func XAdESSignedDataObjectPropertiesBuild(signatureElement *xmldom.Node, xadesPaths definition.XAdESPath) *XAdESSignedDataObjectProperties {
	signedSignatureProperties := xadesSignedDataObjectPropertiesGetSignedSignaturePropertiesDom(signatureElement, xadesPaths)
	return newXAdESSignedDataObjectProperties(signedSignatureProperties, xadesPaths)
}

// xadesSignedDataObjectPropertiesGetSignedSignaturePropertiesDom gets the
// xades:SignedDataObjectProperties element. Port of the protected static
// getSignedSignaturePropertiesDom(Element, XAdESPath).
func xadesSignedDataObjectPropertiesGetSignedSignaturePropertiesDom(signatureElement *xmldom.Node, xadesPaths definition.XAdESPath) *xmldom.Node {
	element, err := utils.XPathUtilsGetElement(signatureElement, xadesPaths.SignedDataObjectPropertiesPath())
	if err != nil {
		return nil
	}
	return element
}
