// Ported from dss-model/.../claim/ClaimArray.java (DSS 6.5.RC1).
package claim

import (
	"reflect"
	"strings"
)

// Array represents an Array encoded (selectively) disclosable claim.
// It is designed for embedding by concrete subtypes: Java's abstract
// `createClaim(Object)` method is ported as the CreateClaim function
// field, which embedders MUST set (typically from their own constructor)
// before calling ListValue/ValueAsString.
type Array struct {
	AbstractClaim

	// value is the content of the array (raw, untyped items - mirrors
	// Java's List<?>).
	value []any

	// CreateClaim implements the abstract createClaim(Object) method:
	// creates a Claim for a single array item. Required non-nil to call
	// ListValue/ValueAsString.
	CreateClaim func(value any) Claim
}

// NewClaimArray ports the constructor with claim name, value, selectively
// disclosable status and parent claim provided.
func NewArray(name string, value []any, selectivelyDisclosable bool, parent Claim, createClaim func(value any) Claim) *Array {
	return NewArrayFull(name, "", value, selectivelyDisclosable, parent, createClaim)
}

// NewClaimArrayFull ports the constructor with claim name, namespace,
// value, selectively disclosable status and parent claim provided.
func NewArrayFull(name, namespace string, value []any, selectivelyDisclosable bool, parent Claim, createClaim func(value any) Claim) *Array {
	return &Array{
		AbstractClaim: NewAbstractClaimFull(name, namespace, selectivelyDisclosable, parent),
		value:         value,
		CreateClaim:   createClaim,
	}
}

// ListValue converts every array item into a Claim via CreateClaim. Ports
// ClaimArray#getListValue.
func (c *Array) ListValue() []Claim {
	if len(c.value) == 0 {
		return []Claim{}
	}
	result := make([]Claim, 0, len(c.value))
	for _, item := range c.value {
		result = append(result, c.CreateClaim(item))
	}
	return result
}

// IsArrayValueType always returns true.
func (c *Array) IsArrayValueType() bool { return true }

// IsNullOrEmpty ports ClaimArray#isNullOrEmpty.
func (c *Array) IsNullOrEmpty() bool { return len(c.value) == 0 }

// ValueAsString ports ClaimArray#getValueAsString.
func (c *Array) ValueAsString() string {
	var sb strings.Builder
	sb.WriteString("[")
	items := c.ListValue()
	for i, item := range items {
		if item.IsStringValueType() {
			sb.WriteString("\"")
		}
		sb.WriteString(item.ValueAsString())
		if item.IsStringValueType() {
			sb.WriteString("\"")
		}
		if i < len(items)-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteString("]")
	return sb.String()
}

// Equals ports ClaimArray#equals (including the AbstractClaim
// super.equals() comparison).
func (c *Array) Equals(other *Array) bool {
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
func (c *Array) String() string { return AbstractString(c) }
