// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/validation/scope/TimestampScopeFinder.java (DSS 6.5.RC1).
package scope

import (
	"github.com/utain/esig/dss/model/scope"
	"github.com/utain/esig/dss/spi/validation"
)

// TimestampScopeFinder is used to find a signature scope for a timestamp.
type TimestampScopeFinder interface {
	// FindTimestampScope returns a timestamp scope for the given TimestampToken. Port of
	// findTimestampScope(TimestampToken).
	FindTimestampScope(timestampToken *validation.TimestampToken) []scope.SignatureScope
}
