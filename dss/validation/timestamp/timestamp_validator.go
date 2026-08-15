// Ported from dss-validation/src/main/java/eu/europa/esig/dss/validation/timestamp/TimestampValidator.java (DSS 6.5.RC1).
package timestamp

import (
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/spi/validation"
)

// TimestampValidator is the interface to be used for timestamp validation.
type TimestampValidator interface {
	// Timestamp returns a single TimestampToken to be validated. Port of getTimestamp().
	Timestamp() *validation.TimestampToken

	// TimestampedData returns the timestamped data. Port of getTimestampedData().
	TimestampedData() model.DSSDocument
}
