package xmldom

import (
	"fmt"
	"strings"
	"unicode/utf16"
	"unicode/utf8"
)

// The output-encoding half of Serialize.
//
// DomUtils.serializeNode copies Document.getXmlEncoding() onto OutputKeys.ENCODING, and
// that single property decides three things at once in Xalan's ToStream:
//
//  1. the text of the XML declaration - Encodings.getMimeEncoding, which for a non-null
//     encoding is convertJava2MimeEncoding: a lookup of the UPPERCASED name in the
//     "java name" table of Encodings.properties, returning that entry's MIME name and the
//     argument unchanged when there is no such entry;
//  2. which characters are written literally and which become decimal character
//     references - EncodingInfo.isInEncoding;
//  3. the charset the bytes are finally written in - Encodings.getWriter.
//
// isInEncoding is not "can the charset represent this character". It encodes the single
// UTF-16 CODE UNIT with the charset and applies EncodingInfo.inEncoding(char, byte[]):
// false when the byte array is empty, false when its FIRST byte is 0, false when its
// first byte is '?' and the character is not '?', true otherwise. That first-byte test is
// why UTF-16BE reports U+00E9 as not-in-encoding (its bytes are 00 E9) while reporting
// U+4E2D as in-encoding (4E 2D), and why UTF-8 reports every high surrogate as
// not-in-encoding (a lone surrogate encodes to '?') and so escapes every supplementary
// character in text and attribute values. No reimplementation would arrive at these rules
// by reasoning; the enc-NN rows of xml/utils/testdata/serialize/goldens.txt pin each one.
//
// Only the encodings decodeSource accepts can reach here, since a declaration naming
// anything else fails at Parse.
type outEncoding struct {
	decl  string // what the XML declaration says
	write func(dst []byte, r rune) []byte
	bom   []byte

	// inEncoding answers isInEncoding for one UTF-16 code unit. Units below 127 never
	// reach it: escapingNotNeeded short-circuits them.
	inEncoding func(u uint16) bool

	// pairInEncoding answers the isInEncoding(high, low) overload, which is consulted
	// only from inside a CDATA section - the one place a supplementary character can be
	// written literally.
	pairInEncoding func(r rune) bool
}

// canonOut folds an encoding name the way Encodings' lookups do, case-insensitively.
func canonOut(name string) string { return strings.ToUpper(strings.TrimSpace(name)) }

func isSurrogateUnit(u uint16) bool { return u >= 0xD800 && u <= 0xDFFF }

// utf8InEncoding: every code unit encodes to a non-zero, non-'?' first byte except a lone
// surrogate, which the encoder replaces with '?'.
func utf8InEncoding(u uint16) bool { return !isSurrogateUnit(u) }

// asciiInEncoding: DEL is the one unit at or above 127 the charset can carry (the
// escapingNotNeeded fast path stops one short of it), and it comes back as the single byte
// 0x7F rather than as '?'. Everything above becomes '?' and so reports not-in-encoding.
func asciiInEncoding(u uint16) bool { return u == 0x7F }

// latin1InEncoding: one byte per unit up to U+00FF; above that the charset encodes to '?'.
func latin1InEncoding(u uint16) bool { return u <= 0xFF }

// utf16beInEncoding: the high byte comes first, so every unit below U+0100 leads with a
// zero byte and is reported as not-in-encoding.
func utf16beInEncoding(u uint16) bool { return u >= 0x100 }

// utf16leInEncoding: the low byte comes first, so a unit whose low byte is zero - U+0100,
// U+0200, ... - is reported as not-in-encoding.
func utf16leInEncoding(u uint16) bool { return byte(u) != 0 }

func alwaysInEncodingUnit(uint16) bool { return true }

func alwaysPairInEncoding(rune) bool { return true }

func neverPairInEncoding(rune) bool { return false }

// utf16bePairInEncoding: the pair's first byte is the high byte of the high surrogate,
// always 0xD8-0xDB.
func utf16bePairInEncoding(rune) bool { return true }

// utf16lePairInEncoding: the pair's first byte is the low byte of the high surrogate.
func utf16lePairInEncoding(r rune) bool {
	h, _ := utf16.EncodeRune(r)
	return byte(h) != 0
}

func writeUTF8(dst []byte, r rune) []byte { return utf8.AppendRune(dst, r) }

// writeXalanUTF8 is WriterToUTF8Buffered, the hand-rolled UTF-8 encoder ToStream installs
// when OutputKeys.ENCODING is spelled exactly "UTF-8" (case-insensitively) - which is both
// the default and what a document declaring UTF-8 selects. It is NOT the JDK's UTF-8
// encoder, and it disagrees with it: for a surrogate pair it computes the leading byte as
//
//	(byte) (0xF0 | (((high + 0x40) >> 8) & 0xf0))
//
// and 0xF0 | anything-masked-with-0xF0 is always 0xF0, so the three high bits of the code
// point are dropped. Code points from U+40000 up are therefore written with a wrong first
// byte - U+10FFFF comes out F0 8F BF BF instead of F4 8F BF BF, which is not even valid
// UTF-8. Reproduced here because these bytes are digested: getNodeBytes then runs them
// through new String(...)/getBytes() and each malformed byte becomes U+FFFD.
//
// Any other spelling of the encoding - "UTF8", "utf8" - takes the OutputStreamWriter path
// instead and gets the JDK's correct encoder. The astral-raw-* rows of
// xml/utils/testdata/serialize/goldens.txt pin both, and the boundary between them.
func writeXalanUTF8(dst []byte, r rune) []byte {
	if r <= 0xFFFF {
		return utf8.AppendRune(dst, r)
	}
	h, l := utf16.EncodeRune(r)
	high, low := int32(h), int32(l)
	return append(dst,
		byte(0xF0|(((high+0x40)>>8)&0xF0)),
		byte(0x80|(((high+0x40)>>2)&0x3F)),
		byte(0x80|((low>>6)&0x0F)+((high<<4)&0x30)),
		byte(0x80|(low&0x3F)),
	)
}

// writeLatin1 and writeASCII substitute '?' for a character outside the charset, which is
// what an OutputStreamWriter's CharsetEncoder does with its default REPLACE action. Only
// unescaped writes - names, comment data, processing-instruction data - can reach it,
// because text and attribute values turn the same characters into character references
// first.
func writeLatin1(dst []byte, r rune) []byte {
	if r > 0xFF || r < 0 {
		return append(dst, '?')
	}
	return append(dst, byte(r))
}

func writeASCII(dst []byte, r rune) []byte {
	if r > 0x7F || r < 0 {
		return append(dst, '?')
	}
	return append(dst, byte(r))
}

// writeTruncatingASCII is WriterToASCI, which ToStream installs when OutputKeys.ENCODING
// is EXACTLY "US-ASCII" or "ASCII" - a case-SENSITIVE String.equals, so a document
// declaring "us-ascii" gets the ordinary OutputStreamWriter and its '?' substitution
// instead. WriterToASCI does neither: it calls OutputStream.write(int) per character,
// which keeps only the low eight bits. U+00E9 comes out as the byte E9, U+4E2D as '-',
// and a supplementary character as the low bytes of its two surrogates. Only the
// raw-write path can reach it, since text and attribute values escape everything above
// 126 first; the graft-NN rows of xml/utils/testdata/serialize/goldens.txt pin it.
func writeTruncatingASCII(dst []byte, r rune) []byte {
	if r <= 0xFFFF {
		return append(dst, byte(r))
	}
	h, l := utf16.EncodeRune(r)
	return append(dst, byte(h), byte(l))
}

func writeUTF16BE(dst []byte, r rune) []byte {
	for _, u := range utf16.Encode([]rune{r}) {
		dst = append(dst, byte(u>>8), byte(u))
	}
	return dst
}

func writeUTF16LE(dst []byte, r rune) []byte {
	for _, u := range utf16.Encode([]rune{r}) {
		dst = append(dst, byte(u), byte(u>>8))
	}
	return dst
}

// javaMimeNames is the part of Encodings.properties' java-name table that the spellings
// decodeSource accepts can hit. Only these are rewritten in the declaration; every other
// spelling - "utf-8", "us-ascii", "latin1", "cp819", "UTF-16LE", ... - is absent from the
// java-name table and is echoed exactly as the document declared it.
var javaMimeNames = map[string]string{
	"UTF8":       "UTF-8",
	"ISO-8859-1": "ISO-8859-1",
	"ISO8859-1":  "ISO-8859-1",
	"ISO8859_1":  "ISO-8859-1",
	"ISO-8859-2": "ISO-8859-2",
	"ISO8859-2":  "ISO-8859-2",
	"ISO8859_2":  "ISO-8859-2",
}

// lookupOutEncoding resolves a declared encoding name to its output behaviour.
//
// A name Java's Charset lookup does not know makes the identity Transformer throw
// UnsupportedEncodingException, which DomUtils wraps in a DSSException. "latin-1" - which
// decodeSource accepts on input and Java has no charset alias for - is the one such name
// reachable here, and it fails on both sides.
func lookupOutEncoding(name string) (outEncoding, error) {
	decl := name
	if mime, ok := javaMimeNames[canonOut(name)]; ok {
		decl = mime
	}
	switch canonOut(name) {
	// Spelled "UTF-8" the serializer uses its own encoder, bug and all; spelled "UTF8" it
	// uses the JDK's. See writeXalanUTF8.
	case "UTF-8":
		return outEncoding{decl: decl, write: writeXalanUTF8,
			inEncoding: utf8InEncoding, pairInEncoding: alwaysPairInEncoding}, nil

	case "UTF8":
		return outEncoding{decl: decl, write: writeUTF8,
			inEncoding: utf8InEncoding, pairInEncoding: alwaysPairInEncoding}, nil

	case "US-ASCII", "ASCII", "ANSI_X3.4-1968", "ISO646-US", "ISO-IR-6":
		write := writeASCII
		if name == "US-ASCII" || name == "ASCII" {
			write = writeTruncatingASCII
		}
		return outEncoding{decl: decl, write: write,
			inEncoding: asciiInEncoding, pairInEncoding: neverPairInEncoding}, nil

	case "ISO-8859-1", "ISO8859-1", "ISO8859_1", "ISO_8859-1", "LATIN1", "L1", "CP819", "IBM819", "ISO-IR-100":
		return outEncoding{decl: decl, write: writeLatin1,
			inEncoding: latin1InEncoding, pairInEncoding: neverPairInEncoding}, nil

	case "ISO-8859-2", "ISO8859-2", "ISO8859_2", "ISO_8859-2", "ISO_8859-2:1987",
		"LATIN2", "L2", "CSISOLATIN2", "ISO-IR-101", "CP912", "IBM912", "IBM-912", "912":
		return outEncoding{decl: decl, write: writeLatin2,
			inEncoding: latin2InEncoding, pairInEncoding: neverPairInEncoding}, nil

	case "UTF-16BE":
		return outEncoding{decl: decl, write: writeUTF16BE,
			inEncoding: utf16beInEncoding, pairInEncoding: utf16bePairInEncoding}, nil

	case "UTF-16LE":
		return outEncoding{decl: decl, write: writeUTF16LE,
			inEncoding: utf16leInEncoding, pairInEncoding: utf16lePairInEncoding}, nil

	// getBytes("UTF-16") emits a big-endian byte-order mark first, which is both why the
	// output starts with FE FF and why every character reports in-encoding: the first
	// byte examined is the mark's.
	case "UTF-16":
		return outEncoding{decl: decl, write: writeUTF16BE, bom: []byte{0xFE, 0xFF},
			inEncoding: alwaysInEncodingUnit, pairInEncoding: alwaysPairInEncoding}, nil
	}
	return outEncoding{}, fmt.Errorf("unsupported output encoding %q", name)
}
