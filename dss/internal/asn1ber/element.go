package asn1ber

import (
	"encoding/asn1"
	"errors"
	"math/big"
	"unicode/utf16"
)

// Element is one parsed ASN.1 element. It replaces BouncyCastle's ASN1Primitive hierarchy
// with exactly what the ported code needs: the identifier, the content and, for a constructed
// element, its children.
//
// An Element is produced by Parse and is read-only: its accessors hand out the parser's own
// slices, which callers must not modify.
type Element struct {
	// class is the identifier's class bits: 0x00 universal, 0x40 application,
	// 0x80 context-specific, 0xC0 private.
	class byte
	// constructed reports whether the element is constructed.
	constructed bool
	// tagNumber is the tag number, decoded from the high-tag-number form when needed.
	tagNumber uint64
	// indefinite reports whether the element used the indefinite-length form.
	indefinite bool
	// content holds the content octets of a primitive element.
	content []byte
	// children holds the components of a constructed element.
	children []*Element
	// encoded is the element's original encoding, tag and length included.
	encoded []byte
}

// Class returns the identifier's class bits: ClassUniversal, ClassApplication,
// ClassContextSpecific or ClassPrivate.
func (e *Element) Class() byte { return e.class }

// TagNumber returns the element's tag number.
func (e *Element) TagNumber() uint64 { return e.tagNumber }

// IsConstructed reports whether the element is constructed.
func (e *Element) IsConstructed() bool { return e.constructed }

// IsIndefinite reports whether the element used the indefinite-length form.
func (e *Element) IsIndefinite() bool { return e.indefinite }

// Content returns the content octets of a primitive element, nil for a constructed one.
func (e *Element) Content() []byte { return e.content }

// Children returns the components of a constructed element, nil for a primitive one.
func (e *Element) Children() []*Element { return e.children }

// Encoded returns the element's original encoding, tag and length included. The bytes are the
// input's own, never re-encoded, so a digest computed over them matches the source exactly.
func (e *Element) Encoded() []byte { return e.encoded }

// IsUniversal reports whether the element carries the given universal tag number.
func (e *Element) IsUniversal(tagNumber uint64) bool {
	return e.class == ClassUniversal && e.tagNumber == tagNumber
}

// IsContextSpecific reports whether the element carries the given context-specific tag number.
func (e *Element) IsContextSpecific(tagNumber uint64) bool {
	return e.class == ClassContextSpecific && e.tagNumber == tagNumber
}

// IsASN1String reports whether the element is one of the types implementing BouncyCastle's
// ASN1String interface.
//
// ASN1BitString and ASN1UniversalString implement it too - their getString() answers "#"
// followed by the upper-case hex of the whole encoding, see AsString - while
// ASN1ObjectDescriptor does NOT (it wraps an ASN1GraphicString without implementing the
// interface), so tag 7 falls through to the hash-encoded branch of the callers.
func (e *Element) IsASN1String() bool {
	if e.class != ClassUniversal {
		return false
	}
	switch e.tagNumber {
	case TagBitString, TagUTF8String, TagNumericString,
		TagPrintableString, TagT61String, TagVideotexString,
		TagIA5String, TagGraphicString, TagVisibleString,
		TagGeneralString, TagUniversalString, TagBMPString:
		return true
	}
	return false
}

// AsString decodes the element's content the way the matching BouncyCastle ASN1String does:
// UTF8String as UTF-8, BMPString as UTF-16BE, BIT STRING and UniversalString as "#" followed
// by the UPPER-case hex of the complete encoding (ASN1BitString#getString and
// ASN1UniversalString#getString are both spelled that way, and neither decodes its content as
// text), and every other string type byte-per-character (ISO-8859-1), which is what
// org.bouncycastle.util.Strings#fromByteArray produces.
func (e *Element) AsString() string {
	switch e.tagNumber {
	case TagUTF8String:
		return string(e.content)
	case TagBitString, TagUniversalString:
		return "#" + hexUpper(e.DEREncoded())
	case TagBMPString:
		if len(e.content)%2 != 0 {
			return ""
		}
		units := make([]uint16, len(e.content)/2)
		for index := range units {
			units[index] = uint16(e.content[2*index])<<8 | uint16(e.content[2*index+1])
		}
		return string(utf16.Decode(units))
	}
	runes := make([]rune, len(e.content))
	for index, b := range e.content {
		runes[index] = rune(b)
	}
	return string(runes)
}

// Octets returns the content of an OCTET STRING, joining the segments of a constructed one.
func (e *Element) Octets() []byte {
	if !e.constructed {
		return e.content
	}
	var joined []byte
	for _, child := range e.children {
		joined = append(joined, child.Octets()...)
	}
	return joined
}

// BitStringOctets returns the value bits of a BIT STRING, i.e. its content without the
// leading count of unused bits, which is what BouncyCastle's ASN1BitString#getOctets returns.
func (e *Element) BitStringOctets() []byte {
	content := e.DERContent()
	if len(content) == 0 {
		return nil
	}
	return content[1:]
}

// Integer returns the value of an INTEGER, decoded as a two's-complement big-endian number.
func (e *Element) Integer() *big.Int {
	value := new(big.Int)
	if len(e.content) == 0 {
		return value
	}
	if e.content[0]&0x80 != 0 {
		// Negative: subtract 2^(8*len) from the unsigned interpretation.
		value.SetBytes(e.content)
		value.Sub(value, new(big.Int).Lsh(big.NewInt(1), uint(8*len(e.content))))
		return value
	}
	return value.SetBytes(e.content)
}

// ObjectIdentifier decodes an OBJECT IDENTIFIER.
func (e *Element) ObjectIdentifier() (asn1.ObjectIdentifier, error) {
	if !e.IsUniversal(TagOID) {
		return nil, errors.New("not an OBJECT IDENTIFIER")
	}
	var oid asn1.ObjectIdentifier
	if _, err := asn1.Unmarshal(e.DEREncoded(), &oid); err != nil {
		return nil, err
	}
	return oid, nil
}

// derBitStringContent returns the content octets of a BIT STRING with the unused bits of the
// final octet cleared, which X.690 clause 11.2.1 requires of DER and which
// org.bouncycastle.asn1.ASN1BitString#getDERContents does. The DL and BER encodings keep the
// octets as they arrived, so only DEREncoded calls this.
func derBitStringContent(content []byte) []byte {
	if len(content) < 2 {
		return content
	}
	unused := content[0]
	if unused == 0 || unused > 7 {
		return content
	}
	masked := make([]byte, len(content))
	copy(masked, content)
	masked[len(masked)-1] &= byte(0xFF) << unused
	return masked
}

// identifierOctets rebuilds the identifier octets of the element.
func (e *Element) identifierOctets(constructed bool) []byte {
	first := e.class
	if constructed {
		first |= Constructed
	}
	if e.tagNumber < 0x1F {
		return []byte{first | byte(e.tagNumber)}
	}
	first |= TagMask
	var trailer []byte
	value := e.tagNumber
	trailer = append(trailer, byte(value&0x7F))
	for value >>= 7; value > 0; value >>= 7 {
		trailer = append([]byte{byte(value&0x7F) | 0x80}, trailer...)
	}
	return append([]byte{first}, trailer...)
}

// DERContent returns the element's content in DER form, collapsing a constructed OCTET or
// BIT STRING into its primitive content.
func (e *Element) DERContent() []byte {
	if !e.constructed {
		return e.content
	}
	if e.class == ClassUniversal && e.tagNumber == TagOctetString {
		return e.Octets()
	}
	if e.class == ClassUniversal && e.tagNumber == TagBitString {
		// A segmented BIT STRING collapses to the concatenation of the segments' value
		// bits, prefixed by the unused-bit count of the final segment (only the final
		// segment may declare unused bits).
		var bits []byte
		unused := byte(0)
		for index, child := range e.children {
			segment := child.DERContent()
			if len(segment) == 0 {
				continue
			}
			bits = append(bits, segment[1:]...)
			if index == len(e.children)-1 {
				unused = segment[0]
			}
		}
		return append([]byte{unused}, bits...)
	}
	if e.class == ClassUniversal && e.tagNumber == TagSet {
		return sortSet(e.children)
	}
	var body []byte
	for _, child := range e.children {
		body = append(body, child.DEREncoded()...)
	}
	return body
}

// DEREncoded returns the DER encoding of the element.
//
// The conversion is a real BER-to-DER normalisation: indefinite lengths become definite,
// constructed OCTET/BIT STRINGs are collapsed into their primitive form, non-minimal length
// encodings are rewritten, BOOLEAN TRUE becomes 0xFF and the components of every SET are
// sorted by their encoding.
func (e *Element) DEREncoded() []byte {
	if !e.constructed {
		content := e.content
		switch {
		case e.class == ClassUniversal && e.tagNumber == TagBoolean && len(content) == 1 && content[0] != 0x00:
			// DER requires TRUE to be all-ones.
			content = []byte{0xFF}
		case e.class == ClassUniversal && e.tagNumber == TagBitString:
			content = derBitStringContent(content)
		}
		return WriteIdentifiedTLV(e.identifierOctets(false), content)
	}
	if e.class == ClassUniversal && e.tagNumber == TagBitString {
		// A constructed string is primitive in DER, and its unused bits have to be zero.
		return WriteIdentifiedTLV(e.identifierOctets(false), derBitStringContent(e.DERContent()))
	}
	if e.class == ClassUniversal && e.tagNumber == TagOctetString {
		// A constructed string is primitive in DER.
		return WriteIdentifiedTLV(e.identifierOctets(false), e.DERContent())
	}
	var body []byte
	if e.class == ClassUniversal && e.tagNumber == TagSet {
		body = sortSet(e.children)
	} else {
		for _, child := range e.children {
			body = append(body, child.DEREncoded()...)
		}
	}
	return WriteIdentifiedTLV(e.identifierOctets(true), body)
}

// DLEncoded returns the DL encoding of the element: definite lengths, and the component
// order of the input preserved (a SET is NOT sorted). A constructed OCTET or BIT STRING
// still collapses to its primitive form, as BouncyCastle's DL encoding of a BEROctetString
// or a BERBitString does.
func (e *Element) DLEncoded() []byte {
	if !e.constructed {
		return WriteIdentifiedTLV(e.identifierOctets(false), e.content)
	}
	if e.class == ClassUniversal && (e.tagNumber == TagOctetString || e.tagNumber == TagBitString) {
		return WriteIdentifiedTLV(e.identifierOctets(false), e.DERContent())
	}
	var body []byte
	for _, child := range e.children {
		body = append(body, child.DLEncoded()...)
	}
	return WriteIdentifiedTLV(e.identifierOctets(true), body)
}

// BEREncoded returns the BER encoding of the element, reproducing what BouncyCastle emits
// for an object read back from the same input:
//
//   - a constructed BIT STRING collapses to its primitive form;
//   - a constructed OCTET STRING becomes an indefinite-length constructed string. It holds
//     the input's own segments when it is the definite-length root - what
//     ASN1InputStream#buildConstructedOctetString hands to BEROctetString - and one segment
//     with the joined content in every other case, which is the BEROctetStringParser path;
//   - a definite-length constructed element encodes as its DL form, subtree included: a
//     BouncyCastle ASN1OutputStream switches to a DL sub-stream as soon as it writes a
//     definite-length object (DLSequence#encode and its siblings), so an indefinite-length
//     element nested inside a definite-length one comes out with a definite length;
//   - an indefinite-length constructed element keeps its indefinite length and encodes its
//     components in BER.
func (e *Element) BEREncoded() []byte {
	return e.berEncoded(true)
}

// berEncoded writes the BER encoding, root telling whether the element is the one the caller
// asked for rather than a component reached through an enclosing indefinite-length element.
func (e *Element) berEncoded(root bool) []byte {
	if !e.constructed {
		return WriteIdentifiedTLV(e.identifierOctets(false), e.content)
	}
	if e.class == ClassUniversal && e.tagNumber == TagBitString {
		return WriteIdentifiedTLV(e.identifierOctets(false), e.DERContent())
	}
	var body []byte
	if e.class == ClassUniversal && e.tagNumber == TagOctetString {
		if root && !e.indefinite {
			for _, child := range e.children {
				body = append(body, child.berEncoded(false)...)
			}
		} else {
			body = WriteTLV(TagOctetString, e.Octets())
		}
		return WriteIndefiniteTLV(e.identifierOctets(true), body)
	}
	if !e.indefinite {
		// Definite length: the whole subtree follows the DL rules, exactly as the DL
		// sub-stream of a BouncyCastle ASN1OutputStream writes it.
		return e.DLEncoded()
	}
	for _, child := range e.children {
		body = append(body, child.berEncoded(false)...)
	}
	return WriteIndefiniteTLV(e.identifierOctets(true), body)
}

// Parse reads one ASN.1 element from the front of the input, supporting both the definite and
// the indefinite length forms, and returns it together with the remaining bytes.
func Parse(input []byte) (*Element, []byte, error) {
	if len(input) < 2 {
		return nil, nil, errors.New("truncated ASN.1 element")
	}
	cursor := 0
	first := input[cursor]
	cursor++
	element := &Element{
		class:       first & ClassMask,
		constructed: first&Constructed != 0,
	}
	if first&TagMask == TagMask {
		var tagNumber uint64
		for {
			if cursor >= len(input) {
				return nil, nil, errors.New("truncated ASN.1 tag")
			}
			b := input[cursor]
			cursor++
			if tagNumber > (1<<56)-1 {
				return nil, nil, errors.New("ASN.1 tag number overflow")
			}
			tagNumber = tagNumber<<7 | uint64(b&0x7F)
			if b&0x80 == 0 {
				break
			}
		}
		element.tagNumber = tagNumber
	} else {
		element.tagNumber = uint64(first & TagMask)
	}

	if cursor >= len(input) {
		return nil, nil, errors.New("truncated ASN.1 length")
	}
	lengthByte := input[cursor]
	cursor++

	if lengthByte == 0x80 {
		if !element.constructed {
			return nil, nil, errors.New("indefinite length on a primitive ASN.1 element")
		}
		element.indefinite = true
		rest := input[cursor:]
		for {
			if len(rest) >= 2 && rest[0] == 0x00 && rest[1] == 0x00 {
				rest = rest[2:]
				break
			}
			child, remaining, err := Parse(rest)
			if err != nil {
				return nil, nil, err
			}
			element.children = append(element.children, child)
			rest = remaining
			if len(rest) < 2 {
				return nil, nil, errors.New("missing end-of-contents octets")
			}
		}
		element.encoded = input[:len(input)-len(rest)]
		return element, rest, nil
	}

	length := 0
	if lengthByte&0x80 == 0 {
		length = int(lengthByte)
	} else {
		count := int(lengthByte & 0x7F)
		if count > 4 || cursor+count > len(input) {
			return nil, nil, errors.New("unsupported or truncated ASN.1 length")
		}
		for _, b := range input[cursor : cursor+count] {
			length = length<<8 | int(b)
		}
		cursor += count
		if length < 0 {
			return nil, nil, errors.New("invalid ASN.1 length")
		}
	}
	if cursor+length > len(input) {
		return nil, nil, errors.New("truncated ASN.1 element")
	}
	body := input[cursor : cursor+length]
	element.encoded = input[:cursor+length]
	if element.constructed {
		rest := body
		for len(rest) > 0 {
			child, remaining, err := Parse(rest)
			if err != nil {
				return nil, nil, err
			}
			element.children = append(element.children, child)
			rest = remaining
		}
	} else {
		element.content = body
	}
	return element, input[cursor+length:], nil
}
