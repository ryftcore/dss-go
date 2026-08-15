// Ported from SimpleReport.xsd (DSS 6.5.RC1): the content model the
// marshaller needs to reproduce one JAXB spelling that encoding/xml cannot
// express on its own (see xml.go's "self-closing tags" quirk).
//
// dss/diagnostic/jaxb derives this table by reflection over the model
// (jaxb_content_model.go there) because DiagnosticData.xsd is too large to
// hand-verify and reuses a handful of element names across genuinely
// different content models. SimpleReport.xsd is small enough to enumerate
// directly and, checked against the schema, binds every element name to
// exactly one content kind everywhere it appears - there is no name that is
// a nested element under one parent and simpleContent under another - so a
// flat, hand-verified table is equivalent in behaviour and far shorter. The
// XSD-completeness sweep in xml_schema_test.go cross-checks this table
// against SimpleReport.xsd so drift is caught rather than silently wrong.
package jaxb

// complexElements holds the names of every element whose JAXB binding is a
// complexType with element (not character-data) content: the RI writes no
// character data at all for these, so an instance with no children collapses
// to a self-closing tag. Every other element name in the schema - every
// plain xs:string/xs:dateTime/enum leaf, and every simpleContent complexType
// (Message, Semantic, SignatureScope, the SignatureLevel/TimestampLevel/
// EAALevel wrappers, and the EAAPayload claim types) - carries character
// data even when that data is the empty string, so the RI keeps <X></X>
// rather than collapsing it; those default to kindText below.
var complexElements = map[string]bool{
	"SimpleReport":          true,
	"ValidationPolicy":      true,
	"PDFAInfo":              true,
	"ValidationMessages":    true,
	"AdESValidationDetails": true,
	"QualificationDetails":  true,
	"Signature":             true,
	"EAASignature":          true,
	"KeyBindingSignature":   true,
	"Timestamps":            true,
	"Timestamp":             true,
	"EvidenceRecords":       true,
	"EvidenceRecord":        true,
	"EAA":                   true,
	"CertificateChain":      true,
	"Certificate":           true,
	"TrustAnchors":          true,
	"TrustAnchor":           true,
	"EAAPayload":            true,
}

// carriesCharData reports whether an empty <name> element would have been
// written by JAXB as <name></name> rather than <name/>. The schema binds
// every element name to one content kind regardless of its enclosing
// element, so the enclosing-element stack that dss/diagnostic/jaxb consults
// is not needed here.
func carriesCharData(name string, _ []string) bool {
	return !complexElements[name]
}
