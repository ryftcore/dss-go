// Ported from dss-model/src/main/java/eu/europa/esig/dss/model/x509/TokenComparator.java (DSS 6.5.RC1).
package model

import "strings"

// TokenComparator compares and sorts tokens by their DSS identifier.
type TokenComparator struct{}

// NewTokenComparator creates the comparator.
func NewTokenComparator() TokenComparator {
	return TokenComparator{}
}

// Compare orders two tokens by their DSS Id string, returning a negative number, zero or a
// positive number as o1 sorts before, equal to, or after o2. Port of compare(Token, Token).
//
// Java returns String#compareTo's UTF-16 code unit difference rather than -1/0/1; DSS Id
// strings are ASCII (a prefix plus hex), so the ordering is identical even though the exact
// magnitudes are not.
func (c TokenComparator) Compare(o1, o2 Token) int {
	return strings.Compare(o1.DSSIDAsString(), o2.DSSIDAsString())
}

// Less adapts Compare to the sort.Slice / slices.SortFunc convention.
func (c TokenComparator) Less(o1, o2 Token) bool {
	return c.Compare(o1, o2) < 0
}
