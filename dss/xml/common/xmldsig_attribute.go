// Ported from dss-xml-common/src/main/java/eu/europa/esig/dss/xml/common/definition/xmldsig/XMLDSigAttribute.java (DSS 6.5.RC1).
package common

// XMLDSigAttribute is an attribute defined in https://www.w3.org/TR/xmldsig-core1/. Ports
// the Java enum: a typed string whose value is the Java name(), with the actual attribute
// name (which differs from the Go/Java constant name for MIME_TYPE) held in a lookup table.
type XMLDSigAttribute string

// XMLDSigAttribute constants, one per XMLDSig schema attribute name.
const (
	XMLDSigAttributeAlgorithm XMLDSigAttribute = "ALGORITHM"
	XMLDSigAttributeEncoding  XMLDSigAttribute = "ENCODING"
	XMLDSigAttributeID        XMLDSigAttribute = "ID"
	XMLDSigAttributeMIMEType  XMLDSigAttribute = "MIME_TYPE"
	XMLDSigAttributeTarget    XMLDSigAttribute = "TARGET"
	XMLDSigAttributeType      XMLDSigAttribute = "TYPE"
	XMLDSigAttributeURI       XMLDSigAttribute = "URI"
)

// xmldsigAttributeNames maps each constant to its wire attribute name (getAttributeName()).
var xmldsigAttributeNames = map[XMLDSigAttribute]string{
	XMLDSigAttributeAlgorithm: "Algorithm",
	XMLDSigAttributeEncoding:  "Encoding",
	XMLDSigAttributeID:        "Id",
	XMLDSigAttributeMIMEType:  "MimeType",
	XMLDSigAttributeTarget:    "Target",
	XMLDSigAttributeType:      "Type",
	XMLDSigAttributeURI:       "URI",
}

// AttributeName implements DSSAttribute. Ports getAttributeName().
func (a XMLDSigAttribute) AttributeName() string {
	return xmldsigAttributeNames[a]
}

var _ DSSAttribute = XMLDSigAttribute("")
