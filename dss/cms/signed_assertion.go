// Ported from dss-cms/src/main/java/eu/europa/esig/dss/cms/asn1/SignedAssertion.java
// (DSS 6.5.RC1).
//
//	SignedAssertion ::= SEQUENCE {
//	    signedAssertionID SIGNED-ASSERTION.&id,
//	    signedAssertion SIGNED-ASSERTION.&Assertion OPTIONAL }
//
//	SIGNED-ASSERTION::= CLASS {
//	    &id OBJECT IDENTIFIER UNIQUE,
//	    &Assertion OPTIONAL }
//	    WITH SYNTAX {
//	    SIGNED-ASSERTION-ID &id
//	    [SIGNED-ASSERTION-TYPE &Assertion] }
package cms

import (
	"encoding/asn1"
	"errors"
	"fmt"

	"github.com/ryftcore/dss-go/dss/internal/asn1ber"
)

// signedAssertionOID is the SignedAssertion OID this port hard-codes, "0.4.0.19122.1.6", per
// the single arc BouncyCastle's own class hard-codes.
var signedAssertionOID = asn1.ObjectIdentifier{0, 4, 0, 19122, 1, 6}

// SignedAssertion is a signed assertion value. Port of the SignedAssertion class.
type SignedAssertion struct {
	// assertion is the SignedAssertion value.
	assertion string
}

// NewSignedAssertion creates the SignedAssertion from a string value.
// Port of SignedAssertion(String).
func NewSignedAssertion(assertion string) *SignedAssertion {
	return &SignedAssertion{assertion: assertion}
}

// ParseSignedAssertion decodes a SignedAssertion from its encoding.
// Port of SignedAssertion.getInstance(Object).
func ParseSignedAssertion(encoded []byte) (*SignedAssertion, error) {
	children, err := cmsSequenceOfSize(encoded, "SignedAssertion", 2, 2)
	if err != nil {
		return nil, err
	}
	if !children[1].IsUniversal(asn1ber.TagPrintableString) {
		return nil, errors.New("malformed SignedAssertion: the assertion is not a PrintableString")
	}
	return &SignedAssertion{assertion: string(children[1].Content())}, nil
}

// String returns the assertion value. Port of #toString.
func (s *SignedAssertion) String() string { return s.assertion }

// DER returns the DER encoding of the SignedAssertion. Port of #toASN1Primitive.
func (s *SignedAssertion) DER() []byte {
	body := asn1ber.EncodeOID(signedAssertionOID)
	body = append(body, asn1ber.WriteTLV(asn1ber.TagPrintableString, []byte(s.assertion))...)
	return asn1ber.WriteSequence(body)
}

// cmsSequenceOfSize parses a SEQUENCE and checks its size, reproducing the "Bad sequence size"
// guard BouncyCastle's ASN1 structures open with; shared by every asn1.go type of this package.
func cmsSequenceOfSize(encoded []byte, name string, min, max int) ([]*asn1ber.Element, error) {
	element, rest, err := asn1ber.Parse(encoded)
	if err != nil {
		return nil, err
	}
	if len(rest) != 0 {
		return nil, fmt.Errorf("extra data found after the %s", name)
	}
	if !element.IsUniversal(asn1ber.TagSequence) || !element.IsConstructed() {
		return nil, fmt.Errorf("malformed %s: not a SEQUENCE", name)
	}
	if len(element.Children()) < min || len(element.Children()) > max {
		return nil, fmt.Errorf("bad sequence size in %s: %d", name, len(element.Children()))
	}
	return element.Children(), nil
}
