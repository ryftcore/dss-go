// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/identifier/IdentifierBuilder.java (DSS 6.5.RC1).
package model

// IdentifierBuilder builds an Identifier.
type IdentifierBuilder interface {
	// Build builds the Identifier. Port of build().
	Build() Identifier
}
