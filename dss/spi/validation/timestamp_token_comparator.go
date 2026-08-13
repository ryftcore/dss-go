// Ported from dss-spi/src/main/java/eu/europa/esig/dss/spi/x509/tsp/TimestampTokenComparator.java (DSS 6.5.RC1).
package validation

// TimestampTokenComparator compares TimestampTokens.
//
// java.io.Serializable is dropped (no Go counterpart).
type TimestampTokenComparator struct{}

// NewTimestampTokenComparator creates the comparator. Port of the default constructor.
func NewTimestampTokenComparator() TimestampTokenComparator {
	return TimestampTokenComparator{}
}

// Compare orders two time-stamp tokens, returning a negative number, zero or a positive number
// as tst1 sorts before, equal to, or after tst2. Port of compare(TimestampToken, TimestampToken).
func (c TimestampTokenComparator) Compare(tst1, tst2 *TimestampToken) int {
	result := c.compareByGenerationTime(tst1, tst2)
	if result == 0 {
		result = c.compareByTokenType(tst1, tst2)
	}
	if result == 0 {
		result = c.compareByManifest(tst1, tst2)
	}
	if result == 0 {
		result = c.compareByCoverage(tst1, tst2)
	}
	if result == 0 {
		result = c.compareByTimestampedReferences(tst1, tst2)
	}
	return result
}

// Less adapts Compare to the sort.Slice / slices.SortFunc convention.
func (c TimestampTokenComparator) Less(tst1, tst2 *TimestampToken) bool {
	return c.Compare(tst1, tst2) < 0
}

// compareByGenerationTime ports the private compareByGenerationTime.
func (c TimestampTokenComparator) compareByGenerationTime(tst1, tst2 *TimestampToken) int {
	return tst1.GenerationTime().Compare(tst2.GenerationTime())
}

// compareByTokenType ports the private compareByTokenType.
func (c TimestampTokenComparator) compareByTokenType(tst1, tst2 *TimestampToken) int {
	tst1Type := tst1.TimeStampType()
	tst2Type := tst2.TimeStampType()
	return tst1Type.Compare(tst2Type)
}

// compareByManifest ports the private compareByManifest.
func (c TimestampTokenComparator) compareByManifest(tst1, tst2 *TimestampToken) int {
	tst1ManifestFile := tst1.ManifestFile()
	tst2ManifestFile := tst2.ManifestFile()
	if tst1ManifestFile != nil && tst1ManifestFile.IsDocumentCovered(tst2.Filename()) {
		return 1
	} else if tst2ManifestFile != nil && tst2ManifestFile.IsDocumentCovered(tst1.Filename()) {
		return -1
	}
	return 0
}

// compareByCoverage ports the private compareByCoverage.
func (c TimestampTokenComparator) compareByCoverage(tst1, tst2 *TimestampToken) int {
	if c.isCoveredByTimestamp(tst1, tst2) {
		return -1
	} else if c.isCoveredByTimestamp(tst2, tst1) {
		return 1
	}
	return 0
}

// isCoveredByTimestamp ports the private isCoveredByTimestamp.
func (c TimestampTokenComparator) isCoveredByTimestamp(tst1, tst2 *TimestampToken) bool {
	tst2References := tst2.TimestampedReferences()
	for _, timestampedReference := range tst2References {
		if tst1.DSSIDAsString() == timestampedReference.ObjectId() {
			return true
		}
	}
	return false
}

// compareByTimestampedReferences ports the private compareByTimestampedReferences.
//
// Java's null checks on both reference lists survive as nil checks: getTimestampedReferences()
// returns the list the token was constructed with, which a caller may leave null.
func (c TimestampTokenComparator) compareByTimestampedReferences(tst1, tst2 *TimestampToken) int {
	tst1References := tst1.TimestampedReferences()
	tst2References := tst2.TimestampedReferences()
	if tst1References != nil && tst2References != nil {
		if len(tst1References) < len(tst2References) {
			return -1
		} else if len(tst1References) > len(tst2References) {
			return 1
		}
	}
	return 0
}
