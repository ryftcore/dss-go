// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/definition/xades111/XAdES111Attribute.java (DSS 6.5.RC1).
package definition

import "github.com/ryftcore/dss-go/dss/xml/common"

// XAdES111Attribute defines attributes for a XAdES 1.1.1 schema.
type XAdES111Attribute string

// XAdES111Attribute constants, one per XAdES 1.1.1 schema attribute name.
const (
	XAdES111AttributeID              XAdES111Attribute = "ID"
	XAdES111AttributeObjectReference XAdES111Attribute = "OBJECT_REFERENCE"
	XAdES111AttributeQualifier       XAdES111Attribute = "QUALIFIER"
	XAdES111AttributeTarget          XAdES111Attribute = "TARGET"
	XAdES111AttributeURI             XAdES111Attribute = "URI"
	XAdES111AttributeURI2            XAdES111Attribute = "URI2"
)

// xades111attributeNames maps each constant to its wire attribute name (getAttributeName()).
var xades111attributeNames = map[XAdES111Attribute]string{
	XAdES111AttributeID:              "Id",
	XAdES111AttributeObjectReference: "ObjectReference",
	XAdES111AttributeQualifier:       "Qualifier",
	XAdES111AttributeTarget:          "Target",
	XAdES111AttributeURI:             "uri",
	XAdES111AttributeURI2:            "URI",
}

// AttributeName implements common.DSSAttribute. Ports getAttributeName().
func (a XAdES111Attribute) AttributeName() string {
	return xades111attributeNames[a]
}

var _ common.DSSAttribute = XAdES111Attribute("")
