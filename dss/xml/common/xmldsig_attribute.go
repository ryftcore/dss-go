// Ported from dss-xml-common/src/main/java/eu/europa/esig/dss/xml/common/definition/xmldsig/XMLDSigAttribute.java (DSS 6.5.RC1).
package common

// XMLDSigAttribute is an attribute defined in https://www.w3.org/TR/xmldsig-core1/. Ports
// the Java enum per PORTING.md's enum convention: a typed string whose value is the Java
// name(), with the actual attribute name (which differs from the Go/Java constant name for
// MIME_TYPE) held in a lookup table.
type XMLDSigAttribute string

// XMLDSigAttribute constants, one per XMLDSig schema attribute name.
const (
	XMLDSigAttribute_ALGORITHM XMLDSigAttribute = "ALGORITHM"
	XMLDSigAttribute_ENCODING  XMLDSigAttribute = "ENCODING"
	XMLDSigAttribute_ID        XMLDSigAttribute = "ID"
	XMLDSigAttribute_MIME_TYPE XMLDSigAttribute = "MIME_TYPE"
	XMLDSigAttribute_TARGET    XMLDSigAttribute = "TARGET"
	XMLDSigAttribute_TYPE      XMLDSigAttribute = "TYPE"
	XMLDSigAttribute_URI       XMLDSigAttribute = "URI"
)

// xmldsigAttributeNames maps each constant to its wire attribute name (getAttributeName()).
var xmldsigAttributeNames = map[XMLDSigAttribute]string{
	XMLDSigAttribute_ALGORITHM: "Algorithm",
	XMLDSigAttribute_ENCODING:  "Encoding",
	XMLDSigAttribute_ID:        "Id",
	XMLDSigAttribute_MIME_TYPE: "MimeType",
	XMLDSigAttribute_TARGET:    "Target",
	XMLDSigAttribute_TYPE:      "Type",
	XMLDSigAttribute_URI:       "URI",
}

// AttributeName implements DSSAttribute. Ports getAttributeName().
func (a XMLDSigAttribute) AttributeName() string {
	return xmldsigAttributeNames[a]
}

var _ DSSAttribute = XMLDSigAttribute("")
