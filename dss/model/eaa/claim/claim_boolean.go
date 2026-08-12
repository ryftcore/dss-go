// Ported from dss-model/.../claim/ClaimBoolean.java (DSS 6.5.RC1).
package claim

import "strconv"

// ClaimBoolean represents a Boolean encoded (selectively) disclosable
// claim.
type ClaimBoolean struct {
	AbstractClaim

	// value is the boolean value of the claim (nil-able, mirrors Java's
	// Boolean).
	value *bool
}

// NewClaimBoolean ports the default constructor.
func NewClaimBoolean(value *bool) *ClaimBoolean {
	return NewClaimBooleanWithName("", value)
}

// NewClaimBooleanWithName ports the constructor with claim name provided.
func NewClaimBooleanWithName(name string, value *bool) *ClaimBoolean {
	return NewClaimBooleanWithDisclosable(name, value, false)
}

// NewClaimBooleanWithDisclosable ports the constructor with claim name and
// selectively disclosable status provided.
func NewClaimBooleanWithDisclosable(name string, value *bool, selectivelyDisclosable bool) *ClaimBoolean {
	return NewClaimBooleanWithParent(name, value, selectivelyDisclosable, nil)
}

// NewClaimBooleanWithParent ports the constructor with claim name,
// selectively disclosable status and parent claim provided.
func NewClaimBooleanWithParent(name string, value *bool, selectivelyDisclosable bool, parent Claim) *ClaimBoolean {
	return NewClaimBooleanFull(name, "", value, selectivelyDisclosable, parent)
}

// NewClaimBooleanFull ports the constructor with claim name, namespace,
// selectively disclosable status and parent claim provided.
func NewClaimBooleanFull(name, namespace string, value *bool, selectivelyDisclosable bool, parent Claim) *ClaimBoolean {
	return &ClaimBoolean{
		AbstractClaim: NewAbstractClaimFull(name, namespace, selectivelyDisclosable, parent),
		value:         value,
	}
}

// BooleanValue returns the boolean value of the claim.
func (c *ClaimBoolean) BooleanValue() *bool { return c.value }

// IsBooleanValueType always returns true.
func (c *ClaimBoolean) IsBooleanValueType() bool { return true }

// ValueAsString ports ClaimBoolean#getValueAsString.
func (c *ClaimBoolean) ValueAsString() string {
	if c.value == nil {
		return ""
	}
	return strconv.FormatBool(*c.value)
}

// IsNullOrEmpty ports ClaimBoolean#isNullOrEmpty.
func (c *ClaimBoolean) IsNullOrEmpty() bool { return c.value == nil }

// Equals ports ClaimBoolean#equals. NOTE upstream's override does not call
// super.equals() (unlike its siblings ClaimByteString/ClaimDate/
// ClaimNumber/ClaimString), so only the value field is compared here,
// matching Java verbatim.
func (c *ClaimBoolean) Equals(other *ClaimBoolean) bool {
	if c == other {
		return true
	}
	if other == nil {
		return false
	}
	if c.value == nil || other.value == nil {
		return c.value == other.value
	}
	return *c.value == *other.value
}

// String ports AbstractClaim#toString, inherited by this claim.
func (c *ClaimBoolean) String() string { return AbstractClaimString(c) }
