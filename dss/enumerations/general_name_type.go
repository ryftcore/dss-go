// Ported from dss-enumerations/.../GeneralNameType.java (DSS 6.5.RC1).
package enumerations

import "fmt"

// GeneralNameType represents possible types of a GeneralName.
type GeneralNameType string

const (
	// GeneralNameTypeOtherName is the otherName GeneralName type.
	GeneralNameTypeOtherName GeneralNameType = "OTHER_NAME"
	// GeneralNameTypeRFC822Name is the rfc822Name GeneralName type.
	GeneralNameTypeRFC822Name GeneralNameType = "RFC822_NAME"
	// GeneralNameTypeDNSName is the dNSName GeneralName type.
	GeneralNameTypeDNSName GeneralNameType = "DNS_NAME"
	// GeneralNameTypeX400Address is the x400Address GeneralName type.
	GeneralNameTypeX400Address GeneralNameType = "X400_ADDRESS"
	// GeneralNameTypeDirectoryName is the directoryName GeneralName type.
	GeneralNameTypeDirectoryName GeneralNameType = "DIRECTORY_NAME"
	// GeneralNameTypeEDIPartyName is the ediPartyName GeneralName type.
	GeneralNameTypeEDIPartyName GeneralNameType = "EDI_PARTY_NAME"
	// GeneralNameTypeUniformResourceIdentifier is the
	// uniformResourceIdentifier GeneralName type.
	GeneralNameTypeUniformResourceIdentifier GeneralNameType = "UNIFORM_RESOURCE_IDENTIFIER"
	// GeneralNameTypeIPAddress is the iPAddress GeneralName type.
	GeneralNameTypeIPAddress GeneralNameType = "IP_ADDRESS"
	// GeneralNameTypeRegisteredID is the registeredID GeneralName type.
	GeneralNameTypeRegisteredID GeneralNameType = "REGISTERED_ID"
)

// generalNameTypeFields holds the (index, label) pair for each constant.
type generalNameTypeFields struct {
	index int
	label string
}

// generalNameTypeData holds the full field tuple for each constant, copied
// verbatim from the Java enum constructors.
var generalNameTypeData = map[GeneralNameType]generalNameTypeFields{
	GeneralNameTypeOtherName:                 {0, "otherName"},
	GeneralNameTypeRFC822Name:                {1, "rfc822Name"},
	GeneralNameTypeDNSName:                   {2, "dNSName"},
	GeneralNameTypeX400Address:               {3, "x400Address"},
	GeneralNameTypeDirectoryName:             {4, "directoryName"},
	GeneralNameTypeEDIPartyName:              {5, "ediPartyName"},
	GeneralNameTypeUniformResourceIdentifier: {6, "uniformResourceIdentifier"},
	GeneralNameTypeIPAddress:                 {7, "iPAddress"},
	GeneralNameTypeRegisteredID:              {8, "registeredID"},
}

// GeneralNameTypeValues returns all constants in declaration order.
func GeneralNameTypeValues() []GeneralNameType {
	return []GeneralNameType{
		GeneralNameTypeOtherName,
		GeneralNameTypeRFC822Name,
		GeneralNameTypeDNSName,
		GeneralNameTypeX400Address,
		GeneralNameTypeDirectoryName,
		GeneralNameTypeEDIPartyName,
		GeneralNameTypeUniformResourceIdentifier,
		GeneralNameTypeIPAddress,
		GeneralNameTypeRegisteredID,
	}
}

// GeneralNameTypeValueOf returns the GeneralNameType matching the given Java enum name.
func GeneralNameTypeValueOf(name string) (GeneralNameType, error) {
	for _, v := range GeneralNameTypeValues() {
		if string(v) == name {
			return v, nil
		}
	}
	return "", fmt.Errorf("no enum constant GeneralNameType.%s", name)
}

// Index gets the index of the GeneralName type.
func (g GeneralNameType) Index() int {
	return generalNameTypeData[g].index
}

// Label gets a human-readable label of the GeneralName type.
func (g GeneralNameType) Label() string {
	return generalNameTypeData[g].label
}

// GeneralNameTypeFromIndex returns a GeneralNameType for the given index if
// it exists, or "" otherwise.
func GeneralNameTypeFromIndex(index int) GeneralNameType {
	for _, v := range GeneralNameTypeValues() {
		if index == v.Index() {
			return v
		}
	}
	return ""
}

// GeneralNameTypeFromLabel returns a GeneralNameType for the given label if
// it exists, or "" otherwise.
func GeneralNameTypeFromLabel(label string) GeneralNameType {
	for _, v := range GeneralNameTypeValues() {
		if label == v.Label() {
			return v
		}
	}
	return ""
}
