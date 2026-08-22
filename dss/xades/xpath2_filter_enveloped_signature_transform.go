// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/reference/XPath2FilterEnvelopedSignatureTransform.java (DSS 6.5.RC1).
package xades

import "github.com/ryftcore/dss-go/dss/xml/common"

const (
	// xPath2FilterEnvelopedSignatureTransformSubtractFilter is the subtract filter. Port of the
	// private SUBTRACT_FILTER constant.
	xPath2FilterEnvelopedSignatureTransformSubtractFilter = "subtract"

	// xPath2FilterEnvelopedSignatureTransformDescendantSignaturePrefix is the all-descendant
	// ds:Signature elements prefix. Port of the private DESCENDANT_SIGNATURE_PREFIX.
	xPath2FilterEnvelopedSignatureTransformDescendantSignaturePrefix = "/descendant::"

	// xPath2FilterEnvelopedSignatureTransformDescendantSignatureSuffix is the all-descendant
	// ds:Signature elements suffix. Port of the private DESCENDANT_SIGNATURE_SUFFIX.
	xPath2FilterEnvelopedSignatureTransformDescendantSignatureSuffix = ":Signature"
)

// XPath2FilterEnvelopedSignatureTransform is the enveloped signature transformation by
// Filter 2.0. It excludes all signatures from the XML.
//
// Java declares the class final; Go has no final types, so the constraint survives only as
// this note - nothing in the port embeds it.
type XPath2FilterEnvelopedSignatureTransform struct {
	XPath2FilterTransform
}

// NewXPath2FilterEnvelopedSignatureTransform ports the default constructor, which delegates
// with XMLDSigNamespace.NS.
func NewXPath2FilterEnvelopedSignatureTransform() *XPath2FilterEnvelopedSignatureTransform {
	return NewXPath2FilterEnvelopedSignatureTransformWithNamespace(common.XMLDSigNS)
}

// NewXPath2FilterEnvelopedSignatureTransformWithNamespace ports
// XPath2FilterEnvelopedSignatureTransform(DSSNamespace). The expression is built from the
// namespace's own prefix, so a signature that uses a non-default xmldsig prefix subtracts on
// that prefix.
func NewXPath2FilterEnvelopedSignatureTransformWithNamespace(
	xmlDSigNamespace *common.DSSNamespace) *XPath2FilterEnvelopedSignatureTransform {
	t := new(XPath2FilterEnvelopedSignatureTransform)
	*t = XPath2FilterEnvelopedSignatureTransform{
		XPath2FilterTransform: newXPath2FilterTransform(xmlDSigNamespace,
			xPath2FilterEnvelopedSignatureTransformDescendantSignaturePrefix+xmlDSigNamespace.Prefix()+
				xPath2FilterEnvelopedSignatureTransformDescendantSignatureSuffix,
			xPath2FilterEnvelopedSignatureTransformSubtractFilter),
	}
	t.InitAbstractTransform(t)
	return t
}
