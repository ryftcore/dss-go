// EncapsulatedContentInfo of RFC 5652 clause 5.2, replacing the eContent handling of
// org.bouncycastle.asn1.cms.ContentInfo and org.bouncycastle.cms.CMSSignedData#getSignedContent.
package cmscore

import (
	"encoding/asn1"
	"fmt"

	"github.com/utain/esig/dss/internal/asn1ber"
)

// EncapsulatedContentInfo is
//
//	EncapsulatedContentInfo ::= SEQUENCE {
//	    eContentType ContentType,
//	    eContent     [0] EXPLICIT OCTET STRING OPTIONAL }
//
// An absent eContent is a detached signature, the content being supplied out of band.
//
// In BER the eContent may be a constructed OCTET STRING split into segments, possibly of
// indefinite length - that is what "openssl cms -sign -stream" and BouncyCastle's
// CMSSignedDataStreamGenerator write. Content joins the segments; ContentElement hands out the
// OCTET STRING as received, since an archive time-stamp digests the encapContentInfo in its
// original encoding.
type EncapsulatedContentInfo struct {
	// EContentType is the eContentType OID, id-data for an ordinary CAdES signature and
	// id-ct-TSTInfo inside a time-stamp token.
	EContentType asn1.ObjectIdentifier

	// content is the eContent OCTET STRING element, nil for a detached signature.
	content *asn1ber.Element
	// element is the parsed EncapsulatedContentInfo, nil when it was built.
	element *asn1ber.Element
	// builtContent holds the content octets of a built EncapsulatedContentInfo.
	builtContent []byte
	// detached records that a built EncapsulatedContentInfo carries no eContent.
	detached bool
}

// NewEncapsulatedContentInfo builds an EncapsulatedContentInfo carrying the given content
// octets. A nil content builds the detached form.
func NewEncapsulatedContentInfo(eContentType asn1.ObjectIdentifier, content []byte) *EncapsulatedContentInfo {
	return &EncapsulatedContentInfo{EContentType: eContentType, builtContent: content, detached: content == nil}
}

// EncapsulatedContentInfoFromElement decodes an already parsed EncapsulatedContentInfo.
func EncapsulatedContentInfoFromElement(element *asn1ber.Element) (*EncapsulatedContentInfo, error) {
	children, err := expectSizedSequence(element, "EncapsulatedContentInfo", 1, 2)
	if err != nil {
		return nil, err
	}
	eContentType, err := expectOID(children[0], "EncapsulatedContentInfo.eContentType")
	if err != nil {
		return nil, err
	}
	info := &EncapsulatedContentInfo{EContentType: eContentType, element: element, detached: true}
	if len(children) > 1 {
		if !children[1].IsContextSpecific(0) {
			return nil, fmt.Errorf("cmscore: EncapsulatedContentInfo.eContent is not [0] EXPLICIT")
		}
		content, err := explicitContent(children[1], "EncapsulatedContentInfo.eContent")
		if err != nil {
			return nil, err
		}
		if !content.IsUniversal(asn1ber.TagOctetString) {
			return nil, fmt.Errorf("cmscore: EncapsulatedContentInfo.eContent is not an OCTET STRING")
		}
		info.content = content
		info.detached = false
	}
	return info, nil
}

// IsDetached reports whether the eContent field is absent.
func (e *EncapsulatedContentInfo) IsDetached() bool { return e.detached }

// Content returns the signed content, joining the segments of a constructed OCTET STRING. It
// is nil for a detached signature.
func (e *EncapsulatedContentInfo) Content() []byte {
	if e.content != nil {
		return e.content.Octets()
	}
	return e.builtContent
}

// ContentElement returns the eContent OCTET STRING as received, nil for a detached signature
// or a built EncapsulatedContentInfo.
func (e *EncapsulatedContentInfo) ContentElement() *asn1ber.Element { return e.content }

// Element returns the parsed EncapsulatedContentInfo, nil when it was built.
func (e *EncapsulatedContentInfo) Element() *asn1ber.Element { return e.element }

// Encoded returns the EncapsulatedContentInfo's original encoding, and its DER encoding when
// it was built rather than parsed.
func (e *EncapsulatedContentInfo) Encoded() []byte {
	if e.element != nil {
		return e.element.Encoded()
	}
	return e.DER()
}

// DER returns the DER encoding of the EncapsulatedContentInfo, in which a segmented eContent
// collapses into a single primitive OCTET STRING.
func (e *EncapsulatedContentInfo) DER() []byte {
	body := asn1ber.EncodeOID(e.EContentType)
	switch {
	case e.content != nil:
		body = append(body, encodeContextTagged(0, e.content.DEREncoded())...)
	case !e.detached:
		body = append(body, encodeContextTagged(0, encodeOctetString(e.builtContent))...)
	}
	return asn1ber.WriteSequence(body)
}
