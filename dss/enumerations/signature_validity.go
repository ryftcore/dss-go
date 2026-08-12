// Ported from dss-enumerations/.../SignatureValidity.java (DSS 6.5.RC1).
//
// Defines result of signature validation for a token.
package enumerations

import "fmt"

// SignatureValidity defines the result of signature validation for a token.
type SignatureValidity string

const (
	// SignatureValidity_VALID: the signature of the token is valid
	// (signing certificate found successfully).
	SignatureValidity_VALID SignatureValidity = "VALID"
	// SignatureValidity_INVALID: the signature of the token is invalid.
	SignatureValidity_INVALID SignatureValidity = "INVALID"
	// SignatureValidity_NOT_EVALUATED: the signature of the token is not
	// evaluated yet.
	SignatureValidity_NOT_EVALUATED SignatureValidity = "NOT_EVALUATED"
)

// SignatureValidityValues returns all constants in declaration order.
func SignatureValidityValues() []SignatureValidity {
	return []SignatureValidity{
		SignatureValidity_VALID,
		SignatureValidity_INVALID,
		SignatureValidity_NOT_EVALUATED,
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
// pointer maps to SignatureValidity_NOT_EVALUATED.
func SignatureValidityGet(isValid *bool) SignatureValidity {
	if isValid == nil {
		return SignatureValidity_NOT_EVALUATED
	} else if *isValid {
		return SignatureValidity_VALID
	}
	return SignatureValidity_INVALID
}
