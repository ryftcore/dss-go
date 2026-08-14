// Ported from dss-jades/src/main/java/eu/europa/esig/dss/jades/validation/JAdESAttribute.java
// (DSS 6.5.RC1).
//
// EtsiUComponent (Java, same package eu.europa.esig.dss.jades.validation, not in this manifest -
// see S6_BRIEF.md's VAL chunk split) extends this class, setting its protected name/value/
// identifier fields directly from its own constructor rather than through NewJAdESAttribute/the
// lazily-computed Identifier(). Its Go counterpart, landing in a sibling chunk of this same
// package, is therefore assumed to reach the unexported name/value/identifier fields below by
// embedding JAdESAttribute and writing to the promoted fields directly - legal because both types
// live in package jades, mirroring the "protected field, same-package subclass" relationship the
// Java source expresses.
package jades

import (
	"github.com/utain/esig/dss/spi/validation"
	"github.com/utain/esig/dss/spi/validation/identifier"
)

// JAdESAttribute represents the JAdES header. Port of the class JAdESAttribute, implementing
// spi/validation.SignatureAttribute.
type JAdESAttribute struct {
	// name if the header. Port of the protected String name field.
	name string

	// value is the component's value. Port of the protected Object value field.
	value any

	// identifier identifies the instance, computed lazily. Port of the protected
	// JAdESAttributeIdentifier identifier field.
	identifier *JAdESAttributeIdentifier
}

// NewJAdESAttribute is the default constructor. Port of the public JAdESAttribute(String, Object)
// constructor.
func NewJAdESAttribute(name string, value any) *JAdESAttribute {
	return &JAdESAttribute{name: name, value: value}
}

// HeaderName gets the header's name. Port of getHeaderName().
func (a *JAdESAttribute) HeaderName() string {
	return a.name
}

// Value gets the value. Port of getValue().
func (a *JAdESAttribute) Value() any {
	return a.value
}

// Identifier gets the attribute identifier, computing it via JAdESAttributeIdentifierBuild on
// first use. Port of getIdentifier(), implementing spi/validation.SignatureAttribute.
func (a *JAdESAttribute) Identifier() identifier.SignatureAttributeIdentifier {
	if a.identifier == nil {
		a.identifier = JAdESAttributeIdentifierBuild(a.name, a.value)
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
func (a *JAdESAttribute) Equals(other *JAdESAttribute) bool {
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

// compile-time assertion: a JAdESAttribute satisfies spi/validation.SignatureAttribute.
var _ validation.SignatureAttribute = (*JAdESAttribute)(nil)
