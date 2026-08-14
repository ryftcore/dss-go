// Ported from dss-model/.../claim/ClaimDrivingPrivileges.java (DSS 6.5.RC1).
package claim

// ClaimDrivingPrivileges represents a "driving_privileges" claim as
// defined in "7.2.4 Categories of vehicles/restrictions/conditions" of
// ISO/IEC 18013-5.
type ClaimDrivingPrivileges interface {
	Claim

	// DrivingPrivileges gets a list of driving privilege claims. Ports
	// ClaimDrivingPrivileges#getDrivingPrivileges.
	DrivingPrivileges() []ClaimDrivingPrivilege
}
