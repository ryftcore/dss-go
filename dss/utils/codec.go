// Ported from dss-utils/.../Utils.java + IUtils.java (DSS 6.5.RC1),
// hex/base64 codec methods, matching org.apache.commons.codec.binary.Hex /
// Base64 semantics.

package utils

import (
	"encoding/base64"
	"encoding/hex"
)

// IsHexEncoded checks if the string is HEX (base16) encoded.
//
// Mirrors commons-codec Hex.decodeHex: requires an even number of valid
// hex digits (0-9, a-f, A-F); the empty string is valid HEX (decodes to
// an empty array).
func IsHexEncoded(hexString string) bool {
	_, err := hex.DecodeString(hexString)
	return err == nil
}

// ToHex transforms the binaries to a HEX-encoded (lowercase) string
// representation. A nil input yields "".
func ToHex(bytes []byte) string {
	return hex.EncodeToString(bytes)
}

// FromHex transforms a HEX-encoded string to a byte array. Returns an
// error if the string is not valid HEX (odd length, or contains
// non-hex-digit characters) — Java throws IllegalArgumentException in
// this case.
func FromHex(hexStr string) ([]byte, error) {
	b, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// base64Alphabet is the set of characters commons-codec's
// Base64.isBase64(byte) accepts as base64 data (standard alphabet + pad).
func isBase64Char(c byte) bool {
	switch {
	case c >= 'A' && c <= 'Z':
		return true
	case c >= 'a' && c <= 'z':
		return true
	case c >= '0' && c <= '9':
		return true
	case c == '+' || c == '/' || c == '=':
		return true
	default:
		return false
	}
}

// isBase64Whitespace mirrors commons-codec's notion of an "ignorable"
// whitespace byte within base64 text: space, \n, \r, \t.
func isBase64Whitespace(c byte) bool {
	switch c {
	case ' ', '\n', '\r', '\t':
		return true
	default:
		return false
	}
}

// IsBase64Encoded checks if the string is base64-encoded.
//
// Mirrors commons-codec Base64.isBase64(String): this is a character-set
// check only (every byte is either a base64 alphabet character, the pad
// character, or ignorable whitespace) — it does NOT validate padding
// correctness or that the length is decodable. The empty string is valid.
func IsBase64Encoded(base64String string) bool {
	for i := 0; i < len(base64String); i++ {
		c := base64String[i]
		if !isBase64Char(c) && !isBase64Whitespace(c) {
			return false
		}
	}
	return true
}

// ToBase64 transforms the binaries to a String base64-encoded
// representation (standard alphabet, padded, no line wrapping). A nil
// input yields "".
func ToBase64(bytes []byte) string {
	return base64.StdEncoding.EncodeToString(bytes)
}

// FromBase64 transforms a base64-encoded string to a byte array.
//
// Mirrors commons-codec Base64.decodeBase64: decoding is LENIENT. Bytes
// that are not part of the base64 alphabet (including whitespace) are
// silently discarded rather than raising an error, and a trailing group
// of characters too short to form a byte is dropped. Unlike Go's
// stdlib encoding/base64 (which is strict), this never errors — matching
// upstream's method signature, which declares no exception either.
func FromBase64(base64Str string) []byte {
	filtered := make([]byte, 0, len(base64Str))
	for i := 0; i < len(base64Str); i++ {
		c := base64Str[i]
		if c == '=' {
			continue
		}
		if (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '+' || c == '/' {
			filtered = append(filtered, c)
		}
	}
	// A trailing group of exactly one leftover char cannot encode any
	// bits into a full byte; drop it (2 leftover chars -> 1 byte, 3
	// leftover chars -> 2 bytes are both valid raw-base64 tails).
	if rem := len(filtered) % 4; rem == 1 {
		filtered = filtered[:len(filtered)-1]
	}
	decoded, err := base64.RawStdEncoding.DecodeString(string(filtered))
	if err != nil {
		// Should not happen given the filtering above, but stay
		// nil/empty-safe rather than panicking.
		return []byte{}
	}
	return decoded
}
