// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/eaa/status/EAARevocationSource.java (DSS 6.5.RC1).
package status

import "github.com/utain/esig/dss/spi/validation"

// EAARevocationSource executes an EAA Status request for the given EAA
// token using the Status List Token mechanism, as defined in IETF Token
// Status List (TSL):
// https://www.ietf.org/archive/id/draft-ietf-oauth-status-list-20.html
type EAARevocationSource interface {
	// EAARevocation gets the resulting revocation token for the given
	// EAA. Ports #getEAARevocation.
	EAARevocation(eaa validation.EAA) validation.EAARevocationToken
}
