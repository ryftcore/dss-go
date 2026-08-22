// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/reference/EnvelopedSignatureTransform.java (DSS 6.5.RC1).
//
// org.apache.xml.security.transforms.Transforms.TRANSFORM_ENVELOPED_SIGNATURE is
// xmldsig.TransformEnvelopedSignature, per internal/xmldsig's doc.go mapping table.
package xades

import (
	"github.com/ryftcore/dss-go/dss/internal/xmldsig"
	"github.com/ryftcore/dss-go/dss/xml/common"
)

// EnvelopedSignatureTransform is used for an Enveloped Signature.
//
// Note: must be followed up by a CanonicalizationTransform.
type EnvelopedSignatureTransform struct {
	AbstractTransform
}

// NewEnvelopedSignatureTransform ports the default constructor
// EnvelopedSignatureTransform(), which - unlike its AbstractTransform siblings - passes
// XMLDSigNamespace.NS explicitly rather than relying on the field initializer.
func NewEnvelopedSignatureTransform() *EnvelopedSignatureTransform {
	t := &EnvelopedSignatureTransform{
		AbstractTransform: newAbstractTransformWithNamespace(common.XMLDSigNS, xmldsig.TransformEnvelopedSignature),
	}
	t.InitAbstractTransform(t)
	return t
}

// NewEnvelopedSignatureTransformWithNamespace ports
// EnvelopedSignatureTransform(DSSNamespace).
func NewEnvelopedSignatureTransformWithNamespace(xmlDSigNamespace *common.DSSNamespace) *EnvelopedSignatureTransform {
	t := &EnvelopedSignatureTransform{
		AbstractTransform: newAbstractTransformWithNamespace(xmlDSigNamespace, xmldsig.TransformEnvelopedSignature),
	}
	t.InitAbstractTransform(t)
	return t
}

// PerformTransform does nothing: on signature creation the new signature does not exist yet,
// so there is nothing to envelope away. Ports performTransform(DSSTransformOutput).
func (t *EnvelopedSignatureTransform) PerformTransform(transformOutput *DSSTransformOutput) (*DSSTransformOutput, error) {
	return transformOutput, nil
}
