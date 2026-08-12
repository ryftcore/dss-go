// Ported from dss-enumerations/.../SignatureProfile.java (DSS 6.5.RC1).
package enumerations

import "strings"

// SignatureProfile represents a "generic" signature augmentation profile
// level, such as "*AdES-BASELINE-T".
type SignatureProfile string

const (
	// SignatureProfile_BASELINE_B is a Basic Signature, incorporating
	// signed and some unsigned qualifying properties when the signature is
	// generated. Example: XAdES-BASELINE-B (ETSI EN 319 132-1).
	SignatureProfile_BASELINE_B SignatureProfile = "BASELINE_B"
	// SignatureProfile_BASELINE_T is a Signature with Time, incorporating a
	// trusted token proving that the signature itself actually existed at
	// a certain date and time. Example: XAdES-BASELINE-T (ETSI EN 319 132-1).
	SignatureProfile_BASELINE_T SignatureProfile = "BASELINE_T"
	// SignatureProfile_BASELINE_LT is a Signature with Long Term Validation
	// Material, incorporating all the material required for validating the
	// signature in the signature document. This level aims to tackle the
	// long term availability of the validation material. Example:
	// XAdES-BASELINE-LT (ETSI EN 319 132-1).
	SignatureProfile_BASELINE_LT SignatureProfile = "BASELINE_LT"
	// SignatureProfile_BASELINE_LTA is a Signature with Long Term
	// Availability and Integrity of Validation Material, incorporating
	// electronic time-stamps that allow validation of the signature long
	// time after its generation. This level aims to tackle the long term
	// availability and integrity of the validation material. Example:
	// XAdES-BASELINE-LTA (ETSI EN 319 132-1).
	SignatureProfile_BASELINE_LTA SignatureProfile = "BASELINE_LTA"
	// SignatureProfile_EXTENDED_BES is a legacy Basic Signature profile.
	// Example: XAdES-E-BES (ETSI TS 319 132-2).
	SignatureProfile_EXTENDED_BES SignatureProfile = "EXTENDED_BES"
	// SignatureProfile_EXTENDED_EPES is a legacy Basic Signature profile
	// with SignaturePolicyIdentifier qualifying property. Example:
	// XAdES-E-EPES (ETSI TS 319 132-2).
	SignatureProfile_EXTENDED_EPES SignatureProfile = "EXTENDED_EPES"
	// SignatureProfile_EXTENDED_T is a legacy Signature with Time profile
	// with a signature timestamp qualifying property. Example: XAdES-E-T
	// (ETSI TS 319 132-2).
	SignatureProfile_EXTENDED_T SignatureProfile = "EXTENDED_T"
	// SignatureProfile_EXTENDED_LT is a legacy Signature with Long Term
	// Validation Material profile built on top of EXTENDED-T with all the
	// material required for validating the signature in the signature
	// document. Example: XAdES-E-LT.
	SignatureProfile_EXTENDED_LT SignatureProfile = "EXTENDED_LT"
	// SignatureProfile_EXTENDED_C is a legacy Signature profile built on
	// top of EXTENDED-T with qualifying properties containing references
	// to certificates and references to certificate status data values.
	// Example: XAdES-E-C (ETSI TS 319 132-2).
	SignatureProfile_EXTENDED_C SignatureProfile = "EXTENDED_C"
	// SignatureProfile_EXTENDED_X is a legacy Signature profile built on
	// top of EXTENDED-C with one or more qualifying properties containing
	// one or more electronic time-stamps. Example: XAdES-E-X (ETSI TS 319
	// 132-2).
	SignatureProfile_EXTENDED_X SignatureProfile = "EXTENDED_X"
	// SignatureProfile_EXTENDED_XL is a legacy Signature profile built on
	// top of EXTENDED-X with qualifying properties that contain
	// certificates and revocation values. Example: XAdES-E-XL (ETSI TS 319
	// 132-2).
	SignatureProfile_EXTENDED_XL SignatureProfile = "EXTENDED_XL"
	// SignatureProfile_EXTENDED_A is a legacy Signature with Long Term
	// Availability and Integrity of Validation Material profile with an
	// archive timestamp qualifying property. Example: XAdES-E-A (ETSI TS
	// 319 132-2).
	SignatureProfile_EXTENDED_A SignatureProfile = "EXTENDED_A"
	// SignatureProfile_EXTENDED_ERS is a Signature with Evidence Record
	// profile built on top of BASELINE-* or EXTENDED-* profile, containing
	// an evidence record unsigned qualified property. Example: XAdES-E-ERS
	// (ETSI TS 319 132-3).
	SignatureProfile_EXTENDED_ERS SignatureProfile = "EXTENDED_ERS"
	// SignatureProfile_EXTENDED_LTV is a legacy Signature with Long Term
	// Availability and Integrity of Validation Material profile with
	// document security store and one or more document time-stamp (PAdES
	// only). Example: PAdES-E-LTV (ETSI TS 119 142-2).
	SignatureProfile_EXTENDED_LTV SignatureProfile = "EXTENDED_LTV"
	// SignatureProfile_AdES defines an AdES profile, which is not BASELINE
	// or EXTENDED.
	SignatureProfile_AdES SignatureProfile = "AdES"
	// SignatureProfile_NOT_ETSI represents an unknown or a not supported
	// signature profile.
	SignatureProfile_NOT_ETSI SignatureProfile = "NOT_ETSI"
)

// SignatureProfileValues returns all constants in declaration order.
func SignatureProfileValues() []SignatureProfile {
	return []SignatureProfile{
		SignatureProfile_BASELINE_B,
		SignatureProfile_BASELINE_T,
		SignatureProfile_BASELINE_LT,
		SignatureProfile_BASELINE_LTA,
		SignatureProfile_EXTENDED_BES,
		SignatureProfile_EXTENDED_EPES,
		SignatureProfile_EXTENDED_T,
		SignatureProfile_EXTENDED_LT,
		SignatureProfile_EXTENDED_C,
		SignatureProfile_EXTENDED_X,
		SignatureProfile_EXTENDED_XL,
		SignatureProfile_EXTENDED_A,
		SignatureProfile_EXTENDED_ERS,
		SignatureProfile_EXTENDED_LTV,
		SignatureProfile_AdES,
		SignatureProfile_NOT_ETSI,
	}
}

// SignatureProfileValueOf returns the constant matching the given Java enum
// name.
func SignatureProfileValueOf(name string) (SignatureProfile, error) {
	for _, v := range SignatureProfileValues() {
		if string(v) == name {
			return v, nil
		}
	}
	return "", &signatureProfileInvalidValueError{name}
}

type signatureProfileInvalidValueError struct {
	name string
}

func (e *signatureProfileInvalidValueError) Error() string {
	return "no enum constant SignatureProfile." + e.name
}

// SignatureProfileValueByName returns the SignatureProfile based on the
// given name, replacing '-' with '_' before matching, mirroring Java's
// valueByName(String).
func SignatureProfileValueByName(name string) (SignatureProfile, error) {
	return SignatureProfileValueOf(strings.ReplaceAll(name, "-", "_"))
}

// String returns the profile name with '_' replaced by '-', mirroring
// Java's toString() override.
func (s SignatureProfile) String() string {
	return strings.ReplaceAll(string(s), "_", "-")
}
