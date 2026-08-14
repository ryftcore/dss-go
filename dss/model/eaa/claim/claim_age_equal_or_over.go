// Ported from dss-model/.../claim/ClaimAgeEqualOrOver.java (DSS 6.5.RC1).
package claim

// ClaimAgeEqualOrOver contains a list of age_over_NN claims.
type ClaimAgeEqualOrOver interface {
	Claim

	// AgeOverNNClaims gets a list of age_over_NN claims. Ports
	// ClaimAgeEqualOrOver#getAgeOverNNClaims.
	AgeOverNNClaims() []ClaimAgeOverNN
}
