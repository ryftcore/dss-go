// Ported from dss-enumerations/.../RevocationType.java (DSS 6.5.RC1).
package enumerations

// RevocationType defines a type of revocation data response.
type RevocationType string

const (
	// RevocationTypeCRL indicates the revocation data was received from
	// CRL response.
	RevocationTypeCRL RevocationType = "CRL"
	// RevocationTypeOCSP indicates the revocation data was received from
	// OCSP response.
	RevocationTypeOCSP RevocationType = "OCSP"
)

// RevocationTypeValues returns all constants in declaration order.
func RevocationTypeValues() []RevocationType {
	return []RevocationType{
		RevocationTypeCRL,
		RevocationTypeOCSP,
	}
}
