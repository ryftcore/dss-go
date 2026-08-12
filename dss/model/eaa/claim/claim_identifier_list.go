// Ported from dss-model/.../claim/ClaimIdentifierList.java (DSS 6.5.RC1).
package claim

// ClaimIdentifierList represents an EAA Identifier List claim.
type ClaimIdentifierList interface {
	ClaimRevocationList

	// Identifier gets the EAA's Identifier's List identifier value, when
	// present. Ports ClaimIdentifierList#getIdentifier.
	Identifier() *ClaimByteString
}
