package asn1ber

import (
	"slices"
	"strings"
)

// ValueToString ports org.bouncycastle.asn1.x500.style.IETFUtils#valueToString: it renders an
// X.500 attribute value the way RFC 4514 wants it, escaping the specials and hash-encoding
// everything that is not a string type.
func ValueToString(element *Element) string {
	var buffer []rune
	// IETFUtils#valueToString takes the string branch for every ASN1String EXCEPT
	// ASN1UniversalString, which it explicitly excludes so that it is hash-encoded like a
	// non-string value.
	if element.IsASN1String() && !element.IsUniversal(TagUniversalString) {
		value := element.AsString()
		if len(value) > 0 && value[0] == '#' {
			buffer = append(buffer, '\\')
		}
		buffer = append(buffer, []rune(value)...)
	} else {
		buffer = append(buffer, '#')
		buffer = append(buffer, []rune(hexLower(element.DEREncoded()))...)
	}

	// The escaping below produces exactly what IETFUtils' in-place StringBuffer#insert loops
	// do, but in linear time: inserting one rune at a time made a value of a few hundred
	// thousand specials or spaces cost minutes.
	index := 0
	if len(buffer) >= 2 && buffer[0] == '\\' && buffer[1] == '#' {
		index += 2
	}
	var escaped []rune
	// Most values (every plain name) have nothing to escape: copy only when one does.
	if slices.ContainsFunc(buffer[index:], isRFC4514Special) {
		escaped = make([]rune, 0, len(buffer)+len(buffer)/4+4)
		escaped = append(escaped, buffer[:index]...)
		for _, r := range buffer[index:] {
			if isRFC4514Special(r) {
				escaped = append(escaped, '\\')
			}
			escaped = append(escaped, r)
		}
		buffer = escaped
	}

	// Every leading space is escaped, then every space of the trailing run of the result.
	leading := 0
	for leading < len(buffer) && buffer[leading] == ' ' {
		leading++
	}
	if leading > 0 {
		escaped = make([]rune, 0, len(buffer)+leading)
		for range leading {
			escaped = append(escaped, '\\', ' ')
		}
		buffer = append(escaped, buffer[leading:]...)
	}
	trailing := 0
	for trailing < len(buffer) && buffer[len(buffer)-1-trailing] == ' ' {
		trailing++
	}
	if trailing > 0 {
		escaped = append(make([]rune, 0, len(buffer)+trailing), buffer[:len(buffer)-trailing]...)
		for range trailing {
			escaped = append(escaped, '\\', ' ')
		}
		buffer = escaped
	}
	return string(buffer)
}

// isRFC4514Special reports whether IETFUtils#valueToString backslash-escapes r.
func isRFC4514Special(r rune) bool {
	switch r {
	case ',', '"', '\\', '+', '=', '<', '>', ';':
		return true
	}
	return false
}

// ASN1ToString reproduces ASN1Primitive#toString for the value types an X.500 attribute can
// carry: a string type yields its text (a BIT STRING and a UniversalString their
// "#"+UPPER-case-hex form, see Element.AsString), an OBJECT IDENTIFIER its dotted form,
// anything else "#" followed by the hex of its DER encoding.
//
// DEVIATION: BouncyCastle renders a constructed value (an attribute whose value is a SEQUENCE
// or a SET, which no X.520 attribute type defines) as its ASN1Dump-style "[a, b]" listing;
// this port hash-encodes it like any other non-string value.
func ASN1ToString(element *Element) string {
	if element.IsASN1String() {
		return element.AsString()
	}
	switch {
	case element.IsUniversal(TagOID):
		if oid, err := element.ObjectIdentifier(); err == nil {
			return oid.String()
		}
	case element.IsUniversal(TagInteger):
		return element.Integer().String()
	case element.IsUniversal(TagOctetString):
		// ASN1OctetString#toString hexes the CONTENT, not the whole encoding.
		return "#" + hexLower(element.Octets())
	case element.IsUniversal(TagBoolean):
		if len(element.content) > 0 && element.content[0] != 0x00 {
			return "TRUE"
		}
		return "FALSE"
	case element.IsUniversal(TagNull):
		return "NULL"
	}
	return "#" + hexLower(element.DEREncoded())
}

// JavaTrim ports java.lang.String#trim, which strips every leading and trailing character
// whose code point is not greater than U+0020.
func JavaTrim(value string) string {
	return strings.TrimFunc(value, func(r rune) bool { return r <= ' ' })
}

// hexUpper renders bytes as upper-case hex without a separator, the way the private hex table
// of ASN1BitString#getString / ASN1UniversalString#getString does.
func hexUpper(data []byte) string {
	const digits = "0123456789ABCDEF"
	var builder strings.Builder
	builder.Grow(2 * len(data))
	for _, b := range data {
		builder.WriteByte(digits[b>>4])
		builder.WriteByte(digits[b&0x0F])
	}
	return builder.String()
}

// hexLower renders bytes as lower-case hex without a separator, the way
// org.bouncycastle.util.encoders.Hex#encode does.
func hexLower(data []byte) string {
	const digits = "0123456789abcdef"
	var builder strings.Builder
	builder.Grow(2 * len(data))
	for _, b := range data {
		builder.WriteByte(digits[b>>4])
		builder.WriteByte(digits[b&0x0F])
	}
	return builder.String()
}
