// Ported from dss-jades/src/main/java/eu/europa/esig/dss/jades/JWSJsonSerializationObject.java
// (DSS 6.5.RC1).
package jades

import (
	"github.com/utain/esig/dss/enumerations"
	"github.com/utain/esig/dss/utils"
)

// JWSJsonSerializationObject is a container with JWS signature attributes. Port of the class
// JWSJsonSerializationObject; java.io.Serializable is dropped (no Go counterpart).
type JWSJsonSerializationObject struct {
	// payload is the JWS payload. Port of the private field of the same name.
	payload string

	// signatures is the list of incorporated signatures. Port of the private field of the same
	// name.
	signatures []*JWS

	// jwsSerializationType defines the JWSSerializationType of the JAdES signature. Port of the
	// private field of the same name.
	jwsSerializationType enumerations.JWSSerializationType

	// structuralValidationErrors is a list of parsing errors, if any occurred. Port of the
	// private field of the same name.
	structuralValidationErrors []string
}

// NewJWSJsonSerializationObject instantiates the object with null values. Port of the default
// constructor.
func NewJWSJsonSerializationObject() *JWSJsonSerializationObject {
	return &JWSJsonSerializationObject{}
}

// Payload gets the base64url-encoded payload of a signature. Port of getPayload(), whose
// lazy "null becomes empty" is preserved by Go's zero value for a string.
func (o *JWSJsonSerializationObject) Payload() string {
	return o.payload
}

// SetPayload sets the base64url-encoded payload of a signature. Port of setPayload(String).
func (o *JWSJsonSerializationObject) SetPayload(encodedPayload string) {
	o.payload = encodedPayload
}

// Signatures returns the JWS signatures. Port of getSignatures(), which lazily creates the list
// so that a caller can add to it; Go returns the slice and callers append through SetSignatures,
// as DSSJsonUtilsToJWSJsonSerializationObject does.
func (o *JWSJsonSerializationObject) Signatures() []*JWS {
	return o.signatures
}

// SetSignatures sets the list of signatures. Port of setSignatures(List).
func (o *JWSJsonSerializationObject) SetSignatures(signatures []*JWS) {
	o.signatures = signatures
}

// AddSignature appends one signature.
//
// It has no Java counterpart because upstream reaches the same effect with
// getSignatures().add(jws) - a mutable view Go's slice return cannot provide, since appending to
// a copy would not be seen by the object.
func (o *JWSJsonSerializationObject) AddSignature(signature *JWS) {
	o.signatures = append(o.signatures, signature)
}

// JWSSerializationType gets the JWSSerializationType used for the signature. Port of
// getJWSSerializationType().
func (o *JWSJsonSerializationObject) JWSSerializationType() enumerations.JWSSerializationType {
	return o.jwsSerializationType
}

// SetJWSSerializationType sets the JWSSerializationType. Port of
// setJWSSerializationType(JWSSerializationType).
func (o *JWSJsonSerializationObject) SetJWSSerializationType(jwsSerializationType enumerations.JWSSerializationType) {
	o.jwsSerializationType = jwsSerializationType
}

// StructuralValidationErrors returns the errors that occurred during the structure (schema)
// validation, empty if none were found. Port of getStructuralValidationErrors().
func (o *JWSJsonSerializationObject) StructuralValidationErrors() []string {
	if utils.IsCollectionNotEmpty(o.structuralValidationErrors) {
		return o.structuralValidationErrors
	}
	return nil
}

// SetStructuralValidationErrors sets the errors that occurred during the structure (schema)
// validation. Port of setStructuralValidationErrors(List).
func (o *JWSJsonSerializationObject) SetStructuralValidationErrors(structuralValidationErrors []string) {
	o.structuralValidationErrors = structuralValidationErrors
}

// IsValid checks whether the signature structure validation succeeded. Port of isValid().
func (o *JWSJsonSerializationObject) IsValid() bool {
	return utils.IsCollectionEmpty(o.structuralValidationErrors)
}
