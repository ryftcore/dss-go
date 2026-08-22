// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/x509/revocation/RevocationTokenIdentifier.java (DSS 6.5.RC1).
package spi

import (
	"github.com/ryftcore/dss-go/dss/model"
)

// RevocationTokenIdentifier is a unique id for a revocation data token (CRL or OCSP).
//
// Java parametrizes the constructor by RevocationToken<?> (a wildcard, i.e. the class itself
// is not generic over R); the Go port only needs the model.Token contract (Encoded()) that
// model.NewTokenIdentifierFromToken digests, so it is not generic over R either.
type RevocationTokenIdentifier struct {
	model.TokenIdentifier
}

// NewRevocationTokenIdentifier builds the identifier of the given revocation token: the
// SHA-256 digest of its DER encoding, rendered with the "R-" prefix.
// Port of the public RevocationTokenIdentifier(RevocationToken<?>) constructor.
func NewRevocationTokenIdentifier(revocationToken model.Token) *RevocationTokenIdentifier {
	return revocationTokenIdentifierNew("R-", revocationToken)
}

// revocationTokenIdentifierNew builds the identifier over the given revocation token with a
// custom prefix. Port of the package-private RevocationTokenIdentifier(String, RevocationToken<?>)
// constructor; the Go port keeps it unexported since it has no caller outside this file.
func revocationTokenIdentifierNew(prefix string, revocationToken model.Token) *RevocationTokenIdentifier {
	return &RevocationTokenIdentifier{
		model.NewTokenIdentifierFromToken("RevocationTokenIdentifier", prefix, revocationToken),
	}
}

// compile-time assertion: a RevocationTokenIdentifier is an Identifier.
var _ model.Identifier = (*RevocationTokenIdentifier)(nil)
