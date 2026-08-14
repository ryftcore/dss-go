// Ported from org.apache.xml.security.utils.XMLUtils and org.apache.xml.security.utils.Constants
// (Apache Santuario xmlsec 3.0.6): the handful of DOM helpers the reference pipeline needs.
package xmldsig

import (
	"encoding/base64"
	"errors"
	"strings"

	"github.com/utain/esig/dss/internal/xmldom"
)

// Namespace URIs. Port of Constants.SignatureSpecNS and the XPath Filter 2.0 namespace of
// transforms.params.XPath2FilterContainer.
const (
	// NamespaceDSig is Constants.SignatureSpecNS.
	NamespaceDSig = "http://www.w3.org/2000/09/xmldsig#"
	// NamespaceExcC14N is InclusiveNamespaces.ExclusiveCanonicalizationNamespace.
	NamespaceExcC14N = "http://www.w3.org/2001/10/xml-exc-c14n#"
	// NamespaceXPathFilter2 is XPath2FilterContainer.XPathFilter2NS.
	NamespaceXPathFilter2 = "http://www.w3.org/2002/06/xmldsig-filter2"
)

// IsDescendantOrSelf reports whether descendantOrSelf is ctx or lies inside it. Port of
// XMLUtils#isDescendantOrSelf, including its one subtlety: an attribute's parent is the
// element that owns it, so an attribute of a ds:Signature counts as inside that signature even
// though it is not among its children.
func IsDescendantOrSelf(ctx, descendantOrSelf *xmldom.Node) bool {
	if ctx == descendantOrSelf {
		return true
	}
	for p := descendantOrSelf; p != nil; p = parentOf(p) {
		if p == ctx {
			return true
		}
	}
	return false
}

// parentOf is Node#getParentNode with the attribute special case. xmldom links an attribute to
// its element through Parent already, so the two collapse; the function exists to make the
// correspondence with the Java visible and to keep IsDescendantOrSelf a literal port.
func parentOf(n *xmldom.Node) *xmldom.Node { return n.Parent }

// NodeSetOf collects the nodes of a subtree, in document order, skipping exclude and its
// subtree. Port of XMLUtils#getSet/getSetRec: an element contributes itself and all of its
// attributes (namespace declarations included, which is why the exclusive canonicalizer can
// ask whether a declaration is in the set), a comment contributes itself only when comments
// are kept, and a DOCTYPE contributes nothing.
func NodeSetOf(root, exclude *xmldom.Node, comments bool) []*xmldom.Node {
	if root == nil {
		return nil
	}
	if exclude != nil && IsDescendantOrSelf(exclude, root) {
		return nil
	}
	var out []*xmldom.Node
	var rec func(n *xmldom.Node)
	rec = func(n *xmldom.Node) {
		if n == exclude {
			return
		}
		switch n.Kind {
		case xmldom.Element:
			out = append(out, n)
			out = append(out, n.Attrs...)
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				rec(c)
			}
		case xmldom.Document:
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				rec(c)
			}
		case xmldom.Comment:
			if comments {
				out = append(out, n)
			}
		default:
			out = append(out, n)
		}
	}
	rec(root)
	return out
}

// selectDSNodes returns the ds: element children of parent with the given local name, in
// document order. Port of XMLUtils#selectDsNodes.
func selectDSNodes(parent *xmldom.Node, local string) []*xmldom.Node {
	return selectNodes(parent, NamespaceDSig, local)
}

// selectNodes returns the element children of parent in the given namespace with the given
// local name. Port of XMLUtils#selectNodes.
func selectNodes(parent *xmldom.Node, space, local string) []*xmldom.Node {
	if parent == nil {
		return nil
	}
	var out []*xmldom.Node
	for c := parent.FirstChild; c != nil; c = c.NextSibling {
		if c.Kind == xmldom.Element && c.Name.Space == space && c.Name.Local == local {
			out = append(out, c)
		}
	}
	return out
}

// selectDSNode returns the index'th ds: element child with the given local name, or nil.
// Port of XMLUtils#selectDsNode.
func selectDSNode(parent *xmldom.Node, local string, index int) *xmldom.Node {
	nodes := selectDSNodes(parent, local)
	if index < len(nodes) {
		return nodes[index]
	}
	return nil
}

// stringFromNode is XMLUtils#getStrFromNode, the reader of a ds:XPath element's expression.
//
// Its shape matters: given a text node it concatenates the character data of EVERY text-node
// sibling, not just the run around it, so an expression interrupted by a comment or a
// processing instruction is still read whole. Given an attribute it is the attribute's value.
func stringFromNode(n *xmldom.Node) string {
	switch n.Kind {
	case xmldom.Text, xmldom.CDATA:
		var sb strings.Builder
		for c := n.Parent.FirstChild; c != nil; c = c.NextSibling {
			if c.Kind == xmldom.Text || c.Kind == xmldom.CDATA {
				sb.WriteString(c.Value)
			}
		}
		return sb.String()
	case xmldom.Attribute:
		return n.Value
	}
	return ""
}

// decodeBase64 is XMLUtils#decode, i.e. java.util.Base64.getMimeDecoder().decode.
//
// The MIME decoder is not the strict one and the difference is observable in real signatures:
// it SKIPS every character outside the base64 alphabet, which is what lets a ds:DigestValue or
// a base64 ds:Object wrapped across lines - or indented with tabs - decode at all.
//
// What it does NOT skip is a malformed final quantum, and that is the half worth porting
// carefully, because it is the difference between rejecting a truncated ds:DigestValue and
// silently accepting a prefix of it. The rules, transcribed from Base64.Decoder#decode0:
//
//   - '=' with no character of the current quantum yet read, or after exactly one, is an error
//     ("wrong 4-byte ending unit" / "last unit does not have enough valid bits");
//   - '=' after exactly two characters requires the IMMEDIATELY following byte to be another
//     '=' - not the next base64 character, the next byte - and is an error otherwise;
//   - '=' after three characters ends the input;
//   - after the padding, non-alphabet bytes are skipped but ANY further base64 character is an
//     error ("incorrect ending byte at N") - so trailing junk is tolerated and trailing data is
//     not.
func decodeBase64(s string) ([]byte, error) {
	var out []byte
	bits, shiftto := 0, 18
	resume := len(s)
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '=' {
			if shiftto == 18 {
				return nil, errors.New("xmldsig: base64 input has a wrong 4-byte ending unit")
			}
			resume = i + 1
			if shiftto == 6 {
				// A quantum of two characters needs "==", and the second '=' must be the very
				// next byte - not the next base64 character.
				if resume >= len(s) || s[resume] != '=' {
					return nil, errors.New("xmldsig: base64 input has a wrong 4-byte ending unit")
				}
				resume++
			}
			break
		}
		v := base64Value(c)
		if v < 0 {
			continue
		}
		bits |= v << shiftto
		shiftto -= 6
		if shiftto < 0 {
			out = append(out, byte(bits>>16), byte(bits>>8), byte(bits))
			bits, shiftto = 0, 18
		}
	}
	switch shiftto {
	case 6:
		out = append(out, byte(bits>>16))
	case 0:
		out = append(out, byte(bits>>16), byte(bits>>8))
	case 12:
		return nil, errors.New("xmldsig: the last base64 unit does not have enough valid bits")
	}
	for i := resume; i < len(s); i++ {
		if base64Value(s[i]) >= 0 {
			return nil, errors.New("xmldsig: base64 input has data after its padding")
		}
	}
	return out, nil
}

// base64Value decodes one character of the standard base64 alphabet, or -1.
func base64Value(c byte) int {
	switch {
	case c >= 'A' && c <= 'Z':
		return int(c - 'A')
	case c >= 'a' && c <= 'z':
		return int(c-'a') + 26
	case c >= '0' && c <= '9':
		return int(c-'0') + 52
	case c == '+':
		return 62
	case c == '/':
		return 63
	}
	return -1
}

// encodeBase64 is XMLUtils#encodeToString.
func encodeBase64(b []byte) string { return base64.StdEncoding.EncodeToString(b) }

// protectAgainstWrappingAttack ports XMLUtils#protectAgainstWrappingAttack: the whole tree from
// startElement down must contain at most one element whose registered ID attribute has the
// given value. Both same-document resolvers call it under secure validation, and it is the
// check that makes an XML signature wrapping attack fail closed rather than silently resolve
// to the attacker's copy.
func protectAgainstWrappingAttack(startElement *xmldom.Node, value string) bool {
	found := 0
	var walk func(n *xmldom.Node)
	walk = func(n *xmldom.Node) {
		if n.Kind == xmldom.Element {
			doc := n.Document()
			for _, a := range n.Attrs {
				if a.Value != value {
					continue
				}
				// Only a registered ID attribute counts, matching Attr#isId. IDAttrs is
				// xmldom's report of the attributes RegisterIDs declared for this element.
				if doc != nil && isIDAttr(doc, n, a) {
					found++
					if found > 1 {
						return
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(startElement)
	return found <= 1
}

func isIDAttr(doc, elem, attr *xmldom.Node) bool {
	for _, a := range doc.IDAttrs(elem) {
		if a == attr {
			return true
		}
	}
	return false
}
