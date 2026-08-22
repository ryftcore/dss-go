// Ported from dss-model/.../claim/ClaimNull.java (DSS 6.5.RC1).
package claim

// Null represents a Null encoded (selectively) disclosable claim.
type Null struct {
	AbstractClaim
}

// NewClaimNull ports the default constructor.
func NewClaimNull() *Null {
	return &Null{AbstractClaim: NewAbstractClaim()}
}

// NewClaimNullWithName ports the constructor with claim header name
// provided.
func NewClaimNullWithName(name string) *Null {
	return &Null{AbstractClaim: NewAbstractClaimWithName(name)}
}

// NewClaimNullWithDisclosable ports the constructor with claim name and
// selectively disclosable status provided.
func NewClaimNullWithDisclosable(name string, selectivelyDisclosable bool) *Null {
	return &Null{AbstractClaim: NewAbstractClaimWithDisclosable(name, selectivelyDisclosable)}
}

// NewClaimNullWithParent ports the constructor with claim name,
// selectively disclosable status and parent claim provided.
func NewClaimNullWithParent(name string, selectivelyDisclosable bool, parent Claim) *Null {
	return &Null{AbstractClaim: NewAbstractClaimWithParent(name, selectivelyDisclosable, parent)}
}

// NewClaimNullFull ports the constructor with claim name, namespace,
// selectively disclosable status and parent claim provided.
func NewClaimNullFull(name, namespace string, selectivelyDisclosable bool, parent Claim) *Null {
	return &Null{AbstractClaim: NewAbstractClaimFull(name, namespace, selectivelyDisclosable, parent)}
}

// ValueAsString always returns "null".
func (c *Null) ValueAsString() string { return "null" }

// IsNullValueType always returns true.
func (c *Null) IsNullValueType() bool { return true }

// IsNullOrEmpty always returns true.
func (c *Null) IsNullOrEmpty() bool { return true }

// String ports AbstractClaim#toString, inherited by this claim.
func (c *Null) String() string { return AbstractClaimString(c) }
