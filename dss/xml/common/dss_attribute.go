// Ported from dss-xml-common/src/main/java/eu/europa/esig/dss/xml/common/definition/DSSAttribute.java (DSS 6.5.RC1).
package common

// DSSAttribute is the XML attribute interface.
type DSSAttribute interface {
	// AttributeName returns the attribute name. Ports getAttributeName().
	AttributeName() string
}

// dssAttributeFromDefinition is the anonymous DSSAttribute the Java static factory
// DSSAttribute.fromDefinition builds.
type dssAttributeFromDefinition struct {
	attributeName string
}

// AttributeName implements DSSAttribute.
func (a *dssAttributeFromDefinition) AttributeName() string {
	return a.attributeName
}

// DSSAttributeFromDefinition creates a DSSAttribute from the given definition. Ports the
// static factory method DSSAttribute.fromDefinition(String).
func DSSAttributeFromDefinition(attributeName string) DSSAttribute {
	return &dssAttributeFromDefinition{attributeName: attributeName}
}
