// Ported from dss-xades/src/main/java/eu/europa/esig/dss/xades/XAdESProfileParameters.java (DSS 6.5.RC1).
//
// Java extends document.signature.ProfileParameters, adding the XAdES-specific profile/builder/
// operationKind/references fields. The Go port embeds document.ProfileParameters by value,
// following the XAdESCounterSignatureParameters/XAdESSignatureParameters precedent of embedding
// the base struct rather than trying to reproduce Java's field-level polymorphism (see
// xades_signature_parameters.go's header for the fuller discussion of that limitation and how
// this type's embedding is exactly what lets XAdESSignatureParameters.GetContext() and the base
// document.AbstractSignatureParameters.GetContext() answer the same underlying values for the
// fields both know about).
//
// java.io.Serializable and serialVersionUID are dropped; hashCode() has no Go counterpart.
package xades

import (
	"fmt"

	"github.com/ryftcore/dss-go/dss/document"
	"github.com/ryftcore/dss-go/dss/enumerations"
)

// XAdESProfileParameters is used to accelerate the signature creation process for XAdES.
type XAdESProfileParameters struct {
	document.ProfileParameters

	// profile is the XAdES creation profile.
	profile XAdESSignatureProfile

	// builder is the builder used to create the signature structure.
	builder SignatureBuilder

	// operationKind indicates the type of the operation to be done.
	operationKind enumerations.SigningOperation

	// references is the list of references created by a reference builder.
	references []*DSSReference
}

// NewXAdESProfileParameters is the default constructor.
func NewXAdESProfileParameters() *XAdESProfileParameters {
	return &XAdESProfileParameters{ProfileParameters: *document.NewProfileParameters()}
}

// Profile returns the current Profile used to generate the signature or its extension. Ports
// getProfile().
func (p *XAdESProfileParameters) Profile() XAdESSignatureProfile {
	return p.profile
}

// SetProfile sets the current Profile used to generate the signature or its extension. Ports
// setProfile(XAdESSignatureProfile).
func (p *XAdESProfileParameters) SetProfile(profile XAdESSignatureProfile) {
	p.profile = profile
}

// Builder gets the signature builder. Ports getBuilder().
func (p *XAdESProfileParameters) Builder() SignatureBuilder {
	return p.builder
}

// SetBuilder sets the signature builder. Ports setBuilder(SignatureBuilder).
func (p *XAdESProfileParameters) SetBuilder(builder SignatureBuilder) {
	p.builder = builder
}

// OperationKind gets the current operation type. Ports getOperationKind().
func (p *XAdESProfileParameters) OperationKind() enumerations.SigningOperation {
	return p.operationKind
}

// SetOperationKind sets the operation kind. Ports setOperationKind(SigningOperation).
func (p *XAdESProfileParameters) SetOperationKind(operationKind enumerations.SigningOperation) {
	p.operationKind = operationKind
}

// References returns a list of references to be incorporated to the signature. Ports
// getReferences().
func (p *XAdESProfileParameters) References() []*DSSReference {
	return p.references
}

// SetReferences sets a list of references to be incorporated into the signature. Ports
// setReferences(List<DSSReference>).
func (p *XAdESProfileParameters) SetReferences(references []*DSSReference) {
	p.references = references
}

// String ports toString(). Java's toString does not call super.toString(), unlike most
// classes in this port; reproduced verbatim.
func (p *XAdESProfileParameters) String() string {
	return fmt.Sprintf("XAdESProfileParameters{profile=%v, builder=%v, operationKind=%v, references=%v}",
		p.profile, p.builder, p.operationKind, p.references)
}

// Equals ports equals(Object).
func (p *XAdESProfileParameters) Equals(other *XAdESProfileParameters) bool {
	if p == other {
		return true
	}
	if other == nil {
		return false
	}
	if !p.ProfileParameters.Equals(&other.ProfileParameters) {
		return false
	}
	if p.profile != other.profile || p.builder != other.builder || p.operationKind != other.operationKind {
		return false
	}
	return xadesProfileParametersReferencesEqual(p.references, other.references)
}

// xadesProfileParametersReferencesEqual ports Objects.equals(references, ...).
func xadesProfileParametersReferencesEqual(a, b []*DSSReference) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
