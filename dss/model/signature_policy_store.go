// Ported from dss-model/.../SignaturePolicyStore.java (DSS 6.5.RC1).
package model

// SignaturePolicyStore represents the SignaturePolicyStore.
//
// SpDocSpecification and DSSDocument are outside this manifest; assumed
// to already exist in this package.
type SignaturePolicyStore struct {
	// id is an optional ID.
	id string

	// spDocSpecification shall identify the technical specification that
	// defines the syntax used for producing the signature policy
	// document.
	spDocSpecification *SpDocSpecification

	// signaturePolicyContent shall contain the base-64 encoded signature
	// policy.
	signaturePolicyContent DSSDocument

	// sigPolDocLocalURI shall have as value the URI referencing a local
	// store where the present document can be retrieved.
	sigPolDocLocalURI string
}

// NewSignaturePolicyStore instantiates the object with null values. Ports
// the default constructor.
func NewSignaturePolicyStore() *SignaturePolicyStore {
	return &SignaturePolicyStore{}
}

// Id gets the optional Id.
func (s *SignaturePolicyStore) Id() string { return s.id }

// SetId sets the optional Id.
func (s *SignaturePolicyStore) SetId(id string) { s.id = id }

// SpDocSpecification gets the SpDocSpecification content.
func (s *SignaturePolicyStore) SpDocSpecification() *SpDocSpecification { return s.spDocSpecification }

// SetSpDocSpecification sets the SpDocSpecification.
func (s *SignaturePolicyStore) SetSpDocSpecification(spDocSpecification *SpDocSpecification) {
	s.spDocSpecification = spDocSpecification
}

// SignaturePolicyContent gets the policy store content.
func (s *SignaturePolicyStore) SignaturePolicyContent() DSSDocument { return s.signaturePolicyContent }

// SetSignaturePolicyContent sets the policy store content.
//
// NOTE: one of SignaturePolicyContent or SigPolDocLocalURI shall be used.
func (s *SignaturePolicyStore) SetSignaturePolicyContent(signaturePolicyContent DSSDocument) {
	s.signaturePolicyContent = signaturePolicyContent
}

// SigPolDocLocalURI gets the SigPolDocLocalURI element value.
func (s *SignaturePolicyStore) SigPolDocLocalURI() string { return s.sigPolDocLocalURI }

// SetSigPolDocLocalURI sets the SigPolDocLocalURI element value, defining
// the local URI where the policy document can be retrieved.
//
// NOTE: one of SignaturePolicyContent or SigPolDocLocalURI shall be used.
func (s *SignaturePolicyStore) SetSigPolDocLocalURI(sigPolDocLocalURI string) {
	s.sigPolDocLocalURI = sigPolDocLocalURI
}
