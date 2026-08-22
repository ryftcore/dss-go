// Ported from dss-enumerations/.../SigDMechanism.java (DSS 6.5.RC1).
//
// NOTE: the Java class logs via SLF4J (org.slf4j.Logger LOG field); that
// field is unused in the Java source itself (dead code) and is not ported.
package enumerations

// SigDMechanism defines a list of algorithm described in ETSI TS 119 182-1
// and TS 119 152-1 for incorporation of 'sigD' dictionary (see 5.2.8 The
// sigD header parameter and 5.2.9 The sigD header parameter, respectively).
type SigDMechanism string

const (
	// SigDMechanismHTTPHeaders is 5.2.8.2 Mechanism HttpHeaders.
	SigDMechanismHTTPHeaders SigDMechanism = "HTTP_HEADERS"
	// SigDMechanismObjectIDByURI is 5.2.8.3.2 Mechanism ObjectIdByURI.
	SigDMechanismObjectIDByURI SigDMechanism = "OBJECT_ID_BY_URI"
	// SigDMechanismObjectIDByURIHash is 5.2.8.3.3 Mechanism
	// ObjectIdByURIHash. NOTE: the default signature creation mechanism
	// used by DSS.
	SigDMechanismObjectIDByURIHash SigDMechanism = "OBJECT_ID_BY_URI_HASH"
	// SigDMechanismNoSigD creates a simple DETACHED signature with
	// omitted payload (without SigD element).
	SigDMechanismNoSigD SigDMechanism = "NO_SIG_D"
)

type sigDMechanismFields struct {
	jadesUri  string
	cbadesUri string
}

// sigDMechanismData holds the (jadesUri, cbadesUri) tuple for each
// constant. NOTE: HTTP_HEADERS has no CB-AdES URI in the Java source
// (constructor passed null); it is represented here as "" for both
// getCBAdESUri()'s null return and forCBAdESUri()'s null-check skip
// (mirrored by the Java forCBAdESUri null guard, since "" never equals a
// caller-supplied non-empty uri argument that forCBAdESUri would compare
// against a null field in Java it would panic on NPE via .equals; here the
// non-nil check is simply an empty-string comparison, which for NO_SIG_D's
// legitimate cbadesUri="" is intentionally still matched — see
// SigDMechanismForCBAdESUri).
var sigDMechanismData = map[SigDMechanism]sigDMechanismFields{
	SigDMechanismHTTPHeaders:       {"http://uri.etsi.org/19182/HttpHeaders", ""},
	SigDMechanismObjectIDByURI:     {"http://uri.etsi.org/19182/ObjectIdByURI", "http://uri.etsi.org/19152/ObjectIdByURI"},
	SigDMechanismObjectIDByURIHash: {"http://uri.etsi.org/19182/ObjectIdByURIHash", "http://uri.etsi.org/19152/ObjectIdByURIHash"},
	SigDMechanismNoSigD:            {"", ""},
}

// sigDMechanismHasCBAdESURI records, per constant, whether Java's cbadesUri
// constructor argument was non-null (HTTP_HEADERS passed null explicitly).
var sigDMechanismHasCBAdESURI = map[SigDMechanism]bool{
	SigDMechanismHTTPHeaders:       false,
	SigDMechanismObjectIDByURI:     true,
	SigDMechanismObjectIDByURIHash: true,
	SigDMechanismNoSigD:            true,
}

// SigDMechanismValues returns all constants in declaration order.
func SigDMechanismValues() []SigDMechanism {
	return []SigDMechanism{
		SigDMechanismHTTPHeaders,
		SigDMechanismObjectIDByURI,
		SigDMechanismObjectIDByURIHash,
		SigDMechanismNoSigD,
	}
}

// JAdESUri returns the JAdES ETSI TS 119 182-1 URI.
func (s SigDMechanism) JAdESUri() string {
	return sigDMechanismData[s].jadesUri
}

// CBAdESUri returns the CB-AdES ETSI TS 119 152-1 URI, and whether one is
// set (HTTP_HEADERS has none, mirroring Java's null field).
func (s SigDMechanism) CBAdESUri() (string, bool) {
	return sigDMechanismData[s].cbadesUri, sigDMechanismHasCBAdESURI[s]
}

// SigDMechanismForJAdESUri returns a SigDMechanism for the given JAdES ETSI
// TS 119 182-1 URI. Returns "" (zero value) if not found, mirroring Java's
// null return.
func SigDMechanismForJAdESUri(uri string) SigDMechanism {
	for _, v := range SigDMechanismValues() {
		if sigDMechanismData[v].jadesUri == uri {
			return v
		}
	}
	return ""
}

// SigDMechanismForCBAdESUri returns a SigDMechanism for the given CB-AdES
// ETSI TS 119 152-1 URI. Returns "" (zero value) if not found, mirroring
// Java's null return. HTTP_HEADERS (whose CB-AdES URI is unset) is skipped,
// matching Java's `sigDMechanism.getCBAdESUri() != null` guard.
func SigDMechanismForCBAdESUri(uri string) SigDMechanism {
	for _, v := range SigDMechanismValues() {
		if !sigDMechanismHasCBAdESURI[v] {
			continue
		}
		if sigDMechanismData[v].cbadesUri == uri {
			return v
		}
	}
	return ""
}
