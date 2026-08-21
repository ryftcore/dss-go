// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/definition/xades132/XAdES132Attribute.java (DSS 6.5.RC1).
package definition

import "github.com/utain/esig/dss/xml/common"

// XAdES132Attribute defines attributes for a XAdES 1.3.2 schema.
type XAdES132Attribute string

// XAdES132Attribute constants, one per XAdES 1.3.2 schema attribute name.
const (
	XAdES132Attribute_ENCODING         XAdES132Attribute = "ENCODING"
	XAdES132Attribute_ID               XAdES132Attribute = "ID"
	XAdES132Attribute_OBJECT_REFERENCE XAdES132Attribute = "OBJECT_REFERENCE"
	XAdES132Attribute_QUALIFIER        XAdES132Attribute = "QUALIFIER"
	XAdES132Attribute_REFERENCED_DATA  XAdES132Attribute = "REFERENCED_DATA"
	XAdES132Attribute_TARGET           XAdES132Attribute = "TARGET"
	XAdES132Attribute_URI              XAdES132Attribute = "URI"
)

// xades132attributeNames maps each constant to its wire attribute name (getAttributeName()).
var xades132attributeNames = map[XAdES132Attribute]string{
	XAdES132Attribute_ENCODING:         "Encoding",
	XAdES132Attribute_ID:               "Id",
	XAdES132Attribute_OBJECT_REFERENCE: "ObjectReference",
	XAdES132Attribute_QUALIFIER:        "Qualifier",
	XAdES132Attribute_REFERENCED_DATA:  "referencedData",
	XAdES132Attribute_TARGET:           "Target",
	XAdES132Attribute_URI:              "URI",
}

// AttributeName implements common.DSSAttribute. Ports getAttributeName().
func (a XAdES132Attribute) AttributeName() string {
	return xades132attributeNames[a]
}

var _ common.DSSAttribute = XAdES132Attribute("")
