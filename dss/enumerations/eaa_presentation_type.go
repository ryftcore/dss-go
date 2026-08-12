// Ported from dss-enumerations/.../EAAPresentationType.java (DSS 6.5.RC1).

package enumerations

import "fmt"

// EAAPresentationType defines a list of format types on which supported EAA
// presentations may be based.
//
// NOTE: This type relates to a format of an EAA presentation document.
type EAAPresentationType string

const (
	// EAAPresentationType_SD_JWT represents an IETF RFC 9901 "Selective Disclosure for
	// JSON Web Tokens" token.
	EAAPresentationType_SD_JWT EAAPresentationType = "SD_JWT"
	// EAAPresentationType_MDOC_DEVICE_RESPONSE represents a DeviceResponse mdoc
	// structure as per ISO/IEC 18013-5 "8.3.2.1.2.2 Device retrieval mdoc response".
	EAAPresentationType_MDOC_DEVICE_RESPONSE EAAPresentationType = "MDOC_DEVICE_RESPONSE"
	// EAAPresentationType_MDOC_ISSUER_SIGNED represents an IssuerSigned mdoc structure
	// as per ISO/IEC 18013-5 "8.3.2.1.2.2 Device retrieval mdoc response".
	EAAPresentationType_MDOC_ISSUER_SIGNED EAAPresentationType = "MDOC_ISSUER_SIGNED"
	// EAAPresentationType_JWS represents a JOSE token, as defined in IETF RFC 7515
	// "JSON Web Signature (JWS)".
	EAAPresentationType_JWS EAAPresentationType = "JWS"
	// EAAPresentationType_X509_AC is the realization of EAA based on X.509 Attribute
	// certificates as specified in IETF RFC 5755.
	EAAPresentationType_X509_AC EAAPresentationType = "X509_AC"
)

// EAAPresentationTypeValues returns all EAAPresentationType constants in declaration order.
func EAAPresentationTypeValues() []EAAPresentationType {
	return []EAAPresentationType{
		EAAPresentationType_SD_JWT,
		EAAPresentationType_MDOC_DEVICE_RESPONSE,
		EAAPresentationType_MDOC_ISSUER_SIGNED,
		EAAPresentationType_JWS,
		EAAPresentationType_X509_AC,
	}
}

// EAAPresentationTypeValueOf returns the EAAPresentationType matching the given Java enum name.
func EAAPresentationTypeValueOf(name string) (EAAPresentationType, error) {
	for _, v := range EAAPresentationTypeValues() {
		if string(v) == name {
			return v, nil
		}
	}
	return "", fmt.Errorf("no enum constant EAAPresentationType.%s", name)
}
