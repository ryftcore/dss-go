// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/random/SecureRandomProvider.java (DSS 6.5.RC1).
//
// DEVIATION: java.security.SecureRandom exposes nextBytes(byte[]) (and generateSeed(int),
// itself implemented in terms of nextBytes). Go's crypto/* and golang.org/x/crypto APIs that
// consume a randomness source (crypto/rsa.GenerateKey, crypto/ecdsa.GenerateKey, ...) all take
// it as an io.Reader, so GetSecureRandom returns an io.Reader rather than a SecureRandom-shaped
// type: io.ReadFull(reader, buf) is nextBytes(buf)'s idiomatic Go equivalent.
package random

import "io"

// SecureRandomProvider provides a pseudo-random byte source, deterministically derived from a
// seed. Ports SecureRandomProvider.
type SecureRandomProvider interface {
	// GetSecureRandom returns an io.Reader producing pseudo-random bytes deterministically
	// derived from seed. Port of getSecureRandom(byte[]).
	GetSecureRandom(seed []byte) io.Reader
}
