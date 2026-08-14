// GeneralName decoding, the read side of org.bouncycastle.asn1.x509.GeneralName. The structure
// itself lives in internal/asn1ber, which already writes it; only the parser is here, because
// the tsa field of a TSTInfo is the one place in this package that needs one.
package cmscore

import (
	"fmt"

	"github.com/utain/esig/dss/internal/asn1ber"
)

// generalNameUniversalTag maps a GeneralName alternative onto the universal tag its value
// carries when written on its own. Implicit tagging erased that tag on the wire, so decoding
// has to restore it for asn1ber.GeneralName.DER to be able to write the value back.
//
//	GeneralName ::= CHOICE {
//	    otherName                 [0] IMPLICIT OtherName,          -- SEQUENCE
//	    rfc822Name                [1] IMPLICIT IA5String,
//	    dNSName                   [2] IMPLICIT IA5String,
//	    x400Address               [3] IMPLICIT ORAddress,          -- SEQUENCE
//	    directoryName             [4] EXPLICIT Name,               -- a CHOICE, hence explicit
//	    ediPartyName              [5] IMPLICIT EDIPartyName,       -- SEQUENCE
//	    uniformResourceIdentifier [6] IMPLICIT IA5String,
//	    iPAddress                 [7] IMPLICIT OCTET STRING,
//	    registeredID              [8] IMPLICIT OBJECT IDENTIFIER }
var generalNameUniversalTag = map[int]byte{
	0: asn1ber.TagSequence | asn1ber.Constructed,
	1: asn1ber.TagIA5String,
	2: asn1ber.TagIA5String,
	3: asn1ber.TagSequence | asn1ber.Constructed,
	5: asn1ber.TagSequence | asn1ber.Constructed,
	6: asn1ber.TagIA5String,
	7: asn1ber.TagOctetString,
	8: asn1ber.TagOID,
}

// ParseGeneralName decodes one alternative of the GeneralName CHOICE. The returned Name holds
// the value's own DER encoding - the X.501 Name SEQUENCE for directoryName, the untagged value
// for the implicitly tagged alternatives - which is the form asn1ber.GeneralName.DER writes
// back, so a parsed GeneralName round-trips.
func ParseGeneralName(element *asn1ber.Element) (*asn1ber.GeneralName, error) {
	if element.Class() != asn1ber.ClassContextSpecific {
		return nil, fmt.Errorf("cmscore: GeneralName is not context-specific")
	}
	tagNo := int(element.TagNumber())
	if tagNo == 4 {
		// directoryName wraps a CHOICE and is therefore explicitly tagged.
		name, err := explicitContent(element, "GeneralName.directoryName")
		if err != nil {
			return nil, err
		}
		return &asn1ber.GeneralName{TagNo: tagNo, Name: name.DEREncoded()}, nil
	}
	identifier, known := generalNameUniversalTag[tagNo]
	if !known {
		return nil, fmt.Errorf("cmscore: unknown GeneralName alternative [%d]", tagNo)
	}
	return &asn1ber.GeneralName{
		TagNo: tagNo,
		Name:  asn1ber.WriteTLV(identifier, element.DERContent()),
	}, nil
}
