// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/definition/xades141/XAdES141Attribute.java (DSS 6.5.RC1).
package definition

import "github.com/ryftcore/dss-go/dss/xml/common"

// XAdES141Attribute defines attributes for a XAdES 1.4.1 schema.
type XAdES141Attribute string

// XAdES141Attribute constants, one per XAdES 1.4.1 schema attribute name.
const (
	XAdES141AttributeID    XAdES141Attribute = "ID"
	XAdES141AttributeOrder XAdES141Attribute = "ORDER"
	XAdES141AttributeURI   XAdES141Attribute = "URI"
)

// xades141attributeNames maps each constant to its wire attribute name (getAttributeName()).
var xades141attributeNames = map[XAdES141Attribute]string{
	XAdES141AttributeID:    "Id",
	XAdES141AttributeOrder: "Order",
	XAdES141AttributeURI:   "URI",
}

// AttributeName implements common.DSSAttribute. Ports getAttributeName().
func (a XAdES141Attribute) AttributeName() string {
	return xades141attributeNames[a]
}

var _ common.DSSAttribute = XAdES141Attribute("")
