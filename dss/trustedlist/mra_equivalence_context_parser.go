// Ported from specs-trusted-list/src/main/java/eu/europa/esig/trustedlist/mra/parsers/MRAEquivalenceContextParser.java (DSS 6.5.RC1).

package trustedlist

import "github.com/utain/esig/dss/enumerations"

// MRAEquivalenceContextParserParse parses the string and returns a
// MRAEquivalenceContext, the empty value if v does not match any known URI.
// slf4j's LOG.warn on an unresolved value is dropped per PORTING.md/
// S9_BRIEF.md's hard rules (slf4j dropped except job alerting semantics,
// which do not apply here).
func MRAEquivalenceContextParserParse(v string) enumerations.MRAEquivalenceContext {
	for _, m := range enumerations.MRAEquivalenceContextValues() {
		if m.URI() == v {
			return m
		}
	}
	return ""
}

// MRAEquivalenceContextParserPrint returns the URI of m, the empty string
// for the empty (Java null) value.
func MRAEquivalenceContextParserPrint(m enumerations.MRAEquivalenceContext) string {
	if m == "" {
		return ""
	}
	return m.URI()
}
