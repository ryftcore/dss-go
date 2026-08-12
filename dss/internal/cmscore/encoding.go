// Encoding helpers shared by the structures of this package: the DER SET OF writer that
// replaces org.bouncycastle.asn1.DERSet, and a handful of parser guards.
package cmscore

import (
	"bytes"
	"encoding/asn1"
	"errors"
	"fmt"
	"math/big"
	"sort"

	"github.com/utain/esig/dss/internal/asn1ber"
)

// derSetOf writes a SET OF from its already DER-encoded members, ordering them as X.690
// clause 11.6 requires and as org.bouncycastle.asn1.ASN1Set#sort implements it: the member
// encodings are compared as unsigned octet strings, a member that is a prefix of another
// sorting first. The identifier octets are given, so the same writer serves the universal
// SET (SignedData.digestAlgorithms, SignedData.signerInfos, an attribute's values) and the
// implicitly tagged ones (SignedData.certificates [0], SignedData.crls [1], a SignerInfo's
// signedAttrs [0] and unsignedAttrs [1]).
func derSetOf(identifier []byte, members [][]byte) []byte {
	ordered := make([][]byte, len(members))
	copy(ordered, members)
	sort.SliceStable(ordered, func(a, b int) bool {
		return bytes.Compare(ordered[a], ordered[b]) < 0
	})
	var body []byte
	for _, member := range ordered {
		body = append(body, member...)
	}
	return asn1ber.WriteIdentifiedTLV(identifier, body)
}

// setIdentifier is the identifier octet of a universal SET.
const setIdentifier = asn1ber.TagSet | asn1ber.Constructed

// contextIdentifier returns the identifier octet of a constructed context-specific element.
func contextIdentifier(tagNumber int) byte {
	return asn1ber.ClassContextSpecific | asn1ber.Constructed | byte(tagNumber)
}

// encodeContextTagged wraps content in a constructed context-specific element.
func encodeContextTagged(tagNumber int, content []byte) []byte {
	return asn1ber.WriteTLV(contextIdentifier(tagNumber), content)
}

// encodeOctetString writes a primitive OCTET STRING.
func encodeOctetString(content []byte) []byte {
	return asn1ber.WriteTLV(asn1ber.TagOctetString, content)
}

// encodeInt writes an INTEGER holding a non-negative Go int.
func encodeInt(value int) []byte {
	return asn1ber.EncodeInteger(big.NewInt(int64(value)))
}

// expectSequence checks that the element is a SEQUENCE holding at least the given number of
// components, and returns them.
func expectSequence(element *asn1ber.Element, name string, minimum int) ([]*asn1ber.Element, error) {
	if element == nil || !element.IsUniversal(asn1ber.TagSequence) || !element.IsConstructed() {
		return nil, fmt.Errorf("cmscore: %s is not a SEQUENCE", name)
	}
	children := element.Children()
	if len(children) < minimum {
		return nil, fmt.Errorf("cmscore: %s holds %d components, at least %d expected", name, len(children), minimum)
	}
	return children, nil
}

// expectSizedSequence checks that the element is a SEQUENCE holding between minimum and
// maximum components, the bounds BouncyCastle's own getInstance methods enforce for the
// structures whose components this package re-encodes one by one.
func expectSizedSequence(element *asn1ber.Element, name string, minimum, maximum int) ([]*asn1ber.Element, error) {
	children, err := expectSequence(element, name, minimum)
	if err != nil {
		return nil, err
	}
	if len(children) > maximum {
		return nil, fmt.Errorf("cmscore: %s holds %d components, at most %d expected", name, len(children), maximum)
	}
	return children, nil
}

// expectSet checks that the element is a SET, definite or indefinite, and returns its members.
func expectSet(element *asn1ber.Element, name string) ([]*asn1ber.Element, error) {
	if element == nil || !element.IsUniversal(asn1ber.TagSet) || !element.IsConstructed() {
		return nil, fmt.Errorf("cmscore: %s is not a SET", name)
	}
	return element.Children(), nil
}

// expectOID decodes an OBJECT IDENTIFIER component.
func expectOID(element *asn1ber.Element, name string) (asn1.ObjectIdentifier, error) {
	oid, err := element.ObjectIdentifier()
	if err != nil {
		return nil, fmt.Errorf("cmscore: %s is not an OBJECT IDENTIFIER: %w", name, err)
	}
	return oid, nil
}

// expectInteger decodes an INTEGER component.
func expectInteger(element *asn1ber.Element, name string) (*big.Int, error) {
	if element == nil || !element.IsUniversal(asn1ber.TagInteger) {
		return nil, fmt.Errorf("cmscore: %s is not an INTEGER", name)
	}
	return element.Integer(), nil
}

// expectSmallInteger decodes an INTEGER component that has to fit a Go int, as every version
// number of RFC 5652 and RFC 3161 does.
func expectSmallInteger(element *asn1ber.Element, name string) (int, error) {
	value, err := expectInteger(element, name)
	if err != nil {
		return 0, err
	}
	if !value.IsInt64() {
		return 0, fmt.Errorf("cmscore: %s is out of range", name)
	}
	return int(value.Int64()), nil
}

// expectOctetString decodes an OCTET STRING component, joining the segments of a constructed
// one as BER allows.
func expectOctetString(element *asn1ber.Element, name string) ([]byte, error) {
	if element == nil || !element.IsUniversal(asn1ber.TagOctetString) {
		return nil, fmt.Errorf("cmscore: %s is not an OCTET STRING", name)
	}
	return element.Octets(), nil
}

// parseOne parses exactly one element out of the input, rejecting trailing data.
func parseOne(input []byte, name string) (*asn1ber.Element, error) {
	if len(input) == 0 {
		return nil, fmt.Errorf("cmscore: %s is empty", name)
	}
	element, rest, err := asn1ber.Parse(input)
	if err != nil {
		return nil, fmt.Errorf("cmscore: cannot parse %s: %w", name, err)
	}
	if len(rest) != 0 {
		return nil, fmt.Errorf("cmscore: %d trailing bytes after %s", len(rest), name)
	}
	return element, nil
}

// explicitContent unwraps an explicitly tagged component, i.e. returns the single element the
// context-specific constructed wrapper holds.
func explicitContent(element *asn1ber.Element, name string) (*asn1ber.Element, error) {
	if !element.IsConstructed() || len(element.Children()) != 1 {
		return nil, fmt.Errorf("cmscore: %s is not an EXPLICIT tagged element", name)
	}
	return element.Children()[0], nil
}

// wrapField prefixes an error with the field it was raised for.
func wrapField(name string, err error) error {
	return fmt.Errorf("%s: %w", name, err)
}

// errNoContent reports a ContentInfo whose optional content field is absent.
var errNoContent = errors.New("cmscore: the ContentInfo carries no content")
