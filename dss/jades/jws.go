// Ported from dss-jades/src/main/java/eu/europa/esig/dss/jades/validation/JWS.java (DSS 6.5.RC1).
//
// The Java package eu.europa.esig.dss.jades.validation is flattened into this one Go package.
//
// Upstream JWS extends org.jose4j.jws.JsonWebSignature. There is no jose4j in Go, so the
// superclass is internal/jose's own JWS - a port of exactly the JsonWebSignature/JsonWebStructure
// surface dss-jades reaches - and this type embeds it. Everything the Java class inherits without
// overriding is therefore reached through the embedded value; everything it declares or overrides
// is below.
package jades

import (
	"github.com/ryftcore/dss-go/dss/enumerations"
	"github.com/ryftcore/dss-go/dss/internal/jose"
	"github.com/ryftcore/dss-go/dss/spi/exception"
)

// JWS is an extension of a JSON Web Signature according to RFC 7515. Port of the class JWS.
//
// java.io.Serializable is dropped (no Go counterpart).
type JWS struct {
	// JWS is the jose4j JsonWebSignature superclass. Embedding is what gives this type the
	// inherited surface (getEncodedPayload, setKey, verifySignature, getHeaders, ...).
	*jose.JWS

	// unprotected is the unprotected header map. Port of the private field of the same name.
	unprotected *jose.Object

	// decodedPayload stores a cached parsed payload as a map, when applicable. Port of the
	// private field of the same name. As upstream, no payload setter invalidates it: the payload
	// must be set (at parse time, or by Signature.jadesSignaturePayload for a detached
	// signature) before the first DecodedPayload call.
	decodedPayload *jose.Object

	// jwsJsonSerializationObject is the parent JWSJsonSerializationObject. Port of the private
	// field of the same name.
	jwsJsonSerializationObject *JWSJsonSerializationObject
}

// NewJWS creates an empty JsonWebSignature. Port of the default constructor JWS().
//
// The CheckCritOverride installed here is the checkCrit() override: upstream deliberately
// separates structure validation from the cryptographic check, leaving the 'crit' rules to
// JAdESBaselineRequirementsChecker. Java gets this through virtual dispatch from
// verifySignature(); Go dispatches statically, so without this hook the embedded base would keep
// running its own check and reject every JAdES signature whose 'crit' names a JAdES header.
func NewJWS() *JWS {
	jws := &JWS{JWS: jose.NewJWS()}
	jws.JWS.CheckCritOverride = func(*jose.JWS) error { return nil }
	return jws
}

// NewJWSFromCompactSerializationParts instantiates a JWSCompactSerialization object (RFC 7515)
// from an array of strings holding the header, the optional payload and the signature. Port of
// the constructor JWS(String[]), whose Objects.requireNonNull becomes a panic carrying the Java
// message and whose JoseException becomes an IllegalInputException.
func NewJWSFromCompactSerializationParts(parts []string) (*JWS, error) {
	if parts == nil {
		panic("Parts part cannot be null")
	}
	jws := NewJWS()
	if err := jws.SetCompactSerializationParts(parts); err != nil {
		return nil, exception.NewIllegalInputExceptionWithCause("Unable to instantiate a compact JWS", err)
	}
	return jws, nil
}

// EncodedHeader returns the base64url-encoded protected header. Port of the getEncodedHeader()
// override, which exists in Java only to widen the inherited protected method to public.
func (j *JWS) EncodedHeader() string {
	return j.JWS.EncodedHeader()
}

// SetPayloadOctets sets the payload binaries depending on the 'b64' header's value. Port of
// setPayloadOctets(byte[]).
//
// The branch mirrors JsonWebSignature.setCompactSerializationParts: with b64=false the bytes are
// the payload, and otherwise they are its base64url encoding. Note that the false branch goes
// through setPayloadBytes, which - as upstream - leaves any cached encoded payload untouched.
func (j *JWS) SetPayloadOctets(payload []byte) {
	// see JsonWebSignature.setCompactSerializationParts(parts)
	if j.IsRfc7797UnencodedPayload() {
		j.SetPayloadBytes(payload)
	} else {
		j.SetEncodedPayload(string(payload))
	}
}

// SignedPayload returns the payload string based on the 'b64' value in the protected header,
// i.e. the actual signed payload value. Port of getSignedPayload().
func (j *JWS) SignedPayload() string {
	if j.IsRfc7797UnencodedPayload() {
		return j.UnverifiedPayload()
	}
	return j.EncodedPayload()
}

// IsRfc7797UnencodedPayload checks whether the signature's payload is 'b64' unencoded (RFC 7797),
// i.e. whether 'b64' is present and set to false. Port of the isRfc7797UnencodedPayload()
// override, which widens the inherited protected method to public.
func (j *JWS) IsRfc7797UnencodedPayload() bool {
	return j.JWS.IsRfc7797UnencodedPayload()
}

// SignatureValue returns the SignatureValue bytes. Port of getSignatureValue().
func (j *JWS) SignatureValue() []byte {
	return j.Signature()
}

// SetSignature sets the SignatureValue bytes. Port of the setSignature(byte[]) override, which
// widens the inherited protected method to public.
func (j *JWS) SetSignature(signature []byte) {
	j.JWS.SetSignature(signature)
}

// SetProtected sets the protected header from its base64url-encoded form. Port of
// setProtected(String).
func (j *JWS) SetProtected(protectedBase64Url string) error {
	return j.SetEncodedHeader(protectedBase64Url)
}

// Unprotected gets the unprotected header map. Port of getUnprotected().
func (j *JWS) Unprotected() *jose.Object {
	if j == nil {
		return nil
	}
	return j.unprotected
}

// SetUnprotected sets the unprotected header map. Port of setUnprotected(Map).
func (j *JWS) SetUnprotected(unprotected *jose.Object) {
	j.unprotected = unprotected
}

// DecodedPayload gets the payload parsed as a map, when applicable, caching the result. Port of
// getDecodedPayload().
//
// Despite the name it is the UNVERIFIED payload that is parsed - upstream calls
// getUnverifiedPayload() - so a caller must not read a signature check into it.
func (j *JWS) DecodedPayload() (*jose.Object, error) {
	if j.decodedPayload == nil {
		unverifiedPayload := j.UnverifiedPayload()
		if unverifiedPayload != "" {
			decoded, err := DSSJsonUtilsParseJSONStringToMap(unverifiedPayload)
			if err != nil {
				return nil, err
			}
			j.decodedPayload = decoded
		}
	}
	return j.decodedPayload, nil
}

// JwsJsonSerializationObject gets the JWSJsonSerializationObject. Port of
// getJwsJsonSerializationObject().
func (j *JWS) JwsJsonSerializationObject() *JWSJsonSerializationObject {
	return j.jwsJsonSerializationObject
}

// SetJwsJsonSerializationObject sets the JWSJsonSerializationObject. Port of
// setJwsJsonSerializationObject(JWSJsonSerializationObject).
func (j *JWS) SetJwsJsonSerializationObject(jwsJsonSerializationObject *JWSJsonSerializationObject) {
	j.jwsJsonSerializationObject = jwsJsonSerializationObject
}

// JwsSerializationType gets the signature's serialization type (compact, flattened, ...). Port of
// getJwsSerializationType().
func (j *JWS) JwsSerializationType() enumerations.JWSSerializationType {
	if j.jwsJsonSerializationObject != nil {
		return j.jwsJsonSerializationObject.JWSSerializationType()
	}
	return enumerations.JWSSerializationTypeCompactSerialization
}

// SetKnownCriticalHeaders sets the values of the 'crit' header that must be known and processed.
// Port of setKnownCriticalHeaders(Collection).
func (j *JWS) SetKnownCriticalHeaders(knownCriticalHeaders []string) {
	j.JWS.SetKnownCriticalHeaders(knownCriticalHeaders)
}

// CheckCrit separates structure validation from the cryptographic check and therefore does
// nothing. Port of the protected checkCrit() override; see
// eu.europa.esig.dss.jades.validation.JAdESBaselineRequirementsChecker for where the work
// actually happens.
//
// NewJWS also installs this as the embedded base's CheckCritOverride, which is what makes the
// override take effect on the base's own self-call.
func (j *JWS) CheckCrit() error {
	// separate structure validation and cryptographic check
	// (see eu.europa.esig.dss.jades.validation.JAdESBaselineRequirementsChecker)
	return nil
}

// ProtectedHeaderValueAsString returns the protected header value under key as a String, or "".
// Port of getProtectedHeaderValueAsString(String).
func (j *JWS) ProtectedHeaderValueAsString(key string) string {
	return DSSJsonUtilsToStringWithHeaderName(j.Headers().Value(key), key)
}

// ProtectedHeaderValueAsNumber returns the protected header value under key as a Number, or nil.
// Port of getProtectedHeaderValueAsNumber(String).
func (j *JWS) ProtectedHeaderValueAsNumber(key string) *jose.Number {
	return DSSJsonUtilsToNumberWithHeaderName(j.Headers().Value(key), key)
}

// ProtectedHeaderValueAsMap returns the protected header value under key as a JSON object, or an
// empty one. Port of getProtectedHeaderValueAsMap(String).
func (j *JWS) ProtectedHeaderValueAsMap(key string) *jose.Object {
	return DSSJsonUtilsToMap(j.Headers().Value(key), key)
}

// ProtectedHeaderValueAsList returns the protected header value under key as a JSON array, or an
// empty one. Port of getProtectedHeaderValueAsList(String).
func (j *JWS) ProtectedHeaderValueAsList(key string) []any {
	return DSSJsonUtilsToList(j.Headers().Value(key), key)
}
