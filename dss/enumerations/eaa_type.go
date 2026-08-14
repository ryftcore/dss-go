// Ported from dss-enumerations/.../EAAType.java (DSS 6.5.RC1).
package enumerations

// EAAType defines a list of EAA types known by the current implementation.
// NOTE: This type relates to a format of an EAA signed token.
type EAAType string

const (
	// EAAType_SD_JWT_VC is a realization of EAA that implements EAA as a
	// JSON Web Signature as specified in IETF RFC 7515, built on IETF
	// SD-JWT VC, which further profiles a Selective Disclosure JSON Web
	// Token as specified in IETF RFC 9901 "Selective Disclosure for JSON
	// Web Tokens".
	EAAType_SD_JWT_VC EAAType = "SD_JWT_VC"
	// EAAType_ISO_IEC_MDOC is a realization of EAA that implements EAA
	// built on the data structures defined in ISO/IEC 18013-5.
	EAAType_ISO_IEC_MDOC EAAType = "ISO_IEC_MDOC"
	// EAAType_W3C_VC is a realization of EAA based on the JSON-LD
	// (specified in W3C Recommendation: "JSON-LD 1.1. A JSON-based
	// Serialization for Linked Data") serialization of W3C Recommendation
	// (15 May 2025): "Verifiable Credentials Data Model v2.0".
	EAAType_W3C_VC EAAType = "W3C_VC"
	// EAAType_X509_AC is a realization of EAA based on X.509 Attribute
	// certificates as specified in IETF RFC 5755.
	EAAType_X509_AC EAAType = "X509_AC"
)

// EAATypeValues returns all constants in declaration order.
func EAATypeValues() []EAAType {
	return []EAAType{
		EAAType_SD_JWT_VC,
		EAAType_ISO_IEC_MDOC,
		EAAType_W3C_VC,
		EAAType_X509_AC,
	}
}
