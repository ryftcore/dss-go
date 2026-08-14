// Ported from org.apache.xml.security.transforms.implementations.TransformC14N,
// TransformC14NWithComments, TransformC14N11, TransformC14N11_WithComments,
// TransformC14NExclusive and TransformC14NExclusiveWithComments (Apache Santuario xmlsec 3.0.6).
package xmldsig

import (
	"bytes"

	"github.com/utain/esig/dss/internal/xmlc14n"
	"github.com/utain/esig/dss/internal/xmldom"
)

// c14nTransform is any of the canonicalization methods used as a ds:Transform. The six
// Santuario classes differ only in which canonicalizer they construct and, for the two
// exclusive ones, in reading the ec:InclusiveNamespaces PrefixList; both differences are data
// here, so one type covers all of them.
//
// The transform always yields OCTETS - "new XMLSignatureInput(writer.toByteArray())" - which
// is what ends a transform chain: whatever node-set filtering happened earlier is resolved
// here, once, and nothing downstream can see the node set again.
type c14nTransform struct{ alg string }

func (t c14nTransform) Algorithm() string { return t.alg }

func (t c14nTransform) Perform(in *Data, element *xmldom.Node, baseURI string, secureValidation bool) (*Data, error) {
	var prefixes []string
	if t.alg == TransformC14NExcl || t.alg == TransformC14NExclWithComments {
		// TransformC14NExclusive: exactly one ec:InclusiveNamespaces child is read; zero
		// leaves the PrefixList null, and more than one is ignored (the length() == 1 guard),
		// so a second element silently disables the parameter. Reproduced, not fixed.
		incl := selectNodes(element, NamespaceExcC14N, "InclusiveNamespaces")
		if len(incl) == 1 {
			prefixes = xmlc14n.ParsePrefixList(incl[0].AttrValue("", "PrefixList"))
		}
	}
	var buf bytes.Buffer
	if err := in.canonicalize(xmlc14n.Algorithm(t.alg), prefixes, &buf); err != nil {
		return nil, err
	}
	// Santuario does not carry the source URI across a canonicalization transform, and the
	// XPath Filter 2.0 transform reads it. Reproduced rather than improved.
	return NewOctetData(buf.Bytes()), nil
}
