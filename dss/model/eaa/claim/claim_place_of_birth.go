// Ported from dss-model/.../claim/ClaimPlaceOfBirth.java (DSS 6.5.RC1).
package claim

// ClaimPlaceOfBirth represents a user's place of birth claim.
type ClaimPlaceOfBirth interface {
	Claim

	// Country gets user's country of birth, represented by 2-letter ISO
	// 3116-1 code, when present. Ports ClaimPlaceOfBirth#getCountry.
	Country() *ClaimString

	// StateOrProvince gets user's state, province, prefecture, or region
	// of birth, when present. Ports
	// ClaimPlaceOfBirth#getStateOrProvince.
	StateOrProvince() *ClaimString

	// City gets user's city or locality of birth, when present. Ports
	// ClaimPlaceOfBirth#getCity.
	City() *ClaimString
}
