// Ported from dss-enumerations/.../RevocationType.java (DSS 6.5.RC1).
package enumerations

// RevocationType defines a type of revocation data response.
type RevocationType string

const (
	// RevocationType_CRL indicates the revocation data was received from
	// CRL response.
	RevocationType_CRL RevocationType = "CRL"
	// RevocationType_OCSP indicates the revocation data was received from
	// OCSP response.
	RevocationType_OCSP RevocationType = "OCSP"
)

// RevocationTypeValues returns all constants in declaration order.
func RevocationTypeValues() []RevocationType {
	return []RevocationType{
		RevocationType_CRL,
		RevocationType_OCSP,
	}
}
