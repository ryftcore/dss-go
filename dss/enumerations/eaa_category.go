// Ported from dss-enumerations/.../EAACategory.java (DSS 6.5.RC1).

package enumerations

import "fmt"

// EAACategory provides a list of EAA category definitions.
type EAACategory string

const (
	// EAACategoryEUQEAA: indication that the attestation has been issued as a
	// qualified electronic attestation of attributes.
	EAACategoryEUQEAA EAACategory = "EU_QEAA"
	// EAACategoryEUPubEAA: indication that the attestation has been issued as an
	// electronic attestation of attributes issued by or on behalf of a public body
	// responsible for an authentic source.
	EAACategoryEUPubEAA EAACategory = "EU_PUBEAA"
)

// eaaCategoryURN maps each EAACategory to its defined URN.
var eaaCategoryURN = map[EAACategory]string{
	EAACategoryEUQEAA:   "urn:etsi:esi:eaa:eu:qualified",
	EAACategoryEUPubEAA: "urn:etsi:esi:eaa:eu:pub",
}

// EAACategoryValues returns all EAACategory constants in declaration order.
func EAACategoryValues() []EAACategory {
	return []EAACategory{EAACategoryEUQEAA, EAACategoryEUPubEAA}
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
