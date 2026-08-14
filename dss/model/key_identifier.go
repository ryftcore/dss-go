// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/identifier/KeyIdentifier.java (DSS 6.5.RC1).
package model

// KeyIdentifier is a unique identifier for a java.security.Key object.
type KeyIdentifier struct {
	IdentifierBase
}

// NewKeyIdentifier builds a KeyIdentifier over the encoded form of the given key.
func NewKeyIdentifier(key Key) *KeyIdentifier {
	return &KeyIdentifier{NewIdentifierBase("KeyIdentifier", "PK-", key.Encoded())}
}

// compile-time interface assertion.
var _ Identifier = (*KeyIdentifier)(nil)
