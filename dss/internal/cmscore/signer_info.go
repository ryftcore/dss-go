// SignerInfo and SignerIdentifier of RFC 5652 clause 5.3, replacing
// org.bouncycastle.asn1.cms.SignerInfo, org.bouncycastle.asn1.cms.SignerIdentifier,
// org.bouncycastle.asn1.cms.IssuerAndSerialNumber and the read side of
// org.bouncycastle.cms.SignerInformation.
package cmscore

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/utain/esig/dss/internal/asn1ber"
)

// IssuerAndSerialNumber is
//
//	IssuerAndSerialNumber ::= SEQUENCE {
//	    issuer       Name,
//	    serialNumber CertificateSerialNumber }
//
// The issuer is kept as its encoding: a distinguished name has to be compared and digested
// byte for byte, and re-encoding one loses the string types the issuing CA chose.
type IssuerAndSerialNumber struct {
	// Issuer is the encoding of the X.501 Name.
	Issuer []byte
	// SerialNumber is the certificate serial number.
	SerialNumber *big.Int
}

// DER returns the DER encoding of the IssuerAndSerialNumber.
func (i *IssuerAndSerialNumber) DER() []byte {
	body := append([]byte{}, i.Issuer...)
	body = append(body, asn1ber.EncodeInteger(i.SerialNumber)...)
	return asn1ber.WriteSequence(body)
}

// SignerIdentifier is
//
//	SignerIdentifier ::= CHOICE {
//	    issuerAndSerialNumber IssuerAndSerialNumber,
//	    subjectKeyIdentifier  [0] SubjectKeyIdentifier }
//
// Exactly one field is set. The subjectKeyIdentifier alternative forces the SignerInfo to
// version 3, see RFC 5652 clause 5.3.
type SignerIdentifier struct {
	// IssuerAndSerialNumber is set for the issuerAndSerialNumber alternative.
	IssuerAndSerialNumber *IssuerAndSerialNumber
	// SubjectKeyIdentifier holds the key identifier octets for the [0] alternative.
	SubjectKeyIdentifier []byte
}

// NewIssuerAndSerialNumberSID builds the issuerAndSerialNumber alternative from the encoding
// of the issuer name and the serial number.
func NewIssuerAndSerialNumberSID(issuer []byte, serialNumber *big.Int) *SignerIdentifier {
	return &SignerIdentifier{IssuerAndSerialNumber: &IssuerAndSerialNumber{Issuer: issuer, SerialNumber: serialNumber}}
}

// NewSubjectKeyIdentifierSID builds the subjectKeyIdentifier alternative.
func NewSubjectKeyIdentifierSID(subjectKeyIdentifier []byte) *SignerIdentifier {
	return &SignerIdentifier{SubjectKeyIdentifier: subjectKeyIdentifier}
}

// IsSubjectKeyIdentifier reports whether the [0] alternative was chosen.
func (s *SignerIdentifier) IsSubjectKeyIdentifier() bool { return s.IssuerAndSerialNumber == nil }

// DER returns the DER encoding of the SignerIdentifier. The subjectKeyIdentifier alternative
// is an implicitly tagged OCTET STRING, hence primitive.
func (s *SignerIdentifier) DER() []byte {
	if s.IssuerAndSerialNumber != nil {
		return s.IssuerAndSerialNumber.DER()
	}
	return asn1ber.WriteTLV(asn1ber.ClassContextSpecific|0, s.SubjectKeyIdentifier)
}

// signerIdentifierFromElement decodes the SignerIdentifier CHOICE.
func signerIdentifierFromElement(element *asn1ber.Element) (*SignerIdentifier, error) {
	if element.IsContextSpecific(0) {
		if element.IsConstructed() {
			// A constructed OCTET STRING is legal in BER; Octets joins its segments.
			return &SignerIdentifier{SubjectKeyIdentifier: element.Octets()}, nil
		}
		return &SignerIdentifier{SubjectKeyIdentifier: element.Content()}, nil
	}
	children, err := expectSizedSequence(element, "SignerInfo.sid", 2, 2)
	if err != nil {
		return nil, err
	}
	serialNumber, err := expectInteger(children[1], "IssuerAndSerialNumber.serialNumber")
	if err != nil {
		return nil, err
	}
	return &SignerIdentifier{IssuerAndSerialNumber: &IssuerAndSerialNumber{
		Issuer:       children[0].Encoded(),
		SerialNumber: serialNumber,
	}}, nil
}

// SignerInfo is
//
//	SignerInfo ::= SEQUENCE {
//	    version            CMSVersion,
//	    sid                SignerIdentifier,
//	    digestAlgorithm    DigestAlgorithmIdentifier,
//	    signedAttrs        [0] IMPLICIT SignedAttributes OPTIONAL,
//	    signatureAlgorithm SignatureAlgorithmIdentifier,
//	    signature          SignatureValue,
//	    unsignedAttrs      [1] IMPLICIT UnsignedAttributes OPTIONAL }
//
// A parsed SignerInfo keeps its own encoding and the encodings of both attribute fields: DSS
// derives a signature identifier from the SignerInfo bytes, and re-serialises the unsigned
// attributes when it augments a signature.
type SignerInfo struct {
	// Version is the CMSVersion, 1 or 3 (see ComputeSignerInfoVersion).
	Version int
	// SID identifies the signing certificate.
	SID *SignerIdentifier
	// DigestAlgorithm is the digest algorithm applied to the content and to signedAttrs.
	DigestAlgorithm *asn1ber.AlgorithmIdentifier
	// SignedAttributes holds the signedAttrs, nil when the field is absent.
	SignedAttributes Attributes
	// SignatureAlgorithm is the signature algorithm identifier.
	SignatureAlgorithm *asn1ber.AlgorithmIdentifier
	// Signature holds the signature octets.
	Signature []byte
	// UnsignedAttributes holds the unsignedAttrs, nil when the field is absent.
	UnsignedAttributes Attributes

	// element is the parsed SignerInfo, nil when it was built.
	element *asn1ber.Element
	// signedAttrsElement and unsignedAttrsElement are the parsed [0] and [1] fields.
	signedAttrsElement   *asn1ber.Element
	unsignedAttrsElement *asn1ber.Element
	// hasSignedAttrs and hasUnsignedAttrs record the presence of the optional fields, which
	// an empty attribute slice cannot express.
	hasSignedAttrs   bool
	hasUnsignedAttrs bool
}

// SignerInfoFromElement decodes an already parsed SignerInfo.
func SignerInfoFromElement(element *asn1ber.Element) (*SignerInfo, error) {
	children, err := expectSequence(element, "SignerInfo", 5)
	if err != nil {
		return nil, err
	}
	version, err := expectSmallInteger(children[0], "SignerInfo.version")
	if err != nil {
		return nil, err
	}
	sid, err := signerIdentifierFromElement(children[1])
	if err != nil {
		return nil, err
	}
	digestAlgorithm, err := asn1ber.AlgorithmIdentifierFromElement(children[2])
	if err != nil {
		return nil, wrapField("cmscore: SignerInfo.digestAlgorithm", err)
	}
	signerInfo := &SignerInfo{
		Version:         version,
		SID:             sid,
		DigestAlgorithm: digestAlgorithm,
		element:         element,
	}

	index := 3
	if children[index].IsContextSpecific(0) {
		if !children[index].IsConstructed() {
			return nil, errors.New("cmscore: SignerInfo.signedAttrs is not constructed")
		}
		attributes, err := attributesFromElement(children[index], "cmscore: SignerInfo.signedAttrs")
		if err != nil {
			return nil, err
		}
		signerInfo.SignedAttributes = attributes
		signerInfo.signedAttrsElement = children[index]
		signerInfo.hasSignedAttrs = true
		index++
	}
	if len(children) < index+2 {
		return nil, errors.New("cmscore: SignerInfo is truncated")
	}

	signatureAlgorithm, err := asn1ber.AlgorithmIdentifierFromElement(children[index])
	if err != nil {
		return nil, wrapField("cmscore: SignerInfo.signatureAlgorithm", err)
	}
	signerInfo.SignatureAlgorithm = signatureAlgorithm
	index++

	signature, err := expectOctetString(children[index], "SignerInfo.signature")
	if err != nil {
		return nil, err
	}
	signerInfo.Signature = signature
	index++

	if index < len(children) {
		if !children[index].IsContextSpecific(1) || !children[index].IsConstructed() {
			return nil, errors.New("cmscore: SignerInfo.unsignedAttrs is not [1] IMPLICIT")
		}
		attributes, err := attributesFromElement(children[index], "cmscore: SignerInfo.unsignedAttrs")
		if err != nil {
			return nil, err
		}
		signerInfo.UnsignedAttributes = attributes
		signerInfo.unsignedAttrsElement = children[index]
		signerInfo.hasUnsignedAttrs = true
		index++
	}
	if index != len(children) {
		return nil, fmt.Errorf("cmscore: SignerInfo holds %d unexpected trailing components", len(children)-index)
	}
	return signerInfo, nil
}

// Element returns the parsed SignerInfo, nil when it was built.
func (s *SignerInfo) Element() *asn1ber.Element { return s.element }

// Encoded returns the SignerInfo's original encoding, and its DER encoding when it was built
// rather than parsed. DSS digests these bytes to identify a signature.
func (s *SignerInfo) Encoded() []byte {
	if s.element != nil {
		return s.element.Encoded()
	}
	return s.DER()
}

// HasSignedAttributes reports whether the signedAttrs field is present. An empty but present
// field is possible in principle, so this is not "len(SignedAttributes) != 0".
func (s *SignerInfo) HasSignedAttributes() bool { return s.hasSignedAttrs }

// HasUnsignedAttributes reports whether the unsignedAttrs field is present.
func (s *SignerInfo) HasUnsignedAttributes() bool { return s.hasUnsignedAttrs }

// SignedAttributesElement returns the parsed signedAttrs field, nil when it is absent or when
// the SignerInfo was built.
func (s *SignerInfo) SignedAttributesElement() *asn1ber.Element { return s.signedAttrsElement }

// UnsignedAttributesElement returns the parsed unsignedAttrs field, nil when it is absent or
// when the SignerInfo was built.
func (s *SignerInfo) UnsignedAttributesElement() *asn1ber.Element { return s.unsignedAttrsElement }

// SignedAttributesRaw returns the signedAttrs field exactly as received, [0] IMPLICIT tag
// included, nil when the field is absent or the SignerInfo was built.
func (s *SignerInfo) SignedAttributesRaw() []byte {
	if s.signedAttrsElement == nil {
		return nil
	}
	return s.signedAttrsElement.Encoded()
}

// SignedAttributesDER returns the signed attributes as a DER SET OF, which RFC 5652
// clause 5.4 makes the input of the signature computation and which
// SignerInformation#getEncodedSignedAttributes returns. It is nil when the field is absent.
//
// The re-encoding is not always the received bytes: a producer that emitted the attributes
// out of DER order, or with a non-minimal length, gets normalised here - exactly as
// BouncyCastle normalises before verifying. SignedAttributesRaw keeps the original.
func (s *SignerInfo) SignedAttributesDER() []byte {
	if !s.hasSignedAttrs {
		return nil
	}
	return s.SignedAttributes.DERSetEncoded()
}

// DER returns the DER encoding of the SignerInfo.
func (s *SignerInfo) DER() []byte {
	body := encodeInt(s.Version)
	body = append(body, s.SID.DER()...)
	body = append(body, s.DigestAlgorithm.DER()...)
	if s.hasSignedAttrs {
		body = append(body, s.SignedAttributes.DERImplicitTagged(0)...)
	}
	body = append(body, s.SignatureAlgorithm.DER()...)
	body = append(body, encodeOctetString(s.Signature)...)
	if s.hasUnsignedAttrs {
		body = append(body, s.UnsignedAttributes.DERImplicitTagged(1)...)
	}
	return asn1ber.WriteSequence(body)
}
