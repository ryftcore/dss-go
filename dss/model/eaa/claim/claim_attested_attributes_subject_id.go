// Ported from dss-model/.../claim/ClaimAttestedAttributesSubjectId.java (DSS 6.5.RC1).
package claim

// AttestedAttributesSubjectId univocally identifies the attribute
// subject.
type AttestedAttributesSubjectId interface {
	Claim

	// FamilyName gets the family name of the attribute subject. Ports
	// ClaimAttestedAttributesSubjectId#getFamilyName.
	FamilyName() *String

	// GivenName gets the given name of the attribute subject. Ports
	// ClaimAttestedAttributesSubjectId#getGivenName.
	GivenName() *String

	// DocumentNumber gets the number of the personal identification data
	// assigned to the attribute subject. Ports
	// ClaimAttestedAttributesSubjectId#getDocumentNumber.
	DocumentNumber() *String
}
