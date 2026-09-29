package asn1ber

import (
	"bytes"
	"encoding/asn1"
	"errors"
	"fmt"
	"math/big"
	"sort"
	"strconv"
	"strings"
)

// WriteTLV builds an element from a single-byte identifier and its content.
func WriteTLV(identifier byte, content []byte) []byte {
	return WriteIdentifiedTLV([]byte{identifier}, content)
}

// WriteIdentifiedTLV builds an element from its identifier octets and content, using the
// shortest definite-length encoding.
func WriteIdentifiedTLV(identifier []byte, content []byte) []byte {
	out := make([]byte, 0, len(identifier)+4+len(content))
	out = append(out, identifier...)
	length := len(content)
	if length < 0x80 {
		out = append(out, byte(length))
	} else {
		var lengthBytes []byte
		for value := length; value > 0; value >>= 8 {
			lengthBytes = append([]byte{byte(value & 0xFF)}, lengthBytes...)
		}
		out = append(out, byte(0x80|len(lengthBytes)))
		out = append(out, lengthBytes...)
	}
	return append(out, content...)
}

// WriteIndefiniteTLV builds a constructed element using the indefinite-length form, i.e. the
// 0x80 length octet followed by the content and the end-of-contents octets.
func WriteIndefiniteTLV(identifier []byte, content []byte) []byte {
	out := make([]byte, 0, len(identifier)+3+len(content))
	out = append(out, identifier...)
	out = append(out, 0x80)
	out = append(out, content...)
	return append(out, 0x00, 0x00)
}

// WriteSequence builds a SEQUENCE from its already encoded components.
func WriteSequence(content []byte) []byte {
	return WriteTLV(TagSequence|Constructed, content)
}

// EncodeOID returns the DER encoding of an OBJECT IDENTIFIER, nil when the value cannot be
// encoded (an OID of fewer than two arcs, or a first arc above 2).
func EncodeOID(oid asn1.ObjectIdentifier) []byte {
	encoded, err := asn1.Marshal(oid)
	if err != nil {
		return nil
	}
	return encoded
}

// EncodeInteger returns the DER encoding of an INTEGER.
func EncodeInteger(value *big.Int) []byte {
	encoded, err := asn1.Marshal(value)
	if err != nil {
		return nil
	}
	return encoded
}

// OIDFromString parses a dotted-decimal OID string, the way Java's
// new ASN1ObjectIdentifier(String) does: it accepts exactly the strings BouncyCastle's
// ASN1ObjectIdentifier#isValidIdentifier accepts - at least two arcs, a first arc of 0, 1 or 2,
// no empty arc, no sign, no leading zero, and a second arc below 40 unless the first arc is 2 -
// and refuses everything else, as the IllegalArgumentException("string X not a valid OID") does
// there. Accepting what BouncyCastle refuses is not harmless: EncodeOID answers nil for an OID
// it cannot encode, so an unchecked "5.3" or "1.2.-3" used to end up as a SEQUENCE whose OID was
// silently missing.
//
// An arc that does not fit an int is refused as well, although BouncyCastle keeps arbitrarily
// large arcs as BigIntegers: asn1.ObjectIdentifier is a []int and cannot hold them.
func OIDFromString(value string) (asn1.ObjectIdentifier, error) {
	if value == "" {
		return nil, errors.New("the OID is not defined")
	}
	if !isValidOIDString(value) {
		return nil, fmt.Errorf("string %s not a valid OID", value)
	}
	parts := strings.Split(value, ".")
	oid := make(asn1.ObjectIdentifier, 0, len(parts))
	for _, part := range parts {
		component, err := strconv.Atoi(part)
		if err != nil {
			return nil, fmt.Errorf("string %s not a valid OID: arc %s: %w", value, part, err)
		}
		oid = append(oid, component)
	}
	return oid, nil
}

// maxOIDStringLength is BouncyCastle's MAX_IDENTIFIER_LENGTH: 4096 content octets, four
// characters each, plus one.
const maxOIDStringLength = 4096*4 + 1

// isValidOIDString is ASN1ObjectIdentifier#isValidIdentifier, including its
// ASN1RelativeOID#isValidIdentifier(identifier, 2) check of the arcs after the first.
func isValidOIDString(identifier string) bool {
	if len(identifier) > maxOIDStringLength || len(identifier) < 3 || identifier[1] != '.' {
		return false
	}
	first := identifier[0]
	if first < '0' || first > '2' {
		return false
	}
	// The arcs after the first: digits separated by single dots, none empty, none with a
	// leading zero. BouncyCastle walks them from the end; so does this.
	digitCount := 0
	for pos := len(identifier) - 1; pos >= 2; pos-- {
		switch ch := identifier[pos]; {
		case ch == '.':
			if digitCount == 0 || (digitCount > 1 && identifier[pos+1] == '0') {
				return false
			}
			digitCount = 0
		case ch >= '0' && ch <= '9':
			digitCount++
		default:
			return false
		}
	}
	// pos is now 1, so identifier[2] is the first digit of the second arc.
	if digitCount == 0 || (digitCount > 1 && identifier[2] == '0') {
		return false
	}
	if first == '2' {
		return true
	}
	// Below a first arc of 0 or 1 the second arc is at most 39.
	if len(identifier) == 3 || identifier[3] == '.' {
		return true
	}
	if len(identifier) == 4 || identifier[4] == '.' {
		return identifier[2] < '4'
	}
	return false
}

// sortSet returns the DER encodings of a SET's components in ascending order, as DER requires
// and as BouncyCastle's ASN1Set#sortElements implements.
func sortSet(children []*Element) []byte {
	encodings := make([][]byte, len(children))
	for index, child := range children {
		encodings[index] = child.DEREncoded()
	}
	sort.SliceStable(encodings, func(a, b int) bool {
		return bytes.Compare(encodings[a], encodings[b]) < 0
	})
	var body []byte
	for _, encoding := range encodings {
		body = append(body, encoding...)
	}
	return body
}
