// Ported from dss-jades/src/main/java/eu/europa/esig/dss/jades/JAdESSigningTimeType.java (DSS 6.5.RC1).
package jades

// SigningTimeType contains possible types for a claimed signing time header incorporation
// within a JAdES signature.
type SigningTimeType string

const (
	// JAdESSigningTimeTypeIAT is the iat header parameter as specified in IETF RFC 7519,
	// clause 4.1.6. Before 2025-07-15T00:00:00Z, this header parameter should be incorporated in
	// new JAdES signatures instead the sigT header parameter specified in clause 5.2.1 of the
	// present document. Starting at 2025-07-15T00:00:00Z, this header parameter shall be
	// incorporated in new JAdES signatures.
	JAdESSigningTimeTypeIAT SigningTimeType = "IAT"

	// JAdESSigningTimeTypeSigT is the sigT header parameter as specified in ETSI TS 119 182-1,
	// clause 5.2.1. Before 2025-07-15T00:00:00Z this header parameter should not be incorporated
	// in new JAdES signatures. Instead, the iaT header parameter should be included. Starting at
	// 2025-07-15T00:00:00Z this header parameter shall not be incorporated in new JAdES
	// signatures.
	JAdESSigningTimeTypeSigT SigningTimeType = "SIG_T"

	// JAdESSigningTimeTypeNone means no signing time header.
	JAdESSigningTimeTypeNone SigningTimeType = "NONE"
)

// SigningTimeTypeValues returns all constants in declaration order.
func SigningTimeTypeValues() []SigningTimeType {
	return []SigningTimeType{
		JAdESSigningTimeTypeIAT,
		JAdESSigningTimeTypeSigT,
		JAdESSigningTimeTypeNone,
	}
}
