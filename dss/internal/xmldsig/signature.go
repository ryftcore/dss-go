// Ported from org.apache.xml.security.signature.XMLSignature (Apache Santuario xmlsec 3.0.6),
// as eu.europa.esig.dss.xades.validation.XAdESSignature and XAdESSignatureIntegrityValidator
// drive it.
package xmldsig

import (
	"crypto"
	"errors"
	"fmt"

	"github.com/utain/esig/dss/internal/xmldom"
)

// XMLSignature is a ds:Signature. Port of XMLSignature(Element, String, boolean), the
// constructor XAdESSignature.getSantuarioSignature calls with secure validation off.
type XMLSignature struct {
	element   *xmldom.Node
	baseURI   string
	signed    *SignedInfo
	followMan bool
}

// ErrNotASignature reports an element that is not a ds:Signature.
var ErrNotASignature = errors.New("xmldsig: not a ds:Signature element")

// NewXMLSignature wraps a ds:Signature element.
//
// The document's ID attributes must already be registered - xmldom's RegisterIDs, which is
// DSS's XAdESDOMDocument.recursiveIdBrowse - or every "#id" reference will fail to resolve.
// Upstream has the same precondition and satisfies it in getSantuarioSignature, one line
// before constructing the XMLSignature.
func NewXMLSignature(element *xmldom.Node, opts *ManifestOptions) (*XMLSignature, error) {
	if element == nil || element.Kind != xmldom.Element ||
		element.Name.Space != NamespaceDSig || element.Name.Local != "Signature" {
		return nil, ErrNotASignature
	}
	siEl := selectDSNode(element, "SignedInfo", 0)
	if siEl == nil {
		return nil, errors.New("xmldsig: the ds:Signature has no ds:SignedInfo")
	}
	si, err := NewSignedInfo(siEl, opts)
	if err != nil {
		return nil, err
	}
	s := &XMLSignature{element: element, signed: si}
	if opts != nil {
		s.baseURI = opts.BaseURI
	}
	return s, nil
}

// Element returns the ds:Signature element.
func (s *XMLSignature) Element() *xmldom.Node { return s.element }

// SignedInfo returns the ds:SignedInfo. Port of getSignedInfo().
func (s *XMLSignature) SignedInfo() *SignedInfo { return s.signed }

// AddResourceResolver appends a resolver to the signature's ds:SignedInfo, which is the
// manifest whose references are being validated. Port of
// XMLSignature#addResourceResolver, which delegates to signedInfo.addResourceResolver.
func (s *XMLSignature) AddResourceResolver(r URIResolver) { s.signed.AddResourceResolver(r) }

// SetFollowNestedManifests ports setFollowNestedManifests: whether CheckSignatureValue also
// validates the references of a ds:Manifest that a reference points at. Default false, as
// upstream.
func (s *XMLSignature) SetFollowNestedManifests(b bool) { s.followMan = b }

// SignatureValue returns the decoded ds:SignatureValue. Port of getSignatureValue().
func (s *XMLSignature) SignatureValue() ([]byte, error) {
	sv := selectDSNode(s.element, "SignatureValue", 0)
	if sv == nil {
		return nil, errors.New("xmldsig: the ds:Signature has no ds:SignatureValue")
	}
	return decodeBase64(sv.TextContent())
}

// KeyInfoElement returns the ds:KeyInfo child, or nil.
func (s *XMLSignature) KeyInfoElement() *xmldom.Node { return selectDSNode(s.element, "KeyInfo", 0) }

// CheckSignatureValue verifies the whole signature against pub: first the cryptographic
// signature over the canonicalized ds:SignedInfo, then every ds:Reference.
// Port of checkSignatureValue(Key).
//
// The order is Santuario's and is not an implementation detail: the references are re-digested
// ONLY if the signature over ds:SignedInfo verified, so a forged signature never causes the
// validator to dereference the URIs it names.
func (s *XMLSignature) CheckSignatureValue(pub crypto.PublicKey) (bool, error) {
	if pub == nil {
		return false, errors.New("xmldsig: no key to check the signature with")
	}
	c14n, err := s.signed.CanonicalizedOctets()
	if err != nil {
		return false, err
	}
	sigValue, err := s.SignatureValue()
	if err != nil {
		return false, err
	}
	ok, err := verifySignature(s.signed.SignatureMethodURI(), pub, c14n, sigValue)
	if err != nil {
		return false, fmt.Errorf("xmldsig: %s: %w", s.signed.SignatureMethodURI(), err)
	}
	if !ok {
		return false, nil
	}
	return s.signed.Verify(s.followMan)
}

// CheckSignatureValueOnly verifies the cryptographic signature over the canonicalized
// ds:SignedInfo and nothing else.
//
// It has no Santuario counterpart because Santuario never separates the two halves, but DSS
// does: XAdESSignature.checkSignatureIntegrity reports "signature intact" and "reference data
// intact" as two independent flags of SignatureCryptographicVerification, and it can only fill
// them in separately if the two checks can be run separately. Upstream gets there by calling
// checkSignatureValue (which does both) and then reading the reference validations back out of
// its own ReferenceValidation list; splitting the call is the same information without
// re-deriving it.
func (s *XMLSignature) CheckSignatureValueOnly(pub crypto.PublicKey) (bool, error) {
	if pub == nil {
		return false, errors.New("xmldsig: no key to check the signature with")
	}
	c14n, err := s.signed.CanonicalizedOctets()
	if err != nil {
		return false, err
	}
	sigValue, err := s.SignatureValue()
	if err != nil {
		return false, err
	}
	ok, err := verifySignature(s.signed.SignatureMethodURI(), pub, c14n, sigValue)
	if err != nil {
		return false, fmt.Errorf("xmldsig: %s: %w", s.signed.SignatureMethodURI(), err)
	}
	return ok, nil
}
