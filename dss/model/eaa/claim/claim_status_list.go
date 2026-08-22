// Ported from dss-model/.../claim/ClaimStatusList.java (DSS 6.5.RC1).
package claim

// StatusList represents an EAA Status List claim.
type StatusList interface {
	RevocationList

	// Index gets the EAA's Status index value, when present. Ports
	// ClaimStatusList#getIndex.
	Index() *Number
}
