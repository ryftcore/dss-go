// Ported from dss-enumerations/.../QWACProfile.java (DSS 6.5.RC1).
package enumerations

// QWACProfile contains QWAC profiles as per ETSI TS 119 411-5.
type QWACProfile string

const (
	// QWACProfileQWAC1 is a Qualified certificate for website
	// authentication based on Approach #1 in clause 1 of the ETSI TS 119
	// 411-5.
	QWACProfileQWAC1 QWACProfile = "QWAC_1"
	// QWACProfileQWAC2 is a Qualified certificate for website
	// authentication based on Approach #2 in clause 1 of the ETSI TS 119
	// 411-5.
	QWACProfileQWAC2 QWACProfile = "QWAC_2"
	// QWACProfileTLSByQWAC2 is a TLS certificate supported by a
	// qualified certificate for website authentication based on Approach
	// #2 in clause 1 of the ETSI TS 119 411-5, through TLS Certificate
	// Binding.
	QWACProfileTLSByQWAC2 QWACProfile = "TLS_BY_QWAC_2"
	// QWACProfileNotQWAC is Not a Qualified certificate for website
	// authentication based on clause 1 of the ETSI TS 119 411-5.
	QWACProfileNotQWAC QWACProfile = "NOT_QWAC"
)

// qwacProfileReadable holds the user-friendly identifier of the QWAC
// certificate type for each constant.
var qwacProfileReadable = map[QWACProfile]string{
	QWACProfileQWAC1:      "1-QWAC",
	QWACProfileQWAC2:      "2-QWAC",
	QWACProfileTLSByQWAC2: "TLS certificate supported by 2-QWAC",
	QWACProfileNotQWAC:    "Not QWAC",
}

// QWACProfileValues returns all constants in declaration order.
func QWACProfileValues() []QWACProfile {
	return []QWACProfile{
		QWACProfileQWAC1,
		QWACProfileQWAC2,
		QWACProfileTLSByQWAC2,
		QWACProfileNotQWAC,
	}
}

// Readable gets the user-friendly label.
func (q QWACProfile) Readable() string {
	return qwacProfileReadable[q]
}

// QWACProfileFromReadable gets the QWACProfile from the readable string.
// Returns "" (the zero value) if readable is empty or not found, mirroring
// Java's null return.
func QWACProfileFromReadable(readable string) QWACProfile {
	if readable == "" {
		return ""
	}
	for _, v := range QWACProfileValues() {
		if qwacProfileReadable[v] == readable {
			return v
		}
	}
	return ""
}
