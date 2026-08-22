// Ported from dss-cms/src/main/java/eu/europa/esig/dss/cms/asn1/SignerAttributeV2.java
// (DSS 6.5.RC1).
//
// "Basic support of ETSI EN 319 122-1 V1.1.1 chapter 5.2.6.1. Based on
// org.bouncycastle.asn1.esf.SignerAttribute. Note : signedAssertions are not supported. Quote
// ETSI : The definition of specific signedAssertions is outside of the scope of the present
// document" (Java doc, kept for context, notwithstanding that this port's own
// NewSignerAttributeV2FromSignedAssertions constructor - like Java's own
// SignerAttributeV2(SignedAssertions) - does build the signedAssertions alternative;
// the "not supported" note is about signedAssertions found while *parsing*, which upstream and
// this port both merely log and skip inside CertifiedAttributesV2, not SignerAttributeV2
// itself).
//
//	SignerAttributeV2 ::= SEQUENCE {
//	    claimedAttributes     [0] ClaimedAttributes OPTIONAL,
//	    certifiedAttributesV2 [1] CertifiedAttributesV2 OPTIONAL,
//	    signedAssertions      [2] SignedAssertions OPTIONAL }
//
// All three tags are EXPLICIT (BouncyCastle's `new DERTaggedObject(tagNo, obj)` two-argument
// form defaults to explicit tagging - verified against a real BouncyCastle 1.84 run, see
// testdata/gen), unlike CMSAlgorithmProtection's implicit [1]/[2] in
// cms_signed_attribute_table_generator.go: the two are unrelated BouncyCastle call sites using
// different DERTaggedObject constructors, and neither should be assumed from the other.
package cms

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/internal/asn1ber"
)

// SignerAttributeV2 is the signer-attributes-v2 value: at most one of ClaimedAttributes,
// CertifiedAttributes and SignedAssertions is ever built by this port's three constructors
// (mirroring Java's three single-argument constructors), though ParseSignerAttributeV2 accepts
// any combination the SEQUENCE's three independent OPTIONAL fields allow. Port of the
// SignerAttributeV2 class.
type SignerAttributeV2 struct {
	// claimedAttributes holds the DER encoding of each Attribute of the claimedAttributes
	// field, hasClaimedAttributes distinguishing "field absent" from "field present but empty".
	claimedAttributes    [][]byte
	hasClaimedAttributes bool

	// certifiedAttributes is set when the certifiedAttributesV2 field is present.
	certifiedAttributes *CertifiedAttributesV2

	// signedAssertions is set when the signedAssertions field is present.
	signedAssertions *SignedAssertions
}

// NewSignerAttributeV2FromClaimedAttributes creates a SignerAttributeV2 from an array of
// claimedAttributes, each the DER encoding of an X.509 Attribute.
// Port of SignerAttributeV2(Attribute[]).
func NewSignerAttributeV2FromClaimedAttributes(claimedAttributes [][]byte) *SignerAttributeV2 {
	return &SignerAttributeV2{claimedAttributes: claimedAttributes, hasClaimedAttributes: true}
}

// NewSignerAttributeV2FromCertifiedAttributes creates a SignerAttributeV2 from
// certifiedAttributes. Port of SignerAttributeV2(CertifiedAttributesV2).
func NewSignerAttributeV2FromCertifiedAttributes(certifiedAttributes *CertifiedAttributesV2) *SignerAttributeV2 {
	return &SignerAttributeV2{certifiedAttributes: certifiedAttributes}
}

// NewSignerAttributeV2FromSignedAssertions creates a SignerAttributeV2 from signedAssertions.
// Port of SignerAttributeV2(SignedAssertions).
func NewSignerAttributeV2FromSignedAssertions(signedAssertions *SignedAssertions) *SignerAttributeV2 {
	return &SignerAttributeV2{signedAssertions: signedAssertions}
}

// ParseSignerAttributeV2 decodes a SignerAttributeV2 from its encoding.
// Port of SignerAttributeV2.getInstance(Object).
func ParseSignerAttributeV2(encoded []byte) (*SignerAttributeV2, error) {
	element, rest, err := asn1ber.Parse(encoded)
	if err != nil {
		return nil, err
	}
	if len(rest) != 0 {
		return nil, fmt.Errorf("extra data found after the SignerAttributeV2")
	}
	if !element.IsUniversal(asn1ber.TagSequence) || !element.IsConstructed() {
		return nil, fmt.Errorf("malformed SignerAttributeV2: not a SEQUENCE")
	}
	attributes := &SignerAttributeV2{}
	for _, child := range element.Children() {
		if child.Class() != asn1ber.ClassContextSpecific {
			return nil, fmt.Errorf("illegal tag in SignerAttributeV2")
		}
		base, err := cmsExplicitBase(child, "SignerAttributeV2")
		if err != nil {
			return nil, err
		}
		switch child.TagNumber() {
		case 0:
			if !base.IsUniversal(asn1ber.TagSequence) || !base.IsConstructed() {
				return nil, fmt.Errorf("malformed SignerAttributeV2.claimedAttributes: not a SEQUENCE")
			}
			for _, attribute := range base.Children() {
				attributes.claimedAttributes = append(attributes.claimedAttributes, attribute.Encoded())
			}
			attributes.hasClaimedAttributes = true
		case 1:
			certifiedAttributes, err := ParseCertifiedAttributesV2(base.Encoded())
			if err != nil {
				return nil, err
			}
			attributes.certifiedAttributes = certifiedAttributes
		case 2:
			// Upstream logs "SAML assertion detected".
			signedAssertions, err := ParseSignedAssertions(base.Encoded())
			if err != nil {
				return nil, err
			}
			attributes.signedAssertions = signedAssertions
		default:
			return nil, fmt.Errorf("illegal tag: %d", child.TagNumber())
		}
	}
	return attributes, nil
}

// ClaimedAttributes returns the DER encoding of each Attribute of the claimedAttributes field,
// nil when the field is absent.
func (s *SignerAttributeV2) ClaimedAttributes() [][]byte { return s.claimedAttributes }

// CertifiedAttributes returns the certifiedAttributesV2 field, nil when absent.
func (s *SignerAttributeV2) CertifiedAttributes() *CertifiedAttributesV2 {
	return s.certifiedAttributes
}

// SignedAssertions returns the signedAssertions field, nil when absent.
func (s *SignerAttributeV2) SignedAssertions() *SignedAssertions { return s.signedAssertions }

// DER returns the DER encoding of the SignerAttributeV2. Port of #toASN1Primitive.
func (s *SignerAttributeV2) DER() []byte {
	var body []byte
	if s.hasClaimedAttributes {
		var claimedBody []byte
		for _, attribute := range s.claimedAttributes {
			claimedBody = append(claimedBody, attribute...)
		}
		body = append(body, asn1ber.WriteTLV(asn1ber.ClassContextSpecific|asn1ber.Constructed|0, asn1ber.WriteSequence(claimedBody))...) //nolint:staticcheck // mirrors upstream SignerAttributeV2#toASN1Primitive: `new DERTaggedObject(0, new DERSequence(...))` - the `0` is the ASN.1 context tag number, kept explicit next to its [1] and [2] siblings.
	}
	if s.certifiedAttributes != nil {
		body = append(body, asn1ber.WriteTLV(asn1ber.ClassContextSpecific|asn1ber.Constructed|1, s.certifiedAttributes.DER())...)
	}
	if s.signedAssertions != nil {
		body = append(body, asn1ber.WriteTLV(asn1ber.ClassContextSpecific|asn1ber.Constructed|2, s.signedAssertions.DER())...)
	}
	return asn1ber.WriteSequence(body)
}
