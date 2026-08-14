// Ported from dss-model/.../claim/ClaimIntegrity.java (DSS 6.5.RC1).
package claim

import "github.com/utain/esig/dss/enumerations"

// ClaimIntegrity represents a claim integrity definition, when
// applicable. This definition is based on W3C Subresource Integrity
// (https://www.w3.org/TR/2016/REC-SRI-20160623/).
type ClaimIntegrity interface {
	Claim

	// DigestAlgorithm gets the Digest Algorithm used to compute claim
	// integrity digest, when present. Ports
	// ClaimIntegrity#getDigestAlgorithm.
	DigestAlgorithm() enumerations.DigestAlgorithm

	// DigestValue gets the claim integrity digest value, when present.
	// NOTE: the digest computation depends on claim semantics. Ports
	// ClaimIntegrity#getDigestValue.
	DigestValue() []byte
}
