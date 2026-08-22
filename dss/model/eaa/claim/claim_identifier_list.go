// Ported from dss-model/.../claim/ClaimIdentifierList.java (DSS 6.5.RC1).
package claim

// IdentifierList represents an EAA Identifier List claim.
type IdentifierList interface {
	RevocationList

	// Identifier gets the EAA's Identifier's List identifier value, when
	// present. Ports ClaimIdentifierList#getIdentifier.
	Identifier() *ByteString
}
