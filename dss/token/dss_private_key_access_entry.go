// Ported from dss-token/src/main/java/eu/europa/esig/dss/token/DSSPrivateKeyAccessEntry.java (DSS 6.5.RC1).
package token

import "crypto"

// DSSPrivateKeyAccessEntry provides an interface to a token connection with an exposed
// (accessible) private key entry.
// NOTE: That does not mean that the cryptographic private key can be extracted.
// The interface is meant to only provide direct access to the private key.
// It is up to the underlying implementation to determine a way the private key can be accessed.
//
// DEVIATION: Java exposes the opaque java.security.PrivateKey, handed to the JCA Signature API
// for the actual signing. Go has no such indirection: a key capable of signing implements
// crypto.Signer directly, so PrivateKey() returns that instead.
type DSSPrivateKeyAccessEntry interface {
	DSSPrivateKeyEntry

	// PrivateKey gets the private key. Port of getPrivateKey().
	PrivateKey() crypto.Signer
}
