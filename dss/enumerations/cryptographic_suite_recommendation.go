// Ported from dss-enumerations/.../CryptographicSuiteRecommendation.java (DSS 6.5.RC1).
//
// The recommendation element shall be used to indicate that a mechanism and
// its parameters are either Recommended (R) or Legacy (L), as defined in
// ETSI TS 119 312 [i.2], clause 3.1.
package enumerations

import "fmt"

// CryptographicSuiteRecommendation represents whether a mechanism and its
// parameters are Recommended or Legacy.
type CryptographicSuiteRecommendation string

const (
	// CryptographicSuiteRecommendationRecommended is the recommended
	// cryptographic algorithm.
	CryptographicSuiteRecommendationRecommended CryptographicSuiteRecommendation = "RECOMMENDED"
	// CryptographicSuiteRecommendationLegacy is the legacy cryptographic
	// algorithm.
	CryptographicSuiteRecommendationLegacy CryptographicSuiteRecommendation = "LEGACY"
)

// cryptographicSuiteRecommendationValues holds the string value for each constant.
var cryptographicSuiteRecommendationValues = map[CryptographicSuiteRecommendation]string{
	CryptographicSuiteRecommendationRecommended: "R",
	CryptographicSuiteRecommendationLegacy:      "L",
}

// CryptographicSuiteRecommendationValues returns all constants in declaration order.
func CryptographicSuiteRecommendationValues() []CryptographicSuiteRecommendation {
	return []CryptographicSuiteRecommendation{
		CryptographicSuiteRecommendationRecommended,
		CryptographicSuiteRecommendationLegacy,
	}
}

// CryptographicSuiteRecommendationValueOf returns the CryptographicSuiteRecommendation matching the given Java enum name.
func CryptographicSuiteRecommendationValueOf(name string) (CryptographicSuiteRecommendation, error) {
	for _, v := range CryptographicSuiteRecommendationValues() {
		if string(v) == name {
			return v, nil
		}
	}
	return "", fmt.Errorf("no enum constant CryptographicSuiteRecommendation.%s", name)
}

// Value gets the value of the recommendation type.
func (c CryptographicSuiteRecommendation) Value() string {
	return cryptographicSuiteRecommendationValues[c]
}

// CryptographicSuiteRecommendationFromValue returns a
// CryptographicSuiteRecommendation by the given value, or "" if not found.
func CryptographicSuiteRecommendationFromValue(value string) CryptographicSuiteRecommendation {
	if value != "" {
		for _, v := range CryptographicSuiteRecommendationValues() {
			if v.Value() == value {
				return v
			}
		}
	}
	return ""
}
