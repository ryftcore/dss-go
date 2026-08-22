// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/validation/analyzer/timestamp/TimestampAnalyzerComparator.java (DSS 6.5.RC1).
//
// java.io.Serializable is dropped (no Go counterpart), per PORTING.md.
package timestamp

import "github.com/ryftcore/dss-go/dss/spi/validation"

// timestampAnalyzerComparatorTimestampComparator is used to compare the timestamps. Port of the
// private static final TimestampTokenComparator timestampComparator field.
var timestampAnalyzerComparatorTimestampComparator = validation.NewTimestampTokenComparator()

// TimestampAnalyzerComparator compares TimestampAnalyzers.
type TimestampAnalyzerComparator struct{}

// NewTimestampAnalyzerComparator instantiates the comparator. Port of the default constructor.
func NewTimestampAnalyzerComparator() TimestampAnalyzerComparator {
	return TimestampAnalyzerComparator{}
}

// Compare orders two timestamp analyzers, returning a negative number, zero or a positive
// number as timestampAnalyzer1 sorts before, equal to, or after timestampAnalyzer2. Port of
// compare(TimestampAnalyzer, TimestampAnalyzer).
func (c TimestampAnalyzerComparator) Compare(timestampAnalyzer1, timestampAnalyzer2 TimestampAnalyzer) int {
	return timestampAnalyzerComparatorTimestampComparator.Compare(
		timestampAnalyzer1.Timestamp(), timestampAnalyzer2.Timestamp())
}

// Less adapts Compare to the sort.Slice / slices.SortFunc convention, matching the precedent
// set by validation.TimestampTokenComparator.Less.
func (c TimestampAnalyzerComparator) Less(timestampAnalyzer1, timestampAnalyzer2 TimestampAnalyzer) bool {
	return c.Compare(timestampAnalyzer1, timestampAnalyzer2) < 0
}
