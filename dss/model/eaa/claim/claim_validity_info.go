// Ported from dss-model/.../claim/ClaimValidityInfo.java (DSS 6.5.RC1).
package claim

// ValidityInfo represents a structure containing information related
// to the validity of the MSO and its signature. The structure corresponds
// to the definition of "ValidityInfo" per "9.1.2.4 Signing method and
// structure for MSO" of ISO/IEC 18013-5.
type ValidityInfo interface {
	Claim

	// Signed gets the timestamp at which the MSO signature was created.
	// Ports ClaimValidityInfo#getSigned.
	Signed() *Date

	// ValidFrom gets the timestamp before which the MSO is not yet
	// valid. Ports ClaimValidityInfo#getValidFrom.
	ValidFrom() *Date

	// ValidUntil gets the timestamp after which the MSO is no longer
	// valid. Ports ClaimValidityInfo#getValidUntil.
	ValidUntil() *Date

	// ExpectedUpdate gets the timestamp at which the issuing authority
	// infrastructure expects to re-sign the MSO. Ports
	// ClaimValidityInfo#getExpectedUpdate.
	ExpectedUpdate() *Date
}
