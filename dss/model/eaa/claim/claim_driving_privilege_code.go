// Ported from dss-model/.../claim/ClaimDrivingPrivilegeCode.java (DSS 6.5.RC1).
package claim

// DrivingPrivilegeCode represents a single Code element entry of the
// "codes" array as defined in the "7.2.4 Categories of
// vehicles/restrictions/conditions" of ISO/IEC 18013-5.
type DrivingPrivilegeCode interface {
	Claim

	// Code gets a code as per ISO/IEC 18013-2 Annex A. Ports
	// ClaimDrivingPrivilegeCode#getCode.
	Code() *String

	// Sign gets a sign as per ISO/IEC 18013-2 Annex A. Ports
	// ClaimDrivingPrivilegeCode#getSign.
	Sign() *String

	// Value gets a value as per ISO/IEC 18013-2 Annex A. Ports
	// ClaimDrivingPrivilegeCode#getValue.
	Value() *String
}
