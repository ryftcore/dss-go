// Ported from dss-enumerations/.../SignaturePolicyType.java (DSS 6.5.RC1).

package enumerations

import "fmt"

// SignaturePolicyType defines available signature policy types.
type SignaturePolicyType string

const (
	// SignaturePolicyTypeNoPolicy: the validation process accepts no policy. No
	// particular treatment is done.
	SignaturePolicyTypeNoPolicy SignaturePolicyType = "NO_POLICY"
	// SignaturePolicyTypeAnyPolicy: the validation process accepts any policy. The
	// used policy is only showed, no particular treatment is done.
	SignaturePolicyTypeAnyPolicy SignaturePolicyType = "ANY_POLICY"
	// SignaturePolicyTypeImplicitPolicy: indicate that the data object(s) being
	// signed and other external data imply the signature policy.
	SignaturePolicyTypeImplicitPolicy SignaturePolicyType = "IMPLICIT_POLICY"
)

// SignaturePolicyTypeValues returns all SignaturePolicyType constants in declaration order.
func SignaturePolicyTypeValues() []SignaturePolicyType {
	return []SignaturePolicyType{
		SignaturePolicyTypeNoPolicy,
		SignaturePolicyTypeAnyPolicy,
		SignaturePolicyTypeImplicitPolicy,
	}
}

// SignaturePolicyTypeValueOf returns the SignaturePolicyType matching the given Java enum name.
func SignaturePolicyTypeValueOf(name string) (SignaturePolicyType, error) {
	for _, v := range SignaturePolicyTypeValues() {
		if string(v) == name {
			return v, nil
		}
	}
	return "", fmt.Errorf("no enum constant SignaturePolicyType.%s", name)
}
