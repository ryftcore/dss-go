// Ported from dss-enumerations/.../GeneralNameType.java (DSS 6.5.RC1).
package enumerations

import "fmt"

// GeneralNameType represents possible types of a GeneralName.
type GeneralNameType string

const (
	// GeneralNameType_OTHER_NAME is the otherName GeneralName type.
	GeneralNameType_OTHER_NAME GeneralNameType = "OTHER_NAME"
	// GeneralNameType_RFC822_NAME is the rfc822Name GeneralName type.
	GeneralNameType_RFC822_NAME GeneralNameType = "RFC822_NAME"
	// GeneralNameType_DNS_NAME is the dNSName GeneralName type.
	GeneralNameType_DNS_NAME GeneralNameType = "DNS_NAME"
	// GeneralNameType_X400_ADDRESS is the x400Address GeneralName type.
	GeneralNameType_X400_ADDRESS GeneralNameType = "X400_ADDRESS"
	// GeneralNameType_DIRECTORY_NAME is the directoryName GeneralName type.
	GeneralNameType_DIRECTORY_NAME GeneralNameType = "DIRECTORY_NAME"
	// GeneralNameType_EDI_PARTY_NAME is the ediPartyName GeneralName type.
	GeneralNameType_EDI_PARTY_NAME GeneralNameType = "EDI_PARTY_NAME"
	// GeneralNameType_UNIFORM_RESOURCE_IDENTIFIER is the
	// uniformResourceIdentifier GeneralName type.
	GeneralNameType_UNIFORM_RESOURCE_IDENTIFIER GeneralNameType = "UNIFORM_RESOURCE_IDENTIFIER"
	// GeneralNameType_IP_ADDRESS is the iPAddress GeneralName type.
	GeneralNameType_IP_ADDRESS GeneralNameType = "IP_ADDRESS"
	// GeneralNameType_REGISTERED_ID is the registeredID GeneralName type.
	GeneralNameType_REGISTERED_ID GeneralNameType = "REGISTERED_ID"
)

// generalNameTypeFields holds the (index, label) pair for each constant.
type generalNameTypeFields struct {
	index int
	label string
}

// generalNameTypeData holds the full field tuple for each constant, copied
// verbatim from the Java enum constructors.
var generalNameTypeData = map[GeneralNameType]generalNameTypeFields{
	GeneralNameType_OTHER_NAME:                  {0, "otherName"},
	GeneralNameType_RFC822_NAME:                 {1, "rfc822Name"},
	GeneralNameType_DNS_NAME:                    {2, "dNSName"},
	GeneralNameType_X400_ADDRESS:                {3, "x400Address"},
	GeneralNameType_DIRECTORY_NAME:              {4, "directoryName"},
	GeneralNameType_EDI_PARTY_NAME:              {5, "ediPartyName"},
	GeneralNameType_UNIFORM_RESOURCE_IDENTIFIER: {6, "uniformResourceIdentifier"},
	GeneralNameType_IP_ADDRESS:                  {7, "iPAddress"},
	GeneralNameType_REGISTERED_ID:               {8, "registeredID"},
}

// GeneralNameTypeValues returns all constants in declaration order.
func GeneralNameTypeValues() []GeneralNameType {
	return []GeneralNameType{
		GeneralNameType_OTHER_NAME,
		GeneralNameType_RFC822_NAME,
		GeneralNameType_DNS_NAME,
		GeneralNameType_X400_ADDRESS,
		GeneralNameType_DIRECTORY_NAME,
		GeneralNameType_EDI_PARTY_NAME,
		GeneralNameType_UNIFORM_RESOURCE_IDENTIFIER,
		GeneralNameType_IP_ADDRESS,
		GeneralNameType_REGISTERED_ID,
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
