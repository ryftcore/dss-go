// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/reference/Base64Transform.java (DSS 6.5.RC1).
//
// org.apache.xml.security.transforms.Transforms.TRANSFORM_BASE64_DECODE is
// xmldsig.TransformBase64Decode, per internal/xmldsig's doc.go mapping table.
package xades

import (
	"github.com/utain/esig/dss/internal/xmldsig"
	"github.com/utain/esig/dss/xml/common"
)

// Base64Transform transforms a reference content to its base64 representation.
//
// NOTE: Not compatible with:
//   - other transformations;
//   - isEmbed(true) parameter;
//   - Manifest signature;
//   - Enveloped signatures.
type Base64Transform struct {
	AbstractTransform
}

// NewBase64Transform ports the default constructor Base64Transform().
func NewBase64Transform() *Base64Transform {
	t := &Base64Transform{AbstractTransform: newAbstractTransform(xmldsig.TransformBase64Decode)}
	t.InitAbstractTransform(t)
	return t
}

// NewBase64TransformWithNamespace ports Base64Transform(DSSNamespace).
func NewBase64TransformWithNamespace(xmlDSigNamespace *common.DSSNamespace) *Base64Transform {
	t := &Base64Transform{
		AbstractTransform: newAbstractTransformWithNamespace(xmlDSigNamespace, xmldsig.TransformBase64Decode),
	}
	t.InitAbstractTransform(t)
	return t
}

// PerformTransform extracts base64-decoded content from a Reference directly, i.e. it returns
// its input untouched. Ports performTransform(DSSTransformOutput).
func (t *Base64Transform) PerformTransform(transformOutput *DSSTransformOutput) (*DSSTransformOutput, error) {
	return transformOutput, nil
}
