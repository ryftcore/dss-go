package pfx

import (
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/asn1"
	"errors"
	"hash"
	"math/big"
)

// pfxHash names the digest RFC 7292 Appendix B derives key material with: its Go constructor,
// its output size u (bits, "HASH FUNCTION" column of Appendix B.2) and its compression input
// size v (bits, same table), both expressed here in bytes since every producer this package
// reads keeps password, salt and digest output byte-aligned.
type pfxHash struct {
	new func() hash.Hash
	u   int // digest size in bytes
	v   int // compression block size in bytes
}

var (
	pfxHashSHA1   = pfxHash{sha1.New, sha1.Size, 64}
	pfxHashSHA224 = pfxHash{sha256.New224, sha256.Size224, 64}
	pfxHashSHA256 = pfxHash{sha256.New, sha256.Size, 64}
	pfxHashSHA384 = pfxHash{sha512.New384, sha512.Size384, 128}
	pfxHashSHA512 = pfxHash{sha512.New, sha512.Size, 128}
)

// bmpStringPassword encodes a password the way RFC 7292 Appendix B.1 requires: UCS-2 (the BMP
// subset of UTF-16, big-endian) with a two-byte zero terminator. Every character DSS's own
// fixtures' passwords use is within the BMP, so a rune that would need a surrogate pair (none
// occur here) is rejected rather than silently mis-encoded.
func bmpStringPassword(password string) ([]byte, error) {
	encoded := make([]byte, 0, 2*len(password)+2)
	for _, r := range password {
		if r > 0xFFFF {
			return nil, errPasswordNotBMP
		}
		encoded = append(encoded, byte(r>>8), byte(r))
	}
	return append(encoded, 0, 0), nil
}

// errPasswordNotBMP reports a password character outside the Basic Multilingual Plane, which
// RFC 7292's BMPString password encoding cannot represent.
var errPasswordNotBMP = errors.New("pfx: password contains a character outside the Basic Multilingual Plane")

// fillWithRepeats returns v-byte-block-aligned repeats of pattern, RFC 7292 Appendix B.2 steps 2
// and 3 ("Concatenate copies of the salt/password together to create a string of length
// v*ceiling(len/v), the final copy truncated as needed"). The empty pattern stays empty.
func fillWithRepeats(pattern []byte, v int) []byte {
	if len(pattern) == 0 {
		return nil
	}
	blocks := (len(pattern) + v - 1) / v
	out := make([]byte, 0, blocks*v)
	for len(out) < blocks*v {
		out = append(out, pattern...)
	}
	return out[:blocks*v]
}

// deriveKeyMaterial implements the RFC 7292 Appendix B.2 pseudorandom-bit generator: purposeID
// selects what the output feeds (1 = encryption/decryption key, 2 = IV, 3 = MAC key), and size
// is the number of output bytes wanted.
func deriveKeyMaterial(h pfxHash, purposeID byte, salt, password []byte, iterations, size int) []byte {
	// Step 1: D, the "diversifier", is v/8 copies of purposeID.
	diversifier := make([]byte, h.v)
	for i := range diversifier {
		diversifier[i] = purposeID
	}

	// Steps 2-4: S = salt repeated to a multiple of v bytes, P = password likewise, I = S||P.
	s := fillWithRepeats(salt, h.v)
	p := fillWithRepeats(password, h.v)
	i := append(append([]byte{}, s...), p...)

	// Step 5: c = ceiling(size/u) hash blocks are needed.
	blocks := (size + h.u - 1) / h.u
	a := make([]byte, 0, blocks*h.u)

	for round := 0; round < blocks; round++ {
		// Step 6.A: Ai = H^iterations(D || I).
		digest := h.new()
		digest.Write(diversifier)
		digest.Write(i)
		ai := digest.Sum(nil)
		for iter := 1; iter < iterations; iter++ {
			digest.Reset()
			digest.Write(ai)
			ai = digest.Sum(nil)
		}
		a = append(a, ai...)

		if round == blocks-1 {
			break // I is only advanced to seed a further round; the last round needs none.
		}

		// Step 6.B: B is v bytes of Ai repeated.
		b := make([]byte, 0, h.v)
		for len(b) < h.v {
			b = append(b, ai...)
		}
		bValue := new(big.Int).SetBytes(b[:h.v])

		// Step 6.C: each v-byte block I_j of I becomes (I_j + B + 1) mod 2^v.
		blockCount := len(i) / h.v
		modulus := new(big.Int).Lsh(big.NewInt(1), uint(8*h.v))
		for j := 0; j < blockCount; j++ {
			block := new(big.Int).SetBytes(i[j*h.v : (j+1)*h.v])
			block.Add(block, bValue)
			block.Add(block, big.NewInt(1))
			block.Mod(block, modulus)
			blockBytes := block.Bytes()
			padded := i[j*h.v : (j+1)*h.v]
			for k := range padded {
				padded[k] = 0
			}
			copy(padded[h.v-len(blockBytes):], blockBytes)
		}
	}

	return a[:size]
}

// pfxHashByOID resolves the plain digest OID a MacData or a PBKDF2 PRF names to its pfxHash,
// nil when unrecognised.
func pfxHashByOID(oid asn1.ObjectIdentifier) *pfxHash {
	for _, candidate := range []struct {
		oid  asn1.ObjectIdentifier
		hash pfxHash
	}{
		{oidSHA1, pfxHashSHA1},
		{oidSHA224, pfxHashSHA224},
		{oidSHA256, pfxHashSHA256},
		{oidSHA384, pfxHashSHA384},
		{oidSHA512, pfxHashSHA512},
	} {
		if candidate.oid.Equal(oid) {
			h := candidate.hash
			return &h
		}
	}
	return nil
}
