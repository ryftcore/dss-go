// Ported from dss-model/.../claim/ClaimString.java (DSS 6.5.RC1).
package claim

// ClaimString represents a String encoded (selectively) disclosable
// claim.
//
// NOTE: since this codebase represents a Java null String as Go's ""
// (see claim.go), IsNullOrEmpty here is indistinguishable from "value is
// the empty string", matching the existing string-nullability convention
// used throughout dss/model.
type ClaimString struct {
	AbstractClaim

	// value is the string value of the claim.
	value string

	// valuePresent tracks whether value was constructed as non-nil, so
	// IsNullOrEmpty can still port Java's `value == null` check precisely.
	valuePresent bool
}

// NewClaimString ports the default constructor.
func NewClaimString(value string) *ClaimString {
	return NewClaimStringWithName("", value)
}

// NewClaimStringWithName ports the constructor with claim header name
// provided.
func NewClaimStringWithName(name, value string) *ClaimString {
	return NewClaimStringWithDisclosable(name, value, false)
}

// NewClaimStringWithDisclosable ports the constructor with claim name and
// selectively disclosable status provided.
func NewClaimStringWithDisclosable(name, value string, selectivelyDisclosable bool) *ClaimString {
	return NewClaimStringWithParent(name, value, selectivelyDisclosable, nil)
}

// NewClaimStringWithParent ports the constructor with claim name,
// selectively disclosable status and parent claim provided.
func NewClaimStringWithParent(name, value string, selectivelyDisclosable bool, parent Claim) *ClaimString {
	return NewClaimStringFull(name, "", value, selectivelyDisclosable, parent)
}

// NewClaimStringFull ports the constructor with claim name, namespace,
// selectively disclosable status and parent claim provided.
func NewClaimStringFull(name, namespace, value string, selectivelyDisclosable bool, parent Claim) *ClaimString {
	return &ClaimString{
		AbstractClaim: NewAbstractClaimFull(name, namespace, selectivelyDisclosable, parent),
		value:         value,
		valuePresent:  true,
	}
}

// StringValue returns the string value of the claim.
func (c *ClaimString) StringValue() string { return c.value }

// IsStringValueType always returns true.
func (c *ClaimString) IsStringValueType() bool { return true }

// ValueAsString returns the string value of the claim.
func (c *ClaimString) ValueAsString() string { return c.value }

// IsNullOrEmpty ports ClaimString#isNullOrEmpty (`value == null`; not
// checking for the empty string, matching upstream verbatim).
func (c *ClaimString) IsNullOrEmpty() bool { return !c.valuePresent }

// Equals ports ClaimString#equals (including the AbstractClaim
// super.equals() comparison).
func (c *ClaimString) Equals(other *ClaimString) bool {
	if c == other {
		return true
	}
	if other == nil {
		return false
	}
	if !c.AbstractClaim.Equals(&other.AbstractClaim) {
		return false
	}
	return c.valuePresent == other.valuePresent && c.value == other.value
}

// String ports AbstractClaim#toString, inherited by this claim.
func (c *ClaimString) String() string { return AbstractClaimString(c) }
