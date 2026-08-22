// Ported from dss-enumerations/.../SignatureTokenType.java (DSS 6.5.RC1).
package enumerations

// SignatureTokenType is the type of API supported by the application.
type SignatureTokenType string

const (
	// SignatureTokenTypePKCS11 is a PKCS11 keystore token.
	SignatureTokenTypePKCS11 SignatureTokenType = "PKCS11"
	// SignatureTokenTypePKCS12 is a PKCS12 keystore token.
	SignatureTokenTypePKCS12 SignatureTokenType = "PKCS12"
	// SignatureTokenTypeMSCAPI is a Windows keystore token.
	SignatureTokenTypeMSCAPI SignatureTokenType = "MSCAPI"
	// SignatureTokenTypeApple is a MacOS keystore token.
	SignatureTokenTypeApple SignatureTokenType = "APPLE"
	// SignatureTokenTypeJKS is a Java keystore token.
	SignatureTokenTypeJKS SignatureTokenType = "JKS"
	// SignatureTokenTypeMOCCA is a MOCCA-implementation keystore token.
	SignatureTokenTypeMOCCA SignatureTokenType = "MOCCA"
)

// SignatureTokenTypeValues returns all constants in declaration order.
func SignatureTokenTypeValues() []SignatureTokenType {
	return []SignatureTokenType{
		SignatureTokenTypePKCS11,
		SignatureTokenTypePKCS12,
		SignatureTokenTypeMSCAPI,
		SignatureTokenTypeApple,
		SignatureTokenTypeJKS,
		SignatureTokenTypeMOCCA,
	}
}
