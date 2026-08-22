// Ported from org.apache.xml.security.signature.SignedInfo (Apache Santuario xmlsec 3.0.6).
package xmldsig

import (
	"bytes"
	"errors"

	"github.com/ryftcore/dss-go/dss/internal/xmlc14n"
	"github.com/ryftcore/dss-go/dss/internal/xmldom"
)

// SignedInfo is a ds:SignedInfo. Port of SignedInfo, which extends Manifest in Java too: the
// reference verification is entirely inherited, and only the canonicalization and the
// signature method are its own.
type SignedInfo struct {
	*Manifest

	c14nBytes []byte
}

// ErrNotASignedInfo reports an element that is not a ds:SignedInfo.
var ErrNotASignedInfo = errors.New("xmldsig: not a ds:SignedInfo element")

// NewSignedInfo wraps a ds:SignedInfo element. Port of SignedInfo(Element, String, boolean).
func NewSignedInfo(element *xmldom.Node, opts *ManifestOptions) (*SignedInfo, error) {
	if element == nil || element.Kind != xmldom.Element ||
		element.Name.Space != NamespaceDSig || element.Name.Local != "SignedInfo" {
		return nil, ErrNotASignedInfo
	}
	m, err := NewManifest(element, opts)
	if err != nil {
		return nil, err
	}
	return &SignedInfo{Manifest: m}, nil
}

// CanonicalizationMethodURI returns the ds:CanonicalizationMethod Algorithm.
// Port of getCanonicalizationMethodURI().
func (si *SignedInfo) CanonicalizationMethodURI() string {
	cm := selectDSNode(si.element, "CanonicalizationMethod", 0)
	if cm == nil {
		return ""
	}
	return cm.AttrValue("", "Algorithm")
}

// SignatureMethodURI returns the ds:SignatureMethod Algorithm.
// Port of getSignatureMethodURI().
func (si *SignedInfo) SignatureMethodURI() string {
	sm := selectDSNode(si.element, "SignatureMethod", 0)
	if sm == nil {
		return ""
	}
	return sm.AttrValue("", "Algorithm")
}

// inclusiveNamespaces returns the exclusive-c14n PrefixList of the ds:CanonicalizationMethod,
// or nil. Port of getInclusiveNamespaces(), whose null/non-null distinction chooses between
// the two canonicalizeSubtree overloads.
func (si *SignedInfo) inclusiveNamespaces() []string {
	cm := selectDSNode(si.element, "CanonicalizationMethod", 0)
	if cm == nil {
		return nil
	}
	alg := cm.AttrValue("", "Algorithm")
	if alg != TransformC14NExcl && alg != TransformC14NExclWithComments {
		return nil
	}
	incl := selectNodes(cm, NamespaceExcC14N, "InclusiveNamespaces")
	if len(incl) == 0 {
		return nil
	}
	return xmlc14n.ParsePrefixList(incl[0].AttrValue("", "PrefixList"))
}

// CanonicalizedOctets returns the octets the ds:SignatureValue is computed over.
// Port of getCanonicalizedOctetStream().
//
// Note what it is NOT: it is a plain SUBTREE canonicalization of the ds:SignedInfo element,
// with no reference processing, no node filtering and no comment exclusion - the comment
// behaviour is entirely the algorithm's, so a ds:CanonicalizationMethod naming a #WithComments
// variant really does sign the comments inside ds:SignedInfo. The result is cached, as
// upstream caches c14nizedBytes, because both the verification and DSS's
// getDataToBeSignedRepresentation ask for it.
func (si *SignedInfo) CanonicalizedOctets() ([]byte, error) {
	if si.c14nBytes != nil {
		return si.c14nBytes, nil
	}
	alg, err := xmlc14n.Resolve(xmlc14n.Algorithm(si.CanonicalizationMethodURI()))
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	err = xmlc14n.Canonicalize(alg, xmlc14n.Input{
		Node:              si.element,
		InclusivePrefixes: si.inclusiveNamespaces(),
	}, &buf)
	if err != nil {
		return nil, err
	}
	si.c14nBytes = buf.Bytes()
	return si.c14nBytes, nil
}

// Verify re-digests every reference. Port of SignedInfo#verify(boolean), which is
// Manifest#verifyReferences under another name.
func (si *SignedInfo) Verify(followManifests bool) (bool, error) {
	return si.VerifyReferences(followManifests)
}
