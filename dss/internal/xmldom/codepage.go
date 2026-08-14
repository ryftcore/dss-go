package xmldom

import "unicode/utf8"

// ISO-8859-2 (Latin-2), the one single-byte code page beyond Latin-1 that the upstream
// XAdES corpus actually declares - validation/Signature-X-HU_MIC-1.xml and
// validation/BaselineBWithCertificateValues.xml both do, and Xerces reads them, so a port
// that rejected the encoding could not validate a Central-European signature at all.
//
// Bytes 0x00-0x9F are identical to Unicode; only the top half differs, so one 96-entry
// table covers the charset in both directions. Every one of the 96 is defined - Latin-2
// has no holes - which is why decoding cannot fail and why the reverse map is exact.
var latin2High = [96]rune{
	0x00A0, 0x0104, 0x02D8, 0x0141, 0x00A4, 0x013D, 0x015A, 0x00A7, // 0xA0
	0x00A8, 0x0160, 0x015E, 0x0164, 0x0179, 0x00AD, 0x017D, 0x017B, // 0xA8
	0x00B0, 0x0105, 0x02DB, 0x0142, 0x00B4, 0x013E, 0x015B, 0x02C7, // 0xB0
	0x00B8, 0x0161, 0x015F, 0x0165, 0x017A, 0x02DD, 0x017E, 0x017C, // 0xB8
	0x0154, 0x00C1, 0x00C2, 0x0102, 0x00C4, 0x0139, 0x0106, 0x00C7, // 0xC0
	0x010C, 0x00C9, 0x0118, 0x00CB, 0x011A, 0x00CD, 0x00CE, 0x010E, // 0xC8
	0x0110, 0x0143, 0x0147, 0x00D3, 0x00D4, 0x0150, 0x00D6, 0x00D7, // 0xD0
	0x0158, 0x016E, 0x00DA, 0x0170, 0x00DC, 0x00DD, 0x0162, 0x00DF, // 0xD8
	0x0155, 0x00E1, 0x00E2, 0x0103, 0x00E4, 0x013A, 0x0107, 0x00E7, // 0xE0
	0x010D, 0x00E9, 0x0119, 0x00EB, 0x011B, 0x00ED, 0x00EE, 0x010F, // 0xE8
	0x0111, 0x0144, 0x0148, 0x00F3, 0x00F4, 0x0151, 0x00F6, 0x00F7, // 0xF0
	0x0159, 0x016F, 0x00FA, 0x0171, 0x00FC, 0x00FD, 0x0163, 0x02D9, // 0xF8
}

// latin2Encode is latin2High inverted, built once.
var latin2Encode = func() map[rune]byte {
	m := make(map[rune]byte, 96)
	for i, r := range latin2High {
		m[r] = byte(0xA0 + i)
	}
	return m
}()

// decodeLatin2 transcodes ISO-8859-2 bytes to UTF-8.
func decodeLatin2(src []byte) []byte {
	out := make([]byte, 0, len(src)+len(src)/2)
	for _, b := range src {
		if b < 0xA0 {
			out = utf8.AppendRune(out, rune(b))
			continue
		}
		out = utf8.AppendRune(out, latin2High[b-0xA0])
	}
	return out
}

// latin2Byte reports the ISO-8859-2 byte for r, if the charset has one.
func latin2Byte(r rune) (byte, bool) {
	if r < 0 {
		return 0, false
	}
	if r < 0xA0 {
		return byte(r), true
	}
	b, ok := latin2Encode[r]
	return b, ok
}

// writeLatin2 is the OutputStreamWriter path for ISO-8859-2: a character the charset
// cannot carry becomes '?', the CharsetEncoder's default REPLACE action. Only unescaped
// writes - names, comment data, processing-instruction data - reach it, since text and
// attribute values turn the same characters into decimal character references first.
func writeLatin2(dst []byte, r rune) []byte {
	if b, ok := latin2Byte(r); ok {
		return append(dst, b)
	}
	return append(dst, '?')
}

// latin2InEncoding answers EncodingInfo.isInEncoding for one UTF-16 code unit: true when
// Latin-2 has a byte for it, which is never 0 and never '?' except for '?' itself.
func latin2InEncoding(u uint16) bool {
	_, ok := latin2Byte(rune(u))
	return ok
}
