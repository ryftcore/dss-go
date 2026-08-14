// Ported from dss-enumerations/.../SignatureTokenType.java (DSS 6.5.RC1).
package enumerations

// SignatureTokenType is the type of API supported by the application.
type SignatureTokenType string

const (
	// SignatureTokenType_PKCS11 is a PKCS11 keystore token.
	SignatureTokenType_PKCS11 SignatureTokenType = "PKCS11"
	// SignatureTokenType_PKCS12 is a PKCS12 keystore token.
	SignatureTokenType_PKCS12 SignatureTokenType = "PKCS12"
	// SignatureTokenType_MSCAPI is a Windows keystore token.
	SignatureTokenType_MSCAPI SignatureTokenType = "MSCAPI"
	// SignatureTokenType_APPLE is a MacOS keystore token.
	SignatureTokenType_APPLE SignatureTokenType = "APPLE"
	// SignatureTokenType_JKS is a Java keystore token.
	SignatureTokenType_JKS SignatureTokenType = "JKS"
	// SignatureTokenType_MOCCA is a MOCCA-implementation keystore token.
	SignatureTokenType_MOCCA SignatureTokenType = "MOCCA"
)

// SignatureTokenTypeValues returns all constants in declaration order.
func SignatureTokenTypeValues() []SignatureTokenType {
	return []SignatureTokenType{
		SignatureTokenType_PKCS11,
		SignatureTokenType_PKCS12,
		SignatureTokenType_MSCAPI,
		SignatureTokenType_APPLE,
		SignatureTokenType_JKS,
		SignatureTokenType_MOCCA,
	}
}
