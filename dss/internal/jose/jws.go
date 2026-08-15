// Ported from org.jose4j.jwx.JsonWebStructure and org.jose4j.jws.JsonWebSignature, flattened
// into one type because dss-jades only ever instantiates the JWS half (jose4j 0.9.6).
package jose

import (
	"bytes"
	"crypto"
	"errors"
	"fmt"
)

// CompactSerializationParts is JsonWebSignature.COMPACT_SERIALIZATION_PARTS: a JWS compact
// serialization has exactly three period-separated parts.
const CompactSerializationParts = 3

// ErrEmptyPart reports a required compact-serialization part that was empty. Port of
// checkNotEmptyPart's JoseException("The Encoded Header cannot be empty.").
var ErrEmptyPart = errors.New("jose: a required compact serialization part is empty")

// JWS is a JSON Web Signature (RFC 7515): a header, a payload and a signature, plus the rules
// that turn them into signing input and back. Port of JsonWebSignature together with the parts
// of its abstract parent JsonWebStructure that dss-jades reaches.
//
// Deliberately absent, because dss-jades has no call site for them: JWE, the compression
// algorithms, the signing half of the algorithm factory (DSS signs through its own
// SignatureValue plumbing and only ever hands jose4j a finished signature), the JWT consumer and
// the key resolvers. Their absence is a decision, not an omission.
//
// The zero value is not usable; construct with NewJWS.
type JWS struct {
	headers *Headers

	payloadBytes    []byte
	encodedPayload  string
	hasEncodedPay   bool
	signature       []byte
	key             crypto.PublicKey
	doKeyValidation bool

	validSignature    *bool
	knownCriticalHdrs map[string]struct{}

	// CheckCritOverride replaces the default 'crit' handling when non-nil.
	//
	// It exists because eu.europa.esig.dss.jades.validation.JWS overrides checkCrit() to do
	// nothing at all ("separate structure validation and cryptographic check"), deferring the
	// work to JAdESBaselineRequirementsChecker - and verifySignature calls checkCrit on itself,
	// i.e. through Java's virtual dispatch. Go dispatches statically, so an embedded base type
	// would silently keep calling its own default and start rejecting every JAdES signature
	// whose 'crit' names a JAdES header. Routing the self-call through this field is what keeps
	// the override effective; see PORTING.md on virtual-dispatch parity.
	CheckCritOverride func(j *JWS) error
}

// NewJWS returns an empty JWS. Port of the JsonWebSignature() constructor, whose only body sets
// the DISALLOW_NONE algorithm constraint - reproduced in resolveAlgorithm rather than as state,
// since DSS never relaxes it.
func NewJWS() *JWS {
	return &JWS{
		headers:         NewHeaders(),
		doKeyValidation: true,
	}
}

// Headers returns the protected header. Port of getHeaders().
func (j *JWS) Headers() *Headers { return j.headers }

// SetHeader sets a protected header parameter. Port of setHeader(String, Object).
func (j *JWS) SetHeader(name string, value any) { j.headers.Put(name, value) }

// Header returns a protected header parameter as a string. Port of getHeader(String).
func (j *JWS) Header(name string) string { return j.headers.StringValue(name) }

// ObjectHeader returns a protected header parameter as it stands. Port of getObjectHeader(String).
func (j *JWS) ObjectHeader(name string) any { return j.headers.Value(name) }

// AlgorithmHeaderValue returns the "alg" header. Port of getAlgorithmHeaderValue().
func (j *JWS) AlgorithmHeaderValue() string { return j.Header(HeaderAlgorithm) }

// SetAlgorithmHeaderValue sets the "alg" header. Port of setAlgorithmHeaderValue(String).
func (j *JWS) SetAlgorithmHeaderValue(alg string) { j.SetHeader(HeaderAlgorithm, alg) }

// ContentTypeHeaderValue returns the "cty" header. Port of getContentTypeHeaderValue().
func (j *JWS) ContentTypeHeaderValue() string { return j.Header(HeaderContentType) }

// KeyIDHeaderValue returns the "kid" header. Port of getKeyIdHeaderValue().
func (j *JWS) KeyIDHeaderValue() string { return j.Header(HeaderKeyID) }

// EncodedHeader returns BASE64URL(UTF8(protected header)). Port of getEncodedHeader().
func (j *JWS) EncodedHeader() string { return j.headers.EncodedHeader() }

// SetEncodedHeader replaces the protected header from its base64url form. Port of
// setEncodedHeader(String), including the "cannot be empty" precondition.
func (j *JWS) SetEncodedHeader(encodedHeader string) error {
	if encodedHeader == "" {
		return fmt.Errorf("%w: Encoded Header", ErrEmptyPart)
	}
	return j.headers.SetEncodedHeader(encodedHeader)
}

// SetPayload sets the payload from a string, encoding it as UTF-8. Port of setPayload(String);
// like jose4j it also drops any cached base64url payload.
func (j *JWS) SetPayload(payload string) {
	j.payloadBytes = []byte(payload)
	j.encodedPayload = ""
	j.hasEncodedPay = false
}

// SetPayloadBytes sets the payload from raw bytes. Port of setPayloadBytes(byte[]).
//
// Note what it does NOT do: unlike setPayload it leaves any previously cached encoded payload
// in place, so EncodedPayload can afterwards return a value unrelated to these bytes. That is
// jose4j's behaviour and JWSJsonSerializationParser depends on the ordering it implies (it sets
// the protected header first and the payload bytes last, with no encoded payload ever set), so
// it is reproduced rather than tidied up.
func (j *JWS) SetPayloadBytes(payloadBytes []byte) { j.payloadBytes = payloadBytes }

// UnverifiedPayloadBytes returns the payload without checking the signature. Port of
// getUnverifiedPayloadBytes().
func (j *JWS) UnverifiedPayloadBytes() []byte { return j.payloadBytes }

// UnverifiedPayload returns the payload as a UTF-8 string without checking the signature. Port
// of getUnverifiedPayload().
func (j *JWS) UnverifiedPayload() string { return UTF8String(j.payloadBytes) }

// SetEncodedPayload sets the payload from its base64url form, caching that form. Port of
// setEncodedPayload(String).
func (j *JWS) SetEncodedPayload(encodedPayload string) {
	j.encodedPayload = encodedPayload
	j.hasEncodedPay = true
	j.payloadBytes = Base64URLDecode(encodedPayload)
}

// EncodedPayload returns the cached base64url payload if there is one, and otherwise encodes the
// payload bytes. Port of getEncodedPayload().
func (j *JWS) EncodedPayload() string {
	if j.hasEncodedPay {
		return j.encodedPayload
	}
	return Base64URLEncode(j.payloadBytes)
}

// Signature returns the raw signature value. Port of the protected getSignature(), which reads
// JsonWebStructure's "integrity" field.
func (j *JWS) Signature() []byte { return j.signature }

// SetSignature sets the raw signature value. Port of the protected setSignature(byte[]).
func (j *JWS) SetSignature(signature []byte) {
	j.signature = signature
	j.validSignature = nil
}

// EncodedSignature returns BASE64URL(signature). Port of getEncodedSignature().
func (j *JWS) EncodedSignature() string { return Base64URLEncode(j.signature) }

// IsRfc7797UnencodedPayload reports whether the "b64" header is present and false, i.e. whether
// the payload appears unencoded in the signing input. Port of isRfc7797UnencodedPayload().
//
// The test is deliberately narrow, exactly as upstream writes it: the value must be a JSON
// boolean and it must be false. A string "false" does not count.
func (j *JWS) IsRfc7797UnencodedPayload() bool {
	b64, ok := j.headers.Value(HeaderBase64URLEncodePayload).(bool)
	return ok && !b64
}

// SetCompactSerializationParts loads a JWS from its three compact parts. Port of
// setCompactSerializationParts(String[]).
//
// The header is set first on purpose: it carries "b64", and the second part is interpreted as
// the payload itself or as its base64url encoding depending on that value.
func (j *JWS) SetCompactSerializationParts(parts []string) error {
	if len(parts) != CompactSerializationParts {
		return fmt.Errorf("jose: a JWS Compact Serialization must have exactly %d parts separated by period ('.') characters",
			CompactSerializationParts)
	}
	if err := j.SetEncodedHeader(parts[0]); err != nil {
		return err
	}
	if j.IsRfc7797UnencodedPayload() {
		j.SetPayload(parts[1])
	} else {
		j.SetEncodedPayload(parts[1])
	}
	j.SetSignature(Base64URLDecode(parts[2]))
	return nil
}

// SigningInputBytes returns the bytes the signature is computed over. Port of the private
// getSigningInputBytes(), i.e. RFC 7797 section 3:
//
//	b64 true   ASCII(BASE64URL(UTF8(protected header)) || '.' || BASE64URL(payload))
//	b64 false  ASCII(BASE64URL(UTF8(protected header)) || '.') || payload
//
// In the second form the payload is appended as raw bytes and never routed through a string,
// which is the difference between validating a detached binary document and corrupting it.
func (j *JWS) SigningInputBytes() []byte {
	if !j.IsRfc7797UnencodedPayload() {
		return ASCIIBytes(CompactSerialize(j.EncodedHeader(), j.EncodedPayload()))
	}
	var buf bytes.Buffer
	buf.Write(ASCIIBytes(j.EncodedHeader()))
	buf.WriteByte(0x2e) // ascii for "."
	buf.Write(j.payloadBytes)
	return buf.Bytes()
}

// Key returns the verification key. Port of getKey().
func (j *JWS) Key() crypto.PublicKey { return j.key }

// SetKey sets the verification key, discarding any cached verification result. Port of
// setKey(Key) together with the onNewKey() hook JsonWebSignature overrides to clear
// validSignature.
func (j *JWS) SetKey(key crypto.PublicKey) {
	j.key = key
	j.validSignature = nil
}

// IsDoKeyValidation reports whether the key is checked against the algorithm's requirements
// before use. Port of isDoKeyValidation().
func (j *JWS) IsDoKeyValidation() bool { return j.doKeyValidation }

// SetDoKeyValidation turns the key checks on or off. Port of setDoKeyValidation(boolean).
// JAdESSignature turns them off before validating ("restrict on key size,..."), so a signature
// made with, say, a 1024-bit RSA key is still reported as cryptographically intact and is left
// for the policy layer to reject.
func (j *JWS) SetDoKeyValidation(doKeyValidation bool) { j.doKeyValidation = doKeyValidation }

// SetKnownCriticalHeaders declares the 'crit' values the caller undertakes to understand. Port
// of setKnownCriticalHeaders(String...).
func (j *JWS) SetKnownCriticalHeaders(headers []string) {
	j.knownCriticalHdrs = make(map[string]struct{}, len(headers))
	for _, h := range headers {
		j.knownCriticalHdrs[h] = struct{}{}
	}
}

// CheckCrit verifies that every value of the 'crit' header is one the caller understands. Port
// of the protected checkCrit(), with JsonWebSignature's isSupportedCriticalHeader folded in
// (only "b64" is supported out of the box).
//
// When CheckCritOverride is set it is called instead; that is how the DSS subclass disables the
// check without this file knowing anything about DSS.
func (j *JWS) CheckCrit() error {
	if j.CheckCritOverride != nil {
		return j.CheckCritOverride(j)
	}
	value := j.headers.Value(HeaderCritical)
	if value == nil {
		return nil
	}
	var criticalHeaders []string
	switch v := value.(type) {
	case []any:
		for _, item := range v {
			s, ok := item.(string)
			if !ok {
				return fmt.Errorf("jose: %s header value not an array of strings", HeaderCritical)
			}
			criticalHeaders = append(criticalHeaders, s)
		}
	case []string:
		criticalHeaders = v
	default:
		return fmt.Errorf("jose: %s header value not an array (%T)", HeaderCritical, value)
	}
	for _, name := range criticalHeaders {
		if _, known := j.knownCriticalHdrs[name]; known {
			continue
		}
		if name == HeaderBase64URLEncodePayload {
			continue
		}
		return fmt.Errorf("jose: unrecognized header '%s' marked as critical", name)
	}
	return nil
}

// VerifySignature checks the signature over the signing input with the configured key. Port of
// verifySignature().
//
// The (bool, error) split is jose4j's boolean/exception split and callers depend on it: false
// means the cryptography said no, an error means the algorithm, the key or the encoding made the
// question unanswerable. The result is cached, and SetKey or SetSignature clears the cache,
// exactly as onNewKey does upstream.
func (j *JWS) VerifySignature() (bool, error) {
	alg, err := j.resolveAlgorithm()
	if err != nil {
		return false, err
	}
	if j.doKeyValidation {
		if err := alg.validateVerificationKey(j.key); err != nil {
			return false, err
		}
	}
	if j.validSignature != nil {
		return *j.validSignature, nil
	}
	if err := j.CheckCrit(); err != nil {
		return false, err
	}
	valid, err := alg.verify(j.signature, j.key, j.SigningInputBytes())
	if err != nil {
		return false, err
	}
	j.validSignature = &valid
	return valid, nil
}

// resolveAlgorithm looks up the "alg" header. Port of the private getAlgorithm(boolean) with the
// constraint check left in: JsonWebSignature's constructor installs DISALLOW_NONE and DSS never
// replaces it, so "none" is rejected here rather than being carried around as configuration.
func (j *JWS) resolveAlgorithm() (jwsAlgorithm, error) {
	name := j.AlgorithmHeaderValue()
	if name == "" {
		return nil, fmt.Errorf("jose: signature algorithm header (%s) not set", HeaderAlgorithm)
	}
	if name == AlgorithmNone {
		return nil, fmt.Errorf("%w: %q is not permitted", ErrUnsupportedAlgorithm, AlgorithmNone)
	}
	return lookupJWSAlgorithm(name)
}

// JwkHeader returns the public key carried by the "jwk" header, or nil when there is none. Port
// of getJwkHeader(), which is Headers.getPublicJwkHeaderValue(JWK, null) - including its refusal
// to accept a JWK that carries private key material.
func (j *JWS) JwkHeader() (crypto.PublicKey, error) {
	value := j.headers.Value(HeaderJWK)
	if value == nil {
		return nil, nil
	}
	params, ok := value.(*Object)
	if !ok {
		return nil, fmt.Errorf("jose: %s header is not a JSON object", HeaderJWK)
	}
	key, hasPrivate, err := ParseJWKPublicKey(params)
	if err != nil {
		return nil, err
	}
	if hasPrivate {
		return nil, fmt.Errorf("jose: %s header contains a private key, which it most definitely should not", HeaderJWK)
	}
	return key, nil
}
