// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/eaa/EAAKeyBindingPayload.java (DSS 6.5.RC1).
package validation

import "github.com/utain/esig/dss/model/eaa/claim"

// EAAKeyBindingPayload represents a key binding payload.
type EAAKeyBindingPayload interface {
	claim.Claim

	// Nonce gets the value of the nonce from the key binding payload. Port of getNonce().
	Nonce() claim.ClaimString

	// IssuedAt gets the issuance date from the key binding payload. Port of getIssuedAt().
	IssuedAt() claim.ClaimDate

	// Audience gets the value of the audience from the key binding payload. Port of getAudience().
	Audience() claim.ClaimString

	// SdHash gets the value of the "sd_hash" from the key binding payload. Port of getSdHash().
	SdHash() claim.ClaimString
}
