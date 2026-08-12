// Ported from dss-enumerations/.../EAACategory.java (DSS 6.5.RC1).

package enumerations

import "fmt"

// EAACategory provides a list of EAA category definitions.
type EAACategory string

const (
	// EAACategory_EU_QEAA: indication that the attestation has been issued as a
	// qualified electronic attestation of attributes.
	EAACategory_EU_QEAA EAACategory = "EU_QEAA"
	// EAACategory_EU_PUBEAA: indication that the attestation has been issued as an
	// electronic attestation of attributes issued by or on behalf of a public body
	// responsible for an authentic source.
	EAACategory_EU_PUBEAA EAACategory = "EU_PUBEAA"
)

// eaaCategoryURN maps each EAACategory to its defined URN.
var eaaCategoryURN = map[EAACategory]string{
	EAACategory_EU_QEAA:   "urn:etsi:esi:eaa:eu:qualified",
	EAACategory_EU_PUBEAA: "urn:etsi:esi:eaa:eu:pub",
}

// EAACategoryValues returns all EAACategory constants in declaration order.
func EAACategoryValues() []EAACategory {
	return []EAACategory{EAACategory_EU_QEAA, EAACategory_EU_PUBEAA}
}

// EAACategoryValueOf returns the EAACategory matching the given Java enum name.
func EAACategoryValueOf(name string) (EAACategory, error) {
	for _, v := range EAACategoryValues() {
		if string(v) == name {
			return v, nil
		}
	}
	return "", fmt.Errorf("no enum constant EAACategory.%s", name)
}

// URN returns the URN defined for the EAA category.
func (e EAACategory) URN() string {
	return eaaCategoryURN[e]
}
