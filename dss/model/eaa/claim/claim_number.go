// Ported from dss-model/.../claim/ClaimNumber.java (DSS 6.5.RC1).
package claim

import (
	"fmt"
	"reflect"
)

// ClaimNumber represents a Number encoded (selectively) disclosable claim.
//
// The value is kept as `any` (nil for Java null), holding whatever
// concrete Go numeric type it was constructed with, mirroring Java's
// polymorphic Number wrapper (Integer/Long/Double/BigDecimal/...).
type ClaimNumber struct {
	AbstractClaim

	// value is the number value of the claim.
	value any
}

// NewClaimNumber ports the default constructor.
func NewClaimNumber(value any) *ClaimNumber {
	return NewClaimNumberWithName("", value)
}

// NewClaimNumberWithName ports the constructor with claim name provided.
func NewClaimNumberWithName(name string, value any) *ClaimNumber {
	return NewClaimNumberWithDisclosable(name, value, false)
}

// NewClaimNumberWithDisclosable ports the constructor with claim name and
// selectively disclosable status provided.
func NewClaimNumberWithDisclosable(name string, value any, selectivelyDisclosable bool) *ClaimNumber {
	return NewClaimNumberWithParent(name, value, selectivelyDisclosable, nil)
}

// NewClaimNumberWithParent ports the constructor with claim name,
// selectively disclosable status and parent claim provided.
func NewClaimNumberWithParent(name string, value any, selectivelyDisclosable bool, parent Claim) *ClaimNumber {
	return NewClaimNumberFull(name, "", value, selectivelyDisclosable, parent)
}

// NewClaimNumberFull ports the constructor with claim name, namespace,
// selectively disclosable status and parent claim provided.
func NewClaimNumberFull(name, namespace string, value any, selectivelyDisclosable bool, parent Claim) *ClaimNumber {
	return &ClaimNumber{
		AbstractClaim: NewAbstractClaimFull(name, namespace, selectivelyDisclosable, parent),
		value:         value,
	}
}

// NumberValue returns the number value of the claim.
func (c *ClaimNumber) NumberValue() any { return c.value }

// IsNumberValueType always returns true.
func (c *ClaimNumber) IsNumberValueType() bool { return true }

// ValueAsString ports ClaimNumber#getValueAsString.
func (c *ClaimNumber) ValueAsString() string {
	if c.value == nil {
		return ""
	}
	return fmt.Sprint(c.value)
}

// IsNullOrEmpty ports ClaimNumber#isNullOrEmpty.
func (c *ClaimNumber) IsNullOrEmpty() bool { return c.value == nil }

// Equals ports ClaimNumber#equals (including the AbstractClaim
// super.equals() comparison).
func (c *ClaimNumber) Equals(other *ClaimNumber) bool {
	if c == other {
		return true
	}
	if other == nil {
		return false
	}
	if !c.AbstractClaim.Equals(&other.AbstractClaim) {
		return false
	}
	return reflect.DeepEqual(c.value, other.value)
}

// String ports AbstractClaim#toString, inherited by this claim.
func (c *ClaimNumber) String() string { return AbstractClaimString(c) }
