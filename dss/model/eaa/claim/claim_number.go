// Ported from dss-model/.../claim/ClaimNumber.java (DSS 6.5.RC1).
package claim

import (
	"fmt"
	"reflect"
)

// Number represents a Number encoded (selectively) disclosable claim.
//
// The value is kept as `any` (nil for Java null), holding whatever
// concrete Go numeric type it was constructed with, mirroring Java's
// polymorphic Number wrapper (Integer/Long/Double/BigDecimal/...).
type Number struct {
	AbstractClaim

	// value is the number value of the claim.
	value any
}

// NewClaimNumber ports the default constructor.
func NewClaimNumber(value any) *Number {
	return NewClaimNumberWithName("", value)
}

// NewClaimNumberWithName ports the constructor with claim name provided.
func NewClaimNumberWithName(name string, value any) *Number {
	return NewClaimNumberWithDisclosable(name, value, false)
}

// NewClaimNumberWithDisclosable ports the constructor with claim name and
// selectively disclosable status provided.
func NewClaimNumberWithDisclosable(name string, value any, selectivelyDisclosable bool) *Number {
	return NewClaimNumberWithParent(name, value, selectivelyDisclosable, nil)
}

// NewClaimNumberWithParent ports the constructor with claim name,
// selectively disclosable status and parent claim provided.
func NewClaimNumberWithParent(name string, value any, selectivelyDisclosable bool, parent Claim) *Number {
	return NewClaimNumberFull(name, "", value, selectivelyDisclosable, parent)
}

// NewClaimNumberFull ports the constructor with claim name, namespace,
// selectively disclosable status and parent claim provided.
func NewClaimNumberFull(name, namespace string, value any, selectivelyDisclosable bool, parent Claim) *Number {
	return &Number{
		AbstractClaim: NewAbstractClaimFull(name, namespace, selectivelyDisclosable, parent),
		value:         value,
	}
}

// NumberValue returns the number value of the claim.
func (c *Number) NumberValue() any { return c.value }

// IsNumberValueType always returns true.
func (c *Number) IsNumberValueType() bool { return true }

// ValueAsString ports ClaimNumber#getValueAsString.
func (c *Number) ValueAsString() string {
	if c.value == nil {
		return ""
	}
	return fmt.Sprint(c.value)
}

// IsNullOrEmpty ports ClaimNumber#isNullOrEmpty.
func (c *Number) IsNullOrEmpty() bool { return c.value == nil }

// Equals ports ClaimNumber#equals (including the AbstractClaim
// super.equals() comparison).
func (c *Number) Equals(other *Number) bool {
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
func (c *Number) String() string { return AbstractClaimString(c) }
