// Ported from dss-cms/src/main/java/eu/europa/esig/dss/cms/asn1/CertifiedAttributesV2.java
// (DSS 6.5.RC1).
//
// "Basic support of ETSI EN 319 122-1 V1.1.1 chapter 5.2.6.1
//
//	CertifiedAttributesV2 ::= SEQUENCE OF CHOICE {
//	    attributeCertificate [0] AttributeCertificate,
//	    otherAttributeCertificate [1] OtherAttributeCertificate
//	}
//
// Note : OtherAttributeCertificate is not supported.
// Quote ETSI : The definition of specific otherAttributeCertificates is outside of the scope of
// the present document." (Java doc, kept for context.)
package cms

import (
	"fmt"

	"github.com/utain/esig/dss/internal/asn1ber"
)

// CertifiedAttributesV2 is a sequence of X.509 attribute certificates.
// Port of the CertifiedAttributesV2 class.
type CertifiedAttributesV2 struct {
	// attributeCertificates holds the DER encoding (untagged) of each attributeCertificate
	// alternative member, in the order they were built or received. An
	// otherAttributeCertificate member is dropped, as Java's own toASN1Primitive drops any
	// values[i] that is not an AttributeCertificate (logging a warning instead).
	attributeCertificates [][]byte
}

// NewCertifiedAttributesV2 builds a CertifiedAttributesV2 from the DER encoding of each
// AttributeCertificate. There is no direct Java constructor to port - upstream only ever
// builds one through the getInstance parse path or via SignerAttributeV2's own
// CertifiedAttributesV2-typed constructor, whose argument is otherwise already a
// CertifiedAttributesV2.
func NewCertifiedAttributesV2(attributeCertificates [][]byte) *CertifiedAttributesV2 {
	return &CertifiedAttributesV2{attributeCertificates: attributeCertificates}
}

// ParseCertifiedAttributesV2 decodes a CertifiedAttributesV2 from its encoding.
// Port of CertifiedAttributesV2.getInstance(Object).
func ParseCertifiedAttributesV2(encoded []byte) (*CertifiedAttributesV2, error) {
	element, rest, err := asn1ber.Parse(encoded)
	if err != nil {
		return nil, err
	}
	if len(rest) != 0 {
		return nil, fmt.Errorf("extra data found after the CertifiedAttributesV2")
	}
	if !element.IsUniversal(asn1ber.TagSequence) || !element.IsConstructed() {
		return nil, fmt.Errorf("malformed CertifiedAttributesV2: not a SEQUENCE")
	}
	attributes := &CertifiedAttributesV2{}
	for _, child := range element.Children() {
		if child.Class() != asn1ber.ClassContextSpecific {
			return nil, fmt.Errorf("illegal tag in CertifiedAttributesV2")
		}
		switch child.TagNumber() {
		case 0:
			base, err := cmsExplicitBase(child, "CertifiedAttributesV2.attributeCertificate")
			if err != nil {
				return nil, err
			}
			attributes.attributeCertificates = append(attributes.attributeCertificates, base.Encoded())
		case 1:
			// Upstream logs "OtherAttributeCertificate detected" and drops it.
		default:
			return nil, fmt.Errorf("illegal tag: %d", child.TagNumber())
		}
	}
	return attributes, nil
}

// AttributeCertificates returns the DER encoding (untagged) of each attributeCertificate
// member.
func (c *CertifiedAttributesV2) AttributeCertificates() [][]byte { return c.attributeCertificates }

// DER returns the DER encoding of the CertifiedAttributesV2. Port of #toASN1Primitive.
func (c *CertifiedAttributesV2) DER() []byte {
	var body []byte
	for _, attributeCertificate := range c.attributeCertificates {
		body = append(body, asn1ber.WriteTLV(asn1ber.ClassContextSpecific|asn1ber.Constructed|0, attributeCertificate)...)
	}
	return asn1ber.WriteSequence(body)
}

// cmsExplicitBase unwraps an explicitly tagged field, i.e. returns the single value the tagged
// object holds. Port of ASN1TaggedObject#getExplicitBaseObject(); shared by every asn1.go type
// of this package that reads an EXPLICIT CHOICE alternative.
func cmsExplicitBase(element *asn1ber.Element, name string) (*asn1ber.Element, error) {
	if !element.IsConstructed() || len(element.Children()) != 1 {
		return nil, fmt.Errorf("malformed %s: not explicitly tagged", name)
	}
	return element.Children()[0], nil
}
