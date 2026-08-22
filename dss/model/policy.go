// Ported from dss-model/.../Policy.java (DSS 6.5.RC1).
package model

import (
	"fmt"
	"reflect"

	"github.com/ryftcore/dss-go/dss/enumerations"
)

// Policy allows defining the signature policy.
//
// UserNotice and SpDocSpecification are outside this manifest; assumed to
// already exist in this package.
type Policy struct {
	// id is the Id of the Policy.
	id string

	// qualifier is the qualifier attribute for XAdES Identifier.
	qualifier enumerations.ObjectIdentifierQualifier

	// description is the Policy description.
	description string

	// documentationReferences is the array of documentation references
	// (used in XAdES/JAdES).
	documentationReferences []string

	// digestAlgorithm is the digest algorithm used to compute the digest.
	digestAlgorithm enumerations.DigestAlgorithm

	// digestValue is the computed digest value.
	digestValue []byte

	// spUri is the Policy URI qualifier.
	spUri string

	// userNotice is the Policy UserNotice qualifier.
	userNotice *UserNotice

	// spDocSpecification is the Policy Document Specification
	// qualifier.
	spDocSpecification *SpDocSpecification

	// hashAsInTechnicalSpecification is used only in JAdES, to indicate
	// that the digest of the signature policy document has been computed
	// as specified in a technical specification.
	hashAsInTechnicalSpecification bool
}

// NewPolicy instantiates the object with null values. Ports the empty
// constructor.
func NewPolicy() *Policy {
	return &Policy{}
}

// Id gets the signature policy (EPES) id.
func (p *Policy) Id() string { return p.id }

// SetId sets the signature policy (EPES) id.
func (p *Policy) SetId(id string) { p.id = id }

// Qualifier gets the identifier qualifier.
func (p *Policy) Qualifier() enumerations.ObjectIdentifierQualifier { return p.qualifier }

// SetQualifier sets the identifier qualifier (used in XAdES only).
func (p *Policy) SetQualifier(qualifier enumerations.ObjectIdentifierQualifier) {
	p.qualifier = qualifier
}

// Description gets the signature policy description.
func (p *Policy) Description() string { return p.description }

// SetDescription sets the signature policy description.
func (p *Policy) SetDescription(description string) { p.description = description }

// DocumentationReferences gets the signature policy documentation
// references.
func (p *Policy) DocumentationReferences() []string { return p.documentationReferences }

// SetDocumentationReferences sets a list of signature documentation
// references.
func (p *Policy) SetDocumentationReferences(documentationReferences ...string) {
	p.documentationReferences = documentationReferences
}

// DigestAlgorithm returns the hash algorithm for the signature policy.
func (p *Policy) DigestAlgorithm() enumerations.DigestAlgorithm { return p.digestAlgorithm }

// SetDigestAlgorithm sets the hash algorithm for the explicit signature
// policy.
func (p *Policy) SetDigestAlgorithm(digestAlgorithm enumerations.DigestAlgorithm) {
	p.digestAlgorithm = digestAlgorithm
}

// DigestValue gets the hash value of the explicit signature policy.
func (p *Policy) DigestValue() []byte { return p.digestValue }

// SetDigestValue sets the hash value of the implicit signature policy.
func (p *Policy) SetDigestValue(digestValue []byte) { p.digestValue = digestValue }

// Spuri gets the SP URI (signature policy URI) qualifier.
func (p *Policy) Spuri() string { return p.spUri }

// SetSpuri sets the SP URI (signature policy URI) qualifier.
func (p *Policy) SetSpuri(spUri string) { p.spUri = spUri }

// UserNotice gets the SP UserNotice qualifier.
func (p *Policy) UserNotice() *UserNotice { return p.userNotice }

// SetUserNotice sets the SP UserNotice qualifier.
func (p *Policy) SetUserNotice(userNotice *UserNotice) { p.userNotice = userNotice }

// SpDocSpecification gets the SP Document Specification qualifier.
func (p *Policy) SpDocSpecification() *SpDocSpecification { return p.spDocSpecification }

// SetSpDocSpecification sets the SP Document Specification qualifier
// identifying the technical specification that defines the syntax used
// for producing the signature policy document.
func (p *Policy) SetSpDocSpecification(spDocSpecification *SpDocSpecification) {
	p.spDocSpecification = spDocSpecification
}

// IsHashAsInTechnicalSpecification gets if the digests of the signature
// policy has been computed as in a technical specification.
func (p *Policy) IsHashAsInTechnicalSpecification() bool { return p.hashAsInTechnicalSpecification }

// SetHashAsInTechnicalSpecification sets if the digests of the signature
// policy has been computed as in a technical specification. If the
// property is set to FALSE, digest of the signature policy is computed in
// a default way (on the policy file).
//
// NOTE: The property is used only in JAdES.
//
// Use SetSpDocSpecification to provide the technical specification.
func (p *Policy) SetHashAsInTechnicalSpecification(hashAsInTechnicalSpecification bool) {
	p.hashAsInTechnicalSpecification = hashAsInTechnicalSpecification
}

// IsEmpty checks if the object's data is not filled.
func (p *Policy) IsEmpty() bool {
	if p.id != "" {
		return false
	}
	if p.qualifier != "" {
		return false
	}
	if p.description != "" {
		return false
	}
	if len(p.documentationReferences) != 0 {
		return false
	}
	if p.digestAlgorithm != "" {
		return false
	}
	if len(p.digestValue) != 0 {
		return false
	}
	if p.spUri != "" {
		return false
	}
	if p.userNotice != nil && !p.userNotice.IsEmpty() {
		return false
	}
	if p.spDocSpecification != nil && p.spDocSpecification.Id() != "" {
		return false
	}
	if p.hashAsInTechnicalSpecification {
		return false
	}
	return true
}

// IsSPQualifierPresent checks if there is a definition at least for one
// signature policy qualifier.
func (p *Policy) IsSPQualifierPresent() bool {
	if p.spUri != "" {
		return true
	}
	if p.userNotice != nil && !p.userNotice.IsEmpty() {
		return true
	}
	if p.spDocSpecification != nil && p.spDocSpecification.Id() != "" {
		return true
	}
	return false
}

// Equals ports Policy#equals.
func (p *Policy) Equals(other *Policy) bool {
	if p == other {
		return true
	}
	if other == nil {
		return false
	}
	if p.hashAsInTechnicalSpecification != other.hashAsInTechnicalSpecification {
		return false
	}
	if p.id != other.id {
		return false
	}
	if p.qualifier != other.qualifier {
		return false
	}
	if p.description != other.description {
		return false
	}
	if !reflect.DeepEqual(p.documentationReferences, other.documentationReferences) {
		return false
	}
	if p.digestAlgorithm != other.digestAlgorithm {
		return false
	}
	if !reflect.DeepEqual(p.digestValue, other.digestValue) {
		return false
	}
	if p.spUri != other.spUri {
		return false
	}
	if !policyUserNoticeEqual(p.userNotice, other.userNotice) {
		return false
	}
	if !policySpDocSpecificationEqual(p.spDocSpecification, other.spDocSpecification) {
		return false
	}
	return true
}

// policyUserNoticeEqual compares two possibly-nil *UserNotice values.
func policyUserNoticeEqual(a, b *UserNotice) bool {
	if a == nil || b == nil {
		return a == b
	}
	return a.Equals(b)
}

// policySpDocSpecificationEqual compares two possibly-nil
// *SpDocSpecification values.
func policySpDocSpecificationEqual(a, b *SpDocSpecification) bool {
	if a == nil || b == nil {
		return a == b
	}
	return a.Equals(b)
}

// String ports Policy#toString.
func (p *Policy) String() string {
	return fmt.Sprintf("Policy {id='%s', qualifier=%v, description='%s', documentationReferences=%v, digestAlgorithm=%v, digestValue=%v, spUri='%s', userNotice=%v, spDocSpecification='%v', hashAsInTechnicalSpecification=%v}",
		p.id, p.qualifier, p.description, p.documentationReferences, p.digestAlgorithm, p.digestValue,
		p.spUri, p.userNotice, p.spDocSpecification, p.hashAsInTechnicalSpecification)
}
