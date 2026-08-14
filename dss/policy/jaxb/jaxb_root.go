// Ported from policy.xsd (DSS 6.5.RC1) via the JAXB classes generated into
// eu.europa.esig.dss.policy.jaxb. See jaxb_common.go's file header for the
// generated-JAXB grouping rule.
package jaxb

import "encoding/xml"

// Namespace is the target namespace of policy.xsd. The schema declares
// elementFormDefault="qualified" and the JAXB RI binds it to the default
// namespace (no prefix), so only the document element carries an xmlns
// declaration - see xml.go and the oracle-backed KATs in jaxb_root_test.go.
const Namespace = "http://dss.esig.europa.eu/validation/policy"

// ConstraintsParameters is the Go form of the generated JAXB class
// ConstraintsParameters (complexType ConstraintsParameters, bound as the
// document element by ObjectFactory#createConstraintsParameters): the root
// element for the DSS validation policy, allowing customization of the
// validation process - add/remove checks and set their criticality.
type ConstraintsParameters struct {
	XMLName xml.Name `xml:"http://dss.esig.europa.eu/validation/policy ConstraintsParameters"`

	Description                    string                     `xml:"Description"`
	ContainerConstraints           *ContainerConstraints      `xml:"ContainerConstraints,omitempty"`
	PDFAConstraints                *PDFAConstraints           `xml:"PDFAConstraints,omitempty"`
	SignatureConstraints           *SignatureConstraints      `xml:"SignatureConstraints,omitempty"`
	CounterSignatureConstraints    *SignatureConstraints      `xml:"CounterSignatureConstraints,omitempty"`
	KeyBindingSignatureConstraints *SignatureConstraints      `xml:"KeyBindingSignatureConstraints,omitempty"`
	Timestamp                      *TimestampConstraints      `xml:"Timestamp,omitempty"`
	Revocation                     *RevocationConstraints     `xml:"Revocation,omitempty"`
	EvidenceRecord                 *EvidenceRecordConstraints `xml:"EvidenceRecord,omitempty"`
	EAA                            *EAAConstraints            `xml:"EAA,omitempty"`
	EAARevocation                  *EAARevocationConstraints  `xml:"EAARevocation,omitempty"`
	Cryptographic                  *CryptographicConstraint   `xml:"Cryptographic,omitempty"`
	Model                          *ModelConstraint           `xml:"Model,omitempty"`
	EIDAS                          *EIDAS                     `xml:"eIDAS,omitempty"`
	Name                           *string                    `xml:"Name,attr,omitempty"`
}

// ---------------------------------------------------------------------------
// The generated ObjectFactory (eu.europa.esig.dss.policy.jaxb.ObjectFactory)
// provides a create*() factory method per JAXB class, plus the
// @XmlElementDecl-annotated createConstraintsParameters(ConstraintsParameters)
// that wraps a value for marshalling as the document element. Go's
// composite literals are the natural equivalent of the trivial "return new
// X()" factories, so this port keeps only the document-element wrapping,
// which JAXBElement made an explicit step; every other ObjectFactory method
// is exercised (as &X{}) throughout this package's KATs and the ../ wrapper
// package instead of being re-declared as a same-signature indirection here.

// NewConstraintsParametersElement returns doc wrapped as the
// ConstraintsParameters document element, ready to pass to Marshal. Ports
// ObjectFactory#createConstraintsParameters(ConstraintsParameters), which
// JAXB requires to associate the QName with a value for a global element
// declared by reference (xsd:element name="ConstraintsParameters"
// type="ConstraintsParameters"); encoding/xml has no equivalent JAXBElement
// wrapper step, so this is the identity function once ConstraintsParameters
// itself already carries XMLName.
func NewConstraintsParametersElement(doc *ConstraintsParameters) *ConstraintsParameters {
	return doc
}
