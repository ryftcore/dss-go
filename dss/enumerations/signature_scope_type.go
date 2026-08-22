// Ported from dss-enumerations/.../SignatureScopeType.java (DSS 6.5.RC1).
package enumerations

// SignatureScopeType defines the used SignatureScope types.
type SignatureScopeType string

const (
	// SignatureScopeTypeFull means the signature covers the complete
	// document.
	SignatureScopeTypeFull SignatureScopeType = "FULL"
	// SignatureScopeTypePartial means the signature covers only a part of
	// the document.
	SignatureScopeTypePartial SignatureScopeType = "PARTIAL"
	// SignatureScopeTypeDigest means the signature covers only the digest
	// of document.
	SignatureScopeTypeDigest SignatureScopeType = "DIGEST"
	// SignatureScopeTypeArchived means the signature covers its bounded
	// archive.
	SignatureScopeTypeArchived SignatureScopeType = "ARCHIVED"
	// SignatureScopeTypeCounterSignature means the signature
	// counter-signs its master signature.
	SignatureScopeTypeCounterSignature SignatureScopeType = "COUNTER_SIGNATURE"
	// SignatureScopeTypeEAASignature is the signature used to issue the
	// EAA.
	SignatureScopeTypeEAASignature SignatureScopeType = "EAA_SIGNATURE"
	// SignatureScopeTypeKeyBindingSignature is the key binding signature
	// used to proof a possession of the key by a Wallet holder.
	SignatureScopeTypeKeyBindingSignature SignatureScopeType = "KEY_BINDING_SIGNATURE"
	// SignatureScopeTypeSignature means the evidence record covers a
	// signature.
	SignatureScopeTypeSignature SignatureScopeType = "SIGNATURE"
)

// SignatureScopeTypeValues returns all constants in declaration order.
func SignatureScopeTypeValues() []SignatureScopeType {
	return []SignatureScopeType{
		SignatureScopeTypeFull,
		SignatureScopeTypePartial,
		SignatureScopeTypeDigest,
		SignatureScopeTypeArchived,
		SignatureScopeTypeCounterSignature,
		SignatureScopeTypeEAASignature,
		SignatureScopeTypeKeyBindingSignature,
		SignatureScopeTypeSignature,
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
