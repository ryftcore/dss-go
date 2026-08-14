// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/definition/xades111/XAdES111Attribute.java (DSS 6.5.RC1).
package definition

import "github.com/utain/esig/dss/xml/common"

// XAdES111Attribute defines attributes for a XAdES 1.1.1 schema.
type XAdES111Attribute string

const (
	XAdES111Attribute_ID               XAdES111Attribute = "ID"
	XAdES111Attribute_OBJECT_REFERENCE XAdES111Attribute = "OBJECT_REFERENCE"
	XAdES111Attribute_QUALIFIER        XAdES111Attribute = "QUALIFIER"
	XAdES111Attribute_TARGET           XAdES111Attribute = "TARGET"
	XAdES111Attribute_URI              XAdES111Attribute = "URI"
	XAdES111Attribute_URI2             XAdES111Attribute = "URI2"
)

// xades111attributeNames maps each constant to its wire attribute name (getAttributeName()).
var xades111attributeNames = map[XAdES111Attribute]string{
	XAdES111Attribute_ID:               "Id",
	XAdES111Attribute_OBJECT_REFERENCE: "ObjectReference",
	XAdES111Attribute_QUALIFIER:        "Qualifier",
	XAdES111Attribute_TARGET:           "Target",
	XAdES111Attribute_URI:              "uri",
	XAdES111Attribute_URI2:             "URI",
}

// AttributeName implements common.DSSAttribute. Ports getAttributeName().
func (a XAdES111Attribute) AttributeName() string {
	return xades111attributeNames[a]
}

var _ common.DSSAttribute = XAdES111Attribute("")
