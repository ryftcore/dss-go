// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/definition/xades122/XAdES122Attribute.java (DSS 6.5.RC1).
package definition

import "github.com/ryftcore/dss-go/dss/xml/common"

// XAdES122Attribute defines attributes for a XAdES 1.2.2 schema.
type XAdES122Attribute string

// XAdES122Attribute constants, one per XAdES 1.2.2 schema attribute name.
const (
	XAdES122AttributeID              XAdES122Attribute = "ID"
	XAdES122AttributeObjectReference XAdES122Attribute = "OBJECT_REFERENCE"
	XAdES122AttributeQualifier       XAdES122Attribute = "QUALIFIER"
	XAdES122AttributeReferencedData  XAdES122Attribute = "REFERENCED_DATA"
	XAdES122AttributeTarget          XAdES122Attribute = "TARGET"
	XAdES122AttributeURI             XAdES122Attribute = "URI"
)

// xades122attributeNames maps each constant to its wire attribute name (getAttributeName()).
var xades122attributeNames = map[XAdES122Attribute]string{
	XAdES122AttributeID:              "Id",
	XAdES122AttributeObjectReference: "ObjectReference",
	XAdES122AttributeQualifier:       "Qualifier",
	XAdES122AttributeReferencedData:  "referencedData",
	XAdES122AttributeTarget:          "Target",
	XAdES122AttributeURI:             "URI",
}

// AttributeName implements common.DSSAttribute. Ports getAttributeName().
func (a XAdES122Attribute) AttributeName() string {
	return xades122attributeNames[a]
}

var _ common.DSSAttribute = XAdES122Attribute("")
