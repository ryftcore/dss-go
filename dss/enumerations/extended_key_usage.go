// Ported from dss-enumerations/.../ExtendedKeyUsage.java (DSS 6.5.RC1).
package enumerations

// ExtendedKeyUsage is the KeyPurposeId object.
//
//	KeyPurposeId ::= OBJECT IDENTIFIER
//
//	id-kp ::= OBJECT IDENTIFIER { iso(1) identified-organization(3)
//	dod(6) internet(1) security(5) mechanisms(5) pkix(7) 3}
type ExtendedKeyUsage string

const (
	// ExtendedKeyUsage_SERVER_AUTH is serverAuth.
	ExtendedKeyUsage_SERVER_AUTH ExtendedKeyUsage = "SERVER_AUTH"
	// ExtendedKeyUsage_CLIENT_AUTH is clientAuth.
	ExtendedKeyUsage_CLIENT_AUTH ExtendedKeyUsage = "CLIENT_AUTH"
	// ExtendedKeyUsage_CODE_SIGNING is codeSigning.
	ExtendedKeyUsage_CODE_SIGNING ExtendedKeyUsage = "CODE_SIGNING"
	// ExtendedKeyUsage_EMAIL_PROTECTION is emailProtection.
	ExtendedKeyUsage_EMAIL_PROTECTION ExtendedKeyUsage = "EMAIL_PROTECTION"
	// ExtendedKeyUsage_TIMESTAMPING is timeStamping.
	// 5,6,7 deprecated by RFC4945.
	ExtendedKeyUsage_TIMESTAMPING ExtendedKeyUsage = "TIMESTAMPING"
	// ExtendedKeyUsage_OCSP_SIGNING is ocspSigning.
	ExtendedKeyUsage_OCSP_SIGNING ExtendedKeyUsage = "OCSP_SIGNING"
	// ExtendedKeyUsage_TSL_SIGNING is the ETSI TS 119 612 OID for TSL
	// signing KeyPurposeID for ExtKeyUsageSyntax:
	//
	//	id-tsl OBJECT IDENTIFIER { itu-t(0) identified-organization(4)
	//	etsi(0) tsl-specification (2231) }
	//	id-tsl-kp OBJECT IDENTIFIER ::= { id-tsl kp(3) }
	//	id-tsl-kp-tslSigning OBJECT IDENTIFIER ::= { id-tsl-kp tsl-signing(0) }
	ExtendedKeyUsage_TSL_SIGNING ExtendedKeyUsage = "TSL_SIGNING"
	// ExtendedKeyUsage_TSL_BINDING is defined by ETSI TS 119 411-5
	// (Annex A (normative): id-kp-tlsBinding EKU specification), defining
	// the KeyPurposeID id-kp-tls-binding:
	//
	//	id-tlsBinding OBJECT IDENTIFIER ::= { itu-t(0) identified-organization(4)
	//	etsi(0) id-qwacImplementation(194115) tls-binding (1) }
	//	id-kp-tls-binding OBJECT IDENTIFIER ::= { id-tlsBinding id-kp-tls-binding(0) }
	ExtendedKeyUsage_TSL_BINDING ExtendedKeyUsage = "TSL_BINDING"
)

type extendedKeyUsageFields struct {
	description string
	oid         string
}

// extendedKeyUsageData holds the (description, oid) tuple for each constant.
var extendedKeyUsageData = map[ExtendedKeyUsage]extendedKeyUsageFields{
	ExtendedKeyUsage_SERVER_AUTH:      {"serverAuth", "1.3.6.1.5.5.7.3.1"},
	ExtendedKeyUsage_CLIENT_AUTH:      {"clientAuth", "1.3.6.1.5.5.7.3.2"},
	ExtendedKeyUsage_CODE_SIGNING:     {"codeSigning", "1.3.6.1.5.5.7.3.3"},
	ExtendedKeyUsage_EMAIL_PROTECTION: {"emailProtection", "1.3.6.1.5.5.7.3.4"},
	ExtendedKeyUsage_TIMESTAMPING:     {"timeStamping", "1.3.6.1.5.5.7.3.8"},
	ExtendedKeyUsage_OCSP_SIGNING:     {"ocspSigning", "1.3.6.1.5.5.7.3.9"},
	ExtendedKeyUsage_TSL_SIGNING:      {"tslSigning", "0.4.0.2231.3.0"},
	ExtendedKeyUsage_TSL_BINDING:      {"tslBinding", "0.4.0.194115.1.0"},
}

// ExtendedKeyUsageValues returns all constants in declaration order.
func ExtendedKeyUsageValues() []ExtendedKeyUsage {
	return []ExtendedKeyUsage{
		ExtendedKeyUsage_SERVER_AUTH,
		ExtendedKeyUsage_CLIENT_AUTH,
		ExtendedKeyUsage_CODE_SIGNING,
		ExtendedKeyUsage_EMAIL_PROTECTION,
		ExtendedKeyUsage_TIMESTAMPING,
		ExtendedKeyUsage_OCSP_SIGNING,
		ExtendedKeyUsage_TSL_SIGNING,
		ExtendedKeyUsage_TSL_BINDING,
	}
}

// OID returns the ExtendedKeyUsage OID.
func (e ExtendedKeyUsage) OID() string {
	return extendedKeyUsageData[e].oid
}

// Description returns the ExtendedKeyUsage description.
func (e ExtendedKeyUsage) Description() string {
	return extendedKeyUsageData[e].description
}

// ExtendedKeyUsageValueOf returns the constant matching the given Java enum
// name.
func ExtendedKeyUsageValueOf(name string) (ExtendedKeyUsage, error) {
	for _, v := range ExtendedKeyUsageValues() {
		if string(v) == name {
			return v, nil
		}
	}
	return "", &extendedKeyUsageInvalidValueError{name}
}

type extendedKeyUsageInvalidValueError struct {
	name string
}

func (e *extendedKeyUsageInvalidValueError) Error() string {
	return "no enum constant ExtendedKeyUsage." + e.name
}

// Compile-time assertion that ExtendedKeyUsage implements OidDescription.
var _ OidDescription = ExtendedKeyUsage_SERVER_AUTH
