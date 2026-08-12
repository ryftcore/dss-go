// Ported from dss-enumerations/.../SignatureScopeType.java (DSS 6.5.RC1).
package enumerations

// SignatureScopeType defines the used SignatureScope types.
type SignatureScopeType string

const (
	// SignatureScopeType_FULL means the signature covers the complete
	// document.
	SignatureScopeType_FULL SignatureScopeType = "FULL"
	// SignatureScopeType_PARTIAL means the signature covers only a part of
	// the document.
	SignatureScopeType_PARTIAL SignatureScopeType = "PARTIAL"
	// SignatureScopeType_DIGEST means the signature covers only the digest
	// of document.
	SignatureScopeType_DIGEST SignatureScopeType = "DIGEST"
	// SignatureScopeType_ARCHIVED means the signature covers its bounded
	// archive.
	SignatureScopeType_ARCHIVED SignatureScopeType = "ARCHIVED"
	// SignatureScopeType_COUNTER_SIGNATURE means the signature
	// counter-signs its master signature.
	SignatureScopeType_COUNTER_SIGNATURE SignatureScopeType = "COUNTER_SIGNATURE"
	// SignatureScopeType_EAA_SIGNATURE is the signature used to issue the
	// EAA.
	SignatureScopeType_EAA_SIGNATURE SignatureScopeType = "EAA_SIGNATURE"
	// SignatureScopeType_KEY_BINDING_SIGNATURE is the key binding signature
	// used to proof a possession of the key by a Wallet holder.
	SignatureScopeType_KEY_BINDING_SIGNATURE SignatureScopeType = "KEY_BINDING_SIGNATURE"
	// SignatureScopeType_SIGNATURE means the evidence record covers a
	// signature.
	SignatureScopeType_SIGNATURE SignatureScopeType = "SIGNATURE"
)

// SignatureScopeTypeValues returns all constants in declaration order.
func SignatureScopeTypeValues() []SignatureScopeType {
	return []SignatureScopeType{
		SignatureScopeType_FULL,
		SignatureScopeType_PARTIAL,
		SignatureScopeType_DIGEST,
		SignatureScopeType_ARCHIVED,
		SignatureScopeType_COUNTER_SIGNATURE,
		SignatureScopeType_EAA_SIGNATURE,
		SignatureScopeType_KEY_BINDING_SIGNATURE,
		SignatureScopeType_SIGNATURE,
	}
}

// SignatureScopeTypeValueOf returns the constant matching the given Java
// enum name.
func SignatureScopeTypeValueOf(name string) (SignatureScopeType, error) {
	for _, v := range SignatureScopeTypeValues() {
		if string(v) == name {
			return v, nil
		}
	}
	return "", &signatureScopeTypeInvalidValueError{name}
}

type signatureScopeTypeInvalidValueError struct {
	name string
}

func (e *signatureScopeTypeInvalidValueError) Error() string {
	return "no enum constant SignatureScopeType." + e.name
}
