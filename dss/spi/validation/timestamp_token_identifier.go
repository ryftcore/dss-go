// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/x509/tsp/TimestampTokenIdentifier.java (DSS 6.5.RC1).
package validation

import (
	"github.com/ryftcore/dss-go/dss/model"
)

// timestampTokenIdentifierPrefix is the default identifier prefix for a time-stamp token.
// Port of the private TIMESTAMP_TOKEN_PREFIX constant.
const timestampTokenIdentifierPrefix = "T-"

// timestampTokenIdentifierClassName is TimestampTokenIdentifier.class.getSimpleName(), which
// model.NewTokenIdentifier records in the identifier and which the validation reports show.
const timestampTokenIdentifierClassName = "TimestampTokenIdentifier"

// TimestampTokenIdentifier is the identifier of a timestamp token.
//
// Java declares the class final; Go has no such modifier, so the type is simply never embedded.
type TimestampTokenIdentifier struct {
	model.TokenIdentifier
}

// NewTimestampTokenIdentifierFromToken builds the identifier of the given time-stamp token:
// the SHA-256 digest of its DER encoding, rendered with the "T-" prefix.
// Port of the TimestampTokenIdentifier(TimestampToken) constructor.
func NewTimestampTokenIdentifierFromToken(timestampToken model.Token) *TimestampTokenIdentifier {
	return &TimestampTokenIdentifier{
		model.NewTokenIdentifierFromToken(timestampTokenIdentifierClassName, timestampTokenIdentifierPrefix, timestampToken),
	}
}

// NewTimestampTokenIdentifier builds a time-stamp identifier from custom binaries.
// Port of the TimestampTokenIdentifier(byte[]) constructor.
func NewTimestampTokenIdentifier(binaries []byte) *TimestampTokenIdentifier {
	return &TimestampTokenIdentifier{
		model.NewTokenIdentifier(timestampTokenIdentifierClassName, timestampTokenIdentifierPrefix, binaries),
	}
}

// compile-time assertion: a TimestampTokenIdentifier is an Identifier.
var _ model.Identifier = (*TimestampTokenIdentifier)(nil)
