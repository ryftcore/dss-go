// Ported from org.apache.xml.security.transforms.implementations.TransformBase64Decode
// (Apache Santuario xmlsec 3.0.6).
package xmldsig

import (
	"strings"

	"github.com/ryftcore/dss-go/dss/internal/xmldom"
)

// base64Transform decodes its input from base64. Port of TransformBase64Decode.
//
// The element branch is the one XAdES actually uses - a ds:Reference pointing at a ds:Object
// that holds base64 text - and it is not "the element's text content": traverseElement walks
// the subtree and concatenates the data of every TEXT node it meets, at any depth, ignoring
// comments, processing instructions and CDATA. A CDATA section inside the ds:Object therefore
// contributes NOTHING, which is a real and reproducible difference from reading the text
// content, and an element nested inside contributes its text as if the element were not there.
type base64Transform struct{}

func (base64Transform) Algorithm() string { return TransformBase64Decode }

func (base64Transform) Perform(in *Data, element *xmldom.Node, baseURI string, secureValidation bool) (*Data, error) {
	if in.IsElement() {
		el := in.Node()
		if el.Kind == xmldom.Text {
			el = el.Parent
		}
		var sb strings.Builder
		traverseElement(el, &sb)
		decoded, err := decodeBase64(sb.String())
		if err != nil {
			return nil, err
		}
		return NewOctetData(decoded), nil
	}
	if in.IsOctetStream() || in.IsNodeSet() {
		b, err := in.Bytes()
		if err != nil {
			return nil, err
		}
		decoded, err := decodeBase64(string(b))
		if err != nil {
			return nil, err
		}
		return NewOctetData(decoded), nil
	}
	return nil, ErrUninitializedData
}

// traverseElement ports TransformBase64Decode#traverseElement: recurse into element children,
// append the data of text children, ignore everything else - CDATA included.
func traverseElement(n *xmldom.Node, sb *strings.Builder) {
	if n == nil {
		return
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		switch c.Kind {
		case xmldom.Element:
			traverseElement(c, sb)
		case xmldom.Text:
			sb.WriteString(c.Value)
		}
	}
}
