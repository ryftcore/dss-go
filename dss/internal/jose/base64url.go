// Ported from org.jose4j.base64url.Base64Url, the bundled
// org.jose4j.base64url.internal.apache.commons.codec.binary.Base64 it delegates to, and
// org.jose4j.lang.StringUtil (jose4j 0.9.6).
package jose

import (
	"strings"
	"unicode/utf8"
)

// base64URLEncodeTable is jose4j's URL_SAFE_ENCODE_TABLE: RFC 4648 section 5, '-' and '_' in
// place of '+' and '/'. DSSJsonUtils copies this table verbatim to implement isBase64UrlEncoded,
// so the two must not drift.
const base64URLEncodeTable = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"

// base64DecodeTable is commons-codec's DECODE_TABLE. It deliberately accepts BOTH alphabets at
// once - '+' and '-' both decode to 62, '/' and '_' both to 63 - because the decoder is shared
// between standard and URL-safe base64. -1 means "not an alphabet character", which the decoder
// silently SKIPS rather than rejecting.
var base64DecodeTable = func() [256]int8 {
	var t [256]int8
	for i := range t {
		t[i] = -1
	}
	for i := 0; i < 26; i++ {
		t['A'+i] = int8(i)
		t['a'+i] = int8(26 + i)
	}
	for i := 0; i < 10; i++ {
		t['0'+i] = int8(52 + i)
	}
	t['+'] = 62
	t['-'] = 62
	t['/'] = 63
	t['_'] = 63
	return t
}()

// Base64URLEncode returns the unpadded base64url encoding of b (RFC 7515 section 2). Port of
// Base64Url.encode(byte[]), which uses new Base64(-1, null, true): line length -1 means no
// chunking, no line separator, and urlSafe=true means the '-'/'_' alphabet with padding omitted.
//
// The encoder is written out rather than delegated to encoding/base64.RawURLEncoding only so
// that this file states the whole contract in one place; the two agree byte for byte, which
// TestBase64URLEncodeMatchesStdlib checks.
func Base64URLEncode(b []byte) string {
	if len(b) == 0 {
		return ""
	}
	var sb strings.Builder
	sb.Grow((len(b)*8 + 5) / 6)
	i := 0
	for ; i+3 <= len(b); i += 3 {
		v := uint32(b[i])<<16 | uint32(b[i+1])<<8 | uint32(b[i+2])
		sb.WriteByte(base64URLEncodeTable[(v>>18)&0x3f])
		sb.WriteByte(base64URLEncodeTable[(v>>12)&0x3f])
		sb.WriteByte(base64URLEncodeTable[(v>>6)&0x3f])
		sb.WriteByte(base64URLEncodeTable[v&0x3f])
	}
	switch len(b) - i {
	case 1:
		v := uint32(b[i])
		sb.WriteByte(base64URLEncodeTable[(v>>2)&0x3f])
		sb.WriteByte(base64URLEncodeTable[(v<<4)&0x3f])
	case 2:
		v := uint32(b[i])<<8 | uint32(b[i+1])
		sb.WriteByte(base64URLEncodeTable[(v>>10)&0x3f])
		sb.WriteByte(base64URLEncodeTable[(v>>4)&0x3f])
		sb.WriteByte(base64URLEncodeTable[(v<<2)&0x3f])
	}
	return sb.String()
}

// Base64URLDecode decodes a base64url (or plain base64) string. Port of Base64Url.decode(String)
// as commons-codec's Base64.decode implements it, and it is LENIENT in three ways that matter:
//
//   - characters outside the alphabet are silently skipped, not rejected: "Zm 9v", "Zm9v." and
//     "Zm9v\n" all decode to "foo", and "!!!!" decodes to the empty slice;
//   - both alphabets are accepted, so "+/8=" and "-_8" decode to the same two bytes;
//   - padding is optional, and the first '=' ends the input; leftover bits that do not fill a
//     whole byte are discarded, so the 1-character input "A" decodes to nothing at all.
//
// The leniency is not incidental. DSSJsonUtils.isBase64UrlEncoded calls this and then separately
// checks that every byte of the input is in the URL-safe alphabet, precisely because the decoder
// alone would accept far too much; a strict decoder here would change which strings DSS treats
// as base64url-encoded 'etsiU' components. Hence encoding/base64.RawURLEncoding is not usable.
//
// There is no error return because commons-codec has no failure mode on a byte array: every
// input decodes to something.
func Base64URLDecode(encoded string) []byte {
	if encoded == "" {
		return []byte{}
	}
	out := make([]byte, 0, len(encoded)*3/4+3)
	var work uint32
	modulus := 0
	for i := 0; i < len(encoded); i++ {
		c := encoded[i]
		if c == '=' {
			// The first pad character ends the stream ("We're done." in commons-codec).
			break
		}
		v := base64DecodeTable[c]
		if v < 0 {
			// Not an alphabet character: skipped entirely, not an error.
			continue
		}
		modulus = (modulus + 1) % 4
		work = work<<6 | uint32(v)
		if modulus == 0 {
			out = append(out, byte(work>>16), byte(work>>8), byte(work))
		}
	}
	// Spare bits: emit only whole bytes and drop the remainder.
	switch modulus {
	case 1:
		// 6 bits, not even one byte: ignored entirely.
	case 2:
		out = append(out, byte(work>>4))
	case 3:
		out = append(out, byte(work>>10), byte(work>>2))
	}
	return out
}

// Base64URLEncodeUTF8 returns the base64url encoding of the UTF-8 bytes of s. Port of
// Base64Url.encodeUtf8ByteRepresentation(String), which is how a serialized protected header
// becomes the first part of a JWS.
func Base64URLEncodeUTF8(s string) string {
	return Base64URLEncode([]byte(s))
}

// Base64URLDecodeToUTF8String decodes encoded and interprets the result as UTF-8. Port of
// Base64Url.decodeToUtf8String(String).
//
// Java's new String(bytes, UTF_8) replaces every malformed sequence with U+FFFD rather than
// failing, and Go's []byte->string conversion keeps the invalid bytes as-is, so the two disagree
// on malformed input. ReplaceInvalidUTF8 makes the Go side match; this matters because
// Headers.setEncodedHeader parses the result as JSON and a JAdES fixture may carry a corrupt
// header that must be rejected by the parser in the same way in both implementations.
func Base64URLDecodeToUTF8String(encoded string) string {
	return UTF8String(Base64URLDecode(encoded))
}

// UTF8String is Java's new String(bytes, StandardCharsets.UTF_8) (org.jose4j.lang.StringUtil
// .newStringUtf8): every byte that is not part of a well-formed UTF-8 sequence becomes the
// replacement character U+FFFD.
func UTF8String(b []byte) string {
	if utf8.Valid(b) {
		return string(b)
	}
	var sb strings.Builder
	sb.Grow(len(b))
	for i := 0; i < len(b); {
		r, size := utf8.DecodeRune(b[i:])
		if r == utf8.RuneError && size == 1 {
			sb.WriteRune(utf8.RuneError)
			i++
			continue
		}
		sb.Write(b[i : i+size])
		i += size
	}
	return sb.String()
}

// ASCIIBytes is Java's String.getBytes(US_ASCII) (org.jose4j.lang.StringUtil.getBytesAscii),
// which DSSJsonUtils.getAsciiBytes exposes and which produces the JWS signing input.
//
// The unrepresentable-character rule is load-bearing: Java's US-ASCII encoder replaces every
// code point above U+007F with a single '?' (0x3F) instead of failing. The signing input is
// built from base64url text, which is pure ASCII, so the branch should never fire in practice -
// but if a caller ever hands it non-ASCII, the two implementations must still produce the same
// bytes, because those bytes are what gets signed.
func ASCIIBytes(s string) []byte {
	out := make([]byte, 0, len(s))
	for _, r := range s {
		if r > 0x7f {
			out = append(out, '?')
			continue
		}
		out = append(out, byte(r))
	}
	return out
}

// IsBase64URLCharacter reports whether b is one of the 64 characters of the URL-safe alphabet.
// Port of the per-byte half of DSSJsonUtils.isBase64UrlEncoded(byte), kept here so the alphabet
// has exactly one definition.
func IsBase64URLCharacter(b byte) bool {
	return strings.IndexByte(base64URLEncodeTable, b) >= 0
}
