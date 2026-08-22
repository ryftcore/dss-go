// Ported from dss-model/.../claim/ClaimStatus.java (DSS 6.5.RC1).
package claim

// Status represents an EAA Status claim.
type Status interface {
	Claim

	/* Token Status List (TSL) draft-ietf-oauth-status-list-20 */

	// StatusList gets the embedded status_list claim value. Ports
	// ClaimStatus#getStatusList.
	StatusList() StatusList

	// IdentifierList gets the embedded identifier_list claim value.
	// Ports ClaimStatus#getIdentifierList.
	IdentifierList() IdentifierList

	/* ETSI TS 119 472-1 status definition */

	// Index gets the EAA's Status index value, when present. Ports
	// ClaimStatus#getIndex.
	Index() *Number

	// Uri gets the EAA's Status URI value, when present. Ports
	// ClaimStatus#getUri.
	Uri() *String

	// Type gets the EAA's Status type value, when present. Ports
	// ClaimStatus#getType.
	Type() *String

	// Purpose gets the EAA's Status purpose value, when present. Ports
	// ClaimStatus#getPurpose.
	Purpose() *String
}
