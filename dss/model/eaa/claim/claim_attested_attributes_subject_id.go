// Ported from dss-model/.../claim/ClaimAttestedAttributesSubjectId.java (DSS 6.5.RC1).
package claim

// ClaimAttestedAttributesSubjectId univocally identifies the attribute
// subject.
type ClaimAttestedAttributesSubjectId interface {
	Claim

	// FamilyName gets the family name of the attribute subject. Ports
	// ClaimAttestedAttributesSubjectId#getFamilyName.
	FamilyName() *ClaimString

	// GivenName gets the given name of the attribute subject. Ports
	// ClaimAttestedAttributesSubjectId#getGivenName.
	GivenName() *ClaimString

	// DocumentNumber gets the number of the personal identification data
	// assigned to the attribute subject. Ports
	// ClaimAttestedAttributesSubjectId#getDocumentNumber.
	DocumentNumber() *ClaimString
}
