// Ported from dss-enumerations/.../SubContext.java (DSS 6.5.RC1).
package enumerations

// SubContext defines the validation sub context.
type SubContext string

const (
	// SubContext_SIGNING_CERT is used for a signing certificate validation.
	SubContext_SIGNING_CERT SubContext = "SIGNING_CERT"
	// SubContext_CA_CERTIFICATE is used for not trusted CA certificates in
	// a chain.
	SubContext_CA_CERTIFICATE SubContext = "CA_CERTIFICATE"
)

// SubContextValues returns all constants in declaration order.
func SubContextValues() []SubContext {
	return []SubContext{
		SubContext_SIGNING_CERT,
		SubContext_CA_CERTIFICATE,
	}
}

// SubContextValueOf returns the constant matching the given Java enum name.
func SubContextValueOf(name string) (SubContext, error) {
	for _, v := range SubContextValues() {
		if string(v) == name {
			return v, nil
		}
	}
	return "", &subContextInvalidValueError{name}
}

type subContextInvalidValueError struct {
	name string
}

func (e *subContextInvalidValueError) Error() string {
	return "no enum constant SubContext." + e.name
}
