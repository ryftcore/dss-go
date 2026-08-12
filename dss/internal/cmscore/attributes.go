// Attribute and the signed/unsigned attribute sets of RFC 5652 clause 5.3, replacing
// org.bouncycastle.asn1.cms.Attribute and org.bouncycastle.asn1.cms.AttributeTable.
package cmscore

import (
	"encoding/asn1"

	"github.com/utain/esig/dss/internal/asn1ber"
)

// Attribute is
//
//	Attribute ::= SEQUENCE {
//	    attrType   OBJECT IDENTIFIER,
//	    attrValues SET OF AttributeValue }
//
// The values are kept as parsed elements: a CAdES attribute value is frequently digested or
// re-parsed by the layers above, which need its exact octets.
type Attribute struct {
	// Type is the attrType OID.
	Type asn1.ObjectIdentifier
	// Values holds the members of attrValues in the order they were received.
	Values []*asn1ber.Element

	// element is the parsed Attribute, nil when it was built.
	element *asn1ber.Element
	// builtValues holds the DER encodings of the values of a built Attribute.
	builtValues [][]byte
}

// NewAttribute builds an Attribute from the DER encodings of its values.
func NewAttribute(attrType asn1.ObjectIdentifier, values ...[]byte) *Attribute {
	return &Attribute{Type: attrType, builtValues: values}
}

// AttributeFromElement decodes an already parsed Attribute.
func AttributeFromElement(element *asn1ber.Element) (*Attribute, error) {
	children, err := expectSequence(element, "Attribute", 2)
	if err != nil {
		return nil, err
	}
	attrType, err := expectOID(children[0], "Attribute.attrType")
	if err != nil {
		return nil, err
	}
	values, err := expectSet(children[1], "Attribute.attrValues")
	if err != nil {
		return nil, err
	}
	return &Attribute{Type: attrType, Values: values, element: element}, nil
}

// Element returns the parsed Attribute, nil when it was built.
func (a *Attribute) Element() *asn1ber.Element { return a.element }

// Encoded returns the Attribute's original encoding, and its DER encoding when it was built
// rather than parsed.
func (a *Attribute) Encoded() []byte {
	if a.element != nil {
		return a.element.Encoded()
	}
	return a.DER()
}

// ValueEncodings returns the original encodings of the attribute values.
func (a *Attribute) ValueEncodings() [][]byte {
	if a.element == nil {
		return a.builtValues
	}
	encodings := make([][]byte, len(a.Values))
	for index, value := range a.Values {
		encodings[index] = value.Encoded()
	}
	return encodings
}

// DER returns the DER encoding of the Attribute, its values ordered as a DER SET OF.
//
// A parsed Attribute is re-encoded from the element it was parsed from rather than rebuilt
// from Type and Values. The two agree for a well-formed Attribute, but only the former keeps
// what a rebuild would drop: RFC 5652 gives the SEQUENCE exactly two components, yet a
// producer that emitted a third one still had it covered by the signature, and
// SignerInformation#getEncodedSignedAttributes - which DER-encodes the received SET rather
// than BouncyCastle's Attribute model - keeps it too. Re-encoding the signed attributes must
// never silently lose an octet the signature was computed over.
func (a *Attribute) DER() []byte {
	if a.element != nil {
		return a.element.DEREncoded()
	}
	body := asn1ber.EncodeOID(a.Type)
	body = append(body, derSetOf([]byte{setIdentifier}, a.builtValues)...)
	return asn1ber.WriteSequence(body)
}

// Attributes is a SignedAttributes or UnsignedAttributes set, i.e. BouncyCastle's
// AttributeTable, kept as a slice because the received order has to survive: re-encoding a
// signed attribute set in the input's order is what lets a caller reproduce the original
// bytes, and DER ordering is applied only where the encoding demands it.
type Attributes []*Attribute

// Get returns the first attribute of the given type, nil when there is none. Port of
// AttributeTable#get(ASN1ObjectIdentifier), which likewise answers the first of several.
func (a Attributes) Get(attrType asn1.ObjectIdentifier) *Attribute {
	for _, attribute := range a {
		if attribute.Type.Equal(attrType) {
			return attribute
		}
	}
	return nil
}

// GetAll returns every attribute of the given type. Port of AttributeTable#getAll.
func (a Attributes) GetAll(attrType asn1.ObjectIdentifier) Attributes {
	var matches Attributes
	for _, attribute := range a {
		if attribute.Type.Equal(attrType) {
			matches = append(matches, attribute)
		}
	}
	return matches
}

// DEREncodings returns the DER encoding of each attribute, in the received order.
func (a Attributes) DEREncodings() [][]byte {
	encodings := make([][]byte, len(a))
	for index, attribute := range a {
		encodings[index] = attribute.DER()
	}
	return encodings
}

// DERSetEncoded returns the attributes as a DER SET OF, i.e. with the universal SET tag and
// the members ordered by their encodings.
//
// This is the value RFC 5652 clause 5.4 makes the input of the signature computation - "the
// message digest is computed on the DER encoding of the SignedAttrs value, with the tag of
// SET OF rather than the [0] IMPLICIT one" - so it is what a CAdES signer signs and what a
// verifier hashes.
func (a Attributes) DERSetEncoded() []byte {
	return derSetOf([]byte{setIdentifier}, a.DEREncodings())
}

// DERImplicitTagged returns the attributes as a DER SET OF carrying the given context-specific
// tag instead of the universal SET tag, i.e. the signedAttrs [0] or unsignedAttrs [1] field of
// a SignerInfo.
func (a Attributes) DERImplicitTagged(tagNumber int) []byte {
	return derSetOf([]byte{contextIdentifier(tagNumber)}, a.DEREncodings())
}

// attributesFromElement decodes a signedAttrs/unsignedAttrs field, whose implicit tagging
// makes it a context-specific constructed element holding the Attribute components.
func attributesFromElement(element *asn1ber.Element, name string) (Attributes, error) {
	attributes := make(Attributes, 0, len(element.Children()))
	for _, child := range element.Children() {
		attribute, err := AttributeFromElement(child)
		if err != nil {
			return nil, wrapField(name, err)
		}
		attributes = append(attributes, attribute)
	}
	return attributes, nil
}
