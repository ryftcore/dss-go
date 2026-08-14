// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/identifier/EntityIdentifier.java (DSS 6.5.RC1).
package model

// EntityIdentifier is a unique id for a public key and subject name combination.
type EntityIdentifier struct {
	IdentifierBase
}

// NewEntityIdentifier builds an EntityIdentifier over the given binaries.
func NewEntityIdentifier(binaries []byte) *EntityIdentifier {
	return &EntityIdentifier{NewIdentifierBase("EntityIdentifier", "EK-", binaries)}
}

// compile-time interface assertion.
var _ Identifier = (*EntityIdentifier)(nil)
