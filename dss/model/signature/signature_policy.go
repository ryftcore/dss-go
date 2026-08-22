// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/signature/SignaturePolicy.java (DSS 6.5.RC1).
package signature

import (
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/model"
)

// SignaturePolicy represents the values of a SignaturePolicy extracted on a signature
// validation.
//
// java.io.Serializable is dropped silently (no Go counterpart).
type SignaturePolicy struct {
	// identifier is the signature policy identifier.
	identifier string

	// description is the policy description.
	description string

	// documentationReferences is the documentation references.
	documentationReferences []string

	// policyContent is the policy content document.
	policyContent model.DSSDocument

	// digest is the digest of the signature policy.
	digest model.Digest

	// hashAsInTechnicalSpecification indicates whether the hash should be computed as
	// specified in a relevant signature specification according to the signature policy
	// format.
	hashAsInTechnicalSpecification bool

	// zeroHash reports whether it is a zero-hash policy.
	zeroHash bool

	// uri is the Signature Policy URI qualifier: a URL where a copy of the signature policy
	// MAY be obtained.
	uri string

	// userNotice is the Signature Policy User Notice qualifier: user notice that should be
	// displayed when the signature is verified.
	userNotice *model.UserNotice

	// docSpecification is the Signature Policy Document Specification qualifier: an
	// identifier of the technical specification that defines the syntax used for producing
	// the signature policy document.
	docSpecification *model.SpDocSpecification

	// validationResult is the validation result of the current signature policy.
	validationResult *SignaturePolicyValidationResult
}

// NewSignaturePolicy is the default constructor for SignaturePolicy: it represents the
// implied policy.
func NewSignaturePolicy() *SignaturePolicy {
	return &SignaturePolicy{identifier: string(enumerations.SignaturePolicyTypeImplicitPolicy)}
}

// NewSignaturePolicyWithIdentifier is the default constructor for SignaturePolicy with the
// policy identifier.
func NewSignaturePolicyWithIdentifier(identifier string) *SignaturePolicy {
	return &SignaturePolicy{identifier: identifier}
}

// Identifier returns the signature policy identifier. Port of getIdentifier().
func (s *SignaturePolicy) Identifier() string {
	return s.identifier
}

// Description gets the description. Port of getDescription().
func (s *SignaturePolicy) Description() string {
	return s.description
}

// SetDescription sets the description (optional). Port of setDescription(String).
func (s *SignaturePolicy) SetDescription(description string) {
	s.description = description
}

// PolicyContent returns a DSSDocument with the signature policy content. Port of
// getPolicyContent().
func (s *SignaturePolicy) PolicyContent() model.DSSDocument {
	return s.policyContent
}

// SetPolicyContent sets the policy document content. Port of setPolicyContent(DSSDocument).
func (s *SignaturePolicy) SetPolicyContent(policyContent model.DSSDocument) {
	s.policyContent = policyContent
}

// Digest gets the Digest. Port of getDigest().
func (s *SignaturePolicy) Digest() model.Digest {
	return s.digest
}

// SetDigest sets the Digest. Port of setDigest(Digest).
func (s *SignaturePolicy) SetDigest(digest model.Digest) {
	s.digest = digest
}

// DocumentationReferences gets the documentation references.
//
// NOTE: optional, used in XAdES.
//
// Port of getDocumentationReferences().
func (s *SignaturePolicy) DocumentationReferences() []string {
	return s.documentationReferences
}

// SetDocumentationReferences sets the documentation references. Port of
// setDocumentationReferences(List<String>).
func (s *SignaturePolicy) SetDocumentationReferences(documentationReferences []string) {
	s.documentationReferences = documentationReferences
}

// TransformsDescription gets a list of Strings describing the 'ds:Transforms' element.
//
// NOTE: XAdES only.
//
// Port of getTransformsDescription(); returns an empty (non-nil) slice by default, matching
// Java's Collections.emptyList().
func (s *SignaturePolicy) TransformsDescription() []string {
	return []string{}
}

// IsZeroHash returns if the policy is a zero-hash (no hash check shall be performed). Port
// of isZeroHash().
func (s *SignaturePolicy) IsZeroHash() bool {
	return s.zeroHash
}

// SetZeroHash sets if the policy is a zero-hash (no hash check shall be performed). Port of
// setZeroHash(boolean).
func (s *SignaturePolicy) SetZeroHash(zeroHash bool) {
	s.zeroHash = zeroHash
}

// IsHashAsInTechnicalSpecification returns if the digest should be computed as specified in
// the relevant technical specification. Port of isHashAsInTechnicalSpecification().
func (s *SignaturePolicy) IsHashAsInTechnicalSpecification() bool {
	return s.hashAsInTechnicalSpecification
}

// SetHashAsInTechnicalSpecification sets whether the digest should be computed as specified
// in a corresponding technical specification. Port of
// setHashAsInTechnicalSpecification(boolean).
func (s *SignaturePolicy) SetHashAsInTechnicalSpecification(hashAsInTechnicalSpecification bool) {
	s.hashAsInTechnicalSpecification = hashAsInTechnicalSpecification
}

// URI returns the signature policy URI (if found), empty when not available. Port of
// getUri().
func (s *SignaturePolicy) URI() string {
	return s.uri
}

// SetURI sets the signature policy URI. Port of setUri(String).
func (s *SignaturePolicy) SetURI(uri string) {
	s.uri = uri
}

// UserNotice gets the user notice that should be displayed when the signature is verified.
// Port of getUserNotice().
func (s *SignaturePolicy) UserNotice() *model.UserNotice {
	return s.userNotice
}

// SetUserNotice sets the user notice that should be displayed when the signature is
// verified. Port of setUserNotice(UserNotice).
func (s *SignaturePolicy) SetUserNotice(userNotice *model.UserNotice) {
	s.userNotice = userNotice
}

// DocSpecification gets the Document Specification Qualifier when present. Port of
// getDocSpecification().
func (s *SignaturePolicy) DocSpecification() *model.SpDocSpecification {
	return s.docSpecification
}

// SetDocSpecification sets the Document Specification qualifier. Port of
// setDocSpecification(SpDocSpecification).
func (s *SignaturePolicy) SetDocSpecification(docSpecification *model.SpDocSpecification) {
	s.docSpecification = docSpecification
}

// ValidationResult gets the validation result of the signature policy. Port of
// getValidationResult().
func (s *SignaturePolicy) ValidationResult() *SignaturePolicyValidationResult {
	return s.validationResult
}

// SetValidationResult sets the signature policy's validation result. Port of
// setValidationResult(SignaturePolicyValidationResult).
func (s *SignaturePolicy) SetValidationResult(validationResult *SignaturePolicyValidationResult) {
	s.validationResult = validationResult
}
