// Ported from dss-model/.../SerializableCounterSignatureParameters.java (DSS 6.5.RC1).
package model

// SerializableCounterSignatureParameters contains the common methods for
// counter signature parameters.
type SerializableCounterSignatureParameters interface {
	SerializableSignatureParameters

	// SignatureIdToCounterSign returns the Id of a signature that needs
	// to be counter signed. Ports #getSignatureIdToCounterSign.
	SignatureIdToCounterSign() string

	// SetSignatureIdToCounterSign sets the Id of a signature to be
	// counter signed.
	//
	// NOTE: the id shall represent the DSS (hash-based) id of a signature
	// or a provided id in the signature document, when available (i.e.
	// XML Id for a XAdES signature).
	SetSignatureIdToCounterSign(signatureId string)
}
