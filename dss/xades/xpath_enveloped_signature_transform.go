// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/reference/XPathEnvelopedSignatureTransform.java (DSS 6.5.RC1).
package xades

import "github.com/utain/esig/dss/xml/common"

const (
	// xPathEnvelopedSignatureTransformNotAncestorOrSelfPrefix is the XPath filter that removes
	// all ds:Signature elements from the XML. Port of the private
	// NOT_ANCESTOR_OR_SELF_PREFIX constant.
	xPathEnvelopedSignatureTransformNotAncestorOrSelfPrefix = "not(ancestor-or-self::"

	// xPathEnvelopedSignatureTransformSignatureSuffix ports the private SIGNATURE_SUFFIX.
	xPathEnvelopedSignatureTransformSignatureSuffix = ":Signature)"
)

// XPathEnvelopedSignatureTransform is the simple enveloped signature transform.
//
// WARN: cannot be used with parallel signatures!
//
// Java declares the class final; Go has no final types, so the constraint survives only as
// this note - nothing in the port embeds it.
type XPathEnvelopedSignatureTransform struct {
	XPathTransform
}

// NewXPathEnvelopedSignatureTransform ports the default constructor, which delegates with
// XMLDSigNamespace.NS.
func NewXPathEnvelopedSignatureTransform() *XPathEnvelopedSignatureTransform {
	return NewXPathEnvelopedSignatureTransformWithNamespace(common.XMLDSigNS)
}

// NewXPathEnvelopedSignatureTransformWithNamespace ports
// XPathEnvelopedSignatureTransform(DSSNamespace). The expression is built from the namespace's
// own prefix, so a signature that uses a non-default xmldsig prefix filters on that prefix.
func NewXPathEnvelopedSignatureTransformWithNamespace(
	xmlDSigNamespace *common.DSSNamespace) *XPathEnvelopedSignatureTransform {
	t := new(XPathEnvelopedSignatureTransform)
	*t = XPathEnvelopedSignatureTransform{
		XPathTransform: *NewXPathTransformWithNamespace(xmlDSigNamespace,
			xPathEnvelopedSignatureTransformNotAncestorOrSelfPrefix+xmlDSigNamespace.Prefix()+
				xPathEnvelopedSignatureTransformSignatureSuffix),
	}
	t.InitAbstractTransform(t)
	return t
}
