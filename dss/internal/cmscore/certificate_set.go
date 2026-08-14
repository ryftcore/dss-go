// CertificateSet and CertificateChoices of RFC 5652 clause 10.2, replacing
// org.bouncycastle.asn1.cms.CertificateChoices and the certificate Store of
// org.bouncycastle.cms.CMSSignedData.
package cmscore

import (
	"github.com/utain/esig/dss/internal/asn1ber"
)

// Alternatives of the CertificateChoices CHOICE.
//
//	CertificateChoices ::= CHOICE {
//	    certificate          Certificate,
//	    extendedCertificate  [0] IMPLICIT ExtendedCertificate,   -- obsolete
//	    v1AttrCert           [1] IMPLICIT AttributeCertificateV1, -- obsolete
//	    v2AttrCert           [2] IMPLICIT AttributeCertificateV2,
//	    other                [3] IMPLICIT OtherCertificateFormat }
const (
	// CertificateChoiceCertificate is the untagged X.509 Certificate alternative.
	CertificateChoiceCertificate = -1
	// CertificateChoiceExtendedCertificate is the obsolete PKCS#6 alternative.
	CertificateChoiceExtendedCertificate = 0
	// CertificateChoiceV1AttrCert is the obsolete version 1 attribute certificate, which
	// forces SignedData to version 3.
	CertificateChoiceV1AttrCert = 1
	// CertificateChoiceV2AttrCert is the version 2 attribute certificate, which forces
	// SignedData to version 4.
	CertificateChoiceV2AttrCert = 2
	// CertificateChoiceOther is OtherCertificateFormat, which forces SignedData to version 5.
	CertificateChoiceOther = 3
)

// CertificateChoice is one member of a CertificateSet. Its encoding is preserved: a
// certificate is identified and validated by its exact DER, never by a re-encoding.
type CertificateChoice struct {
	// TagNo is the chosen alternative, CertificateChoiceCertificate for a plain certificate
	// and the context-specific tag number otherwise.
	TagNo int

	// encoded is the member's encoding, tag included.
	encoded []byte
	// element is the parsed member, nil when it was built.
	element *asn1ber.Element
}

// NewCertificateChoice builds the untagged certificate alternative from a certificate's DER.
func NewCertificateChoice(certificate []byte) CertificateChoice {
	return CertificateChoice{TagNo: CertificateChoiceCertificate, encoded: certificate}
}

// NewTaggedCertificateChoice builds a tagged alternative from the complete encoding of the
// tagged member, i.e. including its context-specific tag.
func NewTaggedCertificateChoice(tagNo int, encoded []byte) CertificateChoice {
	return CertificateChoice{TagNo: tagNo, encoded: encoded}
}

// IsCertificate reports whether the member is a plain X.509 certificate.
func (c CertificateChoice) IsCertificate() bool { return c.TagNo == CertificateChoiceCertificate }

// Encoded returns the member's original encoding.
func (c CertificateChoice) Encoded() []byte { return c.encoded }

// Element returns the parsed member, nil when it was built.
func (c CertificateChoice) Element() *asn1ber.Element { return c.element }

// DER returns the DER encoding of the member.
func (c CertificateChoice) DER() []byte {
	if c.element != nil {
		return c.element.DEREncoded()
	}
	return c.encoded
}

// CertificateSet is the SignedData.certificates field,
//
//	CertificateSet ::= SET OF CertificateChoices
//
// carried by a [0] IMPLICIT tag. The field being optional, an absent one is a nil
// *CertificateSet rather than an empty set.
type CertificateSet struct {
	// Choices holds the members in the order they were received.
	Choices []CertificateChoice

	// element is the parsed [0] IMPLICIT field, nil when it was built.
	element *asn1ber.Element
}

// certificateSetFromElement decodes the [0] IMPLICIT certificates field.
func certificateSetFromElement(element *asn1ber.Element) (*CertificateSet, error) {
	set := &CertificateSet{element: element}
	for _, child := range element.Children() {
		choice := CertificateChoice{TagNo: CertificateChoiceCertificate, encoded: child.Encoded(), element: child}
		if child.Class() == asn1ber.ClassContextSpecific {
			choice.TagNo = int(child.TagNumber())
		}
		set.Choices = append(set.Choices, choice)
	}
	return set, nil
}

// Element returns the parsed certificates field, nil when it was built or absent.
func (c *CertificateSet) Element() *asn1ber.Element {
	if c == nil {
		return nil
	}
	return c.element
}

// Encoded returns the certificates field exactly as received, [0] IMPLICIT tag included, and
// its DER encoding when the set was built rather than parsed. An absent field - a nil
// *CertificateSet - has no encoding.
func (c *CertificateSet) Encoded() []byte {
	if c == nil {
		return nil
	}
	if c.element != nil {
		return c.element.Encoded()
	}
	return c.DER()
}

// Certificates returns the encoding of every plain X.509 certificate in the set, in the
// received order.
func (c *CertificateSet) Certificates() [][]byte {
	if c == nil {
		return nil
	}
	var certificates [][]byte
	for _, choice := range c.Choices {
		if choice.IsCertificate() {
			certificates = append(certificates, choice.Encoded())
		}
	}
	return certificates
}

// AttributeCertificatesV2 returns the encoding of every [2] IMPLICIT AttributeCertificateV2 in
// the set, in the received order.
func (c *CertificateSet) AttributeCertificatesV2() [][]byte {
	if c == nil {
		return nil
	}
	var certificates [][]byte
	for _, choice := range c.Choices {
		if choice.TagNo == CertificateChoiceV2AttrCert {
			certificates = append(certificates, choice.Encoded())
		}
	}
	return certificates
}

// DER returns the DER encoding of the certificates field, [0] IMPLICIT tag included and
// members ordered as a DER SET OF. An absent field - a nil *CertificateSet - has none.
func (c *CertificateSet) DER() []byte {
	if c == nil {
		return nil
	}
	members := make([][]byte, len(c.Choices))
	for index, choice := range c.Choices {
		members[index] = choice.DER()
	}
	return derSetOf([]byte{contextIdentifier(0)}, members)
}
