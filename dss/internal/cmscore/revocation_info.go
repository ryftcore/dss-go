// RevocationInfoChoices and OtherRevocationInfoFormat of RFC 5652 clause 10.2.1, replacing
// org.bouncycastle.asn1.cms.OtherRevocationInfoFormat and the CRL/OCSP Stores of
// org.bouncycastle.cms.CMSSignedData.
package cmscore

import (
	"encoding/asn1"
	"errors"

	"github.com/ryftcore/dss-go/dss/internal/asn1ber"
)

// OtherRevocationInfoFormat is
//
//	OtherRevocationInfoFormat ::= SEQUENCE {
//	    otherRevInfoFormat OBJECT IDENTIFIER,
//	    otherRevInfo       ANY DEFINED BY otherRevInfoFormat }
//
// CAdES uses it to embed OCSP responses in SignedData.crls: an OCSPResponse under
// id-ri-ocsp-response, and - as several producers do and DSS accepts - a bare
// BasicOCSPResponse under id-pkix-ocsp-basic.
type OtherRevocationInfoFormat struct {
	// Format is the otherRevInfoFormat OID.
	Format asn1.ObjectIdentifier
	// Info is the encoding of the otherRevInfo value.
	Info []byte
}

// NewOtherRevocationInfoFormat builds an OtherRevocationInfoFormat from the DER encoding of
// its value.
func NewOtherRevocationInfoFormat(format asn1.ObjectIdentifier, info []byte) *OtherRevocationInfoFormat {
	return &OtherRevocationInfoFormat{Format: format, Info: info}
}

// DER returns the DER encoding of the OtherRevocationInfoFormat, without the [1] IMPLICIT tag
// the RevocationInfoChoice adds. Info is written as it stands, so a parsed structure keeps
// whatever encoding its value arrived in; RevocationInfoChoice.DER normalises instead.
func (o *OtherRevocationInfoFormat) DER() []byte {
	body := asn1ber.EncodeOID(o.Format)
	body = append(body, o.Info...)
	return asn1ber.WriteSequence(body)
}

// RevocationInfoChoice is
//
//	RevocationInfoChoice ::= CHOICE {
//	    crl   CertificateList,
//	    other [1] IMPLICIT OtherRevocationInfoFormat }
//
// Exactly one field is set.
type RevocationInfoChoice struct {
	// CRL is the encoding of the CertificateList for the crl alternative, nil otherwise.
	CRL []byte
	// Other is set for the [1] IMPLICIT alternative, nil otherwise.
	Other *OtherRevocationInfoFormat

	// element is the parsed member, nil when it was built.
	element *asn1ber.Element
}

// NewCRLRevocationInfoChoice builds the crl alternative from a CertificateList's DER.
func NewCRLRevocationInfoChoice(crl []byte) RevocationInfoChoice {
	return RevocationInfoChoice{CRL: crl}
}

// NewOtherRevocationInfoChoice builds the [1] IMPLICIT alternative.
func NewOtherRevocationInfoChoice(other *OtherRevocationInfoFormat) RevocationInfoChoice {
	return RevocationInfoChoice{Other: other}
}

// Encoded returns the member's original encoding, tag included.
func (r RevocationInfoChoice) Encoded() []byte {
	if r.element != nil {
		return r.element.Encoded()
	}
	return r.DER()
}

// Element returns the parsed member, nil when it was built.
func (r RevocationInfoChoice) Element() *asn1ber.Element { return r.element }

// DER returns the DER encoding of the member, tag included.
func (r RevocationInfoChoice) DER() []byte {
	if r.element != nil {
		return r.element.DEREncoded()
	}
	if r.Other != nil {
		// The implicit [1] tag replaces the SEQUENCE tag of OtherRevocationInfoFormat.
		body := asn1ber.EncodeOID(r.Other.Format)
		body = append(body, r.Other.Info...)
		return encodeContextTagged(1, body)
	}
	return r.CRL
}

// RevocationInfoChoices is the SignedData.crls field,
//
//	RevocationInfoChoices ::= SET OF RevocationInfoChoice
//
// carried by a [1] IMPLICIT tag. The field being optional, an absent one is a nil
// *RevocationInfoChoices rather than an empty set.
type RevocationInfoChoices struct {
	// Choices holds the members in the order they were received.
	Choices []RevocationInfoChoice

	// element is the parsed [1] IMPLICIT field, nil when it was built.
	element *asn1ber.Element
}

// revocationInfoChoicesFromElement decodes the [1] IMPLICIT crls field.
func revocationInfoChoicesFromElement(element *asn1ber.Element) (*RevocationInfoChoices, error) {
	choices := &RevocationInfoChoices{element: element}
	for _, child := range element.Children() {
		choice := RevocationInfoChoice{element: child}
		switch {
		case child.IsContextSpecific(1):
			if !child.IsConstructed() {
				return nil, errors.New("cmscore: RevocationInfoChoice.other is not constructed")
			}
			// The [1] tag implicitly replaced the SEQUENCE tag, so the children are the
			// components of OtherRevocationInfoFormat.
			if len(child.Children()) < 2 {
				return nil, errors.New("cmscore: OtherRevocationInfoFormat holds fewer than 2 components")
			}
			format, err := expectOID(child.Children()[0], "OtherRevocationInfoFormat.otherRevInfoFormat")
			if err != nil {
				return nil, err
			}
			choice.Other = &OtherRevocationInfoFormat{
				Format: format,
				Info:   child.Children()[1].Encoded(),
			}
		case child.Class() == asn1ber.ClassUniversal:
			choice.CRL = child.Encoded()
		default:
			return nil, errors.New("cmscore: unknown RevocationInfoChoice alternative")
		}
		choices.Choices = append(choices.Choices, choice)
	}
	return choices, nil
}

// Element returns the parsed crls field, nil when it was built or absent.
func (r *RevocationInfoChoices) Element() *asn1ber.Element {
	if r == nil {
		return nil
	}
	return r.element
}

// Encoded returns the crls field exactly as received, [1] IMPLICIT tag included, and its DER
// encoding when it was built rather than parsed. An absent field - a nil
// *RevocationInfoChoices - has no encoding.
func (r *RevocationInfoChoices) Encoded() []byte {
	if r == nil {
		return nil
	}
	if r.element != nil {
		return r.element.Encoded()
	}
	return r.DER()
}

// CRLs returns the encoding of every CertificateList in the set, in the received order.
// OtherRevocationInfoFormat members are excluded, matching CMSSignedData#getCRLs.
func (r *RevocationInfoChoices) CRLs() [][]byte {
	if r == nil {
		return nil
	}
	var crls [][]byte
	for _, choice := range r.Choices {
		if choice.Other == nil {
			crls = append(crls, choice.CRL)
		}
	}
	return crls
}

// OtherRevocationInfo returns the encoding of the otherRevInfo value of every member carrying
// the given format, in the received order. Port of
// CMSSignedData#getOtherRevocationInfo(ASN1ObjectIdentifier).
func (r *RevocationInfoChoices) OtherRevocationInfo(format asn1.ObjectIdentifier) [][]byte {
	if r == nil {
		return nil
	}
	var values [][]byte
	for _, choice := range r.Choices {
		if choice.Other != nil && choice.Other.Format.Equal(format) {
			values = append(values, choice.Other.Info)
		}
	}
	return values
}

// OCSPResponses returns the OCSPResponse encodings carried under id-ri-ocsp-response, i.e.
// what DSS reads as the OCSP response store.
func (r *RevocationInfoChoices) OCSPResponses() [][]byte {
	return r.OtherRevocationInfo(OIDRIOCSPResponse)
}

// OCSPBasicResponses returns the BasicOCSPResponse encodings carried under id-pkix-ocsp-basic,
// i.e. what DSS reads as the OCSP basic store.
func (r *RevocationInfoChoices) OCSPBasicResponses() [][]byte {
	return r.OtherRevocationInfo(OIDPKIXOCSPBasic)
}

// HasOtherFormat reports whether any member uses the [1] IMPLICIT OtherRevocationInfoFormat
// alternative, which forces SignedData to version 5.
func (r *RevocationInfoChoices) HasOtherFormat() bool {
	if r == nil {
		return false
	}
	for _, choice := range r.Choices {
		if choice.Other != nil {
			return true
		}
	}
	return false
}

// DER returns the DER encoding of the crls field, [1] IMPLICIT tag included and members
// ordered as a DER SET OF. An absent field - a nil *RevocationInfoChoices - has none.
func (r *RevocationInfoChoices) DER() []byte {
	if r == nil {
		return nil
	}
	members := make([][]byte, len(r.Choices))
	for index, choice := range r.Choices {
		members[index] = choice.DER()
	}
	return derSetOf([]byte{contextIdentifier(1)}, members)
}
