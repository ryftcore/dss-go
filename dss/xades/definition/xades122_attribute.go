// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/definition/xades122/XAdES122Attribute.java (DSS 6.5.RC1).
package definition

import "github.com/ryftcore/dss-go/dss/xml/common"

// XAdES122Attribute defines attributes for a XAdES 1.2.2 schema.
type XAdES122Attribute string

// XAdES122Attribute constants, one per XAdES 1.2.2 schema attribute name.
const (
	XAdES122Attribute_ID               XAdES122Attribute = "ID"
	XAdES122Attribute_OBJECT_REFERENCE XAdES122Attribute = "OBJECT_REFERENCE"
	XAdES122Attribute_QUALIFIER        XAdES122Attribute = "QUALIFIER"
	XAdES122Attribute_REFERENCED_DATA  XAdES122Attribute = "REFERENCED_DATA"
	XAdES122Attribute_TARGET           XAdES122Attribute = "TARGET"
	XAdES122Attribute_URI              XAdES122Attribute = "URI"
)

// xades122attributeNames maps each constant to its wire attribute name (getAttributeName()).
var xades122attributeNames = map[XAdES122Attribute]string{
	XAdES122Attribute_ID:               "Id",
	XAdES122Attribute_OBJECT_REFERENCE: "ObjectReference",
	XAdES122Attribute_QUALIFIER:        "Qualifier",
	XAdES122Attribute_REFERENCED_DATA:  "referencedData",
	XAdES122Attribute_TARGET:           "Target",
	XAdES122Attribute_URI:              "URI",
}

// AttributeName implements common.DSSAttribute. Ports getAttributeName().
func (a XAdES122Attribute) AttributeName() string {
	return xades122attributeNames[a]
}

var _ common.DSSAttribute = XAdES122Attribute("")
