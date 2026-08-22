// Stand-ins for the cross-namespace JAXB types 1910202xmlSchema.xsd
// dereferences from specs-xmldsig, specs-xades and specs-trusted-list -
// eu.europa.esig.xmldsig.jaxb.{DigestMethodType,SignatureValueType,SignatureType},
// eu.europa.esig.xades.jaxb.xades132.{DigestAlgAndValueType,SignaturePolicyIdentifierType}
// and eu.europa.esig.trustedlist211.jaxb.tsl.{DigitalIdentityType,TSPInformationType}.
// See doc.go's "Cross-namespace types" section for why these are modelled
// here rather than imported from a port of those modules (none exists;
// they are outside S8B_BRIEF.md's manifest).
//
// DigestMethodType/DigestAlgAndValueType/SignatureValueType are small,
// low-cardinality and well-known (unchanged since the 2002 xmldsig-core and
// 2003 XAdES 1.3.2 recommendations), so they are modelled field-for-field
// from those schemas. SignatureType, SignaturePolicyIdentifierType,
// DigitalIdentityType and TSPInformationType are modelled as raw-XML
// capture stand-ins instead: their content is preserved byte-for-byte on
// unmarshal->marshal without being interpreted, which is what the
// marshal-parity KAT needs and does not require porting those schemas'
// generated-JAXB trees.
//
// # The "ns2:" prefix on ds:-namespaced elements
//
// Every marshal-parity oracle in testdata/oracle declares the XMLDSig
// namespace at the document root with the JAXB-assigned prefix "ns2" (JAXB
// pre-declares every namespace the tree uses, once, at the document
// element, with a prefix generated the first time that namespace is
// encountered while walking the tree - and a validation report's very first
// piece of content, SignatureIdentifier.DigestAlgAndValue, is always the
// first thing to need one). encoding/xml has no equivalent behaviour: left
// to a plain `xml:"<namespace> Local"` field tag, it re-declares the
// namespace as a *local* default xmlns on every single element that needs
// it (`<DigestMethod xmlns="...">` instead of reusing a shared prefix),
// which is well-formed XML but not the JAXB reference implementation's
// bytes. DigestMethodType, DSDigestValue and SignatureValueType below
// implement xml.Marshaler and hard-code the "ns2:" prefix directly onto
// their own element name instead of leaving it to their field tag - a
// value's own MarshalXML overrides whatever name its enclosing field's
// struct tag would otherwise have produced, confirmed empirically and
// pinned by TestDSPrefixOverridesFieldTag in xml_kat_test.go - which is
// enough because every one of these types is used only where xmldsig-core
// content belongs, so there is nothing to disambiguate by tag. Root-level
// namespace declarations (including this "ns2" one) are captured from the
// source document and replayed verbatim by ValidationReportType's own
// MarshalXML - see jaxb_root.go - so this file does not need to invent the
// declaration, only the elements that use it.
package jaxb

import (
	"encoding/xml"

	"github.com/ryftcore/dss-go/dss/internal/xmldsig"
)

// namespaceDSig is the XMLDSig namespace 1910202xmlSchema.xsd's ds:-prefixed
// element refs resolve to. Struct tags below must spell it out as a literal
// (Go struct tags cannot reference a constant), so its provenance is pinned
// against dss/internal/xmldsig.NamespaceDSig - the same constant the CAdES/
// XAdES signing stack uses - by an equality assertion in xml_kat_test.go.
const namespaceDSig = "http://www.w3.org/2000/09/xmldsig#"

var _ = xmldsig.NamespaceDSig // see namespaceDSig above; referenced for the doc link only.

// DigestMethodType stands in for eu.europa.esig.xmldsig.jaxb.DigestMethodType
// (ds:DigestMethodType): the Algorithm attribute is what every use in this
// schema (SACertIDType, SACRLIDType, DigestAlgAndValueType) sets; the
// schema's optional ds:DigestMethod child extension content (HMACOutputLength)
// is not used by any DSS-generated report and is not modelled.
type DigestMethodType struct {
	Algorithm string `xml:"Algorithm,attr"`
}

// MarshalXML writes the element as ns2:DigestMethod - see this file's
// header.
func (d DigestMethodType) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name = xml.Name{Local: "ns2:DigestMethod"}
	if d.Algorithm != "" {
		start.Attr = []xml.Attr{{Name: xml.Name{Local: "Algorithm"}, Value: d.Algorithm}}
	} else {
		start.Attr = nil
	}
	if err := e.EncodeToken(start); err != nil {
		return err
	}
	return e.EncodeToken(start.End())
}

// DSDigestValue is Base64Binary bound to ds:DigestValue - the sibling of
// DigestMethodType in DigestAlgAndValueType, SACertIDType and SACRLIDType.
// It is a distinct type from Base64Binary (used, among other things, for
// this schema's own unprefixed vr:DigestValue,
// SignatureReferenceType.DigestValue) purely so it can carry the "ns2:"
// prefix unconditionally - see this file's header.
type DSDigestValue []byte

// MarshalXML writes the element as ns2:DigestValue - see this file's
// header.
func (v DSDigestValue) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name = xml.Name{Local: "ns2:DigestValue"}
	start.Attr = nil
	text, err := Base64Binary(v).MarshalText()
	if err != nil {
		return err
	}
	if err := e.EncodeToken(start); err != nil {
		return err
	}
	if err := e.EncodeToken(xml.CharData(text)); err != nil {
		return err
	}
	return e.EncodeToken(start.End())
}

// UnmarshalText decodes a base64Binary lexical form, matching Base64Binary.
func (v *DSDigestValue) UnmarshalText(text []byte) error {
	var b Base64Binary
	if err := b.UnmarshalText(text); err != nil {
		return err
	}
	*v = DSDigestValue(b)
	return nil
}

// DigestAlgAndValueType stands in for
// eu.europa.esig.xades.jaxb.xades132.DigestAlgAndValueType (XAdES 1.3.2's
// ds:DigestMethod + ds:DigestValue pair).
type DigestAlgAndValueType struct {
	DigestMethod DigestMethodType `xml:"http://www.w3.org/2000/09/xmldsig# DigestMethod"`
	DigestValue  DSDigestValue    `xml:"http://www.w3.org/2000/09/xmldsig# DigestValue"`
}

// SignatureValueType stands in for eu.europa.esig.xmldsig.jaxb.SignatureValueType
// (ds:SignatureValueType: a base64Binary simpleContent with an optional Id).
// Used only by SignatureIdentifierType.SignatureValue, so - like
// DigestMethodType/DSDigestValue - it hard-codes its own "ns2:" element
// name; see this file's header.
type SignatureValueType struct {
	Id    string       `xml:"Id,attr,omitempty"`
	Value Base64Binary `xml:",chardata"`
}

// MarshalXML writes the element as ns2:SignatureValue - see this file's
// header.
func (v SignatureValueType) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name = xml.Name{Local: "ns2:SignatureValue"}
	if v.Id != "" {
		start.Attr = []xml.Attr{{Name: xml.Name{Local: "Id"}, Value: v.Id}}
	} else {
		start.Attr = nil
	}
	text, err := v.Value.MarshalText()
	if err != nil {
		return err
	}
	if err := e.EncodeToken(start); err != nil {
		return err
	}
	if err := e.EncodeToken(xml.CharData(text)); err != nil {
		return err
	}
	return e.EncodeToken(start.End())
}

// SignatureType stands in for eu.europa.esig.xmldsig.jaxb.SignatureType
// (ds:Signature, ValidationReportType.Signature): content and the element's
// own attributes are captured verbatim (see jaxb_common.go's flattenAttrs),
// and - like DigestMethodType/DSDigestValue/SignatureValueType - its "ns2:"
// element name is hard-coded rather than left to its field tag; see this
// file's header. It is used only at that one position, so there is nothing
// to disambiguate by tag.
type SignatureType struct {
	Attrs   []xml.Attr `xml:",any,attr"`
	Content string     `xml:",innerxml"`
}

// MarshalXML writes the element as ns2:Signature with its captured
// attributes flattened back to literal prefixes.
func (s SignatureType) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name = xml.Name{Local: "ns2:Signature"}
	start.Attr = flattenAttrs(s.Attrs)
	type alias struct {
		Content string `xml:",innerxml"`
	}
	return e.EncodeElement(alias{Content: s.Content}, start)
}

// SignaturePolicyIdentifierType stands in for
// eu.europa.esig.xades.jaxb.xades132.SignaturePolicyIdentifierType
// (SignatureValidationPolicyType.SignaturePolicyIdentifier). The element
// itself ("SignaturePolicyIdentifier") is declared in this schema's own
// default namespace - only its *type* is XAdES's - so, unlike SignatureType,
// its name is left to its field tag; only its content and own attributes
// are captured verbatim (see jaxb_common.go's flattenAttrs).
type SignaturePolicyIdentifierType struct {
	Attrs   []xml.Attr `xml:",any,attr"`
	Content string     `xml:",innerxml"`
}

// MarshalXML writes the element's own attributes (flattened, see
// flattenAttrs) followed by its raw captured content.
func (s SignaturePolicyIdentifierType) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Attr = append(start.Attr, flattenAttrs(s.Attrs)...)
	type alias struct {
		Content string `xml:",innerxml"`
	}
	return e.EncodeElement(alias{Content: s.Content}, start)
}

// DigitalIdentityType stands in for
// eu.europa.esig.trustedlist211.jaxb.tsl.DigitalIdentityType
// (SignatureValidatorType.DigitalId). Like SignaturePolicyIdentifierType,
// the element itself ("DigitalId") is in this schema's own default
// namespace, so only its content and own attributes are captured verbatim.
type DigitalIdentityType struct {
	Attrs   []xml.Attr `xml:",any,attr"`
	Content string     `xml:",innerxml"`
}

// MarshalXML writes the element's own attributes (flattened, see
// flattenAttrs) followed by its raw captured content.
func (d DigitalIdentityType) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Attr = append(start.Attr, flattenAttrs(d.Attrs)...)
	type alias struct {
		Content string `xml:",innerxml"`
	}
	return e.EncodeElement(alias{Content: d.Content}, start)
}

// TSPInformationType stands in for
// eu.europa.esig.trustedlist211.jaxb.tsl.TSPInformationType
// (SignatureValidatorType.TSPInformation). Like SignaturePolicyIdentifierType,
// the element itself ("TSPInformation") is in this schema's own default
// namespace, so only its content and own attributes are captured verbatim.
type TSPInformationType struct {
	Attrs   []xml.Attr `xml:",any,attr"`
	Content string     `xml:",innerxml"`
}

// MarshalXML writes the element's own attributes (flattened, see
// flattenAttrs) followed by its raw captured content.
func (t TSPInformationType) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Attr = append(start.Attr, flattenAttrs(t.Attrs)...)
	type alias struct {
		Content string `xml:",innerxml"`
	}
	return e.EncodeElement(alias{Content: t.Content}, start)
}
