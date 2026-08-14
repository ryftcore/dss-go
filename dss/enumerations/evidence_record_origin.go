// Ported from dss-enumerations/.../EvidenceRecordOrigin.java (DSS 6.5.RC1).
//
// Defines the origin of an Evidence Record.
package enumerations

import "fmt"

// EvidenceRecordOrigin defines the origin of an Evidence Record.
type EvidenceRecordOrigin string

const (
	// EvidenceRecordOrigin_CONTAINER defines an evidence record extracted
	// from an ASiC container.
	EvidenceRecordOrigin_CONTAINER EvidenceRecordOrigin = "CONTAINER"
	// EvidenceRecordOrigin_SIGNATURE defines an evidence record embedded in
	// electronic signature.
	EvidenceRecordOrigin_SIGNATURE EvidenceRecordOrigin = "SIGNATURE"
	// EvidenceRecordOrigin_EXTERNAL is an evidence record provided
	// externally to the validation.
	EvidenceRecordOrigin_EXTERNAL EvidenceRecordOrigin = "EXTERNAL"
)

// EvidenceRecordOriginValues returns all constants in declaration order.
func EvidenceRecordOriginValues() []EvidenceRecordOrigin {
	return []EvidenceRecordOrigin{
		EvidenceRecordOrigin_CONTAINER,
		EvidenceRecordOrigin_SIGNATURE,
		EvidenceRecordOrigin_EXTERNAL,
	}
}

// EvidenceRecordOriginValueOf returns the EvidenceRecordOrigin matching the given Java enum name.
func EvidenceRecordOriginValueOf(name string) (EvidenceRecordOrigin, error) {
	for _, v := range EvidenceRecordOriginValues() {
		if string(v) == name {
			return v, nil
		}
	}
	return "", fmt.Errorf("no enum constant EvidenceRecordOrigin.%s", name)
}
