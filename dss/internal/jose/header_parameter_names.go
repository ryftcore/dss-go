// Ported from org.jose4j.jwx.HeaderParameterNames (jose4j 0.9.6).
package jose

// The registered JOSE header parameter names (RFC 7515 section 4.1, RFC 7518 sections 4.6-4.8
// and RFC 7797). Values are copied verbatim from jose4j; PORTING.md forbids inventing,
// abbreviating or "fixing" any of them, and X509CertificateSHA256Thumbprint in particular really
// does contain a '#'.
const (
	// HeaderAlgorithm is the "alg" (algorithm) header parameter.
	HeaderAlgorithm = "alg"
	// HeaderEncryptionMethod is the "enc" (encryption algorithm) header parameter.
	HeaderEncryptionMethod = "enc"
	// HeaderKeyID is the "kid" (key ID) header parameter.
	HeaderKeyID = "kid"
	// HeaderType is the "typ" (type) header parameter.
	HeaderType = "typ"
	// HeaderContentType is the "cty" (content type) header parameter.
	HeaderContentType = "cty"
	// HeaderJWKSetURL is the "jku" (JWK Set URL) header parameter.
	HeaderJWKSetURL = "jku"
	// HeaderJWK is the "jwk" (JSON Web Key) header parameter.
	HeaderJWK = "jwk"
	// HeaderX509CertificateChain is the "x5c" (X.509 certificate chain) header parameter.
	HeaderX509CertificateChain = "x5c"
	// HeaderX509CertificateThumbprint is the "x5t" (X.509 certificate SHA-1 thumbprint) header
	// parameter.
	HeaderX509CertificateThumbprint = "x5t"
	// HeaderX509CertificateSHA256Thumbprint is the "x5t#S256" (X.509 certificate SHA-256
	// thumbprint) header parameter.
	HeaderX509CertificateSHA256Thumbprint = "x5t#S256"
	// HeaderX509URL is the "x5u" (X.509 URL) header parameter.
	HeaderX509URL = "x5u"
	// HeaderEphemeralPublicKey is the "epk" (ephemeral public key) header parameter.
	HeaderEphemeralPublicKey = "epk"
	// HeaderAgreementPartyUInfo is the "apu" (agreement PartyUInfo) header parameter.
	HeaderAgreementPartyUInfo = "apu"
	// HeaderAgreementPartyVInfo is the "apv" (agreement PartyVInfo) header parameter.
	HeaderAgreementPartyVInfo = "apv"
	// HeaderZip is the "zip" (compression algorithm) header parameter.
	HeaderZip = "zip"
	// HeaderPBES2SaltInput is the "p2s" (PBES2 salt input) header parameter.
	HeaderPBES2SaltInput = "p2s"
	// HeaderPBES2IterationCount is the "p2c" (PBES2 count) header parameter.
	HeaderPBES2IterationCount = "p2c"
	// HeaderInitializationVector is the "iv" (initialization vector) header parameter.
	HeaderInitializationVector = "iv"
	// HeaderAuthenticationTag is the "tag" (authentication tag) header parameter.
	HeaderAuthenticationTag = "tag"
	// HeaderCritical is the "crit" (critical) header parameter.
	HeaderCritical = "crit"
	// HeaderBase64URLEncodePayload is the "b64" (base64url-encode payload) header parameter of
	// RFC 7797: when present and false, the payload appears in the JWS and in the signing input
	// as itself, with no encoding performed.
	HeaderBase64URLEncodePayload = "b64"
)
