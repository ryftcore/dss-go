// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/definition/xades132/XAdES132Attribute.java (DSS 6.5.RC1).
package definition

import "github.com/ryftcore/dss-go/dss/xml/common"

// XAdES132Attribute defines attributes for a XAdES 1.3.2 schema.
type XAdES132Attribute string

// XAdES132Attribute constants, one per XAdES 1.3.2 schema attribute name.
const (
	XAdES132AttributeEncoding        XAdES132Attribute = "ENCODING"
	XAdES132AttributeID              XAdES132Attribute = "ID"
	XAdES132AttributeObjectReference XAdES132Attribute = "OBJECT_REFERENCE"
	XAdES132AttributeQualifier       XAdES132Attribute = "QUALIFIER"
	XAdES132AttributeReferencedData  XAdES132Attribute = "REFERENCED_DATA"
	XAdES132AttributeTarget          XAdES132Attribute = "TARGET"
	XAdES132AttributeURI             XAdES132Attribute = "URI"
)

// xades132attributeNames maps each constant to its wire attribute name (getAttributeName()).
var xades132attributeNames = map[XAdES132Attribute]string{
	XAdES132AttributeEncoding:        "Encoding",
	XAdES132AttributeID:              "Id",
	XAdES132AttributeObjectReference: "ObjectReference",
	XAdES132AttributeQualifier:       "Qualifier",
	XAdES132AttributeReferencedData:  "referencedData",
	XAdES132AttributeTarget:          "Target",
	XAdES132AttributeURI:             "URI",
}

// AttributeName implements common.DSSAttribute. Ports getAttributeName().
func (a XAdES132Attribute) AttributeName() string {
	return xades132attributeNames[a]
}

var _ common.DSSAttribute = XAdES132Attribute("")
