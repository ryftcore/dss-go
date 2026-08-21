// Ported from dss-model/.../eaa/DisclosureValidation.java (DSS 6.5.RC1).
package eaa

import (
	"reflect"

	"github.com/ryftcore/dss-go/dss/model"
	"github.com/ryftcore/dss-go/dss/model/eaa/claim"
)

// DisclosureValidation represents a validation result of a selectable
// disclosure provided with presentation of Electronic Attestation of
// Attributes.
type DisclosureValidation struct {
	*model.ReferenceValidation

	// disclosure is the disclosure object, when applicable.
	disclosure *ValidationDisclosure

	// namespace is the namespace of the selective disclosure (mdoc).
	namespace string

	// digestId is the unique identifier of the selective disclosure
	// (mdoc). Nil-able (mirrors Java's Long).
	digestId *int64
}

// NewDisclosureValidation ports the default constructor.
func NewDisclosureValidation() *DisclosureValidation {
	return &DisclosureValidation{ReferenceValidation: model.NewReferenceValidation()}
}

// NewDisclosureValidationWithDisclosure ports the constructor with a
// provided disclosure. Panics if disclosure is nil (Java
// Objects.requireNonNull("Disclosure cannot be null!")).
func NewDisclosureValidationWithDisclosure(disclosure *ValidationDisclosure) *DisclosureValidation {
	if disclosure == nil {
		panic("Disclosure cannot be null!")
	}
	return &DisclosureValidation{
		ReferenceValidation: model.NewReferenceValidation(),
		disclosure:          disclosure,
	}
}

// Disclosure gets disclosure when applicable. Ports
// DisclosureValidation#getDisclosure.
func (d *DisclosureValidation) Disclosure() *ValidationDisclosure { return d.disclosure }

// ClaimName gets the provided disclosure name. Ports
// DisclosureValidation#getClaimName.
func (d *DisclosureValidation) ClaimName() string {
	if d.disclosure != nil {
		return d.disclosure.Name()
	}
	return ""
}

// Value gets the original provided disclosure claim value. Ports
// DisclosureValidation#getValue.
func (d *DisclosureValidation) Value() claim.Claim {
	if d.disclosure != nil {
		return d.disclosure.Claim
	}
	return nil
}

// Namespace gets disclosure's namespace (mdoc only). Ports
// DisclosureValidation#getNamespace.
func (d *DisclosureValidation) Namespace() string {
	if d.namespace != "" {
		return d.namespace
	}
	if d.disclosure != nil {
		return d.disclosure.Namespace()
	}
	return ""
}

// SetNamespace sets the namespace of the selective disclosure (mdoc
// only). Ports DisclosureValidation#setNamespace.
func (d *DisclosureValidation) SetNamespace(namespace string) { d.namespace = namespace }

// DigestId gets the digest Id (mdoc only). Ports
// DisclosureValidation#getDigestId.
func (d *DisclosureValidation) DigestId() *int64 {
	if d.digestId != nil {
		return d.digestId
	}
	if d.disclosure != nil {
		return d.disclosure.DigestId()
	}
	return nil
}

// SetDigestId sets the unique identifier of the selective disclosure
// (mdoc only). Ports DisclosureValidation#setDigestId.
func (d *DisclosureValidation) SetDigestId(digestId *int64) { d.digestId = digestId }

// Equals ports DisclosureValidation#equals: NOTE this does not compare
// the embedded ReferenceValidation fields (upstream's override does not
// call super.equals() either), only disclosure/namespace/digestId.
func (d *DisclosureValidation) Equals(other *DisclosureValidation) bool {
	if d == other {
		return true
	}
	if other == nil {
		return false
	}
	digestIdEqual := (d.digestId == nil) == (other.digestId == nil)
	if digestIdEqual && d.digestId != nil {
		digestIdEqual = *d.digestId == *other.digestId
	}
	return reflect.DeepEqual(d.disclosure, other.disclosure) &&
		d.namespace == other.namespace &&
		digestIdEqual
}
