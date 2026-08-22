// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/validation/XAdESSignedDataObjectProperties.java
// (DSS 6.5.RC1).
//
// XAdESSigProperties (Java eu.europa.esig.dss.xades.validation.XAdESSigProperties, the same Go
// package xades) lives alongside this file as xades_sig_properties.go: abstract in Java,
// implementing spi.validation.SignatureProperties[Attribute] over a signature-properties DOM
// element and an XAdESPath. Its constructor is the package-private
// newSigProperties(signatureProperties *xmldom.Node, xadesPaths definition.XAdESPath)
// SigProperties, called directly below - Go has no abstract-class dispatch to register back
// into, and SigProperties calls nothing virtual on its embedder (IsExist/Attributes are
// entirely self-contained), so no Init-style registration is needed.
package xades

import (
	"github.com/ryftcore/dss-go/dss/internal/xmldom"
	"github.com/ryftcore/dss-go/dss/xades/definition"
	"github.com/ryftcore/dss-go/dss/xml/utils"
)

// XAdESSignedDataObjectProperties builds XAdESSignedDataObjectProperties. Port of the class
// SignedDataObjectProperties, extending SigProperties.
type SignedDataObjectProperties struct {
	SigProperties
}

// newXAdESSignedDataObjectProperties is the port of the package-private
// SignedDataObjectProperties(Element, XAdESPath) constructor.
func newSignedDataObjectProperties(signatureProperties *xmldom.Node, xadesPaths definition.XAdESPath) *SignedDataObjectProperties {
	return &SignedDataObjectProperties{
		SigProperties: newSigProperties(signatureProperties, xadesPaths),
	}
}

// XAdESSignedDataObjectPropertiesBuild builds a XAdESSignedDataObjectProperties. Port of the
// static build(Element, XAdESPath).
func SignedDataObjectPropertiesBuild(signatureElement *xmldom.Node, xadesPaths definition.XAdESPath) *SignedDataObjectProperties {
	signedSignatureProperties := xadesSignedDataObjectPropertiesGetSignedSignaturePropertiesDom(signatureElement, xadesPaths)
	return newSignedDataObjectProperties(signedSignatureProperties, xadesPaths)
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
