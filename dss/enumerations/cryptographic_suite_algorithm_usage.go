// Ported from dss-enumerations/.../CryptographicSuiteAlgorithmUsage.java (DSS 6.5.RC1).
package enumerations

// CryptographicSuiteAlgorithmUsage lists possible algorithm-usage element
// URIs as per ETSI TS 119 322. Implements UriBasedEnum.
type CryptographicSuiteAlgorithmUsage string

const (
	// CryptographicSuiteAlgorithmUsage_SIGN_DATA shall be used to indicate
	// that the evaluation is applicable for any signed data.
	CryptographicSuiteAlgorithmUsage_SIGN_DATA CryptographicSuiteAlgorithmUsage = "SIGN_DATA"
	// CryptographicSuiteAlgorithmUsage_SIGN_CERTIFICATES shall be used to
	// indicate that the evaluation is applicable for signing certificates.
	CryptographicSuiteAlgorithmUsage_SIGN_CERTIFICATES CryptographicSuiteAlgorithmUsage = "SIGN_CERTIFICATES"
	// CryptographicSuiteAlgorithmUsage_SIGN_OCSP shall be used to indicate
	// that the evaluation is applicable for signing OCSP responses.
	CryptographicSuiteAlgorithmUsage_SIGN_OCSP CryptographicSuiteAlgorithmUsage = "SIGN_OCSP"
	// CryptographicSuiteAlgorithmUsage_SIGN_TIMESTAMPS shall be used to
	// indicate that the evaluation is applicable for signing timestamps.
	CryptographicSuiteAlgorithmUsage_SIGN_TIMESTAMPS CryptographicSuiteAlgorithmUsage = "SIGN_TIMESTAMPS"
	// CryptographicSuiteAlgorithmUsage_VALIDATE_DATA shall be used to
	// indicate that the evaluation is applicable for the validation of any
	// signed data.
	CryptographicSuiteAlgorithmUsage_VALIDATE_DATA CryptographicSuiteAlgorithmUsage = "VALIDATE_DATA"
	// CryptographicSuiteAlgorithmUsage_VALIDATE_CERTIFICATES shall be used
	// to indicate that the evaluation is applicable for the validation of
	// certificates.
	CryptographicSuiteAlgorithmUsage_VALIDATE_CERTIFICATES CryptographicSuiteAlgorithmUsage = "VALIDATE_CERTIFICATES"
	// CryptographicSuiteAlgorithmUsage_VALIDATE_OCSP shall be used to
	// indicate that the evaluation is applicable for the validation of OCSP
	// responses.
	CryptographicSuiteAlgorithmUsage_VALIDATE_OCSP CryptographicSuiteAlgorithmUsage = "VALIDATE_OCSP"
	// CryptographicSuiteAlgorithmUsage_VALIDATE_TIMESTAMPS shall be used to
	// indicate that the evaluation is applicable for the validation of
	// timestamps.
	CryptographicSuiteAlgorithmUsage_VALIDATE_TIMESTAMPS CryptographicSuiteAlgorithmUsage = "VALIDATE_TIMESTAMPS"
)

// cryptographicSuiteAlgorithmUsageURIs holds the URI for each constant.
var cryptographicSuiteAlgorithmUsageURIs = map[CryptographicSuiteAlgorithmUsage]string{
	CryptographicSuiteAlgorithmUsage_SIGN_DATA:             "http://uri.etsi.org/19322/sign_data",
	CryptographicSuiteAlgorithmUsage_SIGN_CERTIFICATES:     "http://uri.etsi.org/19322/sign_data/sign_certificates",
	CryptographicSuiteAlgorithmUsage_SIGN_OCSP:             "http://uri.etsi.org/19322/sign_data/sign_ocsp",
	CryptographicSuiteAlgorithmUsage_SIGN_TIMESTAMPS:       "http://uri.etsi.org/19322/sign_data/sign_timestamps",
	CryptographicSuiteAlgorithmUsage_VALIDATE_DATA:         "http://uri.etsi.org/19322/sign_data/validate_data",
	CryptographicSuiteAlgorithmUsage_VALIDATE_CERTIFICATES: "http://uri.etsi.org/19322/sign_data/validate_data/validate_certificates",
	CryptographicSuiteAlgorithmUsage_VALIDATE_OCSP:         "http://uri.etsi.org/19322/sign_data/validate_data/validate_ocsp",
	CryptographicSuiteAlgorithmUsage_VALIDATE_TIMESTAMPS:   "http://uri.etsi.org/19322/sign_data/validate_data/validate_timestamps",
}

// CryptographicSuiteAlgorithmUsageValues returns all constants in
// declaration order.
func CryptographicSuiteAlgorithmUsageValues() []CryptographicSuiteAlgorithmUsage {
	return []CryptographicSuiteAlgorithmUsage{
		CryptographicSuiteAlgorithmUsage_SIGN_DATA,
		CryptographicSuiteAlgorithmUsage_SIGN_CERTIFICATES,
		CryptographicSuiteAlgorithmUsage_SIGN_OCSP,
		CryptographicSuiteAlgorithmUsage_SIGN_TIMESTAMPS,
		CryptographicSuiteAlgorithmUsage_VALIDATE_DATA,
		CryptographicSuiteAlgorithmUsage_VALIDATE_CERTIFICATES,
		CryptographicSuiteAlgorithmUsage_VALIDATE_OCSP,
		CryptographicSuiteAlgorithmUsage_VALIDATE_TIMESTAMPS,
	}
}

// URI returns the URI identifier of the algorithm-usage element. Implements
// UriBasedEnum.
func (c CryptographicSuiteAlgorithmUsage) URI() string {
	return cryptographicSuiteAlgorithmUsageURIs[c]
}

// CryptographicSuiteAlgorithmUsageFromURI returns a
// CryptographicSuiteAlgorithmUsage by the given URI, or "" (zero value) if
// uri is empty or unknown, mirroring Java's null return.
func CryptographicSuiteAlgorithmUsageFromURI(uri string) CryptographicSuiteAlgorithmUsage {
	if uri == "" {
		return ""
	}
	for _, v := range CryptographicSuiteAlgorithmUsageValues() {
		if cryptographicSuiteAlgorithmUsageURIs[v] == uri {
			return v
		}
	}
	return ""
}

// compile-time interface assertion.
var _ UriBasedEnum = CryptographicSuiteAlgorithmUsage("")
