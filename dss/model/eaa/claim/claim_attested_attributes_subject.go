// Ported from dss-model/.../claim/ClaimAttestedAttributesSubject.java (DSS 6.5.RC1).
package claim

// AttestedAttributesSubject associates a set of attributes to one
// entity different than the EAA subject.
type AttestedAttributesSubject interface {
	Claim

	// SubjectId gets the identifier of the attribute subject, which shall
	// associate the attributes to this attribute subject. Ports
	// ClaimAttestedAttributesSubject#getSubjectId.
	SubjectId() Claim

	// SubjectPseudonym gets the pseudonym of an attribute subject which
	// shall associate the attributes to this attribute subject. Ports
	// ClaimAttestedAttributesSubject#getSubjectPseudonym.
	SubjectPseudonym() *String

	// Attributes gets the attributes associated to the attribute subject
	// whose identifier appears in the sub_id member or whose pseudonym
	// appears in the sub_aka member. Ports
	// ClaimAttestedAttributesSubject#getAttributes.
	Attributes() *Array
}
