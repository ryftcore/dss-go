// Ported from dss-model/.../claim/ClaimDrivingPrivilegeCodes.java (DSS 6.5.RC1).
package claim

// DrivingPrivilegeCodes represents an array of Code's of the
// DrivingPrivilege element, as defined in "7.2.4 Categories of
// vehicles/restrictions/conditions" of ISO/IEC 18013-5.
type DrivingPrivilegeCodes interface {
	Claim

	// Codes gets a list of codes. Ports
	// ClaimDrivingPrivilegeCodes#getCodes.
	Codes() []DrivingPrivilegeCode
}
