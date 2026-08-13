// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/validation/analyzer/timestamp/TimestampAnalyzer.java (DSS 6.5.RC1).
package timestamp

import (
	"github.com/utain/esig/dss/model"
	"github.com/utain/esig/dss/spi/validation"
)

// TimestampAnalyzer performs processing of a timestamp.
type TimestampAnalyzer interface {
	// Timestamp returns a single TimestampToken to be validated. Port of getTimestamp().
	Timestamp() *validation.TimestampToken

	// TimestampedData returns the timestamped data. Port of getTimestampedData().
	TimestampedData() model.DSSDocument
}
