// Ported from dss-model/.../claim/ClaimByteString.java (DSS 6.5.RC1).
package claim

import (
	"bytes"
	"encoding/base64"
)

// ClaimByteString represents a byte array (selectively) disclosable claim.
type ClaimByteString struct {
	AbstractClaim

	// value is the byte array value of the claim.
	value []byte
}

// NewClaimByteString ports the default constructor.
func NewClaimByteString(value []byte) *ClaimByteString {
	return NewClaimByteStringWithName("", value)
}

// NewClaimByteStringWithName ports the constructor with claim name
// provided.
func NewClaimByteStringWithName(name string, value []byte) *ClaimByteString {
	return NewClaimByteStringWithDisclosable(name, value, false)
}

// NewClaimByteStringWithDisclosable ports the constructor with claim name
// and selectively disclosable status provided.
func NewClaimByteStringWithDisclosable(name string, value []byte, selectivelyDisclosable bool) *ClaimByteString {
	return NewClaimByteStringWithParent(name, value, selectivelyDisclosable, nil)
}

// NewClaimByteStringWithParent ports the constructor with claim name,
// selectively disclosable status and parent claim provided.
func NewClaimByteStringWithParent(name string, value []byte, selectivelyDisclosable bool, parent Claim) *ClaimByteString {
	return NewClaimByteStringFull(name, "", value, selectivelyDisclosable, parent)
}

// NewClaimByteStringFull ports the constructor with claim name, namespace,
// selectively disclosable status and parent claim provided.
func NewClaimByteStringFull(name, namespace string, value []byte, selectivelyDisclosable bool, parent Claim) *ClaimByteString {
	return &ClaimByteString{
		AbstractClaim: NewAbstractClaimFull(name, namespace, selectivelyDisclosable, parent),
		value:         value,
	}
}

// BinaryValue returns the byte array value of the claim.
func (c *ClaimByteString) BinaryValue() []byte { return c.value }

// IsBinaryValueType always returns true.
func (c *ClaimByteString) IsBinaryValueType() bool { return true }

// IsNullOrEmpty ports ClaimByteString#isNullOrEmpty verbatim: NOTE this is
// `value != null` in upstream (the inverse of what the method name
// suggests), copied as-is per the porting rule against "fixing" upstream
// values.
func (c *ClaimByteString) IsNullOrEmpty() bool { return c.value != nil }

// ValueAsString base64-encodes the value. Ports
// ClaimByteString#getValueAsString.
func (c *ClaimByteString) ValueAsString() string {
	return base64.StdEncoding.EncodeToString(c.value)
}

// Equals ports ClaimByteString#equals (including the AbstractClaim
// super.equals() comparison).
func (c *ClaimByteString) Equals(other *ClaimByteString) bool {
	if c == other {
		return true
	}
	if other == nil {
		return false
	}
	if !c.AbstractClaim.Equals(&other.AbstractClaim) {
		return false
	}
	return bytes.Equal(c.value, other.value)
}

// String ports AbstractClaim#toString, inherited by this claim.
func (c *ClaimByteString) String() string { return AbstractClaimString(c) }
