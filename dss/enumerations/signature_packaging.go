// Ported from dss-enumerations/.../SignaturePackaging.java (DSS 6.5.RC1).

package enumerations

import "fmt"

// SignaturePackaging is the packaging method of the signature.
type SignaturePackaging string

const (
	// SignaturePackagingEnveloped: the signature is enveloped to the signed document.
	SignaturePackagingEnveloped SignaturePackaging = "ENVELOPED"
	// SignaturePackagingEnveloping: the signature envelops the signed document.
	SignaturePackagingEnveloping SignaturePackaging = "ENVELOPING"
	// SignaturePackagingDetached: the signature is detached from the signed document.
	SignaturePackagingDetached SignaturePackaging = "DETACHED"
	// SignaturePackagingInternallyDetached: the signature file contains the signed
	// document (XAdES only).
	SignaturePackagingInternallyDetached SignaturePackaging = "INTERNALLY_DETACHED"
)

// SignaturePackagingValues returns all SignaturePackaging constants in declaration order.
func SignaturePackagingValues() []SignaturePackaging {
	return []SignaturePackaging{
		SignaturePackagingEnveloped,
		SignaturePackagingEnveloping,
		SignaturePackagingDetached,
		SignaturePackagingInternallyDetached,
	}
}

// SignaturePackagingValueOf returns the SignaturePackaging matching the given Java enum name.
func SignaturePackagingValueOf(name string) (SignaturePackaging, error) {
	for _, v := range SignaturePackagingValues() {
		if string(v) == name {
			return v, nil
		}
	}
	return "", fmt.Errorf("no enum constant SignaturePackaging.%s", name)
}
