// Ported from specs-trusted-list/src/main/java/eu/europa/esig/trustedlist/mra/parsers/MRAStatusParser.java (DSS 6.5.RC1).

package trustedlist

import "github.com/ryftcore/dss-go/dss/enumerations"

// MRAStatusParserParse parses the string and returns a MRAStatus, the empty
// value if v does not match any known URI. slf4j's LOG.warn on an
// unresolved value is dropped per PORTING.md/S9_BRIEF.md's hard rules.
func MRAStatusParserParse(v string) enumerations.MRAStatus {
	for _, m := range enumerations.MRAStatusValues() {
		if m.URI() == v {
			return m
		}
	}
	return ""
}

// MRAStatusParserPrint returns the URI of m, the empty string for the
// empty (Java null) value.
func MRAStatusParserPrint(m enumerations.MRAStatus) string {
	if m == "" {
		return ""
	}
	return m.URI()
}
