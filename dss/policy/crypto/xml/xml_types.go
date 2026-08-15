// Ported (as plain encoding/xml-tagged structs, not full JAXB) from
// xsd/rfc5698.xsd and xsd/19322algocatxmlschema.xsd - see doc.go for why.
// Each type is named after the JAXB class upstream's XJC would generate for
// the same XSD complexType, so this file reads as a schema-shaped model in
// the same spirit as dss/policy/jaxb and dss/diagnostic/jaxb, even though it
// is hand-written rather than xjc-derived.
package cryptoxml

import "encoding/xml"

// dsscNS is the RFC 5698 ("Digital Signature Service Core Protocols,
// Elements, and Bindings") target namespace every dssc:-prefixed element in
// this file belongs to.
const dsscNS = "urn:ietf:params:xml:ns:dssc"

// algocatNS is the ETSI TS 119 322 v1.1.1 extension namespace
// "MoreDetails"/"AlgorithmUsage"/"Recommendation" belong to.
const algocatNS = "http://uri.etsi.org/19322/v1.1.1#"

// SecuritySuitabilityPolicyType is the Go form of the generated JAXB class
// for the dssc:SecuritySuitabilityPolicyType complexType, the document
// element of an RFC 5698 cryptographic suite catalogue.
//
// The optional trailing ds:Signature element is not modelled: no in-scope
// code path reads it (see doc.go).
type SecuritySuitabilityPolicyType struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:dssc SecuritySuitabilityPolicy"`

	PolicyName      PolicyNameType  `xml:"urn:ietf:params:xml:ns:dssc PolicyName"`
	Publisher       PublisherType   `xml:"urn:ietf:params:xml:ns:dssc Publisher"`
	PolicyIssueDate string          `xml:"urn:ietf:params:xml:ns:dssc PolicyIssueDate"`
	NextUpdate      string          `xml:"urn:ietf:params:xml:ns:dssc NextUpdate,omitempty"`
	Usage           string          `xml:"urn:ietf:params:xml:ns:dssc Usage,omitempty"`
	Algorithm       []AlgorithmType `xml:"urn:ietf:params:xml:ns:dssc Algorithm"`

	// Version/Lang/Id carry xs:attribute defaults ("1"/"en") the JAXB RI
	// applies during unmarshalling when the attribute is absent; Go's
	// encoding/xml does not read XSD defaults, so a present-but-empty
	// value here is resolved to the schema default by
	// buildXMLMetadata/versionOrDefault/langOrDefault in
	// cryptographic_suite_xml_catalogue.go rather than here.
	Version string `xml:"version,attr"`
	Lang    string `xml:"lang,attr"`
	Id      string `xml:"id,attr"`
}

// PolicyNameType is the Go form of the generated JAXB class for the
// dssc:PolicyNameType complexType.
type PolicyNameType struct {
	Name             string `xml:"urn:ietf:params:xml:ns:dssc Name"`
	ObjectIdentifier string `xml:"urn:ietf:params:xml:ns:dssc ObjectIdentifier,omitempty"`
	URI              string `xml:"urn:ietf:params:xml:ns:dssc URI,omitempty"`
}

// PublisherType is the Go form of the generated JAXB class for the
// dssc:PublisherType complexType.
type PublisherType struct {
	Name    string `xml:"urn:ietf:params:xml:ns:dssc Name"`
	Address string `xml:"urn:ietf:params:xml:ns:dssc Address,omitempty"`
	URI     string `xml:"urn:ietf:params:xml:ns:dssc URI,omitempty"`
}

// AlgorithmType is the Go form of the generated JAXB class for the
// dssc:AlgorithmType complexType.
//
// The trailing xs:any wildcard is not modelled: no in-scope code path
// reads it (see doc.go).
type AlgorithmType struct {
	AlgorithmIdentifier AlgorithmIdentifierType `xml:"urn:ietf:params:xml:ns:dssc AlgorithmIdentifier"`
	Evaluation          []EvaluationType        `xml:"urn:ietf:params:xml:ns:dssc Evaluation"`
	Information         *InformationType        `xml:"urn:ietf:params:xml:ns:dssc Information,omitempty"`
}

// AlgorithmIdentifierType is the Go form of the generated JAXB class for
// the dssc:AlgorithmIdentifierType complexType.
type AlgorithmIdentifierType struct {
	Name             string   `xml:"urn:ietf:params:xml:ns:dssc Name"`
	ObjectIdentifier []string `xml:"urn:ietf:params:xml:ns:dssc ObjectIdentifier"`
	URI              []string `xml:"urn:ietf:params:xml:ns:dssc URI,omitempty"`
}

// InformationType is the Go form of the generated JAXB class for the
// dssc:InformationType complexType.
type InformationType struct {
	Text []string `xml:"urn:ietf:params:xml:ns:dssc Text"`
}

// EvaluationType is the Go form of the generated JAXB class for the
// dssc:EvaluationType complexType.
//
// Java's `Object any` (the wildcard xs:any this element declares, holding
// the etsi19322:MoreDetails extension in every document this chunk ports -
// see CryptographicSuiteXmlCatalogue#getExtensionType's
// JAXBElement/instanceof/getName().getLocalPart().equals("MoreDetails")
// dance) becomes a plain namespace-and-name-tagged field: encoding/xml
// matches an element by its qualified name directly, so no runtime type
// check is needed here.
type EvaluationType struct {
	Parameter   []ParameterType `xml:"urn:ietf:params:xml:ns:dssc Parameter,omitempty"`
	Validity    ValidityType    `xml:"urn:ietf:params:xml:ns:dssc Validity"`
	MoreDetails *ExtensionType  `xml:"http://uri.etsi.org/19322/v1.1.1# MoreDetails,omitempty"`
}

// ValidityType is the Go form of the generated JAXB class for the
// dssc:ValidityType complexType. Start/End are xs:date (optionally
// timezone-qualified, per every End value in dss-crypto-suite.xml); parsed
// with parseXSDDate in cryptographic_suite_xml_catalogue.go rather than
// via a custom xml.Unmarshaler, to keep the raw XSD lexical value
// available for that function's error handling.
type ValidityType struct {
	Start string `xml:"urn:ietf:params:xml:ns:dssc Start,omitempty"`
	End   string `xml:"urn:ietf:params:xml:ns:dssc End,omitempty"`
}

// ParameterType is the Go form of the generated JAXB class for the
// dssc:ParameterType complexType.
//
// The trailing xs:any wildcard is not modelled: no in-scope code path
// reads it (see doc.go).
type ParameterType struct {
	Min  *int   `xml:"urn:ietf:params:xml:ns:dssc Min,omitempty"`
	Max  *int   `xml:"urn:ietf:params:xml:ns:dssc Max,omitempty"`
	Name string `xml:"name,attr"`
}

// ExtensionType is the Go form of the generated JAXB class for the
// ETSI TS 119 322 algocat schema's ExtensionType complexType (the
// etsi19322:MoreDetails element's content model).
type ExtensionType struct {
	AlgorithmUsage []string `xml:"http://uri.etsi.org/19322/v1.1.1# AlgorithmUsage,omitempty"`
	Recommendation string   `xml:"http://uri.etsi.org/19322/v1.1.1# Recommendation,omitempty"`
}
