// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/identifier/TokenIdentifier.java (DSS 6.5.RC1).
package model

// TokenIdentifier is a unique id for a Token. It is the port of the abstract Java class of
// the same name and is meant to be embedded by the concrete token identifiers.
type TokenIdentifier struct {
	MultipleDigestIdentifier
}

// NewTokenIdentifierFromToken computes an identifier from the encoded form of the given
// token. Port of the protected TokenIdentifier(String, Token) constructor; className carries
// the Java simple class name of the concrete subclass.
func NewTokenIdentifierFromToken(className, prefix string, token Token) TokenIdentifier {
	return NewTokenIdentifier(className, prefix, token.Encoded())
}

// NewTokenIdentifier builds an identifier from the provided token binaries. Port of the
// protected TokenIdentifier(String, byte[]) constructor.
func NewTokenIdentifier(className, prefix string, binaries []byte) TokenIdentifier {
	return TokenIdentifier{NewMultipleDigestIdentifier(className, prefix, binaries)}
}
