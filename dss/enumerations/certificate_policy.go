// Ported from dss-enumerations/.../CertificatePolicy.java (DSS 6.5.RC1).
package enumerations

import "fmt"

// CertificatePolicy is the certificate policies enumeration, as in ETSI TS
// 101 456. Implements OidDescription.
type CertificatePolicy string

const (
	// CertificatePolicyQCPPublic is a certificate policy for qualified
	// certificates issued to the public.
	CertificatePolicyQCPPublic CertificatePolicy = "QCP_PUBLIC"
	// CertificatePolicyQCPPublicWithSSCD is a certificate policy for
	// qualified certificates issued to the public, requiring use of secure
	// signature-creation devices.
	CertificatePolicyQCPPublicWithSSCD CertificatePolicy = "QCP_PUBLIC_WITH_SSCD"
	// CertificatePolicyNCP: Normalized Certificate Policy.
	CertificatePolicyNCP CertificatePolicy = "NCP"
	// CertificatePolicyNCPP: Normalized Certificate Policy requiring a
	// secure cryptographic device.
	CertificatePolicyNCPP CertificatePolicy = "NCPP"
	// CertificatePolicyLCP: Lightweight Certificate Policy.
	CertificatePolicyLCP CertificatePolicy = "LCP"
	// CertificatePolicyEVCP: Extended Validation Certificate Policy.
	CertificatePolicyEVCP CertificatePolicy = "EVCP"
	// CertificatePolicyDVCP: Domain Validation Certificate Policy.
	CertificatePolicyDVCP CertificatePolicy = "DVCP"
	// CertificatePolicyOVCP: Organizational Validation Certificate Policy.
	CertificatePolicyOVCP CertificatePolicy = "OVCP"
	// CertificatePolicyIVCP: Individual Validation Certificate Policy.
	CertificatePolicyIVCP CertificatePolicy = "IVCP"
	// CertificatePolicyQCPNatural: certificate policy for EU qualified
	// certificates issued to natural persons.
	CertificatePolicyQCPNatural CertificatePolicy = "QCP_NATURAL"
	// CertificatePolicyQCPLegal: certificate policy for EU qualified
	// certificates issued to legal persons.
	CertificatePolicyQCPLegal CertificatePolicy = "QCP_LEGAL"
	// CertificatePolicyQCPNaturalQSCD: certificate policy for EU
	// qualified certificates issued to natural persons with private key
	// related to the certified public key in a QSCD.
	CertificatePolicyQCPNaturalQSCD CertificatePolicy = "QCP_NATURAL_QSCD"
	// CertificatePolicyQCPLegalQSCD: certificate policy for EU qualified
	// certificates issued to legal persons with private key related to the
	// certified public key in a QSCD.
	CertificatePolicyQCPLegalQSCD CertificatePolicy = "QCP_LEGAL_QSCD"
	// CertificatePolicyQCPWeb: certificate policy for EU qualified
	// website authentication certificates based on EVCP.
	CertificatePolicyQCPWeb CertificatePolicy = "QCP_WEB"
	// CertificatePolicyQNCPWeb: certificate policy for EU qualified
	// website authentication certificates based on NCP, and OVCP or IVCP.
	CertificatePolicyQNCPWeb CertificatePolicy = "QNCP_WEB"
	// CertificatePolicyQNCPWebGen: certificate policy for EU qualified
	// website authentication based on NCP and requirements tagged as [WEB]
	// in ETSI EN 319 411-1.
	CertificatePolicyQNCPWebGen CertificatePolicy = "QNCP_WEB_GEN"
)

// certificatePolicyFields holds the (description, oid) pair for each constant.
type certificatePolicyFields struct {
	description string
	oid         string
}

// certificatePolicyData holds the full field tuple for each constant, copied
// verbatim from the Java enum constructors.
var certificatePolicyData = map[CertificatePolicy]certificatePolicyFields{
	CertificatePolicyQCPPublic:         {"qcp-public", "0.4.0.1456.1.2"},
	CertificatePolicyQCPPublicWithSSCD: {"qcp-public-with-sscd", "0.4.0.1456.1.1"},
	CertificatePolicyNCP:               {"normalized-certificate-policy", "0.4.0.2042.1.1"},
	CertificatePolicyNCPP:              {"normalized-certificate-policy-sscd", "0.4.0.2042.1.2"},
	CertificatePolicyLCP:               {"lightweight-certificate-policy", "0.4.0.2042.1.3"},
	CertificatePolicyEVCP:              {"extended-validation-certificate-policy", "0.4.0.2042.1.4"},
	CertificatePolicyDVCP:              {"domain-validation-certificate-policy", "0.4.0.2042.1.6"},
	CertificatePolicyOVCP:              {"organizational-validation-certificate-policy", "0.4.0.2042.1.7"},
	CertificatePolicyIVCP:              {"individual-validation-certificate-policy", "0.4.0.2042.1.8"},
	CertificatePolicyQCPNatural:        {"qcp-natural", "0.4.0.194112.1.0"},
	CertificatePolicyQCPLegal:          {"qcp-legal", "0.4.0.194112.1.1"},
	CertificatePolicyQCPNaturalQSCD:    {"qcp-natural-qscd", "0.4.0.194112.1.2"},
	CertificatePolicyQCPLegalQSCD:      {"qcp-legal-qscd", "0.4.0.194112.1.3"},
	CertificatePolicyQCPWeb:            {"qcp-web", "0.4.0.194112.1.4"},
	CertificatePolicyQNCPWeb:           {"qncp-web", "0.4.0.194112.1.5"},
	CertificatePolicyQNCPWebGen:        {"qncp-web-gen", "0.4.0.194112.1.6"},
}

// CertificatePolicyValues returns all constants in declaration order.
func CertificatePolicyValues() []CertificatePolicy {
	return []CertificatePolicy{
		CertificatePolicyQCPPublic,
		CertificatePolicyQCPPublicWithSSCD,
		CertificatePolicyNCP,
		CertificatePolicyNCPP,
		CertificatePolicyLCP,
		CertificatePolicyEVCP,
		CertificatePolicyDVCP,
		CertificatePolicyOVCP,
		CertificatePolicyIVCP,
		CertificatePolicyQCPNatural,
		CertificatePolicyQCPLegal,
		CertificatePolicyQCPNaturalQSCD,
		CertificatePolicyQCPLegalQSCD,
		CertificatePolicyQCPWeb,
		CertificatePolicyQNCPWeb,
		CertificatePolicyQNCPWebGen,
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
