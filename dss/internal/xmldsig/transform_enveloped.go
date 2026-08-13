// Ported from org.apache.xml.security.transforms.implementations.TransformEnvelopedSignature
// and its EnvelopedNodeFilter (Apache Santuario xmlsec 3.0.6).
package xmldsig

import (
	"errors"

	"github.com/utain/esig/dss/internal/xmldom"
)

// envelopedSignatureTransform removes the enclosing ds:Signature from the input.
// Port of TransformEnvelopedSignature.
//
// It does two things and neither of them makes a node set: it names the ds:Signature as the
// input's exclude node, and it attaches a filter. The exclude node is what a following
// canonicalization transform uses when the input is still an ELEMENT (the common case: a
// same-document "" reference gives a document node, so the chain is
// enveloped-signature then c14n over a subtree with a hole in it). The filter is what a
// following canonicalization uses when some earlier transform already turned the input into a
// node set. Setting both, and NOT calling setNodeSet, is exactly what Santuario does - the
// transform alone never changes the input's kind.
type envelopedSignatureTransform struct{}

func (envelopedSignatureTransform) Algorithm() string { return TransformEnvelopedSignature }

func (envelopedSignatureTransform) Perform(in *Data, element *xmldom.Node, baseURI string, secureValidation bool) (*Data, error) {
	sig, err := searchSignatureElement(element)
	if err != nil {
		return nil, err
	}
	in.SetExcludeNode(sig)
	if err := in.AddNodeFilter(envelopedNodeFilter{exclude: sig}); err != nil {
		return nil, err
	}
	return in, nil
}

// searchSignatureElement walks up from the ds:Transform element to the ds:Signature that
// encloses it. Port of TransformEnvelopedSignature#searchSignatureElement, whose failure -
// "transform.envelopedSignatureTransformNotInSignatureElement" - is what a ds:Transforms
// element detached from its signature produces.
func searchSignatureElement(n *xmldom.Node) (*xmldom.Node, error) {
	for ; n != nil && n.Kind != xmldom.Document; n = n.Parent {
		if n.Name.Space == NamespaceDSig && n.Name.Local == "Signature" {
			return n, nil
		}
	}
	return nil, errors.New("xmldsig: the enveloped-signature transform is not inside a ds:Signature element")
}

// envelopedNodeFilter is EnvelopedNodeFilter. The two methods differ, and the difference is
// the point:
//
//   - IsNodeIncludeDO compares by identity only, because the traversal reaches nodes in
//     document order and the -1 it returns for the ds:Signature element prunes the whole
//     subtree before any descendant is ever offered;
//   - IsNodeInclude has to test descendant-or-self, because it is asked about nodes the
//     traversal did not arrive at from above - attributes, and the parent element again when
//     its end tag is due.
type envelopedNodeFilter struct{ exclude *xmldom.Node }

func (f envelopedNodeFilter) IsNodeIncludeDO(n *xmldom.Node, level int) (int, error) {
	if n == f.exclude {
		return -1, nil
	}
	return 1, nil
}

func (f envelopedNodeFilter) IsNodeInclude(n *xmldom.Node) (int, error) {
	if n == f.exclude || IsDescendantOrSelf(f.exclude, n) {
		return -1, nil
	}
	return 1, nil
}
