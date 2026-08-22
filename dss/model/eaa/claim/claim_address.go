// Ported from dss-model/.../claim/ClaimAddress.java (DSS 6.5.RC1).
package claim

// Address represents an "address" claim.
type Address interface {
	Claim

	// PostalAddress gets the user's full postal or mailing address,
	// formatted, when present. Ports ClaimAddress#getPostalAddress.
	PostalAddress() *String

	// StreetAddress gets the user's street address, when present. The
	// component may include a house number, street name, Post Office
	// Box, and multi-line extended street address information. Ports
	// ClaimAddress#getStreetAddress.
	StreetAddress() *String

	// City gets the user's city or locality address, when present. Ports
	// ClaimAddress#getCity.
	City() *String

	// StateOrProvince gets the user's state or region address, when
	// present. Ports ClaimAddress#getStateOrProvince.
	StateOrProvince() *String

	// PostalCode gets the user's zip code or postal code address, when
	// present. Ports ClaimAddress#getPostalCode.
	PostalCode() *String

	// Country gets the user's country address, when present. Ports
	// ClaimAddress#getCountry.
	Country() *String

	// HouseNumber gets the house number where the user to whom the person
	// identification data relates currently resides, including any affix
	// or suffix, when present (ARF PID Rulebook claim). Ports
	// ClaimAddress#getHouseNumber.
	HouseNumber() *String
}
