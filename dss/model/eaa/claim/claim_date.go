// Ported from dss-model/.../claim/ClaimDate.java (DSS 6.5.RC1).
package claim

import "time"

// RFC3339TimeFormat is the RFC 3339 DateTime format used by default. Ports
// ClaimDate#RFC3339_TIME_FORMAT ("yyyy-MM-dd'T'HH:mm:ss'Z'"), expressed as
// a Go time layout.
const RFC3339TimeFormat = "2006-01-02T15:04:05Z"

// ClaimDate represents a Date encoded (selectively) disclosable claim.
type ClaimDate struct {
	AbstractClaim

	// value is the date value of the claim.
	value *time.Time
}

// NewClaimDate ports the default constructor.
func NewClaimDate(value *time.Time) *ClaimDate {
	return NewClaimDateWithName("", value)
}

// NewClaimDateWithName ports the constructor with claim name provided.
func NewClaimDateWithName(name string, value *time.Time) *ClaimDate {
	return NewClaimDateWithDisclosable(name, value, false)
}

// NewClaimDateWithDisclosable ports the constructor with claim name and
// selectively disclosable status provided.
func NewClaimDateWithDisclosable(name string, value *time.Time, selectivelyDisclosable bool) *ClaimDate {
	return NewClaimDateWithParent(name, value, selectivelyDisclosable, nil)
}

// NewClaimDateWithParent ports the constructor with claim name,
// selectively disclosable status and parent claim provided.
func NewClaimDateWithParent(name string, value *time.Time, selectivelyDisclosable bool, parent Claim) *ClaimDate {
	return NewClaimDateFull(name, "", value, selectivelyDisclosable, parent)
}

// NewClaimDateFull ports the constructor with claim name, namespace,
// selectively disclosable status and parent claim provided.
func NewClaimDateFull(name, namespace string, value *time.Time, selectivelyDisclosable bool, parent Claim) *ClaimDate {
	return &ClaimDate{
		AbstractClaim: NewAbstractClaimFull(name, namespace, selectivelyDisclosable, parent),
		value:         value,
	}
}

// DateValue returns the date value of the claim.
func (c *ClaimDate) DateValue() *time.Time { return c.value }

// IsDateValueType always returns true.
func (c *ClaimDate) IsDateValueType() bool { return true }

// ValueAsString formats the date in RFC3339TimeFormat, UTC. Ports
// ClaimDate#getValueAsString ("N/A" when the value is nil).
func (c *ClaimDate) ValueAsString() string {
	if c.value == nil {
		return "N/A"
	}
	return c.value.UTC().Format(RFC3339TimeFormat)
}

// IsNullOrEmpty ports ClaimDate#isNullOrEmpty.
func (c *ClaimDate) IsNullOrEmpty() bool { return c.value == nil }

// Equals ports ClaimDate#equals (including the AbstractClaim
// super.equals() comparison).
func (c *ClaimDate) Equals(other *ClaimDate) bool {
	if c == other {
		return true
	}
	if other == nil {
		return false
	}
	if !c.AbstractClaim.Equals(&other.AbstractClaim) {
		return false
	}
	if c.value == nil || other.value == nil {
		return c.value == other.value
	}
	return c.value.Equal(*other.value)
}

// String ports AbstractClaim#toString, inherited by this claim.
func (c *ClaimDate) String() string { return AbstractClaimString(c) }
