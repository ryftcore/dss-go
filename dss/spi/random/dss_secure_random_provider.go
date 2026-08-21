// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/random/DSSSecureRandomProvider.java (DSS 6.5.RC1).
//
// Upstream's inner class DSSFixedSecureRandom wraps org.bouncycastle.crypto.prng.FixedSecureRandom:
// it replays a precomputed byte block instead of drawing entropy, and once the block is
// exhausted, re-derives the next block by re-hashing the previous one with digestAlgorithm.
// FixedSecureRandom itself is nothing more than "hand out these bytes sequentially, then
// report exhausted" - a plain byte-buffer cursor - so this port drops the BouncyCastle type
// and keeps only that cursor behaviour directly on dssFixedSecureRandom.
package random

import (
	"io"

	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/spi"
)

// DSSSecureRandomProvider is the default SecureRandomProvider used in DSS, returning a
// deterministic hash-chain byte source keyed on the supplied seed. Ports DSSSecureRandomProvider.
type DSSSecureRandomProvider struct {
	// digestAlgorithm is the DigestAlgorithm used to derive each successive block.
	digestAlgorithm enumerations.DigestAlgorithm
}

// NewDSSSecureRandomProvider creates a DSSSecureRandomProvider using SHA512 to derive blocks
// (512-bit blocks). Port of the default constructor.
func NewDSSSecureRandomProvider() *DSSSecureRandomProvider {
	return &DSSSecureRandomProvider{digestAlgorithm: enumerations.DigestAlgorithm_SHA512}
}

// NewDSSSecureRandomProviderWithDigestAlgorithm creates a DSSSecureRandomProvider deriving
// blocks with digestAlgorithm. Port of DSSSecureRandomProvider(DigestAlgorithm).
func NewDSSSecureRandomProviderWithDigestAlgorithm(digestAlgorithm enumerations.DigestAlgorithm) *DSSSecureRandomProvider {
	if digestAlgorithm == "" {
		panic("DigestAlgorithm cannot be null!")
	}
	return &DSSSecureRandomProvider{digestAlgorithm: digestAlgorithm}
}

// GetSecureRandom returns a deterministic io.Reader seeded from seed. Port of
// getSecureRandom(byte[]).
func (p *DSSSecureRandomProvider) GetSecureRandom(seed []byte) io.Reader {
	return newDSSFixedSecureRandom(p.digestAlgorithm, seed)
}

// dssFixedSecureRandom generates a deterministic keystream by repeatedly re-hashing the seed:
// block[0] = digest(seed), block[n] = digest(block[n-1]); each block's bytes are read out
// sequentially before the next block is derived. Ports the combined behaviour of
// DSSSecureRandomProvider's DSSFixedSecureRandom inner class and the
// org.bouncycastle.crypto.prng.FixedSecureRandom it wraps.
type dssFixedSecureRandom struct {
	digestAlgorithm enumerations.DigestAlgorithm
	// block is the current, fully-derived byte block being read out.
	block []byte
	// index is the read cursor within block; a read is exhausted once index == len(block).
	index int
}

// newDSSFixedSecureRandom creates a dssFixedSecureRandom whose first block is digest(seed).
// Port of the private DSSFixedSecureRandom(byte[] seed) constructor.
func newDSSFixedSecureRandom(digestAlgorithm enumerations.DigestAlgorithm, seed []byte) *dssFixedSecureRandom {
	if seed == nil {
		panic("Seed cannot be null")
	}
	block, err := spi.DSSUtilsDigest(digestAlgorithm, seed)
	if err != nil {
		// digestAlgorithm is expected to already be a supported, validated algorithm by the
		// time a SecureRandomProvider is constructed; MessageDigest.getInstance failing here
		// has no clean Go analogue since GetSecureRandom returns no error (matching Java's
		// nextBytes(byte[]), which is also declared to return void).
		panic(err)
	}
	return &dssFixedSecureRandom{digestAlgorithm: digestAlgorithm, block: block}
}

// Read implements io.Reader, filling p deterministically from the block cursor, re-deriving
// the block by re-hashing whenever it is exhausted. Port of nextBytes(byte[]).
func (r *dssFixedSecureRandom) Read(p []byte) (int, error) {
	offset := 0
	for offset < len(p) {
		if r.index >= len(r.block) {
			next, err := spi.DSSUtilsDigest(r.digestAlgorithm, r.block)
			if err != nil {
				panic(err)
			}
			r.block = next
			r.index = 0
		}

		available := len(r.block) - r.index
		requested := len(p) - offset
		length := available
		if requested < length {
			length = requested
		}

		copy(p[offset:offset+length], r.block[r.index:r.index+length])

		offset += length
		r.index += length
	}
	return len(p), nil
}

// compile-time interface assertion.
var _ io.Reader = (*dssFixedSecureRandom)(nil)
var _ SecureRandomProvider = (*DSSSecureRandomProvider)(nil)
