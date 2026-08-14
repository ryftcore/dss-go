// Ported from dss-enumerations/.../CertificatePolicy.java (DSS 6.5.RC1).
package enumerations

import "fmt"

// CertificatePolicy is the certificate policies enumeration, as in ETSI TS
// 101 456. Implements OidDescription.
type CertificatePolicy string

const (
	// CertificatePolicy_QCP_PUBLIC is a certificate policy for qualified
	// certificates issued to the public.
	CertificatePolicy_QCP_PUBLIC CertificatePolicy = "QCP_PUBLIC"
	// CertificatePolicy_QCP_PUBLIC_WITH_SSCD is a certificate policy for
	// qualified certificates issued to the public, requiring use of secure
	// signature-creation devices.
	CertificatePolicy_QCP_PUBLIC_WITH_SSCD CertificatePolicy = "QCP_PUBLIC_WITH_SSCD"
	// CertificatePolicy_NCP: Normalized Certificate Policy.
	CertificatePolicy_NCP CertificatePolicy = "NCP"
	// CertificatePolicy_NCPP: Normalized Certificate Policy requiring a
	// secure cryptographic device.
	CertificatePolicy_NCPP CertificatePolicy = "NCPP"
	// CertificatePolicy_LCP: Lightweight Certificate Policy.
	CertificatePolicy_LCP CertificatePolicy = "LCP"
	// CertificatePolicy_EVCP: Extended Validation Certificate Policy.
	CertificatePolicy_EVCP CertificatePolicy = "EVCP"
	// CertificatePolicy_DVCP: Domain Validation Certificate Policy.
	CertificatePolicy_DVCP CertificatePolicy = "DVCP"
	// CertificatePolicy_OVCP: Organizational Validation Certificate Policy.
	CertificatePolicy_OVCP CertificatePolicy = "OVCP"
	// CertificatePolicy_IVCP: Individual Validation Certificate Policy.
	CertificatePolicy_IVCP CertificatePolicy = "IVCP"
	// CertificatePolicy_QCP_NATURAL: certificate policy for EU qualified
	// certificates issued to natural persons.
	CertificatePolicy_QCP_NATURAL CertificatePolicy = "QCP_NATURAL"
	// CertificatePolicy_QCP_LEGAL: certificate policy for EU qualified
	// certificates issued to legal persons.
	CertificatePolicy_QCP_LEGAL CertificatePolicy = "QCP_LEGAL"
	// CertificatePolicy_QCP_NATURAL_QSCD: certificate policy for EU
	// qualified certificates issued to natural persons with private key
	// related to the certified public key in a QSCD.
	CertificatePolicy_QCP_NATURAL_QSCD CertificatePolicy = "QCP_NATURAL_QSCD"
	// CertificatePolicy_QCP_LEGAL_QSCD: certificate policy for EU qualified
	// certificates issued to legal persons with private key related to the
	// certified public key in a QSCD.
	CertificatePolicy_QCP_LEGAL_QSCD CertificatePolicy = "QCP_LEGAL_QSCD"
	// CertificatePolicy_QCP_WEB: certificate policy for EU qualified
	// website authentication certificates based on EVCP.
	CertificatePolicy_QCP_WEB CertificatePolicy = "QCP_WEB"
	// CertificatePolicy_QNCP_WEB: certificate policy for EU qualified
	// website authentication certificates based on NCP, and OVCP or IVCP.
	CertificatePolicy_QNCP_WEB CertificatePolicy = "QNCP_WEB"
	// CertificatePolicy_QNCP_WEB_GEN: certificate policy for EU qualified
	// website authentication based on NCP and requirements tagged as [WEB]
	// in ETSI EN 319 411-1.
	CertificatePolicy_QNCP_WEB_GEN CertificatePolicy = "QNCP_WEB_GEN"
)

// certificatePolicyFields holds the (description, oid) pair for each constant.
type certificatePolicyFields struct {
	description string
	oid         string
}

// certificatePolicyData holds the full field tuple for each constant, copied
// verbatim from the Java enum constructors.
var certificatePolicyData = map[CertificatePolicy]certificatePolicyFields{
	CertificatePolicy_QCP_PUBLIC:           {"qcp-public", "0.4.0.1456.1.2"},
	CertificatePolicy_QCP_PUBLIC_WITH_SSCD: {"qcp-public-with-sscd", "0.4.0.1456.1.1"},
	CertificatePolicy_NCP:                  {"normalized-certificate-policy", "0.4.0.2042.1.1"},
	CertificatePolicy_NCPP:                 {"normalized-certificate-policy-sscd", "0.4.0.2042.1.2"},
	CertificatePolicy_LCP:                  {"lightweight-certificate-policy", "0.4.0.2042.1.3"},
	CertificatePolicy_EVCP:                 {"extended-validation-certificate-policy", "0.4.0.2042.1.4"},
	CertificatePolicy_DVCP:                 {"domain-validation-certificate-policy", "0.4.0.2042.1.6"},
	CertificatePolicy_OVCP:                 {"organizational-validation-certificate-policy", "0.4.0.2042.1.7"},
	CertificatePolicy_IVCP:                 {"individual-validation-certificate-policy", "0.4.0.2042.1.8"},
	CertificatePolicy_QCP_NATURAL:          {"qcp-natural", "0.4.0.194112.1.0"},
	CertificatePolicy_QCP_LEGAL:            {"qcp-legal", "0.4.0.194112.1.1"},
	CertificatePolicy_QCP_NATURAL_QSCD:     {"qcp-natural-qscd", "0.4.0.194112.1.2"},
	CertificatePolicy_QCP_LEGAL_QSCD:       {"qcp-legal-qscd", "0.4.0.194112.1.3"},
	CertificatePolicy_QCP_WEB:              {"qcp-web", "0.4.0.194112.1.4"},
	CertificatePolicy_QNCP_WEB:             {"qncp-web", "0.4.0.194112.1.5"},
	CertificatePolicy_QNCP_WEB_GEN:         {"qncp-web-gen", "0.4.0.194112.1.6"},
}

// CertificatePolicyValues returns all constants in declaration order.
func CertificatePolicyValues() []CertificatePolicy {
	return []CertificatePolicy{
		CertificatePolicy_QCP_PUBLIC,
		CertificatePolicy_QCP_PUBLIC_WITH_SSCD,
		CertificatePolicy_NCP,
		CertificatePolicy_NCPP,
		CertificatePolicy_LCP,
		CertificatePolicy_EVCP,
		CertificatePolicy_DVCP,
		CertificatePolicy_OVCP,
		CertificatePolicy_IVCP,
		CertificatePolicy_QCP_NATURAL,
		CertificatePolicy_QCP_LEGAL,
		CertificatePolicy_QCP_NATURAL_QSCD,
		CertificatePolicy_QCP_LEGAL_QSCD,
		CertificatePolicy_QCP_WEB,
		CertificatePolicy_QNCP_WEB,
		CertificatePolicy_QNCP_WEB_GEN,
	}
}

// CertificatePolicyValueOf returns the CertificatePolicy matching the given Java enum name.
func CertificatePolicyValueOf(name string) (CertificatePolicy, error) {
	for _, v := range CertificatePolicyValues() {
		if string(v) == name {
			return v, nil
		}
	}
	return "", fmt.Errorf("no enum constant CertificatePolicy.%s", name)
}

// OID returns the OID. Implements OidDescription.
func (c CertificatePolicy) OID() string {
	return certificatePolicyData[c].oid
}

// Description returns the literal description. Implements OidDescription.
func (c CertificatePolicy) Description() string {
	return certificatePolicyData[c].description
}

// compile-time interface assertion.
var _ OidDescription = CertificatePolicy("")
