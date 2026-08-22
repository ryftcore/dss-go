// Ported from
// dss-pades/src/main/java/eu/europa/esig/dss/pades/validation/RevocationInfoArchival.java
// (DSS 6.5.RC1).
//
// org.bouncycastle.asn1.esf.OtherRevVals, org.bouncycastle.asn1.ocsp.OCSPResponse and
// org.bouncycastle.asn1.x509.CertificateList have no DSS class of their own and no port
// elsewhere in this codebase (PORTING.md's "machinery BouncyCastle provides upstream" rule).
// Following the contract pdf_cms_ocsp_source.go's header already commits this type to
// ("RevocationInfoArchival ... a struct with OcspVals() [][]byte, each element being the DER
// encoding of one ASN.1 OCSPResponse ... found in the RevocationInfoArchival ASN.1 SEQUENCE's
// [1] member"), CrlVals()/OcspVals() return the DER encoding of each CertificateList/OCSPResponse
// directly - Java's CertificateList#getEncoded()/OCSPResponse#getEncoded() collapsed into the
// getter, rather than exposing the intermediate ASN1Sequence/ASN1Encodable BouncyCastle would.
// OtherRevVals has no caller anywhere in this port (upstream never reads it back either - only
// the encode-direction constructor accepts it), so it is kept as the parsed *asn1ber.Element for
// ToASN1Primitive to re-encode.
//
//	RevocationInfoArchival ::= SEQUENCE {
//	  crl [0] EXPLICIT SEQUENCE of CRLs, OPTIONAL
//	  ocsp [1] EXPLICIT SEQUENCE of OCSP Responses, OPTIONAL
//	  otherRevInfo [2] EXPLICIT SEQUENCE of OtherRevInfo, OPTIONAL
//	}
package pades

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/internal/asn1ber"
)

// RevocationInfoArchival ports the ASN1Object subclass eu.europa.esig.dss.pades.validation.
// RevocationInfoArchival.
type RevocationInfoArchival struct {
	// crlVals is the CRL values (tag [0]): the DER encoding of each CertificateList found in
	// the inner SEQUENCE.
	crlVals [][]byte

	// ocspVals is the OCSP values (tag [1]): the DER encoding of each OCSPResponse found in the
	// inner SEQUENCE.
	ocspVals [][]byte

	// otherRevVals is the other revocation values (tag [2]), the raw OtherRevVals SEQUENCE.
	otherRevVals *asn1ber.Element
}

// NewRevocationInfoArchival is the constructor taking the three optional components directly.
// Port of the public RevocationInfoArchival(CertificateList[], OCSPResponse[], OtherRevVals)
// constructor. crlVals and ocspVals are the DER encoding of each CertificateList/OCSPResponse.
func NewRevocationInfoArchival(crlVals, ocspVals [][]byte, otherRevVals *asn1ber.Element) *RevocationInfoArchival {
	return &RevocationInfoArchival{
		crlVals:      crlVals,
		ocspVals:     ocspVals,
		otherRevVals: otherRevVals,
	}
}

// RevocationInfoArchivalGetInstance gets the RevocationInfoArchival object from a parsed ASN.1
// SEQUENCE. Port of the static getInstance(Object); Go has no dual-type instanceof/getInstance
// bridge, so the caller (PAdESUtils.getRevocationInfoArchival, a forward dependency of a sibling
// chunk) is expected to have already turned the attribute value into a parsed *asn1ber.Element
// before calling this. obj == nil ports the "obj != null" guard, returning (nil, nil) as Java's
// getInstance(null) does.
func RevocationInfoArchivalGetInstance(obj *asn1ber.Element) (*RevocationInfoArchival, error) {
	if obj == nil {
		return nil, nil
	}
	return newRevocationInfoArchival(obj)
}

// newRevocationInfoArchival ports the private RevocationInfoArchival(ASN1Sequence) constructor.
func newRevocationInfoArchival(seq *asn1ber.Element) (*RevocationInfoArchival, error) {
	children := seq.Children()
	if len(children) > 3 {
		return nil, fmt.Errorf("Bad sequence size: %d", len(children))
	}

	archival := &RevocationInfoArchival{}
	for _, o := range children {
		// EXPLICIT tagging: the tagged object is constructed and its single child is the
		// underlying value (ASN1TaggedObject#getBaseObject for explicit tagging).
		if len(o.Children()) != 1 {
			return nil, fmt.Errorf("invalid tag: %d", o.TagNumber())
		}
		base := o.Children()[0]

		switch o.TagNumber() {
		case 0:
			for _, crl := range base.Children() {
				archival.crlVals = append(archival.crlVals, crl.DEREncoded())
			}
		case 1:
			for _, ocsp := range base.Children() {
				archival.ocspVals = append(archival.ocspVals, ocsp.DEREncoded())
			}
		case 2:
			archival.otherRevVals = base
		default:
			return nil, fmt.Errorf("invalid tag: %d", o.TagNumber())
		}
	}
	return archival, nil
}

// CrlVals gets the CRL values, the DER encoding of each CertificateList.
// Port of getCrlVals(): CertificateList[].
func (r *RevocationInfoArchival) CrlVals() [][]byte {
	if r.crlVals == nil {
		return [][]byte{}
	}
	return r.crlVals
}

// OcspVals gets the OCSP values, the DER encoding of each OCSPResponse.
// Port of getOcspVals(): OCSPResponse[].
func (r *RevocationInfoArchival) OcspVals() [][]byte {
	if r.ocspVals == nil {
		return [][]byte{}
	}
	return r.ocspVals
}

// OtherRevVals gets the other revocation values. Port of getOtherRevVals().
func (r *RevocationInfoArchival) OtherRevVals() *asn1ber.Element {
	return r.otherRevVals
}

// ToASN1Primitive encodes the RevocationInfoArchival to its DER SEQUENCE. Port of
// toASN1Primitive().
func (r *RevocationInfoArchival) ToASN1Primitive() []byte {
	var body []byte
	if r.crlVals != nil {
		body = append(body, explicitTag(0, derSequenceOfEncoded(r.crlVals))...)
	}
	if r.ocspVals != nil {
		body = append(body, explicitTag(1, derSequenceOfEncoded(r.ocspVals))...)
	}
	if r.otherRevVals != nil {
		body = append(body, explicitTag(2, r.otherRevVals.DEREncoded())...)
	}
	return asn1ber.WriteIdentifiedTLV([]byte{asn1ber.ClassUniversal | asn1ber.Constructed | asn1ber.TagSequence}, body)
}

// derSequenceOfEncoded DER-encodes a SEQUENCE OF the already DER-encoded members.
func derSequenceOfEncoded(members [][]byte) []byte {
	var body []byte
	for _, member := range members {
		body = append(body, member...)
	}
	return asn1ber.WriteIdentifiedTLV([]byte{asn1ber.ClassUniversal | asn1ber.Constructed | asn1ber.TagSequence}, body)
}

// explicitTag wraps the given DER-encoded content in an EXPLICIT context-specific constructed
// tag, mirroring new DERTaggedObject(true, tagNumber, content).
func explicitTag(tagNumber byte, content []byte) []byte {
	return asn1ber.WriteIdentifiedTLV([]byte{asn1ber.ClassContextSpecific | asn1ber.Constructed | tagNumber}, content)
}
