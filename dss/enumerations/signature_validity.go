// Ported from dss-enumerations/.../SignatureValidity.java (DSS 6.5.RC1).
//
// Defines result of signature validation for a token.
package enumerations

import "fmt"

// SignatureValidity defines the result of signature validation for a token.
type SignatureValidity string

const (
	// SignatureValidityValid: the signature of the token is valid
	// (signing certificate found successfully).
	SignatureValidityValid SignatureValidity = "VALID"
	// SignatureValidityInvalid: the signature of the token is invalid.
	SignatureValidityInvalid SignatureValidity = "INVALID"
	// SignatureValidityNotEvaluated: the signature of the token is not
	// evaluated yet.
	SignatureValidityNotEvaluated SignatureValidity = "NOT_EVALUATED"
)

// SignatureValidityValues returns all constants in declaration order.
func SignatureValidityValues() []SignatureValidity {
	return []SignatureValidity{
		SignatureValidityValid,
		SignatureValidityInvalid,
		SignatureValidityNotEvaluated,
	}
}

// SignatureValidityValueOf returns the SignatureValidity matching the given Java enum name.
func SignatureValidityValueOf(name string) (SignatureValidity, error) {
	for _, v := range SignatureValidityValues() {
		if string(v) == name {
			return v, nil
		}
	}
	return "", fmt.Errorf("no enum constant SignatureValidity.%s", name)
}

// SignatureValidityGet returns the SignatureValidity type matching the given
// value. isValid is a pointer to mirror Java's nullable Boolean: a nil
// pointer maps to SignatureValidityNotEvaluated.
func SignatureValidityGet(isValid *bool) SignatureValidity {
	if isValid == nil {
		return SignatureValidityNotEvaluated
	} else if *isValid {
		return SignatureValidityValid
	}
	return SignatureValidityInvalid
}
