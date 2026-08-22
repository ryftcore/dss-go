// Ported from dss-model/.../claim/ClaimBirthDate.java (DSS 6.5.RC1).
package claim

// BirthDate represents an ISO/IEC 23220-2:2026 "6.3.1.3 Date of birth
// structure" data element.
type BirthDate interface {
	Claim

	// BirthDate gets day, month and year on which the holder was born.
	// Unknown parts (i.e., year, month, day) are masked with 1. Ports
	// ClaimBirthDate#getBirthDate.
	BirthDate() *Date

	// ApproximateMask gets an 8 digit flag to denote the location of the
	// mask in YYYYMMDD format. 1 denotes mask. Ports
	// ClaimBirthDate#getApproximateMask.
	ApproximateMask() *String
}
