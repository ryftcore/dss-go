// Ported from dss-model/.../claim/ClaimStatus.java (DSS 6.5.RC1).
package claim

// ClaimStatus represents an EAA Status claim.
type ClaimStatus interface {
	Claim

	/* Token Status List (TSL) draft-ietf-oauth-status-list-20 */

	// StatusList gets the embedded status_list claim value. Ports
	// ClaimStatus#getStatusList.
	StatusList() ClaimStatusList

	// IdentifierList gets the embedded identifier_list claim value.
	// Ports ClaimStatus#getIdentifierList.
	IdentifierList() ClaimIdentifierList

	/* ETSI TS 119 472-1 status definition */

	// Index gets the EAA's Status index value, when present. Ports
	// ClaimStatus#getIndex.
	Index() *ClaimNumber

	// Uri gets the EAA's Status URI value, when present. Ports
	// ClaimStatus#getUri.
	Uri() *ClaimString

	// Type gets the EAA's Status type value, when present. Ports
	// ClaimStatus#getType.
	Type() *ClaimString

	// Purpose gets the EAA's Status purpose value, when present. Ports
	// ClaimStatus#getPurpose.
	Purpose() *ClaimString
}
