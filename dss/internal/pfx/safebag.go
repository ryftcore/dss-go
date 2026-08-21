package pfx

import (
	"encoding/asn1"
	"errors"
	"fmt"
	"unicode/utf16"

	"github.com/ryftcore/dss-go/dss/internal/asn1ber"
)

// safeBag is one parsed, not-yet-decrypted SafeBag component (RFC 7292 section 4.2): its type
// (id), its EXPLICIT bagValue content, and the two PKCS12Attribute values a certBag and the
// (shrouded) keyBag it belongs to are correlated by.
type safeBag struct {
	id           asn1.ObjectIdentifier
	value        *asn1ber.Element
	friendlyName string
	localKeyID   []byte
}

// parseSafeContents decodes a SafeContents ::= SEQUENCE OF SafeBag.
func parseSafeContents(der []byte) ([]safeBag, error) {
	element, rest, err := asn1ber.Parse(der)
	if err != nil {
		return nil, fmt.Errorf("pfx: %w", err)
	}
	if len(rest) != 0 {
		return nil, errors.New("pfx: trailing data after SafeContents")
	}
	if !element.IsUniversal(asn1ber.TagSequence) || !element.IsConstructed() {
		return nil, errors.New("pfx: SafeContents is not a SEQUENCE")
	}
	bags := make([]safeBag, 0, len(element.Children()))
	for _, child := range element.Children() {
		bag, err := parseSafeBag(child)
		if err != nil {
			return nil, err
		}
		bags = append(bags, bag)
	}
	return bags, nil
}

// parseSafeBag decodes one
//
//	SafeBag ::= SEQUENCE {
//	    bagId         OBJECT IDENTIFIER,
//	    bagValue      [0] EXPLICIT ANY DEFINED BY bagId,
//	    bagAttributes SET OF PKCS12Attribute OPTIONAL }
func parseSafeBag(element *asn1ber.Element) (safeBag, error) {
	if !element.IsUniversal(asn1ber.TagSequence) || !element.IsConstructed() {
		return safeBag{}, errors.New("pfx: SafeBag is not a SEQUENCE")
	}
	children := element.Children()
	if len(children) < 2 {
		return safeBag{}, fmt.Errorf("pfx: SafeBag holds %d components, at least 2 expected", len(children))
	}
	id, err := children[0].ObjectIdentifier()
	if err != nil {
		return safeBag{}, err
	}
	if !children[1].IsContextSpecific(0) {
		return safeBag{}, errors.New("pfx: SafeBag.bagValue is not [0] EXPLICIT")
	}
	value, err := explicitContent(children[1])
	if err != nil {
		return safeBag{}, err
	}
	bag := safeBag{id: id, value: value}
	if len(children) > 2 {
		friendlyName, localKeyID, err := parseBagAttributes(children[2])
		if err != nil {
			return safeBag{}, err
		}
		bag.friendlyName = friendlyName
		bag.localKeyID = localKeyID
	}
	return bag, nil
}

// explicitContent unwraps an EXPLICIT-tagged element, returning the single child it wraps.
func explicitContent(element *asn1ber.Element) (*asn1ber.Element, error) {
	if !element.IsConstructed() || len(element.Children()) != 1 {
		return nil, errors.New("pfx: not an EXPLICIT tagged element")
	}
	return element.Children()[0], nil
}

// parseBagAttributes decodes a SafeBag.bagAttributes SET OF PKCS12Attribute, extracting the two
// attributes a certBag and its owning (shrouded) keyBag correlate through: friendlyName (a
// BMPString) and localKeyId (an OCTET STRING). Every other attribute (e.g. Microsoft's CSP
// name) is ignored, as neither openssl nor the JDK's PKCS12 provider needs it to open the
// store.
func parseBagAttributes(element *asn1ber.Element) (friendlyName string, localKeyID []byte, err error) {
	if !element.IsUniversal(asn1ber.TagSet) || !element.IsConstructed() {
		return "", nil, errors.New("pfx: SafeBag.bagAttributes is not a SET")
	}
	for _, attribute := range element.Children() {
		if !attribute.IsUniversal(asn1ber.TagSequence) || !attribute.IsConstructed() || len(attribute.Children()) != 2 {
			return "", nil, errors.New("pfx: PKCS12Attribute is not a two-component SEQUENCE")
		}
		attrID, err := attribute.Children()[0].ObjectIdentifier()
		if err != nil {
			return "", nil, err
		}
		values := attribute.Children()[1]
		if !values.IsUniversal(asn1ber.TagSet) || !values.IsConstructed() || len(values.Children()) != 1 {
			continue // an attribute with no or multiple values carries nothing this package reads.
		}
		value := values.Children()[0]
		switch {
		case attrID.Equal(oidFriendlyName):
			friendlyName = decodeBMPStringAttribute(value)
		case attrID.Equal(oidLocalKeyID):
			localKeyID = value.Octets()
		}
	}
	return friendlyName, localKeyID, nil
}

// decodeBMPStringAttribute decodes a friendlyName attribute value (a BMPString, i.e. UTF-16BE),
// stripping a trailing NUL code point some producers append despite it not being part of the
// ASN.1 BMPString type itself (openssl's own PKCS#12 attribute encoder does this, carrying over
// the RFC 7292 Appendix B.1 password convention). An unexpected value type decodes to "" rather
// than failing the whole load - a cosmetic attribute is not worth refusing an otherwise valid
// store over.
func decodeBMPStringAttribute(element *asn1ber.Element) string {
	if !element.IsUniversal(asn1ber.TagBMPString) {
		return ""
	}
	content := element.Content()
	if len(content) >= 2 && content[len(content)-1] == 0 && content[len(content)-2] == 0 {
		content = content[:len(content)-2]
	}
	if len(content)%2 != 0 {
		return ""
	}
	units := make([]uint16, len(content)/2)
	for i := range units {
		units[i] = uint16(content[2*i])<<8 | uint16(content[2*i+1])
	}
	return string(utf16.Decode(units))
}
