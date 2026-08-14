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
// new ASN1ObjectIdentifier(String) does.
func OIDFromString(value string) (asn1.ObjectIdentifier, error) {
	if value == "" {
		return nil, errors.New("the OID is not defined")
	}
	parts := strings.Split(value, ".")
	oid := make(asn1.ObjectIdentifier, 0, len(parts))
	for _, part := range parts {
		component, err := strconv.Atoi(part)
		if err != nil {
			return nil, fmt.Errorf("string %s not an OID", value)
		}
		oid = append(oid, component)
	}
	return oid, nil
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
