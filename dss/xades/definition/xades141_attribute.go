// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/definition/xades141/XAdES141Attribute.java (DSS 6.5.RC1).
package definition

import "github.com/utain/esig/dss/xml/common"

// XAdES141Attribute defines attributes for a XAdES 1.4.1 schema.
type XAdES141Attribute string

const (
	XAdES141Attribute_ID    XAdES141Attribute = "ID"
	XAdES141Attribute_ORDER XAdES141Attribute = "ORDER"
	XAdES141Attribute_URI   XAdES141Attribute = "URI"
)

// xades141attributeNames maps each constant to its wire attribute name (getAttributeName()).
var xades141attributeNames = map[XAdES141Attribute]string{
	XAdES141Attribute_ID:    "Id",
	XAdES141Attribute_ORDER: "Order",
	XAdES141Attribute_URI:   "URI",
}

// AttributeName implements common.DSSAttribute. Ports getAttributeName().
func (a XAdES141Attribute) AttributeName() string {
	return xades141attributeNames[a]
}

var _ common.DSSAttribute = XAdES141Attribute("")
