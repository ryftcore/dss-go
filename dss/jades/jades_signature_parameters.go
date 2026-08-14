// Ported from dss-jades/src/main/java/eu/europa/esig/dss/jades/JAdESSignatureParameters.java (DSS 6.5.RC1).
package jades

import (
	"fmt"
	"time"

	"github.com/utain/esig/dss/document"
	"github.com/utain/esig/dss/enumerations"
)

// JAdESSignatureParameters holds the parameters to create/extend a JAdES signature.
type JAdESSignatureParameters struct {
	document.AbstractSignatureParameters[*JAdESTimestampParameters]

	// includeCertificateChain defines if certificate chain binaries must be included into the
	// signed header ('x5c' attribute).
	//
	// DEFAULT: TRUE (the certificate chain header will be included into the signed header)
	includeCertificateChain bool

	// includeSignatureType defines if the signature must incorporate its MimeType definition in
	// the signed header ('typ' attribute).
	//
	// DEFAULT: TRUE (the signature MimeType will be included into the signed header)
	includeSignatureType bool

	// signatureType defines a MimeType of the signature to be created, to be provided within a
	// signed header ('typ' attribute).
	//
	// DEFAULT: The type is determined based on the JWS serialization type
	signatureType string

	// includeKeyIdentifier defines whether a 'kid' (key identifier) header parameter should be
	// added to a protected header.
	//
	// NOTE: a signing certificate shall be provided to embed the 'kid' header
	//
	// DEFAULT: TRUE ('kid' header parameter is included into the signed header, provided that
	// the signing-certificate is defined within the signature parameters).
	includeKeyIdentifier bool

	// keyIdentifier is the value of the 'kid' (key identifier) parameter to be embedded within
	// the protected header of the signature.
	//
	// DEFAULT: when not defined and includeKeyIdentifier is enabled, the value of the embedded
	// 'kid' protected header corresponds to the IssuerSerial of the signing-certificate.
	keyIdentifier string

	// contentType defines the value of the 'cty' (content type) header parameter. When set, the
	// value of 'cty' header parameter will be defined with a given String. If not set, the value
	// of 'cty' header parameter will be derived from a MimeType of the signer document (except
	// for detached packaging).
	contentType string

	// x509Url defines a value for the 'x5u' signed header parameter. The value shall refer to a
	// URI where the X.509 public key certificate or certificate chain corresponding to the key
	// used to digitally sign the JWS can be retrieved from.
	//
	// NOTE: use SetSigningCertificate and SetIncludeCertificateChain to disable encapsulation of
	// the signing certificate and certificate chain binaries.
	//
	// DEFAULT: NULL (the 'x5u' header parameter is not included)
	x509Url string

	// base64UrlEncodedPayload defines if the payload has to be base64url encoded. If false,
	// original signed document binaries will be used according to RFC 7797.
	//
	// NOTE: the parameter is independent from base64UrlEncodedEtsiUComponents.
	//
	// Default: TRUE (base64url encoded payload)
	base64UrlEncodedPayload bool

	// base64UrlEncodedEtsiUComponents defines if the items of the 'etsiU' unprotected headers
	// will be incorporated in their corresponding base64url encodings, if FALSE the components
	// will appear as clear JSON instances. The parameter is used for Serialization (or
	// Flattened) format only with an unprotected header. All the components of 'etsiU' header
	// shall appear in the same representation.
	//
	// NOTE: the parameter is independent from base64UrlEncodedPayload.
	//
	// Default: TRUE (base64url encoded etsiU components); nil is Java's null third state
	// (unset).
	base64UrlEncodedEtsiUComponents *bool

	// signingCertificateDigestMethod is the DigestAlgorithm used to create a reference to a
	// signing certificate, namely 'x5t#256' for SHA256 or 'x5t#o' for other algorithms.
	signingCertificateDigestMethod enumerations.DigestAlgorithm

	// jwsSerializationType defines a JWS signature type according to RFC 7515, 3. JSON Web
	// Signature (JWS) Overview.
	//
	// Default: JWSSerializationType.COMPACT_SERIALIZATION
	jwsSerializationType enumerations.JWSSerializationType

	// sigDMechanism defines a used 'sigD' mechanism for a detached signature.
	sigDMechanism enumerations.SigDMechanism

	// jadesSigningTimeType identifies a type of claimed signing time header to be used on JAdES
	// signature creation.
	jadesSigningTimeType JAdESSigningTimeType

	// expirationTime is the value of the 'exp' (expiration time) signed header parameter as per
	// ETSI TS 119 411-5. The value is used for a TLS Certificate Binding signature and contains
	// the expiry date of the binding. nil is Java's null (unset).
	expirationTime *time.Time
}

// NewJAdESSignatureParameters instantiates the object with default parameters. Port of the
// default constructor.
func NewJAdESSignatureParameters() *JAdESSignatureParameters {
	return &JAdESSignatureParameters{
		AbstractSignatureParameters:    document.NewAbstractSignatureParameters[*JAdESTimestampParameters](),
		includeCertificateChain:        true,
		includeSignatureType:           true,
		includeKeyIdentifier:           true,
		base64UrlEncodedPayload:        true,
		signingCertificateDigestMethod: enumerations.DigestAlgorithm_SHA512,
		jwsSerializationType:           enumerations.JWSSerializationType_COMPACT_SERIALIZATION,
		jadesSigningTimeType:           JAdESSigningTimeType_IAT,
	}
}

// SetSignatureLevel overrides AbstractSerializableSignatureParameters#setSignatureLevel,
// restricting the value to the JAdES form. Panics with the Java message when signatureLevel is
// empty or not a JAdES level (IllegalArgumentException upstream).
func (p *JAdESSignatureParameters) SetSignatureLevel(signatureLevel enumerations.SignatureLevel) {
	form, err := signatureLevel.SignatureForm()
	if signatureLevel == "" || err != nil || enumerations.SignatureForm_JAdES != form {
		panic("Only JAdES form is allowed !")
	}
	p.AbstractSignatureParameters.SetSignatureLevel(signatureLevel)
}

// GetContentTimestampParameters overrides AbstractSerializableSignatureParameters, lazily
// instantiating JAdESTimestampParameters. Port of #getContentTimestampParameters.
func (p *JAdESSignatureParameters) GetContentTimestampParameters() *JAdESTimestampParameters {
	if p.ContentTimestampParameters == nil {
		p.ContentTimestampParameters = NewJAdESTimestampParameters()
	}
	return p.ContentTimestampParameters
}

// GetSignatureTimestampParameters overrides AbstractSerializableSignatureParameters, lazily
// instantiating JAdESTimestampParameters. Port of #getSignatureTimestampParameters.
func (p *JAdESSignatureParameters) GetSignatureTimestampParameters() *JAdESTimestampParameters {
	if p.SignatureTimestampParameters == nil {
		p.SignatureTimestampParameters = NewJAdESTimestampParameters()
	}
	return p.SignatureTimestampParameters
}

// GetArchiveTimestampParameters overrides AbstractSerializableSignatureParameters, lazily
// instantiating JAdESTimestampParameters. Port of #getArchiveTimestampParameters.
func (p *JAdESSignatureParameters) GetArchiveTimestampParameters() *JAdESTimestampParameters {
	if p.ArchiveTimestampParameters == nil {
		p.ArchiveTimestampParameters = NewJAdESTimestampParameters()
	}
	return p.ArchiveTimestampParameters
}

// IsIncludeCertificateChain defines if complete certificate chain binaries must be included into
// the signed header ('x5c' attribute). Port of #isIncludeCertificateChain.
func (p *JAdESSignatureParameters) IsIncludeCertificateChain() bool {
	return p.includeCertificateChain
}

// SetIncludeCertificateChain sets if complete certificate chain binaries must be included into
// the signed header. Default: TRUE. Port of #setIncludeCertificateChain.
func (p *JAdESSignatureParameters) SetIncludeCertificateChain(includeCertificateChain bool) {
	p.includeCertificateChain = includeCertificateChain
}

// IsIncludeSignatureType defines if the signature MimeType string must be included into the
// signed header ('typ' attribute). Port of #isIncludeSignatureType.
func (p *JAdESSignatureParameters) IsIncludeSignatureType() bool {
	return p.includeSignatureType
}

// SetIncludeSignatureType sets if the signature MimeType string must be included into the
// signed header ('typ' attribute). Default: TRUE. Port of #setIncludeSignatureType.
func (p *JAdESSignatureParameters) SetIncludeSignatureType(includeSignatureType bool) {
	p.includeSignatureType = includeSignatureType
}

// SignatureType gets the MimeType of the signature, to be incorporated in the signed header
// ('typ' attribute). Port of #getSignatureType.
func (p *JAdESSignatureParameters) SignatureType() string {
	return p.signatureType
}

// SetSignatureType sets the MimeType of the signature to be incorporated within the signed
// header ('typ' attribute). Default: derived from the selected JWS serialization type. Port of
// #setSignatureType.
func (p *JAdESSignatureParameters) SetSignatureType(signatureType string) {
	p.signatureType = signatureType
}

// IsIncludeKeyIdentifier returns whether a 'kid' (key identifier) header parameter should be
// created. Port of #isIncludeKeyIdentifier.
func (p *JAdESSignatureParameters) IsIncludeKeyIdentifier() bool {
	return p.includeKeyIdentifier
}

// SetIncludeKeyIdentifier sets whether a 'kid' (key identifier) header parameter should be
// created within a protected header, provided that a signing-certificate is defined within the
// signature parameters. Default: TRUE. Port of #setIncludeKeyIdentifier.
func (p *JAdESSignatureParameters) SetIncludeKeyIdentifier(includeKeyIdentifier bool) {
	p.includeKeyIdentifier = includeKeyIdentifier
}

// KeyIdentifier gets the value of the 'kid' (key identifier) protected header parameter. Port of
// #getKeyIdentifier.
func (p *JAdESSignatureParameters) KeyIdentifier() string {
	return p.keyIdentifier
}

// SetKeyIdentifier sets the 'kid' value to be incorporated within the signature's protected
// header. Port of #setKeyIdentifier.
func (p *JAdESSignatureParameters) SetKeyIdentifier(keyIdentifier string) {
	p.keyIdentifier = keyIdentifier
}

// ContentType gets value of the 'cty' (content type) signed header parameter is to be included
// in a protected header of the signature. Port of #getContentType.
func (p *JAdESSignatureParameters) ContentType() string {
	return p.contentType
}

// SetContentType sets value of the 'cty' (content type) signed header parameter is to be
// included in a protected header of the signature. Port of #setContentType.
func (p *JAdESSignatureParameters) SetContentType(contentType string) {
	p.contentType = contentType
}

// X509Url returns the value of the 'x5u' header parameter if present. Port of #getX509Url.
func (p *JAdESSignatureParameters) X509Url() string {
	return p.x509Url
}

// SetX509Url sets the value for the 'x5u' signed header parameter. Port of #setX509Url.
func (p *JAdESSignatureParameters) SetX509Url(x509Url string) {
	p.x509Url = x509Url
}

// SigningCertificateDigestMethod sees SetSigningCertificateDigestMethod. Port of
// #getSigningCertificateDigestMethod.
func (p *JAdESSignatureParameters) SigningCertificateDigestMethod() enumerations.DigestAlgorithm {
	return p.signingCertificateDigestMethod
}

// SetSigningCertificateDigestMethod sets the digest method indicating the digest algorithm to be
// used to calculate the certificate digest to define a signing certificate ('x5t#256' for SHA256
// or 'x5t#o' for other algorithms). Default: SHA512 ('x5t#o' attribute will be created). Panics
// with the Java message when signingCertificateDigestMethod is empty (Objects.requireNonNull
// upstream). Port of #setSigningCertificateDigestMethod.
func (p *JAdESSignatureParameters) SetSigningCertificateDigestMethod(signingCertificateDigestMethod enumerations.DigestAlgorithm) {
	if signingCertificateDigestMethod == "" {
		panic("SigningCertificateDigestMethod cannot be null!")
	}
	p.signingCertificateDigestMethod = signingCertificateDigestMethod
}

// JwsSerializationType gets the JWSSerializationType. Port of #getJwsSerializationType.
func (p *JAdESSignatureParameters) JwsSerializationType() enumerations.JWSSerializationType {
	return p.jwsSerializationType
}

// SetJwsSerializationType sets the JWSSerializationType. Default:
// JWSSerializationType.COMPACT_SERIALIZATION. Panics with the Java message when
// jwsSerializationType is empty (Objects.requireNonNull upstream). Port of
// #setJwsSerializationType.
func (p *JAdESSignatureParameters) SetJwsSerializationType(jwsSerializationType enumerations.JWSSerializationType) {
	if jwsSerializationType == "" {
		panic("JWSSerializationType cannot be null!")
	}
	p.jwsSerializationType = jwsSerializationType
}

// SigDMechanism returns a SigDMechanism to use. Port of #getSigDMechanism.
func (p *JAdESSignatureParameters) SigDMechanism() enumerations.SigDMechanism {
	return p.sigDMechanism
}

// SetSigDMechanism sets SigDMechanism to use for a Detached signature. Port of
// #setSigDMechanism.
func (p *JAdESSignatureParameters) SetSigDMechanism(sigDMechanism enumerations.SigDMechanism) {
	p.sigDMechanism = sigDMechanism
}

// JadesSigningTimeType returns the JAdES claimed signing-time header parameters to be used.
// Port of #getJadesSigningTimeType.
func (p *JAdESSignatureParameters) JadesSigningTimeType() JAdESSigningTimeType {
	return p.jadesSigningTimeType
}

// SetJadesSigningTimeType sets the claimed signing-time header parameters to be used.
//
// Requirements ETSI TS 119 182-1, clause 6.3, for iat and sigT: Before 2025-07-15T00:00:00Z the
// generator should include the iat header parameter for indicating the claimed signing time in
// new JAdES signatures and should not include the iat header parameter for indicating the
// claimed signing time in new JAdES signatures. Starting at 2025-07-15T00:00:00Z the generator
// shall include the iat header parameter for indicating the claimed signing time in new JAdES
// signatures.
//
// Default: IAT ('iat' header parameter will be used). Port of #setJadesSigningTimeType.
func (p *JAdESSignatureParameters) SetJadesSigningTimeType(jadesSigningTimeType JAdESSigningTimeType) {
	p.jadesSigningTimeType = jadesSigningTimeType
}

// ExpirationTime gets the expiration time of the signature. NOTE: The signed header is used for
// an ETSI TS 119 411-5 TLS Certificate Binding signature and contains an expiry date of the
// binding. nil is Java's null (unset). Port of #getExpirationTime.
func (p *JAdESSignatureParameters) ExpirationTime() *time.Time {
	return p.expirationTime
}

// SetExpirationTime sets the value for the 'exp' (expiration time) signed header of the
// signature. The claim identifies the expiration time on or after which the signature should
// not be accepted for processing. NOTE: The signed header is used for an ETSI TS 119 411-5 TLS
// Certificate Binding signature and contains an expiry date of the binding. Port of
// #setExpirationTime.
func (p *JAdESSignatureParameters) SetExpirationTime(expirationTime *time.Time) {
	p.expirationTime = expirationTime
}

// IsBase64UrlEncodedPayload gets if base64Url encoded payload shall be used. Port of
// #isBase64UrlEncodedPayload.
func (p *JAdESSignatureParameters) IsBase64UrlEncodedPayload() bool {
	return p.base64UrlEncodedPayload
}

// SetBase64UrlEncodedPayload sets if base64Url encoded payload shall be used. If FALSE, the
// unencoded (original) payload will be used according to RFC 7797.
//
// NOTE: some restrictions for payload content can apply when dealing with unencoded payload. For
// more information please see RFC 7797. The parameter is independent from
// base64UrlEncodedEtsiUComponents.
//
// Default: TRUE (base64Url encoded payload will be used). Port of
// #setBase64UrlEncodedPayload.
func (p *JAdESSignatureParameters) SetBase64UrlEncodedPayload(base64EncodedPayload bool) {
	p.base64UrlEncodedPayload = base64EncodedPayload
}

// IsBase64UrlEncodedEtsiUComponents gets if the instances of the 'etsiU' unprotected header
// shall appear in their corresponding base64url encoding.
//
// Default: TRUE (base64Url encoded etsiU components will be used). Port of
// #isBase64UrlEncodedEtsiUComponents; the *bool return mirrors Java's nullable Boolean (nil is
// Java's null, unset).
func (p *JAdESSignatureParameters) IsBase64UrlEncodedEtsiUComponents() *bool {
	return p.base64UrlEncodedEtsiUComponents
}

// SetBase64UrlEncodedEtsiUComponents sets if the instances of the 'etsiU' header shall appear in
// their corresponding base64url encoding. If FALSE the components of 'etsiU' will appear in
// their clear JSON incorporation. The parameter is used for Serialization (or Flattened) format
// only with unsigned properties.
//
// NOTE: the parameter is independent from base64UrlEncodedPayload.
//
// Default: TRUE (base64url encoded etsiU components). Port of
// #setBase64UrlEncodedEtsiUComponents(boolean); the *bool parameter lets nil clear the value,
// matching the field's nullable-Boolean nature (Java's overload only ever assigns a boxed
// non-null boolean, but every call site in this port passes an explicit pointer).
func (p *JAdESSignatureParameters) SetBase64UrlEncodedEtsiUComponents(base64UrlEncodedEtsiUComponents *bool) {
	p.base64UrlEncodedEtsiUComponents = base64UrlEncodedEtsiUComponents
}

// String ports #toString.
func (p *JAdESSignatureParameters) String() string {
	return fmt.Sprintf("JAdESSignatureParameters [includeCertificateChain=%v, includeSignatureType=%v, includeKeyIdentifier=%v, x509Url='%v', base64UrlEncodedPayload=%v, base64UrlEncodedEtsiUComponents=%v, signingCertificateDigestMethod=%v, jwsSerializationType=%v, sigDMechanism=%v, jadesSigningTimeType=%v] %s",
		p.includeCertificateChain, p.includeSignatureType, p.includeKeyIdentifier, p.x509Url,
		p.base64UrlEncodedPayload, boolPtrString(p.base64UrlEncodedEtsiUComponents),
		p.signingCertificateDigestMethod, p.jwsSerializationType, p.sigDMechanism, p.jadesSigningTimeType,
		p.AbstractSignatureParameters.String())
}

// boolPtrString renders a *bool the way Java's Boolean#toString would: "null" for the unset
// nullable value, "true"/"false" otherwise.
func boolPtrString(b *bool) string {
	if b == nil {
		return "null"
	}
	if *b {
		return "true"
	}
	return "false"
}

// Equals ports #equals.
func (p *JAdESSignatureParameters) Equals(other *JAdESSignatureParameters) bool {
	if p == other {
		return true
	}
	if other == nil {
		return false
	}
	if !p.AbstractSignatureParameters.Equals(&other.AbstractSignatureParameters) {
		return false
	}
	return p.includeCertificateChain == other.includeCertificateChain &&
		p.includeSignatureType == other.includeSignatureType &&
		p.includeKeyIdentifier == other.includeKeyIdentifier &&
		p.base64UrlEncodedPayload == other.base64UrlEncodedPayload &&
		boolPtrEquals(p.base64UrlEncodedEtsiUComponents, other.base64UrlEncodedEtsiUComponents) &&
		p.x509Url == other.x509Url &&
		p.signingCertificateDigestMethod == other.signingCertificateDigestMethod &&
		p.jwsSerializationType == other.jwsSerializationType &&
		p.sigDMechanism == other.sigDMechanism &&
		p.jadesSigningTimeType == other.jadesSigningTimeType
}

// boolPtrEquals ports the Objects.equals(Boolean, Boolean) semantics used by #equals for the
// nullable base64UrlEncodedEtsiUComponents field.
func boolPtrEquals(a, b *bool) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
}
