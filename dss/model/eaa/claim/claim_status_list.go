// Ported from dss-model/.../claim/ClaimStatusList.java (DSS 6.5.RC1).
package claim

// ClaimStatusList represents an EAA Status List claim.
type ClaimStatusList interface {
	ClaimRevocationList

	// Index gets the EAA's Status index value, when present. Ports
	// ClaimStatusList#getIndex.
	Index() *ClaimNumber
}
