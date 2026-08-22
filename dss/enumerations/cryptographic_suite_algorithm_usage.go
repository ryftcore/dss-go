// Ported from dss-enumerations/.../CryptographicSuiteAlgorithmUsage.java (DSS 6.5.RC1).
package enumerations

// CryptographicSuiteAlgorithmUsage lists possible algorithm-usage element
// URIs as per ETSI TS 119 322. Implements UriBasedEnum.
type CryptographicSuiteAlgorithmUsage string

const (
	// CryptographicSuiteAlgorithmUsageSignData shall be used to indicate
	// that the evaluation is applicable for any signed data.
	CryptographicSuiteAlgorithmUsageSignData CryptographicSuiteAlgorithmUsage = "SIGN_DATA"
	// CryptographicSuiteAlgorithmUsageSignCertificates shall be used to
	// indicate that the evaluation is applicable for signing certificates.
	CryptographicSuiteAlgorithmUsageSignCertificates CryptographicSuiteAlgorithmUsage = "SIGN_CERTIFICATES"
	// CryptographicSuiteAlgorithmUsageSignOCSP shall be used to indicate
	// that the evaluation is applicable for signing OCSP responses.
	CryptographicSuiteAlgorithmUsageSignOCSP CryptographicSuiteAlgorithmUsage = "SIGN_OCSP"
	// CryptographicSuiteAlgorithmUsageSignTimestamps shall be used to
	// indicate that the evaluation is applicable for signing timestamps.
	CryptographicSuiteAlgorithmUsageSignTimestamps CryptographicSuiteAlgorithmUsage = "SIGN_TIMESTAMPS"
	// CryptographicSuiteAlgorithmUsageValidateData shall be used to
	// indicate that the evaluation is applicable for the validation of any
	// signed data.
	CryptographicSuiteAlgorithmUsageValidateData CryptographicSuiteAlgorithmUsage = "VALIDATE_DATA"
	// CryptographicSuiteAlgorithmUsageValidateCertificates shall be used
	// to indicate that the evaluation is applicable for the validation of
	// certificates.
	CryptographicSuiteAlgorithmUsageValidateCertificates CryptographicSuiteAlgorithmUsage = "VALIDATE_CERTIFICATES"
	// CryptographicSuiteAlgorithmUsageValidateOCSP shall be used to
	// indicate that the evaluation is applicable for the validation of OCSP
	// responses.
	CryptographicSuiteAlgorithmUsageValidateOCSP CryptographicSuiteAlgorithmUsage = "VALIDATE_OCSP"
	// CryptographicSuiteAlgorithmUsageValidateTimestamps shall be used to
	// indicate that the evaluation is applicable for the validation of
	// timestamps.
	CryptographicSuiteAlgorithmUsageValidateTimestamps CryptographicSuiteAlgorithmUsage = "VALIDATE_TIMESTAMPS"
)

// cryptographicSuiteAlgorithmUsageURIs holds the URI for each constant.
var cryptographicSuiteAlgorithmUsageURIs = map[CryptographicSuiteAlgorithmUsage]string{
	CryptographicSuiteAlgorithmUsageSignData:             "http://uri.etsi.org/19322/sign_data",
	CryptographicSuiteAlgorithmUsageSignCertificates:     "http://uri.etsi.org/19322/sign_data/sign_certificates",
	CryptographicSuiteAlgorithmUsageSignOCSP:             "http://uri.etsi.org/19322/sign_data/sign_ocsp",
	CryptographicSuiteAlgorithmUsageSignTimestamps:       "http://uri.etsi.org/19322/sign_data/sign_timestamps",
	CryptographicSuiteAlgorithmUsageValidateData:         "http://uri.etsi.org/19322/sign_data/validate_data",
	CryptographicSuiteAlgorithmUsageValidateCertificates: "http://uri.etsi.org/19322/sign_data/validate_data/validate_certificates",
	CryptographicSuiteAlgorithmUsageValidateOCSP:         "http://uri.etsi.org/19322/sign_data/validate_data/validate_ocsp",
	CryptographicSuiteAlgorithmUsageValidateTimestamps:   "http://uri.etsi.org/19322/sign_data/validate_data/validate_timestamps",
}

// CryptographicSuiteAlgorithmUsageValues returns all constants in
// declaration order.
func CryptographicSuiteAlgorithmUsageValues() []CryptographicSuiteAlgorithmUsage {
	return []CryptographicSuiteAlgorithmUsage{
		CryptographicSuiteAlgorithmUsageSignData,
		CryptographicSuiteAlgorithmUsageSignCertificates,
		CryptographicSuiteAlgorithmUsageSignOCSP,
		CryptographicSuiteAlgorithmUsageSignTimestamps,
		CryptographicSuiteAlgorithmUsageValidateData,
		CryptographicSuiteAlgorithmUsageValidateCertificates,
		CryptographicSuiteAlgorithmUsageValidateOCSP,
		CryptographicSuiteAlgorithmUsageValidateTimestamps,
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
