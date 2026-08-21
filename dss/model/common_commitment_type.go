// Ported from dss-model/.../CommonCommitmentType.java (DSS 6.5.RC1).
package model

import "github.com/ryftcore/dss-go/dss/enumerations"

// CommonCommitmentType provides a basic implementation of
// enumerations.CommitmentType, allowing creation of a customized
// CommitmentType signed property.
type CommonCommitmentType struct {
	CommonObjectIdentifier

	// signedDataObjects defines signed data objects referenced by the
	// current CommitmentType. Use: OPTIONAL.
	signedDataObjects []string

	// commitmentTypeQualifiers defines custom CommitmentTypeQualifiers
	// list. Use: OPTIONAL.
	commitmentTypeQualifiers []*CommitmentQualifier
}

var _ enumerations.CommitmentType = (*CommonCommitmentType)(nil)

// NewCommonCommitmentType instantiates the object with null values. Ports
// the default constructor.
func NewCommonCommitmentType() *CommonCommitmentType {
	return &CommonCommitmentType{}
}

// SignedDataObjects gets references to signed data objects for the
// current CommitmentType.
func (c *CommonCommitmentType) SignedDataObjects() []string { return c.signedDataObjects }

// SetSignedDataObjects sets signed data objects referenced by the current
// CommitmentType.
//
// When CommitmentType is made for a subset of signed data objects, each
// element of the slice shall refer to one ds:Reference element within the
// ds:SignedInfo element or within a signed ds:Manifest element. When
// CommitmentType is made for all signed data objects, the slice shall be:
//   - empty (default), then AllSignedDataObjects element will be created; or
//   - contain references to all signed data objects (one ObjectReference
//     will be created for each).
//
// Use: OPTIONAL (XAdES only).
func (c *CommonCommitmentType) SetSignedDataObjects(signedDataObjects ...string) {
	c.signedDataObjects = signedDataObjects
}

// CommitmentTypeQualifiers gets the custom CommitmentTypeQualifiers list.
func (c *CommonCommitmentType) CommitmentTypeQualifiers() []*CommitmentQualifier {
	return c.commitmentTypeQualifiers
}

// SetCommitmentTypeQualifiers sets the custom CommitmentTypeQualifiers
// list. Use: OPTIONAL.
func (c *CommonCommitmentType) SetCommitmentTypeQualifiers(commitmentTypeQualifiers ...*CommitmentQualifier) {
	c.commitmentTypeQualifiers = commitmentTypeQualifiers
}
