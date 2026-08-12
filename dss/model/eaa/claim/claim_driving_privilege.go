// Ported from dss-model/.../claim/ClaimDrivingPrivilege.java (DSS 6.5.RC1).
package claim

// ClaimDrivingPrivilege represents a single item of the
// "driving_privileges" claim array, as defined in "7.2.4 Categories of
// vehicles/restrictions/conditions" of ISO/IEC 18013-5.
type ClaimDrivingPrivilege interface {
	Claim

	// VehicleCategoryCode gets a vehicle category code as per ISO/IEC
	// 18013-1 Annex B. Ports
	// ClaimDrivingPrivilege#getVehicleCategoryCode.
	VehicleCategoryCode() *ClaimString

	// IssueDate gets a date of issue encoded as full-date. Ports
	// ClaimDrivingPrivilege#getIssueDate.
	IssueDate() *ClaimDate

	// ExpiryDate gets a date of expiry encoded as full-date. Ports
	// ClaimDrivingPrivilege#getExpiryDate.
	ExpiryDate() *ClaimDate

	// Codes gets an array of code info. Ports
	// ClaimDrivingPrivilege#getCodes.
	Codes() ClaimDrivingPrivilegeCodes
}
