// Ported from dss-model/.../claim/ClaimNull.java (DSS 6.5.RC1).
package claim

// ClaimNull represents a Null encoded (selectively) disclosable claim.
type ClaimNull struct {
	AbstractClaim
}

// NewClaimNull ports the default constructor.
func NewClaimNull() *ClaimNull {
	return &ClaimNull{AbstractClaim: NewAbstractClaim()}
}

// NewClaimNullWithName ports the constructor with claim header name
// provided.
func NewClaimNullWithName(name string) *ClaimNull {
	return &ClaimNull{AbstractClaim: NewAbstractClaimWithName(name)}
}

// NewClaimNullWithDisclosable ports the constructor with claim name and
// selectively disclosable status provided.
func NewClaimNullWithDisclosable(name string, selectivelyDisclosable bool) *ClaimNull {
	return &ClaimNull{AbstractClaim: NewAbstractClaimWithDisclosable(name, selectivelyDisclosable)}
}

// NewClaimNullWithParent ports the constructor with claim name,
// selectively disclosable status and parent claim provided.
func NewClaimNullWithParent(name string, selectivelyDisclosable bool, parent Claim) *ClaimNull {
	return &ClaimNull{AbstractClaim: NewAbstractClaimWithParent(name, selectivelyDisclosable, parent)}
}

// NewClaimNullFull ports the constructor with claim name, namespace,
// selectively disclosable status and parent claim provided.
func NewClaimNullFull(name, namespace string, selectivelyDisclosable bool, parent Claim) *ClaimNull {
	return &ClaimNull{AbstractClaim: NewAbstractClaimFull(name, namespace, selectivelyDisclosable, parent)}
}

// ValueAsString always returns "null".
func (c *ClaimNull) ValueAsString() string { return "null" }

// IsNullValueType always returns true.
func (c *ClaimNull) IsNullValueType() bool { return true }

// IsNullOrEmpty always returns true.
func (c *ClaimNull) IsNullOrEmpty() bool { return true }

// String ports AbstractClaim#toString, inherited by this claim.
func (c *ClaimNull) String() string { return AbstractClaimString(c) }
