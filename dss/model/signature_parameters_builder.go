// Ported from dss-model/.../SignatureParametersBuilder.java (DSS 6.5.RC1).
package model

// SignatureParametersBuilder hides the complexity of a configuration for
// particular usages and simplifies signature creation, generic over the
// SerializableSignatureParameters implementation SP to be created.
//
// SerializableSignatureParameters is outside this manifest; assumed to
// already exist in this package.
type SignatureParametersBuilder[SP SerializableSignatureParameters] interface {
	// Build creates a Signature Parameters instance.
	Build() SP
}
