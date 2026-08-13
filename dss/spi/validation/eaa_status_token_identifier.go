// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/eaa/EAAStatusTokenIdentifier.java (DSS 6.5.RC1).
package validation

import (
	"github.com/utain/esig/dss/model"
)

// eaaStatusTokenIdentifierClassName is EAAStatusTokenIdentifier.class.getSimpleName().
const eaaStatusTokenIdentifierClassName = "EAAStatusTokenIdentifier"

// eaaStatusTokenIdentifierDefaultPrefix is the default identifier prefix. Port of the "ST-"
// literal used by the single-argument constructor.
const eaaStatusTokenIdentifierDefaultPrefix = "ST-"

// EAAStatusTokenIdentifier is an identifier for an EAA Status Token.
//
// Java declares the class final; Go has no such modifier, so the type is simply never embedded.
type EAAStatusTokenIdentifier struct {
	model.TokenIdentifier
}

// NewEAAStatusTokenIdentifier builds the identifier of the given EAA revocation token using
// the default "ST-" prefix. Port of the EAAStatusTokenIdentifier(EAARevocationToken)
// constructor.
func NewEAAStatusTokenIdentifier(eaaRevocationToken *EAARevocationToken) *EAAStatusTokenIdentifier {
	return newEAAStatusTokenIdentifier(eaaStatusTokenIdentifierDefaultPrefix, eaaRevocationToken)
}

// newEAAStatusTokenIdentifier builds the identifier with a custom prefix. Port of the
// package-private EAAStatusTokenIdentifier(String, EAARevocationToken) constructor.
func newEAAStatusTokenIdentifier(prefix string, eaaRevocationToken *EAARevocationToken) *EAAStatusTokenIdentifier {
	return &EAAStatusTokenIdentifier{
		model.NewTokenIdentifierFromToken(eaaStatusTokenIdentifierClassName, prefix, eaaRevocationToken),
	}
}

// compile-time assertion: an EAAStatusTokenIdentifier is an Identifier.
var _ model.Identifier = (*EAAStatusTokenIdentifier)(nil)
