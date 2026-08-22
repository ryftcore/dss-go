// Ported from dss-jades/src/main/java/eu/europa/esig/dss/jades/validation/JAdESAttribute.java
// (DSS 6.5.RC1).
//
// EtsiUComponent (same package) extends this class, reaching the unexported name/value/identifier
// fields via embedding + promoted-field writes - legal because both types live in package jades,
// mirroring Java's "protected field, same-package subclass" relationship.
package jades

import (
	"github.com/ryftcore/dss-go/dss/spi/validation"
	"github.com/ryftcore/dss-go/dss/spi/validation/identifier"
)

// JAdESAttribute represents the JAdES header. Port of the class JAdESAttribute, implementing
// spi/validation.SignatureAttribute.
type Attribute struct {
	// name if the header. Port of the protected String name field.
	name string

	// value is the component's value. Port of the protected Object value field.
	value any

	// identifier identifies the instance, computed lazily. Port of the protected
	// AttributeIdentifier identifier field.
	identifier *AttributeIdentifier
}

// NewJAdESAttribute is the default constructor. Port of the public JAdESAttribute(String, Object)
// constructor.
func NewJAdESAttribute(name string, value any) *Attribute {
	return &Attribute{name: name, value: value}
}

// HeaderName gets the header's name. Port of getHeaderName().
func (a *Attribute) HeaderName() string {
	return a.name
}

// Value gets the value. Port of getValue().
func (a *Attribute) Value() any {
	return a.value
}

// Identifier gets the attribute identifier, computing it via AttributeIdentifierBuild on
// first use. Port of getIdentifier(), implementing spi/validation.SignatureAttribute.
func (a *Attribute) Identifier() identifier.SignatureAttributeIdentifier {
	if a.identifier == nil {
		a.identifier = AttributeIdentifierBuild(a.name, a.value)
	}
	return a.identifier.SignatureAttributeIdentifier
}

// Equals is the port of equals(Object), comparing by Identifier.
//
// Java's equals() also requires getClass() == o.getClass(), which makes a JAdESAttribute and an
// EtsiUComponent sharing the same identifier bytes unequal even though their computed identifiers
// match; EtsiUComponent is not part of this manifest (see the file header), so this narrower,
// same-concrete-type comparison mirrors the CAdESAttribute/XAdESAttribute precedent rather than
// attempting to reproduce that cross-type distinction here.
func (a *Attribute) Equals(other *Attribute) bool {
	if a == other {
		return true
	}
	if other == nil {
		return false
	}
	selfID := a.Identifier()
	otherID := other.Identifier()
	return selfID.Equals(&otherID)
}

// compile-time assertion: a Attribute satisfies spi/validation.SignatureAttribute.
var _ validation.SignatureAttribute = (*Attribute)(nil)
