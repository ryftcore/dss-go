// Ported from dss-enumerations/.../EndorsementType.java (DSS 6.5.RC1).

package enumerations

import "fmt"

// EndorsementType defines available types of a SignerRole element.
type EndorsementType string

const (
	// EndorsementType_CERTIFIED: attributes certified in attribute certificates
	// issued by an Attribute Authority.
	EndorsementType_CERTIFIED EndorsementType = "CERTIFIED"
	// EndorsementType_CLAIMED: attributes claimed by the signer.
	EndorsementType_CLAIMED EndorsementType = "CLAIMED"
	// EndorsementType_SIGNED: assertions signed by a third party.
	EndorsementType_SIGNED EndorsementType = "SIGNED"
)

// endorsementTypeValue maps each EndorsementType to its string value.
var endorsementTypeValue = map[EndorsementType]string{
	EndorsementType_CERTIFIED: "certified",
	EndorsementType_CLAIMED:   "claimed",
	EndorsementType_SIGNED:    "signed",
}

// EndorsementTypeValues returns all EndorsementType constants in declaration order.
func EndorsementTypeValues() []EndorsementType {
	return []EndorsementType{EndorsementType_CERTIFIED, EndorsementType_CLAIMED, EndorsementType_SIGNED}
}

// EndorsementTypeValueOf returns the EndorsementType matching the given Java enum name.
func EndorsementTypeValueOf(name string) (EndorsementType, error) {
	for _, v := range EndorsementTypeValues() {
		if string(v) == name {
			return v, nil
		}
	}
	return "", fmt.Errorf("no enum constant EndorsementType.%s", name)
}

// Value returns the string value of the enumeration.
func (e EndorsementType) Value() string {
	return endorsementTypeValue[e]
}

// EndorsementTypeFromString parses the string value and returns the matching
// EndorsementType, or "" if none matches. Upstream returns null here rather than
// throwing, so this is a plain zero-value return and not an error.
func EndorsementTypeFromString(value string) EndorsementType {
	for _, v := range EndorsementTypeValues() {
		if endorsementTypeValue[v] == value {
			return v
		}
	}
	return ""
}
