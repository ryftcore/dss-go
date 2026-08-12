// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/identifier/X500NameIdentifier.java (DSS 6.5.RC1).
package model

// X500NameIdentifier is a unique identifier for a Relative Distinguished Name (RDN).
type X500NameIdentifier struct {
	IdentifierBase
}

// NewX500NameIdentifier builds an X500NameIdentifier over the DER encoding of the principal.
func NewX500NameIdentifier(x500Principal *X500Principal) *X500NameIdentifier {
	return &X500NameIdentifier{NewIdentifierBase("X500NameIdentifier", "RDN-", x500Principal.Encoded())}
}

// compile-time interface assertion.
var _ Identifier = (*X500NameIdentifier)(nil)
