// Ported from SimpleCertificateReport.xsd (DSS 6.5.RC1): the content model
// the marshaller needs to reproduce one JAXB spelling that encoding/xml
// cannot express on its own (see xml.go's "self-closing tags" quirk). Same
// rationale and design as dss/simplereport/jaxb/jaxb_content_model.go
// (Phase 8b): the schema is small enough, and binds every element name to
// exactly one content kind everywhere it appears, that a flat hand-verified
// table is equivalent to (and far shorter than) deriving it by reflection as
// dss/diagnostic/jaxb does. The XSD-completeness sweep in
// xml_schema_test.go cross-checks this table against
// SimpleCertificateReport.xsd so drift is caught rather than silently
// wrong.
package jaxb

// complexElements holds the names of every element whose JAXB binding is a
// complexType with element (not character-data) content, including every
// @XmlElementWrapper collection element: the RI writes no character data at
// all for these, so an instance with no children collapses to a
// self-closing tag. Every other element name in the schema - every plain
// xs:string/xs:dateTime/xs:boolean/enum leaf (including the URI-based
// ListType/ServiceTypeIdentifier/ServiceStatus enums, which marshal through
// encoding.TextMarshaler like any other adapter), and every simpleContent
// complexType (Message, SignatureScope) - carries character data even when
// that data is the empty string, so the RI keeps <X></X> rather than
// collapsing it; those default to kindText below.
var complexElements = map[string]bool{
	"SimpleCertificateReport":          true,
	"ValidationPolicy":                 true,
	"ConnectionDetails":                true,
	"Certificate":                      true,
	"ChainItem":                        true,
	"Chain":                            true,
	"subject":                          true,
	"keyUsages":                        true,
	"extendedKeyUsages":                true,
	"ocspUrls":                         true,
	"crlUrls":                          true,
	"aiaUrls":                          true,
	"cpsUrls":                          true,
	"pdsUrls":                          true,
	"qualificationDetailsAtIssuance":   true,
	"qualificationDetailsAtValidation": true,
	"qwacDetails":                      true,
	"X509ValidationDetails":            true,
	"Details":                          true,
	"certificateApprovalStatusAtIssuanceTime":   true,
	"certificateApprovalStatusAtValidationTime": true,
	"certificateApprovalStatus":                 true,
	"revocation":                                true,
	"trustAnchors":                              true,
	"trustAnchor":                               true,
	"TLSBindingSignature":                       true,
}

// carriesCharData reports whether an empty <name> element would have been
// written by JAXB as <name></name> rather than <name/>. The schema binds
// every element name to one content kind regardless of its enclosing
// element, so the enclosing-element stack that dss/diagnostic/jaxb consults
// is not needed here.
func carriesCharData(name string, _ []string) bool {
	return !complexElements[name]
}
