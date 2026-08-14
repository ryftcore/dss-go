// Ported from dss-enumerations/.../ValidationLevel.java (DSS 6.5.RC1).

package enumerations

import "fmt"

// ValidationLevel is the target validation level as per EN 319 102-1.
//
// NOTE: the validation process "stops" processing on the chosen level.
type ValidationLevel string

const (
	// ValidationLevel_BASIC_SIGNATURES: validation as per "5.3 Validation process for
	// Basic Signatures".
	ValidationLevel_BASIC_SIGNATURES ValidationLevel = "BASIC_SIGNATURES"
	// ValidationLevel_TIMESTAMPS: validation as per "5.4 Time-stamp validation
	// building block".
	ValidationLevel_TIMESTAMPS ValidationLevel = "TIMESTAMPS"
	// ValidationLevel_LONG_TERM_DATA: validation as per "5.5 Validation process for
	// Signatures with Time and Signatures with Long-Term Validation Material".
	ValidationLevel_LONG_TERM_DATA ValidationLevel = "LONG_TERM_DATA"
	// ValidationLevel_ARCHIVAL_DATA: validation as per "5.6 Validation process for
	// Signatures providing Long Term Availability and Integrity of Validation Material".
	ValidationLevel_ARCHIVAL_DATA ValidationLevel = "ARCHIVAL_DATA"
)

// ValidationLevelValues returns all ValidationLevel constants in declaration order.
func ValidationLevelValues() []ValidationLevel {
	return []ValidationLevel{
		ValidationLevel_BASIC_SIGNATURES,
		ValidationLevel_TIMESTAMPS,
		ValidationLevel_LONG_TERM_DATA,
		ValidationLevel_ARCHIVAL_DATA,
	}
}

// ValidationLevelValueOf returns the ValidationLevel matching the given Java enum name.
func ValidationLevelValueOf(name string) (ValidationLevel, error) {
	for _, v := range ValidationLevelValues() {
		if string(v) == name {
			return v, nil
		}
	}
	return "", fmt.Errorf("no enum constant ValidationLevel.%s", name)
}
