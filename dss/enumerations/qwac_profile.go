// Ported from dss-enumerations/.../QWACProfile.java (DSS 6.5.RC1).
package enumerations

// QWACProfile contains QWAC profiles as per ETSI TS 119 411-5.
type QWACProfile string

const (
	// QWACProfile_QWAC_1 is a Qualified certificate for website
	// authentication based on Approach #1 in clause 1 of the ETSI TS 119
	// 411-5.
	QWACProfile_QWAC_1 QWACProfile = "QWAC_1"
	// QWACProfile_QWAC_2 is a Qualified certificate for website
	// authentication based on Approach #2 in clause 1 of the ETSI TS 119
	// 411-5.
	QWACProfile_QWAC_2 QWACProfile = "QWAC_2"
	// QWACProfile_TLS_BY_QWAC_2 is a TLS certificate supported by a
	// qualified certificate for website authentication based on Approach
	// #2 in clause 1 of the ETSI TS 119 411-5, through TLS Certificate
	// Binding.
	QWACProfile_TLS_BY_QWAC_2 QWACProfile = "TLS_BY_QWAC_2"
	// QWACProfile_NOT_QWAC is Not a Qualified certificate for website
	// authentication based on clause 1 of the ETSI TS 119 411-5.
	QWACProfile_NOT_QWAC QWACProfile = "NOT_QWAC"
)

// qwacProfileReadable holds the user-friendly identifier of the QWAC
// certificate type for each constant.
var qwacProfileReadable = map[QWACProfile]string{
	QWACProfile_QWAC_1:        "1-QWAC",
	QWACProfile_QWAC_2:        "2-QWAC",
	QWACProfile_TLS_BY_QWAC_2: "TLS certificate supported by 2-QWAC",
	QWACProfile_NOT_QWAC:      "Not QWAC",
}

// QWACProfileValues returns all constants in declaration order.
func QWACProfileValues() []QWACProfile {
	return []QWACProfile{
		QWACProfile_QWAC_1,
		QWACProfile_QWAC_2,
		QWACProfile_TLS_BY_QWAC_2,
		QWACProfile_NOT_QWAC,
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
