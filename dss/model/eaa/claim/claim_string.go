// Ported from dss-model/.../claim/ClaimString.java (DSS 6.5.RC1).
package claim

// String represents a String encoded (selectively) disclosable
// claim.
//
// NOTE: since this codebase represents a Java null String as Go's ""
// (see claim.go), IsNullOrEmpty here is indistinguishable from "value is
// the empty string", matching the existing string-nullability convention
// used throughout dss/model.
type String struct {
	AbstractClaim

	// value is the string value of the claim.
	value string

	// valuePresent tracks whether value was constructed as non-nil, so
	// IsNullOrEmpty can still port Java's `value == null` check precisely.
	valuePresent bool
}

// NewClaimString ports the default constructor.
func NewString(value string) *String {
	return NewStringWithName("", value)
}

// NewClaimStringWithName ports the constructor with claim header name
// provided.
func NewStringWithName(name, value string) *String {
	return NewStringWithDisclosable(name, value, false)
}

// NewClaimStringWithDisclosable ports the constructor with claim name and
// selectively disclosable status provided.
func NewStringWithDisclosable(name, value string, selectivelyDisclosable bool) *String {
	return NewStringWithParent(name, value, selectivelyDisclosable, nil)
}

// NewClaimStringWithParent ports the constructor with claim name,
// selectively disclosable status and parent claim provided.
func NewStringWithParent(name, value string, selectivelyDisclosable bool, parent Claim) *String {
	return NewStringFull(name, "", value, selectivelyDisclosable, parent)
}

// NewClaimStringFull ports the constructor with claim name, namespace,
// selectively disclosable status and parent claim provided.
func NewStringFull(name, namespace, value string, selectivelyDisclosable bool, parent Claim) *String {
	return &String{
		AbstractClaim: NewAbstractClaimFull(name, namespace, selectivelyDisclosable, parent),
		value:         value,
		valuePresent:  true,
	}
}

// StringValue returns the string value of the claim.
func (c *String) StringValue() string { return c.value }

// IsStringValueType always returns true.
func (c *String) IsStringValueType() bool { return true }

// ValueAsString returns the string value of the claim.
func (c *String) ValueAsString() string { return c.value }

// IsNullOrEmpty ports ClaimString#isNullOrEmpty (`value == null`; not
// checking for the empty string, matching upstream verbatim).
func (c *String) IsNullOrEmpty() bool { return !c.valuePresent }

// Equals ports ClaimString#equals (including the AbstractClaim
// super.equals() comparison).
func (c *String) Equals(other *String) bool {
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
func (c *String) String() string { return AbstractString(c) }
