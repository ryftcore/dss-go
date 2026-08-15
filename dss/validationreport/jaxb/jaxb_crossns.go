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
package jaxb

import "github.com/utain/esig/dss/internal/xmldsig"

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

// DigestAlgAndValueType stands in for
// eu.europa.esig.xades.jaxb.xades132.DigestAlgAndValueType (XAdES 1.3.2's
// ds:DigestMethod + ds:DigestValue pair).
type DigestAlgAndValueType struct {
	DigestMethod DigestMethodType `xml:"http://www.w3.org/2000/09/xmldsig# DigestMethod"`
	DigestValue  Base64Binary     `xml:"http://www.w3.org/2000/09/xmldsig# DigestValue"`
}

// SignatureValueType stands in for eu.europa.esig.xmldsig.jaxb.SignatureValueType
// (ds:SignatureValueType: a base64Binary simpleContent with an optional Id).
type SignatureValueType struct {
	Id    string       `xml:"Id,attr,omitempty"`
	Value Base64Binary `xml:",chardata"`
}

// SignatureType stands in for eu.europa.esig.xmldsig.jaxb.SignatureType
// (ds:Signature, ValidationReportType.Signature): captured verbatim, see
// this file's header.
type SignatureType struct {
	Content string `xml:",innerxml"`
}

// SignaturePolicyIdentifierType stands in for
// eu.europa.esig.xades.jaxb.xades132.SignaturePolicyIdentifierType
// (SignatureValidationPolicyType.SignaturePolicyIdentifier): captured
// verbatim, see this file's header.
type SignaturePolicyIdentifierType struct {
	Content string `xml:",innerxml"`
}

// DigitalIdentityType stands in for
// eu.europa.esig.trustedlist211.jaxb.tsl.DigitalIdentityType
// (SignatureValidatorType.DigitalId): captured verbatim, see this file's
// header.
type DigitalIdentityType struct {
	Content string `xml:",innerxml"`
}

// TSPInformationType stands in for
// eu.europa.esig.trustedlist211.jaxb.tsl.TSPInformationType
// (SignatureValidatorType.TSPInformation): captured verbatim, see this
// file's header.
type TSPInformationType struct {
	Content string `xml:",innerxml"`
}
