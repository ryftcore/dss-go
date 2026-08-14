// Ported from JDK javax.security.auth.x500.X500Principal and its sun.security.x509
// backing classes (X500Name, RDN, AVA) of JDK 21, as required by
// dss-model/src/main/java/eu/europa/esig/dss/model/x509/X500PrincipalHelper.java (DSS 6.5.RC1).
//
// dss-model exposes the CANONICAL and RFC2253 string forms of a distinguished name
// (X500PrincipalHelper#getCanonical / #getRFC2253 / #getPrettyPrintRFC2253) and those
// strings end up verbatim in DSS reports and in identity comparisons, so the JDK
// algorithms are reproduced here instead of being approximated with crypto/x509/pkix
// (whose pkix.Name.String() uses a different, incompatible syntax). The DER given to
// the constructor is kept byte-for-byte and is never re-encoded.
package model

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf16"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"golang.org/x/text/unicode/norm"
)

// DER tags used by an RDNSequence, mirroring sun.security.util.DerValue's tag constants.
const (
	x500PrincipalTagOID             byte = 0x06
	x500PrincipalTagSequence        byte = 0x30
	x500PrincipalTagSet             byte = 0x31
	x500PrincipalTagUTF8String      byte = 0x0C
	x500PrincipalTagPrintableString byte = 0x13
	x500PrincipalTagT61String       byte = 0x14
	x500PrincipalTagIA5String       byte = 0x16
	x500PrincipalTagGeneralString   byte = 0x1B
	x500PrincipalTagBMPString       byte = 0x1E
)

// x500PrincipalRFC2253Keywords maps an attribute type OID to the keyword AVAKeyword
// emits for AVA.RFC2253. Only the entries registered as RFC 2253 compliant are listed;
// AVAKeyword registers several keywords per OID and the *last* registration wins in its
// OID -> keyword map, which is why 2.5.4.8 resolves to "ST" (not "S"), 2.5.4.46 to the
// non-compliant "DNQ" and 1.2.840.113549.1.9.1 to the non-compliant "EMAILADDRESS".
// Every OID absent from this map is rendered as its dotted-decimal string.
var x500PrincipalRFC2253Keywords = map[string]string{
	"2.5.4.3":                    "CN",     // commonName
	"2.5.4.6":                    "C",      // countryName
	"2.5.4.7":                    "L",      // localityName
	"2.5.4.8":                    "ST",     // stateName
	"2.5.4.10":                   "O",      // orgName
	"2.5.4.11":                   "OU",     // orgUnitName
	"2.5.4.9":                    "STREET", // streetAddress
	"0.9.2342.19200300.100.1.25": "DC",     // domainComponent
	"0.9.2342.19200300.100.1.1":  "UID",    // userid
}

// x500PrincipalAVA is one AttributeTypeAndValue of a distinguished name.
type x500PrincipalAVA struct {
	// oid is the attribute type in dotted-decimal form.
	oid string
	// valueTag is the DER tag byte of the attribute value.
	valueTag byte
	// valueFull is the complete DER encoding (tag, length and content) of the value.
	valueFull []byte
	// valueData holds the content octets of the value.
	valueData []byte
}

// x500PrincipalRDN is one RelativeDistinguishedName, i.e. a SET OF AttributeTypeAndValue.
type x500PrincipalRDN struct {
	assertions []x500PrincipalAVA
}

// X500Principal represents an X.500 distinguished name, holding the DER encoding it was
// built from. It is the Go stand-in for javax.security.auth.x500.X500Principal.
type X500Principal struct {
	encoded []byte
	// names holds the RDNs in DER order (names[0] is the first RDN of the RDNSequence),
	// the same order sun.security.x509.X500Name stores them in.
	names []x500PrincipalRDN

	rfc2253      string
	rfc2253Valid bool

	canonical      string
	canonicalValid bool
}

// NewX500Principal parses the DER encoding of an X.501 Name (an RDNSequence) and returns
// the corresponding principal. The encoding is retained as given and is returned unchanged
// by Encoded. Port of the X500Principal(byte[] name) constructor; the Java constructor
// throws IllegalArgumentException on malformed input, which becomes an error here.
func NewX500Principal(encoded []byte) (*X500Principal, error) {
	if encoded == nil {
		panic("X500Principal encoding cannot be null")
	}
	names, err := x500PrincipalParse(encoded)
	if err != nil {
		return nil, err
	}
	return &X500Principal{encoded: encoded, names: names}, nil
}

// Encoded returns the distinguished name in ASN.1 DER encoded form, exactly as supplied to
// NewX500Principal. Port of X500Principal#getEncoded().
//
// Unlike the Java method this does not clone: the returned slice is the principal's own DER
// and must not be modified by the caller.
func (p *X500Principal) Encoded() []byte {
	return p.encoded
}

// RFC2253Name returns the RFC 2253 form of the name, emitting only the attribute type
// keywords RFC 2253 defines. Port of getName(X500Principal.RFC2253).
func (p *X500Principal) RFC2253Name() string {
	if !p.rfc2253Valid {
		// The empty OID map can never trigger the keyword validation error.
		name, _ := p.generateRFC2253DN(nil)
		p.rfc2253 = name
		p.rfc2253Valid = true
	}
	return p.rfc2253
}

// RFC2253NameWithOIDMap returns the RFC 2253 form of the name, additionally emitting the
// attribute type keywords given by oidMap (keys are dotted-decimal OIDs, values keywords).
// Entries of oidMap take precedence over the built-in keywords. Port of
// getName(X500Principal.RFC2253, Map).
//
// Java throws IllegalArgumentException when an OID present in the name maps to an
// improperly specified keyword; that becomes the returned error.
func (p *X500Principal) RFC2253NameWithOIDMap(oidMap map[string]string) (string, error) {
	if oidMap == nil {
		panic("provided null OID map")
	}
	if len(oidMap) == 0 {
		return p.RFC2253Name(), nil
	}
	return p.generateRFC2253DN(oidMap)
}

// Canonical returns the canonical form of the name, as specified for
// getName(X500Principal.CANONICAL): the RFC 2253 form restricted to PrintableString and
// UTF8String values, with multi-valued RDNs sorted, internal whitespace collapsed, leading
// and trailing whitespace removed, the whole string upper-cased then lower-cased with the
// US locale, and finally normalized to Unicode Normalization Form KD.
func (p *X500Principal) Canonical() string {
	if !p.canonicalValid {
		name, _ := p.generateRFC2253CanonicalDN()
		p.canonical = name
		p.canonicalValid = true
	}
	return p.canonical
}

// String returns the RFC 2253 form of the name.
//
// DEVIATION: the JDK's X500Principal#toString() emits the RFC 1779 flavoured "DEFAULT"
// format, which dss-model never consumes and which is therefore not ported. Callers that
// print a principal get the RFC 2253 form instead.
func (p *X500Principal) String() string {
	return p.RFC2253Name()
}

// Equals reports whether the two principals denote the same distinguished name, comparing
// their canonical forms. Port of X500Principal#equals(Object), which delegates to
// X500Name#equals and therefore compares RFC 2253 canonical names, NOT the DER encodings.
func (p *X500Principal) Equals(other *X500Principal) bool {
	if p == other {
		return true
	}
	if p == nil || other == nil {
		return false
	}
	// Quick check that the RDN/AVA counts match before canonicalizing, as X500Name#equals does.
	if len(p.names) != len(other.names) {
		return false
	}
	for i := range p.names {
		if len(p.names[i].assertions) != len(other.names[i].assertions) {
			return false
		}
	}
	return p.Canonical() == other.Canonical()
}

// generateRFC2253DN ports sun.security.x509.X500Name#generateRFC2253DN: the RDNs are
// emitted starting with the last element of the sequence and moving backwards, separated
// by a comma.
func (p *X500Principal) generateRFC2253DN(oidMap map[string]string) (string, error) {
	if len(p.names) == 0 {
		return "", nil
	}
	parts := make([]string, 0, len(p.names))
	for i := len(p.names) - 1; i >= 0; i-- {
		s, err := p.names[i].rfc2253String(false, oidMap)
		if err != nil {
			return "", err
		}
		parts = append(parts, s)
	}
	return strings.Join(parts, ","), nil
}

// generateRFC2253CanonicalDN ports sun.security.x509.X500Name#getRFC2253CanonicalName.
func (p *X500Principal) generateRFC2253CanonicalDN() (string, error) {
	if len(p.names) == 0 {
		return "", nil
	}
	parts := make([]string, 0, len(p.names))
	for i := len(p.names) - 1; i >= 0; i-- {
		s, err := p.names[i].rfc2253String(true, nil)
		if err != nil {
			return "", err
		}
		parts = append(parts, s)
	}
	return strings.Join(parts, ","), nil
}

// rfc2253String ports sun.security.x509.RDN#toRFC2253StringInternal. Adjoining AVAs of a
// multi-valued RDN are separated by '+'; in canonical mode they are first sorted with
// AVAComparator (AVAs carrying a standard RFC 2253 keyword come first, ordered
// alphabetically, followed by the OID-keyword AVAs ordered numerically).
func (r *x500PrincipalRDN) rfc2253String(canonical bool, oidMap map[string]string) (string, error) {
	if len(r.assertions) == 1 {
		if canonical {
			return r.assertions[0].rfc2253CanonicalString()
		}
		return r.assertions[0].rfc2253String(oidMap)
	}

	toOutput := r.assertions
	if canonical {
		toOutput = make([]x500PrincipalAVA, len(r.assertions))
		copy(toOutput, r.assertions)
		// Pre-compute the canonical strings the comparator sorts on; an error here would
		// also surface below, so failures are deferred rather than swallowed.
		keys := make([]string, len(toOutput))
		for i := range toOutput {
			keys[i], _ = toOutput[i].rfc2253CanonicalString()
		}
		idx := make([]int, len(toOutput))
		for i := range idx {
			idx[i] = i
		}
		sort.SliceStable(idx, func(a, b int) bool {
			ia, ib := idx[a], idx[b]
			aHas := toOutput[ia].hasRFC2253Keyword()
			bHas := toOutput[ib].hasRFC2253Keyword()
			if aHas == bHas {
				return x500PrincipalCompareJavaStrings(keys[ia], keys[ib]) < 0
			}
			return aHas
		})
		sorted := make([]x500PrincipalAVA, len(toOutput))
		for i, j := range idx {
			sorted[i] = toOutput[j]
		}
		toOutput = sorted
	}

	parts := make([]string, 0, len(toOutput))
	for i := range toOutput {
		var (
			s   string
			err error
		)
		if canonical {
			s, err = toOutput[i].rfc2253CanonicalString()
		} else {
			s, err = toOutput[i].rfc2253String(oidMap)
		}
		if err != nil {
			return "", err
		}
		parts = append(parts, s)
	}
	return strings.Join(parts, "+"), nil
}

// hasRFC2253Keyword ports AVA#hasRFC2253Keyword().
func (a *x500PrincipalAVA) hasRFC2253Keyword() bool {
	_, ok := x500PrincipalRFC2253Keywords[a.oid]
	return ok
}

// rfc2253String ports sun.security.x509.AVA#toRFC2253String(Map).
func (a *x500PrincipalAVA) rfc2253String(oidMap map[string]string) (string, error) {
	keyword, err := x500PrincipalKeyword(a.oid, oidMap)
	if err != nil {
		return "", err
	}
	var out strings.Builder
	out.WriteString(keyword)
	out.WriteByte('=')

	first := keyword[0]
	if (first >= '0' && first <= '9') || !x500PrincipalIsDerString(a.valueTag, false) {
		// Values without a string representation are emitted as '#' followed by the
		// lower-case hex of the whole BER encoding of the attribute value.
		out.WriteByte('#')
		out.WriteString(x500PrincipalHexLower(a.valueFull))
		return out.String(), nil
	}

	valStr, err := x500PrincipalDecodeValue(a.valueTag, a.valueData, false)
	if err != nil {
		return "", err
	}

	// NOTE: the JDK also escapes '=' and '#' anywhere in the value, and null as "\00".
	const escapees = ",=+<>#;\"\\"
	var escaped []rune
	for _, c := range valStr {
		switch {
		case x500PrincipalIsPrintableStringChar(c) || strings.ContainsRune(escapees, c):
			if strings.ContainsRune(escapees, c) {
				escaped = append(escaped, '\\')
			}
			escaped = append(escaped, c)
		case c == 0:
			escaped = append(escaped, '\\', '0', '0')
		default:
			escaped = append(escaped, c)
		}
	}

	// Escape leading and trailing whitespace, where "whitespace" is only ' ' and '\r'.
	lead := 0
	for ; lead < len(escaped); lead++ {
		if escaped[lead] != ' ' && escaped[lead] != '\r' {
			break
		}
	}
	trail := len(escaped) - 1
	for ; trail >= 0; trail-- {
		if escaped[trail] != ' ' && escaped[trail] != '\r' {
			break
		}
	}
	for i, c := range escaped {
		if i < lead || i > trail {
			out.WriteByte('\\')
		}
		out.WriteRune(c)
	}
	return out.String(), nil
}

// rfc2253CanonicalString ports sun.security.x509.AVA#toRFC2253CanonicalString().
func (a *x500PrincipalAVA) rfc2253CanonicalString() (string, error) {
	keyword, err := x500PrincipalKeyword(a.oid, nil)
	if err != nil {
		return "", err
	}
	var out strings.Builder
	out.WriteString(keyword)
	out.WriteByte('=')

	first := keyword[0]
	if (first >= '0' && first <= '9') || !x500PrincipalIsDerString(a.valueTag, true) {
		// In canonical form only PrintableString and UTF8String keep a string
		// representation; every other string type is hex-encoded.
		out.WriteByte('#')
		out.WriteString(x500PrincipalHexLower(a.valueFull))
	} else {
		valStr, err := x500PrincipalDecodeValue(a.valueTag, a.valueData, true)
		if err != nil {
			return "", err
		}

		const escapees = ",+<>;\"\\"
		var buf strings.Builder
		previousWhite := false
		for i, c := range []rune(valStr) {
			leadingHash := i == 0 && c == '#'
			if x500PrincipalIsPrintableStringChar(c) || strings.ContainsRune(escapees, c) || leadingHash {
				if leadingHash || strings.ContainsRune(escapees, c) {
					buf.WriteByte('\\')
				}
				// Collapse runs of whitespace into a single character.
				if !x500PrincipalIsJavaWhitespace(c) {
					previousWhite = false
					buf.WriteRune(c)
				} else if !previousWhite {
					previousWhite = true
					buf.WriteRune(c)
				}
			} else {
				previousWhite = false
				buf.WriteRune(c)
			}
		}
		out.WriteString(x500PrincipalJavaTrim(buf.String()))
	}

	canon := out.String()
	canon = cases.Lower(language.AmericanEnglish).String(cases.Upper(language.AmericanEnglish).String(canon))
	return norm.NFKD.String(canon), nil
}

// x500PrincipalKeyword ports AVAKeyword#getKeyword(oid, AVA.RFC2253, extraOidMap): the
// caller-supplied map wins over the built-in keywords, and an OID without a compliant
// keyword is rendered as its dotted-decimal string.
func x500PrincipalKeyword(oid string, extraOidMap map[string]string) (string, error) {
	if keyword, ok := extraOidMap[oid]; ok {
		if keyword == "" {
			return "", errors.New("keyword cannot be empty")
		}
		keyword = x500PrincipalJavaTrim(keyword)
		if keyword == "" {
			return "", errors.New("keyword cannot be empty")
		}
		for i, c := range []rune(keyword) {
			letter := !(c < 65 || c > 122 || (c > 90 && c < 97))
			if i == 0 {
				if !letter {
					return "", errors.New("keyword does not start with letter")
				}
				continue
			}
			digit := c >= 48 && c <= 57
			if !letter && !digit && c != '_' {
				return "", errors.New("keyword character is not a letter, digit, or underscore")
			}
		}
		return keyword, nil
	}
	if keyword, ok := x500PrincipalRFC2253Keywords[oid]; ok {
		return keyword, nil
	}
	return oid, nil
}

// x500PrincipalIsDerString ports AVA#isDerString(DerValue, boolean).
func x500PrincipalIsDerString(tag byte, canonical bool) bool {
	if canonical {
		return tag == x500PrincipalTagPrintableString || tag == x500PrincipalTagUTF8String
	}
	switch tag {
	case x500PrincipalTagPrintableString, x500PrincipalTagT61String, x500PrincipalTagIA5String,
		x500PrincipalTagGeneralString, x500PrincipalTagBMPString, x500PrincipalTagUTF8String:
		return true
	}
	return false
}

// x500PrincipalDecodeValue ports AVA#getCharset(DerValue, boolean) followed by the
// String construction: 8-bit types are decoded as ISO-8859-1 (not US-ASCII, so that
// non-compliant certificates keep their full byte range), BMPString as UTF-16BE and
// UTF8String as UTF-8.
func x500PrincipalDecodeValue(tag byte, data []byte, canonical bool) (string, error) {
	if canonical {
		switch tag {
		case x500PrincipalTagPrintableString:
			return x500PrincipalDecodeLatin1(data), nil
		case x500PrincipalTagUTF8String:
			return string(data), nil
		}
		return "", fmt.Errorf("unexpected tag: %d", int(tag))
	}
	switch tag {
	case x500PrincipalTagPrintableString, x500PrincipalTagT61String,
		x500PrincipalTagIA5String, x500PrincipalTagGeneralString:
		return x500PrincipalDecodeLatin1(data), nil
	case x500PrincipalTagBMPString:
		return x500PrincipalDecodeUTF16BE(data)
	case x500PrincipalTagUTF8String:
		return string(data), nil
	}
	return "", fmt.Errorf("unexpected tag: %d", int(tag))
}

// x500PrincipalDecodeLatin1 decodes ISO-8859-1 bytes, where each byte is one code point.
func x500PrincipalDecodeLatin1(data []byte) string {
	runes := make([]rune, len(data))
	for i, b := range data {
		runes[i] = rune(b)
	}
	return string(runes)
}

// x500PrincipalDecodeUTF16BE decodes big-endian UTF-16 (a DER BMPString).
func x500PrincipalDecodeUTF16BE(data []byte) (string, error) {
	if len(data)%2 != 0 {
		return "", errors.New("BMPString with an odd number of octets")
	}
	units := make([]uint16, len(data)/2)
	for i := range units {
		units[i] = uint16(data[2*i])<<8 | uint16(data[2*i+1])
	}
	return string(utf16.Decode(units)), nil
}

// x500PrincipalIsPrintableStringChar ports sun.security.util.DerValue#isPrintableStringChar.
func x500PrincipalIsPrintableStringChar(c rune) bool {
	if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') {
		return true
	}
	switch c {
	case ' ', '\'', '(', ')', '+', ',', '-', '.', '/', ':', '=', '?':
		return true
	}
	return false
}

// x500PrincipalIsJavaWhitespace ports java.lang.Character#isWhitespace(int): the Unicode
// space separators except the non-breaking ones, plus the ASCII control characters Java
// treats as whitespace. It deliberately differs from unicode.IsSpace, which also accepts
// U+00A0 and U+0085.
func x500PrincipalIsJavaWhitespace(c rune) bool {
	switch c {
	case '\t', '\n', 0x0B, '\f', '\r', 0x1C, 0x1D, 0x1E, 0x1F:
		return true
	case 0x00A0, 0x2007, 0x202F: // non-breaking spaces
		return false
	}
	return unicode.In(c, unicode.Zs, unicode.Zl, unicode.Zp)
}

// x500PrincipalJavaTrim ports java.lang.String#trim(), which strips every leading and
// trailing character whose code point is not greater than U+0020 - not the Unicode
// whitespace set strings.TrimSpace uses.
func x500PrincipalJavaTrim(s string) string {
	return strings.TrimFunc(s, func(r rune) bool { return r <= ' ' })
}

// x500PrincipalCompareJavaStrings ports java.lang.String#compareTo, which orders by UTF-16
// code unit and therefore differs from Go's byte-wise string comparison for code points
// above the BMP.
func x500PrincipalCompareJavaStrings(a, b string) int {
	ua := utf16.Encode([]rune(a))
	ub := utf16.Encode([]rune(b))
	n := len(ua)
	if len(ub) < n {
		n = len(ub)
	}
	for i := 0; i < n; i++ {
		if ua[i] != ub[i] {
			return int(ua[i]) - int(ub[i])
		}
	}
	return len(ua) - len(ub)
}

// x500PrincipalHexLower renders bytes the way java.util.HexFormat.of() does: lower case,
// two digits per byte, no separator.
func x500PrincipalHexLower(data []byte) string {
	const digits = "0123456789abcdef"
	var sb strings.Builder
	sb.Grow(2 * len(data))
	for _, b := range data {
		sb.WriteByte(digits[b>>4])
		sb.WriteByte(digits[b&0x0F])
	}
	return sb.String()
}

// x500PrincipalTLV is one parsed DER tag-length-value triple.
type x500PrincipalTLV struct {
	tag   byte
	full  []byte
	value []byte
}

// x500PrincipalReadTLV reads one DER element from the front of b and returns it together
// with the remaining bytes.
func x500PrincipalReadTLV(b []byte) (x500PrincipalTLV, []byte, error) {
	var tlv x500PrincipalTLV
	if len(b) < 2 {
		return tlv, nil, errors.New("truncated DER element")
	}
	tag := b[0]
	if tag&0x1F == 0x1F {
		return tlv, nil, errors.New("high tag number form is not supported in a distinguished name")
	}
	lengthByte := b[1]
	var (
		length     int
		headerSize int
	)
	if lengthByte&0x80 == 0 {
		length = int(lengthByte)
		headerSize = 2
	} else {
		numBytes := int(lengthByte & 0x7F)
		if numBytes == 0 {
			return tlv, nil, errors.New("indefinite length is not valid DER")
		}
		if numBytes > 4 || len(b) < 2+numBytes {
			return tlv, nil, errors.New("unsupported or truncated DER length")
		}
		for _, c := range b[2 : 2+numBytes] {
			length = length<<8 | int(c)
		}
		if length < 0 {
			return tlv, nil, errors.New("invalid DER length")
		}
		headerSize = 2 + numBytes
	}
	if len(b) < headerSize+length {
		return tlv, nil, errors.New("truncated DER element")
	}
	tlv.tag = tag
	tlv.full = b[:headerSize+length]
	tlv.value = b[headerSize : headerSize+length]
	return tlv, b[headerSize+length:], nil
}

// x500PrincipalParse decodes an RDNSequence, keeping the RDNs in DER order.
func x500PrincipalParse(encoded []byte) ([]x500PrincipalRDN, error) {
	seq, rest, err := x500PrincipalReadTLV(encoded)
	if err != nil {
		return nil, err
	}
	if seq.tag != x500PrincipalTagSequence {
		return nil, errors.New("X500Principal encoding is not a SEQUENCE")
	}
	if len(rest) != 0 {
		return nil, errors.New("trailing bytes after the X500Principal encoding")
	}

	var names []x500PrincipalRDN
	body := seq.value
	for len(body) > 0 {
		var set x500PrincipalTLV
		set, body, err = x500PrincipalReadTLV(body)
		if err != nil {
			return nil, err
		}
		if set.tag != x500PrincipalTagSet {
			return nil, errors.New("RelativeDistinguishedName is not a SET")
		}
		rdn := x500PrincipalRDN{}
		avaBody := set.value
		for len(avaBody) > 0 {
			var ava x500PrincipalTLV
			ava, avaBody, err = x500PrincipalReadTLV(avaBody)
			if err != nil {
				return nil, err
			}
			if ava.tag != x500PrincipalTagSequence {
				return nil, errors.New("AVA not a sequence")
			}
			oidTLV, avaRest, err := x500PrincipalReadTLV(ava.value)
			if err != nil {
				return nil, err
			}
			if oidTLV.tag != x500PrincipalTagOID {
				return nil, errors.New("AVA type is not an OBJECT IDENTIFIER")
			}
			oid, err := x500PrincipalDecodeOID(oidTLV.value)
			if err != nil {
				return nil, err
			}
			valueTLV, avaRest, err := x500PrincipalReadTLV(avaRest)
			if err != nil {
				return nil, err
			}
			if len(avaRest) != 0 {
				return nil, fmt.Errorf("AVA, extra bytes = %d", len(avaRest))
			}
			if valueTLV.tag == x500PrincipalTagBMPString && len(valueTLV.value)%2 != 0 {
				return nil, errors.New("BMPString with an odd number of octets")
			}
			rdn.assertions = append(rdn.assertions, x500PrincipalAVA{
				oid:       oid,
				valueTag:  valueTLV.tag,
				valueFull: valueTLV.full,
				valueData: valueTLV.value,
			})
		}
		if len(rdn.assertions) == 0 {
			return nil, errors.New("RelativeDistinguishedName is empty")
		}
		names = append(names, rdn)
	}
	return names, nil
}

// x500PrincipalDecodeOID renders the content octets of an OBJECT IDENTIFIER as the
// dotted-decimal string ObjectIdentifier#toString() produces.
func x500PrincipalDecodeOID(data []byte) (string, error) {
	if len(data) == 0 {
		return "", errors.New("empty OBJECT IDENTIFIER")
	}
	var (
		sb    strings.Builder
		value uint64
		start = true
		any   bool
	)
	for i, b := range data {
		if value > (1<<57)-1 {
			return "", errors.New("OBJECT IDENTIFIER component overflow")
		}
		value = value<<7 | uint64(b&0x7F)
		if b&0x80 != 0 {
			if i == len(data)-1 {
				return "", errors.New("truncated OBJECT IDENTIFIER")
			}
			continue
		}
		if start {
			first := value / 40
			if first > 2 {
				first = 2
			}
			sb.WriteString(strconv.FormatUint(first, 10))
			sb.WriteByte('.')
			sb.WriteString(strconv.FormatUint(value-first*40, 10))
			start = false
		} else {
			sb.WriteByte('.')
			sb.WriteString(strconv.FormatUint(value, 10))
		}
		value = 0
		any = true
	}
	if !any {
		return "", errors.New("truncated OBJECT IDENTIFIER")
	}
	return sb.String(), nil
}
