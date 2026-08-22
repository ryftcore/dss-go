// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/reference/DSSTransformOutput.java (DSS 6.5.RC1).
//
// Java's DSSTransformOutput wraps org.apache.xml.security.signature.XMLSignatureInput. Per the
// internal/xmldsig doc.go mapping table, XMLSignatureInput is xmldsig.Data and
// XMLSignatureInput#getBytes is Data.Bytes, so this wrapper carries a *xmldsig.Data.
package xades

import (
	"github.com/ryftcore/dss-go/dss/internal/xmldom"
	"github.com/ryftcore/dss-go/dss/internal/xmldsig"
	"github.com/ryftcore/dss-go/dss/model"
)

// DSSTransformOutput represents an xmldsig.Data (Santuario XMLSignatureInput) wrapper.
type DSSTransformOutput struct {
	// data is the cached XMLSignatureInput.
	data *xmldsig.Data
}

// NewDSSTransformOutput instantiates the object from an xmldsig.Data.
// Ports DSSTransformOutput(XMLSignatureInput).
func NewDSSTransformOutput(data *xmldsig.Data) *DSSTransformOutput {
	return &DSSTransformOutput{data: data}
}

// NewDSSTransformOutputFromNode instantiates the object from a Node.
// Ports DSSTransformOutput(Node), which delegates to new XMLSignatureInput(node).
func NewDSSTransformOutputFromNode(node *xmldom.Node) *DSSTransformOutput {
	return NewDSSTransformOutput(xmldsig.NewNodeData(node))
}

// xmlSignatureInput returns the wrapped xmldsig.Data. Ports the protected
// getXmlSignatureInput(); Java's "protected" is package visibility here, which is what
// ComplexTransform - the only caller - needs.
func (o *DSSTransformOutput) xmlSignatureInput() *xmldsig.Data {
	return o.data
}

// Bytes returns the bytes after performing transforms. Ports getBytes(), whose IOException /
// XMLSecurityException wrapping in a DSSException becomes the returned *model.DSSError with
// the same message.
func (o *DSSTransformOutput) Bytes() ([]byte, error) {
	b, err := o.data.Bytes()
	if err != nil {
		return nil, model.NewDSSErrorMessageCause(
			"Cannot extract Transform output bytes. Reason : ["+err.Error()+"]", err)
	}
	return b, nil
}
