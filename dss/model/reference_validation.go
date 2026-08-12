// Ported from dss-model/.../ReferenceValidation.java (DSS 6.5.RC1).
package model

import "github.com/utain/esig/dss/enumerations"

// ReferenceValidation stores individual reference validations.
//
// For XAdES, that means reference tag(s) validation.
// For CAdES, that means message-digest validation.
type ReferenceValidation struct {
	// typ is the type of the Reference.
	typ enumerations.DigestMatcherType

	// found is whether the pointed reference is found.
	found bool

	// intact is whether the pointed reference is intact.
	intact bool

	// digest is the digest value embedded in the reference element.
	digest Digest

	// id is the unique identifier of the reference. (E.g. for XAdES:
	// reference Id attribute value).
	id string

	// uri is the reference to the original document. (E.g. for XAdES:
	// reference URI attribute value).
	uri string

	// dataObjectReferences lists data object references covered by the
	// current reference validation (e.g. used in JAdES for
	// SigDMechanism.OBJECT_ID_BY_URI).
	dataObjectReferences []string

	// document is the matching document (when applicable).
	document DSSDocument

	// transforms lists used transforms to compute digest of the reference.
	transforms []string

	// isDuplicated is whether the reference points to more than one
	// element.
	isDuplicated bool

	// errorMessages lists errors occurred during the reference validation.
	errorMessages []string

	// dependentReferenceValidations lists dependent ReferenceValidations
	// (used in case of manifest type for manifest entries).
	dependentReferenceValidations []*ReferenceValidation
}

// NewReferenceValidation instantiates the object with null values (an
// empty errorMessages list, per Java field initializer). Ports the default
// constructor.
func NewReferenceValidation() *ReferenceValidation {
	return &ReferenceValidation{errorMessages: []string{}}
}

// Type returns the type of the validated reference.
func (r *ReferenceValidation) Type() enumerations.DigestMatcherType { return r.typ }

// SetType sets the type of the reference.
func (r *ReferenceValidation) SetType(typ enumerations.DigestMatcherType) { r.typ = typ }

// IsFound gets if the reference's data has been found.
func (r *ReferenceValidation) IsFound() bool { return r.found }

// SetFound sets if the reference's data has been found.
func (r *ReferenceValidation) SetFound(found bool) { r.found = found }

// IsIntact gets if the digest of a referenced document matches the one
// defined in the reference.
func (r *ReferenceValidation) IsIntact() bool { return r.intact }

// SetIntact sets if the digest value of a referenced document matches.
func (r *ReferenceValidation) SetIntact(intact bool) { r.intact = intact }

// Digest gets the incorporated Digest.
func (r *ReferenceValidation) Digest() Digest { return r.digest }

// SetDigest sets the reference's Digest.
func (r *ReferenceValidation) SetDigest(digest Digest) { r.digest = digest }

// Id gets the unique identifier of a reference. (E.g. for XAdES: reference
// Id attribute value).
func (r *ReferenceValidation) Id() string { return r.id }

// SetId sets the unique identifier of a reference. (E.g. for XAdES:
// reference Id attribute value).
func (r *ReferenceValidation) SetId(id string) { r.id = id }

// Uri gets the reference to the original document. (E.g. for XAdES:
// reference URI attribute value).
func (r *ReferenceValidation) Uri() string { return r.uri }

// SetUri sets the reference to the original document. (E.g. for XAdES:
// reference URI attribute value).
func (r *ReferenceValidation) SetUri(uri string) { r.uri = uri }

// DataObjectReferences gets extracted data object reference URIs, covered
// by the current reference. Example: JAdES signatures with
// SigDMechanism.OBJECT_ID_BY_URI.
func (r *ReferenceValidation) DataObjectReferences() []string { return r.dataObjectReferences }

// SetDataObjectReferences sets extracted data object reference URIs,
// covered by the current reference.
func (r *ReferenceValidation) SetDataObjectReferences(dataObjectReferences []string) {
	r.dataObjectReferences = dataObjectReferences
}

// Document gets the matching document.
func (r *ReferenceValidation) Document() DSSDocument { return r.document }

// SetDocument sets the matching document.
func (r *ReferenceValidation) SetDocument(document DSSDocument) { r.document = document }

// TransformationNames returns the list of transformations contained in the
// reference.
func (r *ReferenceValidation) TransformationNames() []string { return r.transforms }

// SetTransformationNames sets the list of transforms for the reference.
func (r *ReferenceValidation) SetTransformationNames(transforms []string) { r.transforms = transforms }

// IsDuplicated returns if the referenced data is ambiguous.
func (r *ReferenceValidation) IsDuplicated() bool { return r.isDuplicated }

// SetDuplicated sets if the referenced data is ambiguous.
func (r *ReferenceValidation) SetDuplicated(isDuplicated bool) { r.isDuplicated = isDuplicated }

// DependentValidations returns the list of dependent validations from this
// ReferenceValidation, lazily initializing it. Note: used to contain
// manifest entries.
func (r *ReferenceValidation) DependentValidations() []*ReferenceValidation {
	if r.dependentReferenceValidations == nil {
		r.dependentReferenceValidations = []*ReferenceValidation{}
	}
	return r.dependentReferenceValidations
}

// ErrorMessages gets error messages occurred during the reference
// validation.
func (r *ReferenceValidation) ErrorMessages() []string { return r.errorMessages }

// SetErrorMessages sets error messages occurred during the reference
// validation.
func (r *ReferenceValidation) SetErrorMessages(errorMessages []string) {
	r.errorMessages = errorMessages
}
