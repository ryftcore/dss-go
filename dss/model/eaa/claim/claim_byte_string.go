// Ported from dss-model/.../claim/ClaimByteString.java (DSS 6.5.RC1).
package claim

import (
	"bytes"
	"encoding/base64"
)

// ByteString represents a byte array (selectively) disclosable claim.
type ByteString struct {
	AbstractClaim

	// value is the byte array value of the claim.
	value []byte
}

// NewClaimByteString ports the default constructor.
func NewClaimByteString(value []byte) *ByteString {
	return NewClaimByteStringWithName("", value)
}

// NewClaimByteStringWithName ports the constructor with claim name
// provided.
func NewClaimByteStringWithName(name string, value []byte) *ByteString {
	return NewClaimByteStringWithDisclosable(name, value, false)
}

// NewClaimByteStringWithDisclosable ports the constructor with claim name
// and selectively disclosable status provided.
func NewClaimByteStringWithDisclosable(name string, value []byte, selectivelyDisclosable bool) *ByteString {
	return NewClaimByteStringWithParent(name, value, selectivelyDisclosable, nil)
}

// NewClaimByteStringWithParent ports the constructor with claim name,
// selectively disclosable status and parent claim provided.
func NewClaimByteStringWithParent(name string, value []byte, selectivelyDisclosable bool, parent Claim) *ByteString {
	return NewClaimByteStringFull(name, "", value, selectivelyDisclosable, parent)
}

// NewClaimByteStringFull ports the constructor with claim name, namespace,
// selectively disclosable status and parent claim provided.
func NewClaimByteStringFull(name, namespace string, value []byte, selectivelyDisclosable bool, parent Claim) *ByteString {
	return &ByteString{
		AbstractClaim: NewAbstractClaimFull(name, namespace, selectivelyDisclosable, parent),
		value:         value,
	}
}

// BinaryValue returns the byte array value of the claim.
func (c *ByteString) BinaryValue() []byte { return c.value }

// IsBinaryValueType always returns true.
func (c *ByteString) IsBinaryValueType() bool { return true }

// IsNullOrEmpty ports ClaimByteString#isNullOrEmpty verbatim: NOTE this is
// `value != null` in upstream (the inverse of what the method name
// suggests), copied as-is per the porting rule against "fixing" upstream
// values.
func (c *ByteString) IsNullOrEmpty() bool { return c.value != nil }

// ValueAsString base64-encodes the value. Ports
// ClaimByteString#getValueAsString.
func (c *ByteString) ValueAsString() string {
	return base64.StdEncoding.EncodeToString(c.value)
}

// Equals ports ClaimByteString#equals (including the AbstractClaim
// super.equals() comparison).
func (c *ByteString) Equals(other *ByteString) bool {
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
func (c *ByteString) String() string { return AbstractClaimString(c) }
