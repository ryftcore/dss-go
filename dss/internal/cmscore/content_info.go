// ContentInfo of RFC 5652 clause 3, replacing org.bouncycastle.asn1.cms.ContentInfo.
package cmscore

import (
	"encoding/asn1"
	"fmt"

	"github.com/ryftcore/dss-go/dss/internal/asn1ber"
)

// ContentInfo is the RFC 5652 outermost structure
//
//	ContentInfo ::= SEQUENCE {
//	    contentType ContentType,
//	    content     [0] EXPLICIT ANY DEFINED BY contentType }
//
// The content is kept as the parsed element it was received as, never re-encoded, so that a
// caller can hand out the original octets or ask asn1ber for a different encoding of the very
// same subtree.
type ContentInfo struct {
	// ContentType is the contentType OID.
	ContentType asn1.ObjectIdentifier

	// content is the element inside the [0] EXPLICIT wrapper, nil when the optional content
	// field is absent.
	content *asn1ber.Element
	// element is the parsed ContentInfo, nil for a ContentInfo built rather than parsed.
	element *asn1ber.Element
	// builtContent holds the DER encoding of the content of a built ContentInfo.
	builtContent []byte
}

// NewContentInfo builds a ContentInfo around the DER encoding of a content value.
func NewContentInfo(contentType asn1.ObjectIdentifier, content []byte) *ContentInfo {
	return &ContentInfo{ContentType: contentType, builtContent: content}
}

// ParseContentInfo decodes a ContentInfo, tolerating the BER encodings a streaming producer
// emits: an indefinite-length ContentInfo SEQUENCE and an indefinite-length [0] wrapper. Trailing
// bytes after the ContentInfo are rejected, matching BC's own CMSSignedData(byte[]) constructor
// (testdata/adversarial/craft-bad-trailing.p7s's bc-oracle.txt entry is parseOK=false for exactly
// this). A caller that instead wants BC's other, lenient entry point - reading one ASN.1 object
// from a stream and not caring what follows, which is what a document-level parse needs (see
// ParseContentInfoTolerateTrailingBytes) - does not go through this function.
func ParseContentInfo(input []byte) (*ContentInfo, error) {
	element, err := parseOne(input, "ContentInfo")
	if err != nil {
		return nil, err
	}
	return ContentInfoFromElement(element)
}

// ParseContentInfoTolerateTrailingBytes decodes a ContentInfo the same way ParseContentInfo
// does, except it does not require the input to be fully consumed.
//
// This is the entry point for parsing a whole document (UtilsParseToCMSBinaries's caller uses
// it, not ParseContentInfo) rather than an isolated byte slice: BC's own
// CMSSignedData(InputStream) constructor - built on ASN1InputStream#readObject(), which reads
// exactly one top-level ASN.1 object off the stream and stops - never looks at what follows it
// either, and upstream DSS relies on exactly that leniency. A real fixture depends on it:
// dss-cades/src/test/resources/validation/dss-1188/Test.bin.sig (DSS-1188) is a well-formed
// 4149-byte ContentInfo followed by 15859 bytes of unrelated trailing data, and upstream DSS
// parses and validates it without complaint - see cades/testdata/upstream/validation/dss-1188/.
// BC is inconsistent about this across its two constructors (CMSSignedData(byte[]) does reject
// trailing bytes, as ParseContentInfo's own doc comment notes), and DSS's own CMS-loading path
// goes through the lenient one; this function exists to let cms's document-parsing entry point
// mirror that specific choice without loosening ParseContentInfo for every other caller
// (SignedData/TimeStampResp/TSTInfo already go through the strict parseOne too, since those are
// always parsed from an already-isolated byte slice - an OCTET STRING's content, never "the rest
// of a file" - where trailing bytes remain a real error).
func ParseContentInfoTolerateTrailingBytes(input []byte) (*ContentInfo, error) {
	if len(input) == 0 {
		return nil, fmt.Errorf("cmscore: ContentInfo is empty")
	}
	element, _, err := asn1ber.Parse(input)
	if err != nil {
		return nil, fmt.Errorf("cmscore: cannot parse ContentInfo: %w", err)
	}
	return ContentInfoFromElement(element)
}

// ContentInfoFromElement decodes an already parsed ContentInfo.
func ContentInfoFromElement(element *asn1ber.Element) (*ContentInfo, error) {
	children, err := expectSizedSequence(element, "ContentInfo", 1, 2)
	if err != nil {
		return nil, err
	}
	contentType, err := expectOID(children[0], "ContentInfo.contentType")
	if err != nil {
		return nil, err
	}
	contentInfo := &ContentInfo{ContentType: contentType, element: element}
	if len(children) > 1 {
		if !children[1].IsContextSpecific(0) {
			return nil, fmt.Errorf("cmscore: ContentInfo.content is not [0] EXPLICIT")
		}
		content, err := explicitContent(children[1], "ContentInfo.content")
		if err != nil {
			return nil, err
		}
		contentInfo.content = content
	}
	return contentInfo, nil
}

// Content returns the element carried by the content field, nil when it is absent.
func (c *ContentInfo) Content() *asn1ber.Element { return c.content }

// ContentEncoded returns the original encoding of the content field's value, or - for a built
// ContentInfo - the DER it was built from. It is nil when the content field is absent.
func (c *ContentInfo) ContentEncoded() []byte {
	if c.content != nil {
		return c.content.Encoded()
	}
	return c.builtContent
}

// Element returns the parsed ContentInfo, nil when it was built rather than parsed.
func (c *ContentInfo) Element() *asn1ber.Element { return c.element }

// Encoded returns the ContentInfo's original encoding, and its DER encoding when it was built
// rather than parsed.
func (c *ContentInfo) Encoded() []byte {
	if c.element != nil {
		return c.element.Encoded()
	}
	return c.DER()
}

// IsDefiniteLength reports whether the ContentInfo SEQUENCE used the definite-length form,
// which is what ContentInfo#isDefiniteLength answers and what DSS turns into the choice
// between the DL and the BER encoding (see CMSUtils#getContentInfoEncoding). A ContentInfo
// that was built rather than parsed is definite, DER being all this package writes.
func (c *ContentInfo) IsDefiniteLength() bool {
	return c.element == nil || !c.element.IsIndefinite()
}

// DER returns the DER encoding of the ContentInfo.
func (c *ContentInfo) DER() []byte {
	body := asn1ber.EncodeOID(c.ContentType)
	if content := c.derContent(); content != nil {
		body = append(body, encodeContextTagged(0, content)...)
	}
	return asn1ber.WriteSequence(body)
}

// derContent returns the DER encoding of the content field's value, nil when absent.
func (c *ContentInfo) derContent() []byte {
	if c.content != nil {
		return c.content.DEREncoded()
	}
	return c.builtContent
}
