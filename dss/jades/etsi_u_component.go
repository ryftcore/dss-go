// Ported from dss-jades/src/main/java/eu/europa/esig/dss/jades/validation/EtsiUComponent.java (DSS 6.5.RC1).
//
// EtsiUComponent extends JAdESAttribute in Java (jades_attribute.go, same package): this port
// embeds it, reaching Attribute's unexported
// name/value/identifier fields directly from newEtsiUComponent - legal because both types live in
// package jades, mirroring the "protected field, same-package subclass" relationship the Java
// source expresses (see jades_attribute.go's own file header, which already documents this same
// expectation from EtsiUComponent's side).
package jades

// EtsiUComponent represents an item of the 'etsiU' header array. Port of the class
// EtsiUComponent, extending Attribute.
//
// java.io.Serializable is dropped (no Go counterpart).
type EtsiUComponent struct {
	Attribute

	// base64UrlEncoded reports whether the component is a base64url encoded instance.
	base64UrlEncoded bool

	// component is the component in its original representation.
	component any
}

// newEtsiUComponent is the port of the package-private constructor
// EtsiUComponent(Object, String, Object, AttributeIdentifier).
func newEtsiUComponent(component any, headerName string, value any, identifier *AttributeIdentifier) *EtsiUComponent {
	c := &EtsiUComponent{
		Attribute:        *NewJAdESAttribute(headerName, value),
		component:        component,
		base64UrlEncoded: DSSJsonUtilsIsStringFormat(component),
	}
	c.identifier = identifier
	return c
}

// EtsiUComponentBuild builds an EtsiUComponent from the 'etsiU' array entry. Port of the static
// build(Object, int).
func EtsiUComponentBuild(component any, order int) *EtsiUComponent {
	m, ok := DSSJsonUtilsParseEtsiUComponent(component)
	if ok && m != nil {
		keys := m.Keys()
		if len(keys) == 0 {
			return nil
		}
		headerName := keys[0]
		value := m.Value(headerName)
		identifier := AttributeIdentifierBuildWithOrder(headerName, value, &order)
		return newEtsiUComponent(component, headerName, value, identifier)
	}
	return nil
}

// EtsiUComponentBuildFromValue builds the EtsiUComponent from the given parameters. Port of the
// static build(String, Object, boolean, AttributeIdentifier).
func EtsiUComponentBuildFromValue(headerName string, value any, base64UrlEncoded bool,
	identifier *AttributeIdentifier) *EtsiUComponent {
	component := etsiUComponentCreate(headerName, value, base64UrlEncoded)
	return newEtsiUComponent(component, headerName, value, identifier)
}

// etsiUComponentCreate returns an 'etsiU' component in the defined representation. Port of the
// private static createEtsiUComponent(String, Object, boolean).
func etsiUComponentCreate(name string, value any, base64UrlEncoded bool) any {
	jsonObject := NewJsonObject()
	jsonObject.Put(name, value)
	if base64UrlEncoded {
		return DSSJsonUtilsToBase64UrlObject(jsonObject)
	}
	return jsonObject
}

// Component gets the attribute in its 'etsiU' member representation. Port of getComponent().
func (c *EtsiUComponent) Component() any {
	return c.component
}

// IsBase64UrlEncoded gets if the component is base64url encoded. Port of isBase64UrlEncoded().
func (c *EtsiUComponent) IsBase64UrlEncoded() bool {
	return c.base64UrlEncoded
}
