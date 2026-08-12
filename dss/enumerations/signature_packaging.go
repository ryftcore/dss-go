// Ported from dss-enumerations/.../SignaturePackaging.java (DSS 6.5.RC1).

package enumerations

import "fmt"

// SignaturePackaging is the packaging method of the signature.
type SignaturePackaging string

const (
	// SignaturePackaging_ENVELOPED: the signature is enveloped to the signed document.
	SignaturePackaging_ENVELOPED SignaturePackaging = "ENVELOPED"
	// SignaturePackaging_ENVELOPING: the signature envelops the signed document.
	SignaturePackaging_ENVELOPING SignaturePackaging = "ENVELOPING"
	// SignaturePackaging_DETACHED: the signature is detached from the signed document.
	SignaturePackaging_DETACHED SignaturePackaging = "DETACHED"
	// SignaturePackaging_INTERNALLY_DETACHED: the signature file contains the signed
	// document (XAdES only).
	SignaturePackaging_INTERNALLY_DETACHED SignaturePackaging = "INTERNALLY_DETACHED"
)

// SignaturePackagingValues returns all SignaturePackaging constants in declaration order.
func SignaturePackagingValues() []SignaturePackaging {
	return []SignaturePackaging{
		SignaturePackaging_ENVELOPED,
		SignaturePackaging_ENVELOPING,
		SignaturePackaging_DETACHED,
		SignaturePackaging_INTERNALLY_DETACHED,
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
