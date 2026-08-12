// Ported from dss-enumerations/.../SignaturePolicyType.java (DSS 6.5.RC1).

package enumerations

import "fmt"

// SignaturePolicyType defines available signature policy types.
type SignaturePolicyType string

const (
	// SignaturePolicyType_NO_POLICY: the validation process accepts no policy. No
	// particular treatment is done.
	SignaturePolicyType_NO_POLICY SignaturePolicyType = "NO_POLICY"
	// SignaturePolicyType_ANY_POLICY: the validation process accepts any policy. The
	// used policy is only showed, no particular treatment is done.
	SignaturePolicyType_ANY_POLICY SignaturePolicyType = "ANY_POLICY"
	// SignaturePolicyType_IMPLICIT_POLICY: indicate that the data object(s) being
	// signed and other external data imply the signature policy.
	SignaturePolicyType_IMPLICIT_POLICY SignaturePolicyType = "IMPLICIT_POLICY"
)

// SignaturePolicyTypeValues returns all SignaturePolicyType constants in declaration order.
func SignaturePolicyTypeValues() []SignaturePolicyType {
	return []SignaturePolicyType{
		SignaturePolicyType_NO_POLICY,
		SignaturePolicyType_ANY_POLICY,
		SignaturePolicyType_IMPLICIT_POLICY,
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
